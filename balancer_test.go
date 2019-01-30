package balancer_test

import (
	"fmt"
	. "github.com/pdbrito/balancer"
	"github.com/shopspring/decimal"
	"reflect"
	"testing"
)

func TestAccount_Balance(t *testing.T) {
	holdings := map[Asset]Holding{
		"ETH": {
			Amount: decimal.NewFromFloat(20),
		},
		"BTC": {
			Amount: decimal.NewFromFloat(0.5),
		},
	}

	targetIndex := map[Asset]decimal.Decimal{
		"ETH": decimal.NewFromFloat(0.3),
		"BTC": decimal.NewFromFloat(0.7),
	}

	pricelist := map[Asset]decimal.Decimal{
		"ETH": decimal.NewFromFloat(200),
		"BTC": decimal.NewFromFloat(5000),
	}

	Account := NewAccount(holdings, pricelist)

	got, err := Account.Balance(targetIndex)

	if err != nil {
		t.Error("got an error but didn't want one")
	}

	want := map[Asset]Trade{
		"ETH": {Action: "sell", Amount: decimal.NewFromFloat(10.25)},
		"BTC": {Action: "buy", Amount: decimal.NewFromFloat(0.41)},
	}

	assertSameTrades(t, got, want)
}

func TestNewHoldings(t *testing.T) {
	got, err := NewHoldings(map[Asset]Holding{
		"ETH": {Amount: decimal.NewFromFloat(5)},
	})

	if err != nil {
		t.Error("got an error but didn't want one")
	}

	want := Holdings{"ETH": {Amount: decimal.NewFromFloat(5)}}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestErrInvalidAssetAmount_Error(t *testing.T) {
	asset := Asset("ETH")
	amount := decimal.NewFromFloat(-5)

	err := ErrInvalidHoldingAmount{Asset: asset, Amount: amount}

	want := "ETH needs positive amount, not -5"
	got := err.Error()

	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}

func TestNewHoldings_ErrorsOnNonPositiveHoldingAmount(t *testing.T) {
	asset := Asset("ETH")
	amount := decimal.NewFromFloat(-5)

	_, err := NewHoldings(map[Asset]Holding{
		asset: {Amount: amount},
	})

	want := ErrInvalidHoldingAmount{Asset: asset, Amount: amount}

	if err != want {
		t.Errorf("got %v, want %v", err, want)
	}
}

func TestNewHoldings_ErrorsOnInvalidInput(t *testing.T) {
	testCases := []struct {
		name     string
		holdings map[Asset]Holding
		err      error
	}{
		{
			name:     "holdings must not be empty",
			holdings: map[Asset]Holding{},
			err:      ErrEmptyHoldings,
		},
		{
			name: "holding assets should be uppercase and unique",
			holdings: map[Asset]Holding{
				"eth": {
					Amount: decimal.NewFromFloat(5),
				},
			},
			err: ErrInvalidAsset,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewHoldings(tt.holdings)

			if err == nil {
				t.Error("wanted an error but didn't get one")
			}

			if err != tt.err {
				t.Errorf("got %v want %v", err, tt.err)
			}
		})
	}
}

func TestAccount_Balance_IntoNewAssets(t *testing.T) {
	holdings := map[Asset]Holding{
		"ETH": {
			Amount: decimal.NewFromFloat(42),
		},
	}

	targetIndex := map[Asset]decimal.Decimal{
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

	Account := NewAccount(holdings, pricelist)

	got, err := Account.Balance(targetIndex)

	if err != nil {
		t.Error("got an error but didn't want one")
	}

	want := map[Asset]Trade{
		"ETH":  {Action: "sell", Amount: decimal.NewFromFloat(33.6)},
		"BTC":  {Action: "buy", Amount: decimal.NewFromFloat(0.84)},
		"IOTA": {Action: "buy", Amount: decimal.NewFromFloat(5600)},
		"BAT":  {Action: "buy", Amount: decimal.NewFromFloat(14000)},
		"XLM":  {Action: "buy", Amount: decimal.NewFromFloat(8400)},
	}

	assertSameTrades(t, got, want)
}

func TestAccount_Balance_ErrorsWhenTargetIndexIsInvalid(t *testing.T) {
	testCases := []struct {
		name        string
		targetIndex map[Asset]decimal.Decimal
	}{
		{
			name: "target index does not sum to 1",
			targetIndex: map[Asset]decimal.Decimal{
				"ETH": decimal.NewFromFloat(0.2),
				"BTC": decimal.NewFromFloat(0.2),
			},
		},
		{
			name:        "target index is empty",
			targetIndex: map[Asset]decimal.Decimal{},
		},
		{
			name: "target index has an asset missing from the pricelist",
			targetIndex: map[Asset]decimal.Decimal{
				"ETH": decimal.NewFromFloat(0.8),
				"BAT": decimal.NewFromFloat(0.2),
			},
		},
	}

	holdings := map[Asset]Holding{
		"ETH": {
			Amount: decimal.NewFromFloat(42),
		},
	}

	pricelist := map[Asset]decimal.Decimal{
		"ETH": decimal.NewFromFloat(200),
		"BTC": decimal.NewFromFloat(2000),
	}

	Account := NewAccount(holdings, pricelist)

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Account.Balance(tt.targetIndex)

			if err == nil {
				t.Error("wanted an error but didn't get one")
			}
		})
	}
}

func assertSameTrades(t *testing.T, got map[Asset]Trade, want map[Asset]Trade) {
	t.Helper()

	if len(got) != len(want) {
		t.Errorf("got %d trades want %d", len(got), len(want))
	}

	for asset, wantTrade := range want {
		gotTrade, exists := got[asset]
		if !exists {
			t.Fatalf("asset %s missing from trade list", asset)
		}
		if gotTrade.Action != wantTrade.Action {
			t.Fatalf(
				"got a trade action of %s, want %s for asset %s",
				gotTrade.Action,
				wantTrade.Action,
				asset,
			)
		}
		if !gotTrade.Amount.Equal(wantTrade.Amount) {
			t.Fatalf(
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
		},
		"BTC": {
			Amount: decimal.NewFromFloat(0.5),
		},
	}

	targetIndex := map[Asset]decimal.Decimal{
		"ETH": decimal.NewFromFloat(0.5),
		"BTC": decimal.NewFromFloat(0.5),
	}

	pricelist := map[Asset]decimal.Decimal{
		"ETH": decimal.NewFromFloat(350),
		"BTC": decimal.NewFromFloat(5000),
	}

	Account := NewAccount(holdings, pricelist)

	requiredTrades, _ := Account.Balance(targetIndex)

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
			},
			"BTC": {
				Amount: decimal.NewFromFloat(0.5),
			},
		}
		targetIndex := map[Asset]decimal.Decimal{
			"ETH": decimal.NewFromFloat(0.3),
			"BTC": decimal.NewFromFloat(0.7),
		}

		pricelist := map[Asset]decimal.Decimal{
			"ETH": decimal.NewFromFloat(200),
			"BTC": decimal.NewFromFloat(5000),
		}

		Account := NewAccount(holdings, pricelist)

		_, _ = Account.Balance(targetIndex)
	}
}
