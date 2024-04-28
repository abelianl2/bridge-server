package server

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/abelianl2/bridge-server/config"
	"github.com/gin-gonic/gin"
	"github.com/sunjiangjun/xlog"
	"github.com/tidwall/gjson"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Service struct {
	db  *gorm.DB
	log *xlog.XLog
}

func NewService(config config.Config, log *xlog.XLog) *Service {
	dsn := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", config.DB.User, config.DB.Password, config.DB.Addr, config.DB.DbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()

	if err != nil {
		panic(err)
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	return &Service{
		db:  db,
		log: log,
	}
}

func (s *Service) SaveTx(ctx *gin.Context) {
	b, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		s.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}
	root := gjson.ParseBytes(b)
	fromNetwork := root.Get("from_network").String()
	fromAddress := root.Get("from_address").String()
	toNetwork := root.Get("to_network").String()
	toAddress := root.Get("to_address").String()
	d := Deposit{fromNetwork, fromAddress, toNetwork, toAddress}

	err = s.db.Create(d).Error
	if err != nil {
		s.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}
	s.Success(ctx, string(b), nil, ctx.Request.RequestURI)
}

func (s *Service) GetToAddress(ctx *gin.Context) {
	fromNetwork := ctx.Query("from_network")
	fromAddress := ctx.Query("from_address")
	d := Deposit{}
	err := s.db.Model(Deposit{}).Where("from_network=? and from_address=?", fromNetwork, fromAddress).First(&d).Error
	if err != nil {
		s.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}
	s.Success(ctx, ctx.Request.RequestURI, d, ctx.Request.RequestURI)
}

type Deposit struct {
	FromNetwork string `json:"from_network" gorm:"from_network"`
	FromAddress string `json:"from_address" gorm:"from_address"`
	ToNetwork   string `json:"to_network" gorm:"to_network"`
	ToAddress   string `json:"to_address" gorm:"to_address"`
}

func (d *Deposit) TableName() string {
	return "deposit"
}

const (
	SUCCESS = 0
	FAIL    = 1
)

func (s *Service) Success(c *gin.Context, req string, resp interface{}, path string) {
	req = strings.Replace(req, "\t", "", -1)
	req = strings.Replace(req, "\n", "", -1)
	if v, ok := resp.(string); ok {
		resp = strings.Replace(v, "\n", "", -1)
	}
	s.log.Printf("path=%v,req=%v,resp=%v\n", path, req, resp)
	mp := make(map[string]interface{})
	mp["code"] = SUCCESS
	mp["message"] = "ok"
	mp["data"] = resp
	c.JSON(200, mp)
}

func (s *Service) Error(c *gin.Context, req string, path string, err string) {
	req = strings.Replace(req, "\t", "", -1)
	req = strings.Replace(req, "\n", "", -1)
	s.log.Errorf("path=%v,req=%v,err=%v\n", path, req, err)
	mp := make(map[string]interface{})
	mp["code"] = FAIL
	mp["message"] = err
	mp["data"] = ""
	c.JSON(200, mp)
}
