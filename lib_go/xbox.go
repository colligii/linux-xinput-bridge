package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"os"
	"syscall"
)

// ==========================================
// 1. STRUCTS DE MAPEAMENTO JSON (NODE -> GO)
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
}

// ==========================================
// 2. STRUCTS E CONSTANTES NATIVAS DO UINPUT
// ==========================================

type inputID struct {
	Bustype uint16
	Vendor  uint16
	Product uint16
	Version uint16
}

type uinputUserDev struct {
	Name       [80]byte
	ID         inputID
	EffectsMax uint32
	AbsMax     [64]int32
	AbsMin     [64]int32
	AbsFuzz    [64]int32
	AbsFlat    [64]int32
}

type inputEvent struct {
	Time  syscall.Timeval
	Type  uint16
	Code  uint16
	Value int32
}

const (
	// Códigos IOCTL corrigidos e calculados para Linux x86_64 / ARM64
	uiSetEvBit  = 0x40045564 // UI_SET_EVBIT
	uiSetKeyBit = 0x40045565 // UI_SET_KEYBIT
	uiSetAbsBit = 0x40045566 // UI_SET_ABSBIT
	uiDevCreate = 0x5501     // UI_DEV_CREATE

	evSyn = 0x00
	evKey = 0x01
	evAbs = 0x03

	// IDs oficiais de botões do Linux (input-event-codes.h)
	btnA     = 304
	btnB     = 305
	btnX     = 307
	btnY     = 308
	btnLB    = 310
	btnRB    = 311
	btnBack  = 314
	btnStart = 315
	btnHome  = 316

	// IDs oficiais de eixos do Linux
	axisLX    = 0x00 // ABS_X
	axisLY    = 0x01 // ABS_Y
	axisRX    = 0x03 // ABS_RX
	axisRY    = 0x04 // ABS_RY
	axisLT    = 0x02 // ABS_Z (Gatilho Esquerdo)
	axisRT    = 0x05 // ABS_RZ (Gatilho Direito)
	axisDpadX = 16   // ABS_HAT0X
	axisDpadY = 17   // ABS_HAT0Y
)

// Ponteiro global para o descritor do uinput
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

func initController() bool {
	var err error
	uinputFd, err = os.OpenFile("/dev/uinput", os.O_WRONLY|syscall.O_NONBLOCK, 0660)
	if err != nil {
		println("Erro ao abrir /dev/uinput:", err.Error())
		return false
	}
	fd := uinputFd.Fd()

	// Configura os tipos de eventos suportados pelo dispositivo
	ioctl(fd, uiSetEvBit, uintptr(evKey))
	ioctl(fd, uiSetEvBit, uintptr(evAbs))

	// Habilita os botões na tabela de capacidades do Kernel
	botoes := []uintptr{btnA, btnB, btnX, btnY, btnLB, btnRB, btnBack, btnStart, btnHome}
	for _, btn := range botoes {
		ioctl(fd, uiSetKeyBit, btn)
	}

	// Habilita os eixos absolutos e direcionais
	eixos := []uintptr{axisLX, axisLY, axisRX, axisRY, axisLT, axisRT, axisDpadX, axisDpadY}
	for _, eix := range eixos {
		ioctl(fd, uiSetAbsBit, eix)
	}

	// Monta os metadados do dispositivo virtual
	var userDev uinputUserDev
	copy(userDev.Name[:], "Virtual Xbox 360 Controller")
	userDev.ID.Bustype = 0x03  // BUS_USB
	userDev.ID.Vendor = 0x045e  // Microsoft Corporation
	userDev.ID.Product = 0x028e // Xbox 360 Controller

	// Aplica limites para analógicos principais (-32768 a 32767)
	eixosPrincipais := []int{axisLX, axisLY, axisRX, axisRY}
	for _, eix := range eixosPrincipais {
		userDev.AbsMin[eix] = -32768
		userDev.AbsMax[eix] = 32767
	}

	// Gatilhos analógicos (0 a 1023)
	userDev.AbsMin[axisLT] = 0
	userDev.AbsMax[axisLT] = 1023
	userDev.AbsMin[axisRT] = 0
	userDev.AbsMax[axisRT] = 1023

	// Eixos do D-Pad digital mapeado em ABS (-1 a 1)
	userDev.AbsMin[axisDpadX] = -1
	userDev.AbsMax[axisDpadX] = 1
	userDev.AbsMin[axisDpadY] = -1
	userDev.AbsMax[axisDpadY] = 1

	// Escreve as propriedades físicas na pipeline do dispositivo virtual
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, userDev)
	_, err = uinputFd.Write(buf.Bytes())
	if err != nil {
		println("Erro ao escrever uinput_user_dev:", err.Error())
		uinputFd.Close()
		return false
	}

	// Cria o dispositivo mapeado no sistema operacional (/dev/input/eventX)
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

	var ev inputEvent
	ev.Type = typ
	ev.Code = code
	ev.Value = value

	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, ev)
	_, _ = uinputFd.Write(buf.Bytes())
}

// ==========================================
// 4. FUNÇÃO PRINCIPAL (LOOP DE ENTRADA STDIN)
// ==========================================

func main() {
	if !initController() {
		println("Erro ao iniciar o controle uinput. Certifique-se de usar SUDO.")
		os.Exit(1)
	}
	writeEvent(evSyn, 0, 0)

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

		println("Go leu -> A:", p.A, "LX:", p.Left.X, "LY:", p.Left.Y, "RX:", p.Right.X)

		// Escreve os estados de botões
		writeEvent(evKey, btnA, p.A)
		writeEvent(evKey, btnB, p.B)
		writeEvent(evKey, btnX, p.X)
		writeEvent(evKey, btnY, p.Y)
		writeEvent(evKey, btnLB, p.LB)
		writeEvent(evKey, btnRB, p.RB)
		writeEvent(evKey, btnStart, p.Start)
		writeEvent(evKey, btnBack, p.Back)
		writeEvent(evKey, btnHome, p.Home)

		analogLT := p.LT * 32767
		analogRT := p.RT * 32767

		// Escreve os analógicos e gatilhos
		writeEvent(evAbs, axisLX, p.Left.X)
		writeEvent(evAbs, axisLY, p.Left.Y)
		writeEvent(evAbs, axisRX, p.Right.X)
		writeEvent(evAbs, axisRY, p.Right.Y)
		writeEvent(evAbs, axisLT, analogLT)
		writeEvent(evAbs, axisRT, analogRT)
		writeEvent(evAbs, axisDpadX, p.Dpad.X)
		writeEvent(evAbs, axisDpadY, p.Dpad.Y)

		// Despacha todos os inputs pendentes ao kernel em um único tick síncrono
		writeEvent(evSyn, 0, 0)
	}
}