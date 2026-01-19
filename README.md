# ClipP2P

A peer-to-peer universal clipboard for your terminal. Copy on one machine, paste on another.


## Features

- **Zero Config** - Automatically finds peers on your local network
- **Real-time Sync** - Copy text on one device, instantly available on other peers
- **Bubbletea Dashboard** - See connection status and sync history
- **Encrypted** - All traffic encrypted via libp2p's secure channels

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/owenHochwald/clipp2p.git
cd clipp2p

# Build
make build

# (Optional) Install to your PATH
sudo mv clipp2p /usr/local/bin/
# Or for user-only install:
mv clipp2p ~/.local/bin/
```

## Usage

```bash
# Run the app
clipp2p

# Or if not in PATH
./clipp2p
```

### Keyboard Controls

| Key | Action |
|-----|--------|
| `q` | Quit |
| `s` | Toggle sync on/off |
| `c` | Clear history |

### Multi-Device Setup

1. Run `clipp2p` on each device connected to the same local network
2. Devices automatically discover each other via mDNS
3. Copy text on any device - it syncs to all connected peers

## How It Works

```
┌──────────────────────────────────────────────────────────────┐
│                        Your Network                          │
│                                                              │
│   ┌─────────┐      mDNS Discovery      ┌─────────┐           │
│   │ Machine │◄────────────────────────►│ Machine │           │
│   │    A    │                          │    B    │           │
│   └────┬────┘                          └────┬────┘           │
│        │                                    │                │
│        │         libp2p streams             │                │
│        │◄──────────────────────────────────►│                │
│        │      (encrypted JSON messages)     │                │ 
│                                                              │
└──────────────────────────────────────────────────────────────┘
```

**Flow:**
1. **Discovery** - mDNS broadcasts your node's presence on the local network
2. **Connection** - When another ClipP2P node is found, a TCP connection is established
3. **Watching** - Each node polls its local clipboard for changes (every 500ms)
4. **Sync** - When clipboard changes, the content is broadcast to all connected peers
5. **Write** - Receiving peers automatically update their local clipboard

## License

Apache 2.0
