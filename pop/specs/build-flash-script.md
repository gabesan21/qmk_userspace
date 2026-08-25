---
id: build-flash-script
project: QMK Userspace
domain: build
kind: contract
status: draft
implementation: planned
origin: "1.2"
created: 2026-08-24
updated: 2026-08-24
supersedes: []
superseded_by:
---

# Spec — Script de build/flash (`scripts/build_flash.sh`)

## Contract

O userspace tem um único caminho documentado de build/flash local: o script `scripts/build_flash.sh` na raiz do repo, executado pelo dono. O script **detecta e executa as correções documentadas** — ao falhar um check de preflight, ele mesmo roda `git submodule update --init modules/bastardkb`, `qmk config user.qmk_home=<path>` (firmware já clonado no local padrão) ou `qmk setup -b main bastardkb/bastardkb-qmk`, revalida e segue; só aborta se a correção não resolver. A instalação do QMK CLI segue como decisão do dono (detecta e instrui).

## Expected behavior

- Dado o repo clonado, quando o dono roda `scripts/build_flash.sh`, então o script exibe um CLI com hierarquia clara (banner, seções por passo, checklist alinhado, resumo final) e executa preflight → compile; qualquer falha de ambiente dispara a correção documentada, que é executada e revalidada antes de seguir.
- Dado um preflight aprovado, quando o compile termina, então o caminho do `.uf2` gerado é exibido.
- Dado o compile concluído sem `--flash`, quando o script termina, então nenhuma gravação no teclado foi tentada e o comando para gravar é exibido.
- Dado `--flash` e o teclado em bootloader, quando o volume `RPI-RP2` aparece (até 60 s, com spinner de espera), então o `.uf2` é copiado para ele e o sucesso é confirmado.

## Invariants

- Flash nunca é implícito: somente atrás da flag explícita `--flash`.
- Toda falha de preflight dispara o auto-fix documentado; a correção é revalidada e, se não resolver, o script termina imediatamente com o comando corretivo exato (fail-fast, saída `1`).
- O auto-fix é limitado às correções documentadas (init do submódulo; `qmk config user.qmk_home` para firmware no local padrão; `qmk setup -b main` do fork); o script não instala programas — o QMK CLI ausente é apenas detectado e instruído.
- Mensagens ao usuário em pt-BR; cores desligadas fora de terminal ou com `NO_COLOR`.
- Nenhum artefato de build (`*.uf2`, `*.hex`, `*.bin`) é gerado dentro do repo (o compile escreve em `user.qmk_home`).

## Interfaces

- **Input (CLI):** `-k|--keyboard <alvo>` (padrão `bastardkb/dilemma/3x5_3_procyon`), `-m|--keymap <nome>` (padrão `vendor`), `--flash` (opt-in), `-h|--help`. Sem argumentos posicionais.
- **Output:** banner de abertura, passos numerados (preflight, compile, flash) em seções com régua, checklist com `✔`/`✘`/`➜`/`⚙` em colunas alinhadas, caixa de instruções de bootloader, spinner na espera do volume e resumo final; erros e correções em stderr; `--help` documenta todas as opções, os padrões, os códigos de saída e exemplos.
- **Códigos de saída:** `0` sucesso · `1` falha de preflight · `2` uso inválido · `3` falha no compile · `4` falha no flash.
- **Ambiente consumido:** `qmk` CLI no PATH; submódulo `modules/bastardkb` inicializado; `user.qmk_home` configurado e existente; volume UF2 `RPI-RP2` montado em `/run/media/$USER`, `/media/$USER` ou `/mnt` (modo bootloader).

## Errors and limits

- **qmk ausente:** erro + instrução de instalação (pacman ou pip); saída `1`. Sem auto-fix.
- **Submódulo não inicializado/vazio:** o script executa `git submodule update --init modules/bastardkb` e revalida; se não resolver, erro + comando corretivo; saída `1`.
- **`user.qmk_home` ausente ou inválido:** se o firmware já existe no local padrão (`~/qmk_firmware`), o script grava a config com `qmk config user.qmk_home=<path>` — `qmk setup` com clone existente (opção "keep") não grava a config; senão executa `qmk setup -b main bastardkb/bastardkb-qmk`. Revalida; se não resolver, erro + comando corretivo; saída `1`. A leitura tolera o formato real do qmk 1.2 (`user.qmk_home=<valor> (origem)`).
- **Compile falha ou `.uf2` não encontrado:** erro apontando a saída do compile/`qmk doctor`; saída `3`.
- **Volume UF2 não aparece em 60 s ou cópia falha:** erro + instrução de confirmar o modo bootloader; saída `4`. Sem automação de montagem além da detecção simples do volume.

## Conformance criteria

- [ ] `bash -n` e leitura confirmam bash válido, shebang e fail-fast (`set -euo pipefail`).
- [ ] Os 3 checks de preflight existem; os dois documentados (submódulo e `user.qmk_home`) executam o auto-fix, revalidam e falham rápido se persistir; `qmk` ausente só instrui.
- [ ] O default é `bastardkb/dilemma/3x5_3_procyon` + `vendor`, sobrescrevível por `-k`/`-m` e documentado no `--help`.
- [ ] Sem `--flash`, nenhuma gravação é tentada; com `--flash`, instruções de bootloader são exibidas antes da espera pelo volume.
- [ ] Execução ponta-a-ponta pelo dono (compile gera o `.uf2`; `--flash` grava via UF2) — checklist humano da task de origem.

## Out of scope

- Instalação de programas (QMK CLI e afins) e qualquer correção de ambiente além das duas documentadas: instruída, jamais executada pelo script.
- Flash por outro mecanismo (dfu, avrdude etc.) e automação de montagem de volumes.
- Build remoto/CI — proibido por decisão do projeto (ver [[pop/PROJECT|PROJECT]]).

## Open questions

- O nome do volume do bootloader (`RPI-RP2`) e seus pontos de montagem cobrem o sistema do dono; validar na primeira execução real (task 1.2.2).

## Related references

- [[pop/PROJECT|PROJECT]] — *target padrão, submódulo e proibição de compile remoto*.
- [[pop/roadmap/1-organizacao|Epoch 1]] — *Phase 1.2, tasks 1.2.1 (origem) e 1.2.2 (verificação ponta-a-ponta)*.
