package creatures

import "netopiland/domain"

// DuplicateDemon copies your identity, blocking you until you prove who you are.
// Effect: sets Blocked=true. Only the identify action clears it.
type DuplicateDemon struct{}

func (d DuplicateDemon) Name() string { return "Duplicate Demon" }

func (d DuplicateDemon) Description() string {
	return "A creature that tries to copy your identity. After the encounter, there might be two entities claiming to be you."
}

func (d DuplicateDemon) Execute(t *domain.Token) domain.ChallengeResult {
	t.Blocked = true
	return domain.ChallengeResult{
		Message: "The Duplicate Demon copies your identity! You are blocked. Use 'identify' to prove who you are.",
		Passed:  false,
	}
}
