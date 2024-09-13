package gui

import (
	"slices"
	"strconv"
	"strings"
)

type Option func(string) string

// FormatInt форматирует число для вывода в графический интерфейс
func FormatInt(i int, opts ...Option) string {
	str := strconv.Itoa(i)
	str = addSpaces(str, 3)

	for _, opt := range opts {
		str = opt(str)
	}

	return str
}

func addMinus(str string) string {
	if str == "0" {
		return str
	}

	return "-" + str
}

// addSpaces добавляет пробелы каждые n символов начиная с конца строки
func addSpaces(str string, n int) string {
	length := len(str)
	nums := make([]string, 0)

	for i := 1; i < length; i++ {
		idx := length - i
		if i%n == 0 {
			nums = append(nums, str[idx:])
			str = str[:idx]
		}
	}

	nums = append(nums, str)
	slices.Reverse(nums)

	return strings.Join(nums, " ")
}
