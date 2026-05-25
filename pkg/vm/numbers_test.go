package vm

import (
	"math"
	"testing"
)

func TestToIntFloatSpecialValues(t *testing.T) {
	tests := []struct {
		name string
		in   Value
		want int
	}{
		{name: "float nan", in: Float(math.NaN()), want: 0},
		{name: "float positive infinity", in: Float(math.Inf(1)), want: int(maxIntValue)},
		{name: "float negative infinity", in: Float(math.Inf(-1)), want: int(minIntValue)},
		{name: "float32 nan", in: Float32(float32(math.NaN())), want: 0},
		{name: "float32 positive infinity", in: Float32(float32(math.Inf(1))), want: int(maxIntValue)},
		{name: "float32 negative infinity", in: Float32(float32(math.Inf(-1))), want: int(minIntValue)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := ToInt(tt.in)
			if !ok {
				t.Fatal("ToInt returned ok=false")
			}
			if got != tt.want {
				t.Fatalf("ToInt(%v) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}
