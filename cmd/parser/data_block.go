package main

// DataBlock is used for sections that do not contain game data and are rather meta/technical
// TODO: Consider dropping since it's not really used
type DataBlock struct {
	Version  uint32
	Size     uint32
	Data     string
	Checksum string
}

func NewDataBlock(reader *BinaryReader) *DataBlock {
	size := reader.OpUint32()
	version := reader.OpUint32()
	block := &DataBlock{
		Size:    size,
		Version: version,
	}
	/*block.Data =*/ reader.OpSkip(int(block.Size) - 16)
	/*block.Checksum =*/ reader.OpSkip(16)

	return block
}
