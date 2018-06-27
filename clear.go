package main 

import (
    _ "github.com/mattn/go-sqlite3"
    "database/sql"
    "fmt"
    "time"
)

func main() {
	db, err1 := sql.Open("sqlite3", "db/data.db")
  	checkErr(err1)
    stmt, _ := db.Prepare("delete from address")
    _,errExec:=stmt.Exec()
    stmt.Close()
    db.Close()
    if errExec==nil{
		fmt.Println("数据库清除成功")
    }else{
    	fmt.Println("数据库清除失败")	
    }
    time.Sleep(time.Duration(2)*time.Second)
}


func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}
