package controllers

import "project_bluebell/models"

// 后面结构体名称后面加上下划线表示这个是 接口文档 中使用的结构体
// 帖子信息接口返回参数
type ResPostList_ struct {
	Code    ResCode                 `json:"code"`    // 状态码
	Message string                  `json:"message"` // 提示信息
	Data    []*models.ApiPostDetail `json:"data"`    // 数据
}

type ResToken_ struct {
	Code    ResCode `json:"code"`    // 状态码
	Message string  `json:"message"` // 提示信息
	Data    string  `json:"data"`    // 数据
}

// 没有数据的响应返回参数
type ResNil_ struct {
	Code    ResCode `json:"code"`    // 状态码
	Message string  `json:"message"` // 提示信息
	Data    string  `json:"data"`    // 数据
}

// 查询帖子接口传入参数
type ReqPageSize_ struct {
	Page string `json:"page"` // 页码
	Size string `json:"size"` // 每页数据量
}

// 查询所有社区返回参数
type ResCommunity_ struct {
	Code    ResCode             `json:"code"`    // 状态码
	Message string              `json:"message"` // 提示信息
	Data    []*models.Community `json:"data"`    // 数据
}

// 查询指定社区的返回参数
type ResCommunityDetail_ struct {
	Code    ResCode                   `json:"code"`    // 状态码
	Message string                    `json:"message"` // 提示信息
	Data    []*models.CommunityDetail `json:"data"`    // 数据
}
