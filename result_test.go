package result

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	type args struct {
		value  int
		rError error
	}
	tests := []struct {
		name      string
		args      args
		wantValue int
		wantError error
		wantPanic bool
	}{
		{
			name:      "Result with value and empty error",
			args:      args{value: 1, rError: nil},
			wantValue: 1,
			wantError: nil,
			wantPanic: false,
		},
		{
			name:      "Result with error and empty value",
			args:      args{value: *new(int), rError: fmt.Errorf("error")},
			wantValue: *new(int),
			wantError: fmt.Errorf("error"),
			wantPanic: false,
		},
		{
			name:      "Result with error and value",
			args:      args{value: 1, rError: fmt.Errorf("error")},
			wantValue: 1,
			wantError: fmt.Errorf("error"),
			wantPanic: false,
		},
		{
			name:      "Result with empty error and empty value",
			args:      args{value: *new(int), rError: fmt.Errorf("error")},
			wantValue: *new(int),
			wantError: fmt.Errorf("error"),
			wantPanic: false,
		},
		{
			name:      "Result with strict config",
			args:      args{value: 1, rError: fmt.Errorf("error")},
			wantValue: 1,
			wantError: fmt.Errorf("error"),
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newResult(tt.args.value, tt.args.rError)

			value, err := got.Unwrap()
			assert.Equal(t, tt.wantValue, value)
			assert.Equal(t, tt.wantError, err)
		})
	}
}
