package infra

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"xiyu.com/common"
)

type UserDo struct {
	bun.BaseModel `bun:"table:tao_market"`
	Id            int64     `bun:"id,autoincrement"`
	OpenId        string    `bun:"openid,notnull"`
	UserType      string    `bun:"usertype,notnull"`
	Nickname      string    `bun:"nickname,notnull"`
	Avatar        string    `bun:"avatar,notnull"`
	CreateTime    time.Time `bun:"createtime,notnull,default:current_timestamp"`
	UpdateTime    time.Time `bun:"updatetime,notnull,default:current_timestamp"`
}

type TaoDb struct {
	common.TAutoCloseable
	db  *bun.DB
	ctx *context.Context
}

func (dbCon *TaoDb) Close() {
	if dbCon.db != nil {
		err := dbCon.db.Close()
		if err != nil {
			return
		}
	}
}

var dbCon *TaoDb
var dbOnce sync.Once

func GetDbCon() *TaoDb {
	dbOnce.Do(func() {
		dbCon = &TaoDb{}
	})
	return dbCon
}

func (dbCon *TaoDb) Connect(ctx *context.Context) error {
	dsn := common.GlbBaInfa.Conf.Infra.DbDns
	sqlDb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	dbCon.db = bun.NewDB(sqlDb, pgdialect.New())
	dbCon.ctx = ctx
	common.Logger.Info("connect to:%s success", dsn)
	common.TaddItem(dbCon)
	return nil
}

func FindUser() {

}
