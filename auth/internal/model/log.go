package model

const (
	LogActionCreate = "create"
	LogActionUpdate = "update"
	LogActionDelete = "delete"
)

type UserLog struct {
	UserID int64
	Action string
}
