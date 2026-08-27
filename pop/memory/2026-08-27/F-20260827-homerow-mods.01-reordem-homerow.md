---
task: F-20260827-homerow-mods
entry: 01-reordem-homerow
---

# Reordem dos homerow mods

O humano pediu os homerow mods da Base na ordem esquerda `LALT - LCTRL - LSHIFT - LMOD`, direito espelhado. "LMOD" não existe como keycode QMK; interpretado como **LGUI** (Super/Win/Cmd), a tecla que ocupava o mindinho na ordem anterior. Se a intenção era outra, corrigir depois.

Antes: esq. `LGUI LALT LCTL LSFT` / dir. `LSFT LCTL LALT LGUI`.
Depois: esq. `LALT LCTL LSFT LGUI` / dir. `LGUI LSFT LCTL LALT`.

Alterado só o keymap da Base (`keyboards/bastardkb/dilemma/3x5_3_procyon/keymaps/vendor/keymap.c`). As linhas de mods crus das camadas Function, Navigation, Pointer, Numeral e Symbols **mantiveram a ordem antiga** — inconsistência consciente, registrada aqui para decisão futura.

Gate de importação: a Epoch 1 (Organização) ainda está `in progress`, mas o comando explícito do humano é soberano e sobrepõe o gate apenas neste escopo — mesmo precedente de [[F-20260824-limpeza-teclados.02-desvio-gate-importacao]].

## Evidence

- `keyboards/bastardkb/dilemma/3x5_3_procyon/keymaps/vendor/keymap.c` — linha da camada Base com a nova ordem.
- `visual-layouts/layers_dillema_3x5_3_procyon.md` — diagrama da Base sincronizado.
- Build local: `qmk compile -kb bastardkb/dilemma/3x5_3_procyon -km vendor` → `bastardkb_dilemma_3x5_3_procyon_vendor.uf2` gerado na raiz.
