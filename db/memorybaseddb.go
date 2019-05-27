package db

type Memorybaseddb struct {
	nodes      map[int64][]byte
	leafLength int64
}

func NewMemorybaseddb(nodes map[int64][]byte, leafLength int64) *Memorybaseddb {
	db := Memorybaseddb{nodes: nodes, leafLength: leafLength}
	return &db
}
func (db Memorybaseddb) Get(index int64) ([]byte, bool) {
	value, ok := db.nodes[index]
	return value, ok
}
func (db Memorybaseddb) Set(index int64, value []byte) {
	db.nodes[index] = value
}
func (db Memorybaseddb) GetLeafLength() int64 {
	return db.leafLength
}
func (db Memorybaseddb) SetLeafLength(_value int64) {
	db.leafLength = _value
}
