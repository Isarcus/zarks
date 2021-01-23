package zmath

import (
	"fmt"
	"math"
)

// Set represents a slice of type float64
type Set []float64

// ToLinear copies by value a [][]float64 to []float64. Any Changes made to the returned Set can then
// be committed to the [][]float64 with the To2D function.
func ToLinear(some2DData [][]float64) Set {
	linearData := make(Set, 0)
	for i := range some2DData {
		linearData = append(linearData, some2DData[i]...)
	}
	return linearData
}

// To2D copies by value a []float64 into a pre-existing [][]float64
func (data Set) To2D(to [][]float64) [][]float64 {
	idx := 0
	for i := 0; i < len(to); i++ {
		for j := 0; j < len(to[i]); j++ {
			to[i][j] = data[idx]
			idx++
		}
	}
	return to
}

// IndicesBetween returns the number of indices between the two values provided. Assumes the Set is pre-sorted.
func (data Set) IndicesBetween(min, max float64) int {
	idxMin := data.IndexOfClosest(min)
	idxMax := data.IndexOfClosest(max)
	return idxMax - idxMin
}

// IndexOf returns the index of the desired item in the Set, or -1 if not found
func (data Set) IndexOf(value float64) int {
	for i, item := range data {
		if item == value {
			return i
		}
	}
	return -1
}

// IndexOfClosest returns the index of the item closest to the desired value. Assumes the Set is pre-sorted.
func (data Set) IndexOfClosest(value float64) int {
	idx := data.IndexOf(value)
	if idx != -1 { // if data contains the exact value (yay)
		return idx
	}

	if value < data.GetMin() {
		return 0
	} else if value > data.GetMax() {
		return len(data) - 1
	}

	for i, item := range data {
		if item > value {
			return i
		}
	}
	return -1
}

// GetMin returns the min of a Set
func (data Set) GetMin() float64 {
	min := data[0]
	for _, item := range data {
		if min > item {
			min = item
		}
	}
	return min
}

// GetMax returns the max of a Set
func (data Set) GetMax() float64 {
	max := data[0]
	for _, item := range data {
		if max < item {
			max = item
		}
	}
	return max
}

// GetRange returns the range of a Set. This is computationally faster than GetMaxOf() - GetMinOf() if you
// only need the range for a given task.
func (data Set) GetRange() float64 {
	min, max := data[0], data[0]
	for _, item := range data {
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
func (data Set) GetMedian() float64 {
	dataCopy := make(Set, 0, len(data))
	copy(dataCopy, data)
	dataCopy.Sort()
	var median float64
	idx := len(data) / 2
	if len(data)%2 == 1 { // if len(data) is odd
		median = data[idx]
	} else {
		median = (data[idx] + data[idx-1]) / 2.0
	}
	return median
}

// GetMean returns the mean of a Set
func (data Set) GetMean() float64 {
	var sum float64
	for _, item := range data {
		sum += item
	}
	return sum / float64(len(data))
}

// GetVariance returns the variance of a Set
func (data Set) GetVariance() float64 {
	var squareSum float64
	mean := data.GetMean()
	for _, item := range data {
		squareSum += math.Pow(mean-item, 2)
	}
	return squareSum / float64(len(data))
}

// GetStd returns the standard deviation of a Set
func (data Set) GetStd() float64 {
	return math.Sqrt(data.GetVariance())
}

// Interpolate linearly adjusts the data to a new min and max
func (data Set) Interpolate(newMin, newMax float64) Set {
	oldMin := data.GetMin()
	oldMax := data.GetMax()
	oldRange := oldMax - oldMin
	newRange := newMax - newMin

	for i := range data {
		data[i] = ((data[i]-oldMin)/oldRange)*newRange + newMin
	}

	return data
}

// Zero zeroes a Set
func (data Set) Zero() {
	for i := 0; i < len(data); i++ {
		data[i] = 0
	}
}

// Sort sorts a Set with a quicksort implementation
func (data Set) Sort() { // this is copied code but it works beautifully
	if len(data) < 2 {
		return
	}

	left, right := 0, len(data)-1
	pivotIndex := len(data) / 2

	data[pivotIndex], data[right] = data[right], data[pivotIndex]

	for i := range data {
		if data[i] < data[right] {
			data[i], data[left] = data[left], data[i]
			left++
		}
	}

	data[left], data[right] = data[right], data[left]

	// Recursively call function
	Set(data[:left]).Sort()
	Set(data[left+1:]).Sort()
}

// Analysis contains basic statistical analysis about a Set
type Analysis struct {
	Size, _  int
	Min      float64
	Max      float64
	Range    float64
	Median   float64
	Mean     float64
	Variance float64
	Std      float64
}

// GetAnalysisOf returns an Analysis struct according to the contents of the Set
func GetAnalysisOf(data Set) *Analysis {
	min, max := data.GetMin(), data.GetMax()
	variance := data.GetVariance()
	return &Analysis{
		Size:     len(data),
		Min:      min,
		Max:      max,
		Range:    max - min,
		Median:   data.GetMedian(),
		Mean:     data.GetMean(),
		Variance: variance,
		Std:      math.Sqrt(variance),
	}
}

// PrintAnalysis outputs the analysis to the terminal
func PrintAnalysis(a *Analysis) {
	fmt.Println("N:        ", a.Size)
	fmt.Println("Min:      ", a.Min)
	fmt.Println("Max:      ", a.Max)
	fmt.Println("Range:    ", a.Range)
	fmt.Println("Median:   ", a.Median)
	fmt.Println("Mean:     ", a.Mean)
	fmt.Println("Variance: ", a.Variance)
	fmt.Println("Std:      ", a.Std)
}

// PrintBasicHistogram prints out what percentages of the dataset are within certain histogram ranges
func (data Set) PrintBasicHistogram() {
	sorted := make(Set, len(data))
	copy(sorted, data)
	sorted.Sort()

	std := data.GetStd()
	mean := data.GetMean()

	for s := -3.; s < 3; s++ {
		min, max := mean+s*std, mean+(s+1)*std
		gap := 100.0 * (float64(sorted.IndicesBetween(min, max)) / float64(len(sorted)))
		fmt.Printf("Between %v and %v stds: %v%% of the data\n", s, s+1, gap)
	}
}
