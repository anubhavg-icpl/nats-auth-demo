# NATS Authorization & Multi-Tenancy Demo

A comprehensive Go project demonstrating NATS authorization patterns and multi-tenancy features based on official NATS documentation.

## 📋 Features

This project demonstrates:

1. **Basic Authorization**
   - User-based permissions (admin, client, service roles)
   - Default permissions for users
   - Subject-level publish/subscribe controls

2. **Allow/Deny Rules**
   - Explicit allow and deny lists
   - Read-only user patterns
   - Subject wildcard permissions

3. **Allow Responses**
   - Service responders with temporary reply permissions
   - Single vs streaming response patterns
   - Time-limited response permissions

4. **Queue Permissions**
   - Queue-specific authorization
   - Queue group restrictions
   - Load balancing across queue members

5. **Account Isolation**
   - Multi-tenancy with complete account isolation
   - Independent subject namespaces per account
   - Secure tenant separation

6. **Account Exports/Imports**
   - Public and private stream exports
   - Public and private service exports
   - Subject remapping and prefixing
   - Cross-account communication

7. **No Auth User**
   - Default account assignment for unauthenticated clients
   - Simplified connection patterns

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or higher
- NATS Server 2.10+ ([Download](https://github.com/nats-io/nats-server/releases))

### Installation

1. Clone or navigate to the project:
```bash
cd /home/anubhavg/Desktop/nats
```

2. Install dependencies:
```bash
go mod download
```

3. Build the demo application:
```bash
go build -o nats-demo cmd/main.go
```

### Running the Demos

Each demo requires a NATS server running with a specific configuration file.

#### Option 1: Interactive Menu

Run the demo application:
```bash
./nats-demo
```

The interactive menu will guide you through each demo and prompt you to start the appropriate NATS server configuration.

#### Option 2: Manual Execution

1. Start NATS server with desired configuration:
```bash
# For basic authorization demo
nats-server -c config/basic-auth.conf

# For allow/deny demo
nats-server -c config/allow-deny.conf

# For allow responses demo
nats-server -c config/allow-responses.conf

# For queue permissions demo
nats-server -c config/queue-permissions.conf

# For accounts demo (demos 5, 6, 7)
nats-server -c config/accounts.conf
```

2. Run the demo application and select the corresponding demo from the menu.

## 📚 Demo Details

### 1. Basic Authorization (Port 4222)

**Config:** `config/basic-auth.conf`

Demonstrates:
- Admin user with full access (`admin:admin123`)
- Client user with request permissions (`client:client123`)
- Service user with response permissions (`service:service123`)
- Other user with default permissions (`other:other123`)

**Key Concepts:**
- Variable-based permission sets
- Request/response patterns
- Default permissions

### 2. Allow/Deny Rules (Port 4223)

**Config:** `config/allow-deny.conf`

Demonstrates:
- Explicit allow lists for publish/subscribe
- Explicit deny lists to block specific subjects
- Read-only user (can only subscribe)
- Limited user with mixed permissions

**Key Concepts:**
- Deny takes precedence over allow
- Wildcard subject matching
- Subject carving patterns

### 3. Allow Responses (Port 4224)

**Config:** `config/allow-responses.conf`

Demonstrates:
- Service with single response permission
- Service with streaming responses (max 5, 1m expiry)
- Service with mixed permissions (explicit + responses)

**Key Concepts:**
- Temporary publish permissions for replies
- Response limits and expiration
- Service responder patterns

### 4. Queue Permissions (Port 4225)

**Config:** `config/queue-permissions.conf`

Demonstrates:
- Queue-only subscriptions
- Queue group restrictions
- Environment-based queue separation (dev vs prod)
- Load balancing across queue members

**Key Concepts:**
- Queue-specific authorization
- Wildcard queue group matching
- Message distribution patterns

### 5. Account Isolation (Port 4226)

**Config:** `config/accounts.conf`

Demonstrates:
- Complete isolation between accounts
- Independent subject namespaces
- Account A, B, and C with different users

**Key Concepts:**
- Multi-tenancy
- Subject namespace isolation
- Account boundaries

### 6. Account Exports/Imports (Port 4226)

**Config:** `config/accounts.conf`

Demonstrates:
- Public stream export (`puba.>`)
- Private stream export (`b.>` - only for Account B)
- Public service export (`pubq.>`)
- Private service export (`q.b` - only for Account B)
- Subject remapping (prefix and to)

**Key Concepts:**
- Stream vs service exports
- Public vs private exports
- Subject transformation
- Cross-account communication

### 7. No Auth User (Port 4226)

**Config:** `config/accounts.conf`

Demonstrates:
- Connecting without credentials
- Automatic assignment to default account
- Inheriting default account permissions

**Key Concepts:**
- Simplified authentication
- Default account setup
- Guest access patterns

## 📁 Project Structure

```
nats/
├── cmd/
│   └── main.go              # Interactive demo application
├── config/
│   ├── basic-auth.conf      # Basic authorization config
│   ├── allow-deny.conf      # Allow/deny rules config
│   ├── allow-responses.conf # Allow responses config
│   ├── queue-permissions.conf # Queue permissions config
│   └── accounts.conf        # Multi-tenancy accounts config
├── examples/
│   ├── basic_auth.go        # Basic authorization demo
│   ├── allow_deny.go        # Allow/deny demo
│   ├── allow_responses.go   # Allow responses demo
│   ├── queue_permissions.go # Queue permissions demo
│   └── accounts.go          # Accounts and exports/imports demo
├── go.mod
└── README.md
```

## 🔑 User Credentials Reference

### Basic Auth Server (port 4222)
- `admin:admin123` - Full access
- `client:client123` - Requestor permissions
- `service:service123` - Responder permissions
- `other:other123` - Default permissions

### Allow/Deny Server (port 4223)
- `admin:admin123` - Full access
- `limited:limited123` - Limited with allow/deny rules
- `readonly:readonly123` - Read-only access

### Allow Responses Server (port 4224)
- `client:client123` - Request maker
- `service_single:service123` - Single response service
- `service_stream:service456` - Streaming response service
- `service_mixed:service789` - Mixed permissions service

### Queue Permissions Server (port 4225)
- `queue_only:queue123` - Queue-only subscriber
- `queue_restricted:queue456` - Restricted queue access

### Accounts Server (port 4226)
- `user_a:pass_a` - Account A
- `user_b:pass_b` - Account B
- `user_c:pass_c` - Account C

## 🐳 Docker Support

Start NATS server with Docker:

```bash
# Basic authorization
docker run -p 4222:4222 -v $(pwd)/config:/config nats:latest -c /config/basic-auth.conf

# Accounts
docker run -p 4226:4226 -v $(pwd)/config:/config nats:latest -c /config/accounts.conf
```

Or use docker-compose:

```bash
docker-compose up
```

## 📖 Documentation References

- [NATS Authorization](https://docs.nats.io/running-a-nats-service/configuration/securing_nats/authorization)
- [Multi-Tenancy using Accounts](https://docs.nats.io/running-a-nats-service/configuration/securing_nats/accounts)
- [NATS Go Client](https://github.com/nats-io/nats.go)
- [NATS by Example](https://natsbyexample.com)

## 🔍 Testing

The demos include built-in verification:
- ✓ indicates successful operations
- ✗ indicates correctly denied operations (demonstrating authorization working)

Watch the output carefully to understand how each authorization rule works.

## 🛠️ Troubleshooting

### Connection Refused
- Ensure NATS server is running with the correct config file
- Check the port matches the demo (each config uses a different port)
- Verify no firewall is blocking the port

### Permission Denied Errors
- This is expected! The demos show both successful and denied operations
- Denied operations demonstrate that authorization is working correctly

### Authentication Timeout
- Ensure you're using the correct username/password
- Check the config file has the user defined
- For no_auth_user, ensure the user exists in the account

## 📝 License

MIT License - feel free to use this code for learning and demonstration purposes.

## 🤝 Contributing

This is a demonstration project. Feel free to:
- Add more demo scenarios
- Improve the examples
- Fix bugs or issues
- Enhance documentation

## 🙏 Acknowledgments

Based on official NATS documentation:
- NATS.io team for excellent documentation
- NATS community for examples and best practices
