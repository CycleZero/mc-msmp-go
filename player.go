package mcmsmpgo

import "github.com/CycleZero/mc-msmp-go/dto/subdto"

func (c *MsmpClient) Players() {
	err := c.SendRequest("minecraft:players", nil)
	if err != nil {
		return
	}
}

func (c *MsmpClient) PlayersKick(player subdto.PlayerDto, reason string) {
	param := map[string]interface{}{
		"player": player,
		"reason": reason,
	}
	err := c.SendRequest("minecraft:players/kick", param)
	if err != nil {
		return
	}
}
