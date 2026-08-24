# QMK Userspace (Bastard Keyboards)

- **Status:** planning
- **Priority:** medium
- **Created on:** 2026-08-23
- **Roadmap:** [[pop/ROADMAP|Roadmap]]

## Goal

Userspace QMK pessoal para os teclados Bastard Keyboards: keymaps `vendor` versionados, compilados em firmware (`.uf2`) pela CI, prontos para flash.

## Context

Fork de [Bastardkb/qmk_userspace](https://github.com/Bastardkb/qmk_userspace) criado em 2026-08 — **ainda sem customizações próprias** (nenhum commit do dono; histórico é 100% upstream). O repo segue o padrão QMK userspace: o conteúdo fica fora do `qmk_firmware`, que é referenciado externamente.

Mapa do recon (2026-08-23):

- **14 build targets** em `qmk.json`, todos keymap `vendor`: Charybdis (3x5, 3x6, 4x6 + variantes `_left`, revisão `splinktegrated_rev1`), Dilemma (3x5_2, 3x5_3, 4x6_4 + variantes `_procyon`), Scylla, Skeletyl e TBK Mini.
- **Keymaps em C** sob `keyboards/bastardkb/<modelo>/<variante>/keymaps/vendor/` (`keymap.c` + `config.h`), com combos e tap dance.
- **`modules/bastardkb`** é submódulo git de [Bastardkb/qmk_modules](https://github.com/Bastardkb/qmk_modules) (o módulo de pointing); o workflow `update-submodule.yml` atualiza a referência automaticamente via `repository_dispatch`.
- **Build:** local via `qmk compile` (requer `qmk config user.qmk_home` apontando para o firmware); CI em `.github/workflows/build_binaries.yaml` usa o workflow reutilizável do QMK contra `bastardkb/bastardkb-qmk@main` e publica os binários.
- **Suporte a devcontainer** (`.devcontainer/`) para build sem toolchain local.
- Frágil/a observar: o repo depende de moving targets (submódulo e firmware upstream em `main`); sem testes locais — a validação é o build.

## Folder structure

Anatomia padrão (ver AGENTS.md da raiz do vault). O conteúdo (firmware userspace: `keyboards/`, `layouts/`, `users/`, `modules/`) vive na raiz do projeto — a própria pasta **é o repositório** (`uni-repo`).

## Agent harness

- **Type and repositories:** `uni-repo`, declarado no [[AGENTS|AGENTS do projeto]] — PR branch `main`.
- **Worktree per task:** yes (padrão) — worktrees em `pop/worktrees/`, gitignored.
- **Tools and restrictions:** builds locais exigem QMK CLI configurado (`qmk config user.qmk_home`); o submódulo `modules/bastardkb` precisa estar inicializado (`git submodule update --init`). Nunca commitar artefatos de build (`*.hex`, `*.bin`, `*.uf2` — já gitignored).
- **Tasks critical by default?** no — crítica é a tarefa que muda o comportamento de uma tecla de boot/reset ou o pipeline de build/CI.
- **Skills:** `pop/skills/` vazio por ora — nasce na Epoch 1 (skill de build/flash).

## Decisions

- **2026-08-23:** importado como `uni-repo` com PR branch `main`; idioma do projeto pt-BR. Fork ainda espelha o upstream — as customizações pessoais de keymap serão planejadas pelo `plan-roadmap` após a Epoch 1.
- **2026-08-24:** o teclado do dono é o **Dilemma V3 (36 keys, trackpad Procyon)** → target principal `bastardkb/dilemma/3x5_3_procyon`, keymap `vendor`. Toolchain local ausente (`qmk`, `clang-format`, submódulo) — build hoje só via CI ou devcontainer; decisão sobre instalar toolchain local fica para a Epoch 1 (fase 1.2).
