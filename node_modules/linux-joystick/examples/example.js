const { JoystickDevice, listDevices } = require('linux-joystick');

const devicePath = listDevices()[0];
const joystick = new JoystickDevice(devicePath);

joystick.on('button_pressed', (event) => {
    console.log('Pressed button:', event);
});

joystick.on('button_released', (event) => {
    console.log('Released button:', event);
});

joystick.on('button_changed', (event) => {
    console.log('Button changed:', event);
});

joystick.on('axis_changed', (event) => {
    // Axis value ranges from -32767 to +32767
    console.log('Axis changed:', event);
});

console.log("Button 0:", joystick.getButton(0));
console.log("Axis 0:", joystick.getAxis(0));