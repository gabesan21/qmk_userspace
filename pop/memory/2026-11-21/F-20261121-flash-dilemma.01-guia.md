---
task: F-20261121-flash-dilemma
entry: 01-guia
---

# Preparação e flash guiados

O ponto de entrada instala dependências ausentes no Arch só após confirmação do pacman. Compila `vendor` uma vez e reaproveita o UF2 para gravar direita e esquerda por padrão, sem tocar o USB em `--compile-only`. Na gravação, exige um único `RPI-RP2` gravável com identificação RP2040 e aguarda a saída do bootloader antes de avançar.

## Evidência

- [[flash.sh]] — ponto de entrada do dono.
- [[scripts/build_flash.sh]] — preparação do QMK e gravação nas metades.
- [[pop/specs/build-flash-script|Build e flash]] — uso e teste físico.
