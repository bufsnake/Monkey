package other

func Exist(source []string, data string) bool {
	for i := 0; i < len(source); i++ {
		if source[i] == data {
			return true
		}
	}
	return false
}
