package clock

import "time"

type Clocker interface {
	Now() time.Time
}

type RealClocker struct{}

func (c *RealClocker) Now() time.Time {
	return time.Now()
}

type FixedClocker struct{}

func (c *FixedClocker) Now() time.Time {
	return time.Date(2022, 5, 10, 12, 24, 56, 0, time.UTC)
}
