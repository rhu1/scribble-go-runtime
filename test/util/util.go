package util

const LOCALHOST = "127.0.0.1"

func Copy(data int, i int) int {
	return data
}

// No generics or polymorphism
func CopyString(data string, i int) string {
	return data
}

func CopyBates(data []byte, i int) []byte {
	bs := make([]byte, len(data))
	copy(bs, data)
	return bs
}

func Sum(xs []int) int {
	res := 0
	for i := 0; i < len(xs); i++ {
		res = res + xs[i]
	}
	return res
}

func UnaryReduce(xs []int) int {
	return xs[0]
}

func UnaryReduceString(xs []string) string {
	return xs[0]
}

func UnaryReduceBates(xs [][]byte) []byte {
	return xs[0]
}
