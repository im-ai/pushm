package main

type S struct {
	M *int
}

func main() {
	var x S
	var i int
	ref(&i, &x)
	// go tool compile -S test1.go
}

func ref(i *int, s *S) {
	s.M = i
}
