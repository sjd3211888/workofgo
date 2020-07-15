package sccsql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

//var Db *sqlx.DB

func init() {

}

func Execsqlcmd(DB *sqlx.DB, sqlcmd string, async bool) (ret int) {
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
		go execsccsql(DB, sqlcmd)
		return 0
	} else {
		return execsccsql(DB, sqlcmd)
	}
	return 0
}

func ConnectMysql(username string, password string, ipAddress string, port int, dbName string, charset string) *sqlx.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", username, password, ipAddress, port, dbName, charset)
	Db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("mysql connect failed, detail is [%v]", err.Error())
	}
	Db.SetMaxIdleConns(8)
	Db.SetMaxOpenConns(8)
	return Db
}
func Ping(Db *sqlx.DB) {
	err := Db.Ping()
	if err != nil {
		fmt.Println("ping failed")
	} else {
		fmt.Println("ping success")
	}
}

func SelectData(Db *sqlx.DB, sqlcmd string) (sccresult []map[string]string) {

	//定义结构体切片，用来存放多条查询记录

	rows, err := Db.Queryx(sqlcmd)
	if err != nil {
		fmt.Printf("query faied, error:[%v]", err.Error())
		return
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
