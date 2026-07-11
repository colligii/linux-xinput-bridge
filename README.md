# Linux XInput Bridge

A Linux controller bridge that connects physical gamepads with virtual controllers using the Linux input subsystem.

The project is designed to improve controller compatibility on Linux by reading raw controller events, processing them, and exposing them through a virtual gamepad interface.

## Overview

Many controllers work differently on Linux depending on how they expose their input events. Some games expect a specific controller format, especially Xbox-compatible input.

Linux XInput Bridge provides a flexible layer between the physical controller and the virtual controller.

```
Physical Controller
        |
        v
Linux evdev (/dev/input/eventX)
        |
        v
controller-reader
        |
        v
Node.js Bridge
        |
        v
xbox-driver
        |
        v
Virtual Controller
```

## Features

* 🎮 Read controller input through Linux evdev
* 🔄 Convert raw events into normalized controller data
* ⚙️ Configurable controller profiles
* 🕹️ Virtual Xbox controller support
* 🔌 Node.js and Go integration
* 🚀 Low-level input processing using Go
* 🐧 Built for Linux gaming environments

## Architecture

The project uses a hybrid architecture:

* **Go** handles low-level controller operations
* **Node.js** handles orchestration, configuration, and application logic

### Go Components

Located in:

```
lib_go/
```

The Go layer contains two binaries.

### controller-reader

Responsible for reading physical controller events.

Responsibilities:

* Open Linux input devices
* Read events from `/dev/input/eventX`
* Parse buttons and analog inputs
* Send controller state updates

Build:

```bash
go build -o controller-reader read_controller.go
```

---

### xbox-driver

Responsible for virtual controller communication.

Responsibilities:

* Receive controller events
* Translate controller states
* Expose a virtual Xbox-compatible controller

Build:

```bash
go build -o xbox-driver xbox.go
```

---

## Requirements

* Linux operating system
* Go 1.20+
* Node.js 18+
* Access to `/dev/input`

## Installation

Clone the repository:

```bash
git clone https://github.com/colligii/linux-xinput-bridge.git

cd linux-xinput-bridge
```

Install Node dependencies:

```bash
npm install
```

Build Go components:

```bash
cd lib_go

chmod +x buildGo.sh
./buildGo.sh
```

This will generate:

```
lib_go/
├── xbox-driver
└── controller-reader
```

## Configuration

Controller settings are stored in:

```
defaultConfig.json
```

Example:

```json
{
    "controller": "gamesir-g7",
    "evtest": "/dev/input/event16"
}
```

Configuration defines:

* Controller profile
* Linux input device path

## Running

Start the bridge:

```bash
npm start
```

The application will:

1. Load the controller configuration
2. Start the Go processes
3. Read controller events
4. Forward input to the virtual controller

## Supported Controllers

Currently tested:

* GameSir G7
* Xbox-compatible controllers

More controller profiles can be added through custom mappings.

## Development

Build Go components manually:

```bash
cd lib_go
./buildGo.sh
```

Run the Node.js bridge:

```bash
npm start
```

## Future Goals

Planned improvements:

* [ ] Automatic controller detection
* [ ] GUI configuration tool
* [ ] More controller profiles
* [ ] Better force feedback support
* [ ] Controller sharing over network
* [ ] Remote local multiplayer support

## Motivation

Linux has a powerful input subsystem, but controller compatibility can still be inconsistent between devices and games.

This project aims to provide a customizable bridge that allows different controllers to work reliably through a unified virtual controller interface.

## License

MIT License
