package middleware

import (
	"asset-management/app/dao"
	"asset-management/app/define"
	"asset-management/utils"
	"bytes"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var customLog *logrus.Logger

func init() {
	customFormatter := new(logrus.JSONFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customLog = logrus.New()
	customLog.SetFormatter(customFormatter)
	customLog.AddHook(dao.MysqlHook)
}

func getOperatorInfo(ctx *utils.Context) *define.UserBasicInfo {
	userInfo, exists := ctx.Get("user")
	if exists {
		if userInfo, ok := userInfo.(define.UserBasicInfo); ok {
			return &userInfo
		}
	}
	return nil
}

type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w CustomResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w CustomResponseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func LogMiddleware() utils.HandlerFunc {
	return func(ctx *utils.Context) {
		if ctx.Request.URL.Path == "/user/login" {
			blw := &CustomResponseWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
			ctx.Writer = blw
			ctx.Next()
			resData := utils.ResponseData{}
			_ = json.Unmarshal(blw.body.Bytes(), &resData)
			if ctx.Writer.Status() == 200 {
				loginRes := resData.Data.(map[string]interface{})
				userInfo := loginRes["user"].(map[string]interface{})
				customLog.WithFields(logrus.Fields{
					"method":        ctx.Request.Method,
					"url":           ctx.Request.URL.Path,
					"status":        ctx.Writer.Status(),
					"error_code":    resData.Error.Code,
					"error_message": resData.Error.Message,
					"user_id":       uint(userInfo["user_id"].(float64)),
					"username":      userInfo["username"].(string),
					"entity_id":     uint(userInfo["entity_id"].(float64)),
					"department_id": uint(userInfo["department_id"].(float64)),
				}).Info("Successfully login")
			} else {
				customLog.WithFields(logrus.Fields{
					"method":        ctx.Request.Method,
					"url":           ctx.Request.URL.Path,
					"status":        ctx.Writer.Status(),
					"error_code":    resData.Error.Code,
					"error_message": resData.Error.Message,
					"user_id":       0,
					"username":      "",
					"entity_id":     0,
					"department_id": 0,
				})
				customLog.Info("Login failed")
			}
		} else if ctx.Request.Method != "GET" {
			userInfo := getOperatorInfo(ctx)
			blw := &CustomResponseWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
			ctx.Writer = blw
			ctx.Next()
			resData := utils.ResponseData{}
			_ = json.Unmarshal(blw.body.Bytes(), &resData)
			if ctx.Writer.Status() == 200 {
				customLog.WithFields(logrus.Fields{
					"method":        ctx.Request.Method,
					"url":           ctx.Request.URL.Path,
					"status":        ctx.Writer.Status(),
					"error_code":    resData.Error.Code,
					"error_message": resData.Error.Message,
					"user_id":       userInfo.UserID,
					"username":      userInfo.UserName,
					"entity_id":     userInfo.EntityID,
					"department_id": userInfo.DepartmentID,
				}).Info("Operation succeed")
			} else {
				customLog.WithFields(logrus.Fields{
					"method":        ctx.Request.Method,
					"url":           ctx.Request.URL.Path,
					"status":        ctx.Writer.Status(),
					"error_code":    resData.Error.Code,
					"error_message": resData.Error.Message,
					"user_id":       userInfo.UserID,
					"username":      userInfo.UserName,
					"entity_id":     userInfo.EntityID,
					"department_id": userInfo.DepartmentID,
				}).Info("Operation fail")
			}
		}
	}
}
