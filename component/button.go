package component

type Button interface {
	SetFrontCall(frontcall map[string]interface{}) error
	SetLabel(Label string) error
}

type WxAppButton struct {
	ID           string                 `json:"ID"`
	Label        string                 `json:"Label"`
	Describe     string                 `json:"Describe"`
	OpenType     string                 `json:"OpenType"`
	Sort         int64                  `json:"Sort"`
	Color        string                 `json:"Color"` //颜色
	Size         string                 `json:"Size"`  //字体
	Disable      bool                   `json:"Disable"`
	Notification string                 `json:"-"`
	FrontCall    map[string]interface{} `json:"FrontCall"` //前端功能调用
}

func (this *WxAppButton) OnInit(ID, Label, Describe string, Sort int64) error {
	this.ID = ID
	this.Label = Label
	this.Describe = Describe
	this.Sort = Sort
	return nil
}
func (this *WxAppButton) SetFrontCall(frontcall map[string]interface{}) error {
	this.FrontCall = frontcall
	return nil
}

func (this *WxAppButton) SetLabel(Label string) error {
	this.Label = Label
	return nil
}
