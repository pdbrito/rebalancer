package balancer_test

import (
	"fmt"
	. "github.com/pdbrito/balancer"
	"github.com/pdbrito/randomSum"
	"github.com/shopspring/decimal"
	"math/rand"
	"reflect"
	"strconv"
	"testing"
	"testing/quick"
)

type fakeAccount struct {
	Holdings    map[Asset]Holding
	TargetIndex map[Asset]decimal.Decimal
	Pricelist   map[Asset]decimal.Decimal
}

func (f fakeAccount) Generate(rand *rand.Rand, size int) reflect.Value {
	holdings := generateHoldingsNumbering(rand.Intn(99) + 1)

	return reflect.ValueOf(fakeAccount{
		Holdings:    holdings,
		TargetIndex: generateTargetIndexForHoldings(holdings),
		Pricelist:   generatePricelistForHoldings(holdings),
	})
}

func TestBalancer_ResultingIndexEqualToTargetIndex(t *testing.T) {
	assertion := func(f fakeAccount) bool {
		Account := NewAccount(f.Holdings, f.Pricelist)
		trades, err := Account.Balance(f.TargetIndex)

		if err != nil {
			return false
		}

		holdingsAfter := execute(trades, f.Holdings)

		resultingIndex := calculateIndex(holdingsAfter, f.Pricelist)

		return indexesAreEqual(f.TargetIndex, resultingIndex)
	}

	if err := quick.Check(assertion, nil); err != nil {
		if e, ok := err.(*quick.CheckError); ok {
			for _, value := range e.In {
				f := value.(fakeAccount)
				Account := NewAccount(f.Holdings, f.Pricelist)
				trades, _ := Account.Balance(f.TargetIndex)
				fmt.Printf("Holdings: %v\n", f.Holdings)
				fmt.Printf("Target index: %v\n", f.TargetIndex)
				fmt.Printf("Trades: %v\n", trades)
			}
		}

		t.Error(err)
	}
}

func calculateIndex(holdings map[Asset]Holding, pricelist map[Asset]decimal.Decimal) map[Asset]decimal.Decimal {
	index := make(map[Asset]decimal.Decimal)
	value := value(holdings, pricelist)
	for asset, holding := range holdings {
		index[asset] = pricelist[asset].Mul(holding.Amount).Div(value)
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

func generateHoldingsNumbering(n int) map[Asset]Holding {
	holdings := make(map[Asset]Holding)
	for i := 0; i < n; i++ {
		assetKey := strconv.Itoa(i)
		holdings[Asset(assetKey)] = Holding{
			Amount: decimal.NewFromFloat(rand.Float64() * 1000),
		}
	}
	return holdings
}

func generateTargetIndexForHoldings(holdings map[Asset]Holding) map[Asset]decimal.Decimal {
	targetIndex := make(map[Asset]decimal.Decimal)

	numberOfAssets := len(holdings)
	indexValues := randomSum.NIntsTotaling(numberOfAssets, 100)

	i := 0
	for asset := range holdings {
		targetIndex[asset] = decimal.New(int64(indexValues[i]), -2)
		i++
	}

	return targetIndex
}

func value(holdings map[Asset]Holding, pricelist map[Asset]decimal.Decimal) (sum decimal.Decimal) {
	for asset, holding := range holdings {
		sum = sum.Add(holding.Amount.Mul(pricelist[asset]))
	}
	return sum
}

func execute(trades map[Asset]Trade, holdings map[Asset]Holding) map[Asset]Holding {
	res := map[Asset]Holding{}

	for asset, trade := range trades {

		modifier := decimal.New(1, 0)
		if trade.Action == "sell" {
			modifier = decimal.New(-1, 0)
		}

		quantityAfterTrade := holdings[asset].Amount.Add(trade.Amount.Mul(modifier))

		res[asset] = Holding{
			Amount: quantityAfterTrade,
		}
	}
	return res
}

func generatePricelistForHoldings(holdings map[Asset]Holding) map[Asset]decimal.Decimal {
	pricelist := map[Asset]decimal.Decimal{}
	for asset := range holdings {
		pricelist[asset] = decimal.NewFromFloat(rand.Float64() * 1000)
	}
	return pricelist
}
