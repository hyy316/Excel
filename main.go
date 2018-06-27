package main

import (
    "fmt"
    "html/template"
    "log"
    "net/http"
    "net/url"
    "strings"
    "encoding/json"  
    "Excel/basic"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    "github.com/Luxurioust/excelize"
    "strconv"
    "io/ioutil"
    "time"

)


func sayhelloName(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()       //解析url传递的参数，对于POST则解析响应包的主体（request body）
    //注意:如果没有调用ParseForm方法，下面无法获取表单的数据
    fmt.Println(r.Form) //这些信息是输出到服务器端的打印信息
    fmt.Println("path", r.URL.Path)
    fmt.Println("scheme", r.URL.Scheme)
    fmt.Println(r.Form["url_long"])
    for k, v := range r.Form {
        fmt.Println("key:", k)
        fmt.Println("val:", strings.Join(v, ""))
    }
    fmt.Fprintf(w, "Hello astaxie!") //这个写入到w的是输出到客户端的
}



func login(w http.ResponseWriter, r *http.Request) {
/*	db, err := sql.Open("sqlite3", ".db/data.db")
    checkErr(err)

    //插入数据
    stmt, err := db.Prepare("INSERT INTO address(Start, End, Username) values(?,?,?)")
    checkErr(err)
    res, err := stmt.Exec(1234, 2341, "qwerqrqw")
    checkErr(err)
    id, err := res.LastInsertId()
    checkErr(err)

    fmt.Println(id)
	*/
	//r.ParseForm()
        //请求的是登录数据，那么执行登录的逻辑判断
	/* var usr Usr 
	   usr.Name=r.Form["username"]
	   usr.Passwd=r.Form["password"]
	   var result basic.Result*/
	//var result basic.Result
   // fmt.Println("method:", r.Method) //获取请求的方法
    if r.Method == "GET" {
        t, _ := template.ParseFiles("index.html")
        log.Println(t.Execute(w, nil))
    } else {
    //result.Code = 101  
	//result.Data = usr
    //result.Message = "用户名或密码不正确"
 	//bytes, _ := json.Marshal(result)  
    //fmt.Fprint(w, string(bytes)) 
    //w.Write(bytes) 
    //}
	}
}


func upload(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json
	
	t1 := time.Now()

	db, err1 := sql.Open("sqlite3", "db/data.db")
    checkErr(err1)
   	defer db.Close()
    tx, err := db.Begin()
    if err != nil {
            fmt.Println("Beginx error:", err)
            panic(err)
    }
    sql := "INSERT INTO address(Start, End, Username,Attribution,Manager,Ascription,Use,Accuse,Principal,Phone) values(?,?,?,?,?,?,?,?,?,?)"
	    file, _, err := r.FormFile("uploadfile")
	    if err != nil {
	        fmt.Println(err)
	        return
	    }
	    xlsx, err := excelize.OpenReader(file)
	    defer file.Close()
	    var result basic.Result
	    data1 := xlsx.GetRows("数据源1")
	    data2 := xlsx.GetRows("数据源2")
	    for i:=1;i<len(data2);i++{
	    	var address basic.Address
	    	address.Start=data2[i][3]
	    	address.End =data2[i][4]  
			address.Username =data1[i][5] 	
			address.Attribution =data1[i][6]	
			address.Manager =data1[i][7]	
			address.Ascription =data2[i][8]	
			address.Use =data2[i][9]	
			address.Accuse=data2[i][10] 
			address.Principal = data2[i][19]
			address.Phone =data2[i][20]
		   	//stmt.Exec(address.Start,address.End,address.Username,address.Attribution,address.Manager,address.Ascription,address.Use,address.Accuse, address.Principal,address.Phone)
		   	tx.Exec(sql, address.Start,address.End,address.Username,address.Attribution,address.Manager,address.Ascription,address.Use,address.Accuse, address.Principal,address.Phone)
		   	//fmt.Println(i)
		    }
		    errtx:=tx.Commit()
		    elapsed := time.Since(t1)
    		fmt.Println("App elapsed: ", elapsed)
		 	if errtx!=nil{
		 		fmt.Println(errtx)
		 		result.Code = 201  
		    	result.Message = "导入失败"
		 	}else{
		 		result.Code = 200  
		    	result.Message = "成功导入   "+strconv.Itoa(len(data2))+"   条数据"
		 	}
		 	bytes, _ := json.Marshal(result)  
			w.Write(bytes)
}


func query(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json
	
    body, _ := ioutil.ReadAll(r.Body)
    var s basic.Iptype
    json.Unmarshal(body, &s)
	ip :=s.Ip
	ips:=strings.Replace(ip, ".", "", -1)
	db, err1 := sql.Open("sqlite3", "db/data.db")
    checkErr(err1)
    defer db.Close()
    sql := "SELECT * FROM address WHERE start<=?AND end>=?"
    tx, _ := db.Begin()
    rows,_:=tx.Query(sql,ips,ips)
    var addresses= []basic.Address{}
    for rows.Next() {
    	var address basic.Address
    	address.Ip=ip	
	    if err := rows.Scan(&address.Id,&address.Start,&address.End,&address.Username,&address.Attribution,&address.Manager,&address.Ascription,&address.Use,&address.Accuse,&address.Principal,&address.Phone); err != nil {
	        log.Fatal(err)
	    }
	    addresses=append(addresses,address)
	}
	export(addresses)
	var result basic.Result
	result.Code = 200  
	result.Data = addresses
    result.Message = "查询成功"
 	bytes, _ := json.Marshal(result)  
 	fmt.Println("执行单IP查询...")
	w.Write(bytes) 
}



