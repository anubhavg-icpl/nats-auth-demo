# Architecture Overview

## 🏛️ System Architecture

### Demo Application Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         Demo Application                        │
│                         (nats-demo)                             │
├─────────────────────────────────────────────────────────────────┤
│  ┌───────────────────────────────────────────────────────────┐ │
│  │                   Interactive Menu                        │ │
│  │                   (cmd/main.go)                           │ │
│  └───────────────────────┬───────────────────────────────────┘ │
│                          │                                       │
│                          ▼                                       │
│  ┌───────────────────────────────────────────────────────────┐ │
│  │                   Example Modules                         │ │
│  ├───────────────────────────────────────────────────────────┤ │
│  │  • basic_auth.go       - Authorization patterns           │ │
│  │  • allow_deny.go       - Allow/deny rules                 │ │
│  │  • allow_responses.go  - Service responders               │ │
│  │  • queue_permissions.go - Queue patterns                  │ │
│  │  • accounts.go         - Multi-tenancy                    │ │
│  └───────────────────────┬───────────────────────────────────┘ │
│                          │                                       │
└──────────────────────────┼───────────────────────────────────────┘
                           │
                           ▼ NATS Protocol
┌──────────────────────────────────────────────────────────────────┐
│                      NATS Server Layer                           │
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐            │
│  │ Basic Auth  │  │ Allow/Deny  │  │   Allow     │            │
│  │   Server    │  │   Server    │  │  Responses  │            │
│  │  :4222      │  │  :4223      │  │  :4224      │            │
│  └─────────────┘  └─────────────┘  └─────────────┘            │
│                                                                  │
│  ┌─────────────┐  ┌─────────────────────────────────────┐     │
│  │   Queue     │  │         Accounts Server             │     │
│  │Permissions  │  │  (Multi-Tenancy)                    │     │
│  │  :4225      │  │  :4226                              │     │
│  └─────────────┘  └─────────────────────────────────────┘     │
│                                                                  │
└──────────────────────────────────────────────────────────────────┘
```

## 🔐 Authorization Flow

### Basic Authorization Pattern

```
┌──────────┐                ┌──────────────┐                ┌──────────┐
│  Client  │                │ NATS Server  │                │ Subject  │
│          │                │              │                │  Space   │
└────┬─────┘                └──────┬───────┘                └────┬─────┘
     │                             │                             │
     │ 1. Connect with credentials │                             │
     ├────────────────────────────>│                             │
     │                             │                             │
     │                             │ 2. Validate credentials     │
     │                             │    & load permissions       │
     │                             │                             │
     │ 3. Connection established   │                             │
     │<────────────────────────────┤                             │
     │                             │                             │
     │ 4. Publish to "req.a"       │                             │
     ├────────────────────────────>│                             │
     │                             │                             │
     │                             │ 5. Check publish permissions│
     │                             │    for "req.a"              │
     │                             │                             │
     │                             │ 6. Forward if allowed       │
     │                             ├────────────────────────────>│
     │                             │                             │
     │ 7. Publish to "admin.cmd"   │                             │
     ├────────────────────────────>│                             │
     │                             │                             │
     │                             │ 8. Check permissions        │
     │                             │    (DENIED)                 │
     │                             │                             │
     │ 9. Permission denied error  │                             │
     │<────────────────────────────┤                             │
     │                             │                             │
```

## 🏢 Multi-Tenancy Architecture

### Account Isolation

```
┌─────────────────────────────────────────────────────────────────┐
│                         NATS Server                             │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌──────────────────────┐  ┌──────────────────────┐           │
│  │    Account A         │  │    Account B         │           │
│  ├──────────────────────┤  ├──────────────────────┤           │
│  │ Users:               │  │ Users:               │           │
│  │   • user_a           │  │   • user_b           │           │
│  │                      │  │                      │           │
│  │ Subject Space:       │  │ Subject Space:       │           │
│  │   • private.data     │  │   • private.data     │           │
│  │   • app.events       │  │   • app.events       │           │
│  │   • ...              │  │   • ...              │           │
│  │                      │  │                      │           │
│  │ Exports:             │  │ Imports:             │           │
│  │   • puba.> (public)  │  │   • b.> from A       │           │
│  │   • pubq.> (public)  │  │   • q.b from A       │           │
│  │   • b.> (private→B)  │  │                      │           │
│  │   • q.b (private→B)  │  │                      │           │
│  └──────────────────────┘  └──────────────────────┘           │
│           ▲                          ▲                         │
│           │                          │                         │
│           │    ISOLATED - No direct communication             │
│           │    (except via exports/imports)                   │
│           │                          │                         │
│           ▼                          ▼                         │
│  ┌──────────────────────────────────────────────┐             │
│  │    Account C                                 │             │
│  ├──────────────────────────────────────────────┤             │
│  │ Users:                                       │             │
│  │   • user_c                                   │             │
│  │                                              │             │
│  │ Subject Space:                               │             │
│  │   • private.data (different from A & B)      │             │
│  │   • from_a.puba.> (imported & prefixed)      │             │
│  │   • Q (remapped from pubq.C)                 │             │
│  │                                              │             │
│  │ Imports:                                     │             │
│  │   • puba.> from A (prefix: from_a)           │             │
│  │   • pubq.C from A (remap to: Q)              │             │
│  └──────────────────────────────────────────────┘             │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Export/Import Flow

