package balancer

import (
	"github.com/shopspring/decimal"
)

type Asset string

type Holding struct {
	Quantity decimal.Decimal
	Value    decimal.Decimal
}

type Trade struct {
	Action string
	Amount decimal.Decimal
}

func Balance(holdings map[Asset]Holding, index map[Asset]decimal.Decimal) map[Asset]Trade {
	//validate assumptions; only unique assets etc
	totalHoldings := decimal.NewFromFloat(0)
	for _, holding := range holdings {
		totalHoldings = totalHoldings.Add(holding.Value.Mul(holding.Quantity))
	}

	trades := map[Asset]Trade{}

	for asset, weight := range index {
		amountRequired := totalHoldings.Mul(weight).Div(holdings[asset].Value).Sub(holdings[asset].Quantity)
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
