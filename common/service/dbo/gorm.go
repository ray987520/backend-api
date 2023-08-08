package dbo

import (
	"GamePoolApi/common/enum/innertrace"
	"GamePoolApi/common/enum/thirdparty"
	iface "GamePoolApi/common/interface"
	"GamePoolApi/common/service/serializer"
	"GamePoolApi/common/service/tracer"
	"GamePoolApi/common/service/zaplog"
	"fmt"
	"time"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

/////////////////////////////
////    SQL DB服務(Gorm)
/////////////////////////////

type GormDB struct {
}

const (
	connectionError = "gorm open connection error:%v"  //open sql connection error
	dbInstanceError = "gorm get sql instance error:%v" //get sql connection pool error
	dbStat          = "gorm db stat:%v"                //sql db stat message
)

var (
	sqlDB            *gorm.DB      //gorm instance
	sqlConnectString string        //gorm sqlConnectString
	maxOpenConns     int           //gorm max open connection,max pool size
	maxIdleConns     int           //gorm max idle connection
	maxIdleSecond    time.Duration //grom max idle time,second
)

// 取GormDB實例
func GetSqlDb(cfg iface.IAppConfig) *GormDB {
	gormInit(cfg)
	return &GormDB{}
}

// 初始化,建立sql db連線與實例
func gormInit(cfg iface.IAppConfig) {
	sqlConnectString = cfg.GetSqlConnectString()
	maxOpenConns = cfg.GetSqlMaxOpenConns()
	maxIdleConns = cfg.GetSqlMaxIdleConns()
	maxIdleSecond = cfg.GetSqlMaxIdleSecond()
	//gorm連接sql server
	db, err := gorm.Open(sqlserver.Open(sqlConnectString), &gorm.Config{})
	db.Use(dbresolver.Register(dbresolver.Config{
		Sources:           []gorm.Dialector{sqlserver.Open(sqlConnectString)}, //主host
		Replicas:          []gorm.Dialector{sqlserver.Open(sqlConnectString)}, //副本host
		Policy:            dbresolver.RandomPolicy{},
		TraceResolverMode: true,
	}).
		SetMaxOpenConns(maxOpenConns).
		SetMaxIdleConns(maxIdleConns).
		SetConnMaxIdleTime(maxIdleSecond * time.Second))
	//連接失敗,panic error
	if err != nil {
		err = fmt.Errorf(connectionError, err)
		zaplog.Errorw(innertrace.ExternalServiceError, innertrace.FunctionNode, thirdparty.GormInit, innertrace.TraceNode, tracer.DefaultTraceId, innertrace.DataNode, err)
		panic(err)
	}

	//嘗試取出DB的connection pool,錯誤就是連線但連不上DB
	sqlDb, err := db.DB()
	//取得connection pool異常,panic error
	if err != nil {
		err = fmt.Errorf(dbInstanceError, err)
		zaplog.Errorw(innertrace.ExternalServiceError, innertrace.FunctionNode, thirdparty.GormInit, innertrace.TraceNode, tracer.DefaultTraceId, innertrace.DataNode, err)
		panic(err)
	}

	//初始化後列印DB狀態
	data := serializer.JsonMarshal(tracer.DefaultTraceId, sqlDb.Stats())
	zaplog.Infow(innertrace.InfoNode, innertrace.FunctionNode, thirdparty.GormInit, innertrace.TraceNode, tracer.DefaultTraceId, innertrace.DataNode, fmt.Sprintf(dbStat, string(data)))

	//把connection傳給全域變數
	sqlDB = db
}

// sql raw call stored procedure,輸出生效行數
func (gormDB *GormDB) CallSP(traceCode string, model interface{}, sqlString string, params ...interface{}) (rowsAffect int64) {
	s := sqlDB.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Raw(sqlString, params...)
	})
	fmt.Println(s)
	tx := sqlDB.Raw(sqlString, params...).Scan(model)
	if tx.Error != nil {
		zaplog.Errorw(innertrace.ExternalServiceError, innertrace.FunctionNode, thirdparty.SqlSelect, innertrace.TraceNode, traceCode, innertrace.DataNode, tx.Error)
		return -1
	}
	return tx.RowsAffected
}

