# Stopwatch

Stopwatch is a lightweight, threadsafe package for Golang that measures event durations within the importing application. It is intended to easily capture event timings to publish as [Prometheus](https://github.com/prometheus/client_golang/) Counters with little overhead.

Stopwatch utilizes the monotonic system clock for duration calcuations, thus insulating its results from the effects of unexpected clock changes, like leap second adjustments and NTP syncs.

## Overview

Stopwatch works like a physical stopwatch. You can Start() it, Pause() it, and record Split() times, and then do things with the resulting telemetry.

When you have an application performing many time-consuming operations, it is useful to record and trend the time taken to perform those operations. Stopwatch makes this easy by insulating your code from repetitive time calculation exercises. Instead, inject a Split() where you would normally demarcate event boundaries, and issue a Pause() after the last event, and all timings are calculated for you.

## Examples

### Basic Usage
```go
func main() {

	// Initiate Randomizer for Work Simulation
	rand.Seed(time.Now().UnixNano())

	// Create the stopwatch with .New() and later call .Start()
	// or use .AutoStart() to return a newly-started stopwatch
	fmt.Println("Starting Stopwatch")
	sw := stopwatch.AutoStart()

	// Sleep to create a duration for the first split
	// Ideally you would be doing some kind of work here.
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(3000)))

	// Name and mark a Split to record the first measurement
	// Duration is saved to the stopwatch by Split Name, and returned in MS
	dur := sw.Split("Split1")
	fmt.Printf("Split1 Time: [%f]\n", dur)

	// Sleep some more to simulate more work
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(3000)))
	fmt.Printf("Split2 Time: [%f]\n", sw.Split("Split2"))

	// Sleep some more to simulate more work
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(3000)))

	// Pause the stopwatch and record the Split measurement
	// Pause is useful during code execution that should not be timed
	fmt.Println("End Split3, Pausing Stopwatch")
	sw.Pause("Split3")

	// Sleep some more to simulate more work
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(3000)))

	// This will be ingored by the stopwatch, you can't Split while paused.
	sw.Split("Bad Split")

	// Resume the Stopwatch
	fmt.Println("Resuming Stopwatch")
	sw.Resume()

	// Sleep some more to simulate more work
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(3000)))
	
	// Pause to create the last split
	sw.Pause("Split4")
	fmt.Println("End Split4, Pausing Stopwatch")

	// Print all of the splits
	fmt.Println("\nSplits:")
	fmt.Printf("%v\n\n", sw.Splits)

	// Get the duration of a specific split:
	fmt.Printf("Duration of Split3: [%f]\n\n", sw.Splits["Split3"].ElapsedMS)

	// Calculate a virtual split time from Split1.StartTime to Split3.EndTime
	// Note this does not currently subtract any time that the stopwatch was paused
	elapsed, _ := sw.ElapsedMS("Split1", "Split3")
	fmt.Printf("Virtual Split Time, Split1 to Split3: [%f]\n\n", elapsed)

}
```

### Usage with Prometheus
```go
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
	// You call Split() as a param to Prometheus.Counter.Add()
	fmt.Println("End Split2")
	Split2.Add(sw.Split("Split2"))

	// Sleep some more to simulate more work
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(3000)))

	// Pause the stopwatch and record the next measurement
	// Pause is useful to create gaps between measurements
	fmt.Println("End Split3")
	sw.Pause("Split3")

	// This will be ingored by the stopwatch
	// You can't Split while paused.
	sw.Split("Bad Split")

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
```