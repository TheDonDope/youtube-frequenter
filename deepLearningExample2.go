package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/NOX73/go-neural/learn"
	"github.com/NOX73/go-neural/persist"
)

func main() {

	network := persist.FromFile("train_sin_2_10_1.json")
	// network := neural.NewNetwork(2, []int{2, 10, 1})
	// network.RandomizeSynapses()
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	if true {
		upTo := 100000000
		for i := 0; i < upTo; i++ {
			x := r1.Float64()
			y := r1.Float64()
			learn.Learn(network, []float64{x, y}, []float64{math.Sin(x + y)}, 0.1)
			fmt.Println(strconv.Itoa(i) + "/" + strconv.Itoa(upTo))
		}
	}
	persist.ToFile("train_sin_2_10_1.json", network)
	// fmt.Println(learn.Evaluation(network, []float64{x, y}, []float64{math.Sin(x + y)}))
	x := r1.Float64()
	y := r1.Float64()
	result := network.Calculate([]float64{x, y})
	idealResult := math.Sin(x + y)
	fmt.Println(fmt.Sprintf("%v should be: %v", result[0], idealResult))
}
