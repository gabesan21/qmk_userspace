#!/usr/bin/env bash
#
# build_flash.sh — build e flash do firmware do QMK userspace (Bastard Keyboards).
#
# Caminho único e documentado de build/flash do projeto
# (contrato durável: pop/specs/build-flash-script.md). O script detecta
# problemas de ambiente e instrui a correção exata — ele NÃO configura nada
# sozinho: `qmk setup` e o clone do firmware são decisões do dono.
#
# Uso: scripts/build_flash.sh [opções]   (ver --help)

set -euo pipefail

# ---------------------------------------------------------------------------
# Padrões do projeto (sobrescrevíveis por argumento)
# ---------------------------------------------------------------------------
readonly KEYBOARD_PADRAO="bastardkb/dilemma/3x5_3_procyon"
readonly KEYMAP_PADRAO="vendor"
readonly SUBMODULO="modules/bastardkb"
readonly FIRMWARE_REPO="bastardkb/bastardkb-qmk"
readonly ESPERA_BOOTLOADER_SEG=60

KEYBOARD="$KEYBOARD_PADRAO"
KEYMAP="$KEYMAP_PADRAO"
FLASH=0

# ---------------------------------------------------------------------------
# Saída visual: cores só em terminal e respeitando NO_COLOR
# ---------------------------------------------------------------------------
if [[ -t 1 && -z "${NO_COLOR:-}" ]]; then
    readonly C_RESET=$'\033[0m' C_BOLD=$'\033[1m' C_AZUL=$'\033[34m'
    readonly C_VERDE=$'\033[32m' C_VERMELHO=$'\033[31m' C_AMARELO=$'\033[33m'
else
    readonly C_RESET='' C_BOLD='' C_AZUL='' C_VERDE='' C_VERMELHO='' C_AMARELO=''
fi
readonly S_OK="${C_VERDE}✔${C_RESET}" S_ERRO="${C_VERMELHO}✘${C_RESET}" S_INFO="${C_AZUL}➜${C_RESET}"

passo()   { printf '\n%s%s== %s ==%s\n' "$C_BOLD" "$C_AZUL" "$1" "$C_RESET"; }
ok()      { printf '%s %s\n' "$S_OK" "$1"; }
info()    { printf '%s %s\n' "$S_INFO" "$1"; }
erro()    { printf '%s %s\n' "$S_ERRO" "$1" >&2; }
corrija() { printf '   %sCorrija com:%s %s\n' "$C_AMARELO" "$C_RESET" "$1" >&2; }

uso() {
    cat <<EOF
Uso: scripts/build_flash.sh [opções]

Compila o firmware do QMK userspace e, somente se pedido, grava no teclado via UF2.

Opções:
  -k, --keyboard <alvo>   Teclado alvo (padrão: $KEYBOARD_PADRAO)
  -m, --keymap <nome>     Keymap       (padrão: $KEYMAP_PADRAO)
      --flash             Grava o firmware via UF2 após o compile (opt-in)
  -h, --help              Mostra esta ajuda

Exemplos:
  scripts/build_flash.sh                                # compile do target padrão
  scripts/build_flash.sh -k bastardkb/dilemma/4x6_4_procyon
  scripts/build_flash.sh --flash                        # compile + flash (teclado em bootloader)
EOF
}

