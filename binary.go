package polar

import (
	"bytes"
	"encoding/binary"
)

type buffer struct {
	bytes.Buffer
}

func (b *buffer) WriteUInt64(i uint64) error {
	return binary.Write(b, binary.BigEndian, i)
}

func (b *buffer) WriteInt64(i int64) error {
	return binary.Write(b, binary.BigEndian, i)
}

func (b *buffer) WriteInt32(i int32) error {
	return binary.Write(b, binary.BigEndian, i)
}

func (b *buffer) WriteInt16(i int16) error {
	return binary.Write(b, binary.BigEndian, i)
}

func (b *buffer) WriteInt8(i int8) error {
	return binary.Write(b, binary.BigEndian, i)
}

func (b *buffer) WriteVarInt(i int32) error {
	//value := int64(i)
	//if (value & (0xFFFFFFFF << 7)) == 0 {
	//	b.Grow(1)
	//	return b.WriteByte(byte(value))
	//} else if (value & (0xFFFFFFFF << 14)) == 0 {
	//	b.Grow(2)
	//	return b.WriteInt16(int16((value&0x7F|0x80)<<8 | (value >> 7)))
	//} else if (value & (0xFFFFFFFF << 21)) == 0 {
	//	b.Grow(3)
	//	return b.WriteBytes([]byte{
	//		byte(value&0x7F | 0x80),
	//		byte((value>>7)&0x7F | 0x80),
	//		byte(value >> 14),
	//	})
	//} else if (value & (0xFFFFFFFF << 28)) == 0 {
	//	b.Grow(4)
	//	//return b.WriteInt32(int32((value & 0x7F | 0x80) << 24 | (((value >> 7) & 0x7F | 0x80) << 16) | ((value >> 14) & 0x7F | 0x80) << 8 | (value >> 21)))
	//	return b.WriteBytes([]byte{
	//		byte(value&0x7F | 0x80),
	//		byte((value>>7)&0x7F | 0x80),
	//		byte((value>>14)&0x7F | 0x80),
	//		byte(value >> 21),
	//	})
	//} else {
	//	b.Grow(5)
	//	return b.WriteBytes([]byte{
	//		byte(value&0x7F | 0x80),
	//		byte((value>>7)&0x7F | 0x80),
	//		byte((value>>14)&0x7F | 0x80),
	//		byte((value>>21)&0x7F | 0x80),
	//		byte(value >> 28),
	//	})
	//}

	var buf [binary.MaxVarintLen32]byte
	n := binary.PutUvarint(buf[:], uint64(uint32(i)))
	_, err := b.Write(buf[:n])
	return err

	// // what even is going on here???
	//	b.Grow(binary.MaxVarintLen32)
	//	temp := make([]byte, binary.MaxVarintLen32)
	//	n := binary.PutVarint(temp, int64(i))
	//	_, err := b.Write(temp[:n])
	//	return err
}

func (b *buffer) WriteBool(v bool) error {
	return binary.Write(b, binary.BigEndian, v)
}

func (b *buffer) WriteBytes(v []byte) error {
	_, err := b.Write(v)
	return err
}

func (b *buffer) WriteString(s string) error {
	rawString := []byte(s)
	if err := b.WriteVarInt(int32(len(rawString))); err != nil {
		return err
	}
	_, err := b.Write(rawString)
	return err
}
