#!/bin/bash

stringToHex() {
    echo -n "$1" | xxd -p | tr -d '\n' | sed 's/^/0x/'
}

input1='{"path":"createContract","payload":{"symbol":"VOLT","address":"0x0000000000000000000000000000000000000001"}}'
input2='{"path":"createContract","payload":{"symbol":"STABLECOIN","address":"0x0000000000000000000000000000000000000002"}}'
input3='{"path":"createStation","payload":{"owner":"stationOwner01","consumption":100,"interest_rate":50,"latitude":40.7128,"longitude":-74.0060}}'
input4='{"path":"createStation","payload":{"owner":"stationOwner02","consumption":100,"interest_rate":10,"latitude":40.7128,"longitude":-74.0060}}'

input5='{"path":"createOrder","payload":{"station_id":1}}'
input6='{"path":"createOrder","payload":{"station_id":2}}'
input7='{"path":"createOrder","payload":{"station_id":1}}'
input8='{"path":"createOrder","payload":{"station_id":2}}'
input9='{"path":"createOrder","payload":{"station_id":1}}'

expires_at=$(($(date +%s) + 5))
input10=$(printf '{"path":"createAuction","payload":{"interest_rate":"1000","expires_at":%d,"debt_issued":%d}}' "$expires_at" 2)

input11='{"path":"createBid","payload":{"interest_rate":"100"}}'
input12='{"path":"createBid","payload":{"interest_rate":"500"}}'
input13='{"path":"createBid","payload":{"interest_rate":"200"}}'
input14='{"path":"createBid","payload":{"interest_rate":"300"}}'
input15='{"path":"createBid","payload":{"interest_rate":"200"}}'

hexInput1=$(stringToHex "$input1")
hexInput2=$(stringToHex "$input2")
hexInput3=$(stringToHex "$input3")
hexInput4=$(stringToHex "$input4")
hexInput5=$(stringToHex "$input5")
hexInput6=$(stringToHex "$input6")
hexInput7=$(stringToHex "$input7")
hexInput8=$(stringToHex "$input8")
hexInput9=$(stringToHex "$input9")
hexInput10=$(stringToHex "$input10")
hexInput11=$(stringToHex "$input11")
hexInput12=$(stringToHex "$input12")
hexInput13=$(stringToHex "$input13")
hexInput14=$(stringToHex "$input14")
hexInput15=$(stringToHex "$input15")

cast send 0x59b22D57D4f067708AB0c00552767405926dc768 "addInput(address,bytes)(bytes32)" 0xab7528bb862fb57e8a2bcd567a2e929a0be56a5e $hexInput1 --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
cast send 0x59b22D57D4f067708AB0c00552767405926dc768 "addInput(address,bytes)(bytes32)" 0xab7528bb862fb57e8a2bcd567a2e929a0be56a5e $hexInput2 --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
cast send 0x59b22D57D4f067708AB0c00552767405926dc768 "addInput(address,bytes)(bytes32)" 0xab7528bb862fb57e8a2bcd567a2e929a0be56a5e $hexInput3 --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
cast send 0x59b22D57D4f067708AB0c00552767405926dc768 "addInput(address,bytes)(bytes32)" 0xab7528bb862fb57e8a2bcd567a2e929a0be56a5e $hexInput4 --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80

cast send 0x92C6bcA388E99d6B304f1Af3c3Cd749Ff0b591e2 "approve(address,uint256)" 0x9C21AEb2093C32DDbC53eEF24B873BDCd1aDa1DB 9000ether --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
cast send 0x59b22D57D4f067708AB0c00552767405926dc768 "addInput(address,bytes)(bytes32)" 0xab7528bb862fb57e8a2bcd567a2e929a0be56a5e $hexInput5 --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
cast send 0x59b22D57D4f067708AB0c00552767405926dc768 "addInput(address,bytes)(bytes32)" 0xab7528bb862fb57e8a2bcd567a2e929a0be56a5e $hexInput6 --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
cast send 0x59b22D57D4f067708AB0c00552767405926dc768 "addInput(address,bytes)(bytes32)" 0xab7528bb862fb57e8a2bcd567a2e929a0be56a5e $hexInput7 --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
cast send 0x59b22D57D4f067708AB0c00552767405926dc768 "addInput(address,bytes)(bytes32)" 0xab7528bb862fb57e8a2bcd567a2e929a0be56a5e $hexInput8 --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
cast send 0x59b22D57D4f067708AB0c00552767405926dc768 "addInput(address,bytes)(bytes32)" 0xab7528bb862fb57e8a2bcd567a2e929a0be56a5e $hexInput9 --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80

cast send 0x59b22D57D4f067708AB0c00552767405926dc768 "addInput(address,bytes)(bytes32)" 0xab7528bb862fb57e8a2bcd567a2e929a0be56a5e $hexInput10 --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80

cast send 0x59b22D57D4f067708AB0c00552767405926dc768 "addInput(address,bytes)(bytes32)" 0xab7528bb862fb57e8a2bcd567a2e929a0be56a5e $hexInput11 --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
cast send 0x59b22D57D4f067708AB0c00552767405926dc768 "addInput(address,bytes)(bytes32)" 0xab7528bb862fb57e8a2bcd567a2e929a0be56a5e $hexInput12 --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
cast send 0x59b22D57D4f067708AB0c00552767405926dc768 "addInput(address,bytes)(bytes32)" 0xab7528bb862fb57e8a2bcd567a2e929a0be56a5e $hexInput13 --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
cast send 0x59b22D57D4f067708AB0c00552767405926dc768 "addInput(address,bytes)(bytes32)" 0xab7528bb862fb57e8a2bcd567a2e929a0be56a5e $hexInput14 --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
cast send 0x59b22D57D4f067708AB0c00552767405926dc768 "addInput(address,bytes)(bytes32)" 0xab7528bb862fb57e8a2bcd567a2e929a0be56a5e $hexInput15 --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80