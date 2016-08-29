package main

import (
	"math"
	"sync"
	"time"
)

type stats struct {
	sync.Mutex
	count int
	max   time.Duration
	min   time.Duration

	// online mean & variance algorithm from https://en.wikipedia.org/wiki/Algorithms_for_calculating_variance
	K   time.Duration
	ex  time.Duration
	ex2 time.Duration
}

func (s *stats) record(c *car) {
	s.Lock()
	defer s.Unlock()
	duration := c.done.Sub(c.arrival)

	if duration > s.max {
		s.max = duration
	}
	// min can't really be zero
	if duration < s.min || s.min == 0 {
		s.min = duration
	}

	// online mean & variance algorithm from https://en.wikipedia.org/wiki/Algorithms_for_calculating_variance
	if s.count == 0 {
		s.K = duration
	}
	s.count++
	diff := duration - s.K
	s.ex += diff
	s.ex2 += diff * diff
}

func (s *stats) mean() time.Duration {
	return s.K + s.ex/time.Duration(s.count)
}

func (s *stats) variance() float64 {
	return float64(s.ex2-(s.ex*s.ex/time.Duration(s.count))) / float64(s.count-1)
}

func (s *stats) stdev() time.Duration {
	return time.Duration(math.Sqrt(s.variance()))
}
