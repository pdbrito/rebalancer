package balancer_test

import (
	"fmt"
	. "github.com/pdbrito/balancer"
	"github.com/shopspring/decimal"
	"testing"
)

func TestAccount_Balance(t *testing.T) {
	holdings := map[Asset]Holding{
		"ETH": {
			Amount: decimal.NewFromFloat(20),
			Price:  decimal.NewFromFloat(200),
		},
		"BTC": {
			Amount: decimal.NewFromFloat(0.5),
			Price:  decimal.NewFromFloat(5000)},
	}

	Account := NewAccount(holdings)

	index := map[Asset]decimal.Decimal{
		"ETH": decimal.NewFromFloat(0.3),
		"BTC": decimal.NewFromFloat(0.7),
	}

	pricelist := map[Asset]decimal.Decimal{
		"ETH": decimal.NewFromFloat(200),
		"BTC": decimal.NewFromFloat(5000),
	}

	got := Account.Balance(index, pricelist)
	want := map[Asset]Trade{
		"ETH": {Action: "sell", Amount: decimal.NewFromFloat(10.25)},
		"BTC": {Action: "buy", Amount: decimal.NewFromFloat(0.41)},
	}

	assertSameTrades(t, got, want)
}

func TestAccount_Balance_IntoNewAssets(t *testing.T) {
	holdings := map[Asset]Holding{
		"ETH": {
			Amount: decimal.NewFromFloat(42),
			Price:  decimal.NewFromFloat(200),
		},
	}

	index := map[Asset]decimal.Decimal{
		"ETH":  decimal.NewFromFloat(0.2),
		"BTC":  decimal.NewFromFloat(0.2),
		"IOTA": decimal.NewFromFloat(0.2),
		"BAT":  decimal.NewFromFloat(0.2),
		"XLM":  decimal.NewFromFloat(0.2),
	}

	pricelist := map[Asset]decimal.Decimal{
		"ETH":  decimal.NewFromFloat(200),
		"BTC":  decimal.NewFromFloat(2000),
		"IOTA": decimal.NewFromFloat(0.3),
		"BAT":  decimal.NewFromFloat(0.12),
		"XLM":  decimal.NewFromFloat(0.2),
	}

	Account := NewAccount(holdings)

	got := Account.Balance(index, pricelist)
	want := map[Asset]Trade{
		"ETH":  {Action: "sell", Amount: decimal.NewFromFloat(33.6)},
		"BTC":  {Action: "buy", Amount: decimal.NewFromFloat(0.84)},
		"IOTA": {Action: "buy", Amount: decimal.NewFromFloat(5600)},
		"BAT":  {Action: "buy", Amount: decimal.NewFromFloat(14000)},
		"XLM":  {Action: "buy", Amount: decimal.NewFromFloat(8400)},
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

func ExampleAccount_Balance() {
	holdings := map[Asset]Holding{
		"ETH": {
			Amount: decimal.NewFromFloat(20),
			Price:  decimal.NewFromFloat(350),
		},
		"BTC": {
			Amount: decimal.NewFromFloat(0.5),
			Price:  decimal.NewFromFloat(5000),
		},
	}

	targetIndex := map[Asset]decimal.Decimal{
		"ETH": decimal.NewFromFloat(0.5),
		"BTC": decimal.NewFromFloat(0.5),
	}

	Account := NewAccount(holdings)

	pricelist := map[Asset]decimal.Decimal{
		"ETH": decimal.NewFromFloat(350),
		"BTC": decimal.NewFromFloat(5000),
	}

	requiredTrades := Account.Balance(targetIndex, pricelist)

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
				Price:  decimal.NewFromFloat(200)},
			"BTC": {
				Amount: decimal.NewFromFloat(0.5),
				Price:  decimal.NewFromFloat(5000),
			},
		}
		index := map[Asset]decimal.Decimal{
			"ETH": decimal.NewFromFloat(0.3),
			"BTC": decimal.NewFromFloat(0.7),
		}

		Account := NewAccount(holdings)

		pricelist := map[Asset]decimal.Decimal{
			"ETH": decimal.NewFromFloat(200),
			"BTC": decimal.NewFromFloat(5000),
		}

		_ = Account.Balance(index, pricelist)
	}
}
