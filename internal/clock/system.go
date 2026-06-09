// Package clock provides the system-backed Clock implementation.
package clock

import (
	"time"

	"github.com/develonaut/todui/internal/ports"
)

// System is a ports.Clock backed by the operating system clock.
type System struct{}

// Now returns the current local time.
func (System) Now() time.Time { return time.Now() }

var _ ports.Clock = System{}
