package commands

import (
	"redis-lite/resp"
	"redis-lite/store"
	"strconv"
	"strings"
	"time"
)

func pingHandler (args []resp.Value, s *store.Store) (resp.Value) {
	if (len(args) > 1) {
		return resp.Value{Typ: "error", Str: "ERR : wrong number of arguments for PING"}
	}

	if len(args) == 0 {
		return resp.Value{Typ: "string", Str: "PONG"}
	}

	return resp.Value{Typ: "bulk", Bulk: args[0].Bulk}
}

func setHandler (args[] resp.Value, s *store.Store) (resp.Value) {
	if (len(args) < 2) {
		return resp.Value{Typ: "error", Str: "ERR : wrong number of arguments for SET"}
	}

	key, value := args[0].Bulk, args[1].Bulk


	if len(args) == 4 {

		flag := strings.ToUpper(args[2].Bulk)
		ttl, err := strconv.Atoi(args[3].Bulk)

		if (err != nil){
			return resp.Value{Typ: "error", Str: "ERR : invalid expire time"}
		}

		switch flag {
			case "EX":
				s.SetWithTTL(key, value, time.Duration(ttl)*time.Second)
				
			case "PX":
				s.SetWithTTL(key, value, time.Duration(ttl)*time.Millisecond)
				
			default:
				return resp.Value{Typ: "error", Str :"ERR : unsupported option :" + flag}
		}

		return resp.Value{Typ: "string", Str: "OK"}
	}

	s.Set(key, value)
	return resp.Value{Typ: "string", Str: "OK"}
}

func getHandler (args[] resp.Value, s *store.Store) (resp.Value) {
	if len(args) != 1 {
		return resp.Value{Typ: "error", Str: "ERR : wrong number of arguments for GET"}
	}

	val, ok := s.Get(args[0].Bulk)
	if !ok {
		return resp.Value{Typ: "null"}
	}

	return resp.Value{Typ: "bulk", Bulk: val.(string)}
}

func ttlHandler (args[] resp.Value, s *store.Store) (resp.Value) {
	if len(args) != 1 {
		return resp.Value{Typ: "error", Str: "ERR : wrong number of arguments for TTL"}
	}

	res := s.TTL(args[0].Bulk)
	return resp.Value{Typ: "integer", Num: res}
}