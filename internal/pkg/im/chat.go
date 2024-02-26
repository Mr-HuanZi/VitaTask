package im

import (
	"VitaTaskGo/pkg/config"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var host = ""

func init() {
	host = config.Get().Gateway.Host
}

func SendUser(userid string, msg string) error {
	postData := map[string]interface{}{
		"msg":    msg,
		"userid": userid,
	}
	jsonData, jsonErr := json.Marshal(postData)
	if jsonErr != nil {
		return jsonErr
	}
	resp, err := http.Post(host+"/send/user", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("failed: http error: %d", resp.StatusCode))
	}

	return nil
}

func SendUsers(users []string, msg string) error {
	postData := map[string]interface{}{
		"msg":   msg,
		"users": users,
	}
	jsonData, jsonErr := json.Marshal(postData)
	if jsonErr != nil {
		return jsonErr
	}
	resp, err := http.Post(host+"/send/users", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("failed: http error: %d", resp.StatusCode))
	}

	return nil
}
