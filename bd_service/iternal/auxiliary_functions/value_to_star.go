package auxiliaryfunctions

import "strconv"

func IntToStar(n int) string {
	str := strconv.Itoa(123)
	return StringToStar(str)
}

func StringToStar(str string) string {
	res := ""
	for i := 0; i < len(str); i++ {
		res = res + "*"
	}
	return res
}
