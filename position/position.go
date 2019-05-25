package position

type Position struct {
	Index     int64 // index
	Height    int64 // height
	Rightness bool  // whether it is a right (or left) child
}
