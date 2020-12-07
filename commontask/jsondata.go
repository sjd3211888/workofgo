package commtask

type Cctaskuser struct {
	CcUsers string `json:"cc" binding:"required"`
}
type Advancedinfi struct {
	Taskbegintime string `json:"taskbegintime"`
	Taskendtime   string `json:"taskendtime"`
	Tasktype      string `json:"tasktype"`
	Taskegrade    string `json:"taskegrade"`
	Taskapprover  string `json:"taskapprover"`
}
type CreateCommtask struct {
	// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
	Taskname string       `json:"taskname" binding:"required"`
	Executor string       `json:"exectutor" binding:"required"`
	Creater  string       `json:"creater" binding:"required"`
	Cc       []Cctaskuser `json:"Cctaskuser"`
	Textinfo string       `json:"textinfo" binding:"required"`
	Filepath string       `json:"filepath"`
	Advanced Advancedinfi `json:"Advancedinfi"`
}
type Docommtask struct {
	Taskid   string `json:"taskid" binding:"required"`
	Executor string `json:"exectutor" binding:"required"`
	Textinfo string `json:"textinfo" binding:"required"`
	Filepath string `json:"filepath"`
}
type Approvetask struct {
	Taskid       string `json:"taskid" binding:"required"`
	Approveornot string `json:"approveornot"`
}
type Tasker struct {
	Taskuser string `json:"taskuser" binding:"required"`
	Pagenum  int    `json:"pagenum"`
}

type Querytask struct {
	Taskid string `json:"taskid" binding:"required"`
}

type Commenttask struct {
	Taskid      string `json:"taskid" binding:"required"`
	Commenter   string `json:"commenter" binding:"required"`
	Commentinfo string `json:"commentinfo" binding:"required"`
	Commentpath string `json:"commentpath"`
}
