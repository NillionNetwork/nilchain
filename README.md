# nilliond

nilliond is the Coordination Layer for the Nillion Network. It coordinates the payment of blind
computations and storage operations performed on the network. It is built using the [Cosmos
SDK](https://github.com/cosmos/cosmos-sdk), a framework for building PoS blockchain applications,
which was chosen for the inter-connectivity, speed and sovereignty its ecosystem provides.

## Building

```
make install
```

```
make build
```

### Scripts

```
cd scripts
```

```
# setup and start single chain locally
sh init.sh
# setup and start two chains locally and create an ibc client/connection/channel 
sh hermes.sh
# naive setup and start single chain for deploying devnet
sh testnet.sh
```
