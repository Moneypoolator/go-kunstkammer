package models

import (
	"fmt"
)

type User struct {
	ID        int    `json:"id"`
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func PrintUser(user User) {
	fmt.Println("User Details:")
	fmt.Printf("ID: %d\n", user.ID)
	fmt.Printf("FullName: %s\n", user.FullName)
	fmt.Printf("Email: %s\n", user.Email)
	fmt.Printf("First Name: %s\n", user.FirstName)
	fmt.Printf("Last Name: %s\n", user.LastName)
}

// FindUserByEmail ищет пользователя по email в массиве
func FindUserByEmail(users []User, email string) (*User, error) {
	for _, user := range users {
		if user.Email == email {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("user with email '%s' not found", email)
}
