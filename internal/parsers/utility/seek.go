package utility

func Seek(data *[]byte, current int, symbolByte byte) int {
	dataLength := len(*data)
	for current < dataLength && (*data)[current] != symbolByte {
		current++
	}

	if current == dataLength {
		return -1
	}

	return current
}
