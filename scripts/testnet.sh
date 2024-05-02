#!/usr/bin/env bash

# Set the home directory for your chain
HOMEDIR="$HOME/.nillionapp"
NODE_IP="0.0.0.0"
RPC_LADDR="$NODE_IP:26648"
GRPC_ADDR="$NODE_IP:9081"

# Remove existing state to ensure a clean initialization
rm -rf "$HOMEDIR"

NILLIOND_BIN=$(which nilliond)

# Initialize the chain
$NILLIOND_BIN init test --chain-id demo --default-denom anillion --home "$HOMEDIR"

# Configure other settings (chain ID, keyring-backend)
$NILLIOND_BIN config set client chain-id demo --home "$HOMEDIR"
$NILLIOND_BIN config set client keyring-backend test --home "$HOMEDIR"

# Add keys for users
$NILLIOND_BIN keys add alice --home "$HOMEDIR"
$NILLIOND_BIN keys add bob --home "$HOMEDIR"

# Add genesis accounts and create a default validator
$NILLIOND_BIN genesis add-genesis-account alice 10000000anillion --keyring-backend test --home "$HOMEDIR"
$NILLIOND_BIN genesis add-genesis-account bob 1000anillion --keyring-backend test --home "$HOMEDIR"

# Create a default validator and collect genesis transactions
$NILLIOND_BIN genesis gentx alice 1000000anillion --chain-id demo --home "$HOMEDIR"
$NILLIOND_BIN genesis collect-gentxs --home "$HOMEDIR"

# Start the chain
$NILLIOND_BIN start \
       --home "$HOMEDIR" \
       --rpc.laddr tcp://${RPC_LADDR} \
       --grpc.address ${GRPC_ADDR} \
       --address tcp://${NODE_IP}:26635 \
       --p2p.laddr tcp://${NODE_IP}:26636 \
       --grpc-web.enable=false \
       --log_level trace \
       --trace \
       &> $HOMEDIR/logs &