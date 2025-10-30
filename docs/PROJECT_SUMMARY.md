# NATS Authorization & Multi-Tenancy Demo - Project Summary

## 🎯 Project Overview

This is a comprehensive Go demonstration project showcasing NATS authorization patterns and multi-tenancy features, built according to official NATS documentation.

## 📦 What's Included

### Configuration Files (5 total)
Located in `config/` directory:

1. **basic-auth.conf** (Port 4222)
   - Demonstrates user roles: admin, client, service, other
   - Shows default permissions
   - Request/response patterns

2. **allow-deny.conf** (Port 4223)
   - Explicit allow/deny lists
   - Read-only user pattern
   - Subject wildcard matching

3. **allow-responses.conf** (Port 4224)
   - Service responder patterns
   - Single vs streaming responses
   - Time-limited permissions

4. **queue-permissions.conf** (Port 4225)
   - Queue-specific authorization
   - Environment separation (dev/prod)
   - Queue group patterns

5. **accounts.conf** (Port 4226)
   - Multi-tenancy with 3 accounts (A, B, C)
   - Public and private exports/imports
   - Subject remapping
   - No-auth user configuration

### Go Examples (5 files)
Located in `examples/` directory:

1. **basic_auth.go** - Basic authorization demo
2. **allow_deny.go** - Allow/deny rules demo
3. **allow_responses.go** - Service responder demo
4. **queue_permissions.go** - Queue permissions demo
5. **accounts.go** - Account isolation and export/import demos

### Main Application
- **cmd/main.go** - Interactive menu-driven demo application

### Supporting Files
- **docker-compose.yml** - Multi-server Docker setup
- **Makefile** - Build and run automation
- **README.md** - Complete documentation
- **QUICKSTART.md** - Fast-start guide
- **.gitignore** - Git configuration
- **go.mod/go.sum** - Go dependencies

## 🚀 Quick Usage

### Using Docker (Recommended for beginners)
```bash
# Start all NATS servers
make docker-up

# Run the demo
./nats-demo

# Stop servers when done
make docker-down
```

### Using Local NATS Server
```bash
# Terminal 1: Start a NATS server
nats-server -c config/basic-auth.conf

# Terminal 2: Run demo
./nats-demo
# Select the matching demo from menu
```

### Using Makefile
```bash
make build          # Build the application
make run            # Build and run
make docker-up      # Start all servers
make docker-down    # Stop all servers
make clean          # Clean build artifacts
```

## 📊 Demo Coverage

### Authorization Features
- ✅ User-based permissions
- ✅ Role-based access control
- ✅ Default permissions
- ✅ Allow/deny lists
- ✅ Wildcard subject matching
- ✅ Request/response patterns
- ✅ Temporary reply permissions
- ✅ Queue-specific permissions
- ✅ Environment separation

### Multi-Tenancy Features
- ✅ Account isolation
- ✅ Public stream exports
- ✅ Private stream exports
- ✅ Public service exports
- ✅ Private service exports
- ✅ Subject prefixing
- ✅ Subject remapping
- ✅ Cross-account communication
- ✅ No-auth user setup

## 🔑 Key Concepts Demonstrated

### 1. Authorization Hierarchy
```
Admin (full access)
  ├── Can publish to any subject
  └── Can subscribe to any subject

Client (requestor)
  ├── Can publish to specific subjects
  └── Can subscribe to response subjects

Service (responder)
  ├── Can subscribe to request subjects
  └── Can publish responses

Other (default)
  ├── Can publish to SANDBOX.*
  └── Can subscribe to PUBLIC.> and _INBOX.>
```

### 2. Account Isolation
```
Account A
  └── Users in A can only see A's subjects

Account B
  └── Users in B can only see B's subjects

Account C
  └── Users in C can only see C's subjects

Communication between accounts requires explicit exports/imports
```

