package dates

import (
	"reflect"
	"testing"
)

func Test_getParts(t *testing.T) {
	type args struct {
		dateStr string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "false dates",
			args: args{"0/0/0000"},
			want: []string{"0", "0", "0000"},
		},
		{
			name: "simple dates",
			args: args{"1/1/1990"},
			want: []string{"1", "1", "1990"},
		},
		{
			name: "full dates",
			args: args{"01/01/1990"},
			want: []string{"01", "01", "1990"},
		},
		{
			name: "bad day",
			args: args{"633/1/1990"},
			want: nil,
		},
		{
			name: "bad day with letters",
			args: args{"-33/1/1990"},
			want: nil,
		},
		{
			name: "bad month",
			args: args{"11/123/1990"},
			want: nil,
		},
		{
			name: "bad year",
			args: args{"01/01/190"},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getParts(tt.args.dateStr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getParts() = %v, want %v", got, tt.want)
			}
		})
	}
}

// func Test_validate2(t *testing.T) {
// 	type args struct {
// 		parts []string
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want []int
// 	}{
// 		{
// 			name: "validate dates",
// 			args: args{[]string{"01", "01", "1900"}},
// 			want: []int{1, 1, 1900},
// 		},
// 		{
// 			name: "too old date",
// 			args: args{[]string{"01", "01", "1899"}},
// 			want: nil,
// 		},
// 		{
// 			name: "invalidate day",
// 			args: args{[]string{"00", "01", "1900"}},
// 			want: nil,
// 		},
// 		{
// 			name: "invalidate day",
// 			args: args{[]string{"32", "01", "1900"}},
// 			want: nil,
// 		},
// 		{
// 			name: "invalidate month",
// 			args: args{[]string{"01", "00", "1900"}},
// 			want: nil,
// 		},
// 		{
// 			name: "invalidate month",
// 			args: args{[]string{"01", "13", "1900"}},
// 			want: nil,
// 		},
// 		{
// 			name: "too old date",
// 			args: args{[]string{"01", "01", "1899"}},
// 			want: nil,
// 		},
// 		{
// 			name: "too advanced date",
// 			args: args{[]string{"01", "01", "3000"}},
// 			want: nil,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := validate(tt.args.parts); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("validate() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

func TestValidate(t *testing.T) {
	type args struct {
		parts []string
	}
	tests := []struct {
		name    string
		args    args
		want    []int
		wantErr bool
	}{
		{
			name:    "validate earliest day",
			args:    args{[]string{"01", "01", "1900"}},
			want:    []int{1, 1, 1900},
			wantErr: false,
		},
		{
			name:    "validate latest day",
			args:    args{[]string{"31", "12", "2999"}},
			want:    []int{31, 12, 2999},
			wantErr: false,
		},
		{
			name:    "invalidate day: 00",
			args:    args{[]string{"00", "01", "1900"}},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalidate day: 32",
			args:    args{[]string{"32", "01", "1900"}},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalidate month: 00",
			args:    args{[]string{"01", "00", "1900"}},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalidate month: 13",
			args:    args{[]string{"01", "13", "1900"}},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "too old date",
			args:    args{[]string{"01", "01", "1899"}},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "too advanced date",
			args:    args{[]string{"01", "01", "3000"}},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validate(tt.args.parts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
