package examples

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

// DemoBasicAuth demonstrates basic authorization with different user roles
func DemoBasicAuth() {
	fmt.Println("\n=== Basic Authorization Demo ===")
	
	// Admin user - has full access
	fmt.Println("\n1. Testing Admin User (full access):")
	adminConn, err := nats.Connect("nats://admin:admin123@localhost:4222")
	if err != nil {
		log.Printf("Admin connection failed: %v", err)
		return
	}
	defer adminConn.Close()
	
	// Admin can publish anywhere
	if err := adminConn.Publish("any.subject", []byte("Admin message")); err != nil {
		log.Printf("Admin publish failed: %v", err)
	} else {
		fmt.Println("✓ Admin published to 'any.subject'")
	}
	
	// Admin can subscribe anywhere
	sub, err := adminConn.SubscribeSync("any.subject")
	if err != nil {
		log.Printf("Admin subscribe failed: %v", err)
	} else {
		fmt.Println("✓ Admin subscribed to 'any.subject'")
		sub.Unsubscribe()
	}
	
	// Client user - requestor role
	fmt.Println("\n2. Testing Client User (requestor role):")
	clientConn, err := nats.Connect("nats://client:client123@localhost:4222")
	if err != nil {
		log.Printf("Client connection failed: %v", err)
		return
	}
	defer clientConn.Close()
	
	// Client can publish to request subjects
	if err := clientConn.Publish("req.a", []byte("Request message")); err != nil {
		log.Printf("Client publish to req.a failed: %v", err)
	} else {
		fmt.Println("✓ Client published to 'req.a'")
	}
	
	// Client cannot publish to other subjects
	if err := clientConn.Publish("other.subject", []byte("Should fail")); err != nil {
		fmt.Printf("✗ Client correctly denied publishing to 'other.subject': %v\n", err)
	}
	
	// Client can subscribe to inbox (for responses)
	inboxSub, err := clientConn.SubscribeSync("_INBOX.>")
	if err != nil {
		log.Printf("Client subscribe to _INBOX failed: %v", err)
	} else {
		fmt.Println("✓ Client subscribed to '_INBOX.>'")
		inboxSub.Unsubscribe()
	}
	
	// Service user - responder role
	fmt.Println("\n3. Testing Service User (responder role):")
	serviceConn, err := nats.Connect("nats://service:service123@localhost:4222")
	if err != nil {
		log.Printf("Service connection failed: %v", err)
		return
	}
	defer serviceConn.Close()
	
	// Service can subscribe to request subjects
	reqSub, err := serviceConn.SubscribeSync("req.a")
	if err != nil {
		log.Printf("Service subscribe to req.a failed: %v", err)
	} else {
		fmt.Println("✓ Service subscribed to 'req.a'")
		reqSub.Unsubscribe()
	}
	
	// Service can publish to inbox (responses)
	if err := serviceConn.Publish("_INBOX.test123", []byte("Response message")); err != nil {
		log.Printf("Service publish to _INBOX failed: %v", err)
	} else {
		fmt.Println("✓ Service published to '_INBOX.test123'")
	}
	
	// Other user - default permissions
	fmt.Println("\n4. Testing Other User (default permissions):")
	otherConn, err := nats.Connect("nats://other:other123@localhost:4222")
	if err != nil {
		log.Printf("Other connection failed: %v", err)
		return
	}
	defer otherConn.Close()
	
	// Other can publish to SANDBOX subjects
	if err := otherConn.Publish("SANDBOX.test", []byte("Sandbox message")); err != nil {
		log.Printf("Other publish to SANDBOX failed: %v", err)
	} else {
		fmt.Println("✓ Other published to 'SANDBOX.test'")
	}
	
	// Other can subscribe to PUBLIC subjects
	pubSub, err := otherConn.SubscribeSync("PUBLIC.announcements")
	if err != nil {
		log.Printf("Other subscribe to PUBLIC failed: %v", err)
	} else {
		fmt.Println("✓ Other subscribed to 'PUBLIC.announcements'")
		pubSub.Unsubscribe()
	}
	
	fmt.Println("\n=== Basic Authorization Demo Complete ===")
}
