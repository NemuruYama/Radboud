package main

import (
	"crypto/des"
	"fmt"
	"math/rand"
)

func main() {
	key := make([]byte, 8)
	rand.Read(key)
	key[len(key)-1] |= 0b00000001 // Set last bit to one
	// key[len(key)-1] &= 0b11111110 // Set last bit to zero
	fmt.Printf("Key:\t\t\t\t\t%08b\n", key)
	c, err := des.NewCipher(key)
	if err != nil {
		panic(err)
	}
	src := make([]byte, des.BlockSize)
	rand.Read(src)
	dst := make([]byte, des.BlockSize)

	fmt.Printf("Source:\t\t\t\t\t%08b\n", src)

	c.Encrypt(dst, src)
	fmt.Printf("Output DES(src):\t\t%08b\n", dst)
	fmt.Println()

	src_c := make([]byte, des.BlockSize)
	for i := range src {
		src_c[i] = ^src[i]
	}
	dst = make([]byte, des.BlockSize)

	fmt.Printf("Source complement:\t\t%08b\n", src_c)

	c.Encrypt(dst, src_c)
	fmt.Printf("Output:\t\t\t\t\t%08b\n", dst)
	fmt.Println()

	// Complement
	for i := range key {
		key[i] = ^key[i]
	}
	fmt.Printf("Complement key:\t\t\t%08b\n", key)
	c, err = des.NewCipher(key)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Source:\t\t\t\t\t%08b\n", src)
	c.Encrypt(dst, src)
	fmt.Printf("Output:\t\t\t\t\t%08b\n", dst)

	for i := range dst {
		dst[i] = ^dst[i]
	}
	fmt.Printf("Output complement:\t\t%08b\n", dst)
}
