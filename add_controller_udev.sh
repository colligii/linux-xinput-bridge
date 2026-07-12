#!/bin/bash
set -e

RULE_FILE="/etc/udev/rules.d/99-controle-go-principal.rules"

sudo rm -f "$RULE_FILE"

sudo udevadm control --reload-rules && sudo udevadm trigger --subsystem-match=input

sudo tee "$RULE_FILE" > /dev/null <<'EOF'
SUBSYSTEM=="input", DEVPATH=="/devices/virtual/input/*", ENV{ID_INPUT_JOYSTICK}=="1", MODE="0666", SYMLINK+="input/controle_principal", TAG+="uaccess", TAG+="seat"
EOF

sudo udevadm control --reload-rules && sudo udevadm trigger --subsystem-match=input

echo "Regra do udev atualizada com sucesso."