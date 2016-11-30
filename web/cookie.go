package web

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/utils"
	"encoding/base64"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func setUserCookie(w http.ResponseWriter, userId string, username string, userType int) error {
	encCookie, err := encryptUserCookie(userId, username, userType)
	if err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "SESS",
		Value:    encCookie,
		Path:     "/",
		Expires:  time.Now().AddDate(1, 0, 0),
		MaxAge:   365 * 24 * 60 * 60,
		HttpOnly: true,
	})
	return nil
}

func getUserCookie(r *http.Request) (userId string, username string, userType int, err error) {
	var cookie *http.Cookie
	if cookie, err = r.Cookie("SESS"); err != nil {
		return
	}
	userId, username, userType, err = decryptUserCookie(cookie.Value)
	return
}

func clearUserCookie(w http.ResponseWriter) error {
	http.SetCookie(w, &http.Cookie{
		Name:   "SESS",
		Path:   "/",
		MaxAge: -1,
	})
	return nil
}

func encryptUserCookie(userId string, username string, userType int) (string, error) {
	key := []byte("U)Z9D>F]+LP4p.v4")
	cookie := userId + "|" + username + "|" + strconv.Itoa(userType)
	aesCookie, err := utils.AesEncrypt([]byte(cookie), key, key)
	if err != nil {
		return "", err
	}
	base64Cookie := base64.StdEncoding.EncodeToString(aesCookie)
	return base64Cookie, nil
}

func decryptUserCookie(encCookie string) (string, string, int, error) {
	key := []byte("U)Z9D>F]+LP4p.v4")
	aesCookie, err := base64.StdEncoding.DecodeString(encCookie)
	if err != nil {
		return "", "", 0, err
	}

	decCookie, err := utils.AesDecrypt(aesCookie, key, key)
	if err != nil {
		return "", "", 0, err
	}

	splits := strings.SplitN(string(decCookie), "|", 3)

	userType, err := strconv.Atoi(splits[2])
	if err != nil {
		return "", "", 0, err
	}
	userId := splits[0]
	username := splits[1]
	return userId, username, userType, nil
}
