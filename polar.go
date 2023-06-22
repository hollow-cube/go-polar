package polar

const MagicNumber = 0x506F6C72 // `Polr`

const LatestVersion int16 = 2

type Compression int

const (
	CompressionNone Compression = iota
	CompressionZstd
)

type World struct {
	Version     int16
	Compression Compression

	MinSection int8
	MaxSection int8
	chunks     map[ChunkIndex]*Chunk
}

type Chunk struct {
	X int32
	Z int32

	Sections      []Section
	BlockEntities []BlockEntity

	//todo heightmaps
}

type Section struct {
	Empty bool

	BlockPalette []string
	BlockStates  []uint64

	BiomePalette []string
	BiomeStates  []uint64

	BlockLight []byte // 2048 bytes, or nil
	SkyLight   []byte // 2048 bytes, or nil
}

type BlockEntity struct {
	ChunkPos int32
	ID       string // Can be empty string if missing
	Data     []byte // NBT data, or nil if not present
}

func (b *BlockEntity) Pos() (int, int, int) {
	return ChunkPosToXYZ(b.ChunkPos)
}
