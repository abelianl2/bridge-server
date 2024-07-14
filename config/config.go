package config

/**

{
  "RootPath": "/api/chain",
  "Port": 9002,
  "DB": {
    "addr": "192.168.75.18",
    "user": "liang",
    "password": "tyjBweA3++y18kJQ",
    "dbName": "bigdata_base"
  }
}
*/

type Config struct {
	RootPath               string `json:"RootPath" gorm:"column:RootPath"`
	Port                   int    `json:"Port" gorm:"column:Port"`
	Https                  bool   `json:"Https" gorm:"column:Https"`
	DepositContractAddress string `json:"DepositContractAddress" gorm:"column:DepositContractAddress"`
	HookUri                string `json:"HookUri" gorm:"column:HookUri"`
	DepositUri             string `json:"DepositUri" gorm:"column:DepositUri"`
	DB                     DB     `json:"DB" gorm:"column:DB"`
	LogLevel               int    `json:"LogLevel"`
}

type DB struct {
	Password string `json:"password" gorm:"column:password"`
	DbName   string `json:"dbName" gorm:"column:dbName"`
	Addr     string `json:"addr" gorm:"column:addr"`
	User     string `json:"user" gorm:"column:user"`
}
