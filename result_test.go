package result

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	type args struct {
		value  int
		rError int
		c      C
	}
	tests := []struct {
		name      string
		args      args
		wantValue int
		wantError int
		wantPanic bool
	}{
		{
			name:      "Result with value and empty error",
			args:      args{value: 1, rError: *new(int)},
			wantValue: 1,
			wantError: *new(int),
			wantPanic: false,
		},
		{
			name:      "Result with error and empty value",
			args:      args{value: *new(int), rError: *new(int)},
			wantValue: *new(int),
			wantError: *new(int),
			wantPanic: false,
		},
		{
			name:      "Result with error and value",
			args:      args{value: 1, rError: 1},
			wantValue: 1,
			wantError: 1,
			wantPanic: false,
		},
		{
			name:      "Result with empty error and empty value",
			args:      args{value: *new(int), rError: *new(int)},
			wantValue: *new(int),
			wantError: *new(int),
			wantPanic: false,
		},
		{
			name:      "Result with strict config",
			args:      args{value: 1, rError: 1, c: C{Strict: true}},
			wantValue: 1,
			wantError: 1,
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() { recover() }()

			got := newResult(tt.args.value, tt.args.rError, tt.args.c)

			if tt.wantPanic {
				assert.Panics(t, func() { got.Unwrap() })
			}

			value, err := got.Unwrap()
			assert.Equal(t, tt.wantValue, value)
			assert.Equal(t, tt.wantError, err)
		})
	}
}