```
Account A (Exporter)              Account C (Importer)
┌─────────────────┐              ┌─────────────────┐
│                 │              │                 │
│  Publisher      │              │  Subscriber     │
│    │            │              │      ▲          │
│    │ Publish to │              │      │ Receive  │
│    ▼ puba.event │              │      │ as       │
│  puba.event     │              │  from_a.puba.   │
│                 │              │      event      │
└────────┬────────┘              └────────▲────────┘
         │                                │
         │                                │
         ▼                                │
    ┌────────────────────────────────────────┐
    │      NATS Server                       │
    │                                        │
    │  1. Receive on puba.event in Account A │
    │  2. Check export rules (public)        │
    │  3. Forward to Account C imports       │
    │  4. Apply prefix: from_a               │
    │  5. Deliver as from_a.puba.event       │
    └────────────────────────────────────────┘

Service Export/Import (with remapping)

Account A (Service Provider)      Account C (Service Consumer)
┌─────────────────┐              ┌─────────────────┐
│                 │              │                 │
│  Service        │              │  Client         │
│  Responder      │              │                 │
│    ▲            │              │    │ Request to │
│    │ Listen on  │              │    ▼ Q          │
│pubq.C           │              │   Q             │
│    │            │              │                 │
│    │ Respond    │              │    ▲            │
│    │            │              │    │ Response   │
└────┼────────────┘              └────┼────────────┘
     │                                │
     ▼                                │
┌────────────────────────────────────────────────┐
│      NATS Server                               │
│                                                │
│  1. Client publishes to Q in Account C         │
│  2. Server remaps Q → pubq.C                   │
│  3. Forward to Account A's pubq.C              │
│  4. Service in A responds                      │
│  5. Response routed back to C's client         │
└────────────────────────────────────────────────┘
```

## 🔄 Request/Response Pattern

### Standard Request/Response

```
┌──────────┐                ┌──────────────┐                ┌──────────┐
│ Requestor│                │ NATS Server  │                │ Responder│
│ (Client) │                │              │                │ (Service)│
└────┬─────┘                └──────┬───────┘                └────┬─────┘
     │                             │                             │
     │ 1. Subscribe to req.a       │                             │
     │                             │<────────────────────────────┤
     │                             │                             │
     │ 2. Subscribe to _INBOX.xyz  │                             │
     ├────────────────────────────>│                             │
     │                             │                             │
     │ 3. Publish request          │                             │
     │    Subject: req.a           │                             │
     │    Reply-To: _INBOX.xyz     │                             │
     ├────────────────────────────>│                             │
     │                             │                             │
     │                             │ 4. Route to subscriber      │
     │                             ├────────────────────────────>│
     │                             │                             │
     │                             │                             │
     │                             │ 5. Publish response to      │
     │                             │    _INBOX.xyz               │
     │                             │<────────────────────────────┤
     │                             │                             │
     │ 6. Receive response         │                             │
     │<────────────────────────────┤                             │
     │                             │                             │
```

### Allow Responses Pattern

```
┌──────────┐                ┌──────────────┐                ┌──────────┐
│  Client  │                │ NATS Server  │                │ Service  │
│          │                │              │                │(w/allow_ │
│          │                │              │                │responses)│
└────┬─────┘                └──────┬───────┘                └────┬─────┘
     │                             │                             │
     │                             │ Service permissions:        │
     │                             │   subscribe: "requests.*"   │
     │                             │   allow_responses: true     │
     │                             │   (implicit deny publish)   │
     │                             │                             │
     │ 1. Subscribe to requests.*  │                             │
     │                             │<────────────────────────────┤
     │                             │                             │
     │ 2. Request                  │                             │
     │    Reply-To: _INBOX.abc123  │                             │
     ├────────────────────────────>│                             │
     │                             ├────────────────────────────>│
     │                             │                             │
     │                             │ 3. Temp permission granted  │
     │                             │    to publish to            │
     │                             │    _INBOX.abc123            │
     │                             │                             │
     │                             │ 4. Publish response         │
     │                             │<────────────────────────────┤
     │                             │                             │
     │ 5. Receive response         │                             │
     │<────────────────────────────┤                             │
     │                             │                             │
     │                             │ 6. Temp permission revoked  │
     │                             │                             │
     │                             │ 7. Try to publish again     │
     │                             │    (DENIED - no permission) │
     │                             │  X──────────────────────────┤
```

