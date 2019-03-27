package api

//
// Slow software only Sub-Tweak Key derivation routines.
//

// A.2 RCON constants
var rcons = [STKCount]byte{
	0x2f, 0x5e, 0xbc, 0x63, 0xc6, 0x97, 0x35, 0x6a,
	0xd4, 0xb3, 0x7d, 0xfa, 0xef, 0xc5, 0x91, 0x39,
	0x72,
}

func H(t *[STKSize]byte) {
	t[0], t[1], t[2], t[3], t[4], t[5], t[6], t[7], t[8], t[9], t[10], t[11], t[12], t[13], t[14], t[15] = t[7], t[0], t[13], t[10], t[11], t[4], t[1], t[14], t[15], t[8], t[5], t[2], t[3], t[12], t[9], t[6]
}

func lfsr2(t *[STKSize]byte) {
	for i, x := range t {
		// x7 || x6 || x5 || x4 || x3 || x2 || x1 || x0 ->
		// x6 || x5 || x4 || x3 || x2 || x1 || x0 || x7 ^ x5
		x7, x5 := x>>7, (x>>5)&1
		t[i] = (x << 1) | (x7 ^ x5)
	}
}

func lfsr3(t *[STKSize]byte) {
	for i, x := range t {
		// x7 || x6 || x5 || x4 || x3 || x2 || x1 || x0 ->
		// x0 ^ x6 || x7 || x6 || x5 || x4 || x3 || x2 || x1
		x0, x6 := x&1, (x>>6)&1
		t[i] = (x >> 1) | ((x0 ^ x6) << 7)
	}
}

func xorRC(t *[STKSize]byte, i int) {
	rcon := rcons[i]
	var rc = [STKSize]byte{
		1, 2, 4, 8,
		rcon, rcon, rcon, rcon,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}

	// Note: No need to XOR the last 8 bytes, so skip that.
	XORBytes(t[:], t[:], rc[:], 8)
}

func STKDeriveK(key []byte, derivedKs *[STKCount][STKSize]byte) {
	var tk2, tk3 [STKSize]byte

	copy(tk2[:], key[16:32]) // Tk2 = W2
	copy(tk3[:], key[0:16])  // Tk3 = W3

	// i == 0
	XORBytes(derivedKs[0][:], tk2[:], tk3[:], STKSize)
	xorRC(&derivedKs[0], 0)

	// i == 1 ... i == 16
	for i := 1; i <= Rounds; i++ {
		// Tk2(i+1) = h(LFSR2(Tk2(i)))
		lfsr2(&tk2)
		H(&tk2)

		// Tk3(i+1) = h(LFSR3(Tk3(i)))
		lfsr3(&tk3)
		H(&tk3)

		XORBytes(derivedKs[i][:], tk2[:], tk3[:], STKSize)
		xorRC(&derivedKs[i], i)
	}
}