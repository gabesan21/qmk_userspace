---
id: build-flash-script
project: QMK Userspace
domain: build
kind: contract
status: draft
implementation: implemented
origin: "1.2"
created: 2026-08-24
updated: 2026-11-21
supersedes: []
superseded_by:
---

# Spec — Build e flash local do Dilemma V3

## Contract

Na raiz do repositório, rode `./flash.sh` em um terminal para compilar o keymap `vendor` do Dilemma V3 (3×5+3 Procyon) uma vez e gravar a metade direita e depois a esquerda, com o mesmo UF2. Use `./flash.sh --compile-only` para gerar só o `.uf2`.

O guia prepara as dependências e delega a compilação e o flash a `scripts/build_flash.sh`. Na metade ligada por USB, segure a tecla inferior externa (`/` na direita ou `Z` na esquerda) para ativar Ponteiro e, sem soltar, pressione a tecla superior externa da mesma metade (`BOOT`). O atalho depende do firmware já gravado. Mantenha o TRRS conectado, use um cabo USB de dados e não retire o USB durante a escrita. Ao terminar, reconecte à direita e teste as teclas. Se o atalho não funcionar, consulte o [guia oficial de recuperação](https://docs.bastardkb.com/fw/flashing.html).

## Expected behavior

- O guia verifica as aplicações necessárias para o modo escolhido. No Arch, oferece instalar as ausentes via `sudo pacman -S --needed`; em outras distribuições, informa as dependências que precisam ser instaladas.
- O motor inicializa o submódulo se necessário, prepara/configura o checkout QMK no primeiro uso e exibe o caminho do `.uf2` após o build.
- No modo de gravação, aguarda até 60 s pelo volume `RPI-RP2` da metade indicada. Se necessário, tenta montá-lo com `udisksctl`. Copia o UF2 somente para um volume único, gravável e identificado como bootloader RP2040. Confere o desaparecimento do volume antes de concluir ou avançar para a outra metade.

## Invariants

- `./flash.sh` é o ponto de entrada interativo para gravar. `--compile-only` nunca procura nem grava volume USB. O motor, quando chamado diretamente, só grava com `--flash` ou `--both`.
- A preparação não apaga nem atualiza automaticamente um checkout de firmware existente. O script inteiro não deve ser executado como root; a instalação de pacotes, quando necessária, pede confirmação pelo `pacman`.
- Mensagens ao usuário em pt-BR; cores desligadas fora de terminal ou com `NO_COLOR`.
- Nenhum artefato de build (`*.uf2`, `*.hex`, `*.bin`) é gerado dentro do repo (o compile escreve em `user.qmk_home`).

## Interfaces

- **Guia na raiz:** `./flash.sh [--compile-only]` ou `./flash.sh --help`.
- **Motor:** `scripts/build_flash.sh [-k|--keyboard <alvo>] [-m|--keymap <nome>] [--flash | --both]`; os padrões são `bastardkb/dilemma/3x5_3_procyon` e `vendor`. `--both` é específico desse modelo. Sem argumentos posicionais.
- **Output:** passos numerados, caminho do UF2, instruções de bootloader e resumo final; erros em stderr.
- **Códigos de saída:** `0` sucesso · `1` falha de preflight · `2` uso inválido · `3` falha no compile · `4` falha no flash.
- **Ambiente:** Bash em Linux, QMK CLI e toolchain de compilação; para gravação, `lsblk`, `findmnt` e montagem de volume USB (`udisksctl` quando não há montagem automática).

## Errors and limits

- **Dependência ausente:** o guia tenta instalar no Arch após confirmação ou orienta a instalar manualmente; falha ou recusa interrompe antes do build. Ao chamar só o motor, QMK CLI ausente é um erro com instrução de instalação.
- **Submódulo não inicializado/vazio:** o script executa `git submodule update --init modules/bastardkb` e revalida; se não resolver, erro + comando corretivo; saída `1`.
- **`user.qmk_home` ausente ou inválido:** reutiliza um checkout válido em `~/qmk_firmware` ou chama `qmk setup -H ~/qmk_firmware -b main -y bastardkb/bastardkb-qmk`. Grava `user.qmk_home` explicitamente e revalida; saída `1` se o ambiente continuar inválido.
- **Compile falha ou `.uf2` não encontrado:** erro apontando a saída do compile/`qmk doctor`; saída `3`.
- **Volume ausente, ambíguo, inválido ou não montável; erro de escrita ou reinício não observado:** mensagem de diagnóstico e saída `4`, sem avançar para a outra metade.

## Conformance criteria

- [ ] `bash -n` e leitura confirmam bash válido, shebang e fail-fast (`set -euo pipefail`).
- [ ] O guia verifica dependências antes do build; o motor prepara o submódulo e o firmware.
- [ ] O padrão é Dilemma V3 + `vendor`; `--compile-only` não acessa USB.
- [ ] Flash exige um único volume RP2040; no guia da raiz, a direita deve sair do bootloader antes de iniciar a esquerda.
- [ ] Teste físico: o volume desaparece após gravar, o teclado volta a funcionar com USB à direita e as camadas atualizadas respondem.

## Out of scope

- Instalação automática de pacotes em distribuições sem `pacman`.
- Flash por outro mecanismo (DFU, avrdude etc.).
- Build remoto/CI — proibido por decisão do projeto (ver [[pop/PROJECT|PROJECT]]).

## Open questions

- Confirmar a montagem e o reinício com o Dilemma V3 físico nas duas metades.

## Related references

- [[pop/PROJECT|PROJECT]] — *target padrão, submódulo e proibição de compile remoto*.
- [[pop/roadmap/1-organizacao|Epoch 1]] — *Phase 1.2, tasks 1.2.1 (origem) e 1.2.2 (verificação ponta-a-ponta)*.
