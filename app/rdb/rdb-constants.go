package rdb

const (
	// Op codes
	RDB_OPCODE_EOF            = 0xFF
	RDB_OPCODE_SELECT_DB      = 0xFE
	RDB_OPCODE_EXPIRE_TIME    = 0xFD
	RDB_OPCODE_EXPIRE_TIME_MS = 0xFC
	RDB_OPCODE_RESIZE_DB      = 0xFB
	RDB_OPCODE_AUX            = 0xFA

	// Encoding types
	RDB_ENCODING_STRING_ENCODING             = 0
	RDB_ENCODING_LIST_ENCODING               = 1
	RDB_ENCODING_SET_ENCODING                = 2
	RDB_ENCODING_SORTED_SET_ENCODING         = 3
	RDB_ENCODING_HASH_ENCODING               = 4
	RDB_ENCODING_ZIPMAP_ENCODING             = 9
	RDB_ENCODING_ZIPLIST_ENCODING            = 10
	RDB_ENCODING_INTSET_ENCODING             = 11
	RDB_ENCODING_SORTED_SET_ZIPLIST_ENCODING = 12
	RDB_ENCODING_HASHMAP_ZIPLIST_ENCODING    = 13
	RDB_ENCODING_LIST_QUICKLIST_ENCODING     = 14

	// Length constants
	RDB_VERSION_NUMBER_LENGTH        = 4
	RDB_EXPIRE_TIME_MS_BYTES_LENGTH  = 8
	RDB_EXPIRE_TIME_SEC_BYTES_LENGTH = 4

	// Magic redis string
	RDB_MAGIC = "REDIS"
)

var SUPPORTED_ENCODINGS = []byte{RDB_ENCODING_STRING_ENCODING}
