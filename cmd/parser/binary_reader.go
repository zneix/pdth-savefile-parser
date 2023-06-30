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
	// unpack 'C'
	convTypeCharacter convType = iota
	// unpack 'V'
	convTypeUint32
	// unpack 'g'
	convTypeFloat32
)

func (br *BinaryReader) convert(dataType convType, dataAmount int) any {
	br.checkBounds(dataAmount)

	switch dataType {
	// character: 1-byte
	case convTypeCharacter:
		//for i := 0; i < dataAmount; i++ {
		//// since we need a 0 or 1 to create boolean's value here, it's enough if we simply return the binary character by itself
		////result = append(result, br.binary[br.offset : br.offset+1][0])
		//result = append(result, br.binary[br.offset])
		//}
		result := br.binary[br.offset : br.offset+dataAmount]
		br.offset += dataAmount
		return result

	// unsigned integer: 4-bytes
	case convTypeUint32:
		//for i := 0; i < dataAmount; i++ {
		//result = append(result, binary.LittleEndian.Uint32(br.binary[br.offset:br.offset+4]))
		//}
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
		//for i := 0; i < dataAmount; i++ {
		//ourData := binary.LittleEndian.Uint32(br.binary[i*4 : i*4+4])
		//result = append(result, math.Float32frombits(ourData))
		//}
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

	//data := br.binary[br.offset : br.offset+byteAmount]
	//dataHex := hex.EncodeToString(data)
	//dataString, err := hex.DecodeString(dataHex)
	//if err != nil {
	//panic(fmt.Errorf("failed to convert %d bytes at offset %d: %w", byteAmount, br.offset, err))
	//}

	// seems like the whole hex-fiesta is rather unnecessary and this works just fine
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

// OpFindStringSize finds at which index []byte is present in the binary string
// FIXME: comments below are incorrect
// ~~TODO~~: Consider dropping this and instead just use
// bytes.Index or bytes.IndexByte in the 1 place where this method is needed
func (br *BinaryReader) OpFindStringSize(singleByte byte) int { // returns just the position
	return bytes.IndexByte(br.binary[br.offset:], singleByte)
	//pos := strings.Index(string(br.binary[br.offset:]), string(singleByte)) // slow AF

	//if pos != -1 {
	//return pos - br.offset // INCORRECT LOGIC, pos already skips br.offset because of '[br.offset:]' construction
	//}
	//return -1
}
