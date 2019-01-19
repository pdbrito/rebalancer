# balancer
[![Build Status](https://travis-ci.com/pdbrito/balancer.png?branch=master)](https://travis-ci.com/pdbrito/balancer) [![GoDoc](https://godoc.org/github.com/pdbrito/balancer?status.svg)](https://godoc.org/github.com/pdbrito/balancer) [![Go Report Card](https://goreportcard.com/badge/github.com/pdbrito/balancer)](https://goreportcard.com/report/github.com/pdbrito/balancer) [![Codecov](https://codecov.io/gh/pdbrito/balancer/branch/master/graphs/badge.svg)](https://codecov.io/gh/pdbrito/balancer/branch/master/)

Balancer provides guidance with the task of [rebalancing assets](https://en.wikipedia.org/wiki/Rebalancing_investments). 

## Examples

### Balance

Imagine you own 0.5 BTC and 20 ETH.

Imagine the current price of 1 BTC is $5000 and the current price of 1 ETH is $350.

You can model your assets thusly:

```go
holdings := map[Asset]Holding{
    "ETH": {
        Amount: decimal.NewFromFloat(20),
        Price:  decimal.NewFromFloat(350),
    },
    "BTC": {
        Amount: decimal.NewFromFloat(0.5),
        Price:  decimal.NewFromFloat(50000)},
    }
}
```

The current value of all your assets is:

```
0.5 x 5000 + 20 x 350 = 9500
```

The current weighting of each asset is:

```
ETH = 20 * 350 / 9500 = 0.736842...
BTC = 0.5 * 2500 / 9500 = 0.263157...
```

If you wanted to change this to a 50/50 split, first model your target weights:

```go
desiredWeights := map[Asset]decimal.Decimal{
    "ETH": decimal.NewFromFloat(0.5),
    "BTC": decimal.NewFromFloat(0.5),
}
```

Creating a new account from your holdings and calling the Balance method on the
account will return the trades necessary to rebalance your portfolio as a 
`map[Asset]Trade`.

```go
Account := balance.NewAccount(holdings)
requiredTrades := Account.Balance(desiredWeights)
    
for asset, trade := range requiredTrades {
	fmt.Printf("%s %s %s\n", trade.Action, trade.Amount, asset)
}
	
// sell 6.4285714285714286 ETH
// buy 0.45 BTC  
```

### BalanceNew

BalanceNew allows you to balance one or more holdings into several other new 
assets, as long as these new assets are included in a pricelist and passed 
through:

```go
holdings := map[Asset]Holding{
    "ETH": {
        Amount: decimal.NewFromFloat(42),
        Price:  decimal.NewFromFloat(200),
    },
}

desiredWeights := map[Asset]decimal.Decimal{
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

Account := balance.NewAccount(holdings)

requiredTrades := Account.BalanceNew(desiredWeights, pricelist)

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