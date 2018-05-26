package helper

import (
	"database/sql"

	"github.com/miklly/miklly/config"

	_ "github.com/mattn/go-sqlite3"
)

type FillModelFunc func(fileds map[string][]byte)
type DataBase struct {
	db *sql.DB
}

func OpenDataBase() *DataBase {
	db, err := sql.Open(config.DBType, config.DBFile)
	if err != nil {
		config.Log("打开数据库失败!", config.DBType, config.DBFile, err)
		panic(err)
	}
	return &DataBase{
		db: db,
	}
}
func (this *DataBase) BeginTranslation() (*sql.Tx, error) {
	db := OpenDataBase()
	this.db = db.db
	return this.db.Begin()
}

//添加记录,返回主键ID
func (db *DataBase) Insert(sql string, pars ...interface{}) (id int64) {
	stmt, err := db.db.Prepare(sql)
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
func (db *DataBase) Modify(sql string, pars ...interface{}) (num int64) {
	stmt, err := db.db.Prepare(sql)
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
func (db *DataBase) Query(sql string, fillModel FillModelFunc, pars ...interface{}) bool {
	rows, err := db.db.Query("SELECT * FROM userinfo")
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
