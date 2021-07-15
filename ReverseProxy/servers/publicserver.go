package main
import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

func pub_process(Snum int,Rnum int,conn net.Conn){
	defer conn.Close() // 关闭连接
	reader := bufio.NewReader(conn)
	msg,err:=reader.ReadString('\n')
	if err == io.EOF {
		return
	}
	if err != nil {
		fmt.Println("decode msg failed, err:", err)
		return 
	}
	str := strings.Split(msg, " ")
	fmt.Println("参数:",str[1])
	SendParam(Snum,Rnum,str[1],conn)
}

func PublicInitRange(PnumF int,PnumT int){
	if PnumT<PnumF{
		fmt.Println("PnumF不能超过PnumT")
		return
	}
	if PnumT+6000>65535{
		fmt.Println("PnumT不能超过",65535-6000)
		return
	}
	fmt.Println("开放Pnum:",strconv.Itoa(PnumF),"~"+strconv.Itoa(PnumT)+",Rnum:",strconv.Itoa(PnumF+6000),"~"+strconv.Itoa(PnumT+6000)+",Snum:",strconv.Itoa(PnumF+5000),"~"+strconv.Itoa(PnumT+5000))
	for b:=PnumF;b<=PnumT-1;b++  {
		go func(num int){
			go ReceInit(num+6000)
			go SendInit(num+5000)
			//fmt.Println("公网服务器开启","127.0.0.1:"+strconv.Itoa(num))
			//defer fmt.Println("公网服务器关闭")
			listen, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(num))
			if err != nil {
				fmt.Println("listen failed, err:", err)
				return
			}
			for {
				conn, err := listen.Accept() // 建立连接
				if err != nil {
					fmt.Println("accept failed, err:", err)
				}
				go pub_process(num+5000,num+6000,conn) // 启动一个goroutine处理连接
			}
		}(b)
	}
	//fmt.Println("公网服务器开启","127.0.0.1:"+strconv.Itoa(PnumT))
	//defer fmt.Println("公网服务器关闭")
	go ReceInit(PnumT+6000)
	go SendInit(PnumT+5000)
	listen, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(PnumT))
	if err != nil {
		fmt.Println("listen failed, err:", err)
		return
	}
	for {
		conn, err := listen.Accept() // 建立连接
		if err != nil {
			fmt.Println("accept failed, err:", err)
		}
		go pub_process(PnumT+5000,PnumT+6000,conn) // 启动一个goroutine处理连接
	}
}

func PublicInit(Rnum int,Snum int){
	go ReceInit(Rnum)
	go SendInit(Snum)
	fmt.Println("公网服务器开启")
	defer fmt.Println("公网服务器关闭")
	listen, err := net.Listen("tcp", "127.0.0.1:80")
	if err != nil {
		fmt.Println("listen failed, err:", err)
		return
	}
	for {
		conn, err := listen.Accept() // 建立连接
		if err != nil {
			fmt.Println("accept failed, err:", err)
		}
		go pub_process(Snum,Rnum,conn) // 启动一个goroutine处理连接
	}
}