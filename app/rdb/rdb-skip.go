package rdb

import "errors"

func (rdb *Rdb) skipFloat() error {
	length, err := rdb.readByte()

	if err != nil {
		return err
	}
	if length < 253 {
		_, err := rdb.readBytes(int(length))

		return err
	}

	return nil

}

func (rdb *Rdb) skipBinaryDouble() error {
	_, err := rdb.readBytes(8)
	return err
}

func (rdb *Rdb) skipString() (bytesToSkip int, err error) {
	length, isEncoded, err := rdb.parseLength()

	if err != nil {
		return
	}

	if isEncoded {
		// https://rdb.fnordig.de/file_format.html#string-encoding
		switch length {
		// *Integers as String
		// indicates that an 8 bit integer follows
		case 0:
			bytesToSkip = 1
		// indicates that a 16 bit integer follows
		case 1:
			bytesToSkip = 2
		// indicates that a 32 bit integer follows
		case 2:
			bytesToSkip = 4
		// *Compressed Strings
		case 3:
			clenLength, _, err := rdb.parseLength()

			if err != nil {
				return 0, err
			}

			_, _, err = rdb.parseLength()

			if err != nil {
				return 0, err
			}

			bytesToSkip = clenLength
		default:
			return 0, errors.New("unsupported encoding")
		}

	} else {
		bytesToSkip = length

	}

	return bytesToSkip, nil

}

func (rdb *Rdb) skipObject(encType byte) (skipStrings int, err error) {
	switch encType {
	case RDB_ENCODING_STRING_ENCODING:
	case RDB_ENCODING_ZIPMAP_ENCODING:
	case RDB_ENCODING_ZIPLIST_ENCODING:
	case RDB_ENCODING_HASHMAP_ZIPLIST_ENCODING:
	case RDB_ENCODING_INTSET_ENCODING:
		skipStrings = 1
	case RDB_ENCODING_LIST_ENCODING, RDB_ENCODING_SET_ENCODING:
	case RDB_ENCODING_LIST_QUICKLIST_ENCODING:
		skipStrings, _, err = rdb.parseLength()
	case RDB_ENCODING_SORTED_SET_ENCODING:
	case RDB_ENCODING_SORTED_SET_ZIPLIST_ENCODING:
		length, _, err := rdb.parseLength()
		if err != nil {
			return 0, err
		}
		for x := int(0); x < length; x++ {
			if encType == RDB_ENCODING_SORTED_SET_ZIPLIST_ENCODING {
				if _, err := rdb.skipString(); err != nil {
					return 0, err
				}
				if err := rdb.skipBinaryDouble(); err != nil {
					return 0, err
				}
			} else {
				if _, err := rdb.skipString(); err != nil {
					return 0, err
				}
				if err := rdb.skipFloat(); err != nil {
					return 0, err
				}
			}
		}
	case RDB_ENCODING_HASH_ENCODING:
		skipStrings, _, err = rdb.parseLength()
		if err != nil {
			return
		}
		skipStrings *= 2

	default:
		err = errors.New("invalid object type")
		return
	}
	return
}
