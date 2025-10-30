package examples

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nats-io/nkeys"
)

type GeneratedNKey struct {
	Role      string
	Seed      string
	PublicKey string
}

func GenerateNKeysForRoles() ([]GeneratedNKey, error) {
	roles := []string{"Admin", "Client", "Service", "Other"}
	keys := make([]GeneratedNKey, 0, len(roles))

	for _, role := range roles {
		kp, err := nkeys.CreateUser()
		if err != nil {
			return nil, fmt.Errorf("failed to create nkey for %s: %w", role, err)
		}

		seed, err := kp.Seed()
		if err != nil {
			return nil, fmt.Errorf("failed to get seed for %s: %w", role, err)
		}

		publicKey, err := kp.PublicKey()
		if err != nil {
			return nil, fmt.Errorf("failed to get public key for %s: %w", role, err)
		}

		keys = append(keys, GeneratedNKey{
			Role:      role,
			Seed:      string(seed),
			PublicKey: publicKey,
		})
	}

	return keys, nil
}

func SaveNKeysToFile(keys []GeneratedNKey, filename string) error {
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	fmt.Fprintln(f, "# Generated NKeys for NATS Authentication")
	fmt.Fprintln(f, "# Keep the seeds (private keys) secret!")
	fmt.Fprintln(f, "# Only share the public keys with the NATS server")
	fmt.Fprintln(f, "")

	for _, key := range keys {
		fmt.Fprintf(f, "# %s User\n", key.Role)
		fmt.Fprintf(f, "Seed (Private Key):  %s\n", key.Seed)
		fmt.Fprintf(f, "Public Key:          %s\n", key.PublicKey)
		fmt.Fprintln(f, "")
	}

	return nil
}

func GenerateServerConfig(keys []GeneratedNKey, filename string) error {
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	fmt.Fprintln(f, "# Generated NKeys Authentication Configuration")
	fmt.Fprintln(f, "# Auto-generated - modify as needed")
	fmt.Fprintln(f, "")
	fmt.Fprintln(f, "port: 4227")
	fmt.Fprintln(f, "")
	fmt.Fprintln(f, "authorization {")
	fmt.Fprintln(f, "  default_permissions = {")
	fmt.Fprintln(f, "    publish = \"SANDBOX.*\"")
	fmt.Fprintln(f, "    subscribe = [\"PUBLIC.>\", \"_INBOX.>\"]")
	fmt.Fprintln(f, "  }")
	fmt.Fprintln(f, "")
	fmt.Fprintln(f, "  ADMIN = {")
	fmt.Fprintln(f, "    publish = \">\"")
	fmt.Fprintln(f, "    subscribe = \">\"")
	fmt.Fprintln(f, "  }")
	fmt.Fprintln(f, "")
	fmt.Fprintln(f, "  REQUESTOR = {")
	fmt.Fprintln(f, "    publish = [\"req.a\", \"req.b\"]")
	fmt.Fprintln(f, "    subscribe = \"_INBOX.>\"")
	fmt.Fprintln(f, "  }")
	fmt.Fprintln(f, "")
	fmt.Fprintln(f, "  RESPONDER = {")
	fmt.Fprintln(f, "    subscribe = [\"req.a\", \"req.b\"]")
	fmt.Fprintln(f, "    publish = \"_INBOX.>\"")
	fmt.Fprintln(f, "  }")
	fmt.Fprintln(f, "")
	fmt.Fprintln(f, "  users = [")

	for i, key := range keys {
		var permission string
		switch key.Role {
		case "Admin":
			permission = "$ADMIN"
		case "Client":
			permission = "$REQUESTOR"
		case "Service":
			permission = "$RESPONDER"
		default:
			permission = ""
		}

		fmt.Fprintf(f, "    # %s User\n", key.Role)
		if permission != "" {
			fmt.Fprintf(f, "    {nkey: \"%s\", permissions: %s}", key.PublicKey, permission)
		} else {
			fmt.Fprintf(f, "    {nkey: \"%s\"}", key.PublicKey)
		}

		if i < len(keys)-1 {
			fmt.Fprintln(f, ",")
		} else {
			fmt.Fprintln(f, "")
		}
	}

	fmt.Fprintln(f, "  ]")
	fmt.Fprintln(f, "}")

	return nil
}

func DemoNKeyGenerationWithFiles() {
	fmt.Println("\n=== NKey Generation with File Export Demo ===")
	
	fmt.Println("\nGenerating NKey pairs for different roles...")
	keys, err := GenerateNKeysForRoles()
	if err != nil {
		fmt.Printf("Error generating keys: %v\n", err)
		return
	}

	fmt.Println("\n📋 Generated Keys:")
	for _, key := range keys {
		fmt.Printf("\n%s:\n", key.Role)
		fmt.Printf("  Public Key: %s\n", key.PublicKey)
		fmt.Printf("  Seed:       %s\n", key.Seed)
	}

	keysFile := "generated/nkeys.txt"
	configFile := "generated/nkeys-server.conf"

	fmt.Printf("\n💾 Saving keys to: %s\n", keysFile)
	if err := SaveNKeysToFile(keys, keysFile); err != nil {
		fmt.Printf("Error saving keys: %v\n", err)
		return
	}
	fmt.Println("✓ Keys saved successfully")

	fmt.Printf("\n💾 Generating server config: %s\n", configFile)
	if err := GenerateServerConfig(keys, configFile); err != nil {
		fmt.Printf("Error generating config: %v\n", err)
		return
	}
	fmt.Println("✓ Server config generated successfully")

	fmt.Println("\n📖 Next Steps:")
	fmt.Println("  1. Review the generated keys in:", keysFile)
	fmt.Println("  2. Store the seeds (private keys) securely")
	fmt.Println("  3. Start NATS server with:", configFile)
	fmt.Printf("     Command: nats-server -c %s\n", configFile)
	fmt.Println("  4. Use the seeds in your client applications")

	fmt.Println("\n=== Generation Complete ===")
}
