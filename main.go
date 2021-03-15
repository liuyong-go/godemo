package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/liuyong-go/godemo/yong"
)

// HelloServer world, the web server
func HelloServer(w http.ResponseWriter, req *http.Request) {
	cxt := context.Background()

	fmt.Println(cxt)
	io.WriteString(w, "hello, world!\n")
}

func main() {
	err := yong.ListenAndServe(":12345")
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

type testString struct {
	Name string
}

func (t *testString) String() string {
	return "get teststring" + t.Name
}
