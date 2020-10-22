package main

import (
	"fmt"
	"io"
	"bytes"
	"time"
	"io/ioutil"
	"net/http"
	"github.com/ugorji/go/codec"
)

type test struct{
	_struct bool    `codec:",toarray"`
	Name string `codec:"name"`
	Ts int64 `codec:"ts"`
	Attrs map[string]string `codec:"attrs"`
}

func main(){
	var (
		v interface{} // value to decode/encode into
		r io.Reader
		w io.Writer
		b []byte
		mh codec.MsgpackHandle
	)

	testkv := make(map[string]string)
	testkv["foo"] = "bar"

	enc := codec.NewEncoder(w, &mh)
	enc = codec.NewEncoderBytes(&b, &mh)
	err := enc.Encode(&test{Name: "foo", Ts: time.Now().UnixNano(), Attrs: testkv})
	check(err)

	fmt.Printf("%T, %v\n\n", b, b)

	dec := codec.NewDecoder(r, &mh)
	dec = codec.NewDecoderBytes(b, &mh)
	err = dec.Decode(&v) 
	check(err)

	// v is of type []interface{}
	fmt.Printf("%T, %v", v, v)

	url := "http://localhost:24224/"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b)) 
	check(err)
	client := &http.Client{}
	resp, err := client.Do(req)
	check(err)
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(content)
}

func check(err error){
	if err!=nil{
		panic(err)
	}
}
