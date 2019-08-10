package gateway

import (
	"fmt"
	"net"
	"runtime"

	log "github.com/sirupsen/logrus"
)

type Conf struct {
	Address    string //ip:port
	MaxConnNum int    //support max conn
	AcceptNum  int    //cur conn num
}

type TcpServer struct {
	Address    string //ip:port
	MaxConnNum int    //support max conn
	AcceptNum  int    //cur conn num
}

func NewTcpServer(conf Conf) *TcpServer {
	return &TcpServer{
		Address:    conf.Address,
		MaxConnNum: conf.MaxConnNum,
		AcceptNum:  conf.AcceptNum,
	}
}

func (s *TcpServer) Start() {
	addr, err := net.ResolveTCPAddr("tcp", s.Address)
	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal("error listening:", err.Error())
		return
	}

	for i := 0; i < s.AcceptNum; i++ {
		go s.Accept(listener)
	}
	select {}
}

func (s *TcpServer) Accept(listener *net.TCPListener) {
	defer RecoverPanic()
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Fatal(err)
			continue
		}

		err = conn.SetKeepAlive(true)
		if err != nil {
			log.Fatal(err)
		}

		client := NewClient(conn)
		go client.DoConn()
	}
}

// RecoverPanic 恢复panic
func RecoverPanic() {
	err := recover()
	if err != nil {
		log.Fatal(GetPanicInfo())
	}
}

// PrintStaStack 打印Panic堆栈信息
func GetPanicInfo() string {
	buf := make([]byte, 2048)
	n := runtime.Stack(buf, false)
	return fmt.Sprintf("%s", buf[:n])
}
