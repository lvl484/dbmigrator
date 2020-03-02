package main

import "testing"

func TestCompareType(t *testing.T) {
	type args struct {
		coltype string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"version1", args{`double precision`}, true},
		{"version2", args{`numeric`}, true},
		{"version3", args{`real`}, true},
		{"version4", args{`decimal`}, true},
		{"version5", args{`int`}, false},
		{"version6", args{`boolean`}, false},
		{"version7", args{`string`}, false},
		{"version8", args{`text`}, false},
		{"version9", args{``}, false},
		{"version10", args{`nfsndsadsada`}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CompareType(tt.args.coltype); got != tt.want {
				t.Errorf("CompareType() = %v, want %v", got, tt.want)
			}
		})
	}
}
