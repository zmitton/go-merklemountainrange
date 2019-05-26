package db

type Db interface {
	Get(int64) ([]byte, bool)
	Set(int64, []byte)
	GetLeafLength() int64
	SetLeafLength(int64)
}
