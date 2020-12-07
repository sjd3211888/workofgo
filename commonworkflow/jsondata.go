package commconworkflow

type ApproverWorkflowuser struct {
	ApproverUsers string `json:"approverWorkflowuser" binding:"required"`
}
type CcWorkflowluser struct {
	CcUsers string `json:"cc" binding:"required"`
}
type CreateCoworkflow struct {
	// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
	Subject      string                 `json:"subject" binding:"required"`
	Workflowtype string                 `json:"workflowtype"`
	Approve      []ApproverWorkflowuser `json:"approvers" binding:"required"`
	Cc           []CcWorkflowluser      `json:"Ccusers"`
	Creater      string                 `json:"creater" binding:"required"`
	Textinfo     string                 `json:"textinfo"`
	Filepath     string                 `json:"filepath"`
}
type Coworkflowtodo struct {
	Sccid   string `json:"sccid" binding:"required"`
	Pagenum int    `json:"pagenum"`
}
type Coworkflowcc struct {
	Sccid   string `json:"sccid" binding:"required"`
	Pagenum int    `json:"pagenum"`
}
type Approve struct {
	Sccid    string `json:"sccid" binding:"required"`
	Approve  string `json:"approve" binding:"required"`
	Workid   string `json:"workid" binding:"required"`
	Textinfo string `json:"textinfo"`
	Fileptah string `json:"filepath"`
}
type Coworkflowdone struct {
	Sccid   string `json:"sccid" binding:"required"`
	Pagenum int    `json:"pagenum"`
}
type Coworkflowbyworkid struct {
	Workid string `json:"workid" binding:"required"`
}
type Comment struct {
	Commentworkid string `json:"commentworkid" binding:"required"`
	Commentid     string `json:"commentid" binding:"required"`
	Commenter     string `json:"commenter" binding:"required"`
	Textinfo      string `json:"textinfo"`
	Fileptah      string `json:"filepath"`
}
