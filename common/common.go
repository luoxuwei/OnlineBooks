package common


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