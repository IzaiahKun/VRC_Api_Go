package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

var userID string

func sendFriendRq(userID string) {
	authToken, err := loadAuthToken()
	if err != nil {
		fmt.Println(err)
		return
	}

	friendRqApiURL := "https://vrchat.com/api/1/user/" + userID + "/friendRequest"
	fmt.Println(friendRqApiURL)

	req, err := http.NewRequest("POST", friendRqApiURL, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "auth="+authToken)
	req.Header.Add("User-Agent", "golang Client")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Failed to send user a friend request. Status Code:", resp.StatusCode)
		return
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	formattedData := formatJSON(respBody)
	fmt.Println("Status:", resp.Status)
	fmt.Println("Response body:", formattedData)

}

func cancelFriendrq(userID string) {
	authToken, err := loadAuthToken()
	if err != nil {
		fmt.Println(err)
		return
	}

	cancelFriendrqApiURl := "https://vrchat.com/api/1/user/" + userID + "/friendRequest"

	req, err := http.NewRequest("DELETE", cancelFriendrqApiURl, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "auth="+authToken)
	req.Header.Add("User-Agent", "golang Client")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Failed to cancel friend request. Status Code:", resp.StatusCode)
		return
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	formattedData := formatJSON(respBody)
	fmt.Println("Status:", resp.Status)
	fmt.Println("Response body:", formattedData)

}
