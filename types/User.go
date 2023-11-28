package Type

import (
	"errors"
	"strconv"
)

type User struct {
	StudentID  int    `json:"studentID"`
	ProfilePic string `json:"profilePic"`
	Bio        string `json:"bio"`
	Email      string `json:"email"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Phone      int    `json:"phone"`
	Role       string `json:"role"`
	Points     int    `json:"points"`
	Github     string `json:"github"`
	UID        string `json:"uid"`
}

func (u *User) Validate() error {
	studentIDStr := strconv.Itoa(u.StudentID)
	if len(studentIDStr) == 0 {
		return errors.New("StudentID is required")
	}
	if u.Email == "" {
		return errors.New("Email is required")
	}
	if u.Username == "" {
		return errors.New("Username is required")
	}
	if u.Password == "" {
		return errors.New("Password is required")
	}
	if u.FirstName == "" {
		return errors.New("FirstName is required")
	}
	if u.LastName == "" {
		return errors.New("LastName is required")
	}
	if u.Role == "" {
		return errors.New("Role is required")
	}
	if u.UID == "" {
		return errors.New("UID is required")
	}
	return nil
}
