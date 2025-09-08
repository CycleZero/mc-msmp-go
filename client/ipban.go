package client

import "github.com/CycleZero/mc-msmp-go/dto/subdto"

func (c *MsmpClient) IpBansSet(bans []subdto.IpBanDTO) {
	err := c.SendRequest("minecraft:ip_bans/set", bans)
	if err != nil {
		return
	}
}

func (c *MsmpClient) IpBans() {
	err := c.SendRequest("minecraft:ip_bans", nil)
	if err != nil {
		return
	}
}

func (c *MsmpClient) IpBansAdd(ban subdto.IpBanDTO) {
	err := c.SendRequest("minecraft:ip_bans/add", ban)
	if err != nil {
		return
	}
}

func (c *MsmpClient) IpBansRemove(ip string) {
	param := map[string]string{
		"ip": ip,
	}
	err := c.SendRequest("minecraft:ip_bans/remove", param)
	if err != nil {
		return
	}
}

func (c *MsmpClient) IpBansClear() {
	err := c.SendRequest("minecraft:ip_bans/clear", nil)
	if err != nil {
		return
	}
}
