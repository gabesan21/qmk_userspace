package main

// A ordem é a de LAYOUT_split_3x5_3 no keymap vendor: três fileiras
// de dez posições (esquerda, direita), depois os seis polegares.
type layerMap struct {
	name   string
	keys   [36]string
	legend []string
}

const (
	inactive    = "×"
	transparent = "▽"
	modifiers   = "Alt, Ctl=Ctrl, Sh=Shift, Gui=Super: esquerdos; Gr=Alt direito."
)

var layerMaps = [...]layerMap{
	{
		name: "Base",
		keys: [36]string{
			"Q·`", "W·~", "E·'", "R", "T", "Y", "U·\"", "I·^", "O", "P",
			"A·Alt", "S·Ctl", "D·Sh", "F·Gui", "G", "H", "J·Gui", "K·Sh", "L·Ctl", "'·Alt",
			"Z·4", "X·Gr", "C", "V", "B", "N", "M", ",", ".·Gr", "/·4",
			"Esc·3", "Tab·5", "Esp·2", "Ent·6", "Bsp·1", "Mudo",
		},
		legend: []string{
			modifiers,
			"Número ao segurar: 1 Funções · 2 Navegação · 3 Mídia",
			"4 Ponteiro · 5 Números · 6 Símbolos",
			"Esc=Escape  Tab=tabulação  Esp=espaço  Ent=Enter",
			"Bsp=retrocesso  Mudo=alternar som",
			"No us-intl, ` ~ ' \" ^ são acentos mortos ao segurar.",
		},
	},
	{
		name: "Funções",
		keys: [36]string{
			inactive, inactive, inactive, inactive, inactive, "Tela", "F7", "F8", "F9", "F12",
			"Alt", "Ctl", "Sh", "Gui", inactive, "ScrLk", "F4", "F5", "F6", "F11",
			inactive, "Gr", inactive, inactive, inactive, "Pausa", "F1", "F2", "F3", "F10",
			inactive, inactive, inactive, inactive, transparent, inactive,
		},
		legend: []string{
			modifiers,
			"F1–F12: teclas de função. Tela=Print Screen (captura).",
			"ScrLk=Scroll Lock (trava de rolagem)  Pausa=Pause/Break",
		},
	},
	{
		name: "Navegação",
		keys: [36]string{
			inactive, inactive, inactive, inactive, inactive, inactive, inactive, inactive, inactive, inactive,
			"Alt", "Ctl", "Sh", "Gui", inactive, "Caps", "←", "↓", "↑", "→",
			inactive, "Gr", inactive, inactive, inactive, "Ins", "Iníc", "Pg↓", "Pg↑", "Fim",
			inactive, inactive, transparent, "Ent", "Del", inactive,
		},
		legend: []string{
			modifiers,
			"Caps=trava de maiúsculas  Ins=Insert  Setas=cursor",
			"Iníc=Home  Fim=End  Pg↓/Pg↑=página abaixo/acima",
			"Ent=Enter  Del=Delete (apagar à direita)",
		},
	},
	{
		name: "Mídia",
		keys: [36]string{
			inactive, "Luz−", "Luz", "Luz+", inactive, inactive, "Luz−", "Luz", "Luz+", inactive,
			"Ant", "Vol−", "Mudo", "Vol+", "Próx", "Ant", "Vol−", "Mudo", "Vol+", "Próx",
			inactive, inactive, inactive, inactive, inactive, inactive, inactive, inactive, inactive, inactive,
			transparent, "Tocar", "Parar", "Parar", "Tocar", "Mudo",
		},
		legend: []string{
			"Luz=liga/desliga RGB; Luz−/+ = efeito anterior/seguinte.",
			"Ant/Próx=faixa anterior/seguinte  Vol−/+=baixar/subir volume",
			"Mudo=alternar som  Tocar=reproduzir/pausar  Parar=interromper",
		},
	},
	{
		name: "Ponteiro",
		keys: [36]string{
			"Boot", "Limpa", inactive, "DPI", "DPIp", "DPIp", "DPI", inactive, "Limpa", "Boot",
			"Alt", "Ctl", "Sh", "Gui", inactive, inactive, "Gui", "Sh", "Ctl", "Alt",
			transparent, "Rol", "Mira", "B3", inactive, inactive, "B3", "Mira", "Rol", transparent,
			"B3", "B2", "B1", "B1", "B2", "B3",
		},
		legend: []string{
			"Alt, Ctl=Ctrl, Sh=Shift e Gui=Super: modificadores esquerdos.",
			"Boot=entrar no bootloader  Limpa=limpar EEPROM do teclado",
			"DPI=alternar resolução  DPIp=alternar resolução de precisão",
			"Rol=rolagem pelo trackpad  Mira=modo de precisão",
			"B1/B2/B3=botão esquerdo/direito/central do mouse",
		},
	},
	{
		name: "Números",
		keys: [36]string{
			inactive, inactive, inactive, inactive, inactive, "*", "7", "8", "9", "+",
			"Alt", "Ctl", "Sh", "Gui", inactive, "/", "4", "5", "6", "-",
			inactive, "Gr", inactive, inactive, inactive, "NumLk", "1", "2", "3", "=",
			inactive, transparent, inactive, "EntN", inactive, "0",
		},
		legend: []string{
			modifiers,
			"NumLk=trava numérica  EntN=Enter do teclado numérico",
			"Os dígitos são da fileira numérica; não dependem de NumLk.",
		},
	},
	{
		name: "Símbolos",
		keys: [36]string{
			"!", "@", "#", "$", "%", "&", "*", "(", ")", "-·_",
			"Alt", "Ctl", "Sh", "Gui", inactive, inactive, ";·:", "[·{", "]·}", "=·+",
			inactive, "Gr", inactive, inactive, inactive, inactive, "<", ">", "\\·|", "?",
			inactive, inactive, "Esp", transparent, inactive, inactive,
		},
		legend: []string{
			modifiers,
			"Esp=espaço. Nos pares, segurar produz o segundo sinal.",
		},
	},
}
