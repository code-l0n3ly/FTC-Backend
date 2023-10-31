package main

type User struct {
	StudentID  int    `json:"studentID"`
	ProfilePic string `json:"profilePic"`
	Bio        string `json:"bio"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Phone      int    `json:"phone"`
	Role       string `json:"role"`
	Points     int    `json:"points"`
	Github     string `json:"github"`
	UID        string `json:"uid"`
}

type Event struct {
	ID           string `json:"id"`
	Banner       string `json:"banner"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Location     string `json:"location"`
	StartTime    string `json:"startTime"`
	EndTime      string `json:"endTime"`
	OrganizerUID string `json:"organizerUID"`
	WhatsappLink string `json:"whatsappLink"`
}
