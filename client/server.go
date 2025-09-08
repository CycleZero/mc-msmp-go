package client

func (c *MsmpClient) ServerStatus() {
	err := c.SendRequest("minecraft:server/status", nil)
	if err != nil {
		return
	}
}

func (c *MsmpClient) ServerSave() {
	err := c.SendRequest("minecraft:server/save", nil)
	if err != nil {
		return
	}
}

func (c *MsmpClient) ServerStop() {
	err := c.SendRequest("minecraft:server/stop", nil)
	if err != nil {
		return
	}
}

func (c *MsmpClient) ServerSystemMessage(message string) {
	param := map[string]string{
		"message": message,
	}
	err := c.SendRequest("minecraft:server/system_message", param)
	if err != nil {
		return
	}
}

func (c *MsmpClient) ServerSettingsGet(path string) {
	method := "minecraft:serversettings/" + path
	err := c.SendRequest(method, nil)
	if err != nil {
		return
	}
}

func (c *MsmpClient) ServerSettingsSet(path string, value interface{}) {
	method := "minecraft:serversettings/" + path + "/set"
	param := map[string]interface{}{
		"value": value,
	}
	err := c.SendRequest(method, param)
	if err != nil {
		return
	}
}
