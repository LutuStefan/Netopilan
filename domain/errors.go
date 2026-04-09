package domain

import "errors"

var (
	ErrNotEnoughEnergy = errors.New("not enough energy")
	ErrAlreadyAtFinish = errors.New("already at the final gate")
	ErrTokenDead       = errors.New("token has no health remaining")
	ErrTokenBlocked    = errors.New("token is blocked by a Duplicate Demon — use 'identify' to prove your identity")
)
