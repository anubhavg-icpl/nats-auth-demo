package examples

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

// DemoAllowResponses demonstrates service responders with temporary reply permissions
func DemoAllowResponses() {
	fmt.Println("\n=== Allow Responses Demo ===")
	
	// Client that makes requests
	clientConn, err := nats.Connect("nats://client:client123@localhost:4224")
	if err != nil {
		log.Printf("Client connection failed: %v", err)
		return
	}
	defer clientConn.Close()
	
	// Service with single response permission
	fmt.Println("\n1. Testing Service with Single Response Permission:")
	serviceSingleConn, err := nats.Connect("nats://service_single:service123@localhost:4224")
	if err != nil {
		log.Printf("Service single connection failed: %v", err)
		return
	}
	defer serviceSingleConn.Close()
	
	// Setup service handler
	_, err = serviceSingleConn.Subscribe("requests.single", func(msg *nats.Msg) {
		fmt.Printf("  Service received request on 'requests.single'\n")
		
		// Service can respond once
		if err := msg.Respond([]byte("Single response")); err != nil {
			log.Printf("  Service response failed: %v", err)
		} else {
			fmt.Println("  ✓ Service sent single response")
		}
		
		// Trying to respond again should fail
		if err := msg.Respond([]byte("Second response")); err != nil {
			fmt.Printf("  ✗ Service correctly denied second response: %v\n", err)
		}
	})
	if err != nil {
		log.Printf("Service subscribe failed: %v", err)
		return
	}
	
	// Client makes request
	fmt.Println("  Client making request to 'requests.single'...")
	resp, err := clientConn.Request("requests.single", []byte("Request 1"), 2*time.Second)
	if err != nil {
		log.Printf("  Client request failed: %v", err)
	} else {
		fmt.Printf("  ✓ Client received response: %s\n", string(resp.Data))
	}
	
	// Service with stream response permission
	fmt.Println("\n2. Testing Service with Stream Response Permission (max 5, 1m expiry):")
	serviceStreamConn, err := nats.Connect("nats://service_stream:service456@localhost:4224")
	if err != nil {
		log.Printf("Service stream connection failed: %v", err)
		return
	}
	defer serviceStreamConn.Close()
	
	responseCount := 0
	_, err = serviceStreamConn.Subscribe("requests.stream", func(msg *nats.Msg) {
		fmt.Printf("  Service received request on 'requests.stream'\n")
		
		// Service can respond up to 5 times
		for i := 1; i <= 6; i++ {
			time.Sleep(100 * time.Millisecond)
			responseMsg := fmt.Sprintf("Response %d", i)
			
			if err := serviceSingleConn.Publish(msg.Reply, []byte(responseMsg)); err != nil {
				fmt.Printf("  ✗ Response %d failed (expected after 5): %v\n", i, err)
				break
			} else {
				responseCount++
				fmt.Printf("  ✓ Service sent response %d\n", i)
			}
		}
	})
	if err != nil {
		log.Printf("Service stream subscribe failed: %v", err)
		return
	}
	
	// Client makes request and receives multiple responses
	fmt.Println("  Client making request to 'requests.stream'...")
	inbox := nats.NewInbox()
	sub, err := clientConn.SubscribeSync(inbox)
	if err != nil {
		log.Printf("  Client inbox subscribe failed: %v", err)
	} else {
		// Publish request
		if err := clientConn.PublishRequest("requests.stream", inbox, []byte("Stream request")); err != nil {
			log.Printf("  Client request failed: %v", err)
		} else {
			// Receive responses
			time.Sleep(1 * time.Second)
			fmt.Printf("  Client checking for responses...\n")
			for i := 0; i < responseCount; i++ {
				if msg, err := sub.NextMsg(500 * time.Millisecond); err == nil {
					fmt.Printf("  ✓ Client received: %s\n", string(msg.Data))
				}
			}
		}
		sub.Unsubscribe()
	}
	
	// Service with mixed permissions
	fmt.Println("\n3. Testing Service with Mixed Permissions:")
	serviceMixedConn, err := nats.Connect("nats://service_mixed:service789@localhost:4224")
	if err != nil {
		log.Printf("Service mixed connection failed: %v", err)
		return
	}
	defer serviceMixedConn.Close()
	
	_, err = serviceMixedConn.Subscribe("requests.mixed", func(msg *nats.Msg) {
		fmt.Printf("  Service received request on 'requests.mixed'\n")
		
		// Can publish to logs (explicit permission)
		if err := serviceMixedConn.Publish("logs.service", []byte("Log entry")); err != nil {
			log.Printf("  Service log publish failed: %v", err)
		} else {
			fmt.Println("  ✓ Service published to 'logs.service'")
		}
		
		// Can also respond to request (allow_responses)
		if err := msg.Respond([]byte("Mixed response")); err != nil {
			log.Printf("  Service response failed: %v", err)
		} else {
			fmt.Println("  ✓ Service sent response")
		}
	})
	if err != nil {
		log.Printf("Service mixed subscribe failed: %v", err)
		return
	}
	
	fmt.Println("  Client making request to 'requests.mixed'...")
	resp, err = clientConn.Request("requests.mixed", []byte("Mixed request"), 2*time.Second)
	if err != nil {
		log.Printf("  Client request failed: %v", err)
	} else {
		fmt.Printf("  ✓ Client received response: %s\n", string(resp.Data))
	}
	
	fmt.Println("\n=== Allow Responses Demo Complete ===")
}
