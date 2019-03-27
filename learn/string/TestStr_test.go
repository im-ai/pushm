package string

import "testing"

func Test_jionstradd(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jionstradd()
		})
	}
}

func Test_jionstrbuff(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jionstrbuff()
		})
	}
}

//10	 170609750 ns/op
func Benchmark_joinstradd(b *testing.B) {
	benchmarks := []struct {
		name string
	}{
		// TODO: benchmarks
		{name: ""},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				jionstradd()
			}
		})
	}
}

//5000	    227413 ns/op
func Benchmark_jionstrbuff(b *testing.B) {
	benchmarks := []struct {
		name string
	}{
		// TODO: benchmarks
		{name: ""},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				jionstrbuff()
			}
		})
	}
}
