package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type simulation struct {
	// Number of car washing machines
	numWashes int
	// Time taken to process a car
	timeToProcessCar time.Duration
	// Number of cars arriving per second
	arrivalRate int
	// Number of cars to process in the simulation
	iterations int

	// Some stats on the time taken per car, including queuing time
	stats
}

type car struct {
	arrival time.Time
	done    time.Time
}

func (sim *simulation) carwash(wg *sync.WaitGroup, carQueue <-chan *car) {
	defer wg.Done()
	for c := range carQueue {
		sim.processCar(c)
		sim.done(c)
	}
}

func (sim *simulation) processCar(c *car) {
	time.Sleep(sim.timeToProcessCar)
}

func (sim *simulation) done(c *car) {
	c.done = time.Now()
	sim.record(c)
}

func (sim *simulation) run() {

	// Start the car washes
	carQueue := make(chan *car, sim.numWashes*100)
	wg := &sync.WaitGroup{}
	wg.Add(sim.numWashes)
	for i := 0; i < sim.numWashes; i++ {
		go sim.carwash(wg, carQueue)
	}

	// generate random cars
	start := time.Now()
	for i := 0; i < sim.iterations; i++ {
		// Wait for the next car to arrive and add it to the queue
		time.Sleep(time.Duration(rand.ExpFloat64()) * time.Second / time.Duration(sim.arrivalRate))
		c := &car{arrival: time.Now()}
		carQueue <- c
	}

	// Shutdown
	close(carQueue)
	wg.Wait()
	duration := time.Since(start)
	fmt.Printf("%d cars in %s. %s per car\n", sim.iterations, duration, duration/time.Duration(sim.iterations))
	fmt.Printf("Processing time min=%s, max=%s, mean=%s, stddev=%s\n", sim.min, sim.max, sim.mean(), sim.stdev())
	processingTime := time.Duration(sim.iterations) * sim.timeToProcessCar
	fmt.Printf("Total processing time %s. %.2f%% usage\n", processingTime, (100*float64(processingTime))/(float64(sim.numWashes)*float64(duration)))
}
