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
	Id  int    `json:"id"` //是userId或者groupId
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
