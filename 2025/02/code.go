package main

import (
	"math"
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
		return "not implemented"
	}
	// solve part 1 here
	for _, idRange := range idRanges {
		bounds := strings.Split(idRange, "-")
		// Narrow the search space if we can
		lenFirst := len(bounds[0])
		lenLast := len(bounds[1])
		// The range is the same odd number of digits at each end, so we can skip it
		repeat := 2
		if (lenFirst == lenLast) && ((lenFirst % repeat) > 0) {
			continue
		}
		first, _ := strconv.Atoi(bounds[0])
		last, _ := strconv.Atoi(bounds[1])
		// Narrow the lower bound
		if (lenFirst % repeat) > 0 { // first contains an odd number of digits
			j := int(math.Pow(10.0, float64(lenFirst)))
			if j <= last {
				first = j
			} else {
				continue
			}
		}
		// Narrow the upper bound
		if (lenLast % 2) > 0 { // last contains an odd number of digits
			j := int(math.Pow(10.0, float64(lenLast-1)) - 1)
			if j >= first {
				last = j
			} else {
				continue
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
		scale := int(math.Pow(10.0, float64(lenSub))) + 1
		// Compute all the possible "invalid" product ids in our final range
		for i := min; i <= max; i++ {
			checkInt := i * scale
			if checkInt >= first {
				if checkInt <= last {
					// Goldilocks!
					total += checkInt
				} else {
					// we are out of the range, we can stop
					break
				}
			} else {
				// keep going until we are in the range (or we run out of numbers)
				continue
			}
		}
	}
	return strconv.Itoa(total)
}
