package toc

// RuneLen возвращает количество символов (rune) в строке s
func RuneLen(s string) int {
	return len([]rune(s))
}
