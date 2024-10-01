package main

import (
	"fmt"
)

var authCookie string
var authToken string
var avatarID string

var (
	Username string
	Password string
	code     string
)

var choice string

func main() {
	authToken, err := loadAuthToken()
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		if loginSessionCheck(authToken) {
			fmt.Println("Login session is valid")
			break
		}

		fmt.Println("Enter your username: ")
		fmt.Scan(&Username)
		fmt.Println("Enter your password: ")
		fmt.Scan(&Password)
		fmt.Println("Making POST request...")

		Login(Username, Password)
		if authCookie != "" {
			break
		}

		fmt.Println("Login failed. Would you like to retry? (y/n)")
		fmt.Scan(&choice)
		if choice != "y" {
			fmt.Println("Exiting program.")
			return
		}
	}

	for {
		if loginSessionCheck(authToken) {
			fmt.Println("Login session is valid")
			break
		}

		fmt.Println("Enter the 2FA email OTP code: ")
		fmt.Scan(&code)
		VerifyEmailOTP(code)
	}

	fmt.Println("input AvatarID | ")
	fmt.Scan(&avatarID)

	getAvatar(avatarID)

	/*
		for {

			fmt.Println("Enter the Avatar ID you want to select: ")
			fmt.Scan(&avatarID)

			ChangeAvatar(authCookie, avatarID)


			fmt.Println("input userID | ")
			fmt.Scan(&userID)

			sendFriendRq(authCookie, userID)

			fmt.Println("Select a different Avatar? Enter 'y' to continue, 'n' to exit program: ")
			fmt.Scan(&choice)

			if choice != "y" {
				fmt.Println("Exiting program.")
				return
			}

		}
	*/
}
