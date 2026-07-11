package main

/*
#include <linux/input.h>
#include <linux/uinput.h>
#include <string.h>
#include <stdlib.h>

enum {
    GO_EV_SYN = EV_SYN,
    GO_EV_KEY = EV_KEY,
    GO_EV_ABS = EV_ABS,

    GO_BTN_SOUTH = BTN_SOUTH,
    GO_BTN_EAST  = BTN_EAST,
    GO_BTN_NORTH = BTN_NORTH,
    GO_BTN_WEST  = BTN_WEST,
    GO_BTN_TL    = BTN_TL,
    GO_BTN_TR    = BTN_TR,
    GO_BTN_SELECT= BTN_SELECT,
    GO_BTN_START = BTN_START,
    GO_BTN_MODE  = BTN_MODE,

    GO_ABS_X = ABS_X,
    GO_ABS_Y = ABS_Y,
    GO_ABS_Z = ABS_Z,
    GO_ABS_RX = ABS_RX,
    GO_ABS_RY = ABS_RY,
    GO_ABS_RZ = ABS_RZ,
    GO_ABS_HAT0X = ABS_HAT0X,
    GO_ABS_HAT0Y = ABS_HAT0Y,
	GO_BTN_THUMBL = BTN_THUMBL,
    GO_BTN_THUMBR = BTN_THUMBR,

    // PARAMOS AS CONSTANTES ANTIGAS AQUI E ADICIONAMOS AS DE CORREÇÃO DO WIDGET:
    GO_UI_SET_EVBIT   = UI_SET_EVBIT,
    GO_UI_SET_KEYBIT  = UI_SET_KEYBIT,
    GO_UI_SET_ABSBIT  = UI_SET_ABSBIT,
    GO_UI_ABS_SETUP   = UI_ABS_SETUP,
    GO_UI_DEV_SETUP   = UI_DEV_SETUP,
    GO_UI_DEV_CREATE  = UI_DEV_CREATE,
    GO_UI_DEV_DESTROY = UI_DEV_DESTROY,
};
*/
import "C"

import (
	"bufio"
	"encoding/json"
	"os"
	"syscall"
	"unsafe"
)

// ==========================================
// 1. STRUCTS DE MAPEAMENTO JSON
// ==========================================

type AxisLeft struct {
	X int32 `json:"JOYSTICK_X"`
	Y int32 `json:"JOYSTICK_Y"`
}

type AxisRight struct {
	X int32 `json:"JOYSTICK_Z_X"`
	Y int32 `json:"JOYSTICK_Z_Y"`
}

type DpadState struct {
	X int32 `json:"DPAD_X"`
	Y int32 `json:"DPAD_Y"`
}

type XboxPacket struct {
	X     int32     `json:"X"`
	A     int32     `json:"A"`
	B     int32     `json:"B"`
	Y     int32     `json:"Y"`
	LB    int32     `json:"LB"`
	RB    int32     `json:"RB"`
	LT    int32     `json:"LT"`
	RT    int32     `json:"RT"`
	Back  int32     `json:"BACK"`
	Start int32     `json:"START"`
	Home  int32     `json:"HOME"`
	Dpad  DpadState `json:"dpad"`
	Left  AxisLeft  `json:"left"`
	Right AxisRight `json:"right"`
	L3 int32 `json:"L3"`
    R3 int32 `json:"R3"`
}

const (
	evSyn = uint16(C.GO_EV_SYN)
	evKey = uint16(C.GO_EV_KEY)
	evAbs = uint16(C.GO_EV_ABS)

	btnA = uint16(C.GO_BTN_SOUTH)
	btnB = uint16(C.GO_BTN_EAST)
	btnX = uint16(C.GO_BTN_NORTH)
	btnY = uint16(C.GO_BTN_WEST)

	btnLB = uint16(C.GO_BTN_TL)
	btnRB = uint16(C.GO_BTN_TR)

	btnBack  = uint16(C.GO_BTN_SELECT)
	btnStart = uint16(C.GO_BTN_START)
	btnHome  = uint16(C.GO_BTN_MODE)

	axisLX = uint16(C.GO_ABS_X)
	axisLY = uint16(C.GO_ABS_Y)
	axisLT = uint16(C.GO_ABS_Z)

	axisRX = uint16(C.GO_ABS_RX)
	axisRY = uint16(C.GO_ABS_RY)
	axisRT = uint16(C.GO_ABS_RZ)

	axisDpadX = uint16(C.GO_ABS_HAT0X)
	axisDpadY = uint16(C.GO_ABS_HAT0Y)

	btnL3    = uint16(C.GO_BTN_THUMBL)
    btnR3    = uint16(C.GO_BTN_THUMBR)
)

