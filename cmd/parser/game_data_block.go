package main

// GameDataBlock
type GameDataBlock struct {
	reader  *BinaryReader
	version uint32
	data    any
}

func NewGameDataBlock(reader *BinaryReader) *GameDataBlock {
	block := &GameDataBlock{
		reader:  reader,
		version: reader.OpUint32(),
	}
	block.data = block.next()
	return block
}

func (gdb *GameDataBlock) Data() any {
	return gdb.data
}

func (gdb *GameDataBlock) next() any {
	dataType := string(gdb.reader.OpBytes(4))
	switch dataType {
	case "\xDE\xE1\xB9\xAA":
		return gdb.ValueTable()
	case "\x8A\x24\xBC\xB1":
		return gdb.ValueNumber()
	case "\xE0\xFF\x5C\xB4":
		return gdb.ValueString()
	case "\x4F\xC6\xE2\x84":
		return gdb.ValueBool()
	}
	panic("didn't find a match for any of data types BabyRage")
}

func (gdb *GameDataBlock) ValueBool() bool {
	return gdb.reader.OpCharacter() != 0
}

func (gdb *GameDataBlock) ValueString() string {
	size := gdb.reader.OpFindStringSize('\x00')
	if size >= 0 {
		return gdb.reader.OpString(size)
	}

	return ""
}

func (gdb *GameDataBlock) ValueNumber() float32 {
	return gdb.reader.OpFloat32()
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
