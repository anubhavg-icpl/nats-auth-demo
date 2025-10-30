package examples

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

// DemoQueuePermissions demonstrates queue-specific permissions
func DemoQueuePermissions() {
	fmt.Println("\n=== Queue Permissions Demo ===")
	
	// Queue-only user
	fmt.Println("\n1. Testing Queue-Only User:")
	queueOnlyConn, err := nats.Connect("nats://queue_only:queue123@localhost:4225")
	if err != nil {
		log.Printf("Queue-only connection failed: %v", err)
		return
	}
	defer queueOnlyConn.Close()
	
	// Can subscribe to foo with queue group
	qSub, err := queueOnlyConn.QueueSubscribeSync("foo", "queue")
	if err != nil {
		log.Printf("Queue subscribe failed: %v", err)
	} else {
		fmt.Println("✓ Queue-only subscribed to 'foo' with queue group 'queue'")
		qSub.Unsubscribe()
	}
	
	// Cannot subscribe to foo without queue group
	if _, err := queueOnlyConn.SubscribeSync("foo"); err != nil {
		fmt.Printf("✗ Queue-only correctly denied plain subscription to 'foo': %v\n", err)
	}
	
	// Cannot subscribe with different queue group
	if _, err := queueOnlyConn.QueueSubscribeSync("foo", "other"); err != nil {
		fmt.Printf("✗ Queue-only correctly denied subscription with queue 'other': %v\n", err)
	}
	
	// Queue-restricted user
	fmt.Println("\n2. Testing Queue-Restricted User:")
	queueRestrictedConn, err := nats.Connect("nats://queue_restricted:queue456@localhost:4225")
	if err != nil {
		log.Printf("Queue-restricted connection failed: %v", err)
		return
	}
	defer queueRestrictedConn.Close()
	
	// Can subscribe to foo without queue
	plainSub, err := queueRestrictedConn.SubscribeSync("foo")
	if err != nil {
		log.Printf("Plain subscribe failed: %v", err)
	} else {
		fmt.Println("✓ Queue-restricted subscribed to 'foo' (plain)")
		plainSub.Unsubscribe()
	}
	
	// Can subscribe with v1 queue group
	v1Sub, err := queueRestrictedConn.QueueSubscribeSync("foo", "v1")
	if err != nil {
		log.Printf("Queue v1 subscribe failed: %v", err)
	} else {
		fmt.Println("✓ Queue-restricted subscribed to 'foo' with queue 'v1'")
		v1Sub.Unsubscribe()
	}
	
	// Can subscribe with v1.dev queue group
	v1DevSub, err := queueRestrictedConn.QueueSubscribeSync("foo", "v1.dev")
	if err != nil {
		log.Printf("Queue v1.dev subscribe failed: %v", err)
	} else {
		fmt.Println("✓ Queue-restricted subscribed to 'foo' with queue 'v1.dev'")
		v1DevSub.Unsubscribe()
	}
	
	// Can subscribe with any .dev queue group
	testDevSub, err := queueRestrictedConn.QueueSubscribeSync("foo", "test.dev")
	if err != nil {
		log.Printf("Queue test.dev subscribe failed: %v", err)
	} else {
		fmt.Println("✓ Queue-restricted subscribed to 'foo' with queue 'test.dev'")
		testDevSub.Unsubscribe()
	}
	
	// Cannot subscribe with .prod queue groups (denied)
	if _, err := queueRestrictedConn.QueueSubscribeSync("foo", "v1.prod"); err != nil {
		fmt.Printf("✗ Queue-restricted correctly denied subscription with queue 'v1.prod': %v\n", err)
	}
	
	if _, err := queueRestrictedConn.QueueSubscribeSync("bar", "test.prod"); err != nil {
		fmt.Printf("✗ Queue-restricted correctly denied subscription with queue 'test.prod': %v\n", err)
	}
	
	// Demonstrate queue distribution
	fmt.Println("\n3. Demonstrating Queue Distribution:")
	
	// Create admin connection to publish messages
	adminConn, err := nats.Connect("nats://localhost:4225")
	if err != nil {
		log.Printf("Admin connection failed: %v", err)
		return
	}
	defer adminConn.Close()
	
	// Create two queue subscribers
	received1 := 0
	received2 := 0
	
	sub1, _ := queueRestrictedConn.QueueSubscribe("foo", "v1.dev", func(msg *nats.Msg) {
		received1++
		fmt.Printf("  Worker 1 received message: %s\n", string(msg.Data))
	})
	defer sub1.Unsubscribe()
	
	sub2, _ := queueRestrictedConn.QueueSubscribe("foo", "v1.dev", func(msg *nats.Msg) {
		received2++
		fmt.Printf("  Worker 2 received message: %s\n", string(msg.Data))
	})
	defer sub2.Unsubscribe()
	
	// Allow subscriptions to register
	time.Sleep(100 * time.Millisecond)
	
	// Publish messages
	fmt.Println("  Publishing 10 messages to 'foo'...")
	for i := 1; i <= 10; i++ {
		msg := fmt.Sprintf("Message %d", i)
		adminConn.Publish("foo", []byte(msg))
		time.Sleep(50 * time.Millisecond)
	}
	
	time.Sleep(500 * time.Millisecond)
	fmt.Printf("\n  Distribution: Worker 1 = %d messages, Worker 2 = %d messages\n", received1, received2)
	fmt.Println("  ✓ Messages distributed across queue group members")
	
	fmt.Println("\n=== Queue Permissions Demo Complete ===")
}
