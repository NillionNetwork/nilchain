#!/usr/bin/env bash

# Set the home directory for your chain
NODE_IP="localhost"
HOMEDIR="$HOME/.nillionapp"
HOMEDIR1="$HOME/.nillionapp1"
RPC_LADDR="$NODE_IP:26648"
GRPC_ADDR="$NODE_IP:9081"
RPC_LADDR1="$NODE_IP:26649"
GRPC_ADDR1="$NODE_IP:9082"

# Remove existing state to ensure a clean initialization
rm -rf "$HOMEDIR"
rm -rf "$HOMEDIR1"

NILLIOND_BIN=$(which nilchaind)

# Initialize the chain
$NILLIOND_BIN init test --chain-id chainA --default-denom anillion --home "$HOMEDIR"
$NILLIOND_BIN init test --chain-id chainB --default-denom anillion --home "$HOMEDIR1"

# Configure other settings (chain ID, keyring-backend)
$NILLIOND_BIN config set client chain-id chainA --home "$HOMEDIR"
$NILLIOND_BIN config set client keyring-backend test --home "$HOMEDIR"

$NILLIOND_BIN config set client chain-id chainB --home "$HOMEDIR1"
$NILLIOND_BIN config set client keyring-backend test --home "$HOMEDIR1"

# Add keys for users
$NILLIOND_BIN keys add alice --home "$HOMEDIR"
$NILLIOND_BIN keys add bob --home "$HOMEDIR"
$NILLIOND_BIN keys add alice --home "$HOMEDIR1"
$NILLIOND_BIN keys add bob --home "$HOMEDIR1"

# Add genesis accounts and create a default validator
$NILLIOND_BIN genesis add-genesis-account alice 10000000anillion --keyring-backend test --home "$HOMEDIR"
$NILLIOND_BIN genesis add-genesis-account bob 1000anillion --keyring-backend test --home "$HOMEDIR"

# Add genesis accounts and create a default validator
$NILLIOND_BIN genesis add-genesis-account alice 10000000anillion --keyring-backend test --home "$HOMEDIR1"
$NILLIOND_BIN genesis add-genesis-account bob 1000anillion --keyring-backend test --home "$HOMEDIR1"

# Create a default validator and collect genesis transactions
$NILLIOND_BIN genesis gentx alice 1000000anillion --chain-id chainA --home "$HOMEDIR"
$NILLIOND_BIN genesis collect-gentxs --home "$HOMEDIR"

# Create a default validator and collect genesis transactions
$NILLIOND_BIN genesis gentx alice 1000000anillion --chain-id chainB --home "$HOMEDIR1"
$NILLIOND_BIN genesis collect-gentxs --home "$HOMEDIR1"

# Add account in genesis (required by Hermes)
# Create user account keypair
$NILLIOND_BIN keys add test --keyring-backend test --home $HOMEDIR --output json > $HOMEDIR/keypair.json 2>&1
$NILLIOND_BIN genesis add-genesis-account $(jq -r .address $HOMEDIR/keypair.json) 1000000000stake --home $HOMEDIR
$NILLIOND_BIN keys add test2 --keyring-backend test --home $HOMEDIR1 --output json > $HOMEDIR1/keypair.json 2>&1
$NILLIOND_BIN genesis add-genesis-account $(jq -r .address $HOMEDIR1/keypair.json) 1000000000stake --home $HOMEDIR1

# Start the chain
$NILLIOND_BIN start \
       --home "$HOMEDIR" \
       --rpc.laddr tcp://${RPC_LADDR} \
       --grpc.address ${GRPC_ADDR} \
       --address tcp://${NODE_IP}:26635 \
       --p2p.laddr tcp://${NODE_IP}:26636 \
       --grpc.enable=true \
       --grpc-web.enable=true \
       --log_level trace \
       --trace \
       &> $HOMEDIR/logs &

# Start the chain
$NILLIOND_BIN start \
       --home "$HOMEDIR1" \
       --rpc.laddr tcp://${RPC_LADDR1} \
       --grpc.address ${GRPC_ADDR1} \
       --address tcp://${NODE_IP}:26634 \
       --grpc.enable=true \
       --p2p.laddr tcp://${NODE_IP}:26635 \
       --grpc-web.enable=true \
       --log_level trace \
       --trace \
       &> $HOMEDIR1/logs &

sleep 10

######################################HERMES###################################

# Setup Hermes in packet relayer mode
killall hermes 2> /dev/null || true
rm -rf ~/.hermes/

mkdir ~/.hermes/
touch ~/.hermes/config.toml
tee ~/.hermes/config.toml<<EOF
[global]
log_level = "trace"

[mode]

[mode.clients]
enabled = true
refresh = true
misbehaviour = true

[mode.connections]
enabled = true

[mode.channels]
enabled = true

[mode.packets]
enabled = true
clear_interval = 100
clear_on_start = true
tx_confirmation = true
auto_register_counterparty_payee = false


[[chains]]
account_prefix = "nillion"
clock_drift = "5s"
gas_multiplier = 1.1
grpc_addr = "tcp://${GRPC_ADDR}"
id = "chainA"
key_name = "relayer"
max_gas = 2000000
rpc_addr = "http://${RPC_LADDR}"
rpc_timeout = "10s"
store_prefix = "ibc"
trusting_period = "599s"
event_source = { mode = 'push', url = 'ws://${RPC_LADDR}/websocket', batch_delay = '500ms' }

[chains.gas_price]
       denom = "stake"
       price = 0.00

[chains.trust_threshold]
       denominator = "3"
       numerator = "1"

[[chains]]
account_prefix = "nillion"
clock_drift = "5s"
gas_multiplier = 1.1
grpc_addr = "tcp://${GRPC_ADDR1}"
id = "chainB"
key_name = "relayer"
max_gas = 2000000
rpc_addr = "http://${RPC_LADDR1}"
rpc_timeout = "10s"
store_prefix = "ibc"
trusting_period = "599s"
event_source = { mode = 'push', url = 'ws://${RPC_LADDR1}/websocket', batch_delay = '500ms' }

[chains.gas_price]
       denom = "stake"
       price = 0.00

[chains.trust_threshold]
       denominator = "3"
       numerator = "1"
EOF

# Delete all previous keys in relayer
hermes keys delete --chain chainA --all
hermes keys delete --chain chainB --all

# Restore keys to hermes relayer
jq -r .mnemonic $HOMEDIR/keypair.json | hermes keys add --chain chainA --mnemonic-file /dev/stdin
jq -r .mnemonic $HOMEDIR1/keypair.json | hermes keys add --chain chainB --mnemonic-file /dev/stdin

sleep 5

hermes create channel --a-chain chainA --b-chain chainB --a-port transfer --b-port transfer --new-client-connection

sleep 7

hermes -j start &> ~/.hermes/logs &