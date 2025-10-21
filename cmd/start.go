package cmd

import (
	"fmt"
	"github.com/CycleZero/mc-msmp-go"
	"log"
	"sync"
	"time"
)

func Start() {

	url := "ws://msmp.server.poyuan233.cn:8088"
	secret := "MjHrY9yN3WTUKXsgtB1bMxTtvWlnJwVAVEbLFT2z"
	cli := mcmsmpgo.NewMsmpClient(url, secret)
	err := cli.Connect()
	if err != nil {
		log.Println(err)
		return
	}
	defer func(cli *mcmsmpgo.MsmpClient) {
		err := cli.Disconnect()
		if err != nil {
			panic(err)
		}
	}(cli)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		i := 0
		for {
			fmt.Println("i:", i)
			if i >= 20 {
				break
			}
			i++
			cli.ServerStatus()
			time.Sleep(5 * time.Second)
		}
		wg.Done()
	}()
	wg.Wait()
	cli.AllowlistSet("8484", "wdwd")
	log.Println("===========end===========")
}
