package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/anubhavg-icpl/nats-auth-demo/examples"
)

func main() {
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║      NATS Authorization & Multi-Tenancy Demo                 ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n┌────────────────────────────────────────────────────────────┐")
		fmt.Println("│ Select a demo to run:                                      │")
		fmt.Println("├────────────────────────────────────────────────────────────┤")
		fmt.Println("│  1. Basic Authorization                                    │")
		fmt.Println("│     - Admin, Client, Service, and Default permissions      │")
		fmt.Println("│     - Server: localhost:4222                               │")
		fmt.Println("│     - Config: config/basic-auth.conf                       │")
		fmt.Println("│                                                            │")
		fmt.Println("│  2. Allow/Deny Rules                                       │")
		fmt.Println("│     - Explicit allow and deny lists                        │")
		fmt.Println("│     - Read-only user example                               │")
		fmt.Println("│     - Server: localhost:4223                               │")
		fmt.Println("│     - Config: config/allow-deny.conf                       │")
		fmt.Println("│                                                            │")
		fmt.Println("│  3. Allow Responses                                        │")
		fmt.Println("│     - Service responders with reply permissions            │")
		fmt.Println("│     - Single vs streaming responses                        │")
		fmt.Println("│     - Server: localhost:4224                               │")
		fmt.Println("│     - Config: config/allow-responses.conf                  │")
		fmt.Println("│                                                            │")
		fmt.Println("│  4. Queue Permissions                                      │")
		fmt.Println("│     - Queue-specific authorization                         │")
		fmt.Println("│     - Load balancing across queue members                  │")
		fmt.Println("│     - Server: localhost:4225                               │")
		fmt.Println("│     - Config: config/queue-permissions.conf                │")
		fmt.Println("│                                                            │")
		fmt.Println("│  5. Account Isolation                                      │")
		fmt.Println("│     - Multi-tenancy with accounts                          │")
		fmt.Println("│     - Isolated communication contexts                      │")
		fmt.Println("│     - Server: localhost:4226                               │")
		fmt.Println("│     - Config: config/accounts.conf                         │")
		fmt.Println("│                                                            │")
		fmt.Println("│  6. Account Exports/Imports                                │")
		fmt.Println("│     - Public and private streams                           │")
		fmt.Println("│     - Public and private services                          │")
		fmt.Println("│     - Subject remapping                                    │")
		fmt.Println("│     - Server: localhost:4226                               │")
		fmt.Println("│     - Config: config/accounts.conf                         │")
		fmt.Println("│                                                            │")
		fmt.Println("│  7. No Auth User                                           │")
		fmt.Println("│     - Connecting without credentials                       │")
		fmt.Println("│     - Default account assignment                           │")
		fmt.Println("│     - Server: localhost:4226                               │")
		fmt.Println("│     - Config: config/accounts.conf                         │")
		fmt.Println("│                                                            │")
		fmt.Println("│  8. Run All Demos                                          │")
		fmt.Println("│                                                            │")
		fmt.Println("│  0. Exit                                                   │")
		fmt.Println("└────────────────────────────────────────────────────────────┘")
		fmt.Print("\nEnter your choice: ")

		input, _ := reader.ReadString('\n')
		choice := strings.TrimSpace(input)

		switch choice {
		case "1":
			fmt.Println("\n⚠️  Make sure NATS server is running with config/basic-auth.conf")
			fmt.Println("   Command: nats-server -c config/basic-auth.conf")
			fmt.Print("\nPress Enter to continue...")
			reader.ReadString('\n')
			examples.DemoBasicAuth()

		case "2":
			fmt.Println("\n⚠️  Make sure NATS server is running with config/allow-deny.conf")
			fmt.Println("   Command: nats-server -c config/allow-deny.conf")
			fmt.Print("\nPress Enter to continue...")
			reader.ReadString('\n')
			examples.DemoAllowDeny()

		case "3":
			fmt.Println("\n⚠️  Make sure NATS server is running with config/allow-responses.conf")
			fmt.Println("   Command: nats-server -c config/allow-responses.conf")
			fmt.Print("\nPress Enter to continue...")
			reader.ReadString('\n')
			examples.DemoAllowResponses()

		case "4":
			fmt.Println("\n⚠️  Make sure NATS server is running with config/queue-permissions.conf")
			fmt.Println("   Command: nats-server -c config/queue-permissions.conf")
			fmt.Print("\nPress Enter to continue...")
			reader.ReadString('\n')
			examples.DemoQueuePermissions()

		case "5":
			fmt.Println("\n⚠️  Make sure NATS server is running with config/accounts.conf")
			fmt.Println("   Command: nats-server -c config/accounts.conf")
			fmt.Print("\nPress Enter to continue...")
			reader.ReadString('\n')
			examples.DemoAccounts()

		case "6":
			fmt.Println("\n⚠️  Make sure NATS server is running with config/accounts.conf")
			fmt.Println("   Command: nats-server -c config/accounts.conf")
			fmt.Print("\nPress Enter to continue...")
			reader.ReadString('\n')
			examples.DemoAccountExports()

		case "7":
			fmt.Println("\n⚠️  Make sure NATS server is running with config/accounts.conf")
			fmt.Println("   Command: nats-server -c config/accounts.conf")
			fmt.Print("\nPress Enter to continue...")
			reader.ReadString('\n')
			examples.DemoNoAuthUser()

		case "8":
			fmt.Println("\n⚠️  This will run all demos. Make sure you start each NATS server")
			fmt.Println("   configuration as prompted.")
			fmt.Print("\nPress Enter to continue...")
			reader.ReadString('\n')

			fmt.Println("\n" + strings.Repeat("=", 64))
			fmt.Println("Starting Demo 1: Basic Authorization")
			fmt.Println(strings.Repeat("=", 64))
			fmt.Println("Start server: nats-server -c config/basic-auth.conf")
			fmt.Print("Press Enter when ready...")
			reader.ReadString('\n')
			examples.DemoBasicAuth()

			fmt.Println("\n" + strings.Repeat("=", 64))
			fmt.Println("Starting Demo 2: Allow/Deny Rules")
			fmt.Println(strings.Repeat("=", 64))
			fmt.Println("Start server: nats-server -c config/allow-deny.conf")
			fmt.Print("Press Enter when ready...")
			reader.ReadString('\n')
			examples.DemoAllowDeny()

			fmt.Println("\n" + strings.Repeat("=", 64))
			fmt.Println("Starting Demo 3: Allow Responses")
			fmt.Println(strings.Repeat("=", 64))
			fmt.Println("Start server: nats-server -c config/allow-responses.conf")
			fmt.Print("Press Enter when ready...")
			reader.ReadString('\n')
			examples.DemoAllowResponses()

			fmt.Println("\n" + strings.Repeat("=", 64))
			fmt.Println("Starting Demo 4: Queue Permissions")
			fmt.Println(strings.Repeat("=", 64))
			fmt.Println("Start server: nats-server -c config/queue-permissions.conf")
			fmt.Print("Press Enter when ready...")
			reader.ReadString('\n')
			examples.DemoQueuePermissions()

			fmt.Println("\n" + strings.Repeat("=", 64))
			fmt.Println("Starting Demo 5-7: Account Features")
			fmt.Println(strings.Repeat("=", 64))
			fmt.Println("Start server: nats-server -c config/accounts.conf")
			fmt.Print("Press Enter when ready...")
			reader.ReadString('\n')
			examples.DemoAccounts()
			examples.DemoAccountExports()
			examples.DemoNoAuthUser()

			fmt.Println("\n" + strings.Repeat("=", 64))
			fmt.Println("All demos completed!")
			fmt.Println(strings.Repeat("=", 64))

		case "0":
			fmt.Println("\nExiting... Goodbye!")
			return

		default:
			fmt.Println("\n❌ Invalid choice. Please try again.")
		}

		fmt.Print("\nPress Enter to return to menu...")
		reader.ReadString('\n')
	}
}
