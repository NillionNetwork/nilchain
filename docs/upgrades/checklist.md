## Upgrade Checklist Nillion

### Release structure

- Tags/Versioning
  - <https://semver.org/>

**MAJOR** version when you make incompatible API changes and/or state machine breaking changes
**MINOR** version when you add functionality in a backward compatible manner/non state machine breaking changes
**PATCH** version when you make backward compatible bug fixes

### Engineering checklist pre-upgrade

- [ ] Run IBC transfer e2e in CI on the latest tag
- [ ] Write an upgrade test using the example test. This allows you to test the upgrade between two provided tags and run the test in CI easily.
- [ ] Locally test upgrade with mainnet state (<https://docs.cosmos.network/v0.50/build/building-apps/app-testnet>)
- [ ] Manual QA of all newly added functionality
- [ ] Run all previous e2e tests and unit tests locally/in CI
- [ ] Provide sufficient release notes for validators/users/developers
- [ ] If migrations are necessary, implement them and write tests.

### Logistics/Comms checklist for handling upgrade

- [ ] Settle on the target upgrade height, use estimation based on current block time
- [ ] Create Governance proposal with upgrade height, version etc. (chain will halt)
  - Example

```json
{
    "title": "Update nillion to v0.2.1",
    "description": "Update to v0.2.1",
    "summary": "Update binary to version v0.2.1",
    "messages": [
        {
            "@type": "/cosmos.upgrade.v1beta1.MsgSoftwareUpgrade",
            "authority": "nillion10d07y265gmmuvt4z0w9aw880jnsr700jpzdkas",
            "plan": {
                "name": "v0.2.1",
                "time": "0001-01-01T00:00:00Z",
                "height": "1100000",
                "info": "{}",
                "upgraded_client_state": null
            }
        }
    ],
    "deposit": "5000000000unil",
    "expedited": true
}
```

- [ ] Make sure key stakeholders (validators etc. vote)
- [ ] Communicate with Validators and specify tag or release to run for upgrade

**Guides**:

- <https://medium.com/web3-surfers/cosmos-dev-series-cosmos-sdk-based-blockchain-upgrade-b5e99181554c>
- <https://docs.cosmos.network/main/learn/advanced/upgrade>
