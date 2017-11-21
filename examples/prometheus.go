package main

import (
	"fmt"
	"math/rand"
	"time"

	"code.comcast.com/vader/stopwatch"
	"github.com/prometheus/client_golang/prometheus"
)

func main() {

	// Create Prometheus Counters.
	Split1 := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "split1_duration_ms",
			Help: "Duration (ms) of Split 1.",
		},
	)

	Split2 := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "split2_duration_ms",
			Help: "Duration (ms) of Split 2.",
		},
	)

	Split3 := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "split3_duration_ms",
			Help: "Duration (ms) of Split 3.",
		},
	)

	Split4 := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "split4_duration_ms",
			Help: "Duration (ms) of Split 4.",
		},
	)

	prometheus.MustRegister(Split1)
	prometheus.MustRegister(Split2)
	prometheus.MustRegister(Split3)
	prometheus.MustRegister(Split4)

	// Initiate Randomizer for Work Simulation
	rand.Seed(time.Now().UnixNano())

	// Create the stopwatch wtih .New() and later call .Start()
	// or use .AutoStart() to return a started stopwatch
	fmt.Println("Starting Stopwatch")
	sw := stopwatch.AutoStart()

	// Sleep to create a duration for the first split
	// Ideally you would be doing some kind of work here.
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(3000)))

	// Name and mark a split to record the first measurement
	// Duration is saved to the stopwatch and returned in Milliseconds
	fmt.Println("End Split1")
	dur := sw.Split("Split1")
	// Add to the prometheus counter
	Split1.Add(dur)

	// Sleep some more to simulate more work
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(3000)))

	// Name and mark a split to record the next measurement
	// You call Split() as a parm to Prometheus.Counter.Add()
	fmt.Println("End Split2")
	Split2.Add(sw.Split("Split2"))

	// Sleep some more to simulate more work
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(3000)))

	// Pause the stopwatch and record the next measurement
	// Pause is useful to create gaps between measurements
	// Restart with Resume()
	fmt.Println("End Split3")
	sw.Pause("Split3")

	// This will be ingored by the stopwatch
	// You can't Split while paused.
	sw.Split("Bad Spilt")

	// Resume the Stopwatch
	sw.Resume()

	// Sleep some more to simulate more work
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(3000)))
	fmt.Println("End Split4")
	sw.Pause("Split4")

	// Print all of the splits
	fmt.Println("Splits:")
	fmt.Printf("%v\n\n", sw.Splits)

	// Calculate a virtual split time from Split1.StartTime to Split3.EndTime
	fmt.Println("Virtual Split Time, Split1 to Split3:")
	elapsed, _ := sw.ElapsedMS("Split1", "Split3")
	fmt.Printf("%f\n", elapsed)

}
