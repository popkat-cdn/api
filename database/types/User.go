package types

type User struct {
	Username string `db:"username" json:"username" validate:"required" description:"Username"`
	UserID   string `db:"userid" json:"user_id" validate:"required" description:"User ID"`
	Banned   bool   `db:"banned" type:"boolean" json:"banned" validate:"required" description:"Is the user banned?"`
	Avatar   string `db:"avatar" json:"avatar" validate:"required" description:"User Avatar"`
	Token    string `db:"token" json:"token" validate:"required" description:"User Auth Token (temp)"`
}
