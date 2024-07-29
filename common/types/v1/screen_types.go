package v1

type ScreenData struct {
	Name        string `form:"name" json:"name,omitempty"`
	Description string `form:"description" json:"description,omitempty"`
	Data        string `form:"data" json:"data,omitempty"`
	State       *int   `form:"state" json:"state,omitempty"`
	Image       string `form:"image" json:"image,omitempty"`
}

type ScreenImageData struct {
	ScreenId string `form:"screenId" json:"screenId,omitempty"`
	Type     string `form:"type" json:"type,omitempty"`
	Data     string `form:"data" json:"data,omitempty"`
	Filename string `form:"filename" json:"filename,omitempty"`
}
