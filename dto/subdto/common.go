package subdto

type PlayerDto struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type UserBanDto struct {
	Expires string    `json:"expires"`
	Player  PlayerDto `json:"player"`
	Reason  string    `json:"reason"`
	Source  string    `json:"source"`
}

type IpBanDTO struct {
	Expires string `json:"expires"`
	Ip      string `json:"ip"`
	Reason  string `json:"reason"`
	Source  string `json:"source"`
}

type OperatorDto struct {
	BypassesPlayerLimit bool      `json:"bypassesPlayerLimit"`
	PermissionLevel     int       `json:"permissionLevel"`
	Player              PlayerDto `json:"player"`
}

type Version struct {
	Name     string `json:"name"`
	Protocol int    `json:"protocol"`
}
type ServerState struct {
	Player  []PlayerDto `json:"player"`
	Started bool        `json:"started"`
	Version Version     `json:"version"`
}

type TypedRule struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}
