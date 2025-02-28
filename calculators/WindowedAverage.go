package calculators

// Package calculators provides various functions to calculate averages.
//
// WindowedAverage.go contains the implementation of a windowed average calculator.
// This calculator computes the average of values within a specified time window.
//
// Dependencies:
// - context: for managing the lifecycle of the average calculation process.
// - fmt: for formatted I/O operations.
// - time: for handling time-related operations.
// - github.com/benbjohnson/clock: for mocking time in tests.
// - github.com/shopspring/decimal: for high-precision decimal arithmetic.
import (
	"context"
	"fmt"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/shopspring/decimal"
)

// WindowedAverage calculates the windowed average of the prices from the ticks channel.
func WindowedAverage(ctx context.Context, ticks <-chan Tick, d time.Duration, clk clock.Clock) <-chan decimal.Decimal {
	out := make(chan decimal.Decimal)
	go func() {
		defer close(out) // Close the output channel when the function returns
		var window []Tick
		var sum decimal.Decimal

		for {
			select {
			case <-ctx.Done(): // If the context is cancelled, stop the processing and close output channel
				return
			case tick, ok := <-ticks: // If a new tick is received, process it, if the channel is closed, return
				if !ok {
					return
				}
				window = append(window, tick)
				sum = sum.Add(tick.Price)

				// Print the current clock value
				fmt.Println("Current clock value:", clk.Now())

				// Remove old ticks whichs is falling beyond duration d window
				for len(window) > 0 && clk.Since(window[0].Timestamp) > d {
					sum = sum.Sub(window[0].Price)
					window = window[1:]
				}

				if len(window) > 0 {
					out <- sum.Div(decimal.NewFromInt(int64(len(window))))
				}
			}
		}
	}()
	return out
}
