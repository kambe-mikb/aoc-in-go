package main

import (
	"bufio"
	"errors"
	"regexp"
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
	reader := strings.NewReader(input)
	scanner := bufio.NewScanner(reader)

	// Optional: Increase buffer size if lines are very long
	// const maxCapacity = 1024 * 1024 // 1MB
	const modulusSize = 100
	// buf := make([]byte, maxCapacity)
	// scanner.Buffer(buf, maxCapacity)

	// Initial position of the dial is 50
	dial, _ := NewModInt(50, modulusSize)

	re := regexp.MustCompile(`(?P<direction>[LR])(?P<distance>[0-9]+)`)

	if part2 {
		// when you're ready to do part 2, remove this "not implemented" block
		return "not implemented"
	} else {
		count := 0
		for scanner.Scan() {
			line := scanner.Text()
			// Process the line here

			if dial.value == 0 {
				count++
			}
			matched := re.FindStringSubmatch(line)
			iDistance, _ := strconv.Atoi(matched[re.SubexpIndex("distance")])
			distance, _ := NewModInt(iDistance, modulusSize)
			if matched[re.SubexpIndex("direction")] == "L" {
				dial, _ = dial.Sub(distance)
			} else {
				dial, _ = dial.Add(distance)
			}
		}

		if err := scanner.Err(); err != nil {
			panic(err)
		}
		return strconv.Itoa(count)
	}
}

// ModInt represents an integer under a given modulus.
type ModInt struct {
	value   int
	modulus int
}

// NewModInt creates a new modular integer.
func NewModInt(value, modulus int) (*ModInt, error) {
	if modulus <= 0 {
		return nil, errors.New("modulus must be positive")
	}
	v := ((value % modulus) + modulus) % modulus // normalize
	return &ModInt{value: v, modulus: modulus}, nil
}

// Value returns the current value.
func (m *ModInt) Value() int {
	return m.value
}

// Add performs modular addition.
func (m *ModInt) Add(other *ModInt) (*ModInt, error) {
	if m.modulus != other.modulus {
		return nil, errors.New("modulus mismatch")
	}
	return NewModInt(m.value+other.value, m.modulus)
}

// Sub performs modular subtraction.
func (m *ModInt) Sub(other *ModInt) (*ModInt, error) {
	if m.modulus != other.modulus {
		return nil, errors.New("modulus mismatch")
	}
	return NewModInt(m.modulus+m.value-other.value, m.modulus)
}
