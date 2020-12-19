package errorx

import (
	"errors"
	"reflect"
	"testing"
)

func TestMatch(t *testing.T) {
	cases := []struct {
		err1        error
		err2        error
		expectMatch bool
	}{
		{
			err1:        E(errors.New("error")),
			err2:        nil,
			expectMatch: false,
		},
		{
			err1:        E(errors.New("error")),
			err2:        errors.New("error"),
			expectMatch: true,
		},
		{
			err1:        E(errors.New("error")),
			err2:        errors.New("something is different"),
			expectMatch: false,
		},
	}

	for _, val := range cases {
		if match := Match(val.err1, val.err2); match != val.expectMatch {
			t.Errorf("Match() want = %vv, got %v", val.expectMatch, match)
		}
	}
}

func TestE(t *testing.T) {
	tests := []struct {
		name string
		args []interface{}
	}{
		{
			name: "No args",
			args: nil,
		},
		{
			name: "1 layer",
			args: []interface{}{
				Message("message"),
				CodeConflict,
				MetricSuccess,
				Fields{
					"K": "V",
				},
				Op("userService.CreateUser"),
			},
		},
		{
			name: "2 layer with standard error",
			args: []interface{}{
				errors.New("standard-error"),
				Message("message"),
				CodeConflict,
				Op("userService.CreateUser"),
			},
		},
		{
			name: "2 layer",
			args: []interface{}{
				&Error{
					Code:     CodeGateway,
					Message:  "gateway-message",
					Op:       "userGateway.FindUser",
					OpTraces: []Op{"userGateway.FindUser"},
					Err:      errors.New("standard-error"),
				},
				Message("message"),
				CodeConflict,
				Op("userService.CreateUser"),
			},
		},
		{
			name: "Invalid type",
			args: []interface{}{
				123,
			},
		},
		{
			name: "Same code",
			args: []interface{}{
				&Error{
					Code:     CodeConflict,
					Message:  "gateway-message",
					Op:       "userGateway.FindUser",
					OpTraces: []Op{"userGateway.FindUser"},
					Err:      errors.New("standard-error"),
				},
				Message("message"),
				CodeConflict,
				Op("userService.CreateUser"),
			},
		},
		{
			name: "Missing code",
			args: []interface{}{
				&Error{
					Code:    CodeConflict,
					Message: "gateway-message",
					Op:      "userGateway.FindUser",
					Err:     errors.New("standard-error"),
				},
				Message("message"),
				Op("userService.CreateUser"),
			},
		},
		{
			name: "New from string",
			args: []interface{}{
				"this is an error",
				Message("message"),
				Op("userService.CreateUser"),
			},
		},
	}

	enableLog := false
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := E(tt.args...)
			if enableLog {
				t.Logf("got = \n%#v\n\n", got)
			}
		})
	}
}

func TestIs(t *testing.T) {
	type args struct {
		code Code
		err  error
	}
	tests := []struct {
		name string
		args
		want bool
	}{
		{
			name: "Error is nil",
			args: args{
				code: "",
				err:  nil,
			},
			want: false,
		},
		{
			name: "Error is standard error",
			args: args{
				code: "",
				err:  errors.New("standard-error"),
			},
			want: false,
		},
		{
			name: "1 layer",
			args: args{
				code: CodeGateway,
				err: &Error{
					Code:    CodeGateway,
					Message: "",
					Op:      "",
					Err:     nil,
				},
			},
			want: true,
		},
		{
			name: "2 layer",
			args: args{
				code: CodeGateway,
				err: &Error{
					Code:    CodeUnknown,
					Message: "",
					Op:      "",
					Err: &Error{
						Code:    CodeGateway,
						Message: "",
						Op:      "",
						Err:     nil,
					},
				},
			},
			want: true,
		},
		{
			name: "2 layer Unknown",
			args: args{
				code: CodeGateway,
				err: &Error{
					Code:    CodeUnknown,
					Message: "",
					Op:      "",
					Err: &Error{
						Code:    CodeUnknown,
						Message: "",
						Op:      "",
						Err:     nil,
					},
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Is(tt.args.code, tt.args.err)
			if !reflect.DeepEqual(tt.want, got) {
				msg := "\nwant = %#v" + "\ngot  = %#v"
				t.Errorf(msg, tt.want, got)
			}
		})
	}
}
