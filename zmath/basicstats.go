package zmath

import (
	"fmt"
	"math"
)

// Set represents a slice of type float64
type Set []float64

// ToLinear copies by value a [][]float64 to []float64. Any Changes made to the returned Set can then
// be committed to the [][]float64 with the To2D function.
func ToLinear(data Map) Set {
	linearData := make(Set, 0, int(data.Area()))
	for i := range data {
		linearData = append(linearData, data[i]...)
	}
	return linearData
}

// To2D copies by value a []float64 into a pre-existing [][]float64
func (s Set) To2D(to [][]float64) [][]float64 {
	idx := 0
	for i := 0; i < len(to); i++ {
		for j := 0; j < len(to[i]); j++ {
			to[i][j] = s[idx]
			idx++
		}
	}
	return to
}

// IndicesBetween returns the number of indices between the two values provided. Assumes the Set is pre-sorted.
func (s Set) IndicesBetween(min, max float64) int {
	idxMin := s.IndexOfClosest(min)
	idxMax := s.IndexOfClosest(max)
	return idxMax - idxMin
}

// IndexOf returns the lowest index of the desired item in the Set, or -1 if not found
func (s Set) IndexOf(value float64) int {
	for i, item := range s {
		if item == value {
			return i
		}
	}
	return -1
}

// IndexOfClosest returns the index of the item closest to the desired value. Assumes the Set is pre-sorted.
func (s Set) IndexOfClosest(value float64) int {
	idx := s.IndexOf(value)
	if idx != -1 { // if data contains the exact value (yay)
		return idx
	}

	if value < s.GetMin() {
		return 0
	} else if value > s.GetMax() {
		return len(s) - 1
	}

	for i, item := range s {
		if item > value {
			return i
		}
	}
	return -1
}

// GetMin returns the min of a Set
func (s Set) GetMin() float64 {
	min := s[0]
	for _, item := range s {
		if min > item {
			min = item
		}
	}
	return min
}

// GetMax returns the max of a Set
func (s Set) GetMax() float64 {
	max := s[0]
	for _, item := range s {
		if max < item {
			max = item
		}
	}
	return max
}

// GetRange returns the range of a Set. This is computationally faster than GetMaxOf() - GetMinOf() if you
// only need the range for a given task.
func (s Set) GetRange() float64 {
	min, max := s[0], s[0]
	for _, item := range s {
		if min > item {
			min = item
		}
		if max < item {
			max = item
		}
	}
	return max - min
}

// GetMedian returns the median of a Set
func (s Set) GetMedian() float64 {
	dataCopy := make(Set, 0, len(s))
	copy(dataCopy, s)
	dataCopy.Sort()
	var median float64
	idx := len(s) / 2
	if len(s)%2 == 1 { // if len(data) is odd
		median = s[idx]
	} else {
		median = (s[idx] + s[idx-1]) / 2.0
	}
	return median
}

// GetMean returns the mean of a Set
func (s Set) GetMean() float64 {
	var sum float64
	for _, item := range s {
		sum += item
	}
	return sum / float64(len(s))
}

// GetVariance returns the variance of a Set
func (s Set) GetVariance() float64 {
	var squareSum float64
	mean := s.GetMean()
	for _, item := range s {
		squareSum += math.Pow(mean-item, 2)
	}
	return squareSum / float64(len(s))
}

// GetStd returns the standard deviation of a Set
func (s Set) GetStd() float64 {
	return math.Sqrt(s.GetVariance())
}

// Zero zeroes a Set
func (s Set) Zero() Set {
	for i := 0; i < len(s); i++ {
		s[i] = 0
	}
	return s
}

// Copy deepcopies the called Set into a new Set and returns the new one
func (s Set) Copy() Set {
	newSet := make(Set, len(s))
	copy(newSet, s)
	return newSet
}

// Interpolate linearly adjusts the data to a new min and max
func (s Set) Interpolate(newMin, newMax float64) Set {
	oldMin := s.GetMin()
	oldMax := s.GetMax()
	oldRange := oldMax - oldMin
	newRange := newMax - newMin

	for i := range s {
		s[i] = ((s[i]-oldMin)/oldRange)*newRange + newMin
	}

	return s
}

// MakeUniform makes the called Set follow a uniform distribution
func (s Set) MakeUniform() Set {
	var (
		length = float64(len(s)) - 1
		sorted = s.Copy().Sort()
		retSet = make(Set, len(s))
	)

	for i, val := range s {
		retSet[i] = float64(sorted.IndexOf(val)) / length
	}

	copy(s, retSet)
	return s
}

// Sort sorts a Set with a quicksort implementation
func (s Set) Sort() Set { // this is copied code but it works beautifully
	if len(s) < 2 {
		return s
	}

	left, right := 0, len(s)-1
	pivotIndex := len(s) / 2

	s[pivotIndex], s[right] = s[right], s[pivotIndex]

	for i := range s {
		if s[i] < s[right] {
			s[i], s[left] = s[left], s[i]
			left++
		}
	}

	s[left], s[right] = s[right], s[left]

	// Recursively call function
	Set(s[:left]).Sort()
	Set(s[left+1:]).Sort()

	return s
}

// Stats contains basic statistical analysis about a Set
type Stats struct {
	Size, _  int
	Min      float64
	Max      float64
	Range    float64
	Median   float64
	Mean     float64
	Variance float64
	Std      float64
}

// GetAnalysis returns an Analysis struct according to the contents of the Set
func (s Set) GetAnalysis() *Stats {
	min, max := s.GetMin(), s.GetMax()
	variance := s.GetVariance()
	return &Stats{
		Size:     len(s),
		Min:      min,
		Max:      max,
		Range:    max - min,
		Median:   s.GetMedian(),
		Mean:     s.GetMean(),
		Variance: variance,
		Std:      math.Sqrt(variance),
	}
}

// PrintAnalysis outputs the analysis to the terminal
func PrintAnalysis(a *Stats) {
	fmt.Println("N:        ", a.Size)
	fmt.Println("Min:      ", a.Min)
	fmt.Println("Max:      ", a.Max)
	fmt.Println("Range:    ", a.Range)
	fmt.Println("Median:   ", a.Median)
	fmt.Println("Mean:     ", a.Mean)
	fmt.Println("Variance: ", a.Variance)
	fmt.Println("Std:      ", a.Std)
}

// PrintBasicHistogram prints out what percentages of the dataset are within certain histogram ranges.
// Use on a noise map (simplex, perlin, etc.) for some interesting insights!
func (s Set) PrintBasicHistogram() {
	sorted := make(Set, len(s))
	copy(sorted, s)
	sorted.Sort()

	std := s.GetStd()
	mean := s.GetMean()

	for s := -3.; s < 3; s++ {
		min, max := mean+s*std, mean+(s+1)*std
		gap := 100.0 * (float64(sorted.IndicesBetween(min, max)) / float64(len(sorted)))
		fmt.Printf("Between %v and %v stds: %v%% of the data\n", s, s+1, gap)
	}
}
