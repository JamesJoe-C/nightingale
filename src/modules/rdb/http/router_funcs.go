package http

import (
	"fmt"
	"strconv"

	"github.com/didi/nightingale/src/models"
	"github.com/didi/nightingale/src/toolkits/i18n"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/errors"
)

func dangerous(v interface{}) {
	errors.Dangerous(v)
}

func bomb(format string, a ...interface{}) {
	errors.Bomb(i18n.Sprintf(format, a...))
}

func bind(c *gin.Context, ptr interface{}) {
	dangerous(c.ShouldBindJSON(ptr))
}

func urlParamStr(c *gin.Context, field string) string {
	val := c.Param(field)

	if val == "" {
		bomb("url param[%s] is blank", field)
	}

	return val
}

func urlParamInt64(c *gin.Context, field string) int64 {
	strval := urlParamStr(c, field)
	intval, err := strconv.ParseInt(strval, 10, 64)
	if err != nil {
		bomb("cannot convert %s to int64", strval)
	}

	return intval
}

func urlParamInt(c *gin.Context, field string) int {
	return int(urlParamInt64(c, field))
}

func queryStr(c *gin.Context, key string, defaultVal ...string) string {
	val := c.Query(key)
	if val != "" {
		return val
	}

	if len(defaultVal) == 0 {
		bomb("query param[%s] is necessary", key)
	}

	return defaultVal[0]
}

func queryInt(c *gin.Context, key string, defaultVal ...int) int {
	strv := c.Query(key)
	if strv != "" {
		intv, err := strconv.Atoi(strv)
		if err != nil {
			bomb("cannot convert [%s] to int", strv)
		}
		return intv
	}

	if len(defaultVal) == 0 {
		bomb("query param[%s] is necessary", key)
	}

	return defaultVal[0]
}

func queryInt64(c *gin.Context, key string, defaultVal ...int64) int64 {
	strv := c.Query(key)
	if strv != "" {
		intv, err := strconv.ParseInt(strv, 10, 64)
		if err != nil {
			bomb("cannot convert [%s] to int64", strv)
		}
		return intv
	}

	if len(defaultVal) == 0 {
		bomb("query param[%s] is necessary", key)
	}

	return defaultVal[0]
}

func offset(c *gin.Context, limit int) int {
	if limit <= 0 {
		limit = 10
	}

	page := queryInt(c, "p", 1)
	return (page - 1) * limit
}

func renderMessage(c *gin.Context, v interface{}) {
	if v == nil {
		c.JSON(200, gin.H{"err": ""})
		return
	}

	switch t := v.(type) {
	case string:
		c.JSON(200, gin.H{"err": i18n.Sprintf(t)})
	case error:
		c.JSON(200, gin.H{"err": t.Error()})
	}
}

func renderData(c *gin.Context, data interface{}, err error) {
	if err == nil {
		c.JSON(200, gin.H{"dat": data, "err": ""})
		return
	}

	renderMessage(c, err.Error())
}

func renderZeroPage(c *gin.Context) {
	renderData(c, gin.H{
		"list":  []int{},
		"total": 0,
	}, nil)
}

type idsForm struct {
	Ids []int64 `json:"ids"`
}

func loginUsername(c *gin.Context) string {
	value, has := c.Get("username")
	if !has {
		bomb("unauthorized")
	}

	if value == nil {
		bomb("unauthorized")
	}

	return value.(string)
}

func loginUser(c *gin.Context) *models.User {
	username := loginUsername(c)

	user, err := models.UserGet("username=?", username)
	dangerous(err)

	if user == nil {
		bomb("unauthorized")
	}

	return user
}

func loginRoot(c *gin.Context) *models.User {
	value, has := c.Get("user")
	if !has {
		bomb("unauthorized")
	}

	return value.(*models.User)
}

func User(id int64) *models.User {
	user, err := models.UserGet("id=?", id)
	if err != nil {
		bomb("cannot retrieve user[%d]: %v", id, err)
	}

	if user == nil {
		bomb("no such user[%d]", id)
	}

	return user
}

func Team(id int64) *models.Team {
	team, err := models.TeamGet("id=?", id)
	if err != nil {
		bomb("cannot retrieve team[%d]: %v", id, err)
	}

	if team == nil {
		bomb("no such team[%d]", id)
	}

	return team
}

func Role(id int64) *models.Role {
	role, err := models.RoleGet("id=?", id)
	if err != nil {
		bomb("cannot retrieve role[%d]: %v", id, err)
	}

	if role == nil {
		bomb("no such role[%d]", id)
	}

	return role
}

func Node(id int64) *models.Node {
	node, err := models.NodeGet("id=?", id)
	dangerous(err)

	if node == nil {
		bomb("no such node[%d]", id)
	}

	return node
}

func _e(format string, a ...interface{}) error {
	return fmt.Errorf(i18n.Sprintf(format, a...))
}

func _s(format string, a ...interface{}) string {
	return i18n.Sprintf(format, a...)
}
