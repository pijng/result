package result

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	type args struct {
		value  int
		rError error
		c      C
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
			args:      args{value: *new(int), rError: errors.New("oops, we'got problems")},
			wantValue: *new(int),
			wantError: errors.New("oops, we'got problems"),
			wantPanic: false,
		},
		{
			name:      "Result with error and value",
			args:      args{value: 1, rError: errors.New("oops, we'got problems")},
			wantValue: *new(int),
			wantError: errors.New("oops, we'got problems"),
			wantPanic: false,
		},
		{
			name:      "Result with empty error and empty value",
			args:      args{value: *new(int), rError: nil},
			wantValue: *new(int),
			wantError: nil,
			wantPanic: false,
		},
		{
			name:      "Result with strict config",
			args:      args{value: 1, rError: nil, c: C{strict: true}},
			wantValue: 1,
			wantError: nil,
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() { recover() }()

			got := New(tt.args.value, tt.args.rError, tt.args.c)

			if tt.wantPanic {
				assert.Panics(t, func() { got.Value() })
			}

			assert.Equal(t, tt.wantValue, got.Value())
			assert.Equal(t, tt.wantError, got.Error())
		})
	}
}
