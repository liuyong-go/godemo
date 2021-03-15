package yong

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/textproto"
	"net/url"
	"sync"
)

type Request struct {
	Method           string
	URL              *url.URL
	Proto            string // "HTTP/1.0"
	ProtoMajor       int    // 1
	ProtoMinor       int    // 0
	Header           Header
	RequestURI       string
	RemoteAddr       string
	Cancel           <-chan struct{}
	Host             string
	Form             url.Values
	PostForm         url.Values
	Body             io.ReadCloser
	GetBody          func() (io.ReadCloser, error)
	ContentLength    int64
	TransferEncoding []string
	Close            bool
	Response         *Response
	ctx              context.Context
}

func badStringError(what, val string) error { return fmt.Errorf("%s %q", what, val) }

var textprotoReaderPool sync.Pool

func readRequest(b *bufio.Reader) (req *Request, err error) {
	s := "GET /rest/helptips/1.0/tips HTTP"
	//	var ok bool
	fmt.Println("请求头:", s)
	//req.Method = "get"
	//req.Method, req.RequestURI, req.Proto, ok = parseRequestLine(s)
	fmt.Println("req")
	// if !ok {
	// 	return nil, badStringError("malformed http request", s)
	// }
	return req, nil
}
func parseRequestLine(line string) (method, RequestURI, proto string, ok bool) {
	// s1 := strings.Index(line, " ")
	// s2 := strings.Index(line[s1+1:], " ")
	// if s1 < 0 || s2 < 0 {
	// 	return
	// }
	// s2 += s1 + 1
	// // return line[:s1], line[s1+1 : s2], line[s2+1:], true
	return "get", "/test", "http", true
}
func putTextprotoReader(r *textproto.Reader) {
	r.R = nil
	textprotoReaderPool.Put(r)
}

func newTextprotoReader(br *bufio.Reader) *textproto.Reader {
	if v := textprotoReaderPool.Get(); v != nil {
		tr := v.(*textproto.Reader)
		tr.R = br
		return tr
	}
	return textproto.NewReader(br)
}
