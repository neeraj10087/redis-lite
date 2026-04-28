package server

import (
	"fmt"
	"net"
	"redis-lite/store"
	"golang.org/x/sys/unix"
	"log"
)

func Start(addr string, s *store.Store) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	fmt.Println("Redis-Lite listening on", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}
		handleClient(conn, s)
	}
}

func StartAsync(addr string, s *store.Store) {

	maxClients := 20000
    events := make([]unix.Kevent_t, maxClients)

	// creating raw socket
	serverFD, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	if err != nil {
		panic(err)
	}
	defer unix.Close(serverFD)

	if err = unix.SetNonblock(serverFD, true); err != nil {
        log.Fatal(err)
    }

	ip4 := net.ParseIP("127.0.0.1")
	if err = unix.Bind(serverFD, &unix.SockaddrInet4{
		Port: 7379,
		Addr: [4]byte{ip4[12], ip4[13], ip4[14], ip4[15]},
	}); err != nil {
		log.Fatal(err)
	}

	if err = unix.Listen(serverFD, maxClients); err != nil {
        log.Fatal(err)
    }

	// fmt.Println("Redis-Lite listening on", addr)

	// create a kqueue
	kqFD, err := unix.Kqueue()
	if err != nil {
		log.Fatal(err)
	}
	defer unix.Close(kqFD)

	// fmt.Println("kqueue created")

	// register server socket with kqueue
	serverEvent := unix.Kevent_t{Ident: uint64(serverFD),Filter: unix.EVFILT_READ, Flags: unix.EV_ADD}

	_, err = unix.Kevent(kqFD, []unix.Kevent_t{serverEvent}, nil,nil)
	if (err != nil) {
		log.Fatal(err)
	}

	// fmt.Println("registered server socket with kqueue")

	for {
		nevents, err := unix.Kevent(kqFD, nil, events, nil)
		// fmt.Println("recieved events from kqueue")
		if err != nil {
			continue
		}

		for i:=0 ; i < nevents; i++ {
			fd := int(events[i].Ident)
			if (fd == serverFD){
				// event is from server fd 
				// this for a new connection
				clientFd, _, err := unix.Accept(serverFD)
				if err != nil {
                	log.Println("accept error:", err)
                    continue
                }

				unix.SetNonblock(clientFd,true)

				// register this client fd in kqueue
				clientEvent := unix.Kevent_t{Ident: uint64(clientFd), Filter: unix.EVFILT_READ, Flags: unix.EV_ADD}

				_, err = unix.Kevent(kqFD,[]unix.Kevent_t{clientEvent},nil,nil)

			} else {
				handleClientFd(fd, s)
			}
		}
	}


}