package stopwatch

import (
	"fmt"
	"sync"
	"time"
)

type Split struct {
	Name      string
	StartTime time.Time
	EndTime   time.Time
	ElapsedMS float64
}

type Stopwatch struct {
	Splits map[string]Split

	mtx          sync.Mutex
	currentSplit Split
	isRunning    bool
}

func New() *Stopwatch {
	s := &Stopwatch{Splits: map[string]Split{}}
	return s
}

func AutoStart() *Stopwatch {
	return &Stopwatch{Splits: map[string]Split{}, currentSplit: Split{StartTime: time.Now()}, isRunning: true}
}

func (s *Stopwatch) Start() {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.Splits = map[string]Split{}
	s.isRunning = true
	s.currentSplit = Split{StartTime: time.Now()}
}

func (s *Stopwatch) Split(splitName string) float64 {

	now := time.Now()

	s.mtx.Lock()
	defer s.mtx.Unlock()

	var elapsed float64

	if s.isRunning {
		s.currentSplit.EndTime = now
		elapsed = ElapsedMS(s.currentSplit.StartTime, s.currentSplit.EndTime)
		s.currentSplit.ElapsedMS = elapsed
		s.currentSplit.Name = splitName
		s.Splits[splitName] = s.currentSplit
		s.currentSplit = Split{StartTime: time.Now()}
	}
	return float64(elapsed)
}

func (s *Stopwatch) Pause(splitName string) float64 {

	now := time.Now()

	s.mtx.Lock()
	defer s.mtx.Unlock()

	var elapsed float64

	if s.isRunning {
		s.currentSplit.EndTime = now
		elapsed = ElapsedMS(s.currentSplit.StartTime, s.currentSplit.EndTime)
		s.currentSplit.ElapsedMS = elapsed
		s.currentSplit.Name = splitName
		s.Splits[splitName] = s.currentSplit
		s.currentSplit = Split{}
		s.isRunning = false
	}
	return float64(elapsed)
}

func (s *Stopwatch) Resume() {

	now := time.Now()

	s.mtx.Lock()
	defer s.mtx.Unlock()

	if !s.isRunning {
		s.currentSplit = Split{StartTime: now}
		s.isRunning = true
	}
}

func (s *Stopwatch) Elapsed(startSplitName string, endSplitName string) (time.Duration, error) {

	var n time.Duration

	if _, ok := s.Splits[startSplitName]; !ok {
		return n, fmt.Errorf("Invalid startSplitName [%s]", startSplitName)
	}

	if _, ok := s.Splits[endSplitName]; !ok {
		return n, fmt.Errorf("Invalid endSplitName [%s]", endSplitName)
	}

	return Elapsed(s.Splits[startSplitName].StartTime, s.Splits[endSplitName].EndTime), nil
}

func (s *Stopwatch) ElapsedMS(startSplitName string, endSplitName string) (float64, error) {

	dur, err := s.Elapsed(startSplitName, endSplitName)
	if err != nil {
		return float64(0), err
	}
	return DurationToMS(dur), nil
}

func Elapsed(start time.Time, end time.Time) time.Duration {
	return end.Sub(start)
}

func ElapsedMS(start time.Time, end time.Time) float64 {
	return DurationToMS(Elapsed(start, end))
}

func DurationToMS(dur time.Duration) float64 {
	return float64(dur.Nanoseconds() / 1000000)
}
