package demo2

import (
	"testing"
)

func TestDate_Interval(t *testing.T) {
	type fields struct {
		Year  int
		Month int
		Day   int
	}
	type args struct {
		o Date
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name:   "same dates",
			fields: fields{1983, 6, 2},
			args:   args{Date{1983, 6, 2}},
			want:   0,
		},
		{
			name:   "same month",
			fields: fields{1983, 6, 2},
			args:   args{Date{1983, 6, 22}},
			want:   19,
		},
		{
			name:   "same month, out of order",
			fields: fields{1983, 6, 22},
			args:   args{Date{1983, 6, 2}},
			want:   19,
		},
		{
			name:   "same year",
			fields: fields{1984, 7, 4},
			args:   args{Date{1984, 12, 25}},
			want:   173,
		},
		{
			name:   "same year, out of order",
			fields: fields{1984, 7, 4},
			args:   args{Date{1984, 12, 25}},
			want:   173,
		},
		{
			name:   "in a leap year",
			fields: fields{1984, 1, 31},
			args:   args{Date{1984, 3, 1}},
			want:   29,
		},
		{
			name:   "in a leap year, out of order",
			fields: fields{1984, 3, 1},
			args:   args{Date{1984, 1, 31}},
			want:   29,
		},
		{
			name:   "same day, different years",
			fields: fields{1984, 7, 4},
			args:   args{Date{1985, 7, 4}},
			want:   364,
		},
		{
			name:   "same day, different years, with a leap year",
			fields: fields{1983, 7, 4},
			args:   args{Date{1985, 7, 4}},
			want:   730,
		},
		{
			name:   "out of order, the same day and month, one year appart, leap year",
			fields: fields{1984, 8, 3},
			args:   args{Date{1983, 8, 3}},
			want:   365,
		},
		{
			name:   "out of order, multiple year apart",
			fields: fields{1989, 1, 3},
			args:   args{Date{1983, 8, 3}},
			want:   1979,
		},
		{
			name:   "out of order, with typo?",
			fields: fields{1989, 1, 3},
			args:   args{Date{1983, 8, 3}},
			want:   2036,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Date{
				Year:  tt.fields.Year,
				Month: tt.fields.Month,
				Day:   tt.fields.Day,
			}
			if got := d.Interval(tt.args.o); got != tt.want {
				t.Errorf("Date.Interval() = %v, want %v", got, tt.want)
			}
		})
	}
}
