#!/bin/bash

# Trap SIGINT (Ctrl+C) to gracefully terminate all child processes
cleanup() {
    echo "Terminating child processes..."
    kill 0 
    exit 0
}
trap cleanup SIGINT

# Deploy tokens using Forge
forge script ./contracts/script/Token.s.sol \
    --private-key ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 \
    --rpc-url http://localhost:8545 \
    --broadcast \
    --root ./contracts \
    -vv &  # Inicia em segundo plano

# Helper function to convert string to hex
stringToHex() {
    echo -n "$1" | xxd -p | tr -d '\n' | sed 's/^/0x/'
}

# Ethereum addresses
INPUT_BOX="0x59b22D57D4f067708AB0c00552767405926dc768"
DAPP_ADDRESS="0xab7528bb862fb57e8a2bcd567a2e929a0be56a5e"
PORTAL_ADDRESS="0x9C21AEb2093C32DDbC53eEF24B873BDCd1aDa1DB"
STABLECOIN_ADDRESS="0x368B8A7D8A2247489582CC83b502d0A9A185E4E9"
TOKENIZED_RECEIVABLE_ADDRESS="0x8C3dADb62dec908515049dE39A52828681cc4912"
ADMIN_ADDRESS="0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
ADMIN_PRIVATE_KEY="0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
CREATOR_ADDRESS="0x70997970C51812dc3A010C7d01b50e0d17dc79C8"
CREATOR_PRIVATE_KEY="0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d"

DAPP_ADDRESS_RELAY="0xF5DE34d6BbC0446E2a45719E718efEbaaE179daE"

INVESTOR_ADDRESSES=(
    "0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC"
    "0x90F79bf6EB2c4f870365E785982E1f101E93b906"
    "0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65"
    "0x9965507D1a55bcC2695C58ba16FB37d819B0A4dc"
    "0x976EA74026E726554dB657fA54763abd0C3a0aa9"
)

INVESTOR_PRIVATE_KEYS=(
    "0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a"
    "0x7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6"
    "0x47e179ec197488593b187f80a00eb0da91f1b9d0b13f8733639f19c30a34926a"
    "0x8b3a350cf5c34c9194ca85829a2df0ec3153be0318b5e2d3348e872092edffba"
    "0x92db14e403b83dfe3df233f83dfa3a0d7096f21ca9b0d6d6b8d88b2b4ec1564e"
)

# Function to relay DApp address to InputBox
relayDAppAddress() {
    local dappAddress="$1"
    local senderPrivateKey="$2"
    cast send $DAPP_ADDRESS_RELAY "relayDAppAddress(address)" $dappAddress --private-key $senderPrivateKey --rpc-url http://localhost:8545
}

# Function to send an input to the INPUT_BOX contract
sendInput() {
    local payload="$1"
    local senderPrivateKey="$2"
    local hexPayload
    hexPayload=$(stringToHex "$payload")
    cast send $INPUT_BOX "addInput(address,bytes)(bytes32)" $DAPP_ADDRESS $hexPayload --private-key $senderPrivateKey --rpc-url http://localhost:8545
}

# Function to approve ERC20 tokens
approveTokens() {
    local token="$1"
    local spender="$2"
    local amount="$3"
    local privateKey="$4"
    cast send $token "approve(address,uint256)" $spender $amount --private-key $privateKey --rpc-url http://localhost:8545
}

# Function to deposit ERC20 tokens
depositTokens() {
    local token="$1"
    local dapp="$2"
    local amount="$3"
    local execLayerData="$4"
    local privateKey="$5"
    cast send $PORTAL_ADDRESS "depositERC20Tokens(address,address,uint256,bytes)" $token $dapp $amount $execLayerData --private-key $privateKey --rpc-url http://localhost:8545
}

# 0. Relay the DApp address (required before sending inputs)
echo "Relaying DApp address..."
relayDAppAddress $DAPP_ADDRESS $ADMIN_PRIVATE_KEY &

