---
task: F-20260827-homerow-mods
project: qmk-userspace
started: 2026-08-27
finished: 2026-08-27
commit:
pr:
authorization: F-20260827-homerow-mods: direct-fix triage · comando explícito do humano soberano sobre o gate de importação
---

# F-20260827-homerow-mods — reordem dos homerow mods do Dilemma 3x5_3_procyon

- **Delivery:** homerow mods da camada Base reordenados — esquerdo `A`=LALT · `S`=LCTL · `D`=LSFT · `F`=LGUI; direito espelhado (`J`=LGUI · `K`=LSFT · `L`=LCTL · `'`=LALT). Diagrama em `visual-layouts/` sincronizado.
- **Verification:** `qmk compile -kb bastardkb/dilemma/3x5_3_procyon -km vendor` OK; UF2 gerado e copiado para a raiz do userspace.
- **Contract impact:** specs: nenhuma existe ainda (projeto importado; specs nascem work-driven) · DOX: não aplicável.

## Entries

- [[F-20260827-homerow-mods.01-reordem-homerow]] — nova ordem dos mods, interpretação de "LMOD" como LGUI e desvio do gate de importação.

## Links

- **Origin:** comando direto do usuário em 2026-08-27 — *follow for the context that asked for the change*.
