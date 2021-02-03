package znshell

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"zarks/system"
	"zarks/zimg"
	"zarks/zmath"
	"zarks/zmath/noise"
)

type saveType int

const (
	timeDate saveType = iota
	custom
)

const savePrefix = "save/"

// RunShell runs a command-line shell for generating noise images
func RunShell() {
	if _, err := os.Stat(savePrefix); os.IsNotExist(err) {
		os.Mkdir(savePrefix, os.ModeDir)
	}

	var (
		r           = system.NewConsoleReader()
		initialized = false
		cfg         = noise.Config{}
		savePref    = timeDate
	)
	fmt.Println("* * * * * * * * * * * * * * *")
	fmt.Println("* * * ZARKS NOISE SHELL * * *")

	for {
		fmt.Println("* * * * * * * * * * * * * * *")
		c()
		fmt.Println("Please answer the following questions about the image you would like to generate.")
		if initialized && system.QueryYN(r, "Use the same settings as last time?") {
			m := noise.NewSimplexMap(cfg)
			img := zimg.Colorify(m, zimg.SchemeGrayscale)
			system.SaveImage(savePrefix+getTitle(savePref, r)+".png", img)
		} else {
			initialized = true
			cfg = getConfig(r)

			m := noise.NewSimplexMap(cfg)
			img := zimg.Colorify(m, zimg.SchemeGrayscale)
			system.SaveImage(savePrefix+getTitle(savePref, r)+".png", img)
		}

		c()
		if !system.QueryYN(r, "Another one?") {
			break
		}
	}

	fmt.Println("* * * SHELL TERMINATED * * *")
}

func getTitle(t saveType, r *bufio.Reader) string {
	switch t {
	case timeDate:
		return system.GetTimeAsString(time.Now())
	case custom:
		return system.Query(r, "Please enter the preferred title without extensions like .png or .jpg:")
	default:
		return "SaveTitleError"
	}
}

func getConfig(r *bufio.Reader) noise.Config {
	var (
		done    = false
		bounds  zmath.VecInt
		octaves int
		rad     float64
		bsi     float64
	)

	// Dimensions
	for !done {
		c()
		fmt.Println("Please enter the desired dimensions of the image in the form of WIDTHxHEIGHT.")
		wh := system.Query(r, "For example, for an 800 by 500 image, enter '800x500'")
		bounds, done = parseWH(wh)
	}

	done = false
	for !done {
		c()
		fmt.Println("How many octaves would you like to generate?")
		octStr := system.Query(r, "(Enter an integer greater than 0; beyond 8 doesn't make much difference)")
		val, err := strconv.Atoi(octStr)
		if err != nil {
			fmt.Println(err)
		} else {
			done = true
			octaves = zmath.MaxInt(1, val)
		}
	}

	done = false
	for !done {
		c()
		fmt.Println("How big would you like the R-size to be?")
		fmt.Println("This parameter affects the area of influence of each vector in 2D space.")
		rStr := system.Query(r, "Typical values range between 0.5 and 0.7; weird things may happen outside that range!")
		val, err := strconv.ParseFloat(rStr, 64)
		if err != nil {
			fmt.Println(err)
		} else {
			done = true
			rad = val
		}
	}

	done = false
	for !done {
		c()
		fmt.Println("About how many pixels should be allocated per simplex (vector) in the first octave?")
		bsiStr := system.Query(r, "A good number to pick is usually about 0.5x the smaller X/Y dimension size.")
		val, err := strconv.ParseFloat(bsiStr, 64)
		if err != nil {
			fmt.Println(err)
		} else {
			done = true
			bsi = val
		}
	}

	return noise.Config{
		Dimensions:     bounds,
		Normalize:      true,
		Octaves:        octaves,
		R:              rad,
		BoxSizeInitial: bsi,
	}
}

func parseWH(ipt string) (zmath.VecInt, bool) {
	// check for x and X
	idx := strings.Index(ipt, "x")
	idx = zmath.MaxInt(idx, strings.Index(ipt, "X"))
	if idx < 1 || (len(ipt)-idx) < 2 {
		return zmath.ZVI, false
	}

	width, err1 := strconv.Atoi(ipt[:idx])
	if err1 != nil {
		fmt.Println(err1)
		return zmath.ZVI, false
	}
	height, err2 := strconv.Atoi(ipt[idx+1:])
	if err2 != nil {
		fmt.Println(err2)
		return zmath.ZVI, false
	}

	vi := zmath.VI(width, height)
	if vi.Min() < 1 {
		fmt.Println("Dimensions must be > 0 !!!")
		return zmath.ZVI, false
	}

	return vi, true
}

func c(n ...int) {
	if len(n) == 0 {
		fmt.Println()
	} else {
		for i := 0; i < n[0]; i++ {
			fmt.Println()
		}
	}
}

func main() {
	RunShell()
}
