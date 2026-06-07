package websocket

import (
	"sync"
	"time"
)

const (
	DefaultInitialTime  = 600 * 1000 // 10 minutes in milliseconds
	DefaultIncrement    = 5 * 1000  // 5 seconds in milliseconds
	ClockTickInterval   = 100 * time.Millisecond
	ClockUpdateInterval = 1 * time.Second
)

// Clock represents a game clock with Fischer-style increment.
type Clock struct {
	mu sync.Mutex

	whiteMs int64
	blackMs int64
	lastTick time.Time

	stopped bool
	gameOver bool
	timeoutColor string // "white" or "black" if timeout occurred

	stopCh chan struct{}
}

// NewClock creates a new game clock with the given initial time and increment.
func NewClock(initialMs, incrementMs int64) *Clock {
	return &Clock{
		whiteMs:  initialMs,
		blackMs:  initialMs,
		lastTick: time.Now(),
		stopCh:   make(chan struct{}),
	}
}

// NewDefaultClock creates a clock with standard 10min+5s settings.
func NewDefaultClock() *Clock {
	return NewClock(DefaultInitialTime, DefaultIncrement)
}

// WhiteMs returns white's remaining time in milliseconds.
func (c *Clock) WhiteMs() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.whiteMs
}

// BlackMs returns black's remaining time in milliseconds.
func (c *Clock) BlackMs() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.blackMs
}

// IsGameOver returns true if a timeout has occurred.
func (c *Clock) IsGameOver() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.gameOver
}

// TimeoutColor returns the color that timed out, or empty string.
func (c *Clock) TimeoutColor() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.timeoutColor
}

// SwitchPlayer switches which clock is running (called after a move).
// It adds the increment to the player who just moved.
func (c *Clock) SwitchPlayer(activeColor string, incrementMs int64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Tick the active player's clock first
	now := time.Now()
	elapsed := now.Sub(c.lastTick).Milliseconds()
	if activeColor == "white" {
		c.whiteMs -= elapsed
	} else {
		c.blackMs -= elapsed
	}
	c.lastTick = now

	// Add increment to the player who just moved (the active one before switch)
	if activeColor == "white" {
		c.whiteMs += incrementMs
	} else {
		c.blackMs += incrementMs
	}

	// Check for timeout
	if c.whiteMs <= 0 {
		c.whiteMs = 0
		c.gameOver = true
		c.timeoutColor = "white"
	}
	if c.blackMs <= 0 {
		c.blackMs = 0
		c.gameOver = true
		c.timeoutColor = "black"
	}
}

// Start begins ticking the clock for the given active color.
func (c *Clock) Start(activeColor string) {
	c.mu.Lock()
	if c.stopped {
		c.mu.Unlock()
		return
	}
	c.lastTick = time.Now()
	c.mu.Unlock()
}

// Stop stops the clock.
func (c *Clock) Stop() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.stopped = true
	close(c.stopCh)
}

// StopCh returns the stop channel.
func (c *Clock) StopCh() <-chan struct{} {
	return c.stopCh
}

// Tick decrements the active player's clock. Returns true if timeout occurred.
func (c *Clock) Tick(activeColor string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.stopped || c.gameOver {
		return false
	}

	now := time.Now()
	elapsed := now.Sub(c.lastTick).Milliseconds()
	c.lastTick = now

	if activeColor == "white" {
		c.whiteMs -= elapsed
		if c.whiteMs <= 0 {
			c.whiteMs = 0
			c.gameOver = true
			c.timeoutColor = "white"
			return true
		}
	} else {
		c.blackMs -= elapsed
		if c.blackMs <= 0 {
			c.blackMs = 0
			c.gameOver = true
			c.timeoutColor = "black"
			return true
		}
	}
	return false
}

// TickForTest forces a tick for testing (bypasses time check).
func (c *Clock) TickForTest(activeColor string, ms int64) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.gameOver {
		return false
	}

	if activeColor == "white" {
		c.whiteMs -= ms
		if c.whiteMs <= 0 {
			c.whiteMs = 0
			c.gameOver = true
			c.timeoutColor = "white"
			return true
		}
	} else {
		c.blackMs -= ms
		if c.blackMs <= 0 {
			c.blackMs = 0
			c.gameOver = true
			c.timeoutColor = "black"
			return true
		}
	}
	return false
}

// SetTime sets a player's time directly (for testing).
func (c *Clock) SetTime(color string, ms int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if color == "white" {
		c.whiteMs = ms
	} else {
		c.blackMs = ms
	}
}
