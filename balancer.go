// Package balancer provides functionality to balance investment assets to a
// target index. This is accomplished by calculating the current percentage
// allocation of assets and then the trades necessary to match the specified
// index.
package balancer

import (
	"github.com/shopspring/decimal"
)

// An Account has holdings, a pricelist and a calculated value
type Account struct {
	holdings  map[Asset]Holding
	pricelist map[Asset]decimal.Decimal
	value     decimal.Decimal
}

// NewAccount returns a new Account struct
func NewAccount(holdings map[Asset]Holding, pricelist map[Asset]decimal.Decimal) Account {
	totalValue := decimal.Zero
	for asset, holding := range holdings {
		totalValue = totalValue.Add(pricelist[asset].Mul(holding.Amount))
	}
	return Account{holdings: holdings, pricelist: pricelist, value: totalValue}
}

// An Asset is a string type used to identify your assets.
type Asset string

// A Holding represents an Amount.
type Holding struct {
	Amount decimal.Decimal
}

// A Trade represents a buy or sell action of a certain amount
type Trade struct {
	Action string
	Amount decimal.Decimal
}

// Balance will return a map[Asset]Trade which will balance the passed in
// holdings to match the passed in target index.
func (a Account) Balance(index map[Asset]decimal.Decimal) map[Asset]Trade {
	//validate assumptions; only unique assets etc
	trades := map[Asset]Trade{}

	amountRequired := decimal.Zero
	for asset, percentage := range index {

		amountRequired = a.value.Mul(percentage).Div(a.pricelist[asset])

		if holding, ok := a.holdings[asset]; ok {
			amountRequired = amountRequired.Sub(holding.Amount)
		}

		if amountRequired.IsNegative() {
			trades[asset] = Trade{"sell", amountRequired.Abs()}
			continue
		}
		trades[asset] = Trade{"buy", amountRequired.Abs()}
	}

	return trades
}
