package alox

import (
	"path"
	"strings"
)

func ShiftHead(value string) (head, tail string) {
	value = path.Clean("/" + value)
	splitIndex := strings.Index(value[1:], "/") + 1

	if splitIndex <= 0 {
		return value[1:], "/"
	}

	return value[1:splitIndex], value[splitIndex:]
}

func ShiftAndAssertHead(value string, assert func(head string) bool) (passed bool, tail string) {
	head, tail := ShiftHead(value)
	return assert(head), tail
}

func ShiftAndMatchHead(value string, head string) (matched bool, tail string) {
	return ShiftAndAssertHead(value, func(actualHead string) bool {
		return actualHead == head
	})
}

func AssertHead(value string, assert func(head string) bool) (passed bool) {
	head, _ := ShiftHead(value)
	return assert(head)
}

func MatchHead(value string, head string) (matched bool) {
	return AssertHead(value, func(actualHead string) bool {
		return actualHead == head
	})
}

func HasPrefix(value string, prefix string) bool {
	return strings.HasPrefix(path.Clean("/"+value), prefix)
}