### 3. Export/Import Patterns
```
Public Export (any account can import)
  └── Stream: puba.>
  └── Service: pubq.>

Private Export (specific accounts only)
  └── Stream: b.> (only Account B)
  └── Service: q.b (only Account B)

Subject Remapping
  └── Import with prefix: puba.> → from_a.puba.>
  └── Import with mapping: pubq.C → Q
```

## 📝 User Credentials Cheat Sheet

### Basic Auth (4222)
- admin:admin123 - Full access
- client:client123 - Requestor
- service:service123 - Responder
- other:other123 - Default

### Allow/Deny (4223)
- admin:admin123 - Full access
- limited:limited123 - Allow/deny rules
- readonly:readonly123 - Read-only

### Allow Responses (4224)
- client:client123 - Request maker
- service_single:service123 - Single response
- service_stream:service456 - Stream response
- service_mixed:service789 - Mixed permissions

### Queue Permissions (4225)
- queue_only:queue123 - Queue-only
- queue_restricted:queue456 - Restricted

### Accounts (4226)
- user_a:pass_a - Account A
- user_b:pass_b - Account B
- user_c:pass_c - Account C
- (no credentials) - No-auth user → Account A

## 🎓 Learning Outcomes

After running these demos, you'll understand:

1. How to configure user-based authorization in NATS
2. How to use allow/deny rules for fine-grained control
3. How to implement request/response patterns securely
4. How to control queue group access
5. How to set up multi-tenant architectures
6. How to enable controlled cross-account communication
7. How to remap subjects for simpler client code
8. How to provide guest access with no-auth users

## 🏗️ Project Structure
```
nats/
├── cmd/
│   └── main.go                    # Interactive demo app
├── config/
│   ├── basic-auth.conf            # Basic authorization
│   ├── allow-deny.conf            # Allow/deny rules
│   ├── allow-responses.conf       # Service responders
│   ├── queue-permissions.conf     # Queue permissions
│   └── accounts.conf              # Multi-tenancy
├── examples/
│   ├── basic_auth.go              # Auth examples
│   ├── allow_deny.go              # Allow/deny examples
│   ├── allow_responses.go         # Service examples
│   ├── queue_permissions.go       # Queue examples
│   └── accounts.go                # Account examples
├── docker-compose.yml             # Multi-server setup
├── Makefile                       # Build automation
├── README.md                      # Full documentation
├── QUICKSTART.md                  # Quick start guide
├── PROJECT_SUMMARY.md             # This file
├── .gitignore                     # Git configuration
├── go.mod                         # Go dependencies
└── nats-demo                      # Built binary
```

## 🔗 Reference Documentation

- [NATS Authorization Docs](https://docs.nats.io/running-a-nats-service/configuration/securing_nats/authorization)
- [NATS Accounts Docs](https://docs.nats.io/running-a-nats-service/configuration/securing_nats/accounts)
- [NATS Go Client](https://github.com/nats-io/nats.go)
- [NATS by Example](https://natsbyexample.com)

## 💡 Tips for Success

1. **Start Simple**: Begin with Demo 1 (Basic Authorization)
2. **Read the Output**: Watch for ✓ (allowed) and ✗ (denied) operations
3. **Experiment**: Modify config files and see what breaks
4. **Use Docker**: Easiest way to run multiple servers
5. **Check Logs**: Use `make docker-logs` to debug issues
6. **Monitor**: Access monitoring endpoints (e.g., http://localhost:8222/varz)

## 🎯 Next Steps

1. Run all demos to understand each pattern
2. Modify configuration files to experiment
3. Create your own authorization patterns
4. Implement multi-tenancy in your applications
5. Explore JWT-based authentication (advanced)
6. Set up NATS clusters with authorization

## 📧 Support

- GitHub Issues: Report bugs or request features
- NATS Slack: Join the community
- Documentation: Read official NATS docs

---

**Built with ❤️ using NATS and Go**

Version: 1.0.0  
Last Updated: 2025-10-30
