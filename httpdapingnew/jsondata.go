package gzcanteen

type GZworkplace struct {
	// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
	Workplace int `json:"workplace" binding:"required"`
}
