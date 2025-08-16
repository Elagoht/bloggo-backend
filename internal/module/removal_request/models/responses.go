package models

import "time"

type RemovalRequestCard struct {
	Id            int64      `json:"id"`
	PostVersionId int64      `json:"postVersionId"`
	PostTitle     string     `json:"postTitle"`
	RequestedBy   UserInfo   `json:"requestedBy"`
	Note          *string    `json:"note"`
	Status        int64      `json:"status"`
	DecidedBy     *UserInfo  `json:"decidedBy"`
	DecidedAt     *time.Time `json:"decidedAt"`
	CreatedAt     time.Time  `json:"createdAt"`
}

type RemovalRequestDetails struct {
	Id            int64      `json:"id"`
	PostVersionId int64      `json:"postVersionId"`
	PostTitle     string     `json:"postTitle"`
	PostContent   string     `json:"postContent"`
	RequestedBy   UserInfo   `json:"requestedBy"`
	Note          *string    `json:"note"`
	Status        int64      `json:"status"`
	DecidedBy     *UserInfo  `json:"decidedBy"`
	DecidedAt     *time.Time `json:"decidedAt"`
	CreatedAt     time.Time  `json:"createdAt"`
}

type UserInfo struct {
	Id     int64   `json:"id"`
	Name   string  `json:"name"`
	Avatar *string `json:"avatar"`
}