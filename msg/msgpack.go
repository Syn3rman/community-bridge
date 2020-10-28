package main

import (
	"fmt"
	"time"
	"net"

	"github.com/vmihailenco/msgpack/v5"
)

type test struct{
	_msgpack struct{}       `msgpack:",asArray"`
	Name string             `msgpack:"name"`
	Ts int64                `msgpack:"ts"`
	Attrs map[string]string `msgpack:"attrs"`
}

func main(){
	testkv := make(map[string]string)
	testkv["foo"] = "bar"
	t, err := msgpack.Marshal(&test{Name: "foo", Ts: time.Now().UnixNano(), Attrs: testkv})
	check(err)
	fmt.Println(t)

	fmt.Println(time.Now().UnixNano())

	url := "localhost:24224"
	tcpAddr, err := net.ResolveTCPAddr("tcp", url)
	check(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	check(err)

	_, err = conn.Write(t)
	check(err)
}

func check(err error){
	if err!=nil{
		fmt.Println(err)
	}
}
