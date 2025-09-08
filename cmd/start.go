package cmd

import "github.com/CycleZero/mc-msmp-go/client"

func Start() {

	url := "ws://localhost:25576/ws"

	cli := client.NewMsmpClient(url)
	err := cli.Connect()
	if err != nil {
		return
	}
	defer func(cli *client.MsmpClient) {
		err := cli.Disconnect()
		if err != nil {
			panic(err)
		}
	}(cli)

	cli.AllowlistSet("8484", "wdwd")

}
