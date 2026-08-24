#!/usr/bin/env bash
#
# build_flash.sh — build e flash do firmware do QMK userspace (Bastard Keyboards).
#
# Caminho único e documentado de build/flash do projeto (contrato durável:
# pop/specs/build-flash-script.md). O preflight detecta problemas de ambiente
# e executa ele mesmo as correções documentadas — init do submódulo e
# `qmk setup -b main` do fork — revalidando em seguida (fail-fast se a
# correção não resolver). A instalação do QMK CLI segue como instrução ao
# dono, e o flash é sempre opt-in via --flash.
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
readonly LARGURA_CAIXA=64

KEYBOARD="$KEYBOARD_PADRAO"
KEYMAP="$KEYMAP_PADRAO"
FLASH=0

# ---------------------------------------------------------------------------
# Saída visual: cores só em terminal e respeitando NO_COLOR
# ---------------------------------------------------------------------------
if [[ -t 1 && -z "${NO_COLOR:-}" ]]; then
    readonly C_RESET=$'\033[0m'    C_BOLD=$'\033[1m'  C_DIM=$'\033[2m'
    readonly C_AZUL=$'\033[34m'    C_CIANO=$'\033[36m'
    readonly C_VERDE=$'\033[32m'   C_VERMELHO=$'\033[31m'  C_AMARELO=$'\033[33m'
else
    readonly C_RESET='' C_BOLD='' C_DIM='' C_AZUL='' C_CIANO='' \
             C_VERDE='' C_VERMELHO='' C_AMARELO=''
fi
readonly S_OK="${C_VERDE}✔${C_RESET}" S_ERRO="${C_VERMELHO}✘${C_RESET}"
readonly S_INFO="${C_CIANO}➜${C_RESET}" S_FIX="${C_AMARELO}⚙${C_RESET}"
readonly SPINNER_FRAMES=(⠋ ⠙ ⠹ ⠸ ⠼ ⠴ ⠦ ⠧ ⠇ ⠏)

# ---------------------------------------------------------------------------
# Blocos visuais: caixas, réguas de seção, checklist e campos alinhados.
# ${#var} conta caracteres (não bytes), então o preenchimento fica correto
# mesmo com acentos e símbolos multibyte.
# ---------------------------------------------------------------------------

_borda() {  # $1/$2: cantos da caixa; $3: cor do bloco
    local preenchimento
    printf -v preenchimento '%*s' "$LARGURA_CAIXA" ''
    printf '%s%s%s%s%s\n' "$3" "$1" "${preenchimento// /─}" "$2" "$C_RESET"
}

