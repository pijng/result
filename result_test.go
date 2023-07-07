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
	}
	tests := []struct {
		name      string
		args      args
		wantValue int
		wantError error
	}{
		{name: "Result with value and empty error", args: args{value: 1, rError: nil}, wantValue: 1, wantError: nil},
		{name: "Result with error and empty value", args: args{value: *new(int), rError: errors.New("oops, we'got problems")}, wantValue: *new(int), wantError: errors.New("oops, we'got problems")},
		{name: "Result with error and with value", args: args{value: 1, rError: errors.New("oops, we'got problems")}, wantValue: *new(int), wantError: errors.New("oops, we'got problems")},
		{name: "Result with empty error and empty value", args: args{value: *new(int), rError: nil}, wantValue: *new(int), wantError: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.value, tt.args.rError)
			assert.Equal(t, tt.wantValue, got.Value())
			assert.Equal(t, tt.wantError, got.Error())
		})
	}
}
