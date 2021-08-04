package huffman

import (
	"log"
	"testing"
)

func TestDevelop(t *testing.T) {

	test := []byte{'a', 'v', 'b', 'c'}
	test = append(test, 'o')
	log.Printf("%p", test)

	test2 := test
	log.Printf("%p", []byte(string(test2)))

	// test := []string{"red", "alfred", "radixholms@gmail.com"}

	// test := Leaves{
	// 	&Node{Value: "red"},
	// 	&Node{Value: "alfred"},
	// 	&Node{Value: "radixholms@gmail.com"},
	// }
	// enc, dec := test.Build().NewCoder()
	// log.Println(enc)
	// b := enc.EncodeStringSlice([]string{"red", "red", "alfred", "alfred", "radixholms@gmail.com", "radixholms@gmail.com"})
	// log.Printf("%b", b)

	// res := dec.DecodeToStringSlice(b)
	// log.Println(res)
}
