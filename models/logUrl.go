package models

import "time"

type LogUrl struct {
	Id        int       `gorm:"primarykey" json:"id"`
	Created   time.Time `json:"created"`
	Url       string    `json:"url"`
	IpAddress string    `json:"ip_address"`
	Status    string    `json:"status"`
}
