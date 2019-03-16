package main

import "testing"

// 测试
func TestAdd(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
		{"case 0", args{1, 2}, 3},
		{"case 1", args{2, 2}, 4},
		{"case 2", args{5, 2}, 7},
		{"case 3", args{11, 22}, 33},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Add(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

// 压力测试 模板
func Benchmark_Add(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Add(i, i)
	}
}
