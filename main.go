package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const (
	itemsSecretKey = "PBG892FXX982ABC*"
)

type Item struct {
	ItemID             int    `json:"item_id"`
	EditableType       int    `json:"editable_type"`
	ItemCategory       int    `json:"item_category"`
	ActionType         int    `json:"action_type"`
	HitSoundType       int    `json:"hit_sound_type"`
	Name               string `json:"name"`
	Texture            string `json:"texture"`
	TextureHash        int    `json:"texture_hash"`
	ItemKind           int    `json:"item_kind"`
	Val1               int    `json:"val1"`
	TextureX           int    `json:"texture_x"`
	TextureY           int    `json:"texture_y"`
	SpreadType         int    `json:"spread_type"`
	IsStripeyWallpaper int    `json:"is_stripey_wallpaper"`
	CollisionType      int    `json:"collision_type"`
	BreakHits          string `json:"break_hits"`
	DropChance         int    `json:"drop_chance"`
	ClothingType       int    `json:"clothing_type"`
	Rarity             int    `json:"rarity"`
	MaxAmount          int    `json:"max_amount"`
	ExtraFile          string `json:"extra_file"`
	ExtraFileHash      int    `json:"extra_file_hash"`
	AudioVolume        int    `json:"audio_volume"`
	PetName            string `json:"pet_name"`
	PetPrefix          string `json:"pet_prefix"`
	PetSuffix          string `json:"pet_suffix"`
	PetAbility         string `json:"pet_ability"`
	SeedBase           int    `json:"seed_base"`
	SeedOverlay        int    `json:"seed_overlay"`
	TreeBase           int    `json:"tree_base"`
	TreeLeaves         int    `json:"tree_leaves"`
	SeedColor          Color  `json:"seed_color"`
	SeedOverlayColor   Color  `json:"seed_overlay_color"`
	GrowTime           int    `json:"grow_time"`
	Val2               int    `json:"val2"`
	IsRayman           int    `json:"is_rayman"`
	ExtraOptions       string `json:"extra_options"`
	Texture2           string `json:"texture2"`
	ExtraOptions2      string `json:"extra_options2"`
	DataPosition80     string `json:"data_position_80"`
	PunchOptions       string `json:"punch_options"`
	DataVersion12      string `json:"data_version_12"`
	IntVersion13       int    `json:"int_version_13"`
	IntVersion14       int    `json:"int_version_14"`
	DataVersion15      string `json:"data_version_15"`
	StrVersion15       string `json:"str_version_15"`
	StrVersion16       string `json:"str_version_16"`
	IntVersion17       int    `json:"int_version_17"`
	IntVersion18       int    `json:"int_version_18"`
}

type Color struct {
	A int `json:"a"`
	R int `json:"r"`
	G int `json:"g"`
	B int `json:"b"`
}

type ItemsData struct {
	Version   int    `json:"version"`
	ItemCount int    `json:"item_count"`
	Items     []Item `json:"items"`
}

func main() {
	encodePtr := flag.Bool("encode", false, "Encode items.dat")
	decodePtr := flag.Bool("decode", false, "Decode items.dat")
	getInfoPrt := flag.Bool("info", false, "Get information about items.dat")
	filePathPtr := flag.String("file", "", "Path to items.dat file")
	flag.Parse()

	if *encodePtr && *decodePtr && *getInfoPrt {
		fmt.Println("Please choose either encode or decode, not both.")
		return
	}

	if *filePathPtr == "" {
		fmt.Println("Please provide a file path using the -file flag.")
		return
	}

	if *encodePtr {
		encodeItems(*filePathPtr)
	} else if *decodePtr {
		decodeItems(*filePathPtr)
	} else if *getInfoPrt {
		getItemsInfo(*filePathPtr)
	} else {
		fmt.Println("Please choose either encode or decode.")
	}
}

func getItemsInfo(filePath string) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	itemsData, err := decodeItemsData(data)
	if err != nil {
		fmt.Println("Error decoding items.dat:", err)
		return
	}

	fmt.Println("Version:", itemsData.Version)
	fmt.Println("Item count:", itemsData.ItemCount)
}

