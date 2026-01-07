package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	criblcontrolplanesdkgo "github.com/criblio/cribl-control-plane-sdk-go"
	"github.com/criblio/cribl-control-plane-sdk-go/models/components"
	"github.com/criblio/cribl-control-plane-sdk-go/models/operations"
)

// Config holds the SDK client and server configuration
type Config struct {
	Client    *criblcontrolplanesdkgo.CriblControlPlane
	ServerURL string
}

// loginResponse represents the response from the login endpoint
type loginResponse struct {
	Token string `json:"token"`
}

// login authenticates with the Cribl server and returns a bearer token
func login(baseURL, username, password string) (string, error) {
	loginURL := baseURL + "/api/v1/auth/login"

	loginData := map[string]string{
		"username": username,
		"password": password,
	}

	jsonData, err := json.Marshal(loginData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal login data: %w", err)
	}

	req, err := http.NewRequest("POST", loginURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create login request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute login request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body to check what we got
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("login failed with status code: %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	var loginResp loginResponse
	if err := json.Unmarshal(bodyBytes, &loginResp); err != nil {
		return "", fmt.Errorf("failed to decode login response: %w, response body: %s", err, string(bodyBytes))
	}

	if loginResp.Token == "" {
		return "", fmt.Errorf("login response did not contain a token, response: %s", string(bodyBytes))
	}

	return loginResp.Token, nil
}

// newClient creates a new SDK client with configuration from environment variables
func newClient() (*Config, error) {
	// Get base server URL from environment or use default
	baseURL := os.Getenv("CRIBL_SERVER_URL")
	if baseURL == "" {
		baseURL = "http://localhost:9000"
	}

	// Remove trailing slash and ensure we have the base URL without /api/v1
	baseURL = strings.TrimSuffix(baseURL, "/")
	baseURL = strings.TrimSuffix(baseURL, "/api/v1")

	// Construct server URL with /api/v1
	serverURL := baseURL + "/api/v1"

	// Get username and password from environment, default to "admin"
	username := os.Getenv("CRIBL_USERNAME")
	if username == "" {
		username = os.Getenv("CRIBL_USER")
	}
	if username == "" {
		username = "admin"
	}

	password := os.Getenv("CRIBL_PASSWORD")
	if password == "" {
		password = os.Getenv("CRIBL_PASS")
	}
	if password == "" {
		password = "admin"
	}

	// Login to get auth token
	fmt.Println("🔐 Authenticating with Cribl server...")
	authToken, err := login(baseURL, username, password)
	if err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}
	fmt.Println("✅ Authentication successful!")

	// Create SDK client with security
	security := components.Security{
		BearerAuth: &authToken,
	}

	sdk := criblcontrolplanesdkgo.New(
		serverURL,
		criblcontrolplanesdkgo.WithSecurity(security),
	)

	return &Config{
		Client:    sdk,
		ServerURL: serverURL,
	}, nil
}

func main() {
	fmt.Println("🚀 Google Cloud Storage Destination Example")
	fmt.Println("===========================================")

	// Create client
	config, err := newClient()
	if err != nil {
		log.Fatalf("❌ Failed to create client: %v", err)
	}

	fmt.Printf("🌐 Connected to: %s\n", config.ServerURL)
	fmt.Println()

	ctx := context.Background()

	destID := "output-google-cloud-storage"

	// Cleanup: Try to delete destination if it exists
	_, _ = config.Client.Destinations.Delete(ctx, destID, operations.WithServerURL(config.ServerURL))

	// Create Google Cloud Storage destination using direct struct initialization
	createReq := components.CreateOutputGoogleCloudStorage(
		components.OutputGoogleCloudStorage{
			ID:   criblcontrolplanesdkgo.Pointer(destID),
			Type: components.OutputGoogleCloudStorageTypeGoogleCloudStorage,
			SystemFields: []string{
				"cribl_pipe",
			},
			Streamtags: []string{
				"web",
				"sdk",
				"test",
			},
			Endpoint:                criblcontrolplanesdkgo.Pointer("https://storage.googleapis.com"),
			SignatureVersion:        criblcontrolplanesdkgo.Pointer(components.OutputGoogleCloudStorageSignatureVersion("v4")),
			AwsAuthenticationMethod: criblcontrolplanesdkgo.Pointer(components.OutputGoogleCloudStorageAuthenticationMethod("manual")),
			StagePath:               criblcontrolplanesdkgo.Pointer("$CRIBL_HOME/state/outputs/staging"),
			DestPath:                criblcontrolplanesdkgo.Pointer("`my-prefix-name`"),
			VerifyPermissions:       criblcontrolplanesdkgo.Pointer(true),
			ObjectACL:               criblcontrolplanesdkgo.Pointer(components.OutputGoogleCloudStorageObjectACLPrivate),
			ReuseConnections:        criblcontrolplanesdkgo.Pointer(true),
			RejectUnauthorized:      criblcontrolplanesdkgo.Pointer(true),
			AddIDToStagePath:        criblcontrolplanesdkgo.Pointer(true),
			RemoveEmptyDirs:         criblcontrolplanesdkgo.Pointer(true),
			Format:                  criblcontrolplanesdkgo.Pointer(components.OutputGoogleCloudStorageDataFormatJSON),
			BaseFileName:            criblcontrolplanesdkgo.Pointer("`CriblOut`"),
			FileNameSuffix:          criblcontrolplanesdkgo.Pointer("`.${C.env[\"CRIBL_WORKER_ID\"]}.${__format}${__compression === \"gzip\" ? \".gz\" : \"\"}`"),
			MaxFileSizeMB:           criblcontrolplanesdkgo.Pointer(32.0),
			MaxFileOpenTimeSec:      criblcontrolplanesdkgo.Pointer(300.0),
			MaxFileIdleTimeSec:      criblcontrolplanesdkgo.Pointer(30.0),
			MaxOpenFiles:            criblcontrolplanesdkgo.Pointer(100.0),
			HeaderLine:              criblcontrolplanesdkgo.Pointer(""),
			WriteHighWaterMark:      criblcontrolplanesdkgo.Pointer(64.0),
			OnBackpressure:          criblcontrolplanesdkgo.Pointer(components.OutputGoogleCloudStorageBackpressureBehaviorBlock),
			DeadletterEnabled:       criblcontrolplanesdkgo.Pointer(false),
			OnDiskFullBackpressure:  criblcontrolplanesdkgo.Pointer(components.OutputGoogleCloudStorageDiskSpaceProtectionBlock),
			Compress:                criblcontrolplanesdkgo.Pointer(components.OutputGoogleCloudStorageCompression("gzip")),
			CompressionLevel:        criblcontrolplanesdkgo.Pointer(components.OutputGoogleCloudStorageCompressionLevelBestSpeed),
			EmptyDirCleanupSec:      criblcontrolplanesdkgo.Pointer(300.0),
			Description:             criblcontrolplanesdkgo.Pointer("Output Google Cloud Storage"),
			Bucket:                  "`m-bucket-name`",
			Region:                  "US-EAST1",
			AwsAPIKey:               criblcontrolplanesdkgo.Pointer("`my-access-key`"),
			AwsSecretKey:            criblcontrolplanesdkgo.Pointer("`my-secret`"),
		},
	)

	// CREATE: Create destination
	fmt.Println("1️⃣ CREATE Destination")
	fmt.Println("-------------------")
	requestJSON, _ := json.MarshalIndent(createReq, "", "  ")
	fmt.Println("REQUEST:")
	fmt.Println(string(requestJSON))
	fmt.Println()

	createResponse, err := config.Client.Destinations.Create(ctx, createReq, operations.WithServerURL(config.ServerURL))
	if err != nil {
		log.Fatalf("❌ Failed to create destination: %v", err)
	}

	if createResponse.Object == nil || createResponse.Object.Count == nil || *createResponse.Object.Count == 0 {
		log.Fatalf("❌ Create failed: count is 0")
	}
	if len(createResponse.Object.Items) == 0 {
		log.Fatalf("❌ Create failed: no items returned")
	}

	fmt.Println("✅ Destination created successfully!")
	fmt.Println()

	// READ: Read destination
	fmt.Println("2️⃣ READ Destination")
	fmt.Println("------------------")
	fmt.Printf("Reading destination ID: %s\n", destID)
	fmt.Println()

	readResponse, err := config.Client.Destinations.Get(ctx, destID, operations.WithServerURL(config.ServerURL))
	if err != nil {
		log.Fatalf("❌ Failed to read destination: %v", err)
	}

	if readResponse.Object == nil || readResponse.Object.Count == nil || *readResponse.Object.Count != 1 {
		log.Fatalf("❌ Read failed: expected count 1, got %v", readResponse.Object.Count)
	}
	if len(readResponse.Object.Items) == 0 {
		log.Fatalf("❌ Read failed: no items returned")
	}

	retrievedDestination := readResponse.Object.Items[0]

	// Convert response to map for comparison
	retrievedBytes, err := json.Marshal(retrievedDestination)
	if err != nil {
		log.Fatalf("❌ Failed to marshal response: %v", err)
	}
	var responseMap map[string]interface{}
	err = json.Unmarshal(retrievedBytes, &responseMap)
	if err != nil {
		log.Fatalf("❌ Failed to convert response to map: %v", err)
	}

	responseJSON, _ := json.MarshalIndent(responseMap, "", "  ")
	fmt.Println("RESPONSE:")
	fmt.Println(string(responseJSON))
	fmt.Println()

	// Compare REQUEST vs RESPONSE
	fmt.Println("3️⃣ REQUEST vs RESPONSE Comparison")
	fmt.Println("----------------------------------")
	fmt.Println("REQUEST (what we sent):")
	fmt.Println(string(requestJSON))
	fmt.Println()
	fmt.Println("RESPONSE (what we received):")
	fmt.Println(string(responseJSON))
	fmt.Println()

	// Cleanup: Delete destination
	_, err = config.Client.Destinations.Delete(ctx, destID, operations.WithServerURL(config.ServerURL))
	if err != nil {
		log.Fatalf("❌ Failed to delete destination: %v", err)
	}
	fmt.Println("✅ Cleanup completed")

	fmt.Println()
	fmt.Println("🎉 Example completed successfully!")
}
