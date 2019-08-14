package db

const wordsize = 64 // 512 bits support only for now

type Db interface {
	Get(int64) ([]byte, bool)
	Set([]byte, int64)
	GetLeafLength() int64
	SetLeafLength(int64)
}

// type Db64 interface {
// 	Get(int64) ([64]byte, bool)
// 	Set([64]byte, int64)
// 	GetLeafLength() int64
// 	SetLeafLength(int64)
// }

// type Db32 interface {
// 	Get(int64) ([32]byte, bool)
// 	Set([32]byte, int64)
// 	GetLeafLength() int64
// 	SetLeafLength(int64)
// }
