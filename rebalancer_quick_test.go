package rebalancer_test

import (
	"fmt"
	"github.com/pdbrito/randomSum"
	. "github.com/pdbrito/rebalancer"
	"github.com/shopspring/decimal"
	"math/rand"
	"reflect"
	"strconv"
	"testing"
	"testing/quick"
)

type fakeAccount struct {
	Portfolio   map[Asset]decimal.Decimal
	TargetIndex map[Asset]decimal.Decimal
	Pricelist   map[Asset]decimal.Decimal
}

func (f fakeAccount) Generate(rand *rand.Rand, size int) reflect.Value {
	portfolio := generatePortfolio(rand.Intn(99) + 1)

	return reflect.ValueOf(fakeAccount{
		Portfolio:   portfolio,
		TargetIndex: generateTargetIndexForPortfolio(portfolio),
		Pricelist:   generatePricelistForPortfolio(portfolio),
	})
}

func TestRebalance_ResultingIndexEqualToTargetIndex(t *testing.T) {
	assertion := func(f fakeAccount) bool {
		_ = SetPricelist(f.Pricelist)
		account, _ := NewAccount(f.Portfolio)
		trades, err := account.Rebalance(f.TargetIndex)

		if err != nil {
			return false
		}

		portfolioAfter := execute(trades, f.Portfolio)

		resultingIndex := calculateIndex(portfolioAfter, f.Pricelist)

		return indexesAreEqual(f.TargetIndex, resultingIndex)
	}

	if err := quick.Check(assertion, nil); err != nil {
		if e, ok := err.(*quick.CheckError); ok {
			for _, value := range e.In {
				f := value.(fakeAccount)
				account, _ := NewAccount(f.Portfolio)
				trades, _ := account.Rebalance(f.TargetIndex)
				fmt.Printf("Portfolio: %v\n", f.Portfolio)
				fmt.Printf("Target index: %v\n", f.TargetIndex)
				fmt.Printf("Trades: %v\n", trades)
			}
		}

		t.Error(err)
	}
}

func calculateIndex(portfolio map[Asset]decimal.Decimal, pricelist map[Asset]decimal.Decimal) map[Asset]decimal.Decimal {
	index := make(map[Asset]decimal.Decimal)
	value := value(portfolio, pricelist)
	for asset, amount := range portfolio {
		index[asset] = pricelist[asset].Mul(amount).Div(value)
	}
	return index
}

func indexesAreEqual(i1, i2 map[Asset]decimal.Decimal) bool {
	if len(i1) != len(i2) {
		return false
	}
	for asset, amount := range i1 {
		if amount2, ok := i2[asset]; !ok || !amount.Equal(amount2) {
			return false
		}
	}
	return true
}

func generatePortfolio(n int) map[Asset]decimal.Decimal {
	portfolio := make(map[Asset]decimal.Decimal)
	for i := 0; i < n; i++ {
		assetKey := strconv.Itoa(i)
		portfolio[Asset(assetKey)] = decimal.NewFromFloat(rand.Float64() * 1000)
	}
	return portfolio
}

func generateTargetIndexForPortfolio(portfolio map[Asset]decimal.Decimal) map[Asset]decimal.Decimal {
	targetIndex := make(map[Asset]decimal.Decimal)

	numberOfAssets := len(portfolio)
	indexValues := randomSum.NIntsTotaling(numberOfAssets, 100)

	i := 0
	for asset := range portfolio {
		targetIndex[asset] = decimal.New(int64(indexValues[i]), -2)
		i++
	}

	return targetIndex
}

func value(portfolio map[Asset]decimal.Decimal, pricelist map[Asset]decimal.Decimal) (sum decimal.Decimal) {
	for asset, amount := range portfolio {
		sum = sum.Add(amount.Mul(pricelist[asset]))
	}
	return sum
}

func execute(trades map[Asset]Trade, portfolio map[Asset]decimal.Decimal) map[Asset]decimal.Decimal {
	res := map[Asset]decimal.Decimal{}

	for asset, trade := range trades {

		modifier := decimal.New(1, 0)
		if trade.Action == "sell" {
			modifier = decimal.New(-1, 0)
		}

		quantityAfterTrade := portfolio[asset].Add(trade.Amount.Mul(modifier))

		res[asset] = quantityAfterTrade

	}
	return res
}

func generatePricelistForPortfolio(portfolio map[Asset]decimal.Decimal) map[Asset]decimal.Decimal {
	pricelist := map[Asset]decimal.Decimal{}
	for asset := range portfolio {
		pricelist[asset] = decimal.NewFromFloat(rand.Float64() * 1000)
	}
	return pricelist
}
