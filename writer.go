package polar

import (
	"fmt"
	"github.com/Tnze/go-mc/nbt"
	"github.com/klauspost/compress/zstd"
)

func Write(world *World) ([]byte, error) {
	var worldDataBuffer buffer
	if err := writeWorld(&worldDataBuffer, world); err != nil {
		return nil, err
	}

	var wrapBuffer buffer
	if err := wrapBuffer.WriteInt32(MagicNumber); err != nil {
		return nil, err
	}
	if err := wrapBuffer.WriteInt16(world.Version); err != nil {
		return nil, err
	}
	if err := wrapBuffer.WriteInt8(int8(world.Compression)); err != nil {
		return nil, err
	}
	if err := wrapBuffer.WriteVarInt(int32(worldDataBuffer.Len())); err != nil {
		return nil, err
	}

	switch world.Compression {
	case CompressionNone:
		if err := wrapBuffer.WriteBytes(worldDataBuffer.Bytes()); err != nil {
			return nil, err
		}
	case CompressionZstd:
		w, _ := zstd.NewWriter(&wrapBuffer)
		if _, err := w.Write(worldDataBuffer.Bytes()); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported compression type %d", world.Compression)
	}

	return wrapBuffer.Bytes(), nil
}

func writeWorld(b *buffer, world *World) (err error) {
	if err = b.WriteInt8(world.MinSection); err != nil {
		return err
	}
	if err = b.WriteInt8(world.MaxSection); err != nil {
		return err
	}
	if err = b.WriteVarInt(int32(len(world.chunks))); err != nil {
		return err
	}
	for _, chunk := range world.chunks {
		if err = writeChunk(b, chunk, world.SectionCount()); err != nil {
			return fmt.Errorf("failed to write chunk %d, %d: %w", chunk.X, chunk.Z, err)
		}
	}

	return nil
}

func writeChunk(b *buffer, chunk *Chunk, sectionCount int) (err error) {
	if err = b.WriteVarInt(chunk.X); err != nil {
		return err
	}
	if err = b.WriteVarInt(chunk.Z); err != nil {
		return err
	}

	// Section data (padded with empty if missing)
	for i := 0; i < sectionCount; i++ {
		var section *Section
		if i < len(chunk.Sections) {
			section = chunk.Sections[i]
		} else {
			section = &Section{}
		}

		if err = writeSection(b, section); err != nil {
			return fmt.Errorf("failed to write section %d (abs): %w", i, err)
		}
	}

	// Block entities
	if err = b.WriteVarInt(int32(len(chunk.BlockEntities))); err != nil {
		return err
	}
	for _, blockEntity := range chunk.BlockEntities {
		if err = writeBlockEntity(b, blockEntity); err != nil {
			x, y, z := blockEntity.Pos()
			return fmt.Errorf("failed to write block entity %d, %d, %d: %w", x, y, z, err)
		}
	}

	// Heightmaps todo
	if err = b.WriteInt32(0); err != nil {
		return err
	}

	// User data todo
	if err = b.WriteVarInt(0); err != nil {
		return err
	}

	return nil
}

func writeSection(b *buffer, section *Section) (err error) {
	isEmpty := section.IsEmpty()
	if err = b.WriteBool(isEmpty); err != nil {
		return err
	}
	if isEmpty {
		return nil
	}

	// Block palette
	if err = b.WriteVarInt(int32(len(section.BlockPalette))); err != nil {
		return err
	}
	for _, block := range section.BlockPalette {
		if err = b.WriteString(block); err != nil {
			return err
		}
	}
	if len(section.BlockPalette) > 1 {
		if err = b.WriteVarInt(int32(len(section.BlockStates))); err != nil {
			return err
		}
		for _, state := range section.BlockStates {
			if err = b.WriteUInt64(state); err != nil {
				return err
			}
		}
	}

	// Biome palette
	if err = b.WriteVarInt(int32(len(section.BiomePalette))); err != nil {
		return err
	}
	for _, biome := range section.BiomePalette {
		if err = b.WriteString(biome); err != nil {
			return err
		}
	}
	if len(section.BiomePalette) > 1 {
		if err = b.WriteVarInt(int32(len(section.BiomeStates))); err != nil {
			return err
		}
		for _, state := range section.BiomeStates {
			if err = b.WriteUInt64(state); err != nil {
				return err
			}
		}
	}

	// Light data
	if section.BlockLight != nil {
		if err = b.WriteBytes(section.BlockLight); err != nil {
			return err
		}
	}
	if section.SkyLight != nil {
		if err = b.WriteBytes(section.SkyLight); err != nil {
			return err
		}
	}

	return nil
}

func writeBlockEntity(b *buffer, blockEntity *BlockEntity) (err error) {
	if err = b.WriteInt32(blockEntity.ChunkPos); err != nil {
		return err
	}
	hasBlockEntityId := blockEntity.ID != ""
	if err = b.WriteBool(hasBlockEntityId); err != nil {
		return err
	}
	if hasBlockEntityId {
		if err = b.WriteString(blockEntity.ID); err != nil {
			return err
		}
	}

	hasBlockEntityData := blockEntity.Data != nil
	if err = b.WriteBool(hasBlockEntityData); err != nil {
		return err
	}
	if hasBlockEntityData {
		if err = nbt.NewEncoder(b).Encode(blockEntity.Data, ""); err != nil {
			return fmt.Errorf("failed to marshal block entity data: %w", err)
		}
	}

	return nil
}
