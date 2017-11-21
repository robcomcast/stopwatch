package main

import (
	"fmt"
	"math/rand"
	"time"

	"code.comcast.com/vader/stopwatch"
)

func main() {

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
	dur := sw.Split("Split1")
	fmt.Printf("Split1 Time: [%f]\n", dur)

	// Sleep some more to simulate more work
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(3000)))
	fmt.Printf("Split2 Time: [%f]\n", sw.Split("Split2"))

	// Sleep some more to simulate more work
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(3000)))

	// Pause the stopwatch and record the next measurement
	// Pause is useful during code execution that should not be timed
	fmt.Println("End Split3, Pausing Stopwatch")
	sw.Pause("Split3")

	time.Sleep(time.Millisecond * time.Duration(rand.Intn(3000)))
	// This will be ingored by the stopwatch, you can't Split while paused.
	sw.Split("Bad Spilt")

	// Resume the Stopwatch
	fmt.Println("Resuming Stopwatch")
	sw.Resume()

	// Sleep some more to simulate more work
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(3000)))
	sw.Pause("Split4")
	fmt.Println("End Split4, Pausing Stopwatch")

	// Print all of the splits
	fmt.Println("\nSplits:")
	fmt.Printf("%v\n\n", sw.Splits)

	// Get the duration of a specific split:
	fmt.Printf("Duration of Split3: [%f]\n\n", sw.Splits["Split3"].ElapsedMS)

	// Calculate a virtual split time from Split1.StartTime to Split3.EndTime
	elapsed, _ := sw.ElapsedMS("Split1", "Split3")
	fmt.Printf("Virtual Split Time, Split1 to Split3: [%f]\n\n", elapsed)

}
