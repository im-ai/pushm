package string

import "bytes"

func jionstradd() {
	var str string
	for i := 0; i < 1000; i++ {
		str = str + "111111111111111111111"
	}
}

func jionstrbuff() {
	var buff bytes.Buffer
	for i := 0; i < 1000; i++ {
		buff.WriteString("111111111111111111111")
	}
}
