package sccsql

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

//var Db *sqlx.DB
type Mysqlconnectpool struct {
	sqlpoint     *sqlx.DB
	mysqlip      string
	username     string
	password     string
	databasename string
	mysqlport    int
}

func init() {

}
func MysqlRealEscapeString(value string) string {
	var tmpstring string
	tmpstring = strings.Replace(value, "\\0", "\\\\0", -1)
	tmpstring = strings.Replace(value, "\n", "\\n", -1)
	tmpstring = strings.Replace(value, "\r", "\\r", -1)
	tmpstring = strings.Replace(value, "\"", "\\\"", -1)
	tmptab := 0x1a
	strtab := string(tmptab)
	strnewtab := "\\"
	strnewtab = strnewtab + string(0x5a)
	tmpstring = strings.Replace(value, strtab, strnewtab, -1)
	tmpstring = strings.Replace(value, "\\'", "\\\\'", -1)
	tmpstring = strings.Replace(value, "\\", "\\\\", -1)
	return tmpstring
}

func (mysqlpool *Mysqlconnectpool) Initmysql(mysqlip string, username string, password string, databasename string, mysqlport int) {
	if nil == mysqlpool.sqlpoint {
		mysqlpool.mysqlip = mysqlip
		mysqlpool.username = username
		mysqlpool.password = password
		mysqlpool.databasename = databasename
		mysqlpool.mysqlport = mysqlport
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", mysqlpool.username, mysqlpool.password, mysqlpool.mysqlip, mysqlpool.mysqlport, mysqlpool.databasename, "utf8mb4")
		Db, err := sqlx.Connect("mysql", dsn) //线程池不太对 要改
		if err != nil {
			fmt.Printf("mysql connect failed, detail is [%v]", err.Error())
		}
		Db.SetMaxIdleConns(8)
		Db.SetMaxOpenConns(8)
		Db.SetConnMaxLifetime(time.Second * 30)
		mysqlpool.sqlpoint = Db
	} else {
		fmt.Println("mysql point is not null,")
	}
}
func (mysqlpool *Mysqlconnectpool) Execsqlcmd(sqlcmd string, async bool) (ret int) {
	//如果是异步
	execsccsql := func(DB *sqlx.DB, sqlcmd string) int {
		r, err := DB.Exec(sqlcmd)
		if err != nil {
			fmt.Println("exec failed, ", err)
			count, _ := r.RowsAffected()
			if 1 == count {
				return 0
			}
		}
		return -1
	}
	if async {
		go execsccsql(mysqlpool.sqlpoint, sqlcmd)
		return 0
	}
	return execsccsql(mysqlpool.sqlpoint, sqlcmd)
}

func (mysqlpool *Mysqlconnectpool) connectMysql() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", mysqlpool.username, mysqlpool.password, mysqlpool.mysqlip, mysqlpool.mysqlport, mysqlpool.databasename, "utf8mb4")
	Db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("mysql connect failed, detail is [%v]", err.Error())
	}
	Db.SetMaxIdleConns(8)
	Db.SetMaxOpenConns(8)
	mysqlpool.sqlpoint = Db
}
func (mysqlpool *Mysqlconnectpool) Ping() {
	err := mysqlpool.sqlpoint.Ping()
	if err != nil {
		fmt.Println("ping failed")
	} else {
		fmt.Println("ping success")
	}
}

func (mysqlpool *Mysqlconnectpool) SelectData(sqlcmd string) []map[string]string {

	//定义结构体切片，用来存放多条查询记录

	rows, err := mysqlpool.sqlpoint.Queryx(sqlcmd)
	if err != nil {
		fmt.Printf("query faied, error:[%v]", err.Error())
		return nil
	}
	columns, _ := rows.Columns()

	//定义一个切片,长度是字段的个数,切片里面的元素类型是sql.RawBytes
	values := make([]sql.RawBytes, len(columns))
	//定义一个切片,元素类型是interface{} 接口
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		//把sql.RawBytes类型的地址存进去了
		scanArgs[i] = &values[i]
	}
	//获取字段值
	var result []map[string]string
	for rows.Next() {
		res := make(map[string]string)
		rows.Scan(scanArgs...)
		for i, col := range values {

			res[columns[i]] = string(col)
		}

		result = append(result, res)
	}
	rows.Close()
	return result
}
func (mysqlpool *Mysqlconnectpool) Insertmql(sqlcmd string) int64 {
	r, err := mysqlpool.sqlpoint.Exec(sqlcmd)
	if err != nil {
		fmt.Println("exec failed, ", err)
		return 0
	}
	id, err := r.LastInsertId()
	if err != nil {
		fmt.Println("exec failed, ", err)
		return 0
	}
	return id
}
