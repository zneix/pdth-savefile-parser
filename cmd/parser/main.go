package main

import (
	"fmt"
	"os"

	"github.com/json-iterator/go"
)

// SaveFile
type SaveFile struct {
	version uint32
	path    string

	header   *DataBlock
	gamedata *GameDataBlock
	//footer   *DataBlock
}

func NewSaveFile(path string) *SaveFile {
	savefile := &SaveFile{
		path: path,
	}

	content, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("failed to open/read file at %q: %w", path, err))
	}

	filesize := len(content)

	reader := NewBinaryReader(content, filesize, 0)
	savefile.version = reader.OpUint32()
	savefile.header = NewDataBlock(reader)
	savefile.gamedata = NewGameDataBlock(reader)
	//savefile.footer = NewDataBlock(reader) // XXX: Doesn't work yet, leave it unimplemented

	return savefile
}

func (sf *SaveFile) Version() uint32 {
	return sf.version
}

func (sf *SaveFile) Contents() any {
	return sf.gamedata.Data()
}

func main() {
	savefile := NewSaveFile("./save099.sav")
	data := savefile.Contents().(map[any]any)
	//fmt.Printf("%.0f\n", data["ExperienceManager"].(map[any]any)["total"])

	dataJSON, err := jsoniter.Config{
		EscapeHTML:                    false,
		MarshalFloatWith6Digits:       true,
		ObjectFieldMustBeSimpleString: true,
		SortMapKeys:                   true,
	}.Froze().MarshalToString(data)
	if err != nil {
		fmt.Println("fdm:", err)
		return
	}
	fmt.Println(dataJSON)
}
