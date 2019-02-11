// Package rebalancer provides functionality to rebalance a portfolio of
// investment assets to match a target index. This is accomplished by
// calculating the current percentage allocation of assets and determining the
// trades necessary to reallocate funds to match the desired target index.
package rebalancer

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

// ErrAssetMissingFromPricelist indicates an asset without a matching entry in
// the global pricelist.
var ErrAssetMissingFromPricelist = errors.New("asset missing from global pricelist")

// Portfolio contains a map of Assets and their current amount.
type Portfolio map[Asset]decimal.Decimal

// ErrEmptyPortfolio indicates an empty portfolio was passed to NewPortfolio.
var ErrEmptyPortfolio = errors.New("portfolio must not be empty")

// NewPortfolio validates and returns a new Portfolio type.
func NewPortfolio(portfolio map[Asset]decimal.Decimal) (Portfolio, error) {
	if len(portfolio) == 0 {
		return nil, ErrEmptyPortfolio
	}
	for asset, amount := range portfolio {
		if string(asset) != strings.ToUpper(string(asset)) {
			return nil, ErrInvalidAsset
		}
		if _, ok := globalPricelist[asset]; !ok {
			return nil, ErrAssetMissingFromPricelist
		}
		if amount.LessThan(decimal.Zero) || amount.Equal(decimal.Zero) {
			return nil, ErrInvalidAssetAmount{Asset: asset, Amount: amount}
		}
	}
	return portfolio, nil
}

// An Account has portfolio, a pricelist and a calculated total value.
type Account struct {
	portfolio Portfolio
	value     decimal.Decimal
}

// NewAccount validates portfolio and then returns a new Account struct.
func NewAccount(portfolio map[Asset]decimal.Decimal) (Account, error) {
	if len(globalPricelist) == 0 {
		return Account{}, ErrEmptyPricelist
	}
	portfolio, err := NewPortfolio(portfolio)
	if err != nil {
		return Account{}, err
	}
	totalValue := decimal.Zero
	for asset, amount := range portfolio {
		totalValue = totalValue.Add(globalPricelist[asset].Mul(amount))
	}
	return Account{portfolio: portfolio, value: totalValue}, nil
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
		if string(asset) != strings.ToUpper(string(asset)) {
			return nil, ErrInvalidAsset
		}
		if _, ok := globalPricelist[asset]; !ok {
			return nil, ErrAssetMissingFromPricelist
		}
		if percentage.LessThan(decimal.Zero) || percentage.Equal(decimal.Zero) {
			return nil, ErrInvalidAssetAmount{Asset: asset, Amount: percentage}
		}
		indexTotal = indexTotal.Add(percentage)
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

// Rebalance will return a map[Asset]Trade which will balance the account's
// portfolio to match the supplied target index.
func (a Account) Rebalance(targetIndex map[Asset]decimal.Decimal) (map[Asset]Trade, error) {
	targetIndex, err := NewIndex(targetIndex)
	if err != nil {
		return nil, err
	}

	trades := map[Asset]Trade{}
	amountRequired := decimal.Zero

	for asset, percentage := range targetIndex {
		amountRequired = a.value.Mul(percentage).Div(globalPricelist[asset])

		if portfolioAmount, ok := a.portfolio[asset]; ok {
			amountRequired = amountRequired.Sub(portfolioAmount)
		}

		if amountRequired.IsNegative() {
			trades[asset] = Trade{"sell", amountRequired.Abs()}
			continue
		}
		trades[asset] = Trade{"buy", amountRequired.Abs()}
	}

	return trades, nil
}
