package reader

type ReaderResult struct {
	Key   byte
	Value []byte
}

type Reader interface {
	HandleRead(path string) ([]ReaderResult, error)
}

func CreateReader(t string) Reader {
	switch t {
	case "rdb":
		return &RdbReader{}

	default:
		return &RdbReader{}
	}
}
