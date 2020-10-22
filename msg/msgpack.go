package main

import (
	"fmt"
	"bytes"
	// "strings"
	"time"
	"io/ioutil"
	"net/http"

	"github.com/vmihailenco/msgpack/v5"
)

type test struct{
	_msgpack struct{} `msgpack:",asArray"`
	Name string `msgpack:"name"`
	Ts int64 `msgpack:"ts"`
	Attrs map[string]string `msgpack:"attrs"`
}

func main(){
	var buf bytes.Buffer
	enc := msgpack.NewEncoder(&buf)
	testkv := make(map[string]string)
	testkv["foo"] = "bar"
	err := enc.Encode(&test{Name: "foo", Ts: time.Now().UnixNano(), Attrs: testkv})
	if err != nil {
		panic(err)					    
	}
	// Decode the messagepack array into slice of interface
	dec := msgpack.NewDecoder(&buf)
	v, err := dec.DecodeSlice()
	if err != nil {
		    panic(err)
				
	}
	fmt.Printf("%T, %v\n\n", v, v)
	// fmt.Printf("%T, %v", buf, buf)
	url := "http://localhost:24224/"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(v))
	check(err)
	client := &http.Client{}
	resp, err := client.Do(req)
	check(err)
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	fmt.Println("Response body: \n", respBody)
}

func check(err error){
	if err!=nil{
		panic(err)
	}
}