func encodeItems(filePath string) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	itemsData, err := decodeItemsData(data)
	if err != nil {
		fmt.Println("Error decoding items.dat:", err)
		return
	}

	encodedData, err := encodeItemsData(itemsData)
	if err != nil {
		fmt.Println("Error encoding items.dat:", err)
		return
	}

	err = ioutil.WriteFile(filePath, encodedData, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("Items.dat encoded successfully!")
}

func decodeItems(filePath string) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	itemsData, err := decodeItemsData(data)
	if err != nil {
		fmt.Println("Error decoding items.dat:", err)
		return
	}

	err = writeItemsData(itemsData, filePath)
	if err != nil {
		fmt.Println("Error writing decoded data:", err)
		return
	}

	fmt.Println("Items.dat decoded successfully!")
}

func decodeItemsData(data []byte) (*ItemsData, error) {
	itemsData := &ItemsData{}
	itemsData.Version = int(data[0])<<8 | int(data[1])
	itemsData.ItemCount = int(data[2])<<24 | int(data[3])<<16 | int(data[4])<<8 | int(data[5])

	memPos := 6
	for i := 0; i < itemsData.ItemCount; i++ {
		if memPos >= len(data) {
			return nil, fmt.Errorf("reached end of data while decoding item %d", i)
		}
		item := Item{}
		item.ItemID = int(data[memPos])<<24 | int(data[memPos+1])<<16 | int(data[memPos+2])<<8 | int(data[memPos+3])
		memPos += 4
		item.EditableType = int(data[memPos])
		memPos++
		item.ItemCategory = int(data[memPos])
		memPos++
		item.ActionType = int(data[memPos])
		memPos++
		item.HitSoundType = int(data[memPos])
		memPos++
		item.Name = readString(data, memPos, true, item.ItemID)
		memPos += int(data[memPos-1]) + 2
		item.Texture = readString(data, memPos, false, 0)
		memPos += int(data[memPos-1]) + 2
		item.TextureHash = int(data[memPos])<<24 | int(data[memPos+1])<<16 | int(data[memPos+2])<<8 | int(data[memPos+3])
		memPos += 4
		item.ItemKind = int(data[memPos])
		memPos++
		item.Val1 = int(data[memPos])<<24 | int(data[memPos+1])<<16 | int(data[memPos+2])<<8 | int(data[memPos+3])
		memPos += 4
		item.TextureX = int(data[memPos])
		memPos++
		item.TextureY = int(data[memPos])
		memPos++
		item.SpreadType = int(data[memPos])
		memPos++
		item.IsStripeyWallpaper = int(data[memPos])
		memPos++
		item.CollisionType = int(data[memPos])
		memPos++
		item.BreakHits = strconv.Itoa(int(data[memPos]))
		if data[memPos]%6 != 0 {
			item.BreakHits += "r"
		} else {
			item.BreakHits = strconv.Itoa(int(data[memPos]) / 6)
		}
		memPos++
		item.DropChance = int(data[memPos])<<24 | int(data[memPos+1])<<16 | int(data[memPos+2])<<8 | int(data[memPos+3])
		memPos += 4
		item.ClothingType = int(data[memPos])
		memPos++
		item.Rarity = int(data[memPos])<<8 | int(data[memPos+1])
		memPos += 2
		item.MaxAmount = int(data[memPos])
		memPos++
		item.ExtraFile = readString(data, memPos, false, 0)
		memPos += int(data[memPos-1]) + 2
		item.ExtraFileHash = int(data[memPos])<<24 | int(data[memPos+1])<<16 | int(data[memPos+2])<<8 | int(data[memPos+3])
		memPos += 4
		item.AudioVolume = int(data[memPos])<<24 | int(data[memPos+1])<<16 | int(data[memPos+2])<<8 | int(data[memPos+3])
		memPos += 4
		item.PetName = readString(data, memPos, false, 0)
		memPos += int(data[memPos-1]) + 2
		item.PetPrefix = readString(data, memPos, false, 0)
		memPos += int(data[memPos-1]) + 2
		item.PetSuffix = readString(data, memPos, false, 0)
		memPos += int(data[memPos-1]) + 2
		item.PetAbility = readString(data, memPos, false, 0)
		memPos += int(data[memPos-1]) + 2
		item.SeedBase = int(data[memPos])
		memPos++
		item.SeedOverlay = int(data[memPos])
		memPos++
		item.TreeBase = int(data[memPos])
		memPos++
		item.TreeLeaves = int(data[memPos])
		memPos++
		item.SeedColor.A = int(data[memPos])
		memPos++
		item.SeedColor.R = int(data[memPos])
		memPos++
		item.SeedColor.G = int(data[memPos])
		memPos++
		item.SeedColor.B = int(data[memPos])
		memPos++
		item.SeedOverlayColor.A = int(data[memPos])
		memPos++
		item.SeedOverlayColor.R = int(data[memPos])
		memPos++
		item.SeedOverlayColor.G = int(data[memPos])
		memPos++
		item.SeedOverlayColor.B = int(data[memPos])
		memPos++
		memPos += 4 // skip ingredients
		item.GrowTime = int(data[memPos])<<24 | int(data[memPos+1])<<16 | int(data[memPos+2])<<8 | int(data[memPos+3])
		memPos += 4
		item.Val2 = int(data[memPos])<<8 | int(data[memPos+1])
		memPos += 2
		item.IsRayman = int(data[memPos])<<8 | int(data[memPos+1])
		memPos += 2
		item.ExtraOptions = readString(data, memPos, false, 0)
		memPos += int(data[memPos-1]) + 2
		item.Texture2 = readString(data, memPos, false, 0)
		memPos += int(data[memPos-1]) + 2
		item.ExtraOptions2 = readString(data, memPos, false, 0)
		memPos += int(data[memPos-1]) + 2
		item.DataPosition80 = toHexString(data[memPos : memPos+80])
		memPos += 80
		if itemsData.Version >= 11 {
			item.PunchOptions = readString(data, memPos, false, 0)
			memPos += int(data[memPos-1]) + 2
		}
		if itemsData.Version >= 12 {
			item.DataVersion12 = toHexString(data[memPos : memPos+13])
			memPos += 13
		}
		if itemsData.Version >= 13 {
			item.IntVersion13 = int(data[memPos])<<24 | int(data[memPos+1])<<16 | int(data[memPos+2])<<8 | int(data[memPos+3])
			memPos += 4
		}
		if itemsData.Version >= 14 {
			item.IntVersion14 = int(data[memPos])<<24 | int(data[memPos+1])<<16 | int(data[memPos+2])<<8 | int(data[memPos+3])
			memPos += 4
		}
		if itemsData.Version >= 15 {
			item.DataVersion15 = toHexString(data[memPos : memPos+25])
			memPos += 25
			item.StrVersion15 = readString(data, memPos, false, 0)
			memPos += int(data[memPos-1]) + 2
		}
		if itemsData.Version >= 16 {
			item.StrVersion16 = readString(data, memPos, false, 0)
			memPos += int(data[memPos-1]) + 2
		}
		if itemsData.Version >= 17 {
			item.IntVersion17 = int(data[memPos])<<24 | int(data[memPos+1])<<16 | int(data[memPos+2])<<8 | int(data[memPos+3])
			memPos += 4
		}
		if itemsData.Version >= 18 {
			item.IntVersion18 = int(data[memPos])<<24 | int(data[memPos+1])<<16 | int(data[memPos+2])<<8 | int(data[memPos+3])
			memPos += 4
		}
		itemsData.Items = append(itemsData.Items, item)
	}

	return itemsData, nil
}

