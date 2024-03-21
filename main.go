package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

var authCookie string

func Login(user string, password string) {
	loginURL := "https://api.vrchat.cloud/api/1/auth/user"

	req, err := http.NewRequest("GET", loginURL, nil)
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return
	}

	req.SetBasicAuth(user, password)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("User-Agent", "golang Client")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Failed to authenticate. Status Code:", resp.StatusCode)
		return
	}

	defer resp.Body.Close()

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "auth" {
			authCookie = cookie.Value
			break
		}
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	fmt.Println("Status:", resp.Status)
	fmt.Println("Response body:", string(respBody))
}
func VerifyEmailOTP(code string) {
	apiUrl := "https://api.vrchat.cloud/api/1/auth/twofactorauth/emailotp/verify"

	requestBody, err := json.Marshal(map[string]string{
		"code": code,
	})
	if err != nil {
		fmt.Println("Error encoding request body:", err)
		return
	}

	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("User-Agent", "golang Client")
	if authCookie != "" {
		req.Header.Set("Cookie", "auth="+authCookie)
	} else {
		fmt.Println("Missing auth cookie. Please login first.")
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	formattedData := formatJSON(respBody)
	fmt.Println("Status:", resp.Status)
	fmt.Println("Response body:", formattedData)
}

func formatJSON(data []byte) string {
	var out bytes.Buffer
	err := json.Indent(&out, data, "", " ")
	if err != nil {
		fmt.Println("Error formatting JSON:", err)
	}

	return out.String()
}

func main() {
	var Username string
	var Password string

	fmt.Println("Enter your username: ")
	fmt.Scan(&Username)
	fmt.Println("Enter your password: ")
	fmt.Scan(&Password)
	fmt.Println("Making POST request...")

	Login(Username, Password)

	var code string
	fmt.Println("Enter the 2FA email OTP code: ")
	fmt.Scan(&code)
	VerifyEmailOTP(code)
}
