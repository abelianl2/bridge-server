package server

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/abelianl2/bridge-server/config"
	"github.com/gin-gonic/gin"
	uuid2 "github.com/google/uuid"
	"github.com/sigurn/crc16"
	"github.com/sunjiangjun/xlog"
	"github.com/tidwall/gjson"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Service struct {
	db     *gorm.DB
	log    *xlog.XLog
	config config.Config
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
		config: config,
		db:     db,
		log:    log,
	}
}

func (s *Service) NotifyTx(ctx *gin.Context) {
	uuid := ctx.Param("id")
	b, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		s.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}
	root := gjson.ParseBytes(b)
	hash := root.Get("hash").String()

	err = s.db.Model(&Deposit{}).Where("uuid=?", uuid).UpdateColumn("hash", hash).Error
	if err != nil {
		s.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}
	s.Success(ctx, string(b), nil, ctx.Request.RequestURI)
}

func (s *Service) SaveTxAndMemo(ctx *gin.Context) {
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
	amount := root.Get("amount").String()
	uuid := uuid2.New().String()
	d := Deposit{FromNetwork: fromNetwork, FromAddress: fromAddress, ToNetwork: toNetwork, ToAddress: toAddress, UUID: uuid}

	err = s.db.Create(d).Error
	if err != nil {
		s.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}

	m := Memo{
		Action:   "deposit",
		Protocol: "Mable",
		From:     d.FromAddress,                   //l1 from address
		To:       s.config.DepositContractAddress, //l1 to address
		Receipt:  d.ToAddress,                     //l2 mint address
		Value:    amount,                          // mint amount
	}

	bs, err := json.Marshal(m)
	if err != nil {
		s.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}
	// TODO 需要构造deposit txmemo
	memo := []byte(bs)

	// TODO 根据需求改动amount , 这里100000为0.01ABE
	//amount := "100000"

	var RESERVEDFIELD uint16 = 0xFFFF
	var PROTOCOLVERSION byte = 0x10

	table := crc16.MakeTable(crc16.CRC16_XMODEM)
	crcValue := crc16.Checksum(memo, table)

	var memoBytes []byte
	memoBytes = append(memoBytes, 0x00)                                             // 第1字节固定为0x00
	memoBytes = append(memoBytes, PROTOCOLVERSION)                                  // 第2字节为协议版本，暂定0x10
	memoBytes = append(memoBytes, byte(len(memo)>>8), byte(len(memo)&0xFF))         // 第3-4字节为长度
	memoBytes = append(memoBytes, byte(crcValue>>8), byte(crcValue&0xFF))           // 第5-6字节为CRC16校验码
	memoBytes = append(memoBytes, byte(RESERVEDFIELD>>8), byte(RESERVEDFIELD&0xFF)) // 第7-8字节预留字段
	memoBytes = append(memoBytes, memo...)

	/**
	{
	  "recipient": "abe338ce0ce178fb0aca42b4e400cdf395c92cbf9c5c9abd678aa516835f697bd6d280b285815924f862352c5463421c9f8d247f65dc112aa04c25de925bd1d1a334",
	  "amountOfAbel": "1.2111",
	  "amountOfGasFee": "0.09",
	  "memo": "001001ade57affff7b0a09202022616374696f6e223a20226465706f736974222c0a0920202270726f746f636f6c223a20224d61626c65222c0a0920202266726f6d223a2022616265333266356339646436376236663065313133333366633534653462353464316630353435366561306532616263366531343539623035363237316533646536313830663763636134636138383061383833396337326434313239383766666434376437666463613630666365353833386266636265613638646437343131343662222c0a09202022746f223a2022616265333338636530636531373866623061636134326234653430306364663339356339326362663963356339616264363738616135313638333566363937626436643238306232383538313539323466383632333532633534363334323163396638643234376636356463313132616130346332356465393235626431643161333334222c0a0920202272656365697074223a2022307864616331376639353864326565353233613232303632303639393435393763313364383331656337222c0a0920202276616c7565223a202230783231323232323030220a097d",
	  "hook":"www.google.com/0x3b1b43e42bfd2123a5a4a949b4207af992ed46d49c5ef271fb5aae9a126e5629"
	}

	*/

	callBack := CallBack{
		Recipient:      s.config.DepositContractAddress,
		AmountOfGasFee: "0.09",
		AmountOfAbel:   amount,
		Memo:           hex.EncodeToString(memoBytes),
		Hook:           fmt.Sprintf("%v/%v", s.config.HookUri, uuid),
	}

	s.Success(ctx, string(b), callBack, ctx.Request.RequestURI)
}

type Memo struct {
	Protocol string `json:"protocol" gorm:"column:protocol"`
	Action   string `json:"action" gorm:"column:action"`
	From     string `json:"from" gorm:"column:from"`
	Receipt  string `json:"receipt" gorm:"column:receipt"`
	To       string `json:"to" gorm:"column:to"`
	Value    string `json:"value" gorm:"column:value"`
}

type CallBack struct {
	AmountOfGasFee string `json:"amountOfGasFee" gorm:"column:amountOfGasFee"`
	Hook           string `json:"hook" gorm:"column:hook"`
	Recipient      string `json:"recipient" gorm:"column:recipient"`
	Memo           string `json:"memo" gorm:"column:memo"`
	AmountOfAbel   string `json:"amountOfAbel" gorm:"column:amountOfAbel"`
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
	hash := root.Get("hash").String()
	uuid := uuid2.New().String()
	d := Deposit{FromNetwork: fromNetwork, FromAddress: fromAddress, ToNetwork: toNetwork, ToAddress: toAddress, Hash: hash, UUID: uuid}

	err = s.db.Create(d).Error
	if err != nil {
		s.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}
	s.Success(ctx, string(b), d, ctx.Request.RequestURI)
}

func (s *Service) GetToAddress(ctx *gin.Context) {
	hash := ctx.Query("hash")
	d := Deposit{}
	err := s.db.Model(Deposit{}).Where("hash=?", hash).First(&d).Error
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
	UUID        string `json:"uuid" gorm:"uuid"`
	Hash        string `json:"hash" gorm:"hash"`
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
