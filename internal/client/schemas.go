package client

type ProxyClient struct {
	TgID        int64  `json:"tg_id"`
	SubID       string `json:"sub_id"`
	Email       string `json:"email"`
	Enable      bool   `json:"enable"`
	LimitIP     int    `json:"limit_ip"`
	LimitHwid   int    `json:"limit_hwid"`
	ExpiryTime  int    `json:"expiry_time"`
	UsedTraffic int    `json:"used_traffic"`
}
