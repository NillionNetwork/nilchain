#!/usr/bin/env bash

# Set the home directory for your chain
HOMEDIR="$HOME/.nilchainapp"

# Remove existing state to ensure a clean initialization
rm -rf "$HOMEDIR"

NILLIOND_BIN=$(which nilchaind)

# Initialize the chain
$NILLIOND_BIN init test --chain-id demo --default-denom unillion --home "$HOMEDIR"

# Configure other settings (chain ID, keyring-backend)
$NILLIOND_BIN config set client chain-id demo --home "$HOMEDIR"
$NILLIOND_BIN config set client keyring-backend test --home "$HOMEDIR"

# Add keys for users
$NILLIOND_BIN keys add alice --home "$HOMEDIR"
$NILLIOND_BIN keys add bob --home "$HOMEDIR"

# Add genesis accounts and create a default validator
$NILLIOND_BIN genesis add-genesis-account alice 10000000unillion --keyring-backend test --home "$HOMEDIR"
$NILLIOND_BIN genesis add-genesis-account bob 1000unillion --keyring-backend test --home "$HOMEDIR"

# Create a default validator and collect genesis transactions
$NILLIOND_BIN genesis gentx alice 1000000unillion --chain-id demo --home "$HOMEDIR"
$NILLIOND_BIN genesis collect-gentxs --home "$HOMEDIR"

# Start the chain
$NILLIOND_BIN start --home "$HOMEDIR"