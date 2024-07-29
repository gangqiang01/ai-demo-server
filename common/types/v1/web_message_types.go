package v1

// Alert web api
type Alert struct {
	Name         string `form:"name" json:"name" binding:"required"`
	Description  string `form:"description" json:"description"`
	Level        *int64 `form:"level" json:"level" binding:"required"`
	Notification string `form:"notification" json:"notification"`
}

// Alert log api
type AlertLog struct {
	Record string `form:"record" json:"record"`
	Status *int32 `form:"status" json:"status"`
	Level  *int64 `form:"level" json:"level"`
}

// fluxdb config
type InfluxDbConfig struct {
	Persistence *bool  `form:"persistence" json:"persistence" binding:"required"`
	Address     string `form:"address" json:"address" binding:"required"`
	Username    string `form:"username" json:"username" binding:"required"`
	Password    string `form:"password" json:"password" binding:"required"`
	DbName      string `form:"dbName" json:"dbName" binding:"required"`
	Duration    string `form:"duration" json:"duration" binding:"required"`
}

// device tag
type DeviceTag struct {
	Name string `form:"name" json:"name" binding:"required"`
}

// user
type User struct {
	UserName    string `form:"username" json:"username"`
	Password    string `form:"password" json:"password"`
	OldPassword string `form:"oldPassword" json:"oldPassword"`
	Email       string `form:"email" json:"email"`
	Rule        string `form:"rule" json:"rule"`
}

type Auth struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
}

// process monitor
type ThreadInfos struct {
	Name       string  `form:"cmd" json:"cmd"`
	Pid        int64   `form:"pid" json:"pid"`
	CpuLoading float64 `form:"cpuloading" json:"cpuloading"`
	MemLoading float64 `form:"memloading" json:"memloading"`
	Status     string  `form:"status" json:"status"`
	UserName   string  `form:"username" json:"username"`
}

// monitor
type MonitorData struct {
	Ctype     string         `form:"type" json:"type" binding:"required"`
	Delay     string         `form:"delay" json:"delay"`
	Processes []*ThreadInfos `form:"processes" json:"processes"`
	Status    string         `form:"status" json:"status"`
	//< > =
	Condition string `form:"condition" json:"condition"`
	Value     string `form:"value" json:"value"`
}

// ai model
type AiModelData struct {
	Name        string `form:"name" json:"name"`
	CType       string `form:"type" json:"type"`
	Path        string `form:"path" json:"path"`
	Labels      string `form:"labels" json:"labels"`
	Description string `form:"description" json:"description"`
}

// ai detect
type AiDetectData struct {
	Name         string `form:"name" json:"name"`
	Source       string `form:"source" json:"source"`
	Config       string `form:"config" json:"config"`
	CType        string `form:"type" json:"type"`
	AiModelId    string `form:"aiModelId" json:"aiModelId"`
	IsShow       string `form:"isShow" json:"isShow"`
	Notification string `form:"notification" json:"notification"`
	Status       string `form:"status" json:"status"`
}
