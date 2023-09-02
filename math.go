package polar

type ChunkIndex int64

func (c ChunkIndex) X() int {
	return int(int32(c >> 32))
}

func (c ChunkIndex) Z() int {
	return int(int32(c))
}

func ChunkIndexFromXZ(x, z int) ChunkIndex {
	return ChunkIndex((int64(x) << 32) | int64(z&0xFFFFFFFF))
}

func ChunkPosToXYZ(chunkPos int32) (int, int, int) {
	y := (chunkPos & 0x07FFFFF0) >> 4
	if ((chunkPos >> 27) & 1) == 1 {
		y = -y // Sign bit set, invert sign
	}
	return int(chunkPos & 0xF), int(y), int((chunkPos >> 28) & 0xF)
}

func ChunkPosFromXYZ(x, y, z int) int32 {
	return int32((x & 0xF) | ((y & 0x7FFFF) << 4) | ((z & 0xF) << 28))
}