_linha() {  # $1: cor da borda; $2: conteúdo (margem esquerda de 2 espaços)
    local espacos
    espacos=$(( LARGURA_CAIXA - 2 - ${#2} ))
    if (( espacos < 0 )); then espacos=0; fi
    printf '%s│%s  %s%*s%s│%s\n' "$1" "$C_RESET" "$2" "$espacos" '' "$1" "$C_RESET"
}

banner() {
    local cor="$C_BOLD$C_AZUL"
    _borda "╭" "╮" "$cor"
    _linha "$cor" ""
    _linha "$cor" "QMK Userspace — Build & Flash"
    _linha "$cor" "Bastard Keyboards · Dilemma · RP2040/UF2"
    _linha "$cor" ""
    _borda "╰" "╯" "$cor"
}

secao() {  # título de seção com régua preenchendo a largura da caixa
    local titulo="$1" preenchimento regra
    regra=$(( LARGURA_CAIXA + 2 - ${#titulo} - 4 ))
    if (( regra < 0 )); then regra=0; fi
    printf -v preenchimento '%*s' "$regra" ''
    printf '\n%s%s── %s %s%s\n' "$C_BOLD" "$C_CIANO" "$titulo" "${preenchimento// /─}" "$C_RESET"
}

caixa_bootloader() {
    local cor="$C_AMARELO"
    printf '\n'
    _borda "╭" "╮" "$cor"
    _linha "$cor" "Coloque o Dilemma em modo bootloader UF2:"
    _linha "$cor" ""
    _linha "$cor" "1. pressione RESET duas vezes em sequência rápida (double-tap);"
    _linha "$cor" "2. ou segure BOOT enquanto conecta o cabo USB."
    _linha "$cor" ""
    _linha "$cor" "Um volume chamado RPI-RP2 deve aparecer montado no sistema."
    _borda "╰" "╯" "$cor"
    printf '\n'
}

# ---------------------------------------------------------------------------
# Mensagens e linhas alinhadas
# ---------------------------------------------------------------------------

ok()      { printf '  %s %s\n' "$S_OK" "$1"; }
info()    { printf '  %s %s\n' "$S_INFO" "$1"; }
erro()    { printf '  %s %s\n' "$S_ERRO" "$1" >&2; }
corrija() { printf '     %sCorrija com:%s %s\n' "$C_AMARELO" "$C_RESET" "$1" >&2; }
autofix() { printf '     %s Corrigindo agora: %s%s%s\n' "$S_FIX" "$C_BOLD" "$1" "$C_RESET"; }

item() {  # linha de checklist do preflight: símbolo + rótulo em coluna + detalhe
    local simbolo="$1" rotulo="$2" detalhe="$3" espacos
    espacos=$(( 14 - ${#rotulo} ))
    if (( espacos < 1 )); then espacos=1; fi
    printf '  %s %s%s%s%*s%s\n' "$simbolo" "$C_BOLD" "$rotulo" "$C_RESET" "$espacos" '' "$detalhe"
}

campo() {  # rótulo em negrito + valor, alinhados na mesma coluna do checklist
    local rotulo="$1" valor="$2" espacos
    espacos=$(( 14 - ${#rotulo} ))
    if (( espacos < 1 )); then espacos=1; fi
    printf '     %s%s%s%*s%s\n' "$C_BOLD" "$rotulo" "$C_RESET" "$espacos" '' "$valor"
}

# ---------------------------------------------------------------------------
# Ajuda
# ---------------------------------------------------------------------------

uso() {
    banner
    secao "Uso"
    printf '  %sscripts/build_flash.sh%s [opções]\n' "$C_BOLD" "$C_RESET"
    secao "Descrição"
    cat <<EOF
  Compila o firmware do QMK userspace e, somente se pedido (--flash),
  grava no teclado via UF2. O preflight corrige sozinho o submódulo e o
  user.qmk_home (auto-fix documentado) e revalida antes de seguir;
  falhas restantes abortam com o comando corretivo exato.
EOF
    secao "Opções"
    cat <<EOF
  -k, --keyboard <alvo>   Teclado alvo (padrão: $KEYBOARD_PADRAO)
  -m, --keymap <nome>     Keymap       (padrão: $KEYMAP_PADRAO)
      --flash             Grava o firmware via UF2 após o compile (opt-in)
  -h, --help              Mostra esta ajuda
EOF
    secao "Códigos de saída"
    cat <<EOF
  0  sucesso                 3  falha no compile
  1  falha de preflight      4  falha no flash
  2  uso inválido
EOF
    secao "Exemplos"
    cat <<EOF
  scripts/build_flash.sh                              compile do target padrão
  scripts/build_flash.sh -k bastardkb/dilemma/4x6_4_procyon
  scripts/build_flash.sh --flash                      compile + flash (bootloader)
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

banner

# ---------------------------------------------------------------------------
# Passo 1 — Preflight: cada falha dispara o auto-fix documentado, revalida
# e só então aborta (fail-fast) se a correção não resolver.
# ---------------------------------------------------------------------------
secao "Passo 1/3 · Verificação do ambiente (preflight)"

# QMK CLI: sem auto-fix — instalar programas é decisão do dono
if ! command -v qmk >/dev/null 2>&1; then
    item "$S_ERRO" "QMK CLI" "não encontrado no PATH"
    corrija "sudo pacman -S qmk   (Arch)   ou   python3 -m pip install --user qmk"
    exit 1
fi
item "$S_OK" "QMK CLI" "instalado ($(qmk --version))"

submodulo_ok() {
    local status
    status="$(git submodule status "$SUBMODULO" 2>/dev/null || true)"
    [[ -n "$status" && "$status" != -* && -n "$(ls -A "$SUBMODULO" 2>/dev/null)" ]]
}

if submodulo_ok; then
    item "$S_OK" "Submódulo" "$SUBMODULO inicializado"
else
    item "$S_ERRO" "Submódulo" "$SUBMODULO não inicializado"
    autofix "git submodule update --init $SUBMODULO"
    if git submodule update --init "$SUBMODULO" && submodulo_ok; then
        item "$S_OK" "Submódulo" "$SUBMODULO inicializado (corrigido agora)"
    else
        erro "A correção automática do submódulo não resolveu."
        corrija "git submodule update --init $SUBMODULO   (rode manualmente e veja o erro)"
        exit 1
    fi
fi

qmk_home() {
    qmk config user.qmk_home 2>/dev/null | sed -n 's/^user\.qmk_home=//p'
}

QMK_HOME="$(qmk_home || true)"
if [[ -n "$QMK_HOME" && "$QMK_HOME" != "None" && -d "$QMK_HOME" ]]; then
    item "$S_OK" "Firmware" "user.qmk_home: $QMK_HOME"
else
    item "$S_ERRO" "Firmware" "user.qmk_home não configurado (ou o diretório não existe)"
    autofix "qmk setup -b main $FIRMWARE_REPO"
    if qmk setup -b main "$FIRMWARE_REPO"; then
        QMK_HOME="$(qmk_home || true)"
    fi
    if [[ -n "$QMK_HOME" && "$QMK_HOME" != "None" && -d "$QMK_HOME" ]]; then
        item "$S_OK" "Firmware" "user.qmk_home: $QMK_HOME (corrigido agora)"
    else
        erro "A correção automática do user.qmk_home não resolveu."
        corrija "qmk setup -b main $FIRMWARE_REPO   (rode manualmente e veja o erro)"
        exit 1
    fi
fi

# ---------------------------------------------------------------------------
# Passo 2 — Compile (QMK_USERSPACE aponta este repo como overlay do firmware)
# ---------------------------------------------------------------------------
secao "Passo 2/3 · Compile"
campo "Teclado" "$KEYBOARD"
campo "Keymap"  "$KEYMAP"
campo "Overlay" "$REPO_ROOT"
printf '\n  %s %sqmk compile -kb %s -km %s%s\n\n' \
    "$S_INFO" "$C_DIM" "$KEYBOARD" "$KEYMAP" "$C_RESET"

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
ok "Firmware gerado:"
printf '     %s%s%s\n' "$C_VERDE" "$UF2" "$C_RESET"

# ---------------------------------------------------------------------------
# Passo 3 — Flash via UF2 (somente com --flash; nunca implícito)
# ---------------------------------------------------------------------------
secao "Passo 3/3 · Flash via UF2"

resumo() {  # $1: desfecho do flash
    secao "Resumo"
    campo "Preflight" "ambiente verificado"
    campo "Compile"   "$(basename "$UF2")"
    campo "Flash"     "$1"
    campo "Duração"   "${SECONDS}s"
}

if [[ "$FLASH" -eq 0 ]]; then
    info "Flash não solicitado — nada foi gravado no teclado."
    info "Para gravar: $0 --flash   (com o teclado em modo bootloader)"
    resumo "não solicitado (opt-in via --flash)"
    exit 0
fi

caixa_bootloader

encontra_volume_uf2() {
    shopt -s nullglob
    local candidatos=( "/run/media/$USER"/RPI-RP2* "/media/$USER"/RPI-RP2* /mnt/RPI-RP2* )
    if [[ ${#candidatos[@]} -gt 0 ]]; then
        printf '%s\n' "${candidatos[0]}"
    fi
}

volume=""
for ((i = 0; i < ESPERA_BOOTLOADER_SEG; i++)); do
    volume="$(encontra_volume_uf2)"
    if [[ -n "$volume" ]]; then break; fi
    if [[ -t 1 ]]; then
        printf '\r  %s%s%s Aguardando o volume RPI-RP2 (%2ds/%ds)...' \
            "$C_CIANO" "${SPINNER_FRAMES[i % ${#SPINNER_FRAMES[@]}]}" "$C_RESET" \
            "$i" "$ESPERA_BOOTLOADER_SEG"
    elif (( i == 0 )); then
        info "Aguardando o volume RPI-RP2 por até ${ESPERA_BOOTLOADER_SEG}s..."
    fi
    sleep 1
done
if [[ -t 1 ]]; then printf '\r%70s\r' ''; fi

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

resumo "gravado via UF2 em $volume"
