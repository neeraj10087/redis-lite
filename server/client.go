package server

import (
	"bufio"
	"fmt"
	"net"
	"redis-lite/commands"
	"redis-lite/resp"
	"redis-lite/store"
	"strings"
	"golang.org/x/sys/unix"
)

func handleClient(conn net.Conn, s *store.Store) {
	defer conn.Close()
	fmt.Println("Client connected:", conn.RemoteAddr())

	reader := bufio.NewReader(conn)

	for {
		value, err := resp.Parse(reader)
		if err != nil {
			fmt.Println("Client disconnected:", conn.RemoteAddr())
			return
		}

		if len(value.Array) == 0 {
			continue
		}

		cmd := value.Array[0].Bulk
		args := value.Array[1:]

		result := commands.Dispatch(cmd, args, s)
		conn.Write(result.Marshal())
	}
}

func handleClientFd(fd int, s *store.Store) {
	buf := make([]byte, 1024)
	n, err := unix.Read(fd, buf)

	if err != nil || n == 0 {
		unix.Close(fd)
		return
	}

	reader := bufio.NewReader(strings.NewReader(string(buf[:n])))
	value, err := resp.Parse(reader)

	if err != nil || len(value.Array) == 0 {
		unix.Close(fd)
		return
	}

	cmd := value.Array[0].Bulk
	args := value.Array[1:]

	result := commands.Dispatch(cmd, args, s)
	// conn.Write(result.Marshal())
	unix.Write(fd, result.Marshal())

}