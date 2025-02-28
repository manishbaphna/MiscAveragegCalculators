package calculators

// MovingAverage calculates the moving average of the prices from the ticks channel.

import (
	"context"

	"github.com/shopspring/decimal"
)

// ctx is the context that can be used to cancel the operation.
// ticks is the channel that provides the Tick values.
// N is the number of ticks to consider in the moving average.
// The function returns a channel that provides the moving average values.
func MovingAverage(ctx context.Context, ticks <-chan Tick, N int) <-chan decimal.Decimal {
	out := make(chan decimal.Decimal)
	go func() {
		defer close(out)
		var window []decimal.Decimal
		var sum decimal.Decimal

		for {
			select {
			case <-ctx.Done(): // If the context is cancelled, stop the processing and close output channel
				return
			case tick, ok := <-ticks: // If a new tick is received, process it, if the channel is closed, return
				if !ok {
					return
				}
				window = append(window, tick.Price)
				sum = sum.Add(tick.Price)
				if len(window) > N {
					sum = sum.Sub(window[0]) // Subtract the oldest value from the sum which is beyond last N ticks
					window = window[1:]
				}
				out <- sum.Div(decimal.NewFromInt(int64(len(window))))
			}
		}
	}()
	return out
}
