package rebalancer_test

import (
	"fmt"
	. "github.com/pdbrito/rebalancer"
	"github.com/shopspring/decimal"
	"log"
	"reflect"
	"testing"
)

const unexpectedError string = "got an error but didn't want one"
const missingError string = "wanted an error but didn't get one"
const wrongError string = "got an error but expected a different one"

func TestErrInvalidAssetAmount_Error(t *testing.T) {
	asset := Asset("ETH")
	amount := decimal.NewFromFloat(-5)

	err := ErrInvalidAssetAmount{Asset: asset, Amount: amount}

	want := "ETH needs positive amount, not -5"
	got := err.Error()

	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}

func TestSetPricelist(t *testing.T) {
	t.Run("a new pricelist can be set", func(t *testing.T) {
		err := SetPricelist(map[Asset]decimal.Decimal{
			"ETH": decimal.NewFromFloat(200),
			"BTC": decimal.NewFromFloat(5000),
		})

		if err != nil {
			t.Error(unexpectedError)
		}
	})
	t.Run("an empty pricelist cannot be set", func(t *testing.T) {
		err := SetPricelist(map[Asset]decimal.Decimal{})

		if err == nil {
			t.Error(missingError)
		}

		if err != ErrEmptyPricelist {
			t.Error(wrongError)
		}
	})
	t.Run("pricelist asset keys must be uppercase", func(t *testing.T) {
		err := SetPricelist(map[Asset]decimal.Decimal{
			"ETH": decimal.NewFromFloat(200),
			"btc": decimal.NewFromFloat(5000),
		})

		if err == nil {
			t.Error(missingError)
		}

		if err != ErrInvalidAsset {
			t.Error(wrongError)
		}
	})
	t.Run("pricelist entries must have a value above 0", func(t *testing.T) {
		invalidAsset := Asset("BTC")
		invalidAmount := decimal.NewFromFloat(-5)

		err := SetPricelist(map[Asset]decimal.Decimal{
			"ETH":        decimal.NewFromFloat(200),
			invalidAsset: invalidAmount,
		})

		want := ErrInvalidAssetAmount{Asset: invalidAsset, Amount: invalidAmount}

		if err != want {
			t.Errorf("got %v, want %v", err, want)
		}
	})
}

