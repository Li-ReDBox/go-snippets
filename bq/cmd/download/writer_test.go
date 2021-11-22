package main

import (
	"io/ioutil"
	"testing"
)

func Test_save(t *testing.T) {
	type args struct {
		c string
		f string
	}
	tests := []struct {
		name string
		args args
	}{
		// This does not work as save uses solid ioutil.Writefile
		// the only way to check is to read back and compare it with c
		{
			name: "first",
			args: args{
				c: "test",
				f: "test-first.txt",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			save(tt.args.c, tt.args.f)

			c, _ := ioutil.ReadFile(tt.args.f)
			if string(c) != tt.args.c {
				t.Errorf("save saved %s, should be %s", c, tt.args.c)
			}
		})
	}
}
