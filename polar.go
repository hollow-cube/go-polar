package polar

import "fmt"

const MagicNumber = 0x506F6C72 // `Polr`

const LatestVersion int16 = 3

type Compression int8

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

func (w *World) SectionCount() int {
	return int(w.MaxSection - w.MinSection + 1)
}

func (w *World) GetChunk(x, z int) *Chunk {
	if w.chunks == nil {
		return nil
	}
	idx := ChunkIndexFromXZ(x, z)
	return w.chunks[idx]
}

func (w *World) SetChunk(chunk *Chunk) {
	if w.chunks == nil {
		w.chunks = make(map[ChunkIndex]*Chunk)
	}

	idx := ChunkIndexFromXZ(int(chunk.X), int(chunk.Z))
	if _, ok := w.chunks[idx]; ok {
		panic(fmt.Sprintf("chunk %d, %d already exists", chunk.X, chunk.Z))
	}
	w.chunks[idx] = chunk
}

type Chunk struct {
	X int32
	Z int32

	Sections      []*Section
	BlockEntities []*BlockEntity

	//todo heightmaps

	//todo user data
}

type Section struct {
	BlockPalette []string
	BlockStates  []uint64

	BiomePalette []string
	BiomeStates  []uint64

	BlockLight []byte // 2048 bytes, or nil
	SkyLight   []byte // 2048 bytes, or nil
}

func (s *Section) IsEmpty() bool {
	return len(s.BlockPalette) == 0 && len(s.BlockStates) == 0 && s.BlockLight == nil && s.SkyLight == nil
}

type BlockEntity struct {
	ChunkPos int32
	ID       string                 // Can be empty string if missing
	Data     map[string]interface{} // NBT data, or nil if not present. todo type here
}

func (b *BlockEntity) Pos() (int, int, int) {
	return ChunkPosToXYZ(b.ChunkPos)
}
