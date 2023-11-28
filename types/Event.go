package Type

type Event struct {
	ID           string   `json:"id"`
	Banner       string   `json:"banner"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Location     string   `json:"location"`
	StartTime    string   `json:"startTime"`
	EndTime      string   `json:"endTime"`
	OrganizerUID string   `json:"organizerUID"`
	MemberUID    []string `json:"MemberUID"`
	Size         int      `json:"Size"`
	WhatsappLink string   `json:"whatsappLink"`
}
