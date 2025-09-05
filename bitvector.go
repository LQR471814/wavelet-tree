package wavelettree

import (
	"fmt"
	"unsafe"
)

// BitVector is effectively a slice of bits.
type BitVector struct {
	bitlength uint64
	bytes     []byte
}

func NewBitVector(bitlength uint64) BitVector {
	bytelength := bitlength/8 + 1
	vec := BitVector{
		bitlength: bitlength,
		bytes:     make([]byte, bytelength),
	}
	return vec
}

// Get8 allows you to get 1-8 bits from the bitvector at once and return it as
// a uint8
func (v BitVector) Get8(size uint8, i uint64) uint8 {
	if size == 0 || size > 8 {
		panic(fmt.Sprintf("invalid bitsize: %d", size))
	}
	return getbits[uint8](8, size, v, i)
}

// Get16 allows you to get 1-16 bits from the bitvector at once and return it
// as a uint16
func (v BitVector) Get16(size uint8, i uint64) uint16 {
	if size == 0 || size > 16 {
		panic(fmt.Sprintf("invalid bitsize: %d", size))
	}
	return getbits[uint16](16, size, v, i)
}

// Get32 allows you to get 1-32 bits from the bitvector at once and return it
// as a uint32
func (v BitVector) Get32(size uint8, i uint64) uint32 {
	if size == 0 || size > 32 {
		panic(fmt.Sprintf("invalid bitsize: %d", size))
	}
	return getbits[uint32](32, size, v, i)
}

// Get64 allows you to get 1-64 bits from the bitvector at once and return it
// as a uint64
func (v BitVector) Get64(size uint8, i uint64) uint64 {
	if size == 0 || size > 64 {
		panic(fmt.Sprintf("invalid bitsize: %d", size))
	}
	return getbits[uint64](64, size, v, i)
}

// Set8 allows you to get 1-8 bits from the bitvector at once and return it as
// a uint8
func (v BitVector) Set8(size uint8, i uint64, value uint8) {
	if size == 0 || size > 8 {
		panic(fmt.Sprintf("invalid bitsize: %d", size))
	}
	setbits(8, size, v, i, value)
}

// Set16 allows you to get 1-16 bits from the bitvector at once and return it
// as a uint16
func (v BitVector) Set16(size uint8, i uint64, value uint16) {
	if size == 0 || size > 16 {
		panic(fmt.Sprintf("invalid bitsize: %d", size))
	}
	setbits(16, size, v, i, value)
}

// Set32 allows you to get 1-32 bits from the bitvector at once and return it
// as a uint32
func (v BitVector) Set32(size uint8, i uint64, value uint32) {
	if size == 0 || size > 32 {
		panic(fmt.Sprintf("invalid bitsize: %d", size))
	}
	setbits(32, size, v, i, value)
}

// Set64 allows you to get 1-64 bitsv.bytes from the bitvector at once and return it
// as a uint64
func (v BitVector) Set64(size uint8, i uint64, value uint64) {
	if size == 0 || size > 64 {
		panic(fmt.Sprintf("invalid bitsize: %d", size))
	}
	setbits(64, size, v, i, value)
}

// Append8 allows you to get 1-8 bits from the bitvector at once and return it as
// a uint8
func (v BitVector) Append8(size uint8, value uint8) BitVector {
	if size == 0 || size > 8 {
		panic(fmt.Sprintf("invalid bitsize: %d", size))
	}
	return appendbits(8, size, v, value)
}

// Append16 allows you to get 1-16 bits from the bitvector at once and return it
// as a uint16
func (v BitVector) Append16(size uint8, value uint16) BitVector {
	if size == 0 || size > 16 {
		panic(fmt.Sprintf("invalid bitsize: %d", size))
	}
	return appendbits(16, size, v, value)
}

// Append32 allows you to get 1-32 bits from the bitvector at once and return it
// as a uint32
func (v BitVector) Append32(size uint8, value uint32) BitVector {
	if size == 0 || size > 32 {
		panic(fmt.Sprintf("invalid bitsize: %d", size))
	}
	return appendbits(32, size, v, value)
}

// Append64 allows you to get 1-64 bits from the bitvector at once and return it
// as a uint64
func (v BitVector) Append64(size uint8, value uint64) BitVector {
	if size == 0 || size > 64 {
		panic(fmt.Sprintf("invalid bitsize: %d", size))
	}
	return appendbits(64, size, v, value)
}

