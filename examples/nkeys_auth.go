package examples

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
)

type NKeyUser struct {
	Name      string
	Seed      string
	PublicKey string
	CanPublishTo   []string
	CanSubscribeTo []string
}

var predefinedUsers = []NKeyUser{
	{
		Name:           "Admin",
		Seed:           "SUACSSL3UAHUDXKFSNVUZRF5UHPMWZ6BFDTJ7M6USDXIEDNPPQYYYCU3VY",
		PublicKey:      "UDXU4RCSJNZOIQHZNWXHXORDPRTGNJAHAHFRGZNEEJCPQTT2M7NLCNF4",
		CanPublishTo:   []string{"any.subject", "admin.>"},
		CanSubscribeTo: []string{"any.subject", "admin.>"},
	},
	{
		Name:           "Client",
		Seed:           "SUAM42UG6PV55WVNPAHKF65J4SJQNWQVNWQP7H2VQWPQVH2SJQNWQVH2SABC",
		PublicKey:      "UAH42UG6PV55WVNPAHKF65J4SJQNWQVNWQP7H2VQWPQVH2SJQNWQVH2S",
		CanPublishTo:   []string{"req.a", "req.b"},
		CanSubscribeTo: []string{"_INBOX.>"},
	},
	{
		Name:           "Service",
		Seed:           "SUBFJ4RCSJNZOIQHZNWXHXORDPRTGNJAHAHFRGZNEEJCPQTT2M7NLCNF5XYZ",
		PublicKey:      "UBFJ4RCSJNZOIQHZNWXHXORDPRTGNJAHAHFRGZNEEJCPQTT2M7NLCNF5",
		CanPublishTo:   []string{"_INBOX.>"},
		CanSubscribeTo: []string{"req.a", "req.b"},
	},
	{
		Name:           "Other",
		Seed:           "SUCGH5RCSJNZOIQHZNWXHXORDPRTGNJAHAHFRGZNEEJCPQTT2M7NLCNF6DEF",
		PublicKey:      "UCGH5RCSJNZOIQHZNWXHXORDPRTGNJAHAHFRGZNEEJCPQTT2M7NLCNF6",
		CanPublishTo:   []string{"SANDBOX.*"},
		CanSubscribeTo: []string{"PUBLIC.>", "_INBOX.>"},
	},
}

func nkeyOption(seed string) nats.Option {
	return func(o *nats.Options) error {
		kp, err := nkeys.FromSeed([]byte(seed))
		if err != nil {
			return fmt.Errorf("failed to parse seed: %w", err)
		}

		publicKey, err := kp.PublicKey()
		if err != nil {
			return fmt.Errorf("failed to get public key: %w", err)
		}

		o.Nkey = publicKey

		o.SignatureCB = func(nonce []byte) ([]byte, error) {
			sig, err := kp.Sign(nonce)
			if err != nil {
				return nil, fmt.Errorf("failed to sign nonce: %w", err)
			}
			return sig, nil
		}

		return nil
	}
}

func connectWithNKey(user NKeyUser) (*nats.Conn, error) {
	nc, err := nats.Connect(
		"nats://localhost:4227",
		nkeyOption(user.Seed),
		nats.Name(user.Name),
		nats.Timeout(5*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("connection failed for %s: %w", user.Name, err)
	}
	return nc, nil
}

func DemoNKeysAuth() {
	fmt.Println("\n=== NKeys Authentication Demo ===")
	fmt.Println("Demonstrating Ed25519 signature-based authentication")

	fmt.Println("\n📋 Pre-configured test users:")
	for _, user := range predefinedUsers {
		fmt.Printf("\n%s:\n", user.Name)
		fmt.Printf("  Public Key: %s\n", user.PublicKey)
		fmt.Printf("  Seed: %s... (secret)\n", user.Seed[:20])
		fmt.Printf("  Can publish to: %v\n", user.CanPublishTo)
		fmt.Printf("  Can subscribe to: %v\n", user.CanSubscribeTo)
	}

	fmt.Println("\n" + strings.Repeat("─", 60))

	for _, user := range predefinedUsers {
		fmt.Printf("\n🔐 Testing %s User:\n", user.Name)

		nc, err := connectWithNKey(user)
		if err != nil {
			log.Printf("  ✗ Connection failed: %v", err)
			continue
		}
		fmt.Printf("  ✓ Connected using NKey signature authentication\n")

		if len(user.CanPublishTo) > 0 {
			subject := user.CanPublishTo[0]
			if err := nc.Publish(subject, []byte(fmt.Sprintf("Message from %s", user.Name))); err != nil {
				fmt.Printf("  ✗ Publish to '%s' failed: %v\n", subject, err)
			} else {
				fmt.Printf("  ✓ Published to '%s'\n", subject)
			}
		}

		if len(user.CanSubscribeTo) > 0 {
			subject := user.CanSubscribeTo[0]
			sub, err := nc.SubscribeSync(subject)
			if err != nil {
				fmt.Printf("  ✗ Subscribe to '%s' failed: %v\n", subject, err)
			} else {
				fmt.Printf("  ✓ Subscribed to '%s'\n", subject)
				sub.Unsubscribe()
			}
		}

		invalidSubject := "unauthorized.subject"
		if err := nc.Publish(invalidSubject, []byte("test")); err != nil {
			fmt.Printf("  ✓ Correctly denied publishing to '%s'\n", invalidSubject)
		} else {
			fmt.Printf("  ⚠️  Unexpected: allowed to publish to '%s'\n", invalidSubject)
		}

		nc.Close()
	}

	fmt.Println("\n=== Request-Response Pattern with NKeys ===")

	fmt.Println("\n1. Starting service responder...")
	serviceNC, err := connectWithNKey(predefinedUsers[2])
	if err != nil {
		log.Printf("Service connection failed: %v", err)
		return
	}
	defer serviceNC.Close()

	serviceSub, err := serviceNC.Subscribe("req.a", func(m *nats.Msg) {
		response := fmt.Sprintf("Response to: %s", string(m.Data))
		m.Respond([]byte(response))
		fmt.Printf("  ✓ Service responded to request\n")
	})
	if err != nil {
		log.Printf("Service subscribe failed: %v", err)
		return
	}
	defer serviceSub.Unsubscribe()
	fmt.Println("  ✓ Service listening on 'req.a'")

	fmt.Println("\n2. Client making request...")
	clientNC, err := connectWithNKey(predefinedUsers[1])
	if err != nil {
		log.Printf("Client connection failed: %v", err)
		return
	}
	defer clientNC.Close()

	msg, err := clientNC.Request("req.a", []byte("Hello from client"), 2*time.Second)
	if err != nil {
		log.Printf("  ✗ Request failed: %v", err)
	} else {
		fmt.Printf("  ✓ Client received response: %s\n", string(msg.Data))
	}

	fmt.Println("\n=== NKeys Authentication Demo Complete ===")
	fmt.Println("\n🔑 Key Advantages of NKeys:")
	fmt.Println("  • Private keys never leave the client")
	fmt.Println("  • Server only stores public keys")
	fmt.Println("  • Each connection uses a unique challenge-response")
	fmt.Println("  • Immune to replay attacks")
	fmt.Println("  • Based on Ed25519 (faster and more secure than RSA)")
}
