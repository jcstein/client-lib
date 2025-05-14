package main

import (
	"context"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"

	libshare "github.com/celestiaorg/go-square/v2/share"

	"github.com/celestiaorg/celestia-node/api/client"
	"github.com/celestiaorg/celestia-node/blob"
)

func main() {
	// Initialize keyring with new key
	keyname := "my_celes_key"
	kr, err := client.KeyringWithNewKey(client.KeyringConfig{
		KeyName:     keyname,
		BackendName: keyring.BackendTest,
	}, "../../.celestia-light-mocha-4/keys/")
	if err != nil {
		fmt.Println("failed to create keyring:", err)
		return
	}

	// Configure client
	cfg := client.Config{
		// this can also be used with light node running at http://localhost:26658, must use --rpc.skip auth or add auth token
		ReadConfig: client.ReadConfig{
			BridgeDAAddr: "https://clean-wiser-glitter.celestia-mocha.quiknode.pro/1adb9c87496929d258c7c358889d921963011005",
			DAAuthToken:  "",
			EnableDATLS:  true,
		},
		SubmitConfig: client.SubmitConfig{
			DefaultKeyName: keyname,
			Network:        "mocha-4",
			CoreGRPCConfig: client.CoreGRPCConfig{
				Addr:       "clean-wiser-glitter.celestia-mocha.quiknode.pro:9090", // or rpc-mocha.pops.one:9090
				TLSEnabled: true,
				AuthToken:  "1adb9c87496929d258c7c358889d921963011005",
			},
		},
	}

	// Create client with full submission capabilities
	ctx := context.Background()
	client, err := client.New(ctx, cfg, kr)
	if err != nil {
		fmt.Println("failed to create client:", err)
		return
	}

	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	// Submit a blob
	namespace := libshare.MustNewV0Namespace([]byte("example"))
	b, err := blob.NewBlob(libshare.ShareVersionZero, namespace, []byte("data"), nil)
	if err != nil {
		fmt.Println("failed to create blob:", err)
		return
	}
	height, err := client.Blob.Submit(ctx, []*blob.Blob{b}, nil)
	if err != nil {
		fmt.Println("failed to submit blob:", err)
		return
	}
	fmt.Println("submitted blob", height)

	// Retrieve a blob
	retrievedBlob, err := client.Blob.Get(ctx, height, namespace, b.Commitment)
	if err != nil {
		fmt.Println("failed to retrieve blob:", err)
		return
	}
	fmt.Println("retrieved blob", string(retrievedBlob.Data()))
}
