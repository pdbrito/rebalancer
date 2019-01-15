package balancer_test

import (
	"fmt"
	. "github.com/pdbrito/balancer"
	"github.com/shopspring/decimal"
	"testing"
)

func TestBalancer_Balance(t *testing.T) {
	holdings := map[Asset]Holding{
		"ETH": {
			Amount: decimal.NewFromFloat(20),
			Value:  decimal.NewFromFloat(200),
		},
		"BTC": {
			Amount: decimal.NewFromFloat(0.5),
			Value:  decimal.NewFromFloat(5000)},
	}

	index := map[Asset]decimal.Decimal{
		"ETH": decimal.NewFromFloat(0.3),
		"BTC": decimal.NewFromFloat(0.7),
	}

	got := Balance(holdings, index)
	want := map[Asset]Trade{
		"ETH": {Action: "sell", Amount: decimal.NewFromFloat(10.25)},
		"BTC": {Action: "buy", Amount: decimal.NewFromFloat(0.41)},
	}

	assertSameTrades(t, got, want)
}

func assertSameTrades(t *testing.T, got map[Asset]Trade, want map[Asset]Trade) {
	t.Helper()

	if len(got) != len(want) {
		t.Errorf("got %d trades want %d", len(got), len(want))
	}

	for asset, wantTrade := range want {
		gotTrade, exists := got[asset]
		if !exists {
			t.Errorf("asset %s missing from trade list", asset)
			return
		}
		if gotTrade.Action != wantTrade.Action {
			t.Errorf(
				"got a trade action of %s, want %s for asset %s",
				gotTrade.Action,
				wantTrade.Action,
				asset,
			)
		}
		if !gotTrade.Amount.Equal(wantTrade.Amount) {
			t.Errorf(
				"got %v want %v for trade of asset %s",
				gotTrade.Amount,
				wantTrade.Amount,
				asset,
			)
		}
	}
}

func ExampleBalance() {
	holdings := map[Asset]Holding{
		"ETH": {
			Amount: decimal.NewFromFloat(20),
			Value:  decimal.NewFromFloat(350),
		},
		"BTC": {
			Amount: decimal.NewFromFloat(0.5),
			Value:  decimal.NewFromFloat(5000),
		},
	}

	desiredWeighting := map[Asset]decimal.Decimal{
		"ETH": decimal.NewFromFloat(0.5),
		"BTC": decimal.NewFromFloat(0.5),
	}

	requiredTrades := Balance(holdings, desiredWeighting)

	for asset, trade := range requiredTrades {
		fmt.Printf("%s %s %s\n", trade.Action, trade.Amount, asset)
	}

	// Unordered output:
	// sell 6.4285714285714286 ETH
	// buy 0.45 BTC
}

func BenchmarkBalance(b *testing.B) {
	for i := 0; i < b.N; i++ {
		holdings := map[Asset]Holding{
			"ETH": {
				Amount: decimal.NewFromFloat(20),
				Value:  decimal.NewFromFloat(200)},
			"BTC": {
				Amount: decimal.NewFromFloat(0.5),
				Value:  decimal.NewFromFloat(5000),
			},
		}
		index := map[Asset]decimal.Decimal{
			"ETH": decimal.NewFromFloat(0.3),
			"BTC": decimal.NewFromFloat(0.7),
		}

		_ = Balance(holdings, index)
	}
}
