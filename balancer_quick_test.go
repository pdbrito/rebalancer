package balancer_test

import (
	"fmt"
	. "github.com/pdbrito/balancer"
	"github.com/pdbrito/randomSum"
	"github.com/shopspring/decimal"
	"math/rand"
	"strconv"
	"testing"
	"testing/quick"
	"time"
)

func TestBalancer_ResultValueEqualToInput(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	f := func(n int) bool {
		holdingsBefore := generateHoldingsNumbering(rand.Intn(8) + 2)
		index := generateIndexForHoldings(holdingsBefore)

		valueBefore := value(holdingsBefore)

		trades := Balance(holdingsBefore, index)

		holdingsAfter := execute(trades, holdingsBefore)

		valueAfter := value(holdingsAfter)

		sub := valueAfter.Sub(valueBefore)
		fmt.Printf("difference: %v \n", sub)

		return valueAfter.Equal(valueBefore)
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func generateHoldingsNumbering(n int) map[Asset]Holding {
	holdings := make(map[Asset]Holding)
	for i := 0; i < n; i++ {
		assetKey := strconv.Itoa(i)
		holdings[Asset(assetKey)] = Holding{
			Quantity: decimal.New(int64(rand.Intn(9)+1), 0),
			Value:    decimal.New(int64(rand.Intn(9)+1), 0),
		}
	}
	return holdings
}

func generateIndexForHoldings(holdings map[Asset]Holding) map[Asset]decimal.Decimal {
	index := make(map[Asset]decimal.Decimal)

	numberOfAssets := len(holdings)
	indexValues := randomSum.NIntsTotaling(numberOfAssets, 100)

	i := 0
	for asset := range holdings {
		index[asset] = decimal.New(int64(indexValues[i]), -2)
		i++
	}

	return index
}

func value(holdings map[Asset]Holding) (sum decimal.Decimal) {
	for _, holding := range holdings {
		sum = sum.Add(holding.Quantity.Mul(holding.Value))
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

		quantityAfterTrade := holdings[asset].Quantity.Add(trade.Amount.Mul(modifier))

		res[asset] = Holding{
			Quantity: quantityAfterTrade,
			Value:    holdings[asset].Value,
		}
	}
	return res
}
