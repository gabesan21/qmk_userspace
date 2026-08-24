---
task: F-20260824-desativa-ci
entry: 01-remocao-workflow
---

# Remoção do workflow de build remoto

Removido `.github/workflows/build_binaries.yaml` via `git rm`. O arquivo era **herdado do upstream** (Bastardkb/qmk_userspace) com gatilhos `push` + `workflow_dispatch` — não foi criado pelo PoP — e compilava/publicava firmware no GitHub Actions a cada push. A decisão do dono é compile apenas local e sob demanda, então o arquivo saiu por inteiro, sem gatilho substituto. Desvio do gate autorizado por comando explícito (rule 20), mesmo precedente de [[F-20260824-limpeza-teclados.02-desvio-gate-importacao]].

## Evidence

- [[pop/PROJECT|PROJECT]] — *decisão de 2026-08-24: build sempre local, workflow removido*.
- [[README|README.md]] — *o userspace permanece; só a CI de build saiu*.
