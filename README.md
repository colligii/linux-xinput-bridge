# Linux XInput Bridge

An alternative driver and mapping solution for controllers on Linux.

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

If you are unable to use `xboxdrv` or if `xpad` lacks full support for your controller, this project provides a reliable alternative.

`Linux XInput Bridge` is a controller bridge that connects physical gamepads to virtual controllers using the Linux input subsystem (`evdev`/`uinput`). It features a customizable mapping engine, allowing you to intercept, modify, and transform input data exactly to your needs.

The project improves controller compatibility on Linux by reading raw controller events, processing them, and exposing them through a virtual gamepad interface.

## Overview

Many controllers behave differently on Linux depending on how they expose their input events. Some games expect a specific controller format—especially Xbox-compatible input.

Linux XInput Bridge provides a flexible abstraction layer between the physical controller and the virtual controller.

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

* 🎮 Read controller input through Linux `evdev`
* 🔄 Convert raw events into normalized controller data
* ⚙️ Configurable controller profiles
* 🕹️ Virtual Xbox controller support
* 🔌 User-friendly API orchestrated by Node.js
* 🚀 Low-level input processing powered by Go
* 🐧 Tailored for Linux gaming environments

## Architecture

The project utilizes a hybrid architecture:

* **Go** handles low-level controller operations and input streaming.
* **Node.js** manages orchestration, configuration, and application logic.

### Go Components

Located in:

```
lib_go/
```

The Go layer compiles into two distinct binaries.

#### controller-reader

Responsible for reading physical controller events.

**Key Responsibilities:**
* Open Linux input devices.
* Read events from `/dev/input/eventX`.
* Parse buttons and analog inputs.
* Broadcast controller state updates.

Build command:
```bash
go build -o controller-reader read_controller.go
```

---

#### xbox-driver

Responsible for virtual controller communication.

**Key Responsibilities:**
* Receive parsed controller events.
* Translate controller states.
* Expose an Xbox-compatible virtual controller.

Build command:
```bash
go build -o xbox-driver xbox.go
```

---

## Requirements

