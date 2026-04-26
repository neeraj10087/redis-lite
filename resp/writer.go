package resp

import "fmt"

func (v Value) Marshal() []byte {
	switch v.Typ {
		case "string":
			return []byte(fmt.Sprintf("+%s\r\n", v.Str))
		case "error":
			return []byte(fmt.Sprintf("-%s\r\n", v.Str))
		case "integer":
			return []byte(fmt.Sprintf(":%d\r\n", v.Num))
		case "bulk":
			return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(v.Bulk), v.Bulk))
		case "null":
			return []byte("$-1\r\n")
		case "array":
			out := []byte(fmt.Sprintf("*%d\r\n", len(v.Array)))
			for _, elem := range v.Array {
				out = append(out, elem.Marshal()...)
			}
			return out
		default:
			return []byte("-ERR unknown type\r\n")
	}
}