package main

import "math"

// GameDataBlock
type GameDataBlock struct {
	reader  *BinaryReader
	version uint32
	data    any
}

func NewGameDataBlock(reader *BinaryReader) *GameDataBlock {
	//$this->reader = $reader;

	//$this->version = $this->reader->uint32();
	//$this->data = $this->next();
	//fmt.Println("DEBUG game data block, offset before init:", reader.offset)
	block := &GameDataBlock{
		reader:  reader,
		version: reader.OpUint32(),
	}
	//fmt.Println("DEBUG game data block, version:", block.version)
	block.data = block.next()
	return block
}

func (gdb *GameDataBlock) Data() any {
	return gdb.data
}

func (gdb *GameDataBlock) next() any {
	dataType := string(gdb.reader.OpBytes(4))
	//fmt.Printf("dataType: %#v\n", dataType)
	switch dataType {
	case "\xDE\xE1\xB9\xAA":
		//fmt.Println("nexting table")
		return gdb.ValueTable()
	case "\x8A\x24\xBC\xB1":
		//fmt.Println("nexting number")
		return gdb.ValueNumber()
	case "\xE0\xFF\x5C\xB4":
		//fmt.Println("nexting string")
		return gdb.ValueString()
	case "\x4F\xC6\xE2\x84":
		//fmt.Println("nexting bool")
		return gdb.ValueBool()
	}
	panic("didn't find a match for any of data types BabyRage")
}

func (gdb *GameDataBlock) ValueBool() bool {
	char := gdb.reader.OpChar()
	return char != 0
}

func (gdb *GameDataBlock) ValueString() string {
	size := gdb.reader.OpFindStringSize('\x00')
	//fmt.Println("DEBUG STRING found size:", size)
	if size >= 0 {
		return gdb.reader.OpString(size)
	}

	return ""
}

// TODO: consider using generics for either float32|uint32
// or just don't cast some weird float32(int(...)) and return as-is
func (gdb *GameDataBlock) ValueNumber() float32 {
	value := gdb.reader.OpFloat32()
	if value >= 1 {
		if math.Floor(float64(value)) == float64(value) {
			return float32(int(value))
		}
	}
	return value
}

func (gdb *GameDataBlock) ValueTable() any {
	result := make(map[any]any, 0)

	// fetch the table size which is present in the binary data before the table itself
	size := gdb.reader.OpUint32()
	for i := 0; i < int(size); i++ {
		key := gdb.next()
		value := gdb.next()
		result[key] = value
	}
	return result
}
