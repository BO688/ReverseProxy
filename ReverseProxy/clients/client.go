package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"ReverseProxy/protocol"
	"strconv"
	"time"
)
var Check =false
var Snum int
var Rnum int
var Local int
var Remote string
func Recchannel(conn net.Conn){
	tmpBuffer := make([]byte, 0)
	//声明一个管道用于接收解包的数据
	readerChannel := make(chan []byte, 16)
	go RecHelp(readerChannel)
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println(err)
			if Check {
				Init(5)
			}
			return
		}
		Check=true
		tmpBuffer = protocol.Unpack(append(tmpBuffer, buffer[:n]...), readerChannel)
	}
}
func RecHelp(readerChannel chan []byte) {
	for {
		select {
		case data := <-readerChannel: {
			//fmt.Println("client rece:",string(data))
			resp, err := http.Get("http://127.0.0.1:"+strconv.Itoa(Local)+string(data))
			if err != nil {
				//fmt.Printf("get failed, err:%v\n", err)
				SendChannel(err.Error())
				continue
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				//fmt.Printf("read from resp.Body failed, err:%v\n", err)
				SendChannel(err.Error())
				continue
			}
			SendChannel(string(body))
			resp.Body.Close()
			}

		}
	}
}
func SendChannel(str string){
	//fmt.Println("client ret:",str)
	SendConn.Write(protocol.Packet([]byte(str)))
}
var  SendConn net.Conn

func Init(num int) {
	total:=num
	fmt.Println("ClientInit.....")
	for {
		num--
		if(num<0){break}
		SendChannel, err := net.Dial("tcp", Remote+":"+strconv.Itoa(Snum))
		if err != nil {
			fmt.Fprintf(os.Stderr, "第%d次连接SendChannel失败\n"+err.Error(),(total-num))
			time.Sleep(time.Second)
			continue
		}
		ReceChannel, err := net.Dial("tcp", Remote+":"+strconv.Itoa(Rnum))
		if err != nil {
			fmt.Fprintf(os.Stderr, "第%d次连接ReceChannel失败\n"+err.Error(),(total-num))
			time.Sleep(time.Second)
			continue
		}
		SendConn=SendChannel
		fmt.Println("ClientInitSuccess!")
		Recchannel(ReceChannel)
		break
	}
}
func ClientInit(remote string,local int,rnum int,snum int){
	fmt.Println(remote,local,rnum,snum)
	Rnum=rnum
	Snum=snum
	Local=local
	Remote=remote
	Init(1)
}
func main() {
	fmt.Println("本地客户端 需要传入连接的公网IP、本地服务端口、远程写隧道、远程读隧道 格式:->127.0.0.1 8080 13000 14000")
	defer func() {
		err:=recover()
		if err !=nil{
			fmt.Println("参数缺少，启用默认参数")
			//ClientInit("127.0.0.1",8080,13000,14000)
			ClientInit("127.0.0.1",688,5080,6080)
		}
	}()
	remote :=os.Args[1]
	local :=os.Args[2]
	rnum:=os.Args[3]
	snum:=os.Args[4]
	Local,_=strconv.Atoi(local)
	Rnum,_=strconv.Atoi(rnum)
	Snum,_=strconv.Atoi(snum)
	ClientInit(remote,Local,Rnum,Snum)

}
