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

type User struct {
	ID       string `json:"id"`
	UserName string `json:"username"`
	Email    string `json:"email"`
}

type PostGoalResponse struct {
	GoalID  uint   `json:"goal_id"`
	TaskIDs []uint `json:"task_ids"`
}

type SloganResponse struct {
	Slogan string `json:"slogan"`
}

type TaskWithChecks struct {
	TaskID    uint   `json:"task_id"`
	Title     string `json:"title"`
	Details   string `json:"details"`
	Completed bool   `json:"completed"`
}

type DailyCount struct {
	DailyCount int `json:"daily_count"`
}
