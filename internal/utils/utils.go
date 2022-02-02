package utils

func WithPrefix(prefix string, arr []string) []string {
	tmp := make([]string, len(arr))
	for i := 0; i < len(arr); i++ {
		tmp[i] = prefix + arr[i]
	}
	return tmp
}
