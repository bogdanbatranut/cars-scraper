package utils

func ToIntPointer(value int) *int {
	return &value
}

func ToStringPointer(value string) *string {
	return &value
}

func InArrayStr(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}
