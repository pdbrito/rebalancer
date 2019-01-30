// Package balancer provides functionality to balance investment assets to a
// target index. This is accomplished by calculating the current percentage
// allocation of assets and then the trades necessary to match the specified
// target index.
package balancer

import (
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"strings"
)

// An Account has holdings, a pricelist and a calculated value
type Account struct {
	holdings  map[Asset]Holding
	pricelist map[Asset]decimal.Decimal
	value     decimal.Decimal
}

// Holdings are a map[Asset]Holding
type Holdings map[Asset]Holding

// ErrEmptyHoldings indicated an empty holdings was passed to NewHoldings
var ErrEmptyHoldings = errors.New("holdings must not be empty")

// ErrInvalidAsset indicates an Asset is not uppercase: "eth" vs "ETH"
var ErrInvalidAsset = errors.New("holding assets must be uppercase")

// ErrInvalidHoldingAmount indicates an invalid Holding.Amount of 0 or below
type ErrInvalidHoldingAmount struct {
	Asset  Asset
	Amount decimal.Decimal
}

// Error formats the error message for ErrInvalidHoldingAmount
func (e ErrInvalidHoldingAmount) Error() string {
	return fmt.Sprintf("%s needs positive amount, not %s", e.Asset, e.Amount)
}

// NewHoldings validates and returns a new Holdings struct
func NewHoldings(holdings map[Asset]Holding) (Holdings, error) {
	if len(holdings) == 0 {
		return nil, ErrEmptyHoldings
	}
	for asset, holding := range holdings {
		if holding.Amount.LessThan(decimal.Zero) || holding.Amount.Equal(decimal.Zero) {
			return nil, ErrInvalidHoldingAmount{Asset: asset, Amount: holding.Amount}
		}
		if string(asset) != strings.ToUpper(string(asset)) {
			return nil, ErrInvalidAsset
		}
	}
	return holdings, nil
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
func (a Account) Balance(targetIndex map[Asset]decimal.Decimal) (map[Asset]Trade, error) {
	indexTotal := decimal.Zero
	for asset, percentage := range targetIndex {
		indexTotal = indexTotal.Add(percentage)
		if _, ok := a.pricelist[asset]; !ok {
			return nil, fmt.Errorf(
				"targetIndex contains asset missing from the pricelist: %s",
				asset,
			)
		}
	}
	if !indexTotal.Equal(decimal.NewFromFloat(1)) {
		return nil, fmt.Errorf(
			"targetIndex should sum to 1, got %v from %v",
			indexTotal,
			targetIndex,
		)
	}

	trades := map[Asset]Trade{}

	amountRequired := decimal.Zero
	for asset, percentage := range targetIndex {

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

	return trades, nil
}
