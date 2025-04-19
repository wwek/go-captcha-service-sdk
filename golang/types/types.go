package types

type CaptData struct {
	Id                string `json:"id,omitempty"`
	CaptchaKey        string `json:"captcha_key,omitempty"`
	MasterImageBase64 string `json:"master_image_base64,omitempty"`
	ThumbImageBase64  string `json:"thumb_image_base64,omitempty"`
	MasterImageWidth  int64  `json:"master_width,omitempty"`
	MasterImageHeight int64  `json:"master_height,omitempty"`
	ThumbImageWidth   int64  `json:"thumb_width,omitempty"`
	ThumbImageHeight  int64  `json:"thumb_height,omitempty"`
	ThumbImageSize    int64  `json:"thumb_size,omitempty"`
	DisplayX          int64  `json:"display_x,omitempty"`
	DisplayY          int64  `json:"display_y,omitempty"`
}

type CaptNormalDataResponse struct {
	Code    int32       `json:"code" default:"200"`
	Message string      `json:"message" default:""`
	Data    interface{} `json:"data"`
}

type CaptStatusDataResponse struct {
	Code    int32  `json:"code" default:"200"`
	Message string `json:"message" default:""`
	Data    string `json:"status" default:""`
}

type CaptStatusInfo struct {
	Info   interface{} `json:"info"`
	Type   int64       `json:"type"`
	Status int64       `json:"status"`

	ClickDataMaps  map[int]*ClickData
	SlideDataMaps  map[int]*SlideData
	RotateDataMaps map[int]*RotateData
}

// ClickData .. .
type ClickData struct {
	Index  int    `json:"index"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Size   int    `json:"size"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Text   string `json:"text"`
	Shape  string `json:"shape"`
	Angle  int    `json:"angle"`
	Color  string `json:"color"`
	Color2 string `json:"color2"`
}

// SlideData .
type SlideData struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
	Angle  int `json:"angle"`
	// Deprecated: As of 2.1.0, it will be removed, please use [Block.DX].
	TileX int `json:"tile_x"`
	// Deprecated: As of 2.1.0, it will be removed, please use [Block.DY].
	TileY int `json:"tile_y"`
	// Display x,y
	DX int `json:"dx"`
	DY int `json:"dy"`
}

// RotateData ..
type RotateData struct {
	// Deprecated: As of 2.1.0, it will be removed, please use [[CaptchaInstance].GetOptions().GetImageSize()].
	ParentWidth int `json:"parent_width"`
	// Deprecated: As of 2.1.0, it will be removed, please use [[CaptchaInstance].GetOptions().GetImageSize()].
	ParentHeight int `json:"parent_height"`
	Width        int `json:"width"`
	Height       int `json:"height"`
	Angle        int `json:"angle"`
}
