#!/bin/bash

stringToHex() {
    echo -n "$1" | xxd -p | tr -d '\n' | sed 's/^/0x/'
}

input1='{"path":"createContract","payload":{"symbol":"STABLECOIN","address":"0x0000000000000000000000000000000000000002"}}'
input2='{"path":"createUser","payload":{"address":"0x70997970C51812dc3A010C7d01b50e0d17dc79C8","role":"creator","username":"vitalik"}}'

expires_at=$(($(date +%s) + 5))
input3=$(printf '{"path":"createAuction","payload":{"max_interest_rate":"10","expires_at":%d,"debt_issued":%d}}' "$expires_at" 100)

hexInput1=$(stringToHex "$input1")
hexInput2=$(stringToHex "$input2")
hexInput3=$(stringToHex "$input3")

cast send 0x59b22D57D4f067708AB0c00552767405926dc768 "addInput(address,bytes)(bytes32)" 0xab7528bb862fb57e8a2bcd567a2e929a0be56a5e $hexInput1 --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
cast send 0x59b22D57D4f067708AB0c00552767405926dc768 "addInput(address,bytes)(bytes32)" 0xab7528bb862fb57e8a2bcd567a2e929a0be56a5e $hexInput2 --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
cast send 0x59b22D57D4f067708AB0c00552767405926dc768 "addInput(address,bytes)(bytes32)" 0xab7528bb862fb57e8a2bcd567a2e929a0be56a5e $hexInput3 --private-key 0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d