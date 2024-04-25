package patient

type Gender int

const (
	Male   Gender = 0
	Female Gender = 1
)

type Person struct {
	Fullname string `json:"fullname"`
	Birthday string `json:"birthday"`
	Gender   Gender `json:"gender"`
	Guid     string `json:"guid"`
}
