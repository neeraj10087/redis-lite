package resp

import (
	"bufio"
	"strings"
	"testing"
	"fmt"
)

func TestParse(t *testing.T) {
	input := "*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"
	reader := bufio.NewReader(strings.NewReader(input))

	val, err := Parse(reader)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("Type:", val.Typ)
	fmt.Println("Command:", val.Array[0].Bulk)
	fmt.Println("Key:",     val.Array[1].Bulk)
	fmt.Println("Value:",   val.Array[2].Bulk)
}