func encodeItemsData(itemsData *ItemsData) ([]byte, error) {
	encodedData := make([]byte, 0, 2+4+itemsData.ItemCount*213)
	encodedData = append(encodedData, byte(itemsData.Version>>8), byte(itemsData.Version))
	encodedData = append(encodedData, byte(itemsData.ItemCount>>24), byte(itemsData.ItemCount>>16), byte(itemsData.ItemCount>>8), byte(itemsData.ItemCount))

	memPos := 6
	for _, item := range itemsData.Items {
		encodedData = append(encodedData, byte(item.ItemID>>24), byte(item.ItemID>>16), byte(item.ItemID>>8), byte(item.ItemID))
		memPos += 4
		encodedData = append(encodedData, byte(item.EditableType))
		memPos++
		encodedData = append(encodedData, byte(item.ItemCategory))
		memPos++
		encodedData = append(encodedData, byte(item.ActionType))
		memPos++
		encodedData = append(encodedData, byte(item.HitSoundType))
		memPos++
		encodedData = append(encodedData, byte(len(item.Name)>>8), byte(len(item.Name)))
		memPos += 2
		encodedData = append(encodedData, writeString(item.Name, true, item.ItemID)...)
		memPos += len(item.Name)
		encodedData = append(encodedData, byte(len(item.Texture)>>8), byte(len(item.Texture)))
		memPos += 2
		encodedData = append(encodedData, writeString(item.Texture, false, 0)...)
		memPos += len(item.Texture)
		encodedData = append(encodedData, byte(item.TextureHash>>24), byte(item.TextureHash>>16), byte(item.TextureHash>>8), byte(item.TextureHash))
		memPos += 4
		encodedData = append(encodedData, byte(item.ItemKind))
		memPos++
		encodedData = append(encodedData, byte(item.Val1>>24), byte(item.Val1>>16), byte(item.Val1>>8), byte(item.Val1))
		memPos += 4
		encodedData = append(encodedData, byte(item.TextureX), byte(item.TextureY), byte(item.SpreadType), byte(item.IsStripeyWallpaper), byte(item.CollisionType))
		memPos += 5
		if strings.Contains(item.BreakHits, "r") {
			breakHits, _ := strconv.Atoi(item.BreakHits[:len(item.BreakHits)-1])
			encodedData = append(encodedData, byte(breakHits))
		} else {
			breakHits, _ := strconv.Atoi(item.BreakHits)
			encodedData = append(encodedData, byte(breakHits*6))
		}
		memPos++
		encodedData = append(encodedData, byte(item.DropChance>>24), byte(item.DropChance>>16), byte(item.DropChance>>8), byte(item.DropChance))
		memPos += 4
		encodedData = append(encodedData, byte(item.ClothingType))
		memPos++
		encodedData = append(encodedData, byte(item.Rarity>>8), byte(item.Rarity))
		memPos += 2
		encodedData = append(encodedData, byte(item.MaxAmount))
		memPos++
		encodedData = append(encodedData, byte(len(item.ExtraFile)>>8), byte(len(item.ExtraFile)))
		memPos += 2
		encodedData = append(encodedData, writeString(item.ExtraFile, false, 0)...)
		memPos += len(item.ExtraFile)
		encodedData = append(encodedData, byte(item.ExtraFileHash>>24), byte(item.ExtraFileHash>>16), byte(item.ExtraFileHash>>8), byte(item.ExtraFileHash))
		memPos += 4
		encodedData = append(encodedData, byte(item.AudioVolume>>24), byte(item.AudioVolume>>16), byte(item.AudioVolume>>8), byte(item.AudioVolume))
		memPos += 4
		encodedData = append(encodedData, byte(len(item.PetName)>>8), byte(len(item.PetName)))
		memPos += 2
		encodedData = append(encodedData, writeString(item.PetName, false, 0)...)
		memPos += len(item.PetName)
		encodedData = append(encodedData, byte(len(item.PetPrefix)>>8), byte(len(item.PetPrefix)))
		memPos += 2
		encodedData = append(encodedData, writeString(item.PetPrefix, false, 0)...)
		memPos += len(item.PetPrefix)
		encodedData = append(encodedData, byte(len(item.PetSuffix)>>8), byte(len(item.PetSuffix)))
		memPos += 2
		encodedData = append(encodedData, writeString(item.PetSuffix, false, 0)...)
		memPos += len(item.PetSuffix)
		encodedData = append(encodedData, byte(len(item.PetAbility)>>8), byte(len(item.PetAbility)))
		memPos += 2
		encodedData = append(encodedData, writeString(item.PetAbility, false, 0)...)
		memPos += len(item.PetAbility)
		encodedData = append(encodedData, byte(item.SeedBase), byte(item.SeedOverlay), byte(item.TreeBase), byte(item.TreeLeaves), byte(item.SeedColor.A), byte(item.SeedColor.R), byte(item.SeedColor.G), byte(item.SeedColor.B), byte(item.SeedOverlayColor.A), byte(item.SeedOverlayColor.R), byte(item.SeedOverlayColor.G), byte(item.SeedOverlayColor.B))
		memPos += 12
		encodedData = append(encodedData, byte(0), byte(0), byte(0), byte(0))
		memPos += 4 // skip ingredients
		encodedData = append(encodedData, byte(item.GrowTime>>24), byte(item.GrowTime>>16), byte(item.GrowTime>>8), byte(item.GrowTime))
		memPos += 4
		encodedData = append(encodedData, byte(item.Val2>>8), byte(item.Val2))
		memPos += 2
		encodedData = append(encodedData, byte(item.IsRayman>>8), byte(item.IsRayman))
		memPos += 2
		encodedData = append(encodedData, byte(len(item.ExtraOptions)>>8), byte(len(item.ExtraOptions)))
		memPos += 2
		encodedData = append(encodedData, writeString(item.ExtraOptions, false, 0)...)
		memPos += len(item.ExtraOptions)
		encodedData = append(encodedData, byte(len(item.Texture2)>>8), byte(len(item.Texture2)))
		memPos += 2
		encodedData = append(encodedData, writeString(item.Texture2, false, 0)...)
		memPos += len(item.Texture2)
		encodedData = append(encodedData, byte(len(item.ExtraOptions2)>>8), byte(len(item.ExtraOptions2)))
		memPos += 2
		encodedData = append(encodedData, writeString(item.ExtraOptions2, false, 0)...)
		memPos += len(item.ExtraOptions2)
		encodedData = append(encodedData, fromHexString(item.DataPosition80)...)
		memPos += 80
		if itemsData.Version >= 11 {
			encodedData = append(encodedData, byte(len(item.PunchOptions)>>8), byte(len(item.PunchOptions)))
			memPos += 2
			encodedData = append(encodedData, writeString(item.PunchOptions, false, 0)...)
			memPos += len(item.PunchOptions)
		}
		if itemsData.Version >= 12 {
			encodedData = append(encodedData, fromHexString(item.DataVersion12)...)
			memPos += 13
		}
		if itemsData.Version >= 13 {
			encodedData = append(encodedData, byte(item.IntVersion13>>24), byte(item.IntVersion13>>16), byte(item.IntVersion13>>8), byte(item.IntVersion13))
			memPos += 4
		}
		if itemsData.Version >= 14 {
			encodedData = append(encodedData, byte(item.IntVersion14>>24), byte(item.IntVersion14>>16), byte(item.IntVersion14>>8), byte(item.IntVersion14))
			memPos += 4
		}
		if itemsData.Version >= 15 {
			encodedData = append(encodedData, fromHexString(item.DataVersion15)...)
			memPos += 25
			encodedData = append(encodedData, byte(len(item.StrVersion15)>>8), byte(len(item.StrVersion15)))
			memPos += 2
			encodedData = append(encodedData, writeString(item.StrVersion15, false, 0)...)
			memPos += len(item.StrVersion15)
		}
		if itemsData.Version >= 16 {
			encodedData = append(encodedData, byte(len(item.StrVersion16)>>8), byte(len(item.StrVersion16)))
			memPos += 2
			encodedData = append(encodedData, writeString(item.StrVersion16, false, 0)...)
			memPos += len(item.StrVersion16)
		}
		if itemsData.Version >= 17 {
			encodedData = append(encodedData, byte(item.IntVersion17>>24), byte(item.IntVersion17>>16), byte(item.IntVersion17>>8), byte(item.IntVersion17))
			memPos += 4
		}
		if itemsData.Version >= 18 {
			encodedData = append(encodedData, byte(item.IntVersion18>>24), byte(item.IntVersion18>>16), byte(item.IntVersion18>>8), byte(item.IntVersion18))
			memPos += 4
		}
	}

	return encodedData, nil
}

