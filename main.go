/* ThreadedEchoServer
 */
package main

import (
  "net"
  "os"
  // "os/exec"
  "time"
  "fmt"
  "encoding/json"
  "code.google.com/p/goconf/conf"
)

type NostraRequest struct {
  Version string
  Params []string
}

func main() {
  c, err := conf.ReadConfigFile("nostra.conf")
  checkError(err)

  port, err := c.GetString("default", "port")
  service := ":" + port

  tcpAddr, err := net.ResolveTCPAddr("ip4", service)
  checkError(err)

  listener, err := net.ListenTCP("tcp", tcpAddr)
  checkError(err)

  for {
    conn, err := listener.Accept()
    if err != nil {
      continue
    }

    go handleClient(conn)
  }
}

func handleClient(conn net.Conn) {
  defer conn.Close()

  var buf [512]byte
  for {
    n, err := conn.Read(buf[0:])
    if err != nil {
      return
    }

    var request NostraRequest

    json.Unmarshal(buf[0:n], &request)

    if request.Version == "" || len(request.Params) == 0 {
      conn.Write([]byte("{\"code\":-1,\"data\":{\"message\":\"Malformed input\"}}"))
      return
    }

    ret := "{\"code\":0,\"data\":{"

    for _, param := range(request.Params) {
      switch param {
      case "hostname":
        // cmd := exec.Command("hostname")
        // output, err := cmd.Output()

        // hostname := (string(output)[:len(output) - 1])

        hostname, err := os.Hostname()

        if err != nil {
          conn.Write([]byte("{\"code\":-2,\"data\":{\"message\":\"Server derped\"}}"))
          return
        } else {
          ret += "\"hostname\":\"" + hostname + "\","
        }
      case "time":
        time := time.Now().Format(time.RFC3339)

        ret += "\"time\":\"" + time + "\","
      default:
        conn.Write([]byte("{\"code\":-3,\"data\":{\"message\":\"Unknown parameter (" + param + ")\"}}"))

        return
      }
    }

    ret = ret[:len(ret) - 1]
    ret += "}}"

    conn.Write([]byte(ret))

    return
  }
}

func checkError(err error) {
  if err != nil {
    fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
    os.Exit(1)
  }
}
