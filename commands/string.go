package commands

import (
	"redis-lite/resp"
	"redis-lite/store"
	"strconv"
	"strings"
	"time"
)

func pingHandler(args []resp.Value, s *store.Store) resp.Value {
	if len(args) > 1 {
		return resp.Value{Typ: "error", Str: "ERR : wrong number of arguments for PING"}
	}

	if len(args) == 0 {
		return resp.Value{Typ: "string", Str: "PONG"}
	}

	return resp.Value{Typ: "bulk", Bulk: args[0].Bulk}
}

func setHandler(args []resp.Value, s *store.Store) resp.Value {
	if len(args) < 2 {
		return resp.Value{Typ: "error", Str: "ERR : wrong number of arguments for SET"}
	}

	key, value := args[0].Bulk, args[1].Bulk

	if len(args) == 4 {

		flag := strings.ToUpper(args[2].Bulk)
		ttl, err := strconv.Atoi(args[3].Bulk)

		if err != nil {
			return resp.Value{Typ: "error", Str: "ERR : invalid expire time"}
		}

		switch flag {
		case "EX":
			s.SetWithTTL(key, value, time.Duration(ttl)*time.Second)

		case "PX":
			s.SetWithTTL(key, value, time.Duration(ttl)*time.Millisecond)

		default:
			return resp.Value{Typ: "error", Str: "ERR : unsupported option :" + flag}
		}

		return resp.Value{Typ: "string", Str: "OK"}
	}

	s.Set(key, value)
	return resp.Value{Typ: "string", Str: "OK"}
}

func getHandler(args []resp.Value, s *store.Store) resp.Value {
	if len(args) != 1 {
		return resp.Value{Typ: "error", Str: "ERR : wrong number of arguments for GET"}
	}

	val, ok := s.Get(args[0].Bulk)
	if !ok {
		return resp.Value{Typ: "null"}
	}

	return resp.Value{Typ: "bulk", Bulk: val.(string)}
}

func ttlHandler(args []resp.Value, s *store.Store) resp.Value {
	if len(args) != 1 {
		return resp.Value{Typ: "error", Str: "ERR : wrong number of arguments for TTL"}
	}

	res := s.TTL(args[0].Bulk)
	return resp.Value{Typ: "integer", Num: res}
}

func delHandler(args []resp.Value, s *store.Store) resp.Value {
	var deletedKey int = 0
	for _, arg := range args {
		if s.Delete(arg.Bulk) {
			deletedKey++
		}
	}
	return resp.Value{Typ: "integer", Num: deletedKey}
}

func expireHandler(args []resp.Value, s *store.Store) resp.Value {
	if len(args) != 2 {
		return resp.Value{Typ: "error", Str: "ERR : wrong number of arguments for EXPIRE"}
	}
	key := args[0].Bulk
	ttl, err := strconv.Atoi(args[1].Bulk)
	
	if err != nil {
		return resp.Value{Typ: "error", Str: "ERR : invalid expire time"}
	}

	keyObj , ok := s.Get(key)
	if !ok {
		return resp.Value{Typ: "integer", Num: 0}
	}

	s.SetWithTTL(key, keyObj, time.Duration(ttl)*time.Second)
	return resp.Value{Typ: "integer", Num: 1}
}

func helloHandler(args []resp.Value, s *store.Store) resp.Value {
	return resp.Value{
		Typ: "array",
		Array: []resp.Value{
			{Typ: "bulk", Bulk: "server"},
			{Typ: "bulk", Bulk: "redis-lite"},
			{Typ: "bulk", Bulk: "version"},
			{Typ: "bulk", Bulk: "0.0.1"},
			{Typ: "bulk", Bulk: "proto"},
			{Typ: "integer", Num: 2},
			{Typ: "bulk", Bulk: "id"},
			{Typ: "integer", Num: 1},
			{Typ: "bulk", Bulk: "mode"},
			{Typ: "bulk", Bulk: "standalone"},
			{Typ: "bulk", Bulk: "role"},
			{Typ: "bulk", Bulk: "master"},
			{Typ: "bulk", Bulk: "modules"},
			{Typ: "array", Array: []resp.Value{}},
		},
	}
}
