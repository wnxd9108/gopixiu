/*
Copyright 2021 The Pixiu Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package user

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/caoyingjunz/gopixiu/api/server/httpstatus"
	"github.com/caoyingjunz/gopixiu/api/server/httputils"
	"github.com/caoyingjunz/gopixiu/api/types"
	"github.com/caoyingjunz/gopixiu/pkg/pixiu"
	"github.com/caoyingjunz/gopixiu/pkg/util"
)

func (u *userRouter) createUser(c *gin.Context) {
	r := httputils.NewResponse()
	var user types.User
	if err := c.ShouldBindJSON(&user); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}
	if err := pixiu.CoreV1.User().Create(context.TODO(), &user); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}

	httputils.SetSuccess(c, r)
}

// 更新用户属性：
// 不允许更改字段:
// 1. 用户名
// 2. 用户密码 —— 通过修改密码API进行修改
func (u *userRouter) updateUser(c *gin.Context) {
	r := httputils.NewResponse()
	var (
		err  error
		user types.User
	)
	if err = c.ShouldBindJSON(&user); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}
	user.Id, err = util.ParseInt64(c.Param("id"))
	if err != nil {
		httputils.SetFailed(c, r, err)
		return
	}
	if err = pixiu.CoreV1.User().Update(context.TODO(), &user); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}

	httputils.SetSuccess(c, r)
}

func (u *userRouter) deleteUser(c *gin.Context) {
	r := httputils.NewResponse()
	uid, err := util.ParseInt64(c.Param("id"))
	if err != nil {
		httputils.SetFailed(c, r, err)
		return
	}
	if err = pixiu.CoreV1.User().Delete(context.TODO(), uid); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}

	httputils.SetSuccess(c, r)
}

func (u *userRouter) getUser(c *gin.Context) {
	r := httputils.NewResponse()
	uid, err := util.ParseInt64(c.Param("id"))
	if err != nil {
		httputils.SetFailed(c, r, err)
		return
	}
	r.Result, err = pixiu.CoreV1.User().Get(context.TODO(), uid)
	if err != nil {
		httputils.SetFailed(c, r, err)
		return
	}

	httputils.SetSuccess(c, r)
}

func (u *userRouter) listUsers(c *gin.Context) {
	r := httputils.NewResponse()
	var err error
	if r.Result, err = pixiu.CoreV1.User().List(context.TODO()); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}

	httputils.SetSuccess(c, r)
}

// login
// 1. 检验用户名和密码是否正确，
// 2. 返回 token
func (u *userRouter) login(c *gin.Context) {
	r := httputils.NewResponse()
	var (
		user types.User
		err  error
	)
	if err = c.ShouldBindJSON(&user); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}
	if r.Result, err = pixiu.CoreV1.User().Login(context.TODO(), &user); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}

	httputils.SetSuccess(c, r)
}

// TODO
func (u *userRouter) logout(c *gin.Context) {}

func (u *userRouter) resetPassword(c *gin.Context) {}

func (u *userRouter) changePassword(c *gin.Context) {
	r := httputils.NewResponse()

	var idOptions types.IdOptions
	if err := c.ShouldBindUri(&idOptions); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}
	// 解析修改密码的三个参数
	//  1. 当前密码 2. 新密码 3. 确认新密码
	var password types.Password
	if err := c.ShouldBindJSON(&password); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}
	password.UserId = idOptions.Id

	// 需要通过 token 中的 id 判断当前操作的用户和需要修改密码的用户是否是同一个
	// Get the uid from token
	if err := pixiu.CoreV1.User().ChangePassword(context.TODO(), c.GetInt64("userId"), &password); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}

	httputils.SetSuccess(c, r)
}

func (u *userRouter) getButtonsByCurrentUser(c *gin.Context) {
	r := httputils.NewResponse()
	uidStr, exist := c.Get("userId")
	if !exist {
		httputils.SetFailed(c, r, httpstatus.NoUserIdError)
		return
	}

	menuId, err := util.ParseInt64(c.Param("id"))
	if err != nil {
		httputils.SetFailed(c, r, httpstatus.ParamsError)
		return
	}
	uid := uidStr.(int64)
	res, err := pixiu.CoreV1.User().GetButtonsByUserID(c, uid, menuId)
	if err != nil {
		httputils.SetFailed(c, r, httpstatus.OperateFailed)
		return
	}
	r.Result = res
	httputils.SetSuccess(c, r)
}

func (u *userRouter) getLeftMenusByCurrentUser(c *gin.Context) {
	uidStr, exist := c.Get("userId")
	r := httputils.NewResponse()
	if !exist {
		httputils.SetFailed(c, r, httpstatus.NoUserIdError)
		return
	}
	var err error
	uid := uidStr.(int64)
	r.Result, err = pixiu.CoreV1.User().GetLeftMenusByUserID(c, uid)
	if err != nil {
		httputils.SetFailed(c, r, httpstatus.OperateFailed)
		return
	}

	httputils.SetSuccess(c, r)
}

func (u *userRouter) getUserRoles(c *gin.Context) {
	r := httputils.NewResponse()
	uid, err := util.ParseInt64(c.Param("id"))
	if err != nil {
		httputils.SetFailed(c, r, httpstatus.ParamsError)
		return
	}
	result, err := pixiu.CoreV1.User().GetRoleIDByUser(c, uid)
	if err != nil {
		httputils.SetFailed(c, r, httpstatus.OperateFailed)
		return
	}
	r.Result = result
	httputils.SetSuccess(c, r)
}

func (u *userRouter) setUserRoles(c *gin.Context) {
	var roles types.Roles
	r := httputils.NewResponse()
	err := c.ShouldBindJSON(&roles)
	if err != nil {
		httputils.SetFailed(c, r, httpstatus.ParamsError)
		return
	}

	uid, err := util.ParseInt64(c.Param("id"))
	if err != nil {
		httputils.SetFailed(c, r, httpstatus.ParamsError)
		return
	}

	res, err := pixiu.CoreV1.User().Get(c, uid)
	if err != nil || res == nil {
		httputils.SetFailed(c, r, httpstatus.ParamsError)
		return
	}

	err = pixiu.CoreV1.User().SetUserRoles(c, uid, roles.RoleIds)
	if err != nil {
		httputils.SetFailed(c, r, httpstatus.OperateFailed)
		return
	}
	httputils.SetSuccess(c, r)
}
