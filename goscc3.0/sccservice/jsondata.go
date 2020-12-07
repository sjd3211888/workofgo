package coreservice

type Fromdinginfo struct {
	// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
	Sccid      string `json:"fromsccid" binding:"required"`
	Pagenum    int    `json:"pagenum" binding:"required"`
	Dingstatus string `json:"dingstatus" binding:"required"`
}

type Todinginfo struct {
	// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
	Sccid      string `json:"tosccid" binding:"required"`
	Pagenum    int    `json:"pagenum" binding:"required"`
	Dingstatus string `json:"dingstatus" binding:"required"`
}

type Dingfrommsgidindgroupid struct {
	// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
	Pagenum    int    `json:"pagenum" binding:"required"`
	Messageid  int    `json:"messageid"`
	Groupid    int    `json:"groupid"`
	Dingstatus string `json:"dingstatus"  binding:"required"`
}
type Dingbyid struct {
	// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
	Messagetype string `json:"messagetype" binding:"required"`
	Messageid   string `json:"messageid" binding:"required"`
	Groupid     string `json:"groupid" binding:"required"`
}

type Relationding struct {
	// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
	Sccid       string `json:"sccid" binding:"required"`
	Pagenum     int    `json:"pagenum" binding:"required"`
	Dingstatus  string `json:"dingstatus" binding:"required"`
	Sccidstatus int    `json:"sccidstatus"` //0 是和我相关的  1 是我发送的 2 是我接收的
}
type Queryuserdetail struct {
	// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
	Sccid string `json:"sccid" binding:"required"`
}
type Querynearby struct {
	// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
	Distance  int    `json:"distance" binding:"required"`
	Longitude string `json:"longitude" binding:"required"`
	Latitude  string `json:"latitude" binding:"required"`
}
type Moduserdetail struct {
	// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
	Sccid       string `json:"sccid" binding:"required"`
	Post        string `json:"post" binding:"required"`
	Mailbox     string `json:"mailbox" binding:"required"`
	Addr        string `json:"addr" binding:"required"`
	Phone       string `json:"phone" binding:"required"`
	Mobilephone string `json:"mobilephone" binding:"required"`
}
type Querygroupimhistory struct {
	// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
	Groupid int `json:"groupid" binding:"required"`
	Pagenum int `json:"pagenum" binding:"required"`
}
type Querrypersonimhistory struct {
	// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
	Sccid   string `json:"sccid" binding:"required"`
	Peerid  string `json:"peerid" binding:"required"`
	Pagenum int    `json:"pagenum" binding:"required"`
}
type Querygps struct {
	// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
	Sccid           string `json:"sccid" binding:"required"`
	Starttime       int    `json:"starttime" binding:"required"`
	Endtime         int    `json:"endtime" binding:"required"`
	Pagenum         int    `json:"pagenum" binding:"required"`
	Needdescription string `json:"needdescription"`
}
type Reportgps struct {
	// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
	Sccid       string `json:"sccid" binding:"required"`
	Longitude   string `json:"longitude" binding:"required"`
	Latitude    string `json:"latitude" binding:"required"`
	Gps         string `json:"gps" binding:"required"`
	Speed       int    `json:"speed"`
	Angle       string `json:"angle"`
	Description string `json:"description"`
}
type QueryRecntSession struct {
	// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
	Sccid string `json:"sccid" binding:"required"`
}
type Querypersonofflinemsg struct {
	// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
	Sccid string `json:"sccid" binding:"required"`
}
type Quserygroupuserinfo struct {
	// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
	Groupid string `json:"groupid" binding:"required"`
}
type Queryuserinfo struct {
	// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
	Sccid string `json:"sccid" binding:"required"`
}
type Querygroupinfo struct {
	// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
	Sccid string `json:"sccid" binding:"required"`
}
type Querysccdeparmentuser struct {
	// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
	Departmentid   string `json:"departmentid" binding:"required"`
	Onlydispatcher string `json:"onlydispatcher" binding:"required"`
}
type Querysccdeparment struct {
	// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
	Departmentid string `json:"departmentid" binding:"required"`
}

type Querymsgbyid struct {
	// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
	Msgid   int    `json:"msgid" binding:"required"`
	Msgtype string `json:"msgtype" binding:"required"`
}

type MobilephoneInfo struct {
	// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
	Sccid           string `json:"sccid" binding:"required"`
	Token           string `json:"token" binding:"required"`
	Mobilephonetype string `json:"Mobilephonetype" binding:"required"`
}
