package dates

import "testing"

func Test_LeapYear(t *testing.T) {
	type args struct {
		y int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "1900",
			args: args{1900},
			want: false,
		},
		{
			name: "1992",
			args: args{1992},
			want: true,
		},
		{
			name: "2000",
			args: args{2000},
			want: true,
		},
		{
			name: "2400",
			args: args{2400},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LeapYear(tt.args.y); got != tt.want {
				t.Errorf("leapYear() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_DaysSoFar(t *testing.T) {
	type args struct {
		d int
		m int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1/1",
			args: args{1, 1},
			want: 1,
		},
		{
			name: "1/2",
			args: args{1, 2},
			want: 32,
		},
		{
			name: "1/3",
			args: args{1, 3},
			want: 60,
		},
		{
			name: "31/12",
			args: args{31, 12},
			want: 365,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DaysSoFar(tt.args.d, tt.args.m); got != tt.want {
				t.Errorf("daysSoFar() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_LeapYears(t *testing.T) {
	type args struct {
		ys int
		ye int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "same year, not leap year",
			args: args{1900, 1900},
			want: 0,
		},
		{
			name: "same year, leap year",
			args: args{2008, 2008},
			want: 1,
		},
		{
			name: "two leap years",
			args: args{2004, 2008},
			want: 2,
		},
		{
			name: "two leap years in the range",
			args: args{2003, 2009},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LeapYears(tt.args.ys, tt.args.ye); got != tt.want {
				t.Errorf("leapYears() = %v, want %v", got, tt.want)
			}
		})
	}
}
