package commands

import (
	"redis-lite/resp"
	"redis-lite/store"
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
	if (len(args) != 2) {
		return resp.Value{Typ: "error", Str: "ERR : wrong number of arguments for SET"}
	}

	s.Set(args[0].Bulk, args[1].Bulk)
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