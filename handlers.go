package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

func HomeHandler(c *gin.Context) {
	c.String(http.StatusOK, `欢迎来到 OAuth 2.0 示例应用 /login`)
}

func LoginHandler(c *gin.Context) {
	authCodeURL := fmt.Sprintf("%s?response_type=code&client_id=%s&redirect_uri=%s&state=%s", authURL, clientID, redirectURI, state)
	c.Redirect(http.StatusFound, authCodeURL)
}

func CallbackHandler(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.String(http.StatusBadRequest, "授权失败")
		return
	}

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)

	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		c.String(http.StatusInternalServerError, "请求访问令牌失败: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "读取访问令牌响应失败: %v", err)
		return
	}

	var tokenResp map[string]interface{}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		c.String(http.StatusInternalServerError, "解析访问令牌响应失败: %v", err)
		return
	}

	accessToken, ok := tokenResp["access_token"].(string)
	if !ok {
		c.String(http.StatusInternalServerError, "获取访问令牌失败")
		return
	}

	c.Set("access_token", accessToken)
	c.Redirect(http.StatusFound, "/profile")
}

func ProfileHandler(c *gin.Context) {
	accessToken, exists := c.Get("access_token")
	if !exists {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		c.String(http.StatusInternalServerError, "创建请求失败: %v", err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+accessToken.(string))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.String(http.StatusInternalServerError, "请求用户信息失败: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "读取用户信息响应失败: %v", err)
		return
	}

	var userInfo map[string]interface{}
	if err := json.Unmarshal(body, &userInfo); err != nil {
		c.String(http.StatusInternalServerError, "解析用户信息响应失败: %v", err)
		return
	}

	c.JSON(http.StatusOK, userInfo)
}
