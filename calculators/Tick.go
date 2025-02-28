package calculators

import (
	"time"

	"github.com/shopspring/decimal"
)

type Tick struct {
	Price     decimal.Decimal
	Timestamp time.Time
}
