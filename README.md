# balancer
[![Build Status](https://travis-ci.com/pdbrito/balancer.png?branch=master)](https://travis-ci.com/pdbrito/balancer) [![GoDoc](https://godoc.org/github.com/pdbrito/balancer?status.svg)](https://godoc.org/github.com/pdbrito/balancer) [![Go Report Card](https://goreportcard.com/badge/github.com/pdbrito/balancer)](https://goreportcard.com/report/github.com/pdbrito/balancer) [![Codecov](https://codecov.io/gh/pdbrito/balancer/branch/master/graphs/badge.svg)](https://codecov.io/gh/pdbrito/balancer/branch/master/)

Balancer provides guidance with the task of [rebalancing assets](https://en.wikipedia.org/wiki/Rebalancing_investments).

### Example

Imagine you own 0.5 BTC and 20 ETH.

Imagine the current price of 1 BTC is $5000 and the current price of 1 ETH is $350.

You can model your assets thusly:
```go
myPorfolio := map[Asset]Holding{
    "BTC": {
        Amount: decimal.NewFromFloat(0.5),
        Value: decimal.NewFromFloat(5000),
    },
    "ETH": {
        Amount: decimal.NewFromFloat(20),
        Value: decimal.NewFromFloat(350)},
    }
}
```

The current value of all your assets is:

`0.5 x 5000 + 20 x 350 = 9500`

The current weighting of each asset is:

```
BTC = 0.5 * 2500 / 9500 = 0.263157...
ETH = 20 * 350 / 9500 = 0.736842...
```