// sql raw執行select,輸出生效行數
func (gormDB *GormDB) Select(traceCode string, model interface{}, sqlString string, params ...interface{}) (rowsAffect int64) {
	tx := sqlDB.Raw(sqlString, params...).Scan(model)
	if tx.Error != nil {
		zaplog.Errorw(innertrace.ExternalServiceError, innertrace.FunctionNode, thirdparty.SqlSelect, innertrace.TraceNode, traceCode, innertrace.DataNode, tx.Error)
		return -1
	}
	return tx.RowsAffected
}

// sql raw執行update,輸出生效行數
func (gormDB *GormDB) Update(traceCode string, sqlString string, params ...interface{}) (rowsAffect int64) {
	tx := sqlDB.Exec(sqlString, params...)
	if tx.Error != nil {
		zaplog.Errorw(innertrace.ExternalServiceError, innertrace.FunctionNode, thirdparty.SqlUpdate, innertrace.TraceNode, traceCode, innertrace.DataNode, tx.Error)
		return -1
	}
	return tx.RowsAffected
}

// sql raw執行delete,輸出生效行數
func (gormDB *GormDB) Delete(traceCode string, sqlString string, params ...interface{}) (rowsAffect int64) {
	tx := sqlDB.Exec(sqlString, params...)
	if tx.Error != nil {
		zaplog.Errorw(innertrace.ExternalServiceError, innertrace.FunctionNode, thirdparty.SqlDelete, innertrace.TraceNode, traceCode, innertrace.DataNode, tx.Error)
		return -1
	}
	return tx.RowsAffected
}

// sql raw執行insert,輸出生效行數
func (gormDB *GormDB) Create(traceCode string, sqlString string, params ...interface{}) (rowsAffect int64) {
	tx := sqlDB.Exec(sqlString, params...)
	if tx.Error != nil {
		zaplog.Errorw(innertrace.ExternalServiceError, innertrace.FunctionNode, thirdparty.SqlCreate, innertrace.TraceNode, traceCode, innertrace.DataNode, tx.Error)
		return -1
	}
	return tx.RowsAffected
}

// GORM自動transaction,sql raw執行,輸出生效行數
func (gormDB *GormDB) Transaction(traceCode string, sqlStrings []string, params ...[]interface{}) (rowsAffect int64) {
	err := sqlDB.Transaction(func(tx *gorm.DB) error {
		//循序執行所有sql,累加rowsAffect,如果有錯誤gorm會rollback,正常結束會自動commit
		for i, sql := range sqlStrings {
			if params != nil {
				partwork := tx.Exec(sql, params[i]...)
				//部分執行失敗,返回錯誤
				if partwork.Error != nil {
					zaplog.Errorw(innertrace.ExternalServiceError, innertrace.FunctionNode, thirdparty.SqlTransaction, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage(innertrace.ErrorInfoNode, partwork.Error, "sql", sql, "params", params[i]))
					return partwork.Error
				}

				//累加生效行數
				rowsAffect += partwork.RowsAffected
			} else {
				partwork := tx.Exec(sql)
				//部分執行失敗,返回錯誤
				if partwork.Error != nil {
					zaplog.Errorw(innertrace.ExternalServiceError, innertrace.FunctionNode, thirdparty.SqlTransaction, innertrace.TraceNode, traceCode, innertrace.DataNode, tracer.MergeMessage(innertrace.ErrorInfoNode, partwork.Error, "sql", sql))
					return partwork.Error
				}

				//累加生效行數
				rowsAffect += partwork.RowsAffected
			}
		}

		return nil
	})
	//如果grom transaction有異常就返回無更新筆數,因為transaction已經log不再加log
	if err != nil {
		return -1
	}

	return rowsAffect
}
