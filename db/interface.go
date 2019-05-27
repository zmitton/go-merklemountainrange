package db

const wordsize = 64 // 512 bits support only for now

type Db interface {
	Get(int64) ([]byte, bool)
	Set(int64, []byte)
	GetLeafLength() int64
	SetLeafLength(int64)
}