# ---------------------------------------------------------------------------
# Argumentos
# ---------------------------------------------------------------------------
while [[ $# -gt 0 ]]; do
    case "$1" in
        -k|--keyboard)
            [[ $# -ge 2 ]] || { erro "A opção $1 exige um valor."; exit 2; }
            KEYBOARD="$2"; shift 2 ;;
        -m|--keymap)
            [[ $# -ge 2 ]] || { erro "A opção $1 exige um valor."; exit 2; }
            KEYMAP="$2"; shift 2 ;;
        --flash) FLASH=1; shift ;;
        -h|--help) uso; exit 0 ;;
        *) erro "Opção desconhecida: $1"; uso; exit 2 ;;
    esac
done

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$REPO_ROOT"

# ---------------------------------------------------------------------------
# Passo 1 — Preflight: cada falha imprime o comando corretivo exato
# ---------------------------------------------------------------------------
passo "Passo 1/3 — Verificação do ambiente (preflight)"

if ! command -v qmk >/dev/null 2>&1; then
    erro "QMK CLI não encontrado no PATH."
    corrija "sudo pacman -S qmk   (Arch)   ou   python3 -m pip install --user qmk"
    exit 1
fi
ok "qmk instalado: $(qmk --version)"

status_submodulo="$(git submodule status "$SUBMODULO" 2>/dev/null || true)"
if [[ -z "$status_submodulo" || "$status_submodulo" == -* || -z "$(ls -A "$SUBMODULO" 2>/dev/null)" ]]; then
    erro "Submódulo '$SUBMODULO' não inicializado."
    corrija "git submodule update --init $SUBMODULO"
    exit 1
fi
ok "Submódulo $SUBMODULO inicializado."

QMK_HOME="$(qmk config user.qmk_home 2>/dev/null | sed -n 's/^user\.qmk_home=//p' || true)"
if [[ -z "$QMK_HOME" || "$QMK_HOME" == "None" || ! -d "$QMK_HOME" ]]; then
    erro "user.qmk_home não configurado (ou o diretório não existe)."
    corrija "qmk setup $FIRMWARE_REPO   (clona o firmware e aponta user.qmk_home)"
    exit 1
fi
ok "user.qmk_home: $QMK_HOME"

# ---------------------------------------------------------------------------
# Passo 2 — Compile (QMK_USERSPACE aponta este repo como overlay do firmware)
# ---------------------------------------------------------------------------
passo "Passo 2/3 — Compile"
info "Keyboard: $KEYBOARD"
info "Keymap:   $KEYMAP"

if ! QMK_USERSPACE="$REPO_ROOT" qmk compile -kb "$KEYBOARD" -km "$KEYMAP"; then
    erro "Compile falhou. Se a causa for o ambiente, veja: qmk doctor"
    exit 3
fi

UF2="$QMK_HOME/${KEYBOARD//\//_}_${KEYMAP}.uf2"
if [[ ! -f "$UF2" ]]; then
    erro "Firmware esperado não encontrado: $UF2"
    info "Confira o nome do arquivo gerado na saída do compile acima."
    exit 3
fi
ok "Firmware gerado: $UF2"

# ---------------------------------------------------------------------------
# Passo 3 — Flash via UF2 (somente com --flash; nunca implícito)
# ---------------------------------------------------------------------------
passo "Passo 3/3 — Flash via UF2"

if [[ "$FLASH" -eq 0 ]]; then
    info "Flash não solicitado — nada foi gravado no teclado."
    info "Para gravar: $0 --flash   (com o teclado em modo bootloader)"
    exit 0
fi

cat <<EOF

  Coloque o Dilemma em modo bootloader UF2:
    1. pressione o botão RESET duas vezes em sequência rápida (double-tap), ou
    2. segure o botão BOOT enquanto conecta o cabo USB.
  Um volume chamado RPI-RP2 deve aparecer montado no sistema.

EOF

encontra_volume_uf2() {
    shopt -s nullglob
    local candidatos=( "/run/media/$USER"/RPI-RP2* "/media/$USER"/RPI-RP2* /mnt/RPI-RP2* )
    if [[ ${#candidatos[@]} -gt 0 ]]; then
        printf '%s\n' "${candidatos[0]}"
    fi
}

info "Aguardando o volume RPI-RP2 por até ${ESPERA_BOOTLOADER_SEG}s..."
volume=""
for ((i = 0; i < ESPERA_BOOTLOADER_SEG; i++)); do
    volume="$(encontra_volume_uf2)"
    if [[ -n "$volume" ]]; then break; fi
    sleep 1
done
if [[ -z "$volume" ]]; then
    erro "Volume RPI-RP2 não apareceu em ${ESPERA_BOOTLOADER_SEG}s."
    info "Confirme o modo bootloader e rode novamente com --flash."
    exit 4
fi
ok "Volume do bootloader: $volume"

info "Gravando $(basename "$UF2")..."
if ! cp "$UF2" "$volume/"; then
    erro "Falha ao copiar o firmware para $volume."
    exit 4
fi
ok "Firmware gravado — o teclado reinicia sozinho ao fim da cópia."
