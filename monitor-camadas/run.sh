#!/usr/bin/env bash
set -euo pipefail

falha() { printf 'Erro: %s\n' "$1" >&2; exit 1; }
confirmar() {
    local resposta
    [[ -t 0 ]] || return 1
    read -r -p "$1 [s/N] " resposta && [[ $resposta == [sS] ]]
}

go_compativel() {
    local versao
    versao=$(GOTOOLCHAIN=local go env GOVERSION 2>/dev/null) || return 1
    [[ $versao =~ ^go([0-9]+)\.([0-9]+) ]] &&
        (( BASH_REMATCH[1] > 1 || (BASH_REMATCH[1] == 1 && BASH_REMATCH[2] >= 25) ))
}

# O descritor QMK começa com Usage Page FF31, Usage 74 e Collection Application.
# Percorrer itens (não pesquisar bytes nos dados) e exigir uma única coleção:
# uma interface composta não pode receber acesso por causa de apenas um usage.
console_exclusivo() {
    od -An -v -tu1 -- "$1" | awk '
        { for (i = 1; i <= NF; i++) b[++n] = $i }
        END {
            if (n < 8 || b[1] != 6 || b[2] != 49 || b[3] != 255 ||
                b[4] != 9 || b[5] != 116 || b[6] != 161 || b[7] != 1) exit 1
            for (i = 8; i < n; ) {
                p = b[i]; tag = int(p / 4); tamanho = p % 4
                if (p == 254 || tag == 40 || tag == 48) exit 1
                i += 1 + (tamanho == 3 ? 4 : tamanho)
            }
            exit !(i == n && b[n] == 192)
        }'
}

preparar_acesso() {
    local raw dev digest regra
    for raw in /sys/class/hidraw/hidraw*; do
        dev=/dev/${raw##*/}
        [[ -c $dev && ! -r $dev ]] || continue
        grep -Fqx 'HID_ID=0003:0000A8F8:00001836' "$raw/device/uevent" 2>/dev/null || continue
        cat -- "$raw/device/report_descriptor" > "$tmp/descritor" 2>/dev/null || continue
        console_exclusivo "$tmp/descritor" || continue

        if ! command -v udevadm >/dev/null || ! command -v sudo >/dev/null; then
            printf 'Sem acesso a %s. A ativação da regra requer udevadm e sudo.\n' "$dev" >&2
            continue
        fi
        digest=$(sha256sum -- "$tmp/descritor")
        digest=${digest%% *}
        regra=/etc/udev/rules.d/70-dilemma-console-$digest.rules
        # O udev confere os mesmos bytes validados acima a cada conexão, não o índice hidraw.
        printf '%s\n' "SUBSYSTEM==\"hidraw\", ATTRS{idVendor}==\"a8f8\", ATTRS{idProduct}==\"1836\", PROGRAM==\"/usr/bin/sha256sum /sys%p/device/report_descriptor\", RESULT==\"$digest *\", TAG+=\"uaccess\"" > "$tmp/regra"
        if [[ -e $regra ]] && ! cmp -s -- "$regra" "$tmp/regra"; then
            printf 'A regra existente %s é diferente; não será sobrescrita. Consulte o administrador.\n' "$regra" >&2
            continue
        fi
        printf 'Sem leitura de %s. Regra da coleção console para a sessão local ativa:\n%s\n' "$dev" "$regra" >&2
        if ! confirmar 'Permitir instalar/ativar essa regra com sudo?'; then continue; fi
        if [[ ! -f $regra ]]; then
            sudo install -m 0644 -- "$tmp/regra" "$regra" || { printf 'Regra não instalada.\n' >&2; continue; }
        fi
        if ! sudo udevadm control --reload-rules || ! sudo udevadm trigger --settle --action=add "$raw" || [[ ! -r $dev ]]; then
            printf 'Reconecte o USB na sessão local ativa; se o acesso continuar negado, consulte o README.\n' >&2
        fi
    done
}

[[ $(uname -s) == Linux ]] || falha 'O monitor usa hidraw do Linux.'
(( EUID != 0 )) || falha 'Execute como usuário comum. Não inicie a TUI com sudo.'
cd -- "$(dirname -- "$(realpath -- "${BASH_SOURCE[0]}")")"
export GOWORK=off
if ! go_compativel; then
    aviso='Instale Go 1.25 ou superior e confira qual go aparece no PATH.'
    [[ -e /etc/arch-release ]] && command -v pacman >/dev/null && command -v sudo >/dev/null || falha "$aviso"
    confirmar 'Go 1.25+ ausente. Instalar/atualizar o pacote go com sudo pacman -S --needed go?' || falha "$aviso"
    sudo pacman -S --needed go
    hash -r
    go_compativel || falha "$aviso"
fi

cache=${XDG_CACHE_HOME:-$HOME/.cache}/monitor-camadas
mkdir -p -- "$cache"
printf 'Preparando dependências e compilando o monitor…\n' >&2
go mod download
go build -mod=readonly -o "$cache/monitor-camadas" .

tmp=$(mktemp -d)
trap 'rm -rf -- "$tmp"' EXIT
preparar_acesso
rm -rf -- "$tmp"
trap - EXIT
exec "$cache/monitor-camadas" "$@"
