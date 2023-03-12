package filter

import (
	"net/http"
	"path"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"alox.sh"
)

func AssertPath(assert func(string) bool) alox.Filter {
	return func(request *http.Request) bool {
		return assert(path.Clean("/" + request.URL.Path)[1:])
	}
}

func AssertPathSegments(assert func([]string) bool) alox.Filter {
	return AssertPath(func(path string) bool {
		segments := []string{}

		if path != "" {
			segments = strings.Split(path, "/")
		}

		return assert(segments)
	})
}

func Path() pathFilter {
	return func(assert func(string) bool) alox.Filter {
		return AssertPath(assert)
	}
}

func PathSegments() pathSegmentsFilter {
	return func(assert func([]string) bool) alox.Filter {
		return AssertPathSegments(assert)
	}
}

type pathFilter func(func(string) bool) alox.Filter

func (filter pathFilter) Segments() pathSegmentsFilter {
	return PathSegments()
}

func (filter pathFilter) IsRoot() alox.Filter {
	return AssertPath(func(path string) bool {
		return path == "" || path == "/"
	})
}

func (filter pathFilter) MatchPattern(pattern string, params ...interface{}) alox.Filter {
	return filter(func(value string) bool {
		for ; pattern != "" && value != ""; pattern = pattern[1:] {
			switch pattern[0] {
			case '+':
				// '+' matches till next slash in path
				slash := strings.IndexByte(value, '/')
				if slash < 0 {
					slash = len(value)
				}

				segment := value[:slash]
				value = value[slash:]

				switch param := params[0].(type) {
				case *string:
					*param = segment
				case *int:
					n, err := strconv.Atoi(segment)
					if err != nil {
						return false
					}

					*param = n
				case *float64:
					n, err := strconv.ParseFloat(segment, 64)
					if err != nil {
						return false
					}

					*param = n
				case *primitive.ObjectID:
					objectID, err := primitive.ObjectIDFromHex(segment)
					if err != nil || objectID.IsZero() {
						return false
					}

					*param = objectID
				default:
					panic("params must be of type: *string, *int, *float64 or *(go.mongodb.org/mongo-driver/bson/primitive).ObjectID")
				}

				params = params[1:]
			case value[0]:
				// non-'+' pattern byte must match path byte
				value = value[1:]
			default:
				return false
			}
		}

		return value == "" && pattern == ""
	})
}

type pathSegmentsFilter func(func([]string) bool) alox.Filter

func (filter pathSegmentsFilter) LenEq(count int) alox.Filter {
	return filter(func(segments []string) bool { return len(segments) == count })
}

func (filter pathSegmentsFilter) LenGt(count int) alox.Filter {
	return filter(func(segments []string) bool { return len(segments) > count })
}

func (filter pathSegmentsFilter) LenGte(count int) alox.Filter {
	return filter(func(segments []string) bool { return len(segments) >= count })
}

func (filter pathSegmentsFilter) LenLt(count int) alox.Filter {
	return filter(func(segments []string) bool { return len(segments) < count })
}

func (filter pathSegmentsFilter) LenLte(count int) alox.Filter {
	return filter(func(segments []string) bool { return len(segments) <= count })
}
