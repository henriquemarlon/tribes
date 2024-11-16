#!/bin/bash

stringToHex() {
    echo -n "$1" | xxd -p | tr -d '\n' | sed 's/^/0x/'
}

input1='{"path":"echo","payload":{"message":"Hello, Rollmelette!"}}'
hexInput1=$(stringToHex "$input1")

cast send 0x58Df21fE097d4bE5dCf61e01d9ea3f6B81c2E1dB "addInput(address,bytes)(bytes32)" 0xab7528bb862fb57e8a2bcd567a2e929a0be56a5e $hexInput1 --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80