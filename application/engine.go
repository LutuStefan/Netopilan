package application

import (
	"fmt"
	"strings"

	"netopiland/domain"
)

const (
	MoveCost              = 5
	ScanCost              = 5
	ShieldCost            = 20
	IdentifyCost          = 10
	WaitRestore           = 15
	ShieldDuration        = 2
	ApprovalRiskThreshold = 30
)

// Engine is the application service that orchestrates the game.
// It receives all dependencies via constructor injection.
type Engine struct {
	Token           *domain.Token
	Journal         *domain.Journal
	GateEffects     map[domain.Gate]domain.GateEffect
	ZoneEvents      map[domain.Gate]domain.ZoneEvent
	Spawner         *Spawner
	PendingCreature domain.Challenge
	Scanned         bool
	GameOver        bool
	Approved        bool
}

// NewEngine creates a game engine with the provided dependencies.
func NewEngine(
	token *domain.Token,
	journal *domain.Journal,
	gateEffects map[domain.Gate]domain.GateEffect,
	zoneEvents map[domain.Gate]domain.ZoneEvent,
	spawner *Spawner,
) *Engine {
	return &Engine{
		Token:       token,
		Journal:     journal,
		GateEffects: gateEffects,
		ZoneEvents:  zoneEvents,
		Spawner:     spawner,
	}
}

// Move advances the token to the next gate, applies gate effect and zone event.
// Costs 5 energy.
func (e *Engine) Move() (string, error) {
	if e.Token.Position.IsFinish() {
		return "", domain.ErrAlreadyAtFinish
	}
	if e.Token.Energy.Value < MoveCost {
		return "", domain.ErrNotEnoughEnergy
	}

	e.Token.Energy.Add(-MoveCost)

	nextGate, _ := e.Token.Position.Next()
	e.Token.Position = nextGate

	var lines []string
	lines = append(lines, fmt.Sprintf("You advance to %s.", nextGate))
	lines = append(lines, nextGate.Description())

	if effect, ok := e.GateEffects[nextGate]; ok {
		lines = append(lines, fmt.Sprintf("[Gate Effect] %s", effect.Description()))
		msg := effect.Apply(e.Token)
		lines = append(lines, msg)
	}

	if event, ok := e.ZoneEvents[nextGate]; ok {
		if msg, triggered := event.Trigger(e.Token); triggered {
			lines = append(lines, fmt.Sprintf("[Zone Event] %s", msg))
		}
	}

	var creature domain.Challenge
	if e.Scanned {
		creature = e.PendingCreature
	} else {
		creature, _ = e.Spawner.Spawn()
	}
	e.Scanned = false
	e.PendingCreature = nil

	if creature != nil {
		lines = append(lines, fmt.Sprintf("[Creature] A %s appears!", creature.Name()))
		result := creature.Execute(e.Token)
		lines = append(lines, result.Message)
		e.Journal.Record(nextGate, "encounter", creature.Name()+": "+result.Message)
		if result.GameOver {
			e.GameOver = true
		}
	}

	e.Token.TickShield()

	if e.Token.Health.Value <= 0 {
		e.GameOver = true
		lines = append(lines, "Your token has lost all health — DECLINED.")
		e.Journal.Record(nextGate, "death", "Token health reached 0 — declined")
	}

	if !e.GameOver && nextGate.IsFinish() {
		e.GameOver = true
		if e.Token.RiskScore.Value > ApprovalRiskThreshold {
			lines = append(lines, fmt.Sprintf("The Issuer Throne examines your risk score (%d/100)... DECLINED.", e.Token.RiskScore.Value))
			e.Journal.Record(nextGate, "judgment", "Token declined — risk too high")
		} else {
			e.Approved = true
			lines = append(lines, fmt.Sprintf("The Issuer Throne examines your risk score (%d/100)... APPROVED!", e.Token.RiskScore.Value))
			e.Journal.Record(nextGate, "judgment", "Token approved!")
		}
	}

	e.Journal.Record(nextGate, "move", "Advanced to "+nextGate.String())

	return strings.Join(lines, "\n"), nil
}

// Scan reveals information about the next gate. Costs 5 energy.
func (e *Engine) Scan() (string, error) {
	if e.Token.Position.IsFinish() {
		return "", domain.ErrAlreadyAtFinish
	}
	if e.Token.Energy.Value < ScanCost {
		return "", domain.ErrNotEnoughEnergy
	}

	e.Token.Energy.Add(-ScanCost)

	nextGate, _ := e.Token.Position.Next()

	var lines []string
	lines = append(lines, fmt.Sprintf("=== Scanning %s ===", nextGate))
	lines = append(lines, nextGate.Description())

	if effect, ok := e.GateEffects[nextGate]; ok {
		lines = append(lines, fmt.Sprintf("Gate Effect: %s", effect.Description()))
	}

	creature, spawned := e.Spawner.Spawn()
	e.PendingCreature = creature
	e.Scanned = true

	if spawned {
		lines = append(lines, fmt.Sprintf("Creature detected: %s — %s", creature.Name(), creature.Description()))
	} else {
		lines = append(lines, "No creature detected... for now.")
	}

	e.Journal.Record(e.Token.Position, "scan", "Scanned "+nextGate.String())

	return strings.Join(lines, "\n"), nil
}

// Shield activates a protective barrier for 2 turns. Costs 20 energy.
func (e *Engine) Shield() (string, error) {
	if e.Token.Energy.Value < ShieldCost {
		return "", domain.ErrNotEnoughEnergy
	}

	e.Token.Energy.Add(-ShieldCost)
	e.Token.ShieldTTL = ShieldDuration

	e.Journal.Record(e.Token.Position, "shield", fmt.Sprintf("Shield activated for %d turns", ShieldDuration))

	return fmt.Sprintf("Shield activated! You are protected for %d turns.", ShieldDuration), nil
}

// Identify clears the blocked state set by a Duplicate Demon. Costs 10 energy.
func (e *Engine) Identify() (string, error) {
	if e.Token.Energy.Value < IdentifyCost {
		return "", domain.ErrNotEnoughEnergy
	}

	if !e.Token.Blocked {
		return "You are not blocked. No need to identify.", nil
	}

	e.Token.Energy.Add(-IdentifyCost)
	e.Token.Blocked = false

	e.Journal.Record(e.Token.Position, "identify", "Identity confirmed — block cleared")

	return "Identity confirmed! The Duplicate Demon's grip fades. You are free to act.", nil
}

// Wait restores 15 energy and ticks the shield down.
func (e *Engine) Wait() string {
	e.Token.Energy.Add(WaitRestore)
	e.Token.TickShield()

	if e.Scanned {
		creature, spawned := e.Spawner.Spawn()
		if spawned {
			e.PendingCreature = creature
		} else {
			e.PendingCreature = nil
		}
	}

	e.Journal.Record(e.Token.Position, "wait", "Rested and restored energy")

	return fmt.Sprintf("You rest and recover energy. Energy +%d (now %d).", WaitRestore, e.Token.Energy.Value)
}

// JournalView returns a formatted view of the journey log.
func (e *Engine) JournalView() string {
	entries := e.Journal.Entries()
	if len(entries) == 0 {
		return "Your journal is empty. Your journey has just begun."
	}
	var lines []string
	lines = append(lines, "=== Journey Journal ===")
	for _, entry := range entries {
		lines = append(lines, entry.String())
	}
	return strings.Join(lines, "\n")
}
