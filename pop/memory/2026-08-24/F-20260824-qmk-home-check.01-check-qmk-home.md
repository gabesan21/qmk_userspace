---
task: F-20260824-qmk-home-check
entry: 01-check-qmk-home
---

# Check de user.qmk_home: parsing real + auto-fix por qmk config

Dois defeitos: (1) o parsing esperava `None` exato, mas o qmk 1.2 imprime `user.qmk_home=<valor> (origem)` — agora captura só o valor; (2) `qmk setup` com clone existente (opção "keep") usa o diretório sem gravar `user.qmk_home`, e a revalidação sempre falhava. Nova cadeia: firmware no local padrão → `qmk config user.qmk_home=<path>`; ausente → `qmk setup -b main`; revalida e fail-fast. No diagnóstico, a config do dono foi gravada.

## Evidence

- [[scripts/build_flash.sh]] — *`qmk_home`/`qmk_home_ok` e a cadeia de auto-fix (commit 01aca48)*.
- [[pop/specs/build-flash-script|build-flash-script]] — *Errors documenta o comportamento e o formato real*.
