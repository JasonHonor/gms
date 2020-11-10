package utils

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestRemovePart(t *testing.T) {
	data := []byte("this is a test,do it twice!")

	//CopyTo fixed length. Convert slice to fixed array.
	var buf [25]byte
	copy(buf[:], data[:25])

	var param []byte

	//Convert fixed array to slice
	param = buf[:25]

	ret := RemoveStringByTag(param, "do", "it", "X")
	fmt.Printf("Result:%v\n", string(ret))
}

func TestRemovePartBytes(t *testing.T) {
	data, _ := hex.DecodeString("112208080808080808080811223344556677880808080808080808083344")

	//CopyTo fixed length. Convert slice to fixed array.
	var buf [25]byte
	copy(buf[:], data[:25])

	tag, _ := hex.DecodeString("080808080808080808")

	ret := RemoveStrByTagBytes(data,
		tag,
		tag, "X")

	fmt.Printf("Result:%v\n", hex.EncodeToString(ret))
}

func TestMacFormat(t *testing.T) {

	s0 := MacFormat("08:f1:ea:ee:f9")
	fmt.Printf("result:%s\n", s0)

	s1 := MacFormat("08:f1:ea:ee:f9:a8")
	fmt.Printf("result:%s\n", s1)

	s2 := MacFormat("08-f1-ea-ee-f9-a8")
	fmt.Printf("result:%s\n", s2)

	s3 := MacFormat("08.f1.ea.ee.f9.a8")
	fmt.Printf("result:%s\n", s3)

	s4 := MacFormat(" 08.f1.ea.ee.f9.a8 ")
	fmt.Printf("result:%s\n", s4)

}
