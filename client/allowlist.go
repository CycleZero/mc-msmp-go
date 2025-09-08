package client

import (
	"github.com/CycleZero/mc-msmp-go/dto/subdto"
)

func (c *MsmpClient) AllowlistSet(id string, name string) {
	param := []subdto.PlayerDto{
		subdto.PlayerDto{
			Id:   id,
			Name: name,
		},
	}
	err := c.SendRequest("minecraft:allowlist/set", param)
	if err != nil {
		return
	}
}

func (c *MsmpClient) Allowlist() {
	err := c.SendRequest("minecraft:allowlist", nil)
	if err != nil {
		return
	}
}

func (c *MsmpClient) AllowlistAdd(id string, name string) {
	param := subdto.PlayerDto{
		Id:   id,
		Name: name,
	}
	err := c.SendRequest("minecraft:allowlist/add", param)
	if err != nil {
		return
	}
}

func (c *MsmpClient) AllowlistRemove(id string, name string) {
	param := subdto.PlayerDto{
		Id:   id,
		Name: name,
	}
	err := c.SendRequest("minecraft:allowlist/remove", param)
	if err != nil {
		return
	}
}

func (c *MsmpClient) AllowlistClear() {
	err := c.SendRequest("minecraft:allowlist/clear", nil)
	if err != nil {
		return
	}
}
