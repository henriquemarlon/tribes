#!/bin/bash

# Trap SIGINT (Ctrl+C) to gracefully terminate all child processes
cleanup() {
    echo "Terminating child processes..."
    kill 0 
    exit 0
}
trap cleanup SIGINT

# Helper function to convert string to hex
stringToHex() {
    echo -n "$1" | xxd -p | tr -d '\n' | sed 's/^/0x/'
}

# Ethereum addresses
INPUT_BOX="0x59b22D57D4f067708AB0c00552767405926dc768"
DAPP_ADDRESS="0xab7528bb862fb57e8a2bcd567a2e929a0be56a5e"
PORTAL_ADDRESS="0x9C21AEb2093C32DDbC53eEF24B873BDCd1aDa1DB"
ADMIN_ADDRESS="0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
ADMIN_PRIVATE_KEY="0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
CREATOR_ADDRESS="0x70997970C51812dc3A010C7d01b50e0d17dc79C8"
CREATOR_PRIVATE_KEY="0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d"

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

# Deploy the Token contract and capture the deployed address
deployToken() {
    local tokenName="$1"
    local tokenSymbol="$2"    
    # Execute the forge create command and capture the output
    result=$(forge create ./src/Token.sol:Token \
        --private-key $ADMIN_PRIVATE_KEY \
        --rpc-url http://localhost:8545 \
        --root ./contracts \
        --constructor-args "$tokenName" "$tokenSymbol" 2>&1)
    
    # Extract the deployed address from the output
    deployedAddress=$(echo "$result" | grep "Deployed to:" | awk '{print $3}')
    
    # Check if the deployed address was extracted
    if [[ -z "$deployedAddress" ]]; then
        echo "Error: Failed to deploy contract for $tokenName ($tokenSymbol)."
        echo "$result"
        exit 1
    fi
    echo "$deployedAddress"
}

# Relay the DApp address
relayDAppAddress() {
    local dappAddress="$1"
    cast send 0xF5DE34d6BbC0446E2a45719E718efEbaaE179daE "relayDAppAddress(address)" $dappAddress --private-key $ADMIN_PRIVATE_KEY --rpc-url http://localhost:8545
}

# Send input to the INPUT_BOX contract
sendInput() {
    local payload="$1"
    hexPayload=$(stringToHex "$payload")
    cast send $INPUT_BOX "addInput(address,bytes)(bytes32)" $DAPP_ADDRESS $hexPayload --private-key $ADMIN_PRIVATE_KEY --rpc-url http://localhost:8545
}

# Mint tokens to a specified address
mintTokens() {
    local tokenAddress="$1"
    local recipient="$2"
    local amount="$3"
    cast send $tokenAddress "mint(address,uint256)" $recipient $amount --private-key $ADMIN_PRIVATE_KEY --rpc-url http://localhost:8545
    echo "Minted $amount tokens to $recipient on $tokenAddress"
}

# Approve ERC20 tokens
approveTokens() {
    local token="$1"
    local spender="$2"
    local amount="$3"
    local privateKey="$4"
    echo "Approving $amount tokens for spender ($spender)..."
    cast send $token \
        "approve(address,uint256)" \
        $spender $amount \
        --private-key $privateKey \
        --rpc-url http://localhost:8545
}

# Function to deposit ERC20 tokens
depositERC20Tokens() {
    local token="$1"
    local dapp="$2"
    local amount="$3"
    local execLayerData="$4"
    local privateKey="$5"
    echo "Depositing $amount of token ($token) to DApp ($dapp)..."
    cast send $PORTAL_ADDRESS \
        "depositERC20Tokens(address,address,uint256,bytes)" \
        $token $dapp $amount "$(stringToHex "$execLayerData")" \
        --private-key $privateKey \
        --rpc-url http://localhost:8545
}

echo "Relaying DApp address..."
relayDAppAddress $DAPP_ADDRESS
sleep 1

echo "Deploying contracts..."
STABLECOIN_ADDRESS=$(deployToken "Stablecoin" "STABLECOIN")
sleep 1

TOKENIZED_RECEIVABLE_ADDRESS=$(deployToken "Pink" "PINK")
sleep 1

echo "Deployed contracts:"
echo "STABLECOIN_ADDRESS=$STABLECOIN_ADDRESS"
echo "TOKENIZED_RECEIVABLE_ADDRESS=$TOKENIZED_RECEIVABLE_ADDRESS"

