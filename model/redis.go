package model

import (
	"github.com/shudiwsh2009/reservation_thxl_go/config"
	"gopkg.in/redis.v5"
	"log"
)

const (
	REDIS_KEY_TEACHER_RESET_PASSWORD_VERIFY_CODE                  = "thxlfzzx#teacher_reset_password_verify_code_%s"
	REDIS_KEY_ADMIN_CLEAR_ALL_STUDENTS_BINDED_TEACHER_VERIFY_CODE = "thxlfzzx#admin_clear_all_students_bind_teacher_verify_code_%s"
	REDIS_KEY_USER_LOGIN                                          = "thxlfzzx#user_login_%d_%s_%s"
)

func NewRedisClient() *redis.Client {
	var client *redis.Client
	if config.Instance().IsSmockServer() {
		client = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})
	} else {
		client = redis.NewClient(&redis.Options{
			Addr:     config.Instance().RedisAddress,
			Password: config.Instance().RedisPassword,
			DB:       config.Instance().RedisDatabase,
		})
	}
	pong, err := client.Ping().Result()
	if err != nil {
		log.Fatalf("连接Redis失败：%v", err)
	}
	log.Printf("连接Redis成功：%s", pong)
	return client
}
