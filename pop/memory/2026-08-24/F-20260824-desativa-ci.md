---
task: F-20260824-desativa-ci
project: qmk-userspace
started: 2026-08-24
finished: 2026-08-24
commit: <hash do commit final, anotado após o commit>
pr:
authorization: F-20260824-desativa-ci: direct-fix triage (rule 13) · comando explícito do humano soberano sobre o gate de importação (rule 20)
---

# F-20260824-desativa-ci — remoção do workflow de build remoto

- **Delivery:** `.github/workflows/build_binaries.yaml` removido — nenhum compile remoto resta no repo; build é só local e sob demanda.
- **Verification:** `git ls-files .github/workflows/` confirma a remoção; `pop_validate` sem violações.
- **Contract impact:** specs: nenhuma existe ainda (projeto importado) · DOX: não aplicável (não é aplicação).

## Entries

- [[F-20260824-desativa-ci.01-remocao-workflow]] — remoção do workflow herdado do upstream que compilava a cada push.

## Links

- **Origin:** comando direto do usuário em 2026-08-24 ("CI completamente desativado, sem gatilhos; compile só local").
