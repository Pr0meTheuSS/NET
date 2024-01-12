package geometry

type Position struct {
	X, Y int32
}

func Find(list []Position, elem Position) int {
	for i, e := range list {
		if e == elem {
			return i
		}
	}
	return -1
}
