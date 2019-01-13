// Package balancer provides functionality to balance investment assets to a
// target index. This is accomplished by calculating the current percentage
// allocation of assets and then the trades necessary to match the specified
// index.
package balancer

import (
	"github.com/shopspring/decimal"
)

// An Asset is a string type used to identify your assets.
type Asset string

// A Holding represents a current quantity and value.
type Holding struct {
	Quantity decimal.Decimal
	Value    decimal.Decimal
}

// A Trade represents a buy or sell action of a certain amount
type Trade struct {
	Action string
	Amount decimal.Decimal
}

// Balance will return a map[Asset]Trade which will balance the passed in
// holdings to match the passed in index.
func Balance(holdings map[Asset]Holding, index map[Asset]decimal.Decimal) map[Asset]Trade {
	//validate assumptions; only unique assets etc
	totalHoldings := decimal.Zero
	for _, holding := range holdings {
		totalHoldings = totalHoldings.Add(holding.Value.Mul(holding.Quantity))
	}

	trades := map[Asset]Trade{}

	for asset, weight := range index {
		amountRequired :=
			totalHoldings.
				Mul(weight).
				Div(holdings[asset].Value).
				Sub(holdings[asset].Quantity)
		trades[asset] = makeTrade(amountRequired)
	}

	return trades
}
func makeTrade(amount decimal.Decimal) Trade {
	var action string
	if amount.IsNegative() {
		action = "sell"
	} else {
		action = "buy"
	}
	return Trade{action, amount.Abs()}
}
