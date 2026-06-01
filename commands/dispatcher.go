package commands

import (
	"redis-lite/resp"
	"redis-lite/store"
	"strings"
)

type HandlerFunc func(args []resp.Value, s *store.Store) resp.Value

var handlers = map[string]HandlerFunc{
	"PING":  pingHandler,
	"GET":   getHandler,
	"SET":   setHandler,
	"TTL":   ttlHandler,
	"HELLO": helloHandler,
	"DEL":   delHandler,
	"EXPIRE": expireHandler,
}

func Dispatch(cmd string, args []resp.Value, s *store.Store) resp.Value {
	handler, ok := handlers[strings.ToUpper(cmd)]
	if !ok {
		return resp.Value{Typ: "error", Str: "ERR unknown command '" + cmd + "'"}
	}
	return handler(args, s)
}