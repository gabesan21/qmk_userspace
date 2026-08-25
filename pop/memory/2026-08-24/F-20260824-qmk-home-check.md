---
task: F-20260824-qmk-home-check
project: qmk-userspace
started: 2026-08-24
finished: 2026-08-24
commit: 01aca48
pr:
authorization: direct-fix triage (rule 13) — defeito revelado pela execução humana do checklist C8
---

# F-20260824-qmk-home-check — correção do check/auto-fix de user.qmk_home

- **Delivery:** `scripts/build_flash.sh` lê o formato real do qmk 1.2 e, com firmware já em `~/qmk_firmware`, grava a config via `qmk config` em vez de re-rodar `qmk setup` (modo "keep" não grava).
- **Verification:** `bash -n` + shellcheck limpos; sed testado contra `None (None)` e `/path (config)`. Reexecução do script: dono (C8 segue aberto).
- **Contract impact:** [[pop/specs/build-flash-script|build-flash-script]] sincronizada (Contract, Invariants, Errors).

## Entries

- [[F-20260824-qmk-home-check.01-check-qmk-home]] — parsing real + auto-fix por `qmk config` no local padrão.

## Links

- **Origin:** saída real do dono em 2026-08-24 — auto-fix rodava `qmk setup` e a correção nunca resolvia.
