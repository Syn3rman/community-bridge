package main

import (
	"fmt"
	"net"
	"time"

	"github.com/vmihailenco/msgpack/v5"
)

type test struct {
	_msgpack struct{}          `msgpack:",asArray"`
	Name     string            `msgpack:"name"`
	Ts       int64             `msgpack:"ts"`
	Attrs    map[string]string `msgpack:"attrs"`
}

func main() {
	testkv := make(map[string]string)
	testkv["foo"] = "bar"
	t, err := msgpack.Marshal(&test{Name: "foo", Ts: time.Now().UnixNano(), Attrs: testkv})
	check(err)
	fmt.Println(t)

	fmt.Println(time.Now().UnixNano())

	url := "localhost:24224"

	udpconn, _ := net.Dial("udp", url)
	data, _ := msgpack.Marshal([]byte{0x00})
	fmt.Println("data: ", data)
	res, err := udpconn.Write(data)
	check(err)
	fmt.Println(res)
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