echo "Minting tokens to investors and creator..."
mintTokens $TOKENIZED_RECEIVABLE_ADDRESS $CREATOR_ADDRESS 10000000
sleep 1

mintTokens $STABLECOIN_ADDRESS $CREATOR_ADDRESS 10000000
sleep 1

for investor in "${INVESTOR_ADDRESSES[@]}"; do
    mintTokens $STABLECOIN_ADDRESS $investor 10000000
    sleep 1
done

# Create contracts
echo "Creating contracts..."
sendInput '{"path":"createContract","payload":{"symbol":"STABLECOIN","address":"'"$STABLECOIN_ADDRESS"'"}}'
sleep 1
sendInput '{"path":"createContract","payload":{"symbol":"TOKENIZED_RECEIVABLE","address":"'"$TOKENIZED_RECEIVABLE_ADDRESS"'"}}'
sleep 1

# Create users
echo "Creating users..."
sendInput '{"path":"createUser","payload":{"address":"'"$CREATOR_ADDRESS"'","role":"creator"}}'
sleep 1
for i in "${!INVESTOR_ADDRESSES[@]}"; do
    role="non_qualified_investor"
    [ "$i" -lt 2 ] && role="qualified_investor"
    sendInput '{"path":"createUser","payload":{"address":"'"${INVESTOR_ADDRESSES[$i]}"'","role":"'"$role"'"}}'
    sleep 1
done

# Create crowdfunding
echo "Creating crowdfunding..."
current_timestamp=$(date +%s)
closes_at=$((current_timestamp + 60))
maturity_at=$((current_timestamp + 120))
crowdfundingPayload='{"path":"createCrowdfunding","payload":{"max_interest_rate":"10","debt_issued":"100000","closes_at":'"$closes_at"',"maturity_at":'"$maturity_at"'}}'
approveTokens $TOKENIZED_RECEIVABLE_ADDRESS $PORTAL_ADDRESS 10000 $CREATOR_PRIVATE_KEY
sleep 1
depositERC20Tokens $TOKENIZED_RECEIVABLE_ADDRESS $DAPP_ADDRESS 10000 "$crowdfundingPayload" $CREATOR_PRIVATE_KEY
sleep 1
dfundingPayload='{"path":"createCrowdfunding","payload":{"max_interest_rate":"10","debt_issued":"100000","closes_at":'"$closes_at"',"maturity_at":'"$maturity_at"'}}'


# 4. Update crowdfunding to ongoing (sent by admin)
echo "Updating crowdfunding state to 'ongoing'..."
updatePayload='{"path":"updateCrowdfunding","payload":{"id":1,"state":"ongoing"}}'
sendInput "$updatePayload" $ADMIN_PRIVATE_KEY &
sleep 1
wait

# 5. Create orders from investors (sent by each investor)
echo "Creating orders from investors..."
ORDER_AMOUNTS=(60000 52000 2000 3000 400)
INTEREST_RATES=("9" "8" "4" "6" "4")
for i in "${!INVESTOR_ADDRESSES[@]}"; do
    orderPayload='{"path":"createOrder","payload":{"creator":"'"$CREATOR_ADDRESS"'","interest_rate":"'"${INTEREST_RATES[$i]}"'"}}'
    approveTokens $STABLECOIN_ADDRESS $PORTAL_ADDRESS "${ORDER_AMOUNTS[$i]}" "${INVESTOR_PRIVATE_KEYS[$i]}" &
    sleep 1
    depositERC20Tokens $STABLECOIN_ADDRESS $DAPP_ADDRESS "${ORDER_AMOUNTS[$i]}" "$orderPayload" "${INVESTOR_PRIVATE_KEYS[$i]}" &
    sleep 1
done
sleep 1
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
sleep 30

# 9. Settle crowdfunding (sent by creator using stablecoin)
echo "Settling crowdfunding..."
settlePayload='{"path":"settleCrowdfunding","payload":{"crowdfunding_id":1}}'
approveTokens $STABLECOIN_ADDRESS $PORTAL_ADDRESS 108270 $CREATOR_PRIVATE_KEY &
sleep 1
depositERC20Tokens $STABLECOIN_ADDRESS $DAPP_ADDRESS 108270 "$settlePayload" $CREATOR_PRIVATE_KEY &
wait

# echo "All transactions completed successfully!"
