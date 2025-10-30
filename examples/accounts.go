package examples

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

// DemoAccounts demonstrates account isolation and multi-tenancy
func DemoAccounts() {
	fmt.Println("\n=== Account Isolation Demo ===")
	
	// Connect to each account
	connA, err := nats.Connect("nats://user_a:pass_a@localhost:4226")
	if err != nil {
		log.Printf("Account A connection failed: %v", err)
		return
	}
	defer connA.Close()
	
	connB, err := nats.Connect("nats://user_b:pass_b@localhost:4226")
	if err != nil {
		log.Printf("Account B connection failed: %v", err)
		return
	}
	defer connB.Close()
	
	connC, err := nats.Connect("nats://user_c:pass_c@localhost:4226")
	if err != nil {
		log.Printf("Account C connection failed: %v", err)
		return
	}
	defer connC.Close()
	
	// 1. Demonstrate account isolation
	fmt.Println("\n1. Testing Account Isolation:")
	
	// Account B subscribes to private subject
	subB, err := connB.SubscribeSync("private.data")
	if err != nil {
		log.Printf("Account B subscribe failed: %v", err)
		return
	}
	
	// Account A publishes to same subject
	fmt.Println("  Account A publishing to 'private.data'...")
	connA.Publish("private.data", []byte("Message from A"))
	connA.Flush()
	
	// Account B should not receive it (different accounts)
	time.Sleep(200 * time.Millisecond)
	_, err = subB.NextMsg(500 * time.Millisecond)
	if err == nats.ErrTimeout {
		fmt.Println("  ✓ Account B correctly did NOT receive message from Account A")
		fmt.Println("    (Accounts are isolated)")
	} else {
		fmt.Println("  ✗ Account B received message (isolation failed)")
	}
	subB.Unsubscribe()
	
	// Account B publishes to its own subject
	fmt.Println("\n  Account B publishing to 'private.data'...")
	subB2, _ := connB.SubscribeSync("private.data")
	connB.Publish("private.data", []byte("Message from B"))
	connB.Flush()
	
	// Account B should receive its own message
	if msg, err := subB2.NextMsg(500 * time.Millisecond); err == nil {
		fmt.Printf("  ✓ Account B received its own message: %s\n", string(msg.Data))
	}
	subB2.Unsubscribe()
	
	fmt.Println("\n=== Account Isolation Demo Complete ===")
}

