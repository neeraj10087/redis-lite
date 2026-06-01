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

func handleClientFd(fd int, s *store.Store, currentClients *int) {
	buf := make([]byte, 1024)
	n, err := unix.Read(fd, buf)

	if err != nil || n == 0 {
		unix.Close(fd)
		*currentClients--
		fmt.Println("client disconnected, total:", *currentClients)
		return
	}

	reader := bufio.NewReader(strings.NewReader(string(buf[:n])))
	value, err := resp.Parse(reader)
	printCommand(value)

	if err != nil || len(value.Array) == 0 {
		return
	}

	cmd := value.Array[0].Bulk
	args := value.Array[1:]

	result := commands.Dispatch(cmd, args, s)
	unix.Write(fd, result.Marshal())

}

func printCommand(value resp.Value) {
	// iterate over array and print each element
	fmt.Print("Command: ")
	for _, arg := range value.Array {
		fmt.Printf("[%s] ", arg.Bulk)
	}
	fmt.Println()
}