* Linux operating system
* Go 1.20+
* Node.js 22+
* Root privileges (required for accessing `/dev/input` and `/dev/uinput`)

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/colligii/linux-xinput-bridge.git
   cd linux-xinput-bridge
   ```

2. Configure the repository environment:
   ```bash
   npm run configure
   ```
   > **Note:** This CLI utility sets up required script permissions (such as making `configure.sh` executable) and runs the initial setup. Because it alters file permissions, this step may require root privileges. You can review `configure.sh` at any time to inspect its operations.

3. Build the Go components:
   ```bash
   npm run buildGoLib
   ```

Upon completion, the binaries will be generated under:
```
lib_go/
├── xbox-driver
└── controller-reader
```

## Configuration

Controller settings are managed inside:

```
defaultConfig.json
```

### Example Configuration:

```json
{
    "controller": "gamesir-g7",
    "evtest": "/dev/input/event16"
}
```

* **`evtest`**: The raw event input path corresponding to your physical controller. You can find your device's path by running `sudo evtest` in your terminal.
* **`controller`**: The name of the profile JSON file that maps your specific device layout.

## Running the Application

To launch the bridge, execute:

```bash
npm start
```

When started, the application will automatically:
1. Load your active controller configuration.
2. Spawn the underlying Go processes.
3. Intercept raw physical events.
4. Stream and forward the translated states to the virtual controller interface.

## Supported Controllers

Currently verified and tested:

* GameSir G7

*Need another layout? You can easily map additional controllers by introducing custom JSON profiles.*

## Development

To manually build the Go binaries during development:
```bash
cd lib_go
./buildGo.sh
```

To run the Node.js orchestrator:
```bash
npm start
```

## Future Goals

Planned enhancements and features:

* [ ] Graphical user interface (GUI) configuration tool
* [ ] Input streaming/sharing over network pipelines
* [ ] Remote local multiplayer capability

## Custom Mapping Configuration

The project supports custom controller profiles using a JSON configuration file structure. This allows you to remap buttons, handle composite joystick events, and apply custom processing algorithms for triggers and analog sticks.

### Configuration Structure Explained

Below is the breakdown of each step inside a custom configuration profile (e.g., `gamesir-g7.json`):

```json
{
    "btnMapping": {
        "305": "A",
        "307": "Y",
        "306": "B",
        "304": "X",
        "316": "HOME",
        "312": "BACK",
        "313": "START",
        "308": "LB",
        "309": "RB",
        "310": "LT",
        "314": "L3",
        "315": "R3",
        "311": "RT",
        "16": "DPAD_X",
        "17": "DPAD_Y",
        "1;3": "JOYSTICK_Y",
        "0;3": "JOYSTICK_X",
        "5;3": "JOYSTICK_Z_Y",
        "2;3": "JOYSTICK_Z_X"
    },
    "composedKeys": {
        "1;3": true,
        "0;3": true,
        "5;3": true,
        "2;3": true
    },
    "joystickGroup": [
        ["dpad", ["16", "17"]],
        ["left", ["1;3", "0;3"]],
        ["right", ["5;3", "2;3"]]
    ],
    "customButtonValues": [
        ["LT", "../../config/gamesir-g7-triggesrs.js"],
        ["RT", "../../config/gamesir-g7-triggesrs.js"]
    ],
    "joystickDecodeAlgo": [
        ["left", "../../config/gamesir-g7-joystick.js", { "x": "JOYSTICK_X", "y": "JOYSTICK_Y" } ],
        ["right", "../../config/gamesir-g7-joystick.js", { "x": "JOYSTICK_Z_X", "y": "JOYSTICK_Z_Y" } ]
    ]
}
```

### Step-by-Step Breakdown

#### 1. `btnMapping` (Raw Event Translation)
Maps the raw Linux input event codes (obtained from `evtest`) to meaningful virtual Xbox controller button and axis labels (e.g., code `305` becomes button `A`). It supports both standard codes and composite string IDs for axis events.

#### 2. `composedKeys` (Event Filtering)
Flags specific composite event IDs that require special handling. Each key is a combination of the Linux input event code and type fields in the format code;type (for example, 1;3 or 0;3). Setting a key to true instructs the orchestration engine to group and process all events matching that composite identifier together. This is especially useful for noisy input devices where relying on the code field alone would incorrectly merge unrelated events.

#### 3. `joystickGroup` (Axis Grouping)
Combines individual multi-directional axis events into semantic pairs. Keeping the `joystickGroup` array intact is required to bind inputs into logical groups like `dpad`, `left` stick, or `right` stick so they can be parsed into nested controller states.

#### 4. `customButtonValues` (Trigger Interceptors & Normalization)
Intercepts raw action inputs (like `LT` and `RT`) and routes them through an external JavaScript mapper file. Since controllers like the GameSir G7 or other non-native pads lack out-of-the-box driver support, this array mapper handles external files to translate and normalize pressure values. This is essential for racing games and simulation titles where **triggers must accurately scale from `0` to `255`**.

#### 5. `joystickDecodeAlgo` (Analog Stick Processing & Range Conversion)
Binds the grouped directional sticks (`left` / `right`) to an external JavaScript processing file alongside their axis dictionary mapping. This step handles low-level mathematical conversions, such as resolving coordinate overflows, applying thumbstick deadzones, and scaling the data to match the expected virtual joystick output range of **`-32768` to `32767`**.

---

### Verifying the Live Controller State

During development or testing, you can modify the core functions to instantiate a new driver interface. To inspect the live data stream being piped through the orchestrator, uncomment the `console.log` statement located on **line 46 of `index.js`**. 

When inputs are intercepted and parsed, the console will output the unified state schema in the following structured format:

```javascript
{
  X: 0,
  A: 0,
  B: 0,
  Y: 0,
  LB: 0,
  RB: 0,
  LT: 0, // Scaled from 0 to 255
  RT: 0, // Scaled from 0 to 255
  BACK: 0,
  START: 0,
  L3: 0,
  R3: 0,
  HOME: 0,
  dpad: { 
    DPAD_X: 0, // Ranges from -1 to 1
    DPAD_Y: 0  // Ranges from -1 to 1
  },
  left: { 
    JOYSTICK_Y: 0, // Ranges from -32768 to 32767
    JOYSTICK_X: 0  // Ranges from -32768 to 32767
  },
  right: { 
    JOYSTICK_Z_Y: 0, // Ranges from -32768 to 32767
    JOYSTICK_Z_X: 0  // Ranges from -32768 to 32767
  }
}
```
This guarantees that while individual components manipulate low-level buffers, the high-level application always receives a predictable, standard XInput payload where the **D-Pad coordinates map from `-1` to `1`**.
## Common Errors

### Permission denied when accessing `/dev/input/event*`

If you receive a permission error while accessing `/dev/input/event*`, run:

```bash
npm run addUserPermission
```

Then restart the application:

```bash
npm start
```

This will grant the current user the required permissions to access Linux input devices.

## Motivation

While the Linux input subsystem is highly robust, gamepad compatibility remains fragmented across different hardware vendors and game engines. This project delivers a lightweight, customizable bridge that seamlessly unifies disparate hardware under a standard virtual controller interface.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.