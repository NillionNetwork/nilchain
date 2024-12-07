# Meta

## Abstract

Meta is a module that holds metadata related to resources paid for in the PET network.

## State

### Data

MetaData is a map of information associated with resource consumption in the PET network. The data which is stored will not be verified. The module has no concept of what is valid or invalid with the associated data.

* Data: `0x01 | sha256(data) (32 byte) | bytes(data) |`

## Messages

Meta has a single transaction type. It is used to pay for resources in the PET network. All fields are required. The `from_address` is the address of the account that is paying for the resources. The `amount` is the amount of resources being paid for. The `resource` is the type of resource being paid for.

The amount is burned instead of sent to another address at this point in time. In the future this will be changed to be sent to a module account in which the funds are pooled and paid out to workers in the PET network.

```protobuf

message MsgPayFor {
  option (cosmos.msg.v1.signer) = "from_address";
  bytes resource = 1;
  string   from_address                    = 2 [(cosmos_proto.scalar) = "cosmos.AddressString"];
  repeated cosmos.base.v1beta1.Coin amount = 3 [
    (gogoproto.nullable)     = false,
    (amino.dont_omitempty)   = true,
    (amino.encoding)         = "legacy_coins",
    (gogoproto.castrepeated) = "github.com/cosmos/cosmos-sdk/types.Coins"
  ];
}
```

```protobuf
// MsgPayForResponse defines the Msg/PayFor response type.
message MsgPayForResponse{}
```

## Queries

### MetaData

Request:

```protobuf
// QueryMetadataRequest defines the request type for metadata
message QueryMetadataRequest {
  string identifier = 1;
}
```

The RESTAPI endpoint for this query is `/nillion/meta/v1/meta_data/{identifier}`

Response:

```protobuf
// QueryMetadataResponse defines the response type for metadata
message QueryMetadataResponse {
  bytes metadata = 1;
}
```

### AllMetaData

Request:

```protobuf
// QueryAllMetadataRequest defines the request type for all metadata
message QueryAllMetadataRequest {
  cosmos.base.query.v1beta1.PageRequest pagination = 1;
}
```

The RESTAPI endpoint for this query is "/nillion/meta/v1/all_meta_data"

```protobuf
// QueryAllMetadataResponse defines the response type for all metadata
message QueryAllMetadataResponse {
  repeated bytes metadata = 1;
}
```

### CLI

```
nilchaind tx meta pay-for <address-paid-for> 1000anillion /path-to-resource/resource.json --from alice --node tcp://localhost:26657 --chain-id demo
```
