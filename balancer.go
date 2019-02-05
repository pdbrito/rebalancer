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

// An Asset is a string type used to identify your assets. It must be uppercase.
type Asset string

// ErrInvalidAsset indicates an Asset is not uppercase: "eth" vs "ETH".
var ErrInvalidAsset = errors.New("assets must be uppercase")

// ErrInvalidAssetAmount indicates an invalid asset amount of 0 or below.
type ErrInvalidAssetAmount struct {
	Asset  Asset
	Amount decimal.Decimal
}

// Error formats the error message for ErrInvalidAssetAmount.
func (e ErrInvalidAssetAmount) Error() string {
	return fmt.Sprintf("%s needs positive amount, not %s", e.Asset, e.Amount)
}

// globalPricelist contains the current pricelist used for all calculations.
var globalPricelist = Pricelist{}

// Pricelist contains a map of Assets and their current price.
type Pricelist map[Asset]decimal.Decimal

// ErrEmptyPricelist indicates an empty pricelist was passed to NewPricelist.
var ErrEmptyPricelist = errors.New("pricelist must not be empty")

// SetPricelist validates and sets a new Pricelist.
func SetPricelist(pricelist map[Asset]decimal.Decimal) error {
	if len(pricelist) == 0 {
		return ErrEmptyPricelist
	}
	for asset, price := range pricelist {
		if price.LessThan(decimal.Zero) || price.Equal(decimal.Zero) {
			return ErrInvalidAssetAmount{Asset: asset, Amount: price}
		}
		if string(asset) != strings.ToUpper(string(asset)) {
			return ErrInvalidAsset
		}
	}
	globalPricelist = pricelist
	return nil
}

// GlobalPricelist returns the current value of the global pricelist.
func GlobalPricelist() Pricelist {
	return globalPricelist
}

// ClearGlobalPricelist clears the global pricelist.
func ClearGlobalPricelist() {
	globalPricelist = Pricelist{}
}

// Holdings contains a map of Assets and their current quantity.
type Holdings map[Asset]decimal.Decimal

// ErrEmptyHoldings indicates an empty holdings was passed to NewHoldings.
var ErrEmptyHoldings = errors.New("holdings must not be empty")

// NewHoldings validates and returns a new Holdings type.
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
		//if asset it not in pricelist, error
	}
	return holdings, nil
}

// An Account has holdings, a pricelist and a calculated total value.
type Account struct {
	holdings Holdings
	value    decimal.Decimal
}

// NewAccount validates holdings and then returns a new Account struct.
func NewAccount(holdings map[Asset]decimal.Decimal) (Account, error) {
	if len(globalPricelist) == 0 {
		return Account{}, ErrEmptyPricelist
	}
	holdings, err := NewHoldings(holdings)
	if err != nil {
		return Account{}, err
	}
	totalValue := decimal.Zero
	for asset, holding := range holdings {
		totalValue = totalValue.Add(globalPricelist[asset].Mul(holding))
	}
	return Account{holdings: holdings, value: totalValue}, nil
}

// Index contains a map of Assets and their values. Indexes values must
// always sum to 1.
type Index map[Asset]decimal.Decimal

// ErrEmptyIndex indicates an empty index was passed to NewIndex.
var ErrEmptyIndex = errors.New("index must not be empty")

// ErrIndexSumIncorrect indicates that the sum of the values in an index is not
// equal to 1.
var ErrIndexSumIncorrect = errors.New("index values must sum to 1")

// NewIndex validates and returns a new Index type whose values must sum to 1.
func NewIndex(index map[Asset]decimal.Decimal) (Index, error) {
	if len(index) == 0 {
		return nil, ErrEmptyIndex
	}
	indexTotal := decimal.Zero
	for asset, percentage := range index {
		indexTotal = indexTotal.Add(percentage)
		if percentage.LessThan(decimal.Zero) || percentage.Equal(decimal.Zero) {
			return nil, ErrInvalidAssetAmount{Asset: asset, Amount: percentage}
		}
		if string(asset) != strings.ToUpper(string(asset)) {
			return nil, ErrInvalidAsset
		}
		//if asset is not in pricelist, error
	}
	if !indexTotal.Equal(decimal.NewFromFloat(1)) {
		return nil, ErrIndexSumIncorrect
	}
	return index, nil
}

// A Trade represents a buy or sell action of a certain amount.
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
		if _, ok := globalPricelist[asset]; !ok {
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

		amountRequired = a.value.Mul(percentage).Div(globalPricelist[asset])

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
