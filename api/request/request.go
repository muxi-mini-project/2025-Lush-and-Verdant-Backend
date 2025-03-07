package request

type UserRegister struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	Email      string `json:"email"`
	Code       string `json:"code"`
	Device_Num string `json:"device_num"`
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
