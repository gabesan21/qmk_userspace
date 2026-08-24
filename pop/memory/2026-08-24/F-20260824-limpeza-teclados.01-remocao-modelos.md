---
task: F-20260824-limpeza-teclados
entry: 01-remocao-modelos
---

# Remoção dos modelos não-Dilemma

Removidos via `git rm` os diretórios `keyboards/bastardkb/{charybdis,scylla,skeletyl,tbkmini}` e os 9 targets correspondentes no `qmk.json`, que passa a listar apenas os 5 targets Dilemma (`3x5_2`, `3x5_3`, `3x5_3_procyon`, `4x6_4`, `4x6_4_procyon`, todos keymap `vendor`). Motivo: limpar o fork para refletir o uso real — o dono tem apenas o Dilemma V3. O submódulo `modules/bastardkb` permanece: o Dilemma o usa (pointing/trackpad).

## Evidence

- [[qmk.json]] — *a lista de build targets reduzida aos 5 Dilemma*.
- [[README|README.md]] — *o userspace que a limpeza manteve intacto (agora só com keymaps Dilemma)*.
