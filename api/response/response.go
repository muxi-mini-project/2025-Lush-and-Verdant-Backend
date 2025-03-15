package response

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Token struct {
	Token  string `json:"token"`
	UserID string `json:"user_id"`
}

type UpToken struct {
	Token string `json:"token"`
}

type URL struct {
	URL string `json:"url"`
}

type URLs struct {
	URLs []URL `json:"urls"`
}

type Goals struct {
	Goals map[string][]map[string]string `json:"goals"`
}

type PostGoalResponse struct {
	GoalID  string   `json:"goal_id"`
	TaskIDs []string `json:"task_ids"`
}

type SloganResponse struct {
	Slogan string `json:"slogan"`
}

type TaskWithChecks struct {
	TaskID    string `json:"task_id"`
	Title     string `json:"title"`
	Details   string `json:"details"`
	Completed bool   `json:"completed"`
}

type DailyCount struct {
	DailyCount int `json:"daily_count"`
}

type GroupInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"` //任务是否公开
	GroupOwner  string `json:"group_owner"`
}

type GroupInfos struct {
	Nums   int         `json:"nums"`
	Groups []GroupInfo `json:"groups"`
}

type User struct {
	ID         string `json:"id""`
	UserName   string `json:"username"`
	Email      string `json:"email"`
	GoalPublic bool   `json:"goal_public"`
	Slogan     string `json:"slogan"`
}

type Users struct {
	Nums  int    `json:"nums"`
	Users []User `json:"users"`
}

type Message struct {
	From    string `json:"from"`
	To      string `json:"to"` // 可以是用户ID或群ID
	Content string `json:"content"`
	Type    string `json:"type"`
	Time    string `json:"time"`
}

type Messages struct {
	Messages []Message `json:"messages"`
}
