package logger

import "testing"

func TestDebug(t *testing.T) {
	type args struct {
		msg []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "debug",
			args: args{msg: nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Debug(tt.args.msg...)
		})
	}
}

func TestDebugf(t *testing.T) {
	type args struct {
		format string
		args   []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "debug",
			args: args{
				format: "",
				args:   nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Debugf(tt.args.format, tt.args.args...)
		})
	}
}

func TestError(t *testing.T) {
	type args struct {
		msg []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "debug",
			args: args{msg: nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Error(tt.args.msg...)
		})
	}
}

func TestErrorf(t *testing.T) {
	type args struct {
		format string
		args   []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "debug",
			args: args{
				format: "",
				args:   nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Errorf(tt.args.format, tt.args.args...)
		})
	}
}

func TestInfo(t *testing.T) {
	type args struct {
		msg []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "debug",
			args: args{
				msg: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Info(tt.args.msg...)
		})
	}
}

func TestInfof(t *testing.T) {
	type args struct {
		format string
		args   []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "debug",
			args: args{
				format: "",
				args:   nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Infof(tt.args.format, tt.args.args...)
		})
	}
}

func TestWarn(t *testing.T) {
	type args struct {
		msg []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "debug",
			args: args{
				msg: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Warn(tt.args.msg...)
		})
	}
}

func TestWarnf(t *testing.T) {
	type args struct {
		format string
		args   []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "debug",
			args: args{
				format: "",
				args:   nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Warnf(tt.args.format, tt.args.args...)
		})
	}
}
