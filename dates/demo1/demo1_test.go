package demo1

import "testing"

func Test_Count(t *testing.T) {
	type args struct {
		dateLeft  []int
		dateRight []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "same dates",
			args: args{[]int{2, 6, 1983}, []int{2, 6, 1983}},
			want: 0,
		},
		{
			name: "same month",
			args: args{[]int{2, 6, 1983}, []int{22, 6, 1983}},
			want: 19,
		},
		{
			name: "same month, out of order",
			args: args{[]int{22, 6, 1983}, []int{2, 6, 1983}},
			want: 19,
		},
		{
			name: "same year",
			args: args{[]int{4, 7, 1984}, []int{25, 12, 1984}},
			want: 173,
		},
		{
			name: "same year, out of order",
			args: args{[]int{25, 12, 1984}, []int{4, 7, 1984}},
			want: 173,
		},
		{
			name: "in a leap year",
			args: args{[]int{31, 1, 1984}, []int{1, 3, 1984}},
			want: 29,
		},
		{
			name: "in a leap year, out of order",
			args: args{[]int{1, 3, 1984}, []int{31, 1, 1984}},
			want: 29,
		},
		{
			name: "same day, different years",
			args: args{[]int{4, 7, 1984}, []int{4, 7, 1985}},
			want: 364,
		},
		{
			name: "same day, different years, with a leap year",
			args: args{[]int{4, 7, 1983}, []int{4, 7, 1985}},
			want: 730,
		},
		{
			name: "out of order, the same day and month, one year appart, leap year",
			args: args{[]int{3, 8, 1984}, []int{3, 8, 1983}},
			want: 365,
		},
		{
			name: "out of order, multiple year apart",
			args: args{[]int{3, 1, 1989}, []int{3, 8, 1983}},
			want: 1979,
		},
		{
			name: "out of order, with typo?",
			args: args{[]int{3, 1, 1989}, []int{3, 8, 1983}},
			want: 2036,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Interval(tt.args.dateLeft, tt.args.dateRight); got != tt.want {
				t.Errorf("Count() = %v, want %v", got, tt.want)
			}
		})
	}
}