// - bitsize can be any number of bits from 1-64
// - bytesize must be one of 8, 16, 32, or 64
func getbits[T uint8 | uint16 | uint32 | uint64](bytesize, bitsize uint8, v BitVector, i uint64) (result T) {
	if i >= v.bitlength {
		panic(fmt.Sprintf("get bit index out of range: [%d]", i))
	}

	byteslice := *(*[]T)(unsafe.Pointer(&v.bytes))

	byte := i / uint64(bytesize)
	bit := uint8(i % uint64(bytesize))

	allones := ^T(0)

	// mask creates a bit mask that has 1's in the places for the
	// target bit being retrieved.
	var mask T = allones >> (bytesize - bit)

	result = (byteslice[byte] >> bit) & mask

	// overlap threshold stores the index at which the rest of the
	// bits would spill over to the next byte/uint
	var overlapThreshold uint8 = (bytesize + 1) - bitsize

	// amount of bits set in the current byte
	var currentSet = bytesize - bit

	// overlap amount is: amount of bits to be found in the next byte after
	// bits are found in current byte
	var overlapAmount uint8 = bitsize - currentSet

	if bit >= overlapThreshold {
		next := byteslice[byte+1]
		var nextmask T = allones >> overlapAmount
		var overlap T = next & nextmask
		result = result | (overlap << currentSet)
	}

	return
}

func setbits[T uint8 | uint16 | uint32 | uint64](bytesize, bitsize uint8, v BitVector, i uint64, value T) {
	if i >= v.bitlength {
		panic(fmt.Sprintf("set bit index out of range: [%d]", i))
	}

	byteslice := *(*[]T)(unsafe.Pointer(&v.bytes))

	byte := i / uint64(bytesize)
	bit := uint8(i % uint64(bytesize))

	fmt.Println(byte, bit)

	allones := ^T(0)

	// mask creates a bit mask that has 1's in the places for the
	// target bit being retrieved.
	var mask T = allones >> (bytesize - bit)

	value = value & mask

	// bit defines everything including and after a particular bit index to be
	// the set content, therefore the underlying content that should be kept
	// would be the inverse of everything including and after the particular
	// bit index
	underlying := byteslice[byte] & ^mask
	byteslice[byte] = underlying | (value << bit)

	fmt.Printf("set: %016b | underlying(mask): %016b | underlying: %016b\n", value, underlying, byteslice[byte])

	// overlap threshold stores the index at which the rest of the
	// bits would spill over to the next byte/uint
	var overlapThreshold uint8 = (bytesize + 1) - bitsize

	fmt.Printf("%d + 1 - %d = %d\n", bytesize, bitsize, overlapThreshold)

	// amount of bits set in the current byte
	var currentSet uint8 = bytesize - bit

	// overlap amount is: amount of bits to be found in the next byte after
	// bits are found in current byte
	var overlapAmount uint8 = bitsize - currentSet

	if bit >= overlapThreshold {
		next := byteslice[byte+1]
		var nextMask T = allones >> overlapAmount
		nextupSurround := next & (^nextMask)

		fmt.Printf("next: %016b | nextmask: %016b\n", next, nextMask)

		// remove the bits in value that have already been set in the current
		// byte and set those bits in the next byte
		byteslice[byte+1] = nextupSurround | (value >> currentSet)
	}
}

func appendbits[T uint8 | uint16 | uint32 | uint64](bytesize, bitsize uint8, v BitVector, value T) BitVector {
	byteslice := *(*[]T)(unsafe.Pointer(&v.bytes))
	v.bytes = *(*[]byte)(unsafe.Pointer(&byteslice))

	originalEnd := v.bitlength
	v.bitlength += uint64(bitsize)
	byteLen := v.bitlength/uint64(bytesize) + 1

	if int(byteLen) > len(byteslice) {
		byteslice = append(byteslice, 0)
	}
	setbits(bytesize, bitsize, v, originalEnd, value)

	return v
}

// Length returns the bit length of the bitvector.
func (v BitVector) Length() uint64 {
	return v.bitlength
}

func (v BitVector) String() string {
	out := make([]byte, 0, v.bitlength)
	for _, b := range v.bytes {
		mask := uint8(1)
		for range 8 {
			if mask&b > 0 {
				out = append(out, '1')
			} else {
				out = append(out, '0')
			}
			mask <<= 1
		}
	}
	return string(out[:v.bitlength])
}