# 1. Create users (sent by admin)
echo "Creating users..."
sendInput '{"path":"createUser","payload":{"address":"'"$CREATOR_ADDRESS"'","role":"creator"}}' $ADMIN_PRIVATE_KEY &
for i in "${!INVESTOR_ADDRESSES[@]}"; do
    role="non_qualified_investor"
    [ "$i" -lt 2 ] && role="qualified_investor"
    sendInput '{"path":"createUser","payload":{"address":"'"${INVESTOR_ADDRESSES[$i]}"'","role":"'"$role"'"}}' $ADMIN_PRIVATE_KEY &
    sleep 0.5
done
wait 
# 2. Create contracts
echo "Creating contracts..."
createContractPayload='{"path":"createContract","payload":{"symbol":"STABLECOIN","address":"'"$STABLECOIN_ADDRESS"'"}}'
sendInput "$createContractPayload" $ADMIN_PRIVATE_KEY &

createContractPayload='{"path":"createContract","payload":{"symbol":"TOKENIZED_RECEIVABLE","address":"'"$TOKENIZED_RECEIVABLE_ADDRESS"'"}}'
sendInput "$createContractPayload" $ADMIN_PRIVATE_KEY &
wait

# 3. Create crowdfunding (sent by creator)
echo "Creating crowdfunding..."
current_timestamp=$(date +%s)
expires_at=$((current_timestamp + 60))  # Expiração em 60 segundos
maturity_at=$((current_timestamp + 120)) # Maturidade em 120 segundos

crowdfundingPayload='{"path":"createCrowdfunding","payload":{"max_interest_rate":"10","debt_issued":"100000","expires_at":'"$expires_at"',"maturity_at":'"$maturity_at"'}}'
approveTokens $TOKENIZED_RECEIVABLE_ADDRESS $PORTAL_ADDRESS 10000 $CREATOR_PRIVATE_KEY &
depositTokens $TOKENIZED_RECEIVABLE_ADDRESS $DAPP_ADDRESS 10000 "$(stringToHex "$crowdfundingPayload")" $CREATOR_PRIVATE_KEY &
wait

# 4. Update crowdfunding to ongoing (sent by admin)
echo "Updating crowdfunding state to 'ongoing'..."
updatePayload='{"path":"updateCrowdfunding","payload":{"id":1,"state":"ongoing"}}'
sendInput "$updatePayload" $ADMIN_PRIVATE_KEY &
wait

# 5. Create orders from investors (sent by each investor)
echo "Creating orders from investors..."
ORDER_AMOUNTS=(60000 52000 2000 3000 400)
INTEREST_RATES=("9" "8" "4" "6" "4")
for i in "${!INVESTOR_ADDRESSES[@]}"; do
    orderPayload='{"path":"createOrder","payload":{"creator":"'"$CREATOR_ADDRESS"'","interest_rate":"'"${INTEREST_RATES[$i]}"'"}}'
    approveTokens $STABLECOIN_ADDRESS $PORTAL_ADDRESS "${ORDER_AMOUNTS[$i]}" "${INVESTOR_PRIVATE_KEYS[$i]}" &
    depositTokens $STABLECOIN_ADDRESS $DAPP_ADDRESS "${ORDER_AMOUNTS[$i]}" "$(stringToHex "$orderPayload")" "${INVESTOR_PRIVATE_KEYS[$i]}" &
    sleep 0.5
done
wait
# 6. Wait for crowdfunding expiration
echo "Waiting for crowdfunding to expire..."
sleep 60

# 7. Close crowdfunding (sent by admin)
echo "Closing crowdfunding..."
closePayload='{"path":"closeCrowdfunding","payload":{"creator":"'"$CREATOR_ADDRESS"'"}}'
sendInput "$closePayload" $ADMIN_PRIVATE_KEY &
wait

# 8. Wait for crowdfunding maturity
echo "Waiting for crowdfunding maturity..."
sleep 60

# 9. Settle crowdfunding (sent by creator using stablecoin)
echo "Settling crowdfunding..."
settlePayload='{"path":"settleCrowdfunding","payload":{"crowdfunding_id":1}}'
approveTokens $STABLECOIN_ADDRESS $PORTAL_ADDRESS 108600 $CREATOR_PRIVATE_KEY &
depositTokens $STABLECOIN_ADDRESS $DAPP_ADDRESS 108600 "$(stringToHex "$settlePayload")" $CREATOR_PRIVATE_KEY &
wait

echo "All transactions completed successfully!"
