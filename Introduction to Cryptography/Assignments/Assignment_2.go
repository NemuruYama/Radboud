package main

import (
	"fmt"
)

// LSFR cipher implementation with an output keystream length of 16 bits (split into 2 bytes)
func lsfrCipher(length int, key byte, taps []int) (byte, byte) {
	var out uint16
	state := key
	for i := 0; i < 16; i++ {
		var z byte
		state, z = iterateLSFR(length, state, taps)
		out <<= 1
		out |= uint16(z)
	}

	return byte(out >> 8), byte(out)
}

func iterateLSFR(length int, in byte, tapsList []int) (byte, byte) {
	i := length - 1 // index of last bit
	last := in & 1  // last bit (z)
	in >>= 1        // shift input to right

	var taps byte
	for _, tap := range tapsList {
		taps |= last << (i - tap)
	}
	in ^= taps        // add taps
	in |= (last << i) // set first bit to z

	return in, last
}

func calcState0(length int, keystream byte, taps []int) byte {
	var state, s0 byte            // s^0 = (s_0, s_1, s_2, s_3, s_4, s_5, ...)
	for i := 0; i < length; i++ { // Loop through output bits
		z := keystream >> (length - 1 - i) & 1 // Get z from the keystream output

		s_i := (state & 1) ^ z // Get last bit of state and XOR with z to get s_i of state 0
		s0 |= s_i << i         // Set s_i to pos i

		// Overwrite last bit to z
		state &^= 1
		state |= z

		state, _ = iterateLSFR(length, state, taps)
	}

	return s0
}

func main() {
	taps4 := []int{3}
	taps8 := []int{1, 3, 5}
	var Z0 byte = 145
	var Z1 byte = 68

	var key4, key8 byte
	for key4 = 0; key4 <= 0b1111; key4++ {
		fmt.Printf("Testing key for 4-bit LSFR: %04b\n", key4)
		out4_0, out4_1 := lsfrCipher(4, key4, taps4)
		fmt.Printf("\tOutput 4-bit LSFR: %08b %08b\n", out4_0, out4_1)

		// Get the output of the 8-bit LSFR by solving the mod 2^8 equation
		var out8_0 byte = byte((int(Z0) - int(out4_0)) % 256)
		var out8_1 byte = byte((int(Z1) - int(out4_1)) % 256)
		fmt.Printf("\tOutput 8-bit LSFR: %08b %08b\n", out8_0, out8_1)

		key8 = calcState0(8, out8_0, taps8)
		fmt.Printf("\tKey for 8-bit LSFR based on first byte: %08b\n", key8)

		newOut8_0, newOut8_1 := lsfrCipher(8, key8, taps8)
		// Sanity check to make sure key8 from calcState0 is correct for the first output byte
		if newOut8_0 != out8_0 {
			panic(fmt.Sprintf(`Expected %08b, actual %08b (key incorrect)`, out8_0, newOut8_0))
		}

		if newOut8_1 == out8_1 {
			fmt.Printf("\tThe rest of the output matches: %08b%08b == %08b%08b\n\n", out8_0, out8_1, newOut8_0, newOut8_1)
			break
		}

		fmt.Printf("\tThe rest of the output does NOT match: %08b%08b != %08b%08b\n\n", out8_0, out8_1, newOut8_0, newOut8_1)
	}

	fmt.Printf("4-bit key = %04b\n8-bit key = %08b\n", key4, key8)
}