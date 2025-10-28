package test

import (
	"fmt"
	mcmsmpgo "github.com/CycleZero/mc-msmp-go"
	"github.com/CycleZero/mc-msmp-go/dto"
	"log"
	"strconv"
	"sync"
	"testing"
	"time"
)

const MaxReq int = 1024 * 1024

var ResultMap = make(map[string]bool)
var ResultMapLock sync.Mutex
var errStr = ""

func PrintNonResponseNum() {
	l := make([]string, 10)
	ResultMapLock.Lock()
	for i := 1; i <= MaxReq; i++ {
		id := strconv.FormatInt(int64(i), 10)
		if ResultMap[id] == false {
			l = append(l, id)
		}
	}
	ResultMapLock.Unlock()
	fmt.Println("未响应的请求 : ", l)
	fmt.Println("未响应数: ", len(l))
}
func TempHandler(request *dto.MsmpRequest, response dto.MsmpResponse) {
	//rq, _ := json.Marshal(request)
	//rs, _ := json.Marshal(response)
	//fmt.Println("request:", string(rq))
	//fmt.Println("response:", string(rs))
	id := request.ID
	ResultMapLock.Lock()
	if ResultMap[strconv.FormatInt(int64(id), 10)] == false {
		ResultMap[strconv.FormatInt(int64(id), 10)] = true

	} else if ResultMap[strconv.FormatInt(int64(id), 10)] == true {
		errStr += "\n重复请求 : " + strconv.FormatInt(int64(id), 10)
	} else {
		errStr += "\n请求错误 : " + strconv.FormatInt(int64(id), 10)
	}
	ResultMapLock.Unlock()
}

func Test(t *testing.T) {
	url := "ws://msmp.server.poyuan233.cn:8088"
	secret := "MjHrY9yN3WTUKXsgtB1bMxTtvWlnJwVAVEbLFT2z"
	cli := mcmsmpgo.NewMsmpClient(url, secret, &mcmsmpgo.NewClientConfig{
		Handler:       TempHandler,
		Container:     nil,
		AutoReconnect: true,
	})

	err := cli.Connect()
	if err != nil {
		log.Println(err)
		return
	}
	//c := jsonrpc.NewClient(cli.Conn.NetConn())
	defer func(cli *mcmsmpgo.MsmpClient) {
		err := cli.Disconnect()
		if err != nil {
			panic(err)
		}
	}(cli)

	var wg sync.WaitGroup

	goRoutineNum := 32
	for r := 0; r < goRoutineNum; r++ {
		wg.Add(1)
		go func(rid int) {
			i := 0
			total := MaxReq / goRoutineNum
			for {
				//if i%1024* == 0 {
				//	fmt.Println("i:", i, "进度:", i/1024, "/", total/1024)
				//}
				if i > total {
					break
				}
				i++
				cli.ServerStatus()
				//var reply string
				//c.Call("ServerStatus", nil, &reply)
				//time.Sleep(1 * time.Millisecond)
			}
			fmt.Println("================发送完毕==================")
			wg.Done()
		}(r)
	}

	wg.Wait()
	fmt.Println("================进入接收阶段==================")
	for {
		log.Println("错误信息 : " + errStr)
		PrintNonResponseNum()
		time.Sleep(1000 * time.Millisecond)
	}
	//cli.AllowlistSet("8484", "wdwd")
	log.Println("===========end===========")
	log.Println("错误信息 : " + errStr)
}
