package main

type arguments struct {
	config string
}

type config struct {
	Monitor monitorSettings `json:"monitor_settings"`
	IpNet   IpNetSettings   `json:"ipnet_settings"`
}

type monitorSettings struct {
	Device   string `json:"device"`
	Baudrate int    `json:"baudrate"`
}

type IpNetSettings struct {
	Pattern string `json:"pattern"`
	Fuzzy   bool   `json:"fuzzy"`
}
