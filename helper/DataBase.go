package helper

import (
	"database/sql"

	"github.com/miklly/miklly/config"

	_ "github.com/mattn/go-sqlite3"
)

type FillModelFunc func(fileds map[string][]byte)
type DataBase struct {
	db   *sql.DB
	tran *sql.Tx
}

var queryDB *DataBase

func OpenDataBase() *DataBase {
	tm := GetTaskManager()
	if queryDB != nil {
		tm.RestTask("DataBase_Close")
		return queryDB
	}
	db, err := sql.Open(config.DBType, config.DBFile)
	if err != nil {
		config.Log("打开数据库失败!", config.DBType, config.DBFile, err)
		panic(err)
	}
	queryDB := &DataBase{
		db: db,
	}
	err = tm.AddTaskByNow(30, "DataBase_Close", func(args ...interface{}) {
		db := DataBase(args[0])
		db.Close()
	}, queryDB)
	if err != nil {
		config.Log("创建数据库延时关闭失败!", err)
	}
	return queryDB
}
func (this *DataBase) BeginTranslation() (*sql.Tx, error) {
	db := OpenDataBase()
	this.db = db.db
	return this.db.Begin()
}
func (this *DataBase) Commint() bool {
	if this.tran == nil {
		return true
	}
	return !config.HasErr(this.tran.Commit(), "提交事务失败!")
}
func (this *DataBase) Rollback() bool {
	if this.tran == nil {
		return true
	}
	return !config.HasErr(this.tran.Rollback(), "回滚事务失败!")
}
func (this *DataBase) Close() {
	if this.db != nil {
		this.db.Close()
	}
}
func (this *DataBase) getPrepare(sql string) (*sql.Stmt, error) {
	if this.tran != nil {
		return this.tran.Prepare(sql)
	}
	return this.db.Prepare(sql)
}

//添加记录,返回主键ID
func (this *DataBase) Insert(sql string, pars ...interface{}) (id int64) {
	stmt, err := this.getPrepare(sql)
	if err == nil {
		res, err := stmt.Exec(pars)
		if err == nil {
			id, err := res.LastInsertId()
			if err == nil {
				return id
			}
		}
	}
	config.Log("插入数据失败!", sql, pars, err)
	return -1
}

//修正记录:更新或删除,返回修正的记录数目
func (this *DataBase) Modify(sql string, pars ...interface{}) (num int64) {
	stmt, err := this.getPrepare(sql)
	if err == nil {
		res, err := stmt.Exec(pars)
		if err == nil {
			num, err := res.RowsAffected()
			if err == nil {
				return num
			}
		}
	}
	config.Log("修正数据失败!", sql, pars, err)
	return -1
}

//查询记录,返回是否查询成功.
func (this *DataBase) Query(sql string, fillModel FillModelFunc, pars ...interface{}) bool {
	rows, err := this.db.Query(sql, pars)
	defer rows.Close()
	if err == nil {
		cols, err := rows.Columns()
		if err != nil {
			config.Log("获取列名列表失败!", sql)
			return false
		}
		vals := make([]sql.RawBytes, len(cols))
		scanArgs := make([]interface{}, len(vals))
		for i := range vals {
			scanArgs[i] = &vals[i]
		}
		kv := make(map[string][]byte)
		for rows.Next() {
			err == rows.Scan(scanArgs...)
			if err != nil {
				config.Log("提取查询数据行结果失败!", sql, pars)
			}
			for i := range vals {
				kv[cols[i]] = vals[i]
			}
			FillModelFunc(kv)
		}
		return true
	}
	config.Log("查询数据失败!", sql, pars, err)
	return false
}
