package resp

import "bufio"
import "strconv"

// # Prefix | Type            | Format Example                  | Python Type
// # -------|-----------------|---------------------------------|-------------
// # +      | Simple String   | +OK\r\n                         | str
// # -      | Error           | -ERR error msg\r\n              | Exception / str
// # :      | Integer         | :1000\r\n                       | int
// # $      | Bulk String     | $6\r\nfoobar\r\n                | str
// # $      | Null Bulk       | $-1\r\n                         | None
// # *      | Array           | *2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n | list
// # *      | Null Array      | *-1\r\n                         | None

type Value struct {
	Typ   string
	Str   string
	Num   int
	Bulk  string
	Array []Value
}

func Parse(r *bufio.Reader) (Value, error) {
	dataType , err := r.ReadByte()
	if err != nil {
		return Value{} , err
	}

	switch dataType {
      case '+':
          return parseSimpleString(r)
      case '-':
          return parseError(r)
      case ':':
          return parseInteger(r)
      case '$':
          return parseBulkString(r)
      case '*':
          return parseArray(r)
      default:
          return Value{}, nil
	}
}

func parseError(r *bufio.Reader) (Value, error) {
	line, err := readLine(r)
	if err != nil {
		return Value{}, err
	}
	return Value{Typ: "error", Str: line}, nil
}

func parseInteger(r *bufio.Reader) (Value, error) {
	line, err := readLine(r)
	if err != nil {
		return Value{}, err
	}
	num, err := strconv.Atoi(line)
	if err != nil {
		return Value{}, err
	}
	return Value{Typ: "integer", Num: num}, nil
}

func parseBulkString(r *bufio.Reader) (Value, error) {
	line, err := readLine(r);
	if err != nil {
		return Value{}, err
	}

	length,err := strconv.Atoi(line);
	if err != nil {
		return Value{},err
	}

	if length == -1 {
		return Value{Typ:"null"},nil
	}
	buf := make([]byte, length + 2)
	r.Read(buf)

	return Value{Typ: "bulk", Bulk: string(buf[:length])}, nil 

}

func parseArray(r *bufio.Reader) (Value, error) {
	line, err := readLine(r)
	if err != nil {
		return Value{}, err
	}

	count, err := strconv.Atoi(line)
	if err != nil {
		return Value{}, err
	}

	if count == -1 {
		return Value{Typ: "null"}, nil
	}

	array := make([]Value, count)
	for i := 0; i < count; i++ {
		val, err := Parse(r)
		if err != nil {
			return Value{}, err
		}
		array[i] = val
	}

	return Value{Typ: "array", Array: array}, nil
}

func readLine(r *bufio.Reader) (string, error) {
	line, err := r.ReadString('\n')
	if (err != nil) {
		return "",err
	}
	return line[:len(line)-2],nil
}


func parseSimpleString(r *bufio.Reader) (Value,error) {
	line , err := readLine(r)
	if err != nil {
		return Value{}, err
	}
	return Value{Typ :"string",Str :line}, nil
}