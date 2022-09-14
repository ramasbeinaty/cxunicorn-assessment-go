package utils

import (
	"bytes"
	"database/sql/driver"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// ErrParseData signifies that an error occured while parsing SQL data
var ErrParseData = errors.New("unable to parse SQL data")

// TimeArray wraps a time.Time slice to be used as a Postgres array
type TimeArray []time.Time

// type PgTimeArray []pq.NullTime

// Value implements the driver.Valuer interface.
func (a TimeArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	if n := len(a); n > 0 {
		// There will be at least two curly brackets, N bytes of values,
		// and N-1 bytes of delimiters.
		b := make([]byte, 1, 1+2*n)
		b[0] = '{'

		b = strconv.AppendQuote(b, a[0].Format(os.Getenv("TIME_FORMAT")))
		for i := 1; i < n; i++ {
			b = append(b, ',')
			b = strconv.AppendQuote(b, a[i].Format(os.Getenv("TIME_FORMAT")))
		}

		return string(append(b, '}')), nil
	}

	return "{}", nil
}

// Scan implements the sql.Scanner interface
func (a *TimeArray) Scan(src interface{}) error {
	// return pq.GenericArray{A: a}.Scan(src)

	switch src := src.(type) {
	case []byte:
		return a.scanBytes(src)
	case time.Time:
		return a.scanBytes([]byte(os.Getenv("TIME_FORMAT")))
	case nil:
		*a = nil
		return nil
	}
	return fmt.Errorf("pq: cannot convert %T to TimeArray", src)

}

func (a *TimeArray) scanBytes(src []byte) error {
	elems, err := scanLinearArray(src, []byte{','}, "TimeArray")
	if err != nil {
		return err
	}
	if *a != nil && len(elems) == 0 {
		*a = (*a)[:0]
	} else {
		b := make(TimeArray, len(elems))
		for i, v := range elems {
			if b[i], err = time.Parse(os.Getenv("TIME_FORMAT"), string(v)); v == nil || err != nil {
				fmt.Println("ERROR: ", err)
				return fmt.Errorf("pq: parsing array element index %d: cannot convert nil to string", i)
			}
		}
		*a = b
	}
	return nil
}

// parseArray extracts the dimensions and elements of an array represented in
// text format. Only representations emitted by the backend are supported.
// Notably, whitespace around brackets and delimiters is significant, and NULL
// is case-sensitive.
//
// See http://www.postgresql.org/docs/current/static/arrays.html#ARRAYS-IO
func parseArray(src, del []byte) (dims []int, elems [][]byte, err error) {
	var depth, i int

	if len(src) < 1 || src[0] != '{' {
		return nil, nil, fmt.Errorf("pq: unable to parse array; expected %q at offset %d", '{', 0)
	}

Open:
	for i < len(src) {
		switch src[i] {
		case '{':
			depth++
			i++
		case '}':
			elems = make([][]byte, 0)
			goto Close
		default:
			break Open
		}
	}
	dims = make([]int, i)

Element:
	for i < len(src) {
		switch src[i] {
		case '{':
			if depth == len(dims) {
				break Element
			}
			depth++
			dims[depth-1] = 0
			i++
		case '"':
			var elem = []byte{}
			var escape bool
			for i++; i < len(src); i++ {
				if escape {
					elem = append(elem, src[i])
					escape = false
				} else {
					switch src[i] {
					default:
						elem = append(elem, src[i])
					case '\\':
						escape = true
					case '"':
						elems = append(elems, elem)
						i++
						break Element
					}
				}
			}
		default:
			for start := i; i < len(src); i++ {
				if bytes.HasPrefix(src[i:], del) || src[i] == '}' {
					elem := src[start:i]
					if len(elem) == 0 {
						return nil, nil, fmt.Errorf("pq: unable to parse array; unexpected %q at offset %d", src[i], i)
					}
					if bytes.Equal(elem, []byte("NULL")) {
						elem = nil
					}
					elems = append(elems, elem)
					break Element
				}
			}
		}
	}

	for i < len(src) {
		if bytes.HasPrefix(src[i:], del) && depth > 0 {
			dims[depth-1]++
			i += len(del)
			goto Element
		} else if src[i] == '}' && depth > 0 {
			dims[depth-1]++
			depth--
			i++
		} else {
			return nil, nil, fmt.Errorf("pq: unable to parse array; unexpected %q at offset %d", src[i], i)
		}
	}

Close:
	for i < len(src) {
		if src[i] == '}' && depth > 0 {
			depth--
			i++
		} else {
			return nil, nil, fmt.Errorf("pq: unable to parse array; unexpected %q at offset %d", src[i], i)
		}
	}
	if depth > 0 {
		err = fmt.Errorf("pq: unable to parse array; expected %q at offset %d", '}', i)
	}
	if err == nil {
		for _, d := range dims {
			if (len(elems) % d) != 0 {
				err = fmt.Errorf("pq: multidimensional arrays must have elements with matching dimensions")
			}
		}
	}
	return
}

func scanLinearArray(src, del []byte, typ string) (elems [][]byte, err error) {
	dims, elems, err := parseArray(src, del)
	if err != nil {
		return nil, err
	}
	if len(dims) > 1 {
		return nil, fmt.Errorf("pq: cannot convert ARRAY%s to %s", strings.Replace(fmt.Sprint(dims), " ", "][", -1), typ)
	}
	return elems, err
}

// Value implements the driver.Valuer interface
// func (a *TimeArray) Value() (driver.Value, error) { return pq.GenericArray{A: a}.Value() }
