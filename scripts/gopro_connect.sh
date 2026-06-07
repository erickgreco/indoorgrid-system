#!/bin/bash

echo "Scanning for GoPro..."

bluetoothctl --timeout 15 scan on > /dev/null 2>&1 || true
SCAN_PID=$!

DevicesName="GoPro"

sleep 3

MAC=""
for i in $(seq 1 15); do
    MAC=$(bluetoothctl devices | grep "$DevicesName" | cut -d' ' -f2)
    echo "$MAC"
    if [ -n "$MAC" ]; then
        break
    fi
    sleep 1
done

kill $SCAN_PID 2>/dev/null || true
bluetoothctl scan off > /dev/null 2>&1 || true

if [ -z "$MAC" ]; then
    echo "Error: GoPro not found. Make sure it is powered on and in range."
    exit 1
fi

NAME=$(bluetoothctl devices | grep "$MAC" | cut -d' ' -f3-)
echo "Found: $NAME ($MAC)"

if bluetoothctl info "$MAC" 2>/dev/null | grep -q "Trusted: yes"; then
    echo "GoPro already paired and trusted, skipping."
    exit 0
fi

SCRIPT_DIR="$(cd "$(dirname "$BASH_SOURCE")" && pwd)"

echo "Pairing..."
bluetoothctl --agent NoInputNoOutput pair "$MAC"
expect "$SCRIPT_DIR/pair_gopro.sh" "$MAC" 

sleep 5

echo "Trusting..."
bluetoothctl trust "$MAC"

echo ""
echo "Done. GoPro is now paired and trusted."
echo "Run this script only once. The API will connect automatically from now on." 
