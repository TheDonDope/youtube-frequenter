package main

import (
	"fmt"
	"math/rand"
	"time"

	neural "github.com/NOX73/go-neural"
	"github.com/NOX73/go-neural/learn"
	"github.com/NOX73/go-neural/persist"
)

func main() {

	// network := persist.FromFile("train150.json")
	network := neural.NewNetwork(3, []int{3, 10, 1})
	network.RandomizeSynapses()

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	good := 1.
	bad := 1.

	for j := 0.; j < 100000; j++ {
		for i := 0.; i < 1000000; i++ {
			// get random x,y,z
			x := r1.Float64()
			y := r1.Float64()
			z := r1.Float64()

			input := []float64{x, y, z}
			var idealOut []float64
			if x+y+z > 1. {
				idealOut = []float64{1.}
			} else {
				idealOut = []float64{0.}
			}

			learn.Learn(network, input, idealOut, 0.5)
			// e := learn.Evaluation(network, input, idealOut)
			// fmt.Println(fmt.Sprintf("(%v/%v): %.5v", i, maxIter, e))
		}
		// save
		persist.ToFile("train150.json", network)

		x := r1.Float64()
		y := r1.Float64()
		z := r1.Float64()
		result := network.Calculate([]float64{x, y, z})

		if x+y+z > 1. && result[0] > 0.5 {
			good++
			fmt.Println(fmt.Sprintf("%.2v good", result[0]))
		} else {
			bad++
			fmt.Println(fmt.Sprintf("%.2v bad", result[0]))
		}
		fmt.Println(fmt.Sprintf("good-ratio: %.4v", good/(good+bad)))
	}
}