func writeItemsData(itemsData *ItemsData, filePath string) error {
	if strings.HasSuffix(filePath, ".json") {
		data, err := json.MarshalIndent(itemsData, "", "  ")
		if err != nil {
			return err
		}
		return ioutil.WriteFile(filePath, data, 0644)
	} else if strings.HasSuffix(filePath, ".txt") {
		f, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer f.Close()
		w := bufio.NewWriter(f)
		defer w.Flush()

		fmt.Fprintf(w, "//Credit: IProgramInCPP & GrowtopiaNoobs\n//Format: add_item\\%s\n//NOTE: There are several items, for the breakhits part, add 'r'.\n//Example: 184r\n//What does it mean? So, adding 'r' to breakhits makes it raw breakhits, meaning, if you add 'r' to breakhits, when encoding items.dat, the encoder won't multiply it by 6.\n\nversion\\%d\nitemCount\\%d\n\n", strings.Join(getKeys(itemsData.Items[0]), "\\"), itemsData.Version, itemsData.ItemCount)

		for _, item := range itemsData.Items { // Iterate over Items slice
			fmt.Fprintf(w, "add_item\\%s\n", strings.Join(getValues(item), "\\"))
		}
		return nil
	}

	return fmt.Errorf("unsupported file extension: %s", filePath)
}

