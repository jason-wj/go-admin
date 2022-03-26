package authdto

const (
	LoginUserId = "loginUserId"
	UserInfo    = "userinfo"
	RoleId      = "roleId"
	RoleName    = "role"
	RoleIds     = "roleIds"
	RoleKey     = "roleKey"
	UserName    = "userName"
	DataScope   = "dataScope"
)

type Resp struct {
	RequestId string      `json:"requestId"`
	Msg       string      `json:"msg"`
	Code      int         `json:"code"`
	Data      interface{} `json:"data"`
}

type Data struct {
	Token    string      `json:"token"`
	Expire   string      `json:"expire"`
	UserInfo interface{} `json:"userInfo"`
}
