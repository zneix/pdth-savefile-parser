package main

import (
	"flag"
	"fmt"
	"math"
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

func reputation(skillProgress map[any]any, cash float32) float64 {
	var rep float32

	// proceed with regular levels
	for _, v := range skillProgress {
		rep += v.(float32)
	}

	// calculate virtual reputation
	// vrep formula: https://github.com/steam-test1/pdth-source/blob/428335cb05cc26a697ac05fac9456742bf49e769/lib/tweak_data/tweakdata.lua#L1661
	// floor((306750 + sqrt((306750^2) - 517000 * (1333214 - 3088000 - $cash + 101904000/47))) / 5500)
	if (rep == 145 && cash > 700000) || (rep >= 193 && cash > 1333214) {
		return math.Max(float64(rep), math.Floor((306750+math.Sqrt(math.Pow(306750, 2)-517000*(1333214-3088000-float64(cash)+101904000/47)))/5500))
	}

	return float64(rep)
}

var (
	printJSON    *bool
	savefileName *string
)

func init() {
	printJSON = flag.Bool("json", false, "controls whether or not print JSON of all available game data")
	savefileName = flag.String("file", "save099.sav", "path to the PAYDAY's savefile")
}

func main() {
	flag.Parse()

	savefile := NewSaveFile(*savefileName)
	data := savefile.Contents().(map[any]any)

	if !*printJSON {
		cash := data["ExperienceManager"].(map[any]any)["total"]
		fmt.Printf("total cash: %.0f\n", cash)
		fmt.Println("reputation:", reputation(data["UpgradesManager"].(map[any]any)["progress"].(map[any]any), cash.(float32)))
		return
	}

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
