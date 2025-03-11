package response

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Token struct {
	Token string `json:"token"`
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