var (
	uiSetEvBit  = uintptr(C.GO_UI_SET_EVBIT)
	uiSetKeyBit = uintptr(C.GO_UI_SET_KEYBIT)
	uiSetAbsBit = uintptr(C.GO_UI_SET_ABSBIT)

	uiAbsSetup   = uintptr(C.GO_UI_ABS_SETUP)
	uiDevSetup   = uintptr(C.GO_UI_DEV_SETUP)
	uiDevCreate  = uintptr(C.GO_UI_DEV_CREATE)
	uiDevDestroy = uintptr(C.GO_UI_DEV_DESTROY)
)

var uinputFd *os.File

// ==========================================
// 3. FUNÇÕES DE INICIALIZAÇÃO E ESCRITA
// ==========================================

func ioctl(fd uintptr, request uintptr, arg uintptr) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, request, arg)
	if errno != 0 {
		return errno
	}
	return nil
}

func applyDeadzone(value int32, threshold int32) int32 {
	if value > -threshold && value < threshold {
		return 0
	}
	return value
}

func initController() bool {
	
	var err error
	// Executar como root (sudo) é obrigatório para acessar o /dev/uinput
	uinputFd, err = os.OpenFile("/dev/uinput", os.O_WRONLY|syscall.O_NONBLOCK, 0660)
	if err != nil {
		println("Erro ao abrir /dev/uinput (Tente rodar como sudo):", err.Error())
		return false
	}
	fd := uinputFd.Fd()

	// 1. Habilita tipos globais
	ioctl(fd, uiSetEvBit, uintptr(evKey))
	ioctl(fd, uiSetEvBit, uintptr(evAbs))

		// Habilita sincronização
	ioctl(fd, uiSetEvBit, uintptr(C.EV_SYN))

	// Habilita gamepad padrão Xbox
	ioctl(fd, uiSetKeyBit, uintptr(C.BTN_MODE))

	// 2. Habilita botões digitais
	// 2. Habilita o código de cada botão digital individualmente
	// Mudamos o tipo do slice para uint16
	botoes := []uint16{btnA, btnB, btnX, btnY, btnLB, btnRB, btnBack, btnStart, btnHome, btnL3, btnR3}
	for _, btn := range botoes {
		// Fazemos o cast para uintptr aqui na chamada do ioctl
		ioctl(fd, uiSetKeyBit, uintptr(btn))
	}

	// 3. Habilita os códigos dos eixos analógicos
	eixosPrincipais := []uint16{axisLX, axisLY, axisRX, axisRY}
	gatilhos := []uint16{axisLT, axisRT}
	dpads := []uint16{axisDpadX, axisDpadY}

	for _, eix := range eixosPrincipais {
		ioctl(fd, uiSetAbsBit, uintptr(eix))
	}
	for _, eix := range gatilhos {
		ioctl(fd, uiSetAbsBit, uintptr(eix))
	}
	for _, eix := range dpads {
		ioctl(fd, uiSetAbsBit, uintptr(eix))
	}

	// 4. Configura limites dos analógicos principais (-32768 a 32767)
	for _, eix := range eixosPrincipais {
		var setup C.struct_uinput_abs_setup
		setup.code = C.__u16(eix)
		setup.absinfo.value = 0
		setup.absinfo.minimum = -32768
		setup.absinfo.maximum = 32767
		setup.absinfo.fuzz = 256  // Aumentado de 16 para 256 (filtra ruídos de trepidação)
		setup.absinfo.flat = 1024 // Aumentado de 128 para 1024 (região central estável)
		setup.absinfo.resolution = 0

		if err := ioctl(fd, uiAbsSetup, uintptr(unsafe.Pointer(&setup))); err != nil {
			println("Erro no UI_ABS_SETUP (Eixo Principal):", err.Error())
		}
	}

	// 5. Configura limites dos gatilhos LT/RT (0 a 255 comumente no driver do Xbox)
	for _, eix := range gatilhos {
		var setup C.struct_uinput_abs_setup
		setup.code = C.__u16(eix)
		setup.absinfo.value = 0
		setup.absinfo.minimum = 0
		setup.absinfo.maximum = 255 // Mude para 32767 se sua entrada for mapeada em alta resolução
		setup.absinfo.fuzz = 0
		setup.absinfo.flat = 0
		setup.absinfo.resolution = 0

		if err := ioctl(fd, uiAbsSetup, uintptr(unsafe.Pointer(&setup))); err != nil {
			println("Erro no UI_ABS_SETUP (Gatilho):", err.Error())
		}
	}

	// 6. Configura limites do D-pad (-1, 0, 1)
	for _, eix := range dpads {
		var setup C.struct_uinput_abs_setup
		setup.code = C.__u16(eix)
		setup.absinfo.value = 0
		setup.absinfo.minimum = -1
		setup.absinfo.maximum = 1
		setup.absinfo.fuzz = 0
		setup.absinfo.flat = 0
		setup.absinfo.resolution = 0

		if err := ioctl(fd, uiAbsSetup, uintptr(unsafe.Pointer(&setup))); err != nil {
			println("Erro no UI_ABS_SETUP (Dpad):", err.Error())
		}
	}

	// 7. Configura metadados do dispositivo virtual
	var setup C.struct_uinput_setup

	deviceName := C.CString("Microsoft X-Box 360 pad")
	defer C.free(unsafe.Pointer(deviceName))

	C.strncpy(
		&setup.name[0],
		deviceName,
		C.size_t(C.UINPUT_MAX_NAME_SIZE-1),
	)

	// Identidade de um Xbox 360 Controller USB real
	setup.id.bustype = C.BUS_USB
	setup.id.vendor = 0x045e  // Microsoft
	setup.id.product = 0x028e // Xbox 360 Controller
	setup.id.version = 0x0114

	err = ioctl(fd, uiDevSetup, uintptr(unsafe.Pointer(&setup)))
	if err != nil {
		println("Erro no UI_DEV_SETUP:", err.Error())
		uinputFd.Close()
		return false
	}

	// 8. Cria o dispositivo no Kernel do Linux
	err = ioctl(fd, uiDevCreate, 0)
	if err != nil {
		println("Erro ao rodar UI_DEV_CREATE:", err.Error())
		uinputFd.Close()
		return false
	}

	return true
}

