package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/plotutil"
)

type options struct {
	washes       int
	time         int
	arrivalsMin  int
	arrivalsIncr int
	arrivalsMax  int
	iterations   int
}

func (o *options) init() {
	flag.IntVar(&o.washes, "washes", 1, "number of washes (or CPUs)")
	flag.IntVar(&o.time, "time", 20, "number of milliseconds to wash a car")
	flag.IntVar(&o.arrivalsMin, "arrivalsMin", 10, "number of arrivals per second (min)")
	flag.IntVar(&o.arrivalsMax, "arrivalsMax", 10, "number of arrivals per second (max")
	flag.IntVar(&o.arrivalsIncr, "arrivalsIncr", 10, "number of arrivals per second (increment)")
	flag.IntVar(&o.iterations, "iterations", 100, "number of iterations of model")
}

func main() {
	o := &options{}
	o.init()
	flag.Parse()

	maxes := plotter.XYs{}
	means := plotter.XYs{}

	// Ah, yes, ignoring the arrivals parameter and
	for arrivals := o.arrivalsMin; arrivals <= o.arrivalsMax; arrivals += o.arrivalsIncr {
		sim := &simulation{
			numWashes: o.washes,
			// time to process
			timeToProcessCar: time.Duration(o.time) * time.Millisecond,
			// Arrivals per second
			arrivalRate: arrivals,

			iterations: o.iterations,
		}

		sim.run()

		maxes = append(maxes, struct{ X, Y float64 }{
			X: float64(arrivals),
			Y: float64(sim.max / time.Millisecond),
		})

		means = append(means, struct{ X, Y float64 }{
			X: float64(arrivals),
			Y: float64(sim.mean() / time.Millisecond),
		})
	}

	p, err := plot.New()
	if err != nil {
		fmt.Printf("Failed to get a new plot %v\n", err)
	}
	p.X.Label.Text = "Requests per second"
	p.Y.Label.Text = "Response time ms"
	if err := plotutil.AddLinePoints(p, "Max response time (ms)", maxes, "Mean response time (ms)", means); err != nil {
		fmt.Printf("failed to add lines. %v\n", err)
	}
	err = p.Save(400, 300, "plot.png")
	if err != nil {
		fmt.Printf("plot failed. %v\n", err)
	}
}