func getKeys(item Item) []string {
	t := reflect.TypeOf(item)
	var keys []string
	for i := 0; i < t.NumField(); i++ {
		keys = append(keys, t.Field(i).Name)
	}
	return keys
}

func getValues(item Item) []string {
	var values []string
	for _, k := range getKeys(item) {
		v := reflect.ValueOf(item).FieldByName(k).Interface() // Access value through reflection
		switch v.(type) {
		case int:
			values = append(values, strconv.Itoa(v.(int)))
		case string:
			values = append(values, v.(string))
		case Color:
			values = append(values, fmt.Sprintf("%d,%d,%d,%d", v.(Color).A, v.(Color).R, v.(Color).G, v.(Color).B))
		}
	}
	return values
}

func readString(data []byte, memPos int, usingKey bool, itemID int) string {
	strLen := int(data[memPos])<<8 | int(data[memPos+1])
	if strLen == 0 {
		return ""
	}
	if memPos+2+strLen > len(data) {
		return fmt.Sprintf("Error reading string at position %d: out of bounds", memPos)
	}
	result := make([]byte, strLen)
	copy(result, data[memPos+2:memPos+2+strLen])
	if usingKey {
		for i := 0; i < strLen; i++ {
			// Use modulo to restrict the index to the length of itemsSecretKey
			keyIndex := (i + itemID) % len(itemsSecretKey)
			result[i] ^= itemsSecretKey[keyIndex]
		}
	}
	return string(result)
}

func writeString(str string, usingKey bool, itemID int) []byte {
	result := make([]byte, 0, len(str)+2)
	result = append(result, byte(len(str)>>8), byte(len(str)))
	for i := 0; i < len(str); i++ {
		if usingKey {
			result = append(result, byte(str[i]^itemsSecretKey[((i+itemID)%len(itemsSecretKey))]))
		} else {
			result = append(result, byte(str[i]))
		}
	}
	return result
}

func toHexString(data []byte) string {
	var result []string
	for _, b := range data {
		result = append(result, fmt.Sprintf("%02X", b))
	}
	return strings.Join(result, " ")
}

func fromHexString(hexString string) []byte {
	var result []byte
	hexStrings := strings.Split(hexString, " ")
	for _, hexStr := range hexStrings {
		if len(hexStr) == 0 {
			continue
		}
		n, _ := strconv.ParseInt(hexStr, 16, 64)
		result = append(result, byte(n))
	}
	return result
}
