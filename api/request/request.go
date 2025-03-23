package request

type UserRegister struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	Email      string `json:"email"`
	Code       string `json:"code"`
	Device_Num string `json:"device_num"`
}

type UserUpdate struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Visitor struct {
	Device_Num string `json:"device_num"`
}

type ForAlter struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Code     string `json:"code"`
}

type UserCancel struct {
	Email string `json:"email"`
}

type Slogan struct {
	Slogan string `json:"slogan"`
}

type Email struct {
	Email string `json:"email"`
}

type Question struct {
	Topic       string `json:"topic"`
	Description string `json:"description"`
	Cycle       string `json:"cycle"`
}

type Image struct {
	Id  string `json:"id"` //是userId或者groupId
	Url string `json:"url"`
}

type TaskRequest struct {
	Title   string `json:"title"`
	Details string `json:"details"`
}

type PostGoalRequest struct {
	Date  string        `json:"date"`  // 任务所属日期
	Tasks []TaskRequest `json:"tasks"` // 该日期下的任务列表
}

type PostGoalRequests struct {
	Goals []PostGoalRequest `json:"goals"`
}

type GroupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	GroupNum    string `json:"group_num"`  //群号 = id
	Password    string `json:"password"`   //群主设置
	IsPublic    bool   `json:"is_public"`  //任务是否公开
	ExecuteId   string `json:"execute_id"` //执行任务的用户id
}

type ExecuteGroupMember struct {
	UserId   string `json:"user_id"`
	GroupNum string `json:"group_num"`
}

type Message struct {
	From    string `json:"from"`
	To      string `json:"to"` // 可以是用户ID或群ID
	Content string `json:"content"`
	Type    string `json:"type"`
}

type UserHistory struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type GroupHistory struct {
	GroupId string `json:"group_id"`
}

type GroupMember struct {
	GroupId string `json:"group_id"`
	UserId  string `json:"user_id"`
}

type ForestLikeReq struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Action string `json:"action"`
}
