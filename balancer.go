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

// Holdings is an account level store of assets
type Holdings map[Asset]decimal.Decimal

// Pricelist contains a map of Assets and their current price
type Pricelist map[Asset]decimal.Decimal

// An Account has holdings, a pricelist and a calculated value
type Account struct {
	holdings  Holdings
	pricelist map[Asset]decimal.Decimal
	value     decimal.Decimal
}

// ErrEmptyHoldings indicates an empty holdings was passed to NewHoldings
var ErrEmptyHoldings = errors.New("holdings must not be empty")

// ErrEmptyPricelist indicated an empty pricelist was passed to NewPricelist
var ErrEmptyPricelist = errors.New("holdings must not be empty")

// ErrInvalidAsset indicates an Asset is not uppercase: "eth" vs "ETH"
var ErrInvalidAsset = errors.New("holding assets must be uppercase")

// ErrInvalidAssetAmount indicates an invalid asset amount of 0 or below
type ErrInvalidAssetAmount struct {
	Asset  Asset
	Amount decimal.Decimal
}

// Error formats the error message for ErrInvalidAssetAmount
func (e ErrInvalidAssetAmount) Error() string {
	return fmt.Sprintf("%s needs positive amount, not %s", e.Asset, e.Amount)
}

// NewHoldings validates and returns a new Holdings type
func NewHoldings(holdings map[Asset]decimal.Decimal) (Holdings, error) {
	if len(holdings) == 0 {
		return nil, ErrEmptyHoldings
	}
	for asset, holding := range holdings {
		if holding.LessThan(decimal.Zero) || holding.Equal(decimal.Zero) {
			return nil, ErrInvalidAssetAmount{Asset: asset, Amount: holding}
		}
		if string(asset) != strings.ToUpper(string(asset)) {
			return nil, ErrInvalidAsset
		}
	}
	return holdings, nil
}

// NewPricelist validates and returns a new Pricelist type
func NewPricelist(pricelist map[Asset]decimal.Decimal) (Pricelist, error) {
	if len(pricelist) == 0 {
		return nil, ErrEmptyPricelist
	}
	for asset, price := range pricelist {
		if price.LessThan(decimal.Zero) || price.Equal(decimal.Zero) {
			return nil, ErrInvalidAssetAmount{Asset: asset, Amount: price}
		}
		if string(asset) != strings.ToUpper(string(asset)) {
			return nil, ErrInvalidAsset
		}
	}
	return pricelist, nil
}

// NewAccount returns a new Account struct
func NewAccount(holdings map[Asset]decimal.Decimal, pricelist map[Asset]decimal.Decimal) (Account, error) {
	holdings, err := NewHoldings(holdings)
	if err != nil {
		return Account{}, err
	}
	pricelist, err = NewPricelist(pricelist)
	if err != nil {
		return Account{}, err
	}
	totalValue := decimal.Zero
	for asset, holding := range holdings {
		totalValue = totalValue.Add(pricelist[asset].Mul(holding))
	}
	return Account{holdings: holdings, pricelist: pricelist, value: totalValue}, nil
}

// An Asset is a string type used to identify your assets.
type Asset string

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
			amountRequired = amountRequired.Sub(holding)
		}

		if amountRequired.IsNegative() {
			trades[asset] = Trade{"sell", amountRequired.Abs()}
			continue
		}
		trades[asset] = Trade{"buy", amountRequired.Abs()}
	}

	return trades, nil
}
