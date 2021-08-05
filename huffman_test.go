package huffman

import (
	"log"
	"testing"
)

func TestDevelop(t *testing.T) {

	test := Leaves{
		&Node{Value: "red", Weight: 1},
		&Node{Value: "alfred", Weight: 2},
		&Node{Value: "radixholms@gmail.com", Weight: 3},
		&Node{Value: "test1", Weight: 4},
		&Node{Value: "test2", Weight: 14},
		&Node{Value: "test3", Weight: 55},
		&Node{Value: "test4", Weight: 56},
		&Node{Value: "test5", Weight: 70},
		&Node{Value: "test6", Weight: 72},
		&Node{Value: "test7", Weight: 78},
		&Node{Value: "test8", Weight: 102},
	}
	enc, dec := test.Build().NewCoder()
	log.Println(enc)
	log.Println(dec)
	b := enc.EncodeStringSlice([]string{"red", "red", "test8", "alfred", "alfred", "radixholms@gmail.com", "radixholms@gmail.com"})
	log.Printf("%b", b)

	res := dec.DecodeToStringSlice(b)
	log.Println(res)
}
