package utils

import "testing"

func TestMustExist(t *testing.T) {
	type args struct {
		dirpath []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test",
			args: args{dirpath: []string{"./cache"}},
			want: "./cache",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MustExist(tt.args.dirpath...); got != tt.want {
				t.Errorf("MustExist() = %v, want %v", got, tt.want)
			}
		})
	}
}
