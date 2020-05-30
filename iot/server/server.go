package server


import (
	"fmt"
	"net"
	"os"
	
	//	"strings"
)



type IoTServer struct {
}

const (
    CONN_HOST = "192.168.8.100"
    CONN_PORT = "3000"
    CONN_TYPE = "tcp"
)


func NewServer() *IoTServer {
	return &IoTServer{
	}
}



func (s *IoTServer) Run() {
    // Listen for incoming connections.
    l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }
    // Close the listener when the application closes.
    defer l.Close()
    fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
    for {
        // Listen for an incoming connection.
        conn, err := l.Accept()
        if err != nil {
            fmt.Println("Error accepting: ", err.Error())
            os.Exit(1)
        }

			fmt.Printf("Accepted connection\n")
			
			// Handle connections in a new goroutine.
			c := NewConnection(conn)
			go c.Run()
			//go s.handleRequest(conn)
    }
}





