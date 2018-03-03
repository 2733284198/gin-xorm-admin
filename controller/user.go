package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/angao/gin-xorm-admin/db"
	"github.com/angao/gin-xorm-admin/forms"
	"github.com/angao/gin-xorm-admin/models"
	"github.com/angao/gin-xorm-admin/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// UserController handle user request
type UserController struct {
}

// Home user home page
func (UserController) Home(c *gin.Context) {
	r.HTML(c.Writer, http.StatusOK, "system/user/user.html", gin.H{})
}

// List query all user
func (UserController) List(c *gin.Context) {
	var userDao db.UserDao

	var userForm forms.UserForm
	if err := c.Bind(&userForm); err != nil {
		r.JSON(c.Writer, http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	users, err := userDao.List(userForm)
	if err != nil {
		r.JSON(c.Writer, http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	r.JSON(c.Writer, http.StatusOK, gin.H{
		"data": users,
	})
}

// Info is handle user info
func (UserController) Info(c *gin.Context) {
	var userDao db.UserDao
	var err error
	var user *models.UserRole
	session := sessions.Default(c)
	id, ok := session.Get("user_id").(int64)
	if ok {
		user, err = userDao.GetUserRole(id)
		if err != nil {
			r.JSON(c.Writer, http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		r.HTML(c.Writer, http.StatusOK, "container.html", gin.H{
			"user":     user.User,
			"roleName": user.Role.Name,
		})
		return
	}
	r.HTML(c.Writer, http.StatusInternalServerError, "container.html", gin.H{
		"error": err.Error(),
	})
}

// ToAdd handle add user page
func (UserController) ToAdd(c *gin.Context) {
	r.HTML(c.Writer, http.StatusOK, "system/user/user_add.html", gin.H{})
}

// ToEdit handle edit user paget
func (UserController) ToEdit(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		r.JSON(c.Writer, http.StatusBadRequest, gin.H{
			"error": "参数错误",
		})
		return
	}
	var userDao db.UserDao
	pid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		r.JSON(c.Writer, http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	user, err := userDao.GetUserRole(pid)
	if err != nil {
		r.JSON(c.Writer, http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	r.HTML(c.Writer, http.StatusOK, "system/user/user_edit.html", gin.H{
		"user":     user.User,
		"roleName": user.Role.Name,
	})
}

// ToRoleAssign handle user role
func (UserController) ToRoleAssign(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		r.JSON(c.Writer, http.StatusBadRequest, gin.H{
			"error": "参数错误",
		})
		return
	}
	var userDao db.UserDao
	pid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		r.JSON(c.Writer, http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	user, err := userDao.GetUserByID(pid)
	if err != nil {
		r.JSON(c.Writer, http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	r.HTML(c.Writer, http.StatusOK, "system/user/user_roleassign.html", gin.H{
		"user": user,
	})
}

// Add handle save user
func (UserController) Add(c *gin.Context) {
	var userDao db.UserDao
	var userAddForm forms.UserAddForm
	var user models.User

	if err := c.Bind(&userAddForm); err != nil {
		r.JSON(c.Writer, http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	log.Printf("user: %#v\n", userAddForm)
	if userAddForm.Password != userAddForm.RePassword {
		r.JSON(c.Writer, http.StatusBadRequest, gin.H{
			"error": "密码不一致",
		})
		return
	}
	user.Name = userAddForm.Name
	user.Account = userAddForm.Account
	user.Email = userAddForm.Email
	user.Sex = userAddForm.Sex
	randomStr := utils.RandomString(5)
	password, err := utils.Encrypt(userAddForm.Password, randomStr)
	if err != nil {
		r.JSON(c.Writer, http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	user.Password = password
	user.Salt = randomStr
	err = userDao.Save(user)
	if err != nil {
		r.JSON(c.Writer, http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	r.JSON(c.Writer, http.StatusOK, gin.H{
		"message": err.Error(),
	})
}

// Delete 删除用户
func (UserController) Delete(c *gin.Context) {
	id := c.PostForm("userId")
	if id == "" {
		r.JSON(c.Writer, http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	var userDao db.UserDao
	pid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		r.JSON(c.Writer, http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	err = userDao.Delete(pid)
	if err != nil {
		r.JSON(c.Writer, http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	r.JSON(c.Writer, http.StatusOK, "")
}

// Reset password
func (UserController) Reset(c *gin.Context) {
	id := c.PostForm("userId")
	if id == "" {
		r.JSON(c.Writer, http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	var userDao db.UserDao
	pid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		r.JSON(c.Writer, http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	var user *models.User
	user, err = userDao.GetUserByID(pid)
	if err != nil {
		r.JSON(c.Writer, http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	user.Id = pid
	password, err := utils.Encrypt("111111", user.Salt)
	if err != nil {
		r.JSON(c.Writer, http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	user.Password = password
	err = userDao.Update(user)
	if err != nil {
		r.JSON(c.Writer, http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	r.JSON(c.Writer, http.StatusOK, "")
}

// SetRole set user role
func (UserController) SetRole(c *gin.Context) {
	roleIDs := c.PostForm("roleIds")
	userID := c.PostForm("userId")

	if roleIDs == "" || userID == "" {
		r.JSON(c.Writer, http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	var userDao db.UserDao
	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		r.JSON(c.Writer, http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	user := models.User{
		Id:     id,
		RoleId: roleIDs,
	}
	err = userDao.Update(&user)
	if err != nil {
		r.JSON(c.Writer, http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	r.JSON(c.Writer, http.StatusOK, gin.H{
		"message": "success",
	})
}

// Freeze user
func (UserController) Freeze(c *gin.Context) {
	userID := c.PostForm("userId")
	if userID == "" {
		r.JSON(c.Writer, http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	var userDao db.UserDao
	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		r.JSON(c.Writer, http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	user := models.User{
		Id:     id,
		Status: 2,
	}
	err = userDao.Update(&user)
	if err != nil {
		r.JSON(c.Writer, http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	r.JSON(c.Writer, http.StatusOK, gin.H{
		"message": "success",
	})
}

// UnFreeze user
func (UserController) UnFreeze(c *gin.Context) {
	userID := c.PostForm("userId")
	if userID == "" {
		r.JSON(c.Writer, http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	var userDao db.UserDao
	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		r.JSON(c.Writer, http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	user := models.User{
		Id:     id,
		Status: 1,
	}
	err = userDao.Update(&user)
	if err != nil {
		r.JSON(c.Writer, http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	r.JSON(c.Writer, http.StatusOK, gin.H{
		"message": "success",
	})
}