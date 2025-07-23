package dbcache

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql" // 导入MySQL驱动
)

type DbCache struct {
	opt      *Option
	db       *sql.DB
	Store    *Store
	isClosed bool
}
type Store struct {
	InsertList []string
	RemoveList []string
	UpdateList []string
	SQLList    []string
}
type Option struct {
	DbType string
	DSN    string
	KEY    string
	COLUMN string
	Prefix string
}

func NewStorage(dbType string, dsn string) (db *DbCache) {
	// 连接数据库
	cache := &DbCache{
		opt: &Option{
			DbType: dbType,
			DSN:    dsn,
			Prefix: "tb",
			KEY:    "k",
			COLUMN: "v",
		},
	}
	cache.Store = &Store{
		InsertList: make([]string, 0),
		RemoveList: make([]string, 0),
		UpdateList: make([]string, 0),
	}
	go cache.doTask()
	return cache
}
func (c *DbCache) AutoOpen() (bool, error) {
	var err error
	if !c.isClosed {
		if c.db != nil {
			err = c.db.Ping()
			if err == nil {
				return true, nil
			}
		}
	}
	//mysql,"username:password@tcp(localhost:3306)/dbname"
	c.db, err = sql.Open(c.opt.DbType, c.opt.DSN)
	if err != nil {
		return false, err
	}
	err = c.db.Ping()
	if err != nil {
		return false, err
	}
	c.isClosed = false
	return true, nil
}
func (c *DbCache) table(key string) string {
	// 生成表名
	tb := fmt.Sprintf("%s_%s", c.opt.Prefix, key)
	return tb
}
func (c *DbCache) create(tb string) {
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s VARCHAR(256) PRIMARY KEY,%s text)", tb, c.opt.KEY, c.opt.COLUMN)
	_, err := c.exec(sql)
	if err != nil {
		fmt.Println(err.Error(), sql)
	}
}

func (c *DbCache) Delete(key string, value string) {
	c.Store.SQLList = append(c.Store.SQLList, c.get_sql(key, value, DELETE))
}
func (c *DbCache) Add(key string, value string) {
	c.Store.SQLList = append(c.Store.SQLList, c.get_sql(key, value, INSERT))
}

func (c *DbCache) Set(key string, value string) {
	c.Store.SQLList = append(c.Store.SQLList, c.get_sql(key, value, UPDATE))
}

func removeAt(slice []string, index int) []string {
	if index < 0 || index >= len(slice) {
		return slice
	}
	return append(slice[:index], slice[index+1:]...)
}

func (c *DbCache) doTask() {
	var (
		max = 1
	)
	var sqls string
	for i, sql := range c.Store.SQLList {
		sqls += sql
		if i > max {
			ok, err := c.exec(sqls)
			if !ok {
				fmt.Println(err.Error(), sql)
			}
		}
		removeAt(c.Store.InsertList, i)
	}
	time.AfterFunc(time.Second*1, func() { c.doTask() })
}

const (
	INSERT = "insert"
	UPDATE = "update"
	DELETE = "delete"
)

func (c *DbCache) get_sql(key string, val string, action string) string {
	var sql_temp string
	table := c.table(key)
	switch action {
	case INSERT:
		sql_temp = fmt.Sprintf("INSERT INTO `%s` (`%s`)VALUES('%s');", table, c.opt.KEY, val)
	case DELETE:
		sql_temp = fmt.Sprintf("DELETE from `%s` where `%s`='%s';", table, c.opt.KEY, val)
	case UPDATE:
		sql_temp = fmt.Sprintf("UPDATE `%s` set `%s`='%s' where `%s`='%s';", table, c.opt.COLUMN, key, c.opt.KEY, val)
	default:
		sql_temp = ""
	}
	return sql_temp
}
func (c *DbCache) exec(sql_temp string, value ...string) (bool, error) {
	c.AutoOpen()
	// 插入
	sql := fmt.Sprintf(sql_temp, value)
	_, err := c.db.Exec(sql)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (c *DbCache) Get(key string) (arr map[string]string, err error) {
	c.AutoOpen()
	arr = make(map[string]string)
	// 查询
	sql := fmt.Sprintf("SELECT key, value FROM %s", key)
	rows, err := c.db.Query(sql)
	if err != nil {
		return arr, err
	}

	for rows.Next() {
		// 假设有两列，都是字符串类型
		var key, val string
		if err := rows.Scan(&key, &val); err != nil {
			return arr, err
		}
		arr[key] = val
	}

	if err := rows.Err(); err != nil {
		return arr, err
	}
	return arr, nil
}

func (c *DbCache) Close() {
	if c.db != nil {
		if !c.isClosed {
			return
		}
		c.db.Close()
	}
}
