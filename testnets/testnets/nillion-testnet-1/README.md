# nillion-testnet-1 testnet

## Details

| Chain ID          | `nillion-testnet-1`                                                                   |
|-------------------|---------------------------------------------------------------------------------|
| Launch Date       | 20th May 2024                                                                  |
| Nilliond version  | `v0.1.0`                                                                   |
| Genesis           | <https://github.com/NillionNetwork/networks/raw/main/nillion-testnet-1/genesis.tar.bz2> |
| Genesis SHA256SUM | `7f84b50c0ad65582c9c52126cd33443c9f2541436ea4c525106ed9b58f7c9ef4`              |

## Endpoints

Summary of the `nillion-testnet-1` endpoints:

| Service     | Url                                                                                                                                    |
|-------------|----------------------------------------------------------------------------------------------------------------------------------------|
| Public RPC  | <https://rpc.nilliontest1.nillion.com>                                                                                                   |
| Public LCD  | <https://grpc.nilliontest1.nillion.com>                                                                                                  |
| Public gRPC | <https://lcd.nilliontest1.nillion.com>                                                                                                   |
| Faucet      | <https://facet.nilliontest1.nillion.com>                                                                                                 |
| Seed Node   | `1f9a9c694c46bd28ad9ad6126e923993fc6c56b1@137.184.121.105:26656`                                                                       |
| Peers       | `3ab030b7fd75ed895c48bcc899b99c17a396736b@137.184.190.127:26656` <br/> `3dbffa30baab16cc8597df02945dcee0aa0a4581@143.198.139.33:26656` |
| Explorer    | <https://testnet.mintscan.io/nillion-testnet>                                                                          |
| Frontend    |     <https://app.nillion.com>                                                                                                |

### Public Nodes

| Protocol | Url                                 |
|----------|-------------------------------------|
| RPC      | <https://rpc.testnet1.nillion.com>  |
| gRPC     | <https://grpc.testnet1.nillion.com> |
| REST     | <https://lcd.testnet1.nillion.com>  |

### 🌱 Seed

| Node | ID                                                               |
|------|------------------------------------------------------------------|
| Seed | `0g2a9c694c46bd28ad9ad6126e923993fc6c56b1@117.183.181.105:26656` |

Add the Node ID in your `p2p.seeds` section of you `config.toml`:

```toml
#######################################################
###           P2P Configuration Options             ###
#######################################################
[p2p]

# ...

# Comma separated list of seed nodes to connect to
seeds = "0g9a9c694c46bd28ad9ad6126e923993fc6c56b2@132.182.181.105:26656"
```

### 🚰 Faucet

The `nillion-testnet-1` testnet faucet is available at <https://faucet.testnet1.nillion>

## Join the network

Before joining the network, ensure that your system meets the following minimum requirements:

- 4 CPU
- 8GB RAM
- 100GB free disk space

To join the Nillion network, there are various options available. Choose the option that works best for you based on your preferences and technical requirements.

1. (Manual) Manual setup

### Option: Manual setup

1. Download the `nilliond` binary:
2. Initialize the node:
3. Download genesis
4. Download the latest snapshot:
5. Set the seed node in the `config.toml`:
6. Start the node:
