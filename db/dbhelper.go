package db

import (
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"xorm.io/core"

	"helper/common/errno"
	"helper/common/logger"
)

type DbConnConf struct {
	DriverName      string
	ConnMaxLifetime int64
	Prefix          string
	Conn
}
type Conn interface {
	GetDataSourceName() string
}

/**
 * mysql
 */
type MysqlConf struct {
	Host     string
	Username string
	Password string
	Database string
}

func (c *MysqlConf) GetDataSourceName() (dataSourceName string) {
	dataSourceName = c.Username + ":" + c.Password + "@(" + c.Host + ")/" + c.Database + "?charset=utf8mb4&loc=Local"
	return
}

var ConfIns *DbConnConf

func (d *DbConnConf) InitDbEngine() {
	ConfIns = d
	SetDbEngine()
}

var dbEngineIns *xorm.Engine
var SetDbEngineOnce sync.Once

func NewEngineInstance() *xorm.Engine {
	SetDbEngineOnce.Do(SetDbEngine)
	return dbEngineIns
}

func SetDbEngine() {
	if ConfIns == nil {
		logger.Instance.WithField("code", errno.ErrConfig).Panic("db config is nil")
	}

	dbEngine, err := xorm.NewEngine(ConfIns.DriverName, ConfIns.Conn.GetDataSourceName())
	if err != nil {
		logger.Instance.WithField("code", errno.ErrDatabase).Panicf("new db engine error: %v", err)
	}

	if ConfIns.ConnMaxLifetime > 0 {
		dbEngine.DB().SetConnMaxLifetime(time.Duration(ConfIns.ConnMaxLifetime) * time.Second)
	}

	if ConfIns.Prefix != "" {
		tbMapper := core.NewPrefixMapper(core.SnakeMapper{}, ConfIns.Prefix)
		dbEngine.SetTableMapper(tbMapper)
	}

	dbEngineIns = dbEngine
}
