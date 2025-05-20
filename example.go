package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"

	libshare "github.com/celestiaorg/go-square/v2/share"

	"github.com/celestiaorg/celestia-node/api/client"
	"github.com/celestiaorg/celestia-node/blob"
)

func main() {
	// Set up logging with microsecond precision
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.Println("Starting Celestia client example")
	
	// Load environment variables from .env
	log.Println("Loading environment variables")
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file:", err)
		return
	}
	log.Println("Environment variables loaded successfully")

	// Initialize keyring with new key
	keyname := "my_celes_key"
	log.Println("Initializing keyring with key name:", keyname)
	kr, err := client.KeyringWithNewKey(client.KeyringConfig{
		KeyName:     keyname,
		BackendName: keyring.BackendTest,
	}, "../../.celestia-light-mocha-4/keys/")
	if err != nil {
		log.Println("Failed to create keyring:", err)
		return
	}
	log.Println("Keyring successfully initialized")

	// Get authentication credentials from environment variables
	log.Println("Configuring Celestia client")
	quickNodeAuthToken := os.Getenv("QUICKNODE_AUTH_TOKEN")
	if quickNodeAuthToken == "" {
		log.Println("QUICKNODE_AUTH_TOKEN not found in environment variables")
		return
	}
	
	quickNodeBridgeURL := os.Getenv("QUICKNODE_BRIDGE_URL")
	if quickNodeBridgeURL == "" {
		log.Println("QUICKNODE_BRIDGE_URL not found in environment variables")
		return
	}
	
	quickNodeGRPCURL := os.Getenv("QUICKNODE_GRPC_URL")
	if quickNodeGRPCURL == "" {
		log.Println("QUICKNODE_GRPC_URL not found in environment variables")
		return
	}

	log.Println("Setting up QuickNode authentication")
	
	// Initialize empty HTTP headers (not needed with token in URL)
	emptyHeaders := http.Header{}
	
	// Build the full URL with authentication token
	fullBridgeURL := quickNodeBridgeURL + "/" + quickNodeAuthToken
	log.Printf("Using bridge URL: %s", fullBridgeURL)
	
	// Configure Celestia client
	cfg := client.Config{
		ReadConfig: client.ReadConfig{
			BridgeDAAddr: fullBridgeURL,
			HTTPHeader:   emptyHeaders,
			EnableDATLS:  true,
		},
		SubmitConfig: client.SubmitConfig{
			DefaultKeyName: keyname,
			Network:        "mocha-4",
			CoreGRPCConfig: client.CoreGRPCConfig{
				Addr:       quickNodeGRPCURL,
				TLSEnabled: true,
				AuthToken:  quickNodeAuthToken,
			},
		},
	}
	log.Println("Client configuration complete")

	// Initialize Celestia client
	log.Println("Creating Celestia client")
	ctx := context.Background()
	client, err := client.New(ctx, cfg, kr)
	if err != nil {
		log.Println("Failed to create client:", err)
		return
	}
	log.Println("Client successfully created")

	// Set timeout for operations
	log.Println("Setting up context with 1-minute timeout")
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	// Create and submit blob
	log.Println("Creating namespace and blob for submission")
	namespace := libshare.MustNewV0Namespace([]byte("example"))
	log.Println("Namespace created")
	
	data := []byte("data")
	log.Printf("Creating new blob with %d bytes of data", len(data))
	b, err := blob.NewBlob(libshare.ShareVersionZero, namespace, data, nil)
	if err != nil {
		log.Println("Failed to create blob:", err)
		return
	}
	log.Println("Blob successfully created with commitment:", b.Commitment)

	// Submit blob to network
	log.Println("Submitting blob to the network...")
	submitStart := time.Now()
	height, err := client.Blob.Submit(ctx, []*blob.Blob{b}, nil)
	if err != nil {
		log.Println("Failed to submit blob:", err)
		return
	}
	log.Printf("Blob successfully submitted at height %d (took %v)", height, time.Since(submitStart))

	// Retrieve blob from network
	log.Println("Retrieving blob from the network...")
	log.Printf("Querying for blob at height %d with namespace %x and commitment %x", height, namespace, b.Commitment)
	getStart := time.Now()
	retrievedBlob, err := client.Blob.Get(ctx, height, namespace, b.Commitment)
	if err != nil {
		log.Println("Failed to retrieve blob:", err)
		return
	}
	log.Printf("Blob successfully retrieved (took %v)", time.Since(getStart))
	log.Printf("Retrieved data: %s", string(retrievedBlob.Data()))
	log.Println("Example completed successfully")
}
