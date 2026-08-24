# Epoch 1 — Organização

- **Project:** [[pop/PROJECT|QMK Userspace]] · **Roadmap:** [[pop/ROADMAP|Roadmap]]
- **Status:** pending
- **Description:** constrói o harness de conhecimento fiel ao repo importado — mapa do projeto, skills operacionais e notas — sem tocar no conteúdo (gate de importação no [[AGENTS|AGENTS]]).

> One phase per section; under each phase, only still-open tasks — **always one-line descriptions**. On closing the `005_closing` stage, remove a task row only after its canonical memory is valid; preserve the epoch, phase, and other open tasks.

## Recon and forks

Recon determinístico do import (2026-08-23) já consolidado no [[pop/PROJECT|PROJECT]]: 14 targets `vendor`, submódulo `modules/bastardkb`, CI de build contra `bastardkb/bastardkb-qmk@main`.

- [x] ~~RECON NEEDED: quais teclados Bastard o dono realmente possui/usa~~ — **resolvido em 2026-08-24:** o dono usa o **Dilemma V3 (36 keys, trackpad Procyon)** → target `bastardkb/dilemma/3x5_3_procyon`, keymap `vendor`.
- [x] ~~RECON NEEDED: toolchain local de build existe?~~ — **resolvido em 2026-08-24:** não — `qmk` CLI e `clang-format` ausentes, submódulo não inicializado; hoje o build só acontece via CI (ou devcontainer, docker disponível).

## Phase 1.1 — Mapa do projeto

- **Status:** pending
- **Description:** inventário navegável dos 14 targets e keymaps `vendor`, consolidado no PROJECT.

| Task | Description (≤1 line) | Status |
|------|-----------------------|--------|
| `1.1.1-inventario-keymaps` | Tabela dos 14 targets (modelo, variante, revisão) + resumo de cada keymap `vendor` no PROJECT. · size: S | not started |
| `1.1.2-phase-verification` | Roda `pop_validate` e confere o mapa contra o repo. · size: S | not started |

## Phase 1.2 — Skills operacionais

- **Status:** pending
- **Description:** como compilar e flashear o firmware, em `pop/skills/`.

| Task | Description (≤1 line) | Status |
|------|-----------------------|--------|
| `1.2.1-skill-build-flash` | Skill de build local (`qmk compile`, devcontainer) + download do artifact da CI + flash (`.uf2`/bootloader). · size: M | not started |
| `1.2.2-phase-verification` | Valida a skill executando o build de ao menos um target (ou via CI). · size: S | not started |

## Phase 1.3 — Notas e decisões

- **Status:** pending
- **Description:** decisões de origem do fork e estratégia de sync com o upstream, em `pop/notes/`.

| Task | Description (≤1 line) | Status |
|------|-----------------------|--------|
| `1.3.1-notas-decisoes` | Notas em `pop/notes/decisions|references`: origem do fork, auto-update do submódulo, política de sync upstream. · size: S | not started |
| `1.3.2-phase-verification` | Roda `pop_validate` e confere links e frontmatter das notas. · size: S | not started |
