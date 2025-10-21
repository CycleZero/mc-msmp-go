package mcmsmpgo

import "github.com/CycleZero/mc-msmp-go/dto/subdto"

func (c *MsmpClient) BansSet(bans []subdto.UserBanDto) {
	err := c.SendRequest("minecraft:bans/set", bans)
	if err != nil {
		return
	}
}

func (c *MsmpClient) Bans() {
	err := c.SendRequest("minecraft:bans", nil)
	if err != nil {
		return
	}
}

func (c *MsmpClient) BansAdd(ban subdto.UserBanDto) {
	err := c.SendRequest("minecraft:bans/add", ban)
	if err != nil {
		return
	}
}

func (c *MsmpClient) BansRemove(player subdto.PlayerDto) {
	err := c.SendRequest("minecraft:bans/remove", player)
	if err != nil {
		return
	}
}

func (c *MsmpClient) BansClear() {
	err := c.SendRequest("minecraft:bans/clear", nil)
	if err != nil {
		return
	}
}
