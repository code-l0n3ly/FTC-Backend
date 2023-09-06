package main

type User struct {
	StudentID string `json:"studentID"`
	Bio       string `json:"bio"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
	Role      string `json:"role"`
	UID       string `json:"uid"`
}
