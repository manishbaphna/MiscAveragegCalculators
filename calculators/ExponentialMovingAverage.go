package calculators

// ExponentialMovingAverage calculates the exponential moving average of the prices from the ticks channel.
// The function uses the formula:
// EMA = alpha * Price + (1 - alpha) * EMAt-1
// where alpha is the smoothing factor, Price is the current price, and EMAt-1 is the previous EMA value.
// The function returns a channel that provides the exponential moving average values.

import (
	"context"

	"github.com/shopspring/decimal"
)

// ctx is the context that can be used to cancel the operation.
// ticks is the channel that provides the Tick values.
// N is the number of ticks to consider in the moving average.
// alpha is the smoothing factor for the exponential moving average.
// The function returns a channel that provides the moving average values.
func ExponentialMovingAverage(ctx context.Context, ticks <-chan Tick, N int, alpha float64) <-chan decimal.Decimal {
	out := make(chan decimal.Decimal)
	alphaDec := decimal.NewFromFloat(alpha)
	go func() {
		defer close(out)
		var window []decimal.Decimal
		var ema decimal.Decimal
		// a constant factor to be used in the calculation of the impact of the oldest value in the window
		removalFactor := decimal.NewFromFloat(1).Sub(alphaDec).Pow(decimal.NewFromInt(int64(N - 1)))

		for {
			select {
			case <-ctx.Done(): // If the context is cancelled, stop the processing and close output channel
				return
			case tick, ok := <-ticks: // If a new tick is received, process it, if the channel is closed, return
				if !ok {
					return
				}
				window = append(window, tick.Price)

				if len(window) > N {
					ema = ema.Sub(alphaDec.Mul(window[0]).Mul(removalFactor)) // Subtract the oldest value and it's total impact from the ema which is beyond last N ticks
					window = window[1:]
					ema = alphaDec.Mul(tick.Price).Add(decimal.NewFromFloat(1).Sub(alphaDec).Mul(ema))
				} else if len(window) == 1 {
					ema = alphaDec.Mul(tick.Price)
				} else {
					ema = alphaDec.Mul(tick.Price).Add(decimal.NewFromFloat(1).Sub(alphaDec).Mul(ema))
				}
				out <- ema
			}
		}
	}()
	return out
}
