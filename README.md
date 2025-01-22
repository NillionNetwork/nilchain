# nilchain

nilchain is the Coordination Layer for the Nillion Network. It coordinates the payment of blind
computations and storage operations performed on the network. It is built using the [Cosmos
SDK](https://github.com/cosmos/cosmos-sdk), a framework for building PoS blockchain applications,
which was chosen for the interconnectivity, speed, and sovereignty its ecosystem provides.

## Building

```
make install
```

```
make build
```

### Cross-Compiling

Use the `build-cross` target to cross-compile for the following platforms: linux/amd64, linux/arm64,
darwin/amd64, darwin/arm64. Please note that building for linux/arm64 requires an arm64-compatible
version of gcc installed. On Debian-based systems, this would be `aarch64-linux-gnu-gcc`.

### Scripts

```
cd scripts
```

```
# Setup and start single chain locally
sh init.sh
# Setup and start two chains locally and create an ibc client/connection/channel 
sh hermes.sh
# Naive setup and start single chain for deploying devnet
sh testnet.sh
```
