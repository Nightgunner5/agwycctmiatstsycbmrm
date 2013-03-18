package state

type Item struct {
	// Legendary-quality items have their own names.
	Name string

	// Materials cannot be sold by the player.
	Material bool

	Category ItemCategory
}

type ItemCategory uint32

const (
	_ ItemCategory = iota << 28
	Metal
	Stone
	Wood
	Organic

	_ ItemCategory = Metal + iota<<20
	Copper
	Tin
	Bronze
	Lead
	Iron
	Steel
	Cobalt
	Silver
	Gold
	Platinum
	Titanium
	Adamantite

	_ ItemCategory = Stone + iota<<20

	_ ItemCategory = Wood + iota<<20

	_ ItemCategory = Organic + iota<<20
	Feather
	Skin
	Leather
	Fur
	Scale
	Shell
	Bone
	Wool

	_ ItemCategory = iota << 8
	Scrap
	Ore
	Nugget
	Ingot
	Block
	Sheet
	Thread
	Cloth

	// Sizes
	Miniscule ItemCategory = 0
	Tiny      ItemCategory = 12
	Small     ItemCategory = 24
	Medium    ItemCategory = 48
	Big       ItemCategory = 96
	Large     ItemCategory = 128
	Huge      ItemCategory = 192
	Gigantic  ItemCategory = 250
)

func (c ItemCategory) String() string {
	var s string

	switch c & 0x000fff00 {
	case 0:
		s = "widget"
	case Scrap:
		s = "scrap"
	case Ore:
		s = "ore"
	case Nugget:
		s = "nugget"
	case Ingot:
		s = "ingot"
	case Block:
		s = "block"
	case Sheet:
		s = "sheet"
	case Thread:
		s = "thread"
	case Cloth:
		s = "cloth"
	default:
		return "ERROR"
	}

	switch c & 0xfff00000 {
	case 0:
		// no prefix
	case Metal:
		s = "metal " + s
	case Copper:
		s = "copper " + s
	case Tin:
		s = "tin " + s
	case Bronze:
		s = "bronze " + s
	case Lead:
		s = "lead " + s
	case Iron:
		s = "iron " + s
	case Steel:
		s = "steel " + s
	case Cobalt:
		s = "cobalt " + s
	case Silver:
		s = "silver " + s
	case Gold:
		s = "gold " + s
	case Platinum:
		s = "platinum " + s
	case Titanium:
		s = "titanium " + s
	case Adamantite:
		s = "adamantite " + s

	case Stone:
		s = "stone " + s

	case Wood:
		s = "wooden " + s

	case Organic:
		s = "organic " + s
	case Feather:
		s = "feather " + s
	case Skin:
		s = "skin " + s
	case Leather:
		s = "leather " + s
	case Fur:
		s = "fur " + s
	case Scale:
		s = "scale " + s
	case Shell:
		s = "shell " + s
	case Bone:
		s = "bone " + s
	}

	switch size := c & 0x000000ff; {
	case size >= Gigantic:
		s = "gigantic " + s
	case size >= Huge:
		s = "huge " + s
	case size >= Large:
		s = "large " + s
	case size >= Big:
		s = "big " + s
	case size >= Medium:
		s = "medium " + s
	case size >= Small:
		s = "small " + s
	case size >= Tiny:
		s = "tiny " + s
	case size >= Miniscule:
		s = "miniscule " + s
	}

	return s
}