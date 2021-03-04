package yong

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"runtime"
	"sync"
)

type Header map[string][]string
type response struct {
	conn             *conn
	req              *http.Request // request for this response
	reqBody          io.ReadCloser
	cancelCtx        context.CancelFunc // when ServeHTTP exits
	wroteHeader      bool               // reply header has been (logically) written
	wroteContinue    bool               // 100 Continue response was written
	wants10KeepAlive bool               // HTTP/1.0 w/ Connection "keep-alive"
	wantsClose       bool               // HTTP request has Connection "close"

	// canWriteContinue is a boolean value accessed as an atomic int32
	// that says whether or not a 100 Continue header can be written
	// to the connection.
	// writeContinueMu must be held while writing the header.
	// These two fields together synchronize the body reader
	// (the expectContinueReader, which wants to write 100 Continue)
	// against the main writer.

	// handlerHeader is the Header that Handlers get access to,
	// which may be retained and mutated even after WriteHeader.
	// handlerHeader is copied into cw.header at WriteHeader
	// time, and privately mutated thereafter.
	handlerHeader Header
	calledHeader  bool // handler accessed handlerHeader via Header

	written       int64 // number of bytes written in body
	contentLength int64 // explicitly-declared Content-Length; or -1
	status        int   // status code passed to WriteHeader

}
type conn struct {
	server     *Server
	cancelCtx  context.CancelFunc
	rwc        net.Conn
	remoteAddr string
	mu         sync.Mutex
	lastMethod string
	bufr       *bufio.Reader
}

//Server 结构体
type Server struct {
	Addr        string
	mu          sync.Mutex
	listeners   map[*net.Listener]struct{}
	activeConn  map[*conn]struct{}
	doneChan    chan struct{}
	BaseContext func(net.Listener) context.Context
	// ConnContext optionally specifies a function that modifies
	// the context used for a new connection c. The provided ctx
	// is derived from the base context and has a ServerContextKey
	// value.
	ConnContext func(ctx context.Context, c net.Conn) context.Context
}

//ontextKey is a value for use with context.WithValue
type contextKey struct {
	name string
}

var (
	// ServerContextKey is a context key. It can be used in HTTP
	// handlers with Context.Value to access the server that
	// started the handler. The associated value will be of
	// type *Server.
	ServerContextKey = &contextKey{"http-server"}
	// LocalAddrContextKey is a context key. It can be used in
	// HTTP handlers with Context.Value to access the local
	// address the connection arrived on.
	// The associated value will be of type net.Addr.
	LocalAddrContextKey = &contextKey{"local-addr"}
)

var ErrServerClosed = errors.New("http: Server closed")
var ErrAbortHandler = errors.New("net/http: abort Handler")

//ListenAndServe 初始化配置启动服务
func ListenAndServe(addr string) error {
	server := &Server{Addr: addr}
	return server.ListenAndServe()
}

//ListenAndServe listen服务然后开启
func (srv *Server) ListenAndServe() error {
	addr := srv.Addr
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return srv.Server(ln)
}

//Server 接收l上的请求链接，为每一个请求创建协程
func (srv *Server) Server(l net.Listener) error {
	//map 添加net.listener 跟踪
	if !srv.trackListener(&l, true) {
		return ErrServerClosed
	}
	origListener := l
	defer srv.trackListener(&l, false)
	baseCxt := context.Background()
	if srv.BaseContext != nil {
		baseCxt = srv.BaseContext(origListener)
		if baseCxt == nil {
			panic("BaseContext returned a nil context")
		}
	}
	ctx := context.WithValue(baseCxt, ServerContextKey, srv)
	//接收请求
	for {
		rw, err := l.Accept()
		if err != nil {
			select {
			case <-srv.getDoneChan():
				return ErrServerClosed
			default:
			}
		}
		connCtx := ctx
		if cc := srv.ConnContext; cc != nil {
			connCtx = cc(connCtx, rw)
			if connCtx == nil {
				panic("connContext return nil")
			}
		}
		c := srv.newConn(rw)
		go c.serve(connCtx)
	}
}

//Serve a new connection.
func (c *conn) serve(ctx context.Context) {
	c.remoteAddr = c.rwc.RemoteAddr().String()
	ctx = context.WithValue(ctx, LocalAddrContextKey, c.rwc.LocalAddr())
	//处理异常
	defer func() {
		if err := recover(); err != nil && err != ErrAbortHandler {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			fmt.Sprintf("http: panic serving %v: %v\n%s", c.remoteAddr, err, buf)
		}
	}()
	ctx, cancelCtx := context.WithCancel(ctx)
	c.cancelCtx = cancelCtx
	defer cancelCtx()
	for {
		w, err := c.readRequest(ctx)
		if err != nil {
			fmt.Println("read request err", err)
		}
		req := w.req
		fmt.Println("response", req)

	}
}

//Read next request from connection.
func (c *conn) readRequest(ctx context.Context) (w *response, err error) {
	fmt.Println("read request")
	return
}
func (s *Server) getDoneChan() <-chan struct{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.getDoneChanLocked()
}
func (s *Server) getDoneChanLocked() chan struct{} {
	if s.doneChan == nil {
		s.doneChan = make(chan struct{})
	}
	return s.doneChan
}

//trackListener 添加或删除 net.Listener 设置 listeners 跟踪
func (srv *Server) trackListener(ln *net.Listener, add bool) bool {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	if srv.listeners == nil {
		srv.listeners = make(map[*net.Listener]struct{})
	}
	if add {
		srv.listeners[ln] = struct{}{}
	} else {
		delete(srv.listeners, ln)
	}
	return true
}

// Create new connection from rwc.
func (srv *Server) newConn(rwc net.Conn) *conn {
	c := &conn{
		server: srv,
		rwc:    rwc,
	}
	return c
}
