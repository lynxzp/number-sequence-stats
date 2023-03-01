package number_sequence_stats

import (
	"fmt"
	"math"
)

type Number interface {
	float64 | float32 |
		int64 | int32 | int16 | int8 | int |
		uint64 | uint32 | uint16 | uint8 | uint
}

type Stat[T Number] struct {
	min    T
	max    T
	avg    float64
	amount int64
	sum    float64
	sqsum  float64
}

func New[T Number]() *Stat[T] {
	return &Stat[T]{}
}

func (s *Stat[T]) Add(value T) {
	s.amount++
	if s.amount == 1 {
		s.min = value
		s.max = value
		s.avg = float64(value)
		s.sum = float64(value)
		s.sqsum = float64(value) * float64(value)
	} else {
		if value < s.min {
			s.min = value
		}
		if value > s.max {
			s.max = value
		}
		s.sum += float64(value)
		s.sqsum += float64(value) * float64(value)
		//s.avg = (s.avg*float64(s.amount-1) + float64(value)) / float64(s.amount)
		s.avg = s.sum / float64(s.amount)
	}
}

func (s *Stat[T]) Min() T {
	return s.min
}

func (s *Stat[T]) Max() T {
	return s.max
}

func (s *Stat[T]) Avg() float64 {
	return s.avg
}

func (s *Stat[T]) Amount() int64 {
	return s.amount
}

func (s *Stat[T]) Sum() float64 {
	return s.sum
}

func (s *Stat[T]) Rms() float64 {
	return math.Sqrt(s.sqsum / float64(s.amount))
}

func (s *Stat[T]) Stddev() float64 {
	return math.Sqrt((s.sqsum - 2*s.avg*s.sum + float64(s.amount)*s.avg*s.avg) / float64(s.amount))
}

func (s *Stat[T]) String() string {
	return fmt.Sprintf("min: %v, max: %v, avg: %.1f, rms: %.1f, stddev: %1.f sum: %v amount: %v,",
		s.Min(), s.Max(), s.Avg(), s.Rms(), s.Stddev(), s.Sum(), s.Amount())
}
