<div>
    <img src="https://avatars.githubusercontent.com/u/99333377?s=48&v=4" align="left" width="35" style="margin-right: 15px"/>
    <h1>
        Nillion Testnets
    </h1>
    <p> This repository contains information on Nillion public testnets and devnets </p>
    <br>
</div>

| Chain ID                              | Type      | Status     | Version       | Notes                     |
|---------------------------------------|-----------|------------|---------------|---------------------------|
| [nillion-testnet-1](./testnets/nillion-testnet-1) | `testnet` | **Inactive** | `v0.1.0` | Testnet                   |
| [devnet](./devnets/)            | `devnet`  | **Beta**   | `v0.1.0`     | Devnet with mainnet state |

## Testnets

Testnets are a type of blockchain used exclusively for testing purposes. They function as a sandbox environment, allowing developers to test new code and functionalities without worrying about affecting the live blockchain (mainnet). They are persistent environments, meaning that they remain active for extended periods of time.

Testnets come with a range of integrated services, including relayers to other testnets, frontends, explorers, and snapshot services.

### nillion-testnet-1

| Chain ID         | `nillion-testnet-1`                                      |
|------------------|----------------------------------------------------|
| Nilliond version | `v0.1.0`                                      |
| Genesis          | <http://genesis.testnet.nillion.com/genesis.json> |
| RPC              | <https://rpc.testnet.nillion.com>                 |
| gRPC             | <https://grpc.testnet.nillion.com>                |
| REST             | <https://lcd.testnet.nillion.com>                 |
| Faucet           | <https://faucet.testnet.nillion.com>              |
| Explorer         | <https://explorer.testnet.nillion.com>            |
| Snapshots        | <https://snapshots.testnet.nillion.com>           |
| Frontend         | <https://testnet.nillion.com>                     |

#### Join the testnet

Join the testnet following the instructions on the [nillion-testnet-1 page](./testnets/nillion-testnet-1/README.md).

## Devnets

Devnets, short for development networks, are also used for testing new functionalities and code. However, unlike testnets, devnets are temporary environments.

Devnets are ephemeral, which means they do not persist state long term and are recreated on an ongoing basis. This ensures that the testing environment closely mirrors the current state of the latest nillion codebase. Devnets are minimal environments, consisting only of a validator. Unlike testnets, devnets do not feature frontends or relayers to other testnets.

### devnet

| Chain ID         | `devnet`                                                                 |
|------------------|--------------------------------------------------------------------------|
| Nilliond version | `v0.1.0`                                                                |
| Genesis          | <https://nillion.com/devnet/genesis.json>        |
| Starting Height  | <https://nillion.com/devnet/height>              |
| RPC              | <https://rpc.devnet.nillion.com>                                        |
| REST             | <https://lcd.devnet.nillion.com>                                        |
| gRPC             | `grpc.devnet.nillion.com:30090`                                         |
| websocket        | `wss://rpc.devnet.nillion.com:443/websocket`                            |
| Faucet           | <https://faucet.devnet.nillion.com>                                     |
| Seed Node        | `5943e04edc5397018803f73e47f826be016010e1@p2p.devnet.nillion.com:30056` |

## 🆘 Issues and support

If you encounter any issues while joining the Nillion network or have questions about the process, please don't hesitate to reach out for support.

- For general questions and community support, join the [Nillion Discord](https://discord.com/) and ask in the `#testnet` channel.

- For technical issues or bugs related to the testnet, submit a detailed issue report on this repository with a clear description of the problem and any relevant error messages or logs.

## 🙋‍♀️ FAQ

**1) I need some funds on the `nillion-testnet-1` testnet, how can I get them?**

You can request testnet tokens for the `nillion-testnet-1` testnet from the faucet available at <https://faucet.testnet.nillion.com>. Please note that the faucet currently dispenses up to 500 NILLION per day per address.

**2) I am an integrator on the `nillion-testnet-1` testnet. How can I request more funds?**

If you are an integrator needing additional testnet tokens for development or testing purposes, you can request them via [this form](https://form-integrators.testnet.nillion.com). Please provide detailed information about your project and the number of tokens you require, and our team will review your request as soon as possible.

**3) What are the differences between testnets and devnets?**

| **Features** | **Testnets**             | **Devnets**                                      |
|--------------|--------------------------|--------------------------------------------------|
| Persistent   | ✅                        | ❌  |
| State        | Maintains its own state. | Forks of the mainnet, mimicking its state.       |
| Faucet       | ✅                        | ✅                                                |
| Explorer     | ✅                        | ❌                                                |
| Frontend     | ✅                        | ❌                                                |
| Relayers     | ✅                        | ❌                                                |
