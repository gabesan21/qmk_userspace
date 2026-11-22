"""Consulta pontual do firmware Procyon; não inicia a TUI nem altera o teclado."""

from contextlib import ExitStack
import os
from pathlib import Path
import re
import select
import sys
import time

# Descritores das coleções exclusivas VIA (IN/OUT de 32 bytes) e console (IN).
VIA = bytes.fromhex('0660ff0961a1010962150026ff009520750881020963150026ff00952075089102c0')
CONSOLE = bytes.fromhex('0631ff0974a1010975150026ff00952075088102c0')
GET = bytes((0x08, 0, 0x4c))  # custom_get_value, canal do teclado, diagnóstico
ESTADOS = ('hook ainda não executado', 'metade secundária', 'camada fora do mapa',
           'USB inativo', 'fila suspensa', 'buffer parcial ocupado', 'fila cheia',
           'reserva de buffer indisponível', 'relatório enfileirado')


def interfaces():
    grupos = {}
    for node in Path('/sys/class/hidraw').glob('hidraw*'):
        dev = (node / 'device').resolve(strict=True)
        if 'HID_ID=0003:0000A8F8:00001836' not in (dev / 'uevent').read_text().upper().splitlines():
            continue
        desc = (dev / 'report_descriptor').read_bytes()
        if desc in (VIA, CONSOLE):
            grupo = grupos.setdefault(dev.parent.parent, {})
            if desc in grupo:
                raise ValueError('Mais de uma interface idêntica; diagnóstico interrompido.')
            grupo[desc] = node.resolve(strict=True)
    if len(grupos) != 1 or set(next(iter(grupos.values()))) != {VIA, CONSOLE}:
        raise ValueError('Conecte apenas um Dilemma Procyon com console e VIA compatíveis.')
    return next(iter(grupos.values()))


def abrir(node, flags, stack):
    identidade = node.stat()
    fd = os.open('/dev/' + node.name, flags | os.O_NONBLOCK | os.O_NOFOLLOW)
    stack.callback(os.close, fd)
    rdev = os.fstat(fd).st_rdev
    atual = Path(f'/sys/dev/char/{os.major(rdev)}:{os.minor(rdev)}')
    if not os.path.samestat(identidade, atual.stat()):
        raise ValueError('O dispositivo mudou durante a abertura. Rode novamente.')
    return fd


def receber(fd, fim):
    if not select.select([fd], [], [], max(0, fim - time.monotonic()))[0]:
        return None
    dados = os.read(fd, 64)
    if not dados:
        raise OSError('O teclado foi desconectado.')
    return dados


def consultar(via):
    pedido = bytearray(32)
    pedido[:3] = GET
    pedido[16:24] = os.urandom(8)  # Correlaciona a resposta sem ler dados de outros clientes.
    if os.write(via, b'\0' + pedido) != 33:  # Report ID zero para hidraw sem IDs.
        raise OSError('Envio da consulta VIA incompleto.')
    fim = time.monotonic() + 2
    while time.monotonic() < fim:
        r = receber(via, fim)
        if r is None:
            break
        if len(r) != 32 or r[16:24] != pedido[16:24] or r[1:3] != GET[1:3]:
            continue
        if r[0] == 0xff:
            raise ValueError('O VIA respondeu: este firmware não implementa a consulta DLD1.')
        if r[:7] != GET + b'DLD1' or r[7] >= len(ESTADOS):
            raise ValueError('Resposta de diagnóstico incompatível.')
        return r[7], r[8], int.from_bytes(r[9:13], 'little')
    raise TimeoutError('O VIA não respondeu à consulta em 2 s.')


def mostrar(amostra):
    estado, camada, publicados = amostra
    print(f'Firmware: {ESTADOS[estado]}; camada observada: {camada if camada != 255 else "—"}; enfileirados: {publicados}.')


def diagnosticar():
    nodes = interfaces()
    with ExitStack() as stack:
        console = abrir(nodes[CONSOLE], os.O_RDONLY, stack)
        via = abrir(nodes[VIA], os.O_RDWR, stack)
        antes = consultar(via)
        mostrar(antes)
        fim = time.monotonic() + 2
        total = validos = 0
        while time.monotonic() < fim:
            r = receber(console, fim)
            if r is None:
                break
            total += 1
            if re.fullmatch(rb'DLM1:[0-6]\n\x00*', r):
                validos += 1
        depois = consultar(via)
        mostrar(depois)
        print(f'Linux: {total} relatórios de console; {validos} quadros DLM1 completos.')
        if validos:
            print('A telemetria chegou ao host. Se a TUI ainda estiver em espera, investigue o leitor do monitor.')
        elif depois[2] != antes[2]:
            print('O contador avançou sem quadro válido no host. Investigue o transporte do console USB.')
        else:
            print('O contador não avançou. Use os estados acima para investigar o emissor.')


if __name__ == '__main__':
    try:
        diagnosticar()
    except (OSError, ValueError) as e:
        sys.exit(str(e))
