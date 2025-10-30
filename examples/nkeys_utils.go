package examples

import (
	"encoding/base64"
	"fmt"

	"github.com/nats-io/nkeys"
)

type NKeyPair struct {
	Seed      string
	PublicKey string
}

func GenerateUserNKey() (*NKeyPair, error) {
	kp, err := nkeys.CreateUser()
	if err != nil {
		return nil, fmt.Errorf("failed to create user nkey: %w", err)
	}

	seed, err := kp.Seed()
	if err != nil {
		return nil, fmt.Errorf("failed to get seed: %w", err)
	}

	publicKey, err := kp.PublicKey()
	if err != nil {
		return nil, fmt.Errorf("failed to get public key: %w", err)
	}

	return &NKeyPair{
		Seed:      string(seed),
		PublicKey: publicKey,
	}, nil
}

func SignChallenge(seed string, challenge []byte) ([]byte, error) {
	kp, err := nkeys.FromSeed([]byte(seed))
	if err != nil {
		return nil, fmt.Errorf("failed to parse seed: %w", err)
	}

	sig, err := kp.Sign(challenge)
	if err != nil {
		return nil, fmt.Errorf("failed to sign challenge: %w", err)
	}

	return sig, nil
}

func VerifySignature(publicKey string, challenge, signature []byte) error {
	kp, err := nkeys.FromPublicKey(publicKey)
	if err != nil {
		return fmt.Errorf("failed to parse public key: %w", err)
	}

	if err := kp.Verify(challenge, signature); err != nil {
		return fmt.Errorf("signature verification failed: %w", err)
	}

	return nil
}

func PrintNKeyPair(pair *NKeyPair, label string) {
	fmt.Printf("\n%s NKey Pair:\n", label)
	fmt.Printf("├─ Seed (Private Key):  %s\n", pair.Seed)
	fmt.Printf("└─ Public Key:          %s\n", pair.PublicKey)
	fmt.Println("\n⚠️  Keep the seed secret! Only share the public key.")
}

func DemoNKeyGeneration() {
	fmt.Println("\n=== NKey Generation Demo ===")
	fmt.Println("Generating NKey pairs for different users...")

	users := []string{"Admin", "Client", "Service"}
	
	for _, user := range users {
		pair, err := GenerateUserNKey()
		if err != nil {
			fmt.Printf("Error generating %s nkey: %v\n", user, err)
			continue
		}
		PrintNKeyPair(pair, user)
	}

	fmt.Println("\n=== Signature Demo ===")
	fmt.Println("Demonstrating challenge-response authentication...")

	pair, err := GenerateUserNKey()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	challenge := []byte("random-server-challenge-12345")
	fmt.Printf("\nChallenge: %s\n", base64.StdEncoding.EncodeToString(challenge))

	signature, err := SignChallenge(pair.Seed, challenge)
	if err != nil {
		fmt.Printf("Error signing: %v\n", err)
		return
	}
	fmt.Printf("Signature: %s\n", base64.StdEncoding.EncodeToString(signature))

	err = VerifySignature(pair.PublicKey, challenge, signature)
	if err != nil {
		fmt.Printf("✗ Verification failed: %v\n", err)
	} else {
		fmt.Println("✓ Signature verified successfully!")
	}

	fmt.Println("\n=== NKey Generation Demo Complete ===")
}
