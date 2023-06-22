package polar

type ChunkIndex int64

func ChunkPosToXYZ(chunkPos int32) (int, int, int) {
	y := (chunkPos & 0x07FFFFF0) >> 4
	if ((chunkPos >> 27) & 1) == 1 {
		y = -y // Sign bit set, invert sign
	}
	return int(chunkPos & 0xF), int(y), int((chunkPos >> 28) & 0xF)
}
