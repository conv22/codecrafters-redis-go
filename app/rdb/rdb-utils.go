package rdb

import (
	"encoding/binary"
	"errors"
)

func (rdb *Rdb) readBytes(l int) ([]byte, error) {
	buffer := make([]byte, l)

	length, err := rdb.reader.Read(buffer)

	if err != nil {
		return nil, err
	}

	if length != l {
		return nil, errors.New("invalid encoding")
	}

	return buffer, nil

}

func (rdb *Rdb) readByte() (byte, error) {
	bytes, err := rdb.readBytes(1)
	if err != nil {
		return 0, err
	}
	return bytes[0], nil
}

func (rdb *Rdb) readUnsignedShort() (uint16, error) {
	b, err := rdb.readBytes(2)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(b), nil
}

func (rdb *Rdb) readSignedInt() (int32, error) {
	b, err := rdb.readBytes(4)
	if err != nil {
		return 0, err
	}
	return int32(binary.LittleEndian.Uint32(b)), nil
}

func (rdb *Rdb) readUnsignedInt() (uint32, error) {
	b, err := rdb.readBytes(4)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(b), nil
}
