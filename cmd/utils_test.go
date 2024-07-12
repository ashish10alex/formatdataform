package cmd

import (
	"strings"
	"testing"
)

func TestLineCounterV3(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		{"Normal file 3 lines", "Line 1\nLine 2\nLine 3\n", 3, false},
		{"Empty file", "", 0, false},
		{"No newline at end", "Line 1\nLine 2\nLine 3", 2, false},
		{"Single line", "Single line", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			got, err := lineCounterV3(reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("lineCounterV3() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("lineCounterV3() = %v, want %v", got, tt.want)
			}
		})
	}
}
