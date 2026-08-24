---
task: F-20260824-limpeza-teclados
project: qmk-userspace
started: 2026-08-24
finished: 2026-08-24
commit: 16181ac
pr:
authorization: F-20260824-limpeza-teclados: direct-fix triage (rule 13) · comando explícito do humano soberano sobre o gate de importação (rule 20)
---

# F-20260824-limpeza-teclados — remoção dos modelos não-Dilemma

- **Delivery:** fork reduzido aos 5 targets Dilemma; removidos charybdis, scylla, skeletyl e tbkmini.
- **Verification:** `qmk.json` íntegro (JSON válido, 5 targets) e `pop_validate` sem violações; build não executado (toolchain ainda não validada).
- **Contract impact:** specs: nenhuma existe ainda (projeto importado; specs nascem work-driven) · DOX: não aplicável (não é aplicação).

## Entries

- [[F-20260824-limpeza-teclados.01-remocao-modelos]] — remoção dos 4 diretórios de teclados e dos 9 targets no `qmk.json`.
- [[F-20260824-limpeza-teclados.02-desvio-gate-importacao]] — desvio registrado: conteúdo alterado durante a Epoch 1 por comando explícito do humano.

## Links

- **Origin:** comando direto do usuário em 2026-08-24 — *follow for the context that asked for the change*.