// DemoAccountExports demonstrates exporting streams and services
func DemoAccountExports() {
	fmt.Println("\n=== Account Export/Import Demo ===")
	
	connA, err := nats.Connect("nats://user_a:pass_a@localhost:4226")
	if err != nil {
		log.Printf("Account A connection failed: %v", err)
		return
	}
	defer connA.Close()
	
	connB, err := nats.Connect("nats://user_b:pass_b@localhost:4226")
	if err != nil {
		log.Printf("Account B connection failed: %v", err)
		return
	}
	defer connB.Close()
	
	connC, err := nats.Connect("nats://user_c:pass_c@localhost:4226")
	if err != nil {
		log.Printf("Account C connection failed: %v", err)
		return
	}
	defer connC.Close()
	
	// 1. Public Stream Export/Import
	fmt.Println("\n1. Testing Public Stream Export (puba.>):")
	
	// Account C subscribes to imported stream (with prefix)
	subC, err := connC.SubscribeSync("from_a.puba.events")
	if err != nil {
		log.Printf("Account C subscribe failed: %v", err)
		return
	}
	
	// Account A publishes to public stream
	fmt.Println("  Account A publishing to 'puba.events'...")
	connA.Publish("puba.events", []byte("Public event from A"))
	connA.Flush()
	time.Sleep(200 * time.Millisecond)
	
	// Account C should receive it (imported with prefix)
	if msg, err := subC.NextMsg(500 * time.Millisecond); err == nil {
		fmt.Printf("  ✓ Account C received: %s\n", string(msg.Data))
		fmt.Printf("    (Imported as 'from_a.puba.events' - note the prefix)\n")
	} else {
		fmt.Printf("  ✗ Account C did not receive message: %v\n", err)
	}
	subC.Unsubscribe()
	
	// 2. Private Stream Export/Import
	fmt.Println("\n2. Testing Private Stream Export (b.> - only for Account B):")
	
	// Account B subscribes to private imported stream
	subB, err := connB.SubscribeSync("b.data")
	if err != nil {
		log.Printf("Account B subscribe failed: %v", err)
		return
	}
	
	// Account A publishes to private stream
	fmt.Println("  Account A publishing to 'b.data'...")
	connA.Publish("b.data", []byte("Private data for B"))
	connA.Flush()
	time.Sleep(200 * time.Millisecond)
	
	// Account B should receive it
	if msg, err := subB.NextMsg(500 * time.Millisecond); err == nil {
		fmt.Printf("  ✓ Account B received: %s\n", string(msg.Data))
		fmt.Println("    (Private stream - only Account B can import this)")
	} else {
		fmt.Printf("  ✗ Account B did not receive message: %v\n", err)
	}
	subB.Unsubscribe()
	
	// Account C cannot access private stream meant for B
	subC2, err := connC.SubscribeSync("b.data")
	if err != nil {
		log.Printf("Account C subscribe failed: %v", err)
	} else {
		connA.Publish("b.data", []byte("Should not reach C"))
		connA.Flush()
		time.Sleep(200 * time.Millisecond)
		
		if _, err := subC2.NextMsg(500 * time.Millisecond); err == nats.ErrTimeout {
			fmt.Println("  ✓ Account C correctly cannot access private stream for B")
		}
		subC2.Unsubscribe()
	}
	
	// 3. Public Service Export/Import with Remapping
	fmt.Println("\n3. Testing Public Service Export with Remapping:")
	
	// Account A sets up service responder
	_, err = connA.Subscribe("pubq.C", func(msg *nats.Msg) {
		fmt.Printf("  Account A service received request: %s\n", string(msg.Data))
		msg.Respond([]byte("Response from A's service"))
	})
	if err != nil {
		log.Printf("Account A service setup failed: %v", err)
		return
	}
	
	// Account C makes request using remapped subject 'Q' (maps to pubq.C)
	fmt.Println("  Account C making request to 'Q' (remapped to 'pubq.C')...")
	resp, err := connC.Request("Q", []byte("Request from C"), 2*time.Second)
	if err != nil {
		log.Printf("  Account C request failed: %v", err)
	} else {
		fmt.Printf("  ✓ Account C received response: %s\n", string(resp.Data))
		fmt.Println("    (Subject remapping: C publishes to 'Q', A receives on 'pubq.C')")
	}
	
	// 4. Private Service Export/Import
	fmt.Println("\n4. Testing Private Service Export (q.b - only for Account B):")
	
	// Account A sets up private service
	_, err = connA.Subscribe("q.b", func(msg *nats.Msg) {
		fmt.Printf("  Account A private service received request: %s\n", string(msg.Data))
		msg.Respond([]byte("Private response for B"))
	})
	if err != nil {
		log.Printf("Account A private service setup failed: %v", err)
		return
	}
	
	// Account B makes request to private service
	fmt.Println("  Account B making request to 'q.b'...")
	resp, err = connB.Request("q.b", []byte("Request from B"), 2*time.Second)
	if err != nil {
		log.Printf("  Account B request failed: %v", err)
	} else {
		fmt.Printf("  ✓ Account B received response: %s\n", string(resp.Data))
		fmt.Println("    (Private service - only Account B can access)")
	}
	
	fmt.Println("\n=== Account Export/Import Demo Complete ===")
}

// DemoNoAuthUser demonstrates the no_auth_user feature
func DemoNoAuthUser() {
	fmt.Println("\n=== No Auth User Demo ===")
	
	// Connect without credentials
	fmt.Println("\n1. Connecting without credentials (uses no_auth_user):")
	noAuthConn, err := nats.Connect("nats://localhost:4226")
	if err != nil {
		log.Printf("No-auth connection failed: %v", err)
		return
	}
	defer noAuthConn.Close()
	
	fmt.Println("  ✓ Connected successfully without credentials")
	fmt.Println("    (Automatically assigned to user_a in Account A)")
	
	// Should be able to publish to public exports from Account A
	fmt.Println("\n2. Testing access as Account A user:")
	if err := noAuthConn.Publish("puba.test", []byte("Message from no-auth user")); err != nil {
		log.Printf("  No-auth publish failed: %v", err)
	} else {
		fmt.Println("  ✓ Successfully published to 'puba.test'")
		fmt.Println("    (Has same permissions as user_a in Account A)")
	}
	
	fmt.Println("\n=== No Auth User Demo Complete ===")
}
