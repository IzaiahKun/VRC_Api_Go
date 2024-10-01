package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

var authData AuthData

type AuthData struct {
	AuthToken string `json:"auth_token"`
}

func saveSessionCookie(authToken string) error {
	config := AuthData{AuthToken: authToken}

	file, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %v", err)
	}

	err = ioutil.WriteFile("config.json", file, 0644)
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	fmt.Println("Auth token saved to config.json")
	return nil
}

func loadAuthToken() (string, error) {
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		return "", fmt.Errorf("error reading config file: %v", err)
	}

	var config AuthData
	err = json.Unmarshal(file, &config)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	fmt.Println("Loaded auth token from JSON:", config.AuthToken)
	return config.AuthToken, nil
}

func loginSessionCheck(authToken string) bool {
	if authToken == "" {
		fmt.Println("Auth token is empty")
		return false
	}

	getCurrentUser := "https://vrchat.com/api/1/auth/user"

	req, err := http.NewRequest("GET", getCurrentUser, nil)
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return false
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "auth="+authToken)
	req.Header.Add("User-Agent", "golang Client")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return false
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return false
	}

	fmt.Println("Status:", resp.Status)

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Failed to authenticate. Status Code:", resp.StatusCode)
		err = saveSessionCookie("")
		if err != nil {
			fmt.Println("Error resetting auth token:", err)
		}
		return false
	}

	fmt.Println("Response body:", string(respBody))

	return true
}

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
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Failed to authenticate. Status Code:", resp.StatusCode)
		return
	}

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "auth" {
			authCookie = cookie.Value
			authToken = cookie.Value
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

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Failed to authenticate. Status Code:", resp.StatusCode)
		return
	}

	saveSessionCookie(authToken)
}

func VerifyEmailOTP(code string) {
	emailOTPUrl := "https://api.vrchat.cloud/api/1/auth/twofactorauth/emailotp/verify"

	requestBody, err := json.Marshal(map[string]string{
		"code": code,
	})
	if err != nil {
		fmt.Println("Error encoding request body:", err)
		return
	}

	req, err := http.NewRequest("POST", emailOTPUrl, bytes.NewBuffer(requestBody))
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
