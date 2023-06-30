package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// BinaryReader
type BinaryReader struct {
	binary []byte
	offset int
	size   int
}

func NewBinaryReader(binaryData []byte, size, offset int) *BinaryReader {
	if size <= 0 {
		size = len(binaryData)
	}

	return &BinaryReader{
		binary: binaryData,
		size:   size,
		offset: offset,
	}
}

func (br *BinaryReader) Offset() int {
	return br.offset
}

func (br *BinaryReader) checkBounds(byteAmount int) {
	if byteAmount < 0 {
		panic("Byte count must be positive.")
	}

	if br.size-br.offset < byteAmount {
		panic("Not enough data.")
	}
}

type convType int

const (
	convTypeCharacter convType = iota // unpack 'C'
	convTypeUint32                    // unpack 'V'
	convTypeFloat32                   // unpack 'g'
)

func (br *BinaryReader) convert(dataType convType, dataAmount int) any {
	br.checkBounds(dataAmount)

	switch dataType {
	// character: 1-byte
	case convTypeCharacter:
		result := br.binary[br.offset : br.offset+dataAmount]
		br.offset += dataAmount
		return result

	// unsigned integer: 4-bytes
	case convTypeUint32:
		result := make([]uint32, dataAmount)
		_ = binary.Read(
			bytes.NewReader(br.binary[br.offset:br.offset+(4*dataAmount)]),
			binary.LittleEndian,
			&result,
		)
		br.offset += dataAmount
		return result

	// float: 4-bytes
	case convTypeFloat32:
		result := make([]float32, dataAmount)
		_ = binary.Read(bytes.NewReader(br.binary[br.offset:br.offset+(4*dataAmount)]), binary.LittleEndian, &result)
		br.offset += dataAmount

		return result

	default:
		panic(fmt.Sprintf("Unhanled convType '%v' passed to BinaryReader.convert", dataType))
	}
}

func (br *BinaryReader) OpChar() byte {
	return br.convert(convTypeCharacter, 1).([]byte)[0]
}

func (br *BinaryReader) OpUint32() uint32 {
	return br.convert(convTypeUint32, 4).([]uint32)[0]
}

func (br *BinaryReader) OpFloat32() float32 {
	return br.convert(convTypeFloat32, 4).([]float32)[0]
}

func (br *BinaryReader) OpString(byteAmount int) string {
	br.checkBounds(byteAmount)

	result := string(br.binary[br.offset : br.offset+byteAmount])
	br.offset += byteAmount + 1 // include null byte terminator
	return result
}

func (br *BinaryReader) OpBytes(byteAmount int) []byte {
	br.checkBounds(byteAmount)

	result := br.binary[br.offset : br.offset+byteAmount]
	br.offset += byteAmount
	return result
}

// OpDump as far as I can understand is some debugging/printing method
// used to specify amount of string-characters to print
// because we want a string, but we get a hex pairs, we gotta multiply charAmount by 2 to account for that
/*
func (br *BinaryReader) OpDump(charAmount int) { // returns void
	br.checkBounds(charAmount)

	result := hex.EncodeToString(br.binary[br.offset : br.offset+(charAmount*2)])
	fmt.Println(result)
}
*/

func (br *BinaryReader) OpSkip(byteAmount int) { // returns void
	br.checkBounds(byteAmount)
	br.offset += byteAmount
}

// OpFindStringSize finds at which index singleByte is present in the binary string
func (br *BinaryReader) OpFindStringSize(singleByte byte) int {
	return bytes.IndexByte(br.binary[br.offset:], singleByte)
}
