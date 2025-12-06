package main

import (
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/jpillora/puzzler/harness/aoc"
)

func main() {
	aoc.Harness(run)
}

// on code change, run will be executed 4 times:
// 1. with: false (part1), and example input
// 2. with: true (part2), and example input
// 3. with: false (part1), and user input
// 4. with: true (part2), and user input
// the return value of each run is printed to stdout
func run(part2 bool, input string) any {
	total := 0
	// The input contains a trailing newline, which breaks our last atoi() call
	idRanges := strings.Split(strings.TrimSpace(input), ",")
	// when you're ready to do part 2, remove this "not implemented" block
	if part2 {
		for _, idRange := range idRanges {
			bounds := strings.Split(idRange, "-")
			lenFirst := len(bounds[0])
			lenLast := len(bounds[1])
			cache := []int{}
			for repeat := 2; repeat <= lenLast; repeat++ {
				var result int
				// We need to return the cache, since the cache header is passed by value,
				// and when the cache size changes, the header needs to change in this scope as well
				result, cache = findInvalids(repeat, bounds[0], lenFirst, bounds[1], lenLast, cache)
				// fmt.Println(" ", idRange, repeat, result)
				total += result
			}
		}
		return strconv.Itoa(total)
	}
	// solve part 1 here
	for _, idRange := range idRanges {
		bounds := strings.Split(idRange, "-")
		lenFirst := len(bounds[0])
		lenLast := len(bounds[1])
		cache := []int{}
		var result int
		result, _ = findInvalids(2, bounds[0], lenFirst, bounds[1], lenLast, cache)
		total += result
	}
	return strconv.Itoa(total)
}

func findInvalids(repeat int, lowerBound string, lenFirst int, upperBound string, lenLast int, cache []int) (int, []int) {
	total := 0
	// Narrow the search space if we can
	// The range is the same number of digits at each end, and can't be evenly divided by the pattern size we're using, so we can skip it
	if (lenFirst == lenLast) && ((lenFirst % repeat) > 0) {
		return 0, cache
	}
	first, _ := strconv.Atoi(lowerBound)
	last, _ := strconv.Atoi(upperBound)
	// Narrow the lower bound
	if (lenFirst % repeat) > 0 { // number of digits in first can't be evenly divided by the pattern size
		j := int(math.Pow(10.0, float64(lenFirst)))
		if j <= last {
			first = j
		} else {
			return 0, cache
		}
	}
	// Narrow the upper bound
	if (lenLast % repeat) > 0 { // last contains an odd number of digits
		j := int(math.Pow(10.0, float64(lenLast-1)) - 1)
		if j >= first {
			last = j
		} else {
			return 0, cache
		}
	}
	firstString := strconv.Itoa(first)
	lastString := strconv.Itoa(last)
	// this is the length of our substring
	// (note the implicit assumption that firstString and lastString are the same length!)
	lenSub := len(firstString) / repeat
	// create an approximate range of substrings
	min, _ := strconv.Atoi(firstString[:lenSub])
	max, _ := strconv.Atoi(lastString[:lenSub])
	// the scale factor converts our substring into the full string with repeated strings of digits
	scale := 1
	for f := 1; f < repeat; f++ {
		scale += int(math.Pow(10.0, float64(lenSub*f)))
	}
	// scale := int(math.Pow(10.0, float64(lenSub))) + 1
	// Compute all the possible "invalid" product ids in our final range
	for i := min; i <= max; i++ {
		checkInt := i * scale
		if checkInt >= first {
			if checkInt <= last {
				// Goldilocks!
				if !slices.Contains(cache, checkInt) {
					total += checkInt
					cache = append(cache, checkInt)
				}
			} else {
				// we are out of the range, we can stop
				break
			}
		} else {
			// keep going until we are in the range (or we run out of numbers)
			continue
		}
	}
	return total, cache
}
