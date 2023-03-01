package number_sequence_stats

import (
	"fmt"
	"github.com/influxdata/tdigest"
	"github.com/wcharczuk/go-chart"
	"io"
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
	digest *tdigest.TDigest
}

var ErrTDigestDisabled = fmt.Errorf("tdigest is disabled")

func New[T Number](tdigestEnable bool) *Stat[T] {
	s := Stat[T]{}
	if tdigestEnable {
		s.digest = tdigest.NewWithCompression(1000)
	}
	return &s
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
		s.avg = s.sum / float64(s.amount)
	}
	if s.digest != nil {
		s.digest.Add(float64(value), 1)
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
	return math.Sqrt((s.sqsum - s.sum*s.sum/float64(s.amount)) / float64(s.amount))
}

func (s *Stat[T]) Quantile(q float64) float64 {
	if s.digest == nil {
		return 0
	}
	return s.digest.Quantile(q)
}

func (s *Stat[T]) String() string {
	if s.digest == nil {
		return fmt.Sprintf("min: %v, max: %v, avg: %.1f, rms: %.1f, stddev: %1.f sum: %v amount: %v",
			s.Min(), s.Max(), s.Avg(), s.Rms(), s.Stddev(), s.Sum(), s.Amount())
	}
	return fmt.Sprintf("min: %v, max: %v, avg: %.1f, rms: %.1f, stddev: %1.f sum: %v amount: %v\n"+
		"0.1%%: %.1f, 1%%: %.1f, 10%%: %.1f, 50%%: %.1f, 90%%: %.1f, 99%%: %.1f, 99.9%%: %.1f",
		s.Min(), s.Max(), s.Avg(), s.Rms(), s.Stddev(), s.Sum(), s.Amount(),
		s.digest.Quantile(0.001), s.digest.Quantile(0.01), s.digest.Quantile(0.1), s.digest.Quantile(0.5),
		s.digest.Quantile(0.9), s.digest.Quantile(0.99), s.digest.Quantile(0.999))
}

func (s *Stat[T]) DrawPNG(w io.Writer, points int) error {
	if s.digest == nil {
		return ErrTDigestDisabled
	}
	xValues := make([]float64, points)
	yValues := make([]float64, points)
	for i := 0; i < points; i++ {
		xValues[i] = float64(i) / float64(points)
		yValues[i] = s.digest.Quantile(xValues[i])
	}
	graph := chart.Chart{
		Series: []chart.Series{
			chart.ContinuousSeries{
				XValues: xValues,
				YValues: yValues,
			},
		},
	}

	return graph.Render(chart.PNG, w)
}
