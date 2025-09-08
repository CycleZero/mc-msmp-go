package client

import "github.com/CycleZero/mc-msmp-go/dto/subdto"

func (c *MsmpClient) OperatorsSet(operators []subdto.OperatorDto) {
	err := c.SendRequest("minecraft:operators/set", operators)
	if err != nil {
		return
	}
}

func (c *MsmpClient) Operators() {
	err := c.SendRequest("minecraft:operators", nil)
	if err != nil {
		return
	}
}

func (c *MsmpClient) OperatorsAdd(operator subdto.OperatorDto) {
	err := c.SendRequest("minecraft:operators/add", operator)
	if err != nil {
		return
	}
}

func (c *MsmpClient) OperatorsRemove(player subdto.PlayerDto) {
	err := c.SendRequest("minecraft:operators/remove", player)
	if err != nil {
		return
	}
}

func (c *MsmpClient) OperatorsClear() {
	err := c.SendRequest("minecraft:operators/clear", nil)
	if err != nil {
		return
	}
}
