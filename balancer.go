package balancer

import (
	"fmt"
	"math"
)

type Asset string

type Holding struct {
	Quantity float64
	Value  float64
}

type Trade struct {
	Action string
	Amount string
}

func Balance(holdings map[Asset]Holding, index map[Asset]float64) map[Asset]Trade {
	//validate assumptions; only unique assets etc
	totalHoldings := float64(0)
	for _, holding := range holdings {
		totalHoldings += holding.Value * holding.Quantity
	}

	targetAmount := map[Asset]float64{}
	targetHoldings := map[Asset]float64{}

	for asset, weight := range index {
		targetHoldings[asset] = totalHoldings * weight
	}

	for asset, targetHolding := range targetHoldings {
		targetAmount[asset] = targetHolding / holdings[asset].Value
	}

	trades := map[Asset]Trade{}
	for asset, amount := range targetAmount {
		amountRequired := amount - holdings[asset].Quantity
		trades[asset] = makeTrade(amountRequired)
	}

	return trades
}
func makeTrade(amount float64) Trade {
	var action string
	if amount < 0 {
		action = "sell"
	} else {
		action = "buy"
	}
	return Trade{action, fmt.Sprintf("%.2f", math.Abs(amount))}
}