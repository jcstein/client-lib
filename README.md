# Celestia Client Library Test

This repository provides a simple implementation that tests the Celestia client library functionality. It demonstrates how to interact with the Celestia network using the official Go client library.

## Overview

[Celestia](https://celestia.org/) is a modular data availability network that securely scales with the number of users, making it easy for anyone to launch their own blockchain. This repository contains example code for interacting with Celestia nodes using the official client library.

## Prerequisites

- Go 1.19 or later
- Access to a Celestia network (testnet or local node)

## Installation

Clone this repository:

```bash
git clone https://github.com/your-username/client-lib.git
cd client-lib
```

Install dependencies:

```bash
go mod tidy
```

## Configuration

The example code in `example.go` is configured to work with:

- A QuickNode endpoint for reading data from the Celestia network
- The public Mocha-4 testnet for transaction submission
- Existing keys located at `../../.celestia-light-mocha-4/keys/`

You may need to modify the following configurations:

- Read configuration (QuickNode endpoint or local node)
- Submit configuration (consensus node endpoint)
- Keyring location (currently pointing to existing keys)

## Usage

Run the example with:

```bash
go run example.go
```

The example demonstrates:

1. Setting up a keyring with a new key
2. Configuring a client with both read and submit capabilities
3. Submitting a blob (data) to the Celestia network
4. Retrieving the blob from the network to verify successful submission

## Features Demonstrated

- **Client Initialization**: Setting up a full-featured Celestia client
- **Keyring Management**: Creating and using keys for transaction signing
- **Blob Submission**: Submitting data to the Celestia network
- **Blob Retrieval**: Querying and retrieving data from the network

## Example Code Walkthrough

The `example.go` file demonstrates:

```go
// Initialize keyring with existing key
kr, err := client.KeyringWithNewKey(client.KeyringConfig{
    KeyName:     keyname,
    BackendName: keyring.BackendTest,
}, "../../.celestia-light-mocha-4/keys/")

// Configure client with read and submit capabilities
cfg := client.Config{
    ReadConfig: client.ReadConfig{
        BridgeDAAddr: "https://clean-wiser-glitter.celestia-mocha.quiknode.pro/...",
        DAAuthToken:  "",
        EnableDATLS:  false,
    },
    SubmitConfig: client.SubmitConfig{
        DefaultKeyName: keyname,
        Network:        "mocha-4",
        CoreGRPCConfig: client.CoreGRPCConfig{
            Addr:       "celestia-testnet-consensus.itrocket.net:9090",
            TLSEnabled: false,
            AuthToken:  "",
        },
    },
}

// Create client with full submission capabilities
client, err := client.New(ctx, cfg, kr)

// Submit a blob
namespace := libshare.MustNewV0Namespace([]byte("example"))
b, err := blob.NewBlob(libshare.ShareVersionZero, namespace, []byte("data"), nil)
height, err := client.Blob.Submit(ctx, []*blob.Blob{b}, nil)

// Retrieve a blob
retrievedBlob, err := client.Blob.Get(ctx, height, namespace, b.Commitment)
```

## References

- [Celestia Client Library Documentation](https://github.com/celestiaorg/celestia-node/blob/celestia-client-lib/api/client/readme.md)
- [Celestia Node v0.22.2-client-lib-rc Release](https://github.com/celestiaorg/celestia-node/releases/tag/v0.22.2-client-lib-rc)

## License

This project uses the same license as the Celestia Node repository.
