package client

import "github.com/CycleZero/mc-msmp-go/dto/subdto"

func (c *MsmpClient) Gamerules() {
	err := c.SendRequest("minecraft:gamerules", nil)
	if err != nil {
		return
	}
}

func (c *MsmpClient) GamerulesUpdate(rules []subdto.TypedRule) {
	err := c.SendRequest("minecraft:gamerules/update", rules)
	if err != nil {
		return
	}
}
