package tpl

const (
	CreateSuccess = `
创建成功！
{{ .Name }}
公历生日： {{ .SolarBirthday.Format "2006-01-02" }}
农历生日： {{ .LunarBirthday.Format "2006-01-02" }}
`
	UpdateSuccess = `
更新成功！
{{ .Name }}
公历生日： {{ .SolarBirthday.Format "2006-01-02" }}
农历生日： {{ .LunarBirthday.Format "2006-01-02" }}
`
	IncomingBirthday = `
【{{ .Name }}】 快要过生日啦！
公历生日： {{ .SolarBirthday.Format "2006-01-02" }}
农历生日： {{ .LunarBirthday.Format "2006-01-02" }}
`
)