func writeEvent(typ uint16, code uint16, value int32) {
	if uinputFd == nil {
		return
	}
	var ev C.struct_input_event
	// Deixamos ev.time zerado (o Kernel do Linux preenche automaticamente ao receber)
	ev._type = C.__u16(typ)
	ev.code = C.__u16(code)
	ev.value = C.__s32(value)

	// Linha errada 'buf.Bytes()' removida. Escrevemos a struct diretamente no fd.
	data := C.GoBytes(unsafe.Pointer(&ev), C.int(C.sizeof_struct_input_event))
	_, _ = uinputFd.Write(data)
}

// ==========================================
// 4. FUNÇÃO PRINCIPAL
// ==========================================

func main() {
	if !initController() {
		println("Erro ao iniciar o controle uinput. Certifique-se de usar sudo.")
		os.Exit(1)
	}
	
	// Limpa o estado inicial destruindo instâncias antigas pendentes se necessário
	defer func() {
		if uinputFd != nil {
			_ = ioctl(uinputFd.Fd(), uiDevDestroy, 0)
			uinputFd.Close()
		}
	}()

	writeEvent(evSyn, 0, 0)
	println("Aguardando JSON perfeito via Stdin...")
	reader := bufio.NewReader(os.Stdin)

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			break
		}

		var p XboxPacket
		if err := json.Unmarshal(line, &p); err != nil {
			println("Erro de parsing:", err.Error())
			continue
		}

		// Envia botões digitais (0 ou 1)
		writeEvent(evKey, btnA, p.A)
		writeEvent(evKey, btnB, p.B)
		writeEvent(evKey, btnX, p.X)
		writeEvent(evKey, btnY, p.Y)
		writeEvent(evKey, btnLB, p.LB)
		writeEvent(evKey, btnRB, p.RB)
		writeEvent(evKey, btnStart, p.Start)
		writeEvent(evKey, btnBack, p.Back)
		writeEvent(evKey, btnHome, p.Home)

		// 2. Filtra a sensibilidade central (uma zona morta de 2000 a 3000 costuma ser o ideal)
		const deadzone = 2500

		posX := applyDeadzone(p.Left.X, deadzone)
		posY := applyDeadzone(p.Left.Y, deadzone)
		posRX := applyDeadzone(p.Right.X, deadzone)
		posRY := applyDeadzone(p.Right.Y, deadzone)

		// Envia analógicos
		writeEvent(evAbs, axisLX, posX)
		writeEvent(evAbs, axisLY, posY)
		writeEvent(evAbs, axisRX, posRX)
		writeEvent(evAbs, axisRY, posRY)
		writeEvent(evAbs, axisLT, p.LT)
		writeEvent(evAbs, axisRT, p.RT)
		writeEvent(evAbs, axisDpadX, p.Dpad.X)
		writeEvent(evAbs, axisDpadY, p.Dpad.Y)
		writeEvent(evKey, btnL3, p.L3)
		writeEvent(evKey, btnR3, p.R3)

		// Sincroniza o frame de inputs completo
		writeEvent(evSyn, 0, 0)
	}
}