package calculators

import (
	"context"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestWindowedAverage(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ticks := make(chan Tick)
	clk := clock.NewMock()
	go func() {
		defer close(ticks)
		ticks <- Tick{Price: decimal.NewFromFloat(10), Timestamp: clk.Now()}
		clk.Add(20 * time.Second)
		ticks <- Tick{Price: decimal.NewFromFloat(20), Timestamp: clk.Now()}
		clk.Add(20 * time.Second)
		ticks <- Tick{Price: decimal.NewFromFloat(30), Timestamp: clk.Now()}
		clk.Add(30 * time.Second)
		ticks <- Tick{Price: decimal.NewFromFloat(40), Timestamp: clk.Now()}
		clk.Add(40 * time.Second)
		ticks <- Tick{Price: decimal.NewFromFloat(50), Timestamp: clk.Now()}
	}()

	windowedAvgCh := WindowedAverage(ctx, ticks, time.Minute, clk)
	var results []decimal.Decimal
	for avg := range windowedAvgCh {
		results = append(results, avg.Round(2))
	}

	expected := []decimal.Decimal{
		decimal.NewFromFloat(10).Round(2),
		decimal.NewFromFloat(15).Round(2),
		decimal.NewFromFloat(25).Round(2),
		decimal.NewFromFloat(40).Round(2),
		decimal.NewFromFloat(45).Round(2),
	}

	assert.Equal(t, expected, results)
}

// This tests if the function correctly handles the context cancellation while processing the ticks.
// It ensures that the function stops processing further ticks once the context is cancelled.
func TestWindowedAverageWithContextCancellationWithinTicksChannel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	ticks := make(chan Tick)
	clk := clock.NewMock()

	windowedAvgCh := WindowedAverage(ctx, ticks, time.Minute, clk)

	go func() {
		defer close(ticks)
		ticks <- Tick{Price: decimal.NewFromFloat(10), Timestamp: clk.Now()}
		clk.Add(20 * time.Second)
		ticks <- Tick{Price: decimal.NewFromFloat(20), Timestamp: clk.Now()}
		clk.Add(20 * time.Second)
		ticks <- Tick{Price: decimal.NewFromFloat(30), Timestamp: clk.Now()}
		clk.Add(30 * time.Second)
		cancel()
		// Further ticks should not result in any more output
		ticks <- Tick{Price: decimal.NewFromFloat(40), Timestamp: clk.Now()}
		clk.Add(40 * time.Second)
		ticks <- Tick{Price: decimal.NewFromFloat(50), Timestamp: clk.Now()}
	}()

	var results []decimal.Decimal
	for avg := range windowedAvgCh {
		results = append(results, avg.Round(2))
	}

	expected := []decimal.Decimal{
		decimal.NewFromFloat(10).Round(2),
		decimal.NewFromFloat(15).Round(2),
		decimal.NewFromFloat(25).Round(2),
	}

	assert.Equal(t, expected, results)
}

// This tests if the function correctly handles the context cancellation from the consumer side.
// It ensures that the function stops sending values to the channel once the context is cancelled.
func TestWindowedAverageWithContextCancellationWithinOutChannel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	ticks := make(chan Tick)
	clk := clock.NewMock()
	go func() {
		defer close(ticks)
		ticks <- Tick{Price: decimal.NewFromFloat(10), Timestamp: clk.Now()}
		clk.Add(20 * time.Second)
		ticks <- Tick{Price: decimal.NewFromFloat(20), Timestamp: clk.Now()}
		clk.Add(20 * time.Second)
		ticks <- Tick{Price: decimal.NewFromFloat(30), Timestamp: clk.Now()}
		clk.Add(30 * time.Second)

		ticks <- Tick{Price: decimal.NewFromFloat(40), Timestamp: clk.Now()}
		clk.Add(40 * time.Second)
		ticks <- Tick{Price: decimal.NewFromFloat(50), Timestamp: clk.Now()}
	}()

	windowedAvgCh := WindowedAverage(ctx, ticks, time.Minute, clk)
	var results []decimal.Decimal
	for avg := range windowedAvgCh {
		results = append(results, avg.Round(2))

		// No results should be received after the context is cancelled
		if len(results) == 3 {
			cancel()
		}

	}

	expected := []decimal.Decimal{
		decimal.NewFromFloat(10).Round(2),
		decimal.NewFromFloat(15).Round(2),
		decimal.NewFromFloat(25).Round(2),
	}

	assert.Equal(t, expected, results)
}
