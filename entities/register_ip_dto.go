package entities

// one code struct
type OneCode struct {
	Code  int
	Times int
}

// another ip
type OneIp struct {
	Ip         string
	CodeStatus []*OneCode
}

// stat strct
type StatisticDTO struct {
	Ips []*OneIp
}
