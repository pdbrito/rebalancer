package balancer_test

import (
	. "github.com/pdbrito/balancer"
	"reflect"
	"testing"
)

func TestBalancer_Balance(t *testing.T) {
	holdings := map[Asset]Holding{
		"eth": {20, 200},
		"btc": {0.5, 5000},
	}
	index := map[Asset]float64{
		"eth": 0.3,
		"btc": 0.7,
	}

	got := Balance(holdings, index)
	want := map[Asset]Trade{
		"eth": {"sell", "10.25"},
		"btc": {"buy", "0.41"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func BenchmarkBalancer_Balance(b *testing.B) {
	for i := 0; i < b.N; i++ {
		holdings := map[Asset]Holding{
			"eth": {20, 200},
			"btc": {0.5, 5000},
		}
		index := map[Asset]float64{
			"eth": 0.3,
			"btc": 0.7,
		}

		got := Balance(holdings, index)
		want := map[Asset]Trade{
			"eth": {"sell", "10.25"},
			"btc": {"buy", "0.41"},
		}

		if !reflect.DeepEqual(got, want) {
			b.Errorf("got %v want %v", got, want)
		}
	}
}
