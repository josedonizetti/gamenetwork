package main

import (
    "net"
    "fmt"
    "os"
    "time"
)

type Client struct {
  Host string
  Port string
}

func NewClient(host, port string) *Client {
  return &Client{host, port}
}

func (client *Client) Send(packet []byte) {
  address, err := net.ResolveUDPAddr("udp", client.Host+":"+client.Port)
  if err != nil {
    fmt.Println(err.Error())
    os.Exit(1)
  }

  conn, err := net.DialUDP("udp", nil, address)
  if err != nil {
    fmt.Println(err.Error())
    os.Exit(1)
  }

  _, err = conn.Write(packet)
  if err != nil {
    fmt.Println(err.Error())
    os.Exit(1)
  }

  _, err = conn.Write(packet)
  if err != nil {
    fmt.Println(err.Error())
    os.Exit(1)
  }

  _, err = conn.Write(packet)
  if err != nil {
    fmt.Println(err.Error())
    os.Exit(1)
  }

  fmt.Println("sleeping")
  time.Sleep(20 * time.Second)
  fmt.Println("waking")

  _, err = conn.Write(packet)
  if err != nil {
    fmt.Println(err.Error())
    os.Exit(1)
  }

  os.Exit(0)
}


func main() {
  client := NewClient("localhost", "1200")
  client.Send([]byte("client to server"))
}
