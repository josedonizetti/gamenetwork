package main

import (
    "net"
    "fmt"
    "os"
)

type Server struct {
  port string
  clients map[string]chan []byte
}

func (server *Server) Start() {
  address, err := net.ResolveUDPAddr("udp", ":"+server.port)

  if err != nil {
    fmt.Println(err.Error())
    os.Exit(1)
  }

  conn, err := net.ListenUDP("udp", address)
  if err != nil {
    fmt.Println(err.Error())
    os.Exit(1)
  }

  for {
    var buf [512]byte
    // blocking, waiting for connection
    size, addr, err := conn.ReadFromUDP(buf[0:])
    if err != nil {
      continue
    }

    clientId := addr.String()
    fmt.Println("packet receive from client: ", clientId)
    if server.clients[clientId] != nil {
      // Client already exist, passing payload to channel
      server.clients[clientId] <- buf[0:size]
    } else {
      // Client does not exist, creating new and passing payload to channel
      channel := make(chan []byte)
      server.clients[clientId] = channel
      go handleClient(channel)
      channel <- buf[0:size]
    }
  }
}

func handleClient(channel chan []byte) {
  for {
    packet := <-channel
    //validate protocol
    //exit for, if connection timeout
    fmt.Println("payload:", string(packet))
  }
}


func main() {
    server := &Server{"1200", make(map[string]chan []byte)}
    server.Start()
}
