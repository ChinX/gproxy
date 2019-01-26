package server

import (
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/shadowsocks/shadowsocks-go/shadowsocks"

	"github.com/elazarl/goproxy"
	"golang.org/x/net/proxy"
)

var unixAddr = ""

type Server struct {
	http      string
	socks     string
	waitGroup sync.WaitGroup
}

func NewServer(socksAddr, httpAddr string) *Server {
	return &Server{socks: socksAddr, http: httpAddr}
}

func (s *Server) ListenAndProxy(addr, method, password string) (err error) {
	s.waitGroup.Add(1)
	go s.socksListenAndProxy(addr, method, password)
	if s.http != ""{
		s.waitGroup.Add(1)
		go s.httpListenAndProxy()
	}
	s.waitGroup.Wait()
	return
}

func (s *Server) socksListenAndProxy(addr, method, password string) {
	defer s.waitGroup.Done()
	socksProxy, err := NewShadowSocks(addr, method, password)
	if err != nil {
		log.Println(err)
		return
	}
	ln, err := net.Listen("tcp", s.socks)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("starting local socks5 server at %v ...\n",  s.socks)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("accept:", err)
			continue
		}
		go socksProxy.HandleConnection(conn)
	}
}

func (s *Server) httpListenAndProxy() {
	defer s.waitGroup.Done()
	httpProxy := goproxy.NewProxyHttpServer()
	httpProxy.Verbose = true

	dialer, err := proxy.SOCKS5("tcp", s.socks, nil, proxy.Direct)
	if err != nil {
		log.Fatal(err)
	}

	httpProxy.Tr = &http.Transport{Dial: dialer.Dial}

	log.Printf("starting local http server at %v ...\n", s.http)
	err = http.ListenAndServe(s.http, httpProxy)
	if err != nil {
		log.Println(err)
	}
}

func IsFileExists(path string) (bool, error) {
	return shadowsocks.IsFileExists(path)
}
