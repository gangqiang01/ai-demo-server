package v1

const (
	NVRChannelTypeRtmp  string = "rtmp"
	NVRChannelTypeRtsp  string = "rtsp"
	NVRChannelTypeHls   string = "hls"
	NVRChannelTypeFlv   string = "flv"
	NVRChannelTypeOnvif string = "onvif"
)

type NVR_CHANNEL struct {
	Name       string `form:"name" json:"name"`
	AccessType string `form:"accessType" json:"accessType"`
	Status     string `form:"status" json:"status"`
	Config     string `form:"config" json:"config"`
	Duration   *int   `form:"duration" json:"duration"`
	Segment    *int   `form:"segment" json:"segment"`
}

type NVR_CONFIG struct {
	Address  string `form:"address" json:"address"`
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
}
