# NKeys Authentication in NATS

## Overview

NKeys is a highly secure, Ed25519-based public-key signature system for NATS authentication. Unlike traditional password-based authentication, NKeys use asymmetric cryptography where:

- **Private keys (seeds)** never leave the client
- **Public keys** are stored on the server
- Each connection requires signing a unique challenge
- **Immune to replay attacks** due to random server-generated challenges

## Key Advantages

✅ **No Password Storage**: Server only stores public keys, not secrets  
✅ **Challenge-Response**: Each connection uses a unique, random challenge  
✅ **Ed25519 Cryptography**: Faster and more secure than RSA  
✅ **Zero Trust**: Client must prove possession of private key for each connection  
✅ **Replay Attack Prevention**: Signatures can't be reused  

## Architecture

```
┌─────────────┐                    ┌─────────────┐
│   Client    │                    │   Server    │
│             │                    │             │
│  Seed (SK)  │                    │ PubKey (PK) │
└─────────────┘                    └─────────────┘
       │                                  │
       │  1. Connect Request              │
       │─────────────────────────────────>│
       │                                  │
       │  2. Random Challenge (nonce)     │
       │<─────────────────────────────────│
       │                                  │
       │  3. Sign(challenge, SK)          │
       │─────────────────────────────────>│
       │                                  │
       │     4. Verify(sig, challenge, PK)│
       │                                  │
       │  5. Connection Established ✓     │
       │<─────────────────────────────────│
```

## Key Types

NKeys use prefixes to identify key types:

- **`S`** - Seed (private key) - e.g., `SUACSSL3UAHUDXKFSNVUZRF5...`
- **`U`** - User (public key) - e.g., `UDXU4RCSJNZOIQHZNWXHXORDP...`
- **`A`** - Account (for multi-tenancy)
- **`N`** - Server
- **`C`** - Cluster
- **`O`** - Operator

For authentication, we primarily use **User** keys (U/SU prefix).

## Quick Start

### 1. Generate NKeys

Run the demo application:
```bash
./nats-demo
# Select option 9 - Generate NKeys
```

Or use the Go code directly:
```go
package main

import (
    "fmt"
    "github.com/nats-io/nkeys"
)

func main() {
    // Generate user key pair
    kp, _ := nkeys.CreateUser()
    
    // Get seed (private key) - keep secret!
    seed, _ := kp.Seed()
    fmt.Println("Seed:", string(seed))
    
    // Get public key - share with server
    pub, _ := kp.PublicKey()
    fmt.Println("Public:", pub)
}
```

### 2. Configure NATS Server

Create a configuration file (`config/nkeys-auth.conf`):

```conf
port: 4227

authorization {
  users = [
    {
      nkey: "UDXU4RCSJNZOIQHZNWXHXORDPRTGNJAHAHFRGZNEEJCPQTT2M7NLCNF4"
      permissions: {
        publish = ">"
        subscribe = ">"
      }
    }
  ]
}
```

Start the server:
```bash
nats-server -c config/nkeys-auth.conf
```

### 3. Connect with Client

```go
package main

import (
    "github.com/nats-io/nats.go"
    "github.com/nats-io/nkeys"
)

func main() {
    seed := "SUACSSL3UAHUDXKFSNVUZRF5UHPMWZ6BFDTJ7M6USDXIEDNPPQYYYCU3VY"
    
    // Create NKey option
    opt := nats.Nkey("UDXU4RCSJNZOIQHZNWXHXORDPRTGNJAHAHFRGZNEEJCPQTT2M7NLCNF4", 
        func(nonce []byte) ([]byte, error) {
            kp, _ := nkeys.FromSeed([]byte(seed))
            return kp.Sign(nonce)
        })
    
    // Connect
    nc, _ := nats.Connect("nats://localhost:4227", opt)
    defer nc.Close()
    
    // Use connection normally
    nc.Publish("test.subject", []byte("Hello NATS!"))
}
```

## Implementation Details

### Key Generation Process

1. **Create User Key Pair**:
   ```go
   kp, err := nkeys.CreateUser()
   ```

2. **Extract Seed (Private Key)**:
   ```go
   seed, err := kp.Seed()
   // Returns: SUACSSL3UAHUDXKFSNVUZRF5UHPMWZ6BFDTJ7M6USDXIEDNPPQYYYCU3VY
   ```

3. **Extract Public Key**:
   ```go
   publicKey, err := kp.PublicKey()
   // Returns: UDXU4RCSJNZOIQHZNWXHXORDPRTGNJAHAHFRGZNEEJCPQTT2M7NLCNF4
   ```

### Challenge-Response Flow

1. **Client initiates connection** with public key
2. **Server generates random challenge** (nonce)
3. **Client signs challenge** with private key:
   ```go
   kp, _ := nkeys.FromSeed([]byte(seed))
   signature, _ := kp.Sign(challenge)
   ```
4. **Server verifies signature**:
   ```go
   kp, _ := nkeys.FromPublicKey(publicKey)
   err := kp.Verify(challenge, signature)
   ```
5. **Connection established** if verification succeeds

### Security Properties

- **One-way function**: Cannot derive private key from public key
- **Challenge uniqueness**: Each connection uses a new random challenge
- **Signature binding**: Signature is valid only for specific challenge
- **No credential storage**: Private keys stay on client, never transmitted
- **Fast verification**: Ed25519 is optimized for speed

## Usage Examples

### Example 1: Basic Authentication

