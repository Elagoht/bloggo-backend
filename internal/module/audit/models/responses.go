package models

import "time"

type ResponseAuditLog struct {
	Id        int64     `json:"id"`
	UserId    *int64    `json:"userId"`
	Entity    string    `json:"entity"`
	EntityId  int64     `json:"entityId"`
	Action    string    `json:"action"`
	CreatedAt time.Time `json:"createdAt"`
}