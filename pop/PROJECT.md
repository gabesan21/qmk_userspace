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

- **5 build targets** em `qmk.json`, todos Dilemma com keymap `vendor`: `3x5_2`, `3x5_3`, `3x5_3_procyon`, `4x6_4`, `4x6_4_procyon`. Os demais modelos Bastard (charybdis, scylla, skeletyl, tbkmini) foram removidos do fork em 2026-08-24 ([[F-20260824-limpeza-teclados]]).
- **Keymaps em C** sob `keyboards/bastardkb/dilemma/<variante>/keymaps/vendor/` (`keymap.c` + `config.h`), com combos e tap dance.
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
- **2026-08-24:** o teclado do dono é o **Dilemma V3 (36 keys, trackpad Procyon)** → target principal `bastardkb/dilemma/3x5_3_procyon`, keymap `vendor`. Toolchain local instalada via pacman (`qmk`, `clang`); build local ainda não validado (fase 1.2).
- **2026-08-24 (desvio do gate de importação):** por comando explícito do humano (rule 20), removidos os modelos não-Dilemma do fork — charybdis, scylla, skeletyl, tbkmini — restando os 5 targets Dilemma. Direct fix [[F-20260824-limpeza-teclados]]; o gate segue ativo para o restante.
- **2026-08-24:** **todo build é local — compile remoto proibido.** O workflow `build_binaries.yaml` (herdado do upstream) foi **removido** em 2026-08-24 ([[F-20260824-desativa-ci]]), sem nenhum gatilho substituto: nenhum compile acontece no GitHub Actions. Motivo: os macros do teclado vão carregar variáveis sensíveis. Pendência de segurança a decidir na época dos macros: repo público não pode receber esses segredos — tornar o fork privado ou isolá-los num arquivo gitignored (ex.: `secrets.h`).
