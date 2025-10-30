package examples

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

// DemoAllowDeny demonstrates explicit allow and deny rules
func DemoAllowDeny() {
	fmt.Println("\n=== Allow/Deny Authorization Demo ===")
	
	// Limited user - can publish to public and events, but not events.private
	fmt.Println("\n1. Testing Limited User:")
	limitedConn, err := nats.Connect("nats://limited:limited123@localhost:4223")
	if err != nil {
		log.Printf("Limited connection failed: %v", err)
		return
	}
	defer limitedConn.Close()
	
	// Can publish to public subjects
	if err := limitedConn.Publish("public.news", []byte("Public message")); err != nil {
		log.Printf("Limited publish to public.news failed: %v", err)
	} else {
		fmt.Println("✓ Limited published to 'public.news'")
	}
	
	// Can publish to events subjects
	if err := limitedConn.Publish("events.user.login", []byte("Event message")); err != nil {
		log.Printf("Limited publish to events.user.login failed: %v", err)
	} else {
		fmt.Println("✓ Limited published to 'events.user.login'")
	}
	
	// Cannot publish to events.private (explicitly denied)
	if err := limitedConn.Publish("events.private", []byte("Should fail")); err != nil {
		fmt.Printf("✗ Limited correctly denied publishing to 'events.private': %v\n", err)
	}
	
	// Can subscribe to allowed subjects
	sub, err := limitedConn.SubscribeSync("client.notifications")
	if err != nil {
		log.Printf("Limited subscribe to client.notifications failed: %v", err)
	} else {
		fmt.Println("✓ Limited subscribed to 'client.notifications'")
		sub.Unsubscribe()
	}
	
	// Cannot subscribe to disallowed subjects
	if _, err := limitedConn.SubscribeSync("admin.commands"); err != nil {
		fmt.Printf("✗ Limited correctly denied subscribing to 'admin.commands': %v\n", err)
	}
	
	// Read-only user - can only subscribe
	fmt.Println("\n2. Testing Read-Only User:")
	readonlyConn, err := nats.Connect("nats://readonly:readonly123@localhost:4223")
	if err != nil {
		log.Printf("Readonly connection failed: %v", err)
		return
	}
	defer readonlyConn.Close()
	
	// Can subscribe to any subject
	sub2, err := readonlyConn.SubscribeSync("any.subject.here")
	if err != nil {
		log.Printf("Readonly subscribe failed: %v", err)
	} else {
		fmt.Println("✓ Readonly subscribed to 'any.subject.here'")
		sub2.Unsubscribe()
	}
	
	// Cannot publish to any subject
	if err := readonlyConn.Publish("any.subject", []byte("Should fail")); err != nil {
		fmt.Printf("✗ Readonly correctly denied publishing: %v\n", err)
	}
	
	// Admin user - full access
	fmt.Println("\n3. Testing Admin User:")
	adminConn, err := nats.Connect("nats://admin:admin123@localhost:4223")
	if err != nil {
		log.Printf("Admin connection failed: %v", err)
		return
	}
	defer adminConn.Close()
	
	fmt.Println("✓ Admin has full publish/subscribe access to all subjects")
	
	fmt.Println("\n=== Allow/Deny Authorization Demo Complete ===")
}
