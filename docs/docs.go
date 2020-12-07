// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/moduserdetail": {
            "post": {
                "description": "moduserdetail 修改用户详细信息",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "description": "修改用户信息",
                        "name": "article",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/coreservice.Moduserdetail"
                        }
                    }
                ]
            }
        },
        "/queryRecnetSession": {
            "post": {
                "description": "queryRecnetSession 查询最近会话",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "description": "查询最近会话",
                        "name": "article",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/coreservice.QueryRecntSession"
                        }
                    }
                ]
            }
        },
        "/querydepartment": {
            "post": {
                "description": "querydepartment 查询部门信息",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "description": "查询部门信息",
                        "name": "article",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/coreservice.Querysccdeparment"
                        }
                    }
                ]
            }
        },
        "/querydepartmentuser": {
            "post": {
                "description": "querydepartmentuser 查询部门成员信息",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "description": "查询部门成员信息",
                        "name": "article",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/coreservice.Querysccdeparmentuser"
                        }
                    }
                ]
            }
        },
        "/querydingbyfromsccid": {
            "post": {
                "description": "querydingbyfromsccid",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "description": "根据被叫SCCid查询必达的情况",
                        "name": "article",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/coreservice.Fromdinginfo"
                        }
                    }
                ]
            }
        },
        "/querydingbymsgid": {
            "post": {
                "description": "querydingbymsgid",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "description": "如果个人必达 messagetype是0  groupid是0  群组必达 messagetype是1 groupid是群组id",
                        "name": "article",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/coreservice.Dingbyid"
                        }
                    }
                ]
            }
        },
        "/querydingbysccidandgroupid": {
            "post": {
                "description": "querydingbysccidandgroupid",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "description": "根据群组id和msessageid查询群组必达的必达情况",
                        "name": "article",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/coreservice.Dingfrommsgidindgroupid"
                        }
                    }
                ]
            }
        },
        "/querydingbytosccid": {
            "post": {
                "description": "querydingbytosccid",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "description": "根据被叫SCCid查询必达的情况",
                        "name": "article",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/coreservice.Todinginfo"
                        }
                    }
                ]
            }
        },
        "/querygps": {
            "post": {
                "description": "querygps 查询个人的历史轨迹",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "description": "sccid和时间查询历史轨迹",
                        "name": "article",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/coreservice.Querygps"
                        }
                    }
                ]
            }
        },
        "/querygroup": {
            "post": {
                "description": "querygroup 查询群组信息",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "description": "查询群组信息",
                        "name": "article",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/coreservice.Querygroupinfo"
                        }
                    }
                ]
            }
        },
        "/querygroupdingbysccid": {
            "post": {
                "description": "querygroupdingbysccid",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "description": "0 是和我相关的  1 是我发送的 2 是我接收的",
                        "name": "article",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/coreservice.Relationding"
                        }
                    }
                ]
            }
        },
        "/querygrouphistoryim": {
            "post": {
                "description": "querygrouphistoryim 根据群组查询历史消息",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "description": "根据群组查询历史消息",
                        "name": "article",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/coreservice.Querygroupimhistory"
                        }
                    }
                ]
            }
        },
        "/querygroupuser": {
            "post": {
                "description": "querygroupuser 查询群组成员",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "description": "查询群组成员",
                        "name": "article",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/coreservice.Quserygroupuserinfo"
                        }
                    }
                ]
            }
        },
        "/querymsgbymsgid": {
            "post": {
                "description": "querymsgbymsgid",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "description": "根据被叫SCCid查询必达的情况",
                        "name": "article",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/coreservice.Todinginfo"
                        }
                    }
                ]
            }
        },
        "/querynearbyscc": {
            "post": {
                "description": "querynearbyscc 查询附近的人",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "description": "查询附近的人",
                        "name": "article",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/coreservice.Querynearby"
                        }
                    }
                ]
            }
        },
        "/queryofflinemsg": {
            "post": {
                "description": "queryofflinemsg 查询离线消息",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "description": "查询离线消息",
                        "name": "article",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/coreservice.Querypersonofflinemsg"
                        }
                    }
                ]
            }
        },
        "/querypersondingbysccid": {
            "post": {
                "description": "querypersondingbysccid",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "description": "0 是和我相关的  1 是我发送的 2 是我接收的",
                        "name": "article",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/coreservice.Relationding"
                        }
                    }
                ]
            }
        },
        "/querypersonhistoryim": {
            "post": {
                "description": "querypersonhistoryim 查询个人消息",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "description": "查询历史信息",
                        "name": "article",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/coreservice.Querrypersonimhistory"
                        }
                    }
                ]
            }
        },
        "/querysccuserdetail": {
            "post": {
                "description": "querysccuserdetail 查询用户详细信息",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "description": "查询用户详细信息",
                        "name": "article",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/coreservice.Queryuserdetail"
                        }
                    }
                ]
            }
        },
        "/queryuser": {
            "post": {
                "description": "queryuser 查询个人成员详细信息",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "description": "查询个人成员详细信息",
                        "name": "article",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/coreservice.Queryuserinfo"
                        }
                    }
                ]
            }
        },
        "/reportgps": {
            "post": {
                "description": "reportgps 上报轨迹",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "description": "上报轨迹",
                        "name": "article",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/coreservice.Reportgps"
                        }
                    }
                ]
            }
        }
    },
    "definitions": {
        "coreservice.Dingbyid": {
            "type": "object",
            "required": [
                "groupid",
                "messageid",
                "messagetype"
            ],
            "properties": {
                "groupid": {
                    "type": "string"
                },
                "messageid": {
                    "type": "string"
                },
                "messagetype": {
                    "description": "binding:\"required\"修饰的字段，若接收为空值，则报错，是必须字段",
                    "type": "string"
                }
            }
        },
        "coreservice.Dingfrommsgidindgroupid": {
            "type": "object",
            "required": [
                "dingstatus",
                "pagenum"
            ],
            "properties": {
                "dingstatus": {
                    "type": "string"
                },
                "groupid": {
                    "type": "integer"
                },
                "messageid": {
                    "type": "integer"
                },
                "pagenum": {
                    "description": "binding:\"required\"修饰的字段，若接收为空值，则报错，是必须字段",
                    "type": "integer"
                }
            }
        },
        "coreservice.Fromdinginfo": {
            "type": "object",
            "required": [
                "dingstatus",
                "fromsccid",
                "pagenum"
            ],
            "properties": {
                "dingstatus": {
                    "type": "string"
                },
                "fromsccid": {
                    "description": "binding:\"required\"修饰的字段，若接收为空值，则报错，是必须字段",
                    "type": "string"
                },
                "pagenum": {
                    "type": "integer"
                }
            }
        },
        "coreservice.Moduserdetail": {
            "type": "object",
            "required": [
                "addr",
                "mailbox",
                "mobilephone",
                "phone",
                "post",
                "sccid"
            ],
            "properties": {
                "addr": {
                    "type": "string"
                },
                "mailbox": {
                    "type": "string"
                },
                "mobilephone": {
                    "type": "string"
                },
                "phone": {
                    "type": "string"
                },
                "post": {
                    "type": "string"
                },
                "sccid": {
                    "description": "binding:\"required\"修饰的字段，若接收为空值，则报错，是必须字段",
                    "type": "string"
                }
            }
        },
        "coreservice.Querrypersonimhistory": {
            "type": "object",
            "required": [
                "pagenum",
                "peerid",
                "sccid"
            ],
            "properties": {
                "pagenum": {
                    "type": "integer"
                },
                "peerid": {
                    "type": "string"
                },
                "sccid": {
                    "description": "binding:\"required\"修饰的字段，若接收为空值，则报错，是必须字段",
                    "type": "string"
                }
            }
        },
        "coreservice.QueryRecntSession": {
            "type": "object",
            "required": [
                "sccid"
            ],
            "properties": {
                "sccid": {
                    "description": "binding:\"required\"修饰的字段，若接收为空值，则报错，是必须字段",
                    "type": "string"
                }
            }
        },
        "coreservice.Querygps": {
            "type": "object",
            "required": [
                "endtime",
                "pagenum",
                "sccid",
                "starttime"
            ],
            "properties": {
                "endtime": {
                    "type": "integer"
                },
                "needdescription": {
                    "type": "string"
                },
                "pagenum": {
                    "type": "integer"
                },
                "sccid": {
                    "description": "binding:\"required\"修饰的字段，若接收为空值，则报错，是必须字段",
                    "type": "string"
                },
                "starttime": {
                    "type": "integer"
                }
            }
        },
        "coreservice.Querygroupimhistory": {
            "type": "object",
            "required": [
                "groupid",
                "pagenum"
            ],
            "properties": {
                "groupid": {
                    "description": "binding:\"required\"修饰的字段，若接收为空值，则报错，是必须字段",
                    "type": "integer"
                },
                "pagenum": {
                    "type": "integer"
                }
            }
        },
        "coreservice.Querygroupinfo": {
            "type": "object",
            "required": [
                "sccid"
            ],
            "properties": {
                "sccid": {
                    "description": "binding:\"required\"修饰的字段，若接收为空值，则报错，是必须字段",
                    "type": "string"
                }
            }
        },
        "coreservice.Querynearby": {
            "type": "object",
            "required": [
                "distance",
                "latitude",
                "longitude"
            ],
            "properties": {
                "distance": {
                    "description": "binding:\"required\"修饰的字段，若接收为空值，则报错，是必须字段",
                    "type": "integer"
                },
                "latitude": {
                    "type": "string"
                },
                "longitude": {
                    "type": "string"
                }
            }
        },
        "coreservice.Querypersonofflinemsg": {
            "type": "object",
            "required": [
                "sccid"
            ],
            "properties": {
                "sccid": {
                    "description": "binding:\"required\"修饰的字段，若接收为空值，则报错，是必须字段",
                    "type": "string"
                }
            }
        },
        "coreservice.Querysccdeparment": {
            "type": "object",
            "required": [
                "departmentid"
            ],
            "properties": {
                "departmentid": {
                    "description": "binding:\"required\"修饰的字段，若接收为空值，则报错，是必须字段",
                    "type": "string"
                }
            }
        },
        "coreservice.Querysccdeparmentuser": {
            "type": "object",
            "required": [
                "departmentid",
                "onlydispatcher"
            ],
            "properties": {
                "departmentid": {
                    "description": "binding:\"required\"修饰的字段，若接收为空值，则报错，是必须字段",
                    "type": "string"
                },
                "onlydispatcher": {
                    "type": "string"
                }
            }
        },
        "coreservice.Queryuserdetail": {
            "type": "object",
            "required": [
                "sccid"
            ],
            "properties": {
                "sccid": {
                    "description": "binding:\"required\"修饰的字段，若接收为空值，则报错，是必须字段",
                    "type": "string"
                }
            }
        },
        "coreservice.Queryuserinfo": {
            "type": "object",
            "required": [
                "sccid"
            ],
            "properties": {
                "sccid": {
                    "description": "binding:\"required\"修饰的字段，若接收为空值，则报错，是必须字段",
                    "type": "string"
                }
            }
        },
        "coreservice.Quserygroupuserinfo": {
            "type": "object",
            "required": [
                "groupid"
            ],
            "properties": {
                "groupid": {
                    "description": "binding:\"required\"修饰的字段，若接收为空值，则报错，是必须字段",
                    "type": "string"
                }
            }
        },
        "coreservice.Relationding": {
            "type": "object",
            "required": [
                "dingstatus",
                "pagenum",
                "sccid"
            ],
            "properties": {
                "dingstatus": {
                    "type": "string"
                },
                "pagenum": {
                    "type": "integer"
                },
                "sccid": {
                    "description": "binding:\"required\"修饰的字段，若接收为空值，则报错，是必须字段",
                    "type": "string"
                },
                "sccidstatus": {
                    "description": "0 是和我相关的  1 是我发送的 2 是我接收的",
                    "type": "integer"
                }
            }
        },
        "coreservice.Reportgps": {
            "type": "object",
            "required": [
                "gps",
                "latitude",
                "longitude",
                "sccid"
            ],
            "properties": {
                "angle": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "gps": {
                    "type": "string"
                },
                "latitude": {
                    "type": "string"
                },
                "longitude": {
                    "type": "string"
                },
                "sccid": {
                    "description": "binding:\"required\"修饰的字段，若接收为空值，则报错，是必须字段",
                    "type": "string"
                },
                "speed": {
                    "type": "integer"
                }
            }
        },
        "coreservice.Todinginfo": {
            "type": "object",
            "required": [
                "dingstatus",
                "pagenum",
                "tosccid"
            ],
            "properties": {
                "dingstatus": {
                    "type": "string"
                },
                "pagenum": {
                    "type": "integer"
                },
                "tosccid": {
                    "description": "binding:\"required\"修饰的字段，若接收为空值，则报错，是必须字段",
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "",
	Description: "",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
