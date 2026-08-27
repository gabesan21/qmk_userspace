---
author: agent
created: 2026-08-27
---

# Origem do fork — de onde veio este userspace

> De onde veio este repositório e qual é o estado do conteúdo herdado?

Origem: [[1.3.1-notas-decisoes]] · Decisões canônicas em [[pop/PROJECT|PROJECT]] — esta nota só linka, não repete.

## Fato

Este repo é um fork de [Bastardkb/qmk_userspace](https://github.com/Bastardkb/qmk_userspace), criado em 2026-08. Segue o padrão QMK userspace: o conteúdo (keymaps, módulos) fica fora do `qmk_firmware`, que é referenciado externamente.

- **Sem commits próprios na origem:** o histórico era 100% upstream na importação; as customizações pessoais de keymap ficam para depois da Epoch 1 (Decision 2026-08-23 do [[pop/PROJECT|PROJECT]]).
- **Redução de escopo:** os modelos não-Dilemma foram removidos em 2026-08-24, restando **5 build targets** em `qmk.json`, todos Dilemma com keymap `vendor` (`3x5_2`, `3x5_3`, `3x5_3_procyon`, `4x6_4`, `4x6_4_procyon`). Contexto completo: [[F-20260824-limpeza-teclados]] — *follow para por que a remoção foi um desvio do gate de importação*.
- **Teclado do dono:** apenas `bastardkb/dilemma/3x5_3_procyon` com keymap `vendor` é usado de fato — é o target principal (primeiro entry de `qmk.json`). Inventário detalhado do hardware e das camadas: seção *Inventário — Dilemma V3* do [[pop/PROJECT|PROJECT]].

## Consequência prática

Qualquer mudança de conteúdo (keymap pessoal) deve presumir que o restante do repo é upstream congelado: os 4 targets Dilemma que o dono não usa são mantidos como herdados, sem manutenção ativa — a validação real acontece no target `3x5_3_procyon` (ver [[politica-sync-upstream]]).

## Links

- [[pop/PROJECT|PROJECT]] — seção *Decisions* (2026-08-23 e 2026-08-24) — *follow para as decisões canônicas de importação e target principal*.
- [[F-20260824-limpeza-teclados]] — *follow para o registro da redução aos 5 targets Dilemma*.
