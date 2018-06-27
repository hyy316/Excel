package basic

type Address struct{
	Id int64			
	Ip string			`json:"ip"`
	Start  string		
	End  string			
	Username  string	`json:"username"`
	Attribution  string	`json:"attribution"`
	Manager  string		`json:"manager"`
	Ascription  string	`json:"ascription"`
	Use  string			`json:"use"`
	Accuse  string		`json:"accuse"`
	Principal  string	`json:"principal"`
	Phone  string		`json:"phone"`
}

type Iptype struct{
	Ip string	`json:"Ip"`
}
