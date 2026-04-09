package domain

// ClampedValue is a numeric value constrained between Min and Max.
// All mutations go through Add, which enforces the bounds.
type ClampedValue struct {
	Value int
	Min   int
	Max   int
}

// NewClampedValue creates a ClampedValue with the given bounds and initial value.
// The initial value is clamped to [min, max].
func NewClampedValue(initial, min, max int) ClampedValue {
	return ClampedValue{
		Value: clamp(initial, min, max),
		Min:   min,
		Max:   max,
	}
}

// Add adjusts the value by delta and clamps the result to [Min, Max].
func (c *ClampedValue) Add(delta int) {
	c.Value = clamp(c.Value+delta, c.Min, c.Max)
}

func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
