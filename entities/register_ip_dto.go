package entities

type OneCode struct {
	Code  int
	Times int
}

type OneIp struct {
	Ip         string
	CodeStatus []*OneCode
}

type StatisticDTO struct {
	Ips []*OneIp
}
