package xstrings

import (
	"bytes"
	"fmt"
)

type Stringer interface {
	String() string
}

func Join(objs []interface{}, sep string) string {
	var buf bytes.Buffer
	var str Stringer
	var ok bool

	for i, obj := range objs {
		if i > 0 {
			buf.WriteString(sep)
		}

		if str, ok = obj.(Stringer); ok {
			buf.WriteString(str.String())
		} else {
			fmt.Fprint(&buf, obj)
		}
	}
	return buf.String()
}
