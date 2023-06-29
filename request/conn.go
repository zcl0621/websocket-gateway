package request

type ConnectQuery struct {
	UserId   string `form:"user_id" binding:"required"`
	Platform string `form:"-"`
	OS       string `form:"os"`       // ios android windows mac
	Version  string `form:"version"`  // 16.3 10 11 7 xp vista
	AppType  string `form:"app_type"` // higo airschool
	ClientIp string `form:"-"`
}
