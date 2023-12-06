package types

type Service struct {
	Name  string `db:"name" json:"name" validate:"required" description:"Service Name"`
	ID    string `db:"id" json:"id" validate:"required" description:"Service ID"`
	Icon  string `db:"icon" json:"icon" validate:"required" description:"Service Icon"`
	Token string `db:"token" json:"token" validate:"required" description:"Service API Token"`
}
