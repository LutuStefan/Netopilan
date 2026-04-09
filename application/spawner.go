package application

import (
	"math/rand"

	"netopiland/domain"
)

// CreatureWeight pairs a creature with its spawn weight (relative probability).
type CreatureWeight struct {
	Creature domain.Challenge
	Weight   float64
}

// Spawner decides if and which creature appears when entering a gate.
type Spawner struct {
	SpawnChance float64
	Creatures   []CreatureWeight
}

// NewSpawner creates a Spawner with the given spawn chance and weighted creature pool.
func NewSpawner(spawnChance float64, creatures []CreatureWeight) *Spawner {
	return &Spawner{
		SpawnChance: spawnChance,
		Creatures:   creatures,
	}
}

// Spawn rolls for a creature. Returns the creature and true if one spawns,
// or nil and false otherwise.
func (s *Spawner) Spawn() (domain.Challenge, bool) {
	if rand.Float64() >= s.SpawnChance {
		return nil, false
	}
	return s.pickCreature(), true
}

func (s *Spawner) pickCreature() domain.Challenge {
	var totalWeight float64
	for _, cw := range s.Creatures {
		totalWeight += cw.Weight
	}

	roll := rand.Float64() * totalWeight
	var cumulative float64
	for _, cw := range s.Creatures {
		cumulative += cw.Weight
		if roll < cumulative {
			return cw.Creature
		}
	}

	return s.Creatures[0].Creature
}