## 📊 Queue Groups

### Queue Distribution

```
                    ┌──────────────┐
                    │ NATS Server  │
                    └──────┬───────┘
                           │
          Publish to "tasks.process"
                           │
                           ▼
    ┌──────────────────────────────────────────┐
    │         Queue Group "workers"            │
    └──────────────────────────────────────────┘
              │              │              │
    ┌─────────┴────┐  ┌──────┴──────┐  ┌───┴────────┐
    │   Worker 1   │  │   Worker 2  │  │  Worker 3  │
    │              │  │             │  │            │
    │ Receives     │  │ Receives    │  │ Receives   │
    │ messages     │  │ messages    │  │ messages   │
    │ 1, 4, 7...   │  │ 2, 5, 8...  │  │ 3, 6, 9... │
    └──────────────┘  └─────────────┘  └────────────┘

    Messages are load balanced across queue members
```

### Queue Permissions

```
User: queue_restricted
Permissions:
  subscribe:
    allow: ["foo", "foo v1", "foo v1.>", "foo *.dev"]
    deny: ["> *.prod"]

┌─────────────────────────────────────────────────────┐
│               Subscription Attempts                 │
├─────────────────────────────────────────────────────┤
│ foo (plain)              ✓ Allowed                  │
│ foo:v1                   ✓ Allowed (queue v1)       │
│ foo:v1.dev               ✓ Allowed (matches v1.>)   │
│ foo:test.dev             ✓ Allowed (matches *.dev)  │
│ foo:v1.prod              ✗ Denied (matches *.prod)  │
│ bar:test.prod            ✗ Denied (matches *.prod)  │
│ foo:v2                   ✗ Denied (not in allow)    │
└─────────────────────────────────────────────────────┘
```

## 🔒 Permission Evaluation Order

```
When a client attempts to publish or subscribe:

1. Check if user is authenticated
   │
   ├─ NO → Reject (unless no_auth_user configured)
   │
   └─ YES ▼

2. Check if user has permissions defined
   │
   ├─ NO → Use default_permissions
   │
   └─ YES ▼

3. For the requested subject:
   │
   ├─ Check DENY list
   │  └─ Match? → REJECT
   │
   └─ Check ALLOW list
      └─ Match? → ACCEPT
      └─ No match? → REJECT (unless allow is empty)

Special cases:
• allow_responses: Grants temporary publish permission to reply subjects
• Queue permissions: Additional check for queue group name
• Account isolation: Check happens before permission check
```

## 📈 Scalability Considerations

```
┌─────────────────────────────────────────────────────┐
│              Production Deployment                  │
├─────────────────────────────────────────────────────┤
│                                                     │
│  ┌─────────────┐  ┌─────────────┐  ┌────────────┐ │
│  │   NATS      │  │   NATS      │  │   NATS     │ │
│  │  Server 1   │──│  Server 2   │──│  Server 3  │ │
│  │  (Cluster)  │  │  (Cluster)  │  │  (Cluster) │ │
│  └─────────────┘  └─────────────┘  └────────────┘ │
│         │                 │                │        │
│         └─────────────────┴────────────────┘        │
│                     │                               │
│         ┌───────────┴────────────┐                  │
│         │                        │                  │
│    ┌────▼────┐             ┌────▼────┐             │
│    │ Account │             │ Account │             │
│    │    A    │             │    B    │             │
│    │         │             │         │             │
│    │ 1000s   │             │ 1000s   │             │
│    │  of     │             │  of     │             │
│    │ clients │             │ clients │             │
│    └─────────┘             └─────────┘             │
│                                                     │
│  Key benefits:                                      │
│  • High availability through clustering             │
│  • Complete isolation between accounts              │
│  • No permission conflicts                          │
│  • Independent scaling per tenant                   │
│                                                     │
└─────────────────────────────────────────────────────┘
```

---

This architecture enables secure, scalable, and isolated multi-tenant messaging systems with fine-grained access control.
