# balancer
[![Build Status](https://travis-ci.com/pdbrito/balancer.png?branch=master)](https://travis-ci.com/pdbrito/balancer) [![GoDoc](https://godoc.org/github.com/pdbrito/balancer?status.svg)](https://godoc.org/github.com/pdbrito/balancer) [![Go Report Card](https://goreportcard.com/badge/github.com/pdbrito/balancer)](https://goreportcard.com/report/github.com/pdbrito/balancer) [![Codecov](https://codecov.io/gh/pdbrito/balancer/branch/master/graphs/badge.svg)](https://codecov.io/gh/pdbrito/balancer/branch/master/)

Balancer provides guidance with the task of [rebalancing assets](https://en.wikipedia.org/wiki/Rebalancing_investments).

### Example

Image you currently own 0.5 BTC and 20 ETH.

If the current price of 1 BTC is $5000 and the current price of 1 ETH is $350

You can model these assets as follows:
```go
myPorfolio := map[Asset]Holding{
    "BTC": {
        decimal.NewFromFloat(0.5),
        decimal.NewFromFloat(5000),
    },
    "ETH": {
        decimal.NewFromFloat(0.5),
        decimal.NewFromFloat(5000)},
    },
}
```