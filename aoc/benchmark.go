package aoc

import "testing"

func BenchmarkSolutions[T any, In interface{ ~[]byte | ~string }](b *testing.B, year, day int, part1, part2 func(In) T) {
	b.Helper()

	input, err := GetAndCacheInput(b.Context(), nil, NewRequest(year, day).BuildGetInputRequest())
	if err != nil {
		b.Fatalf("error getting input: %v", err)
	}

	typedInput := In(input)

	b.Run("part 1", func(b *testing.B) {
		if part1 == nil {
			b.Skip("no solution provided for part 2")
		}
		for b.Loop() {
			part1(typedInput)
		}
	})

	b.Run("part 2", func(b *testing.B) {
		if part2 == nil {
			b.Skip("no solution provided for part 2")
		}
		for b.Loop() {
			part2(typedInput)
		}
	})
}
