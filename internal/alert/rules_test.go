package alert

import "testing"

func TestCrossed(t *testing.T) {
	tests := []struct {
		name     string
		previous float64
		current  float64
		target   float64
		want     bool
	}{
		{"crossed up", 3480, 3510, 3500, true},
		{"crossed down", 3520, 3490, 3500, true},
		{"still above", 3510, 3520, 3500, false},
		{"still below", 3480, 3490, 3500, false},
		{"lands on target from below", 3480, 3500, 3500, true},
		{"leaves target downward", 3500, 3490, 3500, true},
		{"stays exactly on target", 3500, 3500, 3500, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Crossed(tc.previous, tc.current, tc.target)
			if got != tc.want {
				t.Errorf("Crossed(%v, %v, %v) = %v, want %v",
					tc.previous, tc.current, tc.target, got, tc.want)
			}
		})
	}
}
