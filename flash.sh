#!/usr/bin/env bash
set -euo pipefail

uso() {
    cat <<'EOF'
Uso: ./flash.sh [--compile-only]

Prepara o ambiente, compila o keymap vendor do Dilemma V3 (Procyon) e
orienta a gravação via USB: primeiro a metade direita, depois a esquerda.
No Arch Linux, oferece instalar as dependências ausentes com pacman.

  --compile-only  Só gera o firmware .uf2; não grava no teclado.
  -h, --help      Mostra esta ajuda sem instalar nem compilar nada.
EOF
}

case "${1:-}" in
    '')               opcoes=(--both) ;;
    --compile-only)   opcoes=() ;;
    -h|--help)        uso; exit 0 ;;
    *)                uso >&2; exit 2 ;;
esac
if (( $# > 1 )); then uso >&2; exit 2; fi

falha() { printf 'Erro: %s\n' "$1" >&2; exit 1; }
[[ $EUID != 0 ]] || falha 'Execute como usuário normal. Só o pacman será chamado com sudo.'
if (( ${#opcoes[@]} )); then
    [[ -t 0 ]] || falha 'Para gravar, rode em um terminal interativo.'
fi

printf '\nPreparação · Dependências do Dilemma V3\n'
requisitos=(git:git qmk:qmk make:make gcc:gcc arm-none-eabi-gcc:arm-none-eabi-gcc)
if (( ${#opcoes[@]} )); then
    requisitos+=(lsblk:util-linux findmnt:util-linux udisksctl:udisks2)
fi

faltam=()
for requisito in "${requisitos[@]}"; do
    comando=${requisito%:*}
    pacote=${requisito#*:}
    if ! command -v "$comando" >/dev/null && [[ " ${faltam[*]} " != *" $pacote "* ]]; then
        faltam+=("$pacote")
    fi
done
if (( ${#faltam[@]} )); then
    printf 'Pacotes necessários no Arch: %s\n' "${faltam[*]}"
    command -v pacman >/dev/null && command -v sudo >/dev/null ||
        falha 'Instale os equivalentes com o gerenciador da sua distribuição e rode novamente.'
    printf 'O pacman vai mostrar os pacotes e pedir confirmação antes de instalar.\n'
    sudo pacman -S --needed -- "${faltam[@]}" || falha 'Instalação cancelada ou malsucedida.'
fi
for requisito in "${requisitos[@]}"; do
    comando=${requisito%:*}
    command -v "$comando" >/dev/null || falha "$comando não encontrado no PATH."
done

repo="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
exec "$repo/scripts/build_flash.sh" "${opcoes[@]}"
