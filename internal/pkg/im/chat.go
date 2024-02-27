package im

import (
	"VitaTaskGo/pkg/config"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func SendUser(userid string, msg interface{}) error {
	postData := map[string]interface{}{
		"msg":    msg,
		"userid": userid,
	}
	jsonData, jsonErr := json.Marshal(postData)
	if jsonErr != nil {
		return jsonErr
	}

	// 获取网关host
	host := getHost()
	// 向网关发送请求
	resp, err := http.Post(host+"/gateway/send/user", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("failed: http error: %d", resp.StatusCode))
	}

	return nil
}

func SendUsers(users []string, msg interface{}) error {
	postData := map[string]interface{}{
		"msg":   msg,
		"users": users,
	}
	jsonData, jsonErr := json.Marshal(postData)
	if jsonErr != nil {
		return jsonErr
	}

	// 获取网关host
	host := getHost()
	// 向网关发送请求
	resp, err := http.Post(host+"/gateway/send/users", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("failed: http error: %d", resp.StatusCode))
	}

	return nil
}

// 判断host是否包含了协议
func getHost() string {
	host := config.Get().Gateway.Host
	port := config.Get().Gateway.Port
	if !strings.HasPrefix(host, "http") {
		host = "http://" + host
	}

	return fmt.Sprintf("%s:%d", host, port)
}
