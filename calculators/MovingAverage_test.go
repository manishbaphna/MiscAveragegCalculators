package calculators

import (
	"context"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

// TestMovingAverage tests the MovingAverage function.
func TestMovingAverage(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a channel to send Tick values.
	ticks := make(chan Tick)
	go func() {
		defer close(ticks)
		ticks <- Tick{Price: decimal.NewFromFloat(10), Timestamp: time.Now()}
		ticks <- Tick{Price: decimal.NewFromFloat(20), Timestamp: time.Now()}
		ticks <- Tick{Price: decimal.NewFromFloat(30), Timestamp: time.Now()}
		ticks <- Tick{Price: decimal.NewFromFloat(40), Timestamp: time.Now()}
		ticks <- Tick{Price: decimal.NewFromFloat(50), Timestamp: time.Now()}
	}()

	// Call the MovingAverage function.
	movingAvgCh := MovingAverage(ctx, ticks, 3)
	var results []decimal.Decimal
	for avg := range movingAvgCh {
		results = append(results, avg.Round(2))
	}

	// Define the expected results.
	expected := []decimal.Decimal{
		decimal.RequireFromString("10").Round(2),
		decimal.RequireFromString("15").Round(2),
		decimal.RequireFromString("20").Round(2),
		decimal.RequireFromString("30").Round(2),
		decimal.RequireFromString("40").Round(2),
	}

	// Assert that the results match the expected values.
	assert.Equal(t, expected, results)
}

// This tests if the function correctly handles the context cancellation while processing the ticks.
// It ensures that the function stops processing further ticks once the context is cancelled.
func TestMovingAverageWithContextCancellationWithinTicksChannel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	ticks := make(chan Tick)
	go func() {
		defer close(ticks)
		ticks <- Tick{Price: decimal.NewFromFloat(10), Timestamp: time.Now()}
		ticks <- Tick{Price: decimal.NewFromFloat(20), Timestamp: time.Now()}
		ticks <- Tick{Price: decimal.NewFromFloat(30), Timestamp: time.Now()}
		cancel()
		ticks <- Tick{Price: decimal.NewFromFloat(40), Timestamp: time.Now()}
	}()

	movingAvgCh := MovingAverage(ctx, ticks, 3)
	var results []decimal.Decimal
	for avg := range movingAvgCh {
		results = append(results, avg.Round(2))
	}

	expected := []decimal.Decimal{
		decimal.RequireFromString("10").Round(2),
		decimal.RequireFromString("15").Round(2),
		decimal.RequireFromString("20").Round(2),
	}

	assert.Equal(t, expected, results)
}

// This tests if the function correctly handles the context cancellation from the consumer side.
// It ensures that the function stops sending values to the channel once the context is cancelled.
func TestMovingAverageWithContextCancellationWithinOutChannel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	ticks := make(chan Tick)
	go func() {
		defer close(ticks)
		ticks <- Tick{Price: decimal.NewFromFloat(10), Timestamp: time.Now()}
		ticks <- Tick{Price: decimal.NewFromFloat(20), Timestamp: time.Now()}
		ticks <- Tick{Price: decimal.NewFromFloat(30), Timestamp: time.Now()}
		ticks <- Tick{Price: decimal.NewFromFloat(40), Timestamp: time.Now()}
		ticks <- Tick{Price: decimal.NewFromFloat(50), Timestamp: time.Now()}
	}()

	movingAvgCh := MovingAverage(ctx, ticks, 3)
	var results []decimal.Decimal
	for avg := range movingAvgCh {
		results = append(results, avg.Round(2))

		// No results should be received after the context is cancelled
		if len(results) == 3 {
			cancel()
		}
	}

	expected := []decimal.Decimal{
		decimal.RequireFromString("10").Round(2),
		decimal.RequireFromString("15").Round(2),
		decimal.RequireFromString("20").Round(2),
	}

	assert.Equal(t, expected, results)
}
