package tool

import "strconv"

func UintToString(a uint) string {
	return strconv.Itoa(int(a))
}
func StringToUint(s string) (uint, error) {
	result, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return uint(result), nil
}

// CompareString 比较两个字符串的大小
// max min
func CompareString(a, b string) (int, int) { //
	aInt, err := strconv.Atoi(a)
	if err != nil {
		return 0, 0
	}
	bInt, err := strconv.Atoi(b)
	if err != nil {
		return 0, 0
	}
	if aInt > bInt {
		return aInt, bInt
	} else {
		return bInt, aInt
	}
}
