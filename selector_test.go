package wrr

import (
	"context"
	"fmt"
	"math"
	"testing"
)

func Test_Selector_Pick(t *testing.T) {
	type args struct {
		ctx      context.Context
		selector Selector
		nodes    []WeightedNode
	}
	tests := []struct {
		name           string
		d              *defaultSelector
		args           args
		Round          int
		wantErr        bool
		PrintGotPicked bool
	}{
		{
			name:    "Pick with  default",
			Round:   1000,
			wantErr: false,
			args: args{
				context.Background(),
				NewDefaultSelector(),
				[]WeightedNode{
					NewWeightedNode(NewNode("10.1.6.1"), WithWeigtOpt(6)),
					NewWeightedNode(NewNode("10.1.1.1"), WithWeigtOpt(1)),
					NewWeightedNode(NewNode("10.1.1.2"), WithWeigtOpt(2)),
				},
			},
		},
		{
			name:           "Pick with wrr",
			Round:          1000,
			PrintGotPicked: false,
			wantErr:        false,
			args: args{
				context.Background(),
				NewWeightedSelector(),
				[]WeightedNode{
					NewWeightedNode(NewNode("10.1.6.1"), WithWeigtOpt(6)),
					NewWeightedNode(NewNode("10.1.1.1"), WithWeigtOpt(1)),
					NewWeightedNode(NewNode("10.1.1.2"), WithWeigtOpt(2)),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.args.selector
			i := 0
			counter := make(map[string]int)
			for i < tt.Round {
				got, err := d.Pick(tt.args.ctx, tt.args.nodes)
				if (err != nil) != tt.wantErr {
					t.Errorf("defaultSelector.Pick() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				i++
				counter[got.Addr()]++
				if tt.PrintGotPicked && tt.Round <= 50 {
					t.Logf("round:%v, picked %s", i, got.Addr())
				}
			}
			passed := probCheck(weightProb(tt.args.nodes), pickProbs(counter, tt.Round))
			t.Logf("selector:%v, after %d times, Pick passed check:%v", tt.args.selector.Name(), tt.Round, passed)
		})
	}
}

func weightProb(nodes []WeightedNode) map[string]float64 {
	probs := make(map[string]float64)

	var totalWeight float64

	for _, node := range nodes {
		totalWeight += node.Wight()
	}

	for _, node := range nodes {
		probs[node.Addr()] = float64(node.Wight() / totalWeight)
	}
	return probs
}

func pickProbs(counter map[string]int, count int) map[string]float64 {
	probs := make(map[string]float64)
	for k, v := range counter {
		probs[k] = float64(v) / float64(count)
	}
	return probs
}

func probCheck(w map[string]float64, p map[string]float64) bool {
	passed := true
	for k, v := range w {
		if pv, ok := p[k]; ok {
			gap := math.Abs(pv - v)
			fmt.Printf("prob check, addres:%v, weightProb:%v, PickProb:%v, gap:%v\n", k, v, pv, gap)
			if gap > 0.1 {
				passed = false
			}
		}
	}
	return passed
}