func query2(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json
	db, err1 := sql.Open("sqlite3", "db/data.db")
  	checkErr(err1)
    defer db.Close()
    sql := "SELECT * FROM address WHERE start<=?AND end>=?"
    tx, _ := db.Begin()
	var addresses= []basic.Address{}
		file, _, err := r.FormFile("uploadfile")
	    if err != nil {
	        fmt.Println(err)
	        return
	    }
	    xlsx, err := excelize.OpenReader(file)
	    defer file.Close()
	    data := xlsx.GetRows("Sheet1")

	for i:=1;i<len(data);i++{
		ip:=strings.Replace(data[i][0], ".", "", -1)
	    rows,_:=tx.Query(sql,ip,ip)
		for rows.Next() {
		    var address basic.Address
		    address.Ip=data[i][0]	
			if err := rows.Scan(&address.Id,&address.Start,&address.End,&address.Username,&address.Attribution,&address.Manager,&address.Ascription,&address.Use,&address.Accuse,&address.Principal,&address.Phone); err != nil {
			    log.Fatal(err)
			}
			addresses=append(addresses,address)
		}
	}
	export(addresses)
	var result basic.Result
	result.Code = 200  
	result.Data = addresses
    result.Message = "查询成功"
 	bytes, _ := json.Marshal(result)  
 	fmt.Println("执行批量查询...")
	w.Write(bytes) 
}

func export(datas []basic.Address) error{
    xlsx := excelize.NewFile()

    index := xlsx.NewSheet("Sheet1")
    xlsx.SetCellValue("Sheet1", "A1", "IP地址")
    xlsx.SetCellValue("Sheet1", "B1", "接入用户名")
    xlsx.SetCellValue("Sheet1", "C1", "归属分公司")
    xlsx.SetCellValue("Sheet1", "D1", "客户经理")
    xlsx.SetCellValue("Sheet1", "E1", "综资一次分配中归属")
    xlsx.SetCellValue("Sheet1", "F1", "使用部门")
    xlsx.SetCellValue("Sheet1", "G1", "负责部门")
    xlsx.SetCellValue("Sheet1", "H1", "负责人")
    xlsx.SetCellValue("Sheet1", "I1", "电话")
    for i:=0;i<len(datas);i++{
        xlsx.SetCellValue("Sheet1", "A"+strconv.Itoa(i+2), datas[i].Ip)
        xlsx.SetCellValue("Sheet1", "B"+strconv.Itoa(i+2), datas[i].Username)
        xlsx.SetCellValue("Sheet1", "C"+strconv.Itoa(i+2), datas[i].Attribution)
        xlsx.SetCellValue("Sheet1", "D"+strconv.Itoa(i+2), datas[i].Manager)
        xlsx.SetCellValue("Sheet1", "E"+strconv.Itoa(i+2), datas[i].Ascription)
        xlsx.SetCellValue("Sheet1", "F"+strconv.Itoa(i+2), datas[i].Use)
        xlsx.SetCellValue("Sheet1", "G"+strconv.Itoa(i+2), datas[i].Accuse)
        xlsx.SetCellValue("Sheet1", "H"+strconv.Itoa(i+2), datas[i].Principal)
        xlsx.SetCellValue("Sheet1", "I"+strconv.Itoa(i+2), datas[i].Phone)
    }
    // Set active sheet of the workbook.
    xlsx.SetActiveSheet(index)
    err := xlsx.SaveAs("file/1.xlsx")
    if err != nil {
        return err
    }
    return nil
}

func download(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Disposition", "attachment; filename="+url.QueryEscape("download.xlsx")) 
	http.ServeFile(w, r, "file/1.xlsx")
}


func remove(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json

	db, err1 := sql.Open("sqlite3", "db/data.db")
  	checkErr(err1)
    defer db.Close()
    stmt, _ := db.Prepare("delete from address")
    _,errExec:=stmt.Exec()
    stmt.Close()
  	var result basic.Result 
    if errExec==nil{
		result.Code = 200 
	    result.Message = "清除成功" 	
    }else{
    	checkErr(errExec)
    	result.Code = 201 
	    result.Message = "清除失败" 	
    }
    bytes, _ := json.Marshal(result)
    fmt.Println("执行remove...")
    w.Write(bytes) 
}


func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}

func main() {
	http.HandleFunc("/",login)
	http.HandleFunc("/remove",remove)			//清空数据库
    http.HandleFunc("/upload",upload)			//更新数据库
    http.HandleFunc("/query",query)				//单个IP查询
    http.HandleFunc("/query2",query2)			//多IP查询
    http.HandleFunc("/download",download)		//查询结果导出
    err := http.ListenAndServe(":9090", nil) 	//设置监听的端口
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }

} 	