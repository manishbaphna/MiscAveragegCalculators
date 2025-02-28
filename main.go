package main

import (
	"MiscAveragegCalculators/calculators"
	"context"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

// Tick represents a price at a specific timestamp.
type Tick struct {
	Timestamp time.Time
	Price     decimal.Decimal
}

func main() {
	// Create a context with cancellation.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a channel to send Tick values.
	ticks := make(chan calculators.Tick)

	// Goroutine to send Tick values to the channel.
	go func() {
		for i := 0; i < 10; i++ {
			ticks <- calculators.Tick{Price: decimal.NewFromFloat(float64(i) + 1.0)}
			time.Sleep(100 * time.Millisecond)
		}
		close(ticks)
	}()

	// Define the alpha value for the Exponential Moving Average.
	alpha := 0.1

	// Call the ExponentialMovingAverage function.
	emaChannel := calculators.ExponentialMovingAverage(ctx, ticks, 10, alpha)

	// Print the Exponential Moving Averages received from the channel.
	for ema := range emaChannel {
		fmt.Println("Exponential Moving Average:", ema)
	}
}
