# Netopiland

A CLI game that simulates the journey of a payment token through the Path of Authorization.

## Getting Started

### With Docker (recommended)

```bash
# Terminal 1: build and start
docker compose up --build

# Terminal 2: attach to play
docker attach netopiland
```

To detach without stopping: `Ctrl+P` then `Ctrl+Q`

### Without Docker

```bash
go run .
```

## Running Tests

```bash
go test ./tests/... -v
```

## Game Rules

You are a **payment token** traveling through 5 zones. Your goal is to reach the Issuer Throne with low risk to get **APPROVED**.

### The Path

| Zone | Description |
|------|-------------|
| Merchant Gate | Starting point. Echoes of previous requests linger. |
| Gateway Bridge | A fragile bridge. Messages may be lost or duplicated. |
| Risk Engine Woods | A dark forest that analyzes your behavior and history. |
| Acquirer Pass | A mountain pass where winds shift constantly. |
| Issuer Throne | The final destination. Approved or Declined. |

### Token Stats

| Stat | Range | Default | Description |
|------|-------|---------|-------------|
| Health | 0-100 | 100 | If it reaches 0, your token is DECLINED |
| Energy | 0-100 | 100 | Actions cost energy. Use `wait` to restore |
| Resistance | 0-100 | 30 | Defensive stat |
| Risk Score | 0-100 | 0 | If above 30 at the Issuer Throne, you are DECLINED |

### Actions

| Action | Cost | Description |
|--------|------|-------------|
| `move` | 5 energy | Advance to the next zone |
| `scan` | 5 energy | Peek at the next zone |
| `shield` | 20 energy | Protective barrier for 2 turns |
| `identify` | 10 energy | Clear a Duplicate Demon block |
| `wait` | free | Rest and restore 15 energy |
| `status` | free | Display current state |
| `journal` | free | Review your journey log |
| `help` | free | Show available actions |
| `quit` | free | End the game |

### Creatures

Each zone has a 50% chance of spawning a creature:

| Creature | Chance | Effect |
|----------|--------|--------|
| Fraudster | 31.3% | Attacks 1-3 times. Each: Risk +5, Health -8 |
| Duplicate Demon | 31.3% | Blocks all actions except `identify`, `status`, `help`, `quit` |
| Timeout Spirit | 31.3% | Risk +5 to +10, Health -5 to -15 |
| Decline Guardian | 6% | 30% chance to end the game immediately with DECLINED |

### Win/Lose Conditions

- **APPROVED**: Reach the Issuer Throne with risk score <= 30
- **DECLINED**: Risk score > 30 at the Issuer Throne, health reaches 0, or Decline Guardian triggers
