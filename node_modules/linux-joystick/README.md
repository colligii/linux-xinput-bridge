# linux-joystick

A simple package to get joystick/gamepad information via linux device files, without any dependencies and prebuilt binaries.

## Installation

```
npm install linux-joystick
```

## Usage

```js
const { JoystickDevice, listDevices } = require("linux-joystick");

const devicePath = listDevices()[0];
const joystick = new JoystickDevice(devicePath);

joystick.on("button_pressed", (event) => {
  console.log("Pressed button:", event);
});
```

### Methods:

- getButton(number)
- getAxis(number)

### Events:

- button_changed
- axis_changed
- button_pressed
- button_released
- button_init
- axis_init

### Helper functions

- listDevices() -> Returns an array with all joystick/gamepad devices.
