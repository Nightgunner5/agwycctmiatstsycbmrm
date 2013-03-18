package state

type Item struct {
	// Legendary-quality items have their own names.
	Name string

	Category ItemCategory

	Components *[]Item

	Quality uint16
}

type ItemCategory uint32

const (
	_ ItemCategory = iota << 28
	Metal
	Stone
	Wood
	Organic
	Gem
)

const (
	// (Roughly) in order of value
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
)

const (
	_ ItemCategory = Stone + iota<<20
)

const (
	// (Roughly) in order of value
	_ ItemCategory = Wood + iota<<20
	Birch
	Pine
	Maple
	Walnut
	Larch
	Oak
	Cedar
	Teak
	Mahogany
	Truffula
)

const (
	// Not in any specific order
	_ ItemCategory = Organic + iota<<20
	Feather
	Skin
	Leather
	Fur
	Scale
	Shell
	Bone
	Wool
)

const (
	// (Roughly) in order of value
	_ ItemCategory = Gem + iota<<20
	Jade
	Jasper
	Lapis
	Amethyst
	Emerald
	Sapphire
	Ruby
	Diamond
)

const (
	Material ItemCategory = iota << 18
	Product
)

const (
	// Not in any specific order
	_ ItemCategory = Material + iota<<8
	Scrap
	Ore
	Nugget
	Ingot
	Block
	Sheet
	Thread
	Cloth
	Log
	Plank
	Cluster
)

const (
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

const (
	// Not in any specific order
	_ ItemCategory = Product + iota
)

func (c ItemCategory) String() (s string) {
	switch c & 0x000c0000 {
	case Material:
		switch c & 0x000fff00 {
		case Material:
			defer func() {
				s = s[:len(s)-1]
			}()
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
		case Log:
			s = "log"
		case Plank:
			s = "plank"
		case Cluster:
			s = "cluster"
		default:
			return "ERROR"
		}

	case Product:
		switch {
		case Product:
			s = "widget"
		default:
			return "ERROR"
		}

	default:
		return "ERROR"
	}

	switch c & 0xfff00000 {
	case 0:
		// TODO: error?

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
	case Birch:
		s = "birch " + s
	case Pine:
		s = "pine " + s
	case Maple:
		s = "maple " + s
	case Walnut:
		s = "walnut " + s
	case Larch:
		s = "larch " + s
	case Oak:
		s = "oak " + s
	case Cedar:
		s = "cedar " + s
	case Teak:
		s = "teak " + s
	case Mahogany:
		s = "mahogany " + s
	case Truffula:
		s = "truffula " + s

	case Organic:
		s = "organic " + s
	case Feather:
		s = "feather " + s
	case Skin:
		s = "hide " + s
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

	case Gem:
		s = "jewel " + s
	case Jade:
		s = "jade " + s
	case Jasper:
		s = "jasper " + s
	case Lapis:
		s = "lapis lazuli " + s
	case Amethyst:
		s = "amethyst " + s
	case Emerald:
		s = "emerald " + s
	case Sapphire:
		s = "sapphire " + s
	case Ruby:
		s = "ruby " + s
	case Diamond:
		s = "diamond " + s
	}

	if c&0x000c0000 == Material {
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
	}

	return s
}
