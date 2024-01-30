package model

type Contact struct {
	UserID      uint64 `json:"userid,omitempty"`
	ContactID   uint64 `json:"contactid,omitempty"`
	ContactName string `json:"contactname,omitempty"`
}
