package main

import (
	"fmt"

	"github.com/liuqi0826/common/database"
)

func main() {
	var err error

	var mysql *database.SQLProxy
	var mysqlAddress = database.EditMySQLPath("127.0.0.1", "root", "passwrod", "database", "")
	mysql, err = database.CreateMySQLConnect(mysqlAddress)
	if err != nil {
		fmt.Println(err)
	}
	err = mysql.Ping()
	if err != nil {
		fmt.Println(err)
	}

	var lite3 *database.SQLProxy
	lite3, err = database.CreateMySQLConnect("data/main")
	if err != nil {
		fmt.Println(err)
	}
	err = lite3.Ping()
	if err != nil {
		fmt.Println(err)
	}
}
