package reader

type ReaderResult struct {
	Key   string
	Value string
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
