package web

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/config"
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/utils"
	"encoding/base64"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func setUserCookie(w http.ResponseWriter, userId string, username string, userType int) error {
	if config.Instance().IsSmockServer() {
		http.SetCookie(w, &http.Cookie{
			Name:     "user_id",
			Value:    userId,
			Path:     "/",
			Expires:  time.Now().AddDate(1, 0, 0),
			MaxAge:   365 * 24 * 60 * 60,
			HttpOnly: true,
		})
		http.SetCookie(w, &http.Cookie{
			Name:     "username",
			Value:    username,
			Path:     "/",
			Expires:  time.Now().AddDate(1, 0, 0),
			MaxAge:   365 * 24 * 60 * 60,
			HttpOnly: true,
		})
		http.SetCookie(w, &http.Cookie{
			Name:     "user_type",
			Value:    strconv.Itoa(userType),
			Path:     "/",
			Expires:  time.Now().AddDate(1, 0, 0),
			MaxAge:   365 * 24 * 60 * 60,
			HttpOnly: true,
		})
	} else {
		encCookie, err := encryptUserCookie(userId, username, userType)
		if err != nil {
			return err
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "USER",
			Value:    encCookie,
			Path:     "/",
			Expires:  time.Now().AddDate(1, 0, 0),
			MaxAge:   365 * 24 * 60 * 60,
			HttpOnly: true,
		})
	}
	return nil
}

func getUserCookie(r *http.Request) (userId string, username string, userType int, err error) {
	if config.Instance().IsSmockServer() {
		var cookie *http.Cookie
		if cookie, err = r.Cookie("user_id"); err != nil {
			return
		} else {
			userId = cookie.Value
		}
		if cookie, err = r.Cookie("username"); err != nil {
			return
		} else {
			username = cookie.Value
		}
		if cookie, err = r.Cookie("user_type"); err != nil {
			return
		} else {
			var ut int
			if ut, err = strconv.Atoi(cookie.Value); err != nil {
				return
			} else {
				userType = ut
			}
		}
	} else {
		var cookie *http.Cookie
		if cookie, err = r.Cookie("USER"); err != nil {
			return
		}
		userId, username, userType, err = decryptUserCookie(cookie.Value)
	}
	return
}

func clearUserCookie(w http.ResponseWriter) error {
	if config.Instance().IsSmockServer() {
		http.SetCookie(w, &http.Cookie{
			Name:   "user_id",
			Path:   "/",
			MaxAge: -1,
		})
		http.SetCookie(w, &http.Cookie{
			Name:   "username",
			Path:   "/",
			MaxAge: -1,
		})
		http.SetCookie(w, &http.Cookie{
			Name:   "user_type",
			Path:   "/",
			MaxAge: -1,
		})
	} else {
		http.SetCookie(w, &http.Cookie{
			Name:   "USER",
			Path:   "/",
			MaxAge: -1,
		})
	}
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