```go
func connectWithNKey(seed, publicKey string) (*nats.Conn, error) {
    opt := nats.Nkey(publicKey, func(nonce []byte) ([]byte, error) {
        kp, err := nkeys.FromSeed([]byte(seed))
        if err != nil {
            return nil, err
        }
        return kp.Sign(nonce)
    })
    
    return nats.Connect("nats://localhost:4227", opt)
}
```

### Example 2: With Permissions

Server config:
```conf
authorization {
  ADMIN = {
    publish = ">"
    subscribe = ">"
  }
  
  REQUESTOR = {
    publish = ["req.*"]
    subscribe = "_INBOX.>"
  }
  
  users = [
    {nkey: "UDXU4...", permissions: $ADMIN}
    {nkey: "UAH42...", permissions: $REQUESTOR}
  ]
}
```

### Example 3: Request-Response Pattern

```go
// Service (responder)
serviceNC, _ := connectWithNKey(serviceSeed, servicePubKey)
serviceNC.Subscribe("req.*", func(m *nats.Msg) {
    m.Respond([]byte("response"))
})

// Client (requester)
clientNC, _ := connectWithNKey(clientSeed, clientPubKey)
msg, _ := clientNC.Request("req.test", []byte("request"), time.Second)
fmt.Println("Response:", string(msg.Data))
```

## Running the Demos

The project includes comprehensive demos:

```bash
./nats-demo
```

Choose from:
- **Option 8**: NKeys Authentication Demo
  - Tests all user roles
  - Demonstrates permissions
  - Shows request-response patterns

- **Option 9**: Generate NKeys
  - **Option a**: Generate and display keys
  - **Option b**: Generate and save to files with server config

## File Structure

```
nats/
├── config/
│   └── nkeys-auth.conf          # NKeys server configuration
├── examples/
│   ├── nkeys_utils.go           # Key generation utilities
│   ├── nkeys_auth.go            # Authentication examples
│   └── nkeys_keygen.go          # Key generation with file export
├── generated/                   # Auto-generated keys (gitignored)
│   ├── nkeys.txt               # Generated key pairs
│   └── nkeys-server.conf       # Generated server config
└── docs/
    └── NKEYS_AUTHENTICATION.md  # This file
```

## Best Practices

### Security

1. **Never commit seeds to version control**:
   ```bash
   # Add to .gitignore
   generated/
   *.seed
   *_seed.txt
   ```

2. **Store seeds securely**:
   - Use environment variables
   - Use secret management systems (Vault, AWS Secrets Manager)
   - Encrypt at rest

3. **Rotate keys regularly**:
   - Generate new key pairs periodically
   - Update server configuration
   - Distribute new seeds to clients

### Key Management

1. **Separate keys per environment**:
   ```
   keys/
   ├── dev/
   ├── staging/
   └── production/
   ```

2. **Use descriptive naming**:
   ```
   admin_user_seed.txt
   api_service_seed.txt
   client_app_seed.txt
   ```

3. **Backup seeds securely**:
   - Encrypted backups
   - Multiple secure locations
   - Access audit logs

### Development vs Production

**Development**:
```go
// OK for development - seed in code
seed := "SUACSSL3UAHUDXKFSNVUZRF5..."
```

**Production**:
```go
// Load from environment or secret store
seed := os.Getenv("NATS_NKEY_SEED")
// or
seed, err := secretManager.GetSecret("nats/user/seed")
```

## Comparison with Other Auth Methods

| Feature | Password | Token | NKeys |
|---------|----------|-------|-------|
| Server stores secrets | Yes | Yes | No |
| Vulnerable to replay | Yes | Yes | No |
| Challenge-response | No | No | Yes |
| Key rotation complexity | Low | Low | Medium |
| Security level | Medium | Medium | High |
| Performance | Fast | Fast | Very Fast |
| Credential transmission | Yes | Yes | No |

## Troubleshooting

### Connection Refused

**Problem**: Client can't connect  
**Solution**: 
1. Verify server is running: `ps aux | grep nats-server`
2. Check port: `netstat -an | grep 4227`
3. Verify public key matches seed

### Permission Denied

**Problem**: "Permissions Violation for Publish/Subscribe"  
**Solution**:
1. Check server config has correct public key
2. Verify permissions in server config
3. Ensure subject matches allowed patterns

### Invalid Signature

**Problem**: "Authentication Violation - Signature"  
**Solution**:
1. Verify seed and public key are from same pair
2. Check seed format (should start with 'SU')
3. Ensure nkeys library version matches server

### Key Mismatch

**Problem**: Public key doesn't match seed  
**Solution**:
```go
// Verify key pair
kp, _ := nkeys.FromSeed([]byte(seed))
pub, _ := kp.PublicKey()
fmt.Println("Expected:", expectedPubKey)
fmt.Println("Actual:", pub)
```

## References

- [NATS NKeys Documentation](https://docs.nats.io/running-a-nats-service/configuration/securing_nats/auth_intro/nkey_auth)
- [Ed25519 Signature Scheme](https://ed25519.cr.yp.to/)
- [nkeys Go Library](https://github.com/nats-io/nkeys)
- [NATS Security Best Practices](https://docs.nats.io/running-a-nats-service/configuration/securing_nats)

## Support

For issues or questions:
1. Check the [NATS Documentation](https://docs.nats.io)
2. Visit [NATS Slack Community](https://slack.nats.io)
3. Review [GitHub Issues](https://github.com/nats-io/nats-server/issues)
