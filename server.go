package main

import (
    "net"
    "fmt"
    "os"
    "time"
)

type Message struct {
  id string
  buf []byte
}

type Connection struct {
  id string
  messages chan []byte
  ticker *time.Ticker
  lastTimeReceived int64
  timeout int64
  removeConnection chan string
}

func NewConnection(id string, messages chan []byte, timeout int64, removeConnection chan string) *Connection {
  ticker := time.NewTicker(5 * time.Second)
  return &Connection{id, messages, ticker, 0, timeout, removeConnection}
}

type Server struct {
  port string
  timeout int64
  clients map[string] *Connection
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

  newConnection := make(chan Message)
  removeConnection := make(chan string)
  go handleConnections(server, newConnection, removeConnection);

  for {
    var buf [512]byte
    // blocking, waiting for connection
    size, addr, err := conn.ReadFromUDP(buf[0:])
    if err != nil {
      continue
    }

    clientId := addr.String()
    newConnection <- Message{clientId,buf[0:size]}
  }
}

func handleConnections(server *Server, newConnection chan Message, removeConnection chan string) {
  for {
    select {
    case message := <- newConnection:
      clientId := message.id
      if server.clients[clientId] != nil {
        fmt.Println("New Packet: ", clientId)
        // Client already exist, passing payload to channel
        server.clients[clientId].messages <- message.buf
      } else {
        fmt.Println("New Client:", clientId)
        // Client does not exist, creating new and passing payload to channel
        messages := make(chan []byte)
        connection := NewConnection(clientId, messages, server.timeout, removeConnection)
        server.clients[clientId] = connection

        go handleConnection(connection)
        connection.messages <- message.buf
      }

    case clientId := <- removeConnection:
      fmt.Println("removing client:", clientId)
      delete(server.clients, clientId)
    }
  }
}

func handleConnection(connection *Connection) {
  outer:
  for {
    select {
      case packet := <-connection.messages:
        //validate protocol
        connection.lastTimeReceived = time.Now().Unix()
        fmt.Println("payload:", connection.id, string(packet))
      case time := <-connection.ticker.C:
        fmt.Println("ticker", time)
        diff := time.Unix() - connection.lastTimeReceived
        if diff > connection.timeout {
          connection.removeConnection <- connection.id
          break outer
        }
    }
  }
  close(connection.messages)
}

func main() {
    server := &Server{"1200", 10, make(map[string]*Connection)}
    server.Start()
}