func TestGlobalPricelist(t *testing.T) {
	t.Run("it returns the current value of the global pricelist", func(t *testing.T) {
		pricelist := map[Asset]decimal.Decimal{
			"ETH": decimal.NewFromFloat(222),
			"BTC": decimal.NewFromFloat(5555),
		}

		err := SetPricelist(pricelist)

		if err != nil {
			t.Error(unexpectedError)
		}

		got := GlobalPricelist()
		want := Pricelist(pricelist)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestClearGlobalPricelist(t *testing.T) {
	t.Run("it clears the value of the global pricelist", func(t *testing.T) {
		err := SetPricelist(map[Asset]decimal.Decimal{
			"ETH": decimal.NewFromFloat(5),
		})

		if err != nil {
			t.Error(unexpectedError)
		}

		ClearGlobalPricelist()

		got := GlobalPricelist()
		want := Pricelist{}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestNewPortfolio(t *testing.T) {
	t.Run("portfolio cannot contain an empty map", func(t *testing.T) {
		_, err := NewPortfolio(map[Asset]decimal.Decimal{})

		if err != ErrEmptyPortfolio {
			t.Errorf("got %v want %v", err, ErrEmptyPortfolio)
		}
	})
	t.Run("portfolio cannot contain invalid asset keys", func(t *testing.T) {
		_, err := NewPortfolio(map[Asset]decimal.Decimal{
			"eth": decimal.NewFromFloat(5),
		})

		if err != ErrInvalidAsset {
			t.Errorf("got %v want %v", err, ErrInvalidAsset)
		}
	})
	t.Run("portfolio cannot contain assets missing from the global pricelist", func(t *testing.T) {
		err := SetPricelist(map[Asset]decimal.Decimal{
			"ETH": decimal.NewFromFloat(200),
		})

		if err != nil {
			t.Error(unexpectedError)
		}

		_, err = NewPortfolio(map[Asset]decimal.Decimal{
			"BTC": decimal.NewFromFloat(5),
		})

		if err != ErrAssetMissingFromPricelist {
			t.Errorf("got %v, want %v", err, ErrAssetMissingFromPricelist)
		}
	})
	t.Run("portfolio cannot contain values of zero or less", func(t *testing.T) {
		err := SetPricelist(map[Asset]decimal.Decimal{
			"ETH": decimal.NewFromFloat(200),
		})

		if err != nil {
			t.Error(unexpectedError)
		}

		asset := Asset("ETH")
		amount := decimal.NewFromFloat(-5)

		_, err = NewPortfolio(map[Asset]decimal.Decimal{
			asset: amount,
		})

		want := ErrInvalidAssetAmount{Asset: asset, Amount: amount}

		if err != want {
			t.Errorf("got %v, want %v", err, want)
		}
	})
	t.Run("a new portfolio can be created", func(t *testing.T) {
		got, err := NewPortfolio(map[Asset]decimal.Decimal{
			"ETH": decimal.NewFromFloat(5),
		})

		if err != nil {
			t.Error(unexpectedError)
		}

		want := Portfolio{"ETH": decimal.NewFromFloat(5)}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

func TestNewAccount(t *testing.T) {
	t.Run("account cannot be created if the global pricelist is empty", func(t *testing.T) {
		ClearGlobalPricelist()

		portfolio := Portfolio{
			"ETH": decimal.NewFromFloat(5),
			"BTC": decimal.NewFromFloat(0.5),
		}

		_, err := NewAccount(portfolio)

		if err != ErrEmptyPricelist {
			t.Error(wrongError)
		}
	})
	t.Run("account cannot contain invalid asset keys in its portfolio", func(t *testing.T) {
		_ = SetPricelist(map[Asset]decimal.Decimal{
			"ETH": decimal.NewFromFloat(200),
			"BTC": decimal.NewFromFloat(5000),
		})

		invalidAsset := Asset("ETH")
		invalidAmount := decimal.NewFromFloat(-5)

		portfolio := Portfolio{
			invalidAsset: invalidAmount,
			"BTC":        decimal.NewFromFloat(0.5),
		}

		_, err := NewAccount(portfolio)

		want := ErrInvalidAssetAmount{Asset: invalidAsset, Amount: invalidAmount}

		if err != want {
			t.Error(wrongError)
		}
	})
	t.Run("account cannot contain empty portfolio", func(t *testing.T) {
		_ = SetPricelist(map[Asset]decimal.Decimal{
			"ETH": decimal.NewFromFloat(200),
			"BTC": decimal.NewFromFloat(5000),
		})

		portfolio := Portfolio{}

		_, err := NewAccount(portfolio)

		if err == nil {
			t.Error(missingError)
		}
	})
	t.Run("account cannot contain invalid asset keys in its portfolio", func(t *testing.T) {
		_ = SetPricelist(map[Asset]decimal.Decimal{
			"ETH": decimal.NewFromFloat(200),
			"BTC": decimal.NewFromFloat(5000),
		})

		portfolio := Portfolio{
			"eth": decimal.NewFromFloat(5),
			"BTC": decimal.NewFromFloat(0.5),
		}

		_, err := NewAccount(portfolio)

		if err == nil {
			t.Error(missingError)
		}
	})
	t.Run("a new account can be created", func(t *testing.T) {
		_ = SetPricelist(map[Asset]decimal.Decimal{
			"ETH": decimal.NewFromFloat(200),
			"BTC": decimal.NewFromFloat(5000),
		})

		portfolio := Portfolio{
			"ETH": decimal.NewFromFloat(20),
			"BTC": decimal.NewFromFloat(0.5),
		}

		_, err := NewAccount(portfolio)

		if err != nil {
			t.Error(unexpectedError)
		}
	})
}

func TestNewIndex(t *testing.T) {
	t.Run("index cannot contain an empty map", func(t *testing.T) {
		_, err := NewIndex(map[Asset]decimal.Decimal{})

		if err != ErrEmptyIndex {
			t.Error(wrongError)
		}
	})
	t.Run("index cannot contain invalid asset keys", func(t *testing.T) {
		_, err := NewIndex(map[Asset]decimal.Decimal{
			"ETH": decimal.NewFromFloat(200),
			"btc": decimal.NewFromFloat(5000),
		})

		if err != ErrInvalidAsset {
			t.Error(wrongError)
		}
	})
	t.Run("index cannot contain assets missing from the global pricelist", func(t *testing.T) {
		err := SetPricelist(map[Asset]decimal.Decimal{
			"ETH": decimal.NewFromFloat(200),
		})

		if err != nil {
			t.Error(unexpectedError)
		}

		_, err = NewIndex(map[Asset]decimal.Decimal{
			"BTC": decimal.NewFromFloat(1),
		})

		if err != ErrAssetMissingFromPricelist {
			t.Errorf("got %v, want %v", err, ErrAssetMissingFromPricelist)
		}
	})
	t.Run("index cannot contain values of zero or less", func(t *testing.T) {
		err := SetPricelist(map[Asset]decimal.Decimal{
			"ETH": decimal.NewFromFloat(200),
			"BTC": decimal.NewFromFloat(5000),
		})

		if err != nil {
			t.Error(unexpectedError)
		}

		invalidAsset := Asset("BTC")
		invalidAmount := decimal.NewFromFloat(-5)

		_, err = NewIndex(map[Asset]decimal.Decimal{
			"ETH":        decimal.NewFromFloat(200),
			invalidAsset: invalidAmount,
		})

		want := ErrInvalidAssetAmount{Asset: invalidAsset, Amount: invalidAmount}

		if err != want {
			t.Errorf("got %v, want %v", err, want)
		}
	})
	t.Run("index values must sum to 1", func(t *testing.T) {
		_, err := NewIndex(map[Asset]decimal.Decimal{
			"ETH": decimal.NewFromFloat(0.2),
			"BTC": decimal.NewFromFloat(0.2),
		})

		if err != ErrIndexSumIncorrect {
			t.Error(wrongError)
		}
	})
	t.Run("a new index can be created", func(t *testing.T) {
		got, err := NewIndex(map[Asset]decimal.Decimal{
			"ETH": decimal.NewFromFloat(0.5),
			"BTC": decimal.NewFromFloat(0.5),
		})

		if err != nil {
			t.Error(unexpectedError)
		}

		want := Index{
			"ETH": decimal.NewFromFloat(0.5),
			"BTC": decimal.NewFromFloat(0.5),
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

func TestAccount_Rebalance(t *testing.T) {
	t.Run("rebalance cannot receive an empty index", func(t *testing.T) {
		err := SetPricelist(map[Asset]decimal.Decimal{
			"ETH": decimal.NewFromFloat(200),
			"BTC": decimal.NewFromFloat(5000),
		})

		if err != nil {
			t.Error(unexpectedError)
		}

		account, err := NewAccount(Portfolio{
			"ETH": decimal.NewFromFloat(20),
			"BTC": decimal.NewFromFloat(0.5),
		})

		if err != nil {
			t.Error(unexpectedError)
		}

		_, err = account.Rebalance(Index{})

		if err != ErrEmptyIndex {
			t.Errorf("got %v, want %v", err, ErrEmptyIndex)
		}
	})
	t.Run("rebalance cannot receive an index with invalid asset keys", func(t *testing.T) {
		err := SetPricelist(map[Asset]decimal.Decimal{
			"ETH": decimal.NewFromFloat(200),
			"BTC": decimal.NewFromFloat(5000),
		})

		if err != nil {
			t.Error(unexpectedError)
		}

		account, err := NewAccount(Portfolio{
			"ETH": decimal.NewFromFloat(20),
			"BTC": decimal.NewFromFloat(0.5),
		})

		if err != nil {
			t.Error(unexpectedError)
		}

		_, err = account.Rebalance(Index{
			"btc": decimal.NewFromFloat(0.5),
		})

		if err != ErrInvalidAsset {
			t.Errorf("got %v, want %v", err, ErrInvalidAsset)
		}
	})
	t.Run("rebalance cannot receive an index with assets missing from the global pricelist", func(t *testing.T) {
		err := SetPricelist(map[Asset]decimal.Decimal{
			"ETH": decimal.NewFromFloat(200),
		})

		if err != nil {
			t.Error(unexpectedError)
		}

		account, err := NewAccount(Portfolio{
			"ETH": decimal.NewFromFloat(20),
		})

		if err != nil {
			t.Error(unexpectedError)
		}

		_, err = account.Rebalance(Index{
			"BTC": decimal.NewFromFloat(0.5),
		})

		if err != ErrAssetMissingFromPricelist {
			t.Errorf("got %v, want %v", err, ErrAssetMissingFromPricelist)
		}
	})
	t.Run("rebalance cannot receive an index with values of zero or less", func(t *testing.T) {
		err := SetPricelist(map[Asset]decimal.Decimal{
			"ETH": decimal.NewFromFloat(200),
			"BTC": decimal.NewFromFloat(5000),
		})

		if err != nil {
			t.Error(unexpectedError)
		}

		account, err := NewAccount(Portfolio{
			"ETH": decimal.NewFromFloat(20),
			"BTC": decimal.NewFromFloat(0.5),
		})

		if err != nil {
			t.Error(unexpectedError)
		}

		invalidAsset := Asset("ETH")
		invalidAmount := decimal.NewFromFloat(-0.3)

		_, err = account.Rebalance(Index{
			invalidAsset: invalidAmount,
			"BTC":        decimal.NewFromFloat(0.7),
		})

		want := ErrInvalidAssetAmount{Asset: invalidAsset, Amount: invalidAmount}

		if err != want {
			t.Errorf("got %v, want %v", err, want)
		}
	})
	t.Run("rebalance cannot receive an index whose values don't sum to 1", func(t *testing.T) {
		err := SetPricelist(map[Asset]decimal.Decimal{
			"ETH": decimal.NewFromFloat(200),
			"BTC": decimal.NewFromFloat(5000),
		})

		if err != nil {
			t.Error(unexpectedError)
		}

		account, err := NewAccount(Portfolio{
			"ETH": decimal.NewFromFloat(20),
			"BTC": decimal.NewFromFloat(0.5),
		})

		if err != nil {
			t.Error(unexpectedError)
		}

		_, err = account.Rebalance(Index{
			"BTC": decimal.NewFromFloat(0.7),
			"ETH": decimal.NewFromFloat(0.7),
		})

		if err != ErrIndexSumIncorrect {
			t.Errorf("got %v, want %v", err, ErrIndexSumIncorrect)
		}
	})
	t.Run("rebalance can rebalance an account", func(t *testing.T) {
		err := SetPricelist(map[Asset]decimal.Decimal{
			"ETH": decimal.NewFromFloat(200),
			"BTC": decimal.NewFromFloat(5000),
		})

		if err != nil {
			t.Error(unexpectedError)
		}

		account, err := NewAccount(Portfolio{
			"ETH": decimal.NewFromFloat(20),
			"BTC": decimal.NewFromFloat(0.5),
		})

		if err != nil {
			t.Error(unexpectedError)
		}

		got, err := account.Rebalance(Index{
			"ETH": decimal.NewFromFloat(0.3),
			"BTC": decimal.NewFromFloat(0.7),
		})

		if err != nil {
			t.Error(unexpectedError)
		}

		want := map[Asset]Trade{
			"ETH": {Action: "sell", Amount: decimal.NewFromFloat(10.25)},
			"BTC": {Action: "buy", Amount: decimal.NewFromFloat(0.41)},
		}

		assertSameTrades(t, got, want)
	})
	t.Run("rebalance can rebalance existing assets into new assets", func(t *testing.T) {
		err := SetPricelist(map[Asset]decimal.Decimal{
			"ETH":  decimal.NewFromFloat(200),
			"BTC":  decimal.NewFromFloat(2000),
			"IOTA": decimal.NewFromFloat(0.3),
			"BAT":  decimal.NewFromFloat(0.12),
			"XLM":  decimal.NewFromFloat(0.2),
		})

		if err != nil {
			t.Error(unexpectedError)
		}

		portfolio := map[Asset]decimal.Decimal{
			"ETH": decimal.NewFromFloat(42),
		}

		targetIndex := map[Asset]decimal.Decimal{
			"ETH":  decimal.NewFromFloat(0.2),
			"BTC":  decimal.NewFromFloat(0.2),
			"IOTA": decimal.NewFromFloat(0.2),
			"BAT":  decimal.NewFromFloat(0.2),
			"XLM":  decimal.NewFromFloat(0.2),
		}

		account, err := NewAccount(portfolio)

		if err != nil {
			t.Error(unexpectedError)
		}

		got, err := account.Rebalance(targetIndex)

		if err != nil {
			t.Error(unexpectedError)
		}

		want := map[Asset]Trade{
			"ETH":  {Action: "sell", Amount: decimal.NewFromFloat(33.6)},
			"BTC":  {Action: "buy", Amount: decimal.NewFromFloat(0.84)},
			"IOTA": {Action: "buy", Amount: decimal.NewFromFloat(5600)},
			"BAT":  {Action: "buy", Amount: decimal.NewFromFloat(14000)},
			"XLM":  {Action: "buy", Amount: decimal.NewFromFloat(8400)},
		}

		assertSameTrades(t, got, want)
	})
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

func ExampleAccount_Rebalance() {
	err := SetPricelist(Pricelist{
		"ETH": decimal.NewFromFloat(200),
		"BTC": decimal.NewFromFloat(5000),
	})

	if err != nil {
		log.Fatalf("unexpected error whilst setting pricelist: %v", err)
	}

	account, err := NewAccount(Portfolio{
		"ETH": decimal.NewFromFloat(20),
		"BTC": decimal.NewFromFloat(0.5),
	})

	if err != nil {
		log.Fatalf("unexpected error whilst creating account: %v", err)
	}

	targetIndex := Index{
		"ETH": decimal.NewFromFloat(0.5),
		"BTC": decimal.NewFromFloat(0.5),
	}

	requiredTrades, err := account.Rebalance(targetIndex)

	if err != nil {
		log.Fatalf("unexpected error whilst balancing account: %v", err)
	}

	for asset, trade := range requiredTrades {
		fmt.Printf("%s %s %s\n", trade.Action, trade.Amount, asset)
	}

	// Unordered output:
	// sell 3.75 ETH
	// buy 0.15 BTC
}

func ExampleAccount_Rebalance_intoNewAssets() {
	err := SetPricelist(Pricelist{
		"ETH":  decimal.NewFromFloat(200),
		"BTC":  decimal.NewFromFloat(2000),
		"IOTA": decimal.NewFromFloat(0.3),
		"BAT":  decimal.NewFromFloat(0.12),
		"XLM":  decimal.NewFromFloat(0.2),
	})

	if err != nil {
		log.Fatalf("unexpected error whilst setting pricelist: %v", err)
	}

	account, err := NewAccount(Portfolio{
		"ETH": decimal.NewFromFloat(42),
	})

	if err != nil {
		log.Fatalf("unexpected error whilst creating account: %v", err)
	}

	targetIndex := Index{
		"ETH":  decimal.NewFromFloat(0.2),
		"BTC":  decimal.NewFromFloat(0.2),
		"IOTA": decimal.NewFromFloat(0.2),
		"BAT":  decimal.NewFromFloat(0.2),
		"XLM":  decimal.NewFromFloat(0.2),
	}

	requiredTrades, err := account.Rebalance(targetIndex)

	if err != nil {
		log.Fatalf("unexpected error whilst balancing account: %v", err)
	}

	for asset, trade := range requiredTrades {
		fmt.Printf("%s %s %s\n", trade.Action, trade.Amount, asset)
	}

	// Unordered output:
	// sell 33.6 ETH
	// buy 0.84 BTC
	// buy 5600 IOTA
	// buy 14000 BAT
	// buy 8400 XLM
}

func BenchmarkRebalance(b *testing.B) {
	_ = SetPricelist(map[Asset]decimal.Decimal{
		"ETH": decimal.NewFromFloat(200),
		"BTC": decimal.NewFromFloat(5000),
	})

	for i := 0; i < b.N; i++ {
		portfolio := map[Asset]decimal.Decimal{
			"ETH": decimal.NewFromFloat(20),
			"BTC": decimal.NewFromFloat(0.5),
		}
		targetIndex := map[Asset]decimal.Decimal{
			"ETH": decimal.NewFromFloat(0.3),
			"BTC": decimal.NewFromFloat(0.7),
		}

		Account, _ := NewAccount(portfolio)

		_, _ = Account.Rebalance(targetIndex)
	}
}
