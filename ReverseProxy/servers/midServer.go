package main

import (
	"errors"
	"fmt"
	"net"
	"ReverseProxy/protocol"
	"strconv"
)


var SendConnMap =make(map[int]net.Conn,100)
var ReceConnMap =make(map[int]net.Conn,100)
var ReceChanMap =make(map[int]chan []byte,100)

func SendRegister(num int,conn net.Conn ) error{
	_, ok := SendConnMap[num]
	if ok{
		return errors.New("already has a connection in SendConnMap")
	}else{
		SendConnMap[num]=conn
		return nil
	}
}
func SendUnRegister(num int){
	delete(SendConnMap,num)
}
func ReceUnRegister(num int){
	delete(ReceConnMap,num)
	delete(ReceChanMap,num)
}
func ReceRegister(num int,conn net.Conn ) error{
	_, ok := ReceConnMap[num]
	if ok{
		return errors.New("already has a connection in ReceConnMap")
	}else{
		ReceConnMap[num]=conn
		return nil
	}
}
func SendInit(num int) {
	fmt.Println("Send Init ->0.0.0.0::"+strconv.Itoa(num))
	listen, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(num))
	if err != nil {
		fmt.Println("listen failed, err:", err)
		return
	}
	//fmt.Println("正在连接发送隧道")
	for   {
		conn, err := listen.Accept() // 建立连接
		if err != nil {
			fmt.Println("accept failed, err:", err)
			return
		}
		fmt.Println(num,"Send")
		err=SendRegister(num,conn)
		if err !=nil {
			fmt.Println(err.Error())
			conn.Close()
		}

	}

}
func SendParam(Snum int,Rnum int,msg string, conn net.Conn) {
	_,ok:=SendConnMap[Snum]
	conn.Write([]byte("HTTP/1.1 200 OK  \r\n\r\n"))
	if ok  {
		fmt.Println("midServer Send:", msg)
		_, err := SendConnMap[Snum].Write(protocol.Packet([]byte(msg)))
		if err != nil {
			msg := "<html><header><meta  charset='UTF-8'><title>隧道连接断开</title></header><h1>隧道连接突然断开，需要重新连接</h1></html>"
			conn.Write([]byte(msg))
			fmt.Println("Send Param 移除断开连接:",Snum,Rnum)
			SendUnRegister(Snum)
			ReceUnRegister(Rnum)
		} else {
			_,ok=ReceConnMap[Rnum]
			if ok{
				RecHelp(Rnum,conn)
			}else{
				msg := "<html><header><meta  charset='UTF-8'><title>还没有隧道连接</title></header><h1>还没有隧道连接</h1></html>"
				conn.Write([]byte(msg))
			}
		}
	} else {
		msg := "<html><head><meta http-equiv='content-type' content='text/html;charset=utf-8'><title>还没有隧道连接</title></head><h1>还没有隧道连接</h1></html>"
		//fmt.Println(msg)
		conn.Write([]byte(msg))
	}
}
func ReceResponse(num int) {
	_,ok:=ReceConnMap[num]
	if ok {
		tmpBuffer := make([]byte, 0)
		//声明一个管道用于接收解包的数据
		buffer := make([]byte, 1024)
		for {
			n, err := ReceConnMap[num].Read(buffer)
			if err != nil {
				fmt.Println("Rece：", err.Error())
				fmt.Println("移除断开连接:20000,30000")
				SendUnRegister(num-1000)
				ReceUnRegister(num)
				return
			}
			//fmt.Println(n)
			tmpBuffer = protocol.Unpack(append(tmpBuffer, buffer[:n]...), ReceChanMap[num])
		}
	} else {
		fmt.Println("还没有隧道连接")
	}
}
func RecHelp(num int,conn net.Conn) {
	data := <-ReceChanMap[num]
	fmt.Println("midServer Rec:", string(data))
	conn.Write(data)
}
func ReceInit(num int) {
	fmt.Println("Rece Init ->0.0.0.0:"+strconv.Itoa(num))
	listen, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(num))
	if err != nil {
		fmt.Println("listen failed, err:", err)
		return
	}
	//fmt.Println("正在连接接收隧道")
	for {
		conn, err := listen.Accept() // 建立连接
		if err != nil {
			fmt.Println("accept failed, err:", err)
			return
		}
		fmt.Println(num,"Rece")
		err=ReceRegister(num,conn)
		if err!=nil{
			fmt.Println(err.Error())
			conn.Close()
			continue
		}
		ReceChanMap[num]= make(chan []byte, 16)
		go ReceResponse(num)
	}
}