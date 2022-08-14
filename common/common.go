package common

import (
	beego "github.com/beego/beego/v2/server/web"
	"strings"
)

const WorkingDirectory = "./"

// 用户权限
const (
	// 超级管理员.
	MemberSuperRole = 0
	//普通管理员.
	MemberAdminRole = 1
	//普通用户.
	MemberGeneralRole = 2
)

func Role(role int) string {
	if role == MemberSuperRole {
		return "超级管理员"
	} else if role == MemberAdminRole {
		return "管理员"
	} else if role == MemberGeneralRole {
		return "普通用户"
	} else {
		return ""
	}
}

//图书关系
const (
	// 创始人.
	BookFounder = 0
	//管理
	BookAdmin = 1
	//编辑
	BookEditor = 2
	//普通用户
	BookGeneral = 3
)

func BookRole(role int) string {
	switch role {
	case BookFounder:
		return "创始人"
	case BookAdmin:
		return "管理员"
	case BookEditor:
		return "编辑"
	case BookGeneral:
		return "普通用户"
	default:
		return ""
	}

}

func IsAllowedFileExt(ext string) bool {

	if strings.HasPrefix(ext, ".") {
		ext = string(ext[1:])
	}
	exts := getFileExt()

	for _, item := range exts {
		if strings.EqualFold(item, ext) {
			return true
		}
	}
	return false
}

func getFileExt() []string {
	ext := beego.AppConfig.DefaultString("upload_file_ext", "png|jpg|jpeg|gif|txt|doc|docx|pdf")
	temp := strings.Split(ext, "|")
	exts := make([]string, len(temp))

	i := 0
	for _, item := range temp {
		if item != "" {
			exts[i] = item
			i++
		}
	}
	return exts
}