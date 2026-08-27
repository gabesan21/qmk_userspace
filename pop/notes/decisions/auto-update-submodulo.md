---
author: agent
created: 2026-08-27
---

# Auto-update do submódulo `modules/bastardkb` — o que o workflow realmente faz

> Como a referência do submódulo `modules/bastardkb` é atualizada automaticamente?

Origem: [[1.3.1-notas-decisoes]] · Fonte primária: `.github/workflows/update-submodule.yml` — a nota descreve só o observado ali.

## Fato

`modules/bastardkb` é um submódulo git apontando para [Bastardkb/qmk_modules](https://github.com/Bastardkb/qmk_modules) (ver `.gitmodules`) — o módulo de pointing dos teclados Bastard. A atualização da referência é automática, via o **único workflow restante** do repo, `update-submodule.yml`:

1. **Gatilho:** `repository_dispatch` com o tipo `submodule-updated` — disparado de fora (presumivelmente pelo repo upstream do módulo ao receber push); não há `schedule` nem `push` como gatilho.
2. **Checkout:** `actions/checkout@v4` com `submodules: true` e `fetch-depth: 0`.
3. **Update:** `git submodule update --remote --merge` — atualiza **todos** os submódulos do repo para o último commit remoto (na prática, só existe o `modules/bastardkb`).
4. **Commit:** se houver diff, o bot `github-actions[bot]` commita `chore: update submodule reference to latest` e dá push direto em `main`; sem diff, só loga "Submodule is already up to date!".

## Contexto e limites

- O workflow de build remoto herdado do upstream (`build_binaries.yaml`) foi **removido** em 2026-08-24 — nenhum compile acontece na CI; registro em [[F-20260824-desativa-ci]] — *follow para o motivo (segredos nos macros) e a verificação da remoção*.
- Consequência: o auto-update commita uma nova referência **sem nenhuma validação automática** — a confirmação de que o submódulo novo ainda compila é o build local, descrito em [[politica-sync-upstream]].
- Push direto em `main` pelo bot convive com o PR branch `main` declarado no harness ([[pop/PROJECT|PROJECT]], seção *Agent harness*) — *follow para a configuração do uni-repo*.

## Links

- [[F-20260824-desativa-ci]] — *follow para por que este é o único workflow do repo*.
- [[politica-sync-upstream]] — *follow para como um update automático é validado (build local)*.
