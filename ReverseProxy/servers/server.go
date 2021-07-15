package main

import (
	"fmt"
	"os"
	"strconv"
)
func main()  {
	fmt.Println("服务端 需要传入开放的端口范围格式：->80 90")
	defer func() {
		err:=recover()
		if err !=nil{
			fmt.Println("参数缺少，启用默认参数")
			PublicInitRange(8080,8080)
		}
	}()
	PnumF:=os.Args[1]
	PnumT:=os.Args[2]
	res1,_:=strconv.Atoi(PnumF)
	res2,_:=strconv.Atoi(PnumT)
	 PublicInitRange(res1,res2)
}
