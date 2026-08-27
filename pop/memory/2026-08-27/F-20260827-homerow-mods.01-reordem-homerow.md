---
task: F-20260827-homerow-mods
entry: 01-reordem-homerow
---

# Reordem dos homerow mods

Pedido: esquerdo `LALT-LCTRL-LSHIFT-LMOD`, direito espelhado. "LMOD" não é keycode QMK — interpretado como **LGUI** (Super). Antes: esq. `LGUI LALT LCTL LSFT`; depois: esq. `LALT LCTL LSFT LGUI`, dir. `LGUI LSFT LCTL LALT`. Só a camada Base mudou; as linhas de mods crus das outras camadas mantiveram a ordem antiga (inconsistência consciente). Gate de importação sobreposto por comando explícito do humano (ver [[F-20260824-limpeza-teclados.02-desvio-gate-importacao]]).

## Evidence

- `keyboards/bastardkb/dilemma/3x5_3_procyon/keymaps/vendor/keymap.c` — Base na nova ordem.
- `visual-layouts/layers_dillema_3x5_3_procyon.md` — diagrama sincronizado.
- Build `qmk compile` do target OK → UF2 na raiz.
