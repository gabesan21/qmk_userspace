# Epoch 1 — Organização

- **Project:** [[pop/PROJECT|QMK Userspace]] · **Roadmap:** [[pop/ROADMAP|Roadmap]]
- **Status:** pending
- **Description:** constrói o harness de conhecimento fiel ao repo importado — mapa do projeto, skills operacionais e notas — sem tocar no conteúdo (gate de importação no [[AGENTS|AGENTS]]).

> One phase per section; under each phase, only still-open tasks — **always one-line descriptions**. On closing the `005_closing` stage, remove a task row only after its canonical memory is valid; preserve the epoch, phase, and other open tasks.

## Recon and forks

Recon determinístico do import (2026-08-23) já consolidado no [[pop/PROJECT|PROJECT]]: 5 targets Dilemma `vendor` (demais modelos removidos em [[F-20260824-limpeza-teclados]]), submódulo `modules/bastardkb`, CI de build **removida** em [[F-20260824-desativa-ci]] — compile é só local.

- [x] ~~RECON NEEDED: quais teclados Bastard o dono realmente possui/usa~~ — **resolvido em 2026-08-24:** o dono usa o **Dilemma V3 (36 keys, trackpad Procyon)** → target `bastardkb/dilemma/3x5_3_procyon`, keymap `vendor`.
- [x] ~~RECON NEEDED: toolchain local de build existe?~~ — **resolvido em 2026-08-24:** `qmk` 1.2.0 e `clang` instalados via pacman; **submódulo ainda não inicializado** e `qmk setup` (clone do `bastardkb/bastardkb-qmk`) ainda não executado — primeira validação de build acontece na task 1.2.2.

## Phase 1.1 — Mapa do projeto

- **Status:** completed (2026-08-24)
- **Description:** inventário do único teclado do dono — Dilemma V3 (`3x5_3_procyon`) — e seu keymap `vendor`, consolidado no PROJECT. *(Reduzida dos 5 targets em 2026-08-24: o dono só tem este teclado.)*

| Task | Description (≤1 line) | Status |
|------|-----------------------|--------|

## Phase 1.2 — Build e flash

- **Status:** pending
- **Description:** script executável de build/flash na raiz do repo (`scripts/`) — rodado pelo próprio dono, sem agente.

| Task | Description (≤1 line) | Status |
|------|-----------------------|--------|
| [[1.2.1-script-build-flash]] | `scripts/build_flash.sh`: checks (qmk, submódulo, `user.qmk_home`), compile do target (padrão `3x5_3_procyon:vendor`) e `--flash` via UF2. · size: M | 001_initial_task |
| `1.2.2-phase-verification` | Roda o script de ponta a ponta (build de ao menos um target; flash opcional) e corrige o que falhar. · size: S | not started |

## Phase 1.3 — Notas e decisões

- **Status:** pending
- **Description:** decisões de origem do fork e estratégia de sync com o upstream, em `pop/notes/`.

| Task | Description (≤1 line) | Status |
|------|-----------------------|--------|
| `1.3.1-notas-decisoes` | Notas em `pop/notes/decisions|references`: origem do fork, auto-update do submódulo, política de sync upstream. · size: S | not started |
| `1.3.2-phase-verification` | Roda `pop_validate` e confere links e frontmatter das notas. · size: S | not started |
