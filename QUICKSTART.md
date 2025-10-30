# Quick Start Guide

## 🎯 Get Started in 3 Steps

### Step 1: Install NATS Server

Download and install NATS Server from [GitHub releases](https://github.com/nats-io/nats-server/releases).

**Linux/macOS:**
```bash
# Download latest release
curl -L https://github.com/nats-io/nats-server/releases/latest/download/nats-server-v2.10.5-linux-amd64.zip -o nats-server.zip

# Extract
unzip nats-server.zip

# Move to PATH
sudo mv nats-server-v2.10.5-linux-amd64/nats-server /usr/local/bin/

# Verify
nats-server --version
```

**Using Docker (Alternative):**
```bash
docker pull nats:latest
```

### Step 2: Build the Demo

```bash
cd /home/anubhavg/Desktop/nats

# Download dependencies and build
make build

# Or manually:
go mod download
go build -o nats-demo cmd/main.go
```

### Step 3: Run Your First Demo

**Option A: Using Docker Compose (Easiest)**
```bash
# Start all NATS servers
make docker-up

# Run the demo
./nats-demo
# Select demo #1 from the menu
```

**Option B: Using Local NATS Server**
```bash
# Terminal 1: Start NATS server
nats-server -c config/basic-auth.conf

# Terminal 2: Run demo
./nats-demo
# Select demo #1 from the menu
```

## 📖 What Each Demo Shows

### Demo 1: Basic Authorization
Learn how to set up users with different permission levels.
- ✅ Admin with full access
- ✅ Client with limited publish rights
- ✅ Service with response permissions
- ✅ Default permissions for new users

### Demo 2: Allow/Deny Rules
See explicit permission control in action.
- ✅ Whitelist specific subjects
- ✅ Blacklist sensitive subjects
- ✅ Read-only access patterns

### Demo 3: Allow Responses
Understand service responder patterns.
- ✅ One-time reply permissions
- ✅ Streaming responses with limits
- ✅ Mixed permission patterns

### Demo 4: Queue Permissions
Control queue group access.
- ✅ Queue-specific authorization
- ✅ Environment separation (dev/prod)
- ✅ Load balancing demonstrations

### Demo 5-7: Multi-Tenancy
Master account isolation and cross-account communication.
- ✅ Complete tenant isolation
- ✅ Public and private exports
- ✅ Subject remapping
- ✅ Guest access patterns

## 🎓 Learning Path

1. **Start with Demo 1** - Understand basic permissions
2. **Try Demo 2** - Learn allow/deny patterns  
3. **Explore Demo 3** - Master service patterns
4. **Test Demo 4** - Understand queue groups
5. **Master Demo 5-7** - Learn multi-tenancy

## 🔧 Common Commands

```bash
# Build
make build

# Run
make run

# Start all servers with Docker
make docker-up

# Stop all servers
make docker-down

# View logs
make docker-logs

# Clean build artifacts
make clean

# Format code
make fmt
```

## 🐛 Troubleshooting

**"Connection refused"**
- Make sure NATS server is running: `ps aux | grep nats-server`
- Check the correct port for the demo you're running
- Use `make docker-up` for easy multi-server setup

**"Permission denied" errors**
- This is normal! The demos show both successful and denied operations
- ✗ marks mean authorization is working correctly

**Build errors**
```bash
go mod tidy
go build -o nats-demo cmd/main.go
```

## 📚 Next Steps

- Read the [full README](README.md) for detailed documentation
- Check NATS config files in `config/` directory
- Review Go code in `examples/` directory
- Experiment with modifying permissions
- Try creating your own authorization patterns

## 🎉 Success Indicators

You'll know it's working when you see:
- ✓ Green checkmarks for allowed operations
- ✗ Red X's for correctly denied operations
- Messages flowing between accounts
- Queue distribution across workers

Happy learning! 🚀
