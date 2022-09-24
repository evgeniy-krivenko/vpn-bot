package entity

type Connection struct {
	Id              int    `db:"id"`
	Location        string `db:"location"`
	Port            uint   `db:"port"`
	UserId          int64  `db:"user_id"`
	EncryptedSecret string `db:"encrypted_secret"`
	IpAddress       string `db:"ip_address"`
	ServerId        int    `db:"server_id"`
}

type ConnectionPortCount struct {
	Port  int `db:"port"`
	Count int `db:"count"`
}
