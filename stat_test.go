package number_sequence_stats

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

func rmsStddev[T Number](v []T) (float64, float64) {
	var sum, sqsum float64
	for _, i := range v {
		sum += float64(i)
		sqsum += float64(i) * float64(i)
	}
	avg := sum / float64(len(v))
	rms := math.Sqrt(sqsum / float64(len(v)))

	var stddev float64
	for _, i := range v {
		stddev += (float64(i) - avg) * (float64(i) - avg)
	}
	stddev = math.Sqrt(stddev / float64(len(v)))
	return rms, stddev
}

func floatCompare(a, b float64) bool {
	return math.Abs(a-b) < 0.000000001*math.Max(math.Abs(a), math.Abs(b))
}

func TestStat(t *testing.T) {

	type test struct {
		input  []int
		min    int
		max    int
		avg    float64
		amount int64
		rms    float64
		stddev float64
	}

	tests := make([]test, 0)
	tests = append(tests, test{input: []int{1, 2, 3}, min: 1, max: 3, avg: 2, amount: 3, rms: 2.160246899469287, stddev: 0.816496580927726})
	tests = append(tests, test{input: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, min: 0, max: 9, avg: 4.5, amount: 10, rms: 5.338539126, stddev: 2.87228132332690143})

	n := rand.Intn(10000000)
	var randomizedT test
	randomizedT.input = make([]int, n)
	min := math.MaxInt
	max := 0
	var sum, sqsum float64
	for i := 0; i < n; i++ {
		randomizedT.input[i] = rand.Int()
		if randomizedT.input[i] < min {
			min = randomizedT.input[i]
		}
		if randomizedT.input[i] > max {
			max = randomizedT.input[i]
		}
		sum += float64(randomizedT.input[i])
		sqsum += float64(randomizedT.input[i]) * float64(randomizedT.input[i])
	}
	randomizedT.min = min
	randomizedT.max = max
	randomizedT.amount = int64(n)
	randomizedT.avg = sum / float64(n)
	randomizedT.rms, randomizedT.stddev = rmsStddev(randomizedT.input)
	tests = append(tests, randomizedT)

	for _, tt := range tests {
		s := New[int](false)
		for _, i := range tt.input {
			s.Add(i)
		}
		if s.Min() != tt.min {
			t.Errorf("Min() = %d, want %d", s.Min(), tt.min)
		}
		if s.Max() != tt.max {
			t.Errorf("Max() = %d, want %d", s.Max(), tt.max)
		}
		if floatCompare(s.Avg(), tt.avg) == false {
			t.Errorf("Avg() = %f, want %f", s.Avg(), tt.avg)
		}
		if s.Amount() != tt.amount {
			t.Errorf("Amount() = %d, want %d", s.Amount(), tt.amount)
		}
		if floatCompare(s.Rms(), tt.rms) == false {
			t.Errorf("Rms() = %f, want %f", s.Rms(), tt.rms)
		}
		if floatCompare(s.Stddev(), tt.stddev) == false {
			t.Errorf("Stddev() = %f, want %f", s.Stddev(), tt.stddev)
		}
	}
}

var table = []struct {
	input int
}{
	{input: 1000},
	{input: 1000 * 1000},
	{input: 1000 * 1000 * 1000},
}

func BenchmarkStat(b *testing.B) {
	for _, v := range table {
		b.Run(fmt.Sprintf("input_size_%d", v.input), func(b *testing.B) {

			for i := 0; i < b.N; i++ {
				s := New[int](false)
				for j := 0; j < v.input; j++ {
					s.Add(j)
				}
			}
		})
	}
}
