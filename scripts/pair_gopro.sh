#!/usr/bin/expect -f

# GoPro MAC address used as argument
set mac_address [lindex $argv 0]

# Bluetooth start
spawn bluetoothctl

# Start pairing
expect "#"
send "pair $mac_address\r"

# Wait for success confirmation
expect "Pairing successful"

#Exits program
send "exit\r"

expect eof