# balancer
[![Build Status](https://travis-ci.com/pdbrito/balancer.png?branch=master)](https://travis-ci.com/pdbrito/balancer) [![GoDoc](https://godoc.org/github.com/pdbrito/balancer?status.svg)](https://godoc.org/github.com/pdbrito/balancer) [![Go Report Card](https://goreportcard.com/badge/github.com/pdbrito/balancer)](https://goreportcard.com/report/github.com/pdbrito/balancer) [![Codecov](https://codecov.io/gh/pdbrito/balancer/branch/master/graphs/badge.svg)](https://codecov.io/gh/pdbrito/balancer/branch/master/)

Balancer provides guidance with the task of [rebalancing assets](https://en.wikipedia.org/wiki/Rebalancing_investments). 

## Examples

### Balance

Let's assume the current price of 1 BTC is $5000 and the current price of 1 ETH is $200.

First set the global pricelist to reflect these prices:

```go
err := SetPricelist(Pricelist{
	"ETH": decimal.NewFromFloat(200),
	"BTC": decimal.NewFromFloat(5000),
})

if err != nil {
	log.Fatalf("unexpected error whilst setting pricelist: %v", err)
}
```

If you're holding 20 ETH and 0.5 BTC, you can model your account like this:

```go
account, err := NewAccount(Holdings{
	"ETH": decimal.NewFromFloat(20),
	"BTC": decimal.NewFromFloat(0.5),
})

if err != nil {
	log.Fatalf("unexpected error whilst creating account: %v", err)
}
```

The current value of all your assets is:

```
0.5 x 5000 + 20 x 200 = 6500
```

The current percentage of each asset is:

```
ETH = 20 * 200 / 6500 = 0.615384...
BTC = 0.5 * 5000 / 6500 = 0.384615...
```

If you wanted to change this to a 50/50 split, we need to model a target index:

```go
targetIndex := Index{
	"ETH": decimal.NewFromFloat(0.5),
	"BTC": decimal.NewFromFloat(0.5),
})
```

You can then pass `targetIndex` to your `account.Balance()` and you'll receive the trades necessary to rebalance your 
holdings as a `map[Asset]Trade`.

```go
requiredTrades, err := account.Balance(targetIndex)

if err != nil {
	log.Fatalf("unexpected error whilst balancing account: %v", err)
}

for asset, trade := range requiredTrades {
	fmt.Printf("%s %s %s\n", trade.Action, trade.Amount, asset)
}

// Unordered output:
// sell 3.75 ETH
// buy 0.15 BTC
```

### Balancing into new assets

You can also balance your current holdings into other new 
assets, as long as these new assets are included in the global pricelist:

```go
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

account, err := NewAccount(Holdings{
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

requiredTrades, err := account.Balance(targetIndex)

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
```