package db

type Memorybaseddb struct {
	nodes      map[int64][]byte
	leafLength int64
}

func NewMemorybaseddb(nodes map[int64][]byte, leafLength int64) *Memorybaseddb {
	db := Memorybaseddb{nodes: nodes, leafLength: leafLength}
	return &db
}
func (db *Memorybaseddb) Get(index int64) ([]byte, bool) {
	value, ok := db.nodes[index]
	return value, ok
}
func (db *Memorybaseddb) Set(value []byte, index int64) {
	//require correct wordsize
	db.nodes[index] = value
}
func (db *Memorybaseddb) GetLeafLength() int64 {
	return db.leafLength
}
func (db *Memorybaseddb) SetLeafLength(value int64) {
	db.leafLength = value
	// fmt.Print("HEREREEREq", value, db.leafLength)
}
// func (db *Memorybaseddb) GetNodes() map[int64][]byte {
// 	nodes := map[int64][]byte{}
// 	return nodes
// }
