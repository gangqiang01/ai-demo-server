package api

import (
	"github.com/edgehook/ithings/common/dbm/model"
	v1 "github.com/edgehook/ithings/common/types/v1"
	"github.com/edgehook/ithings/common/utils"
	"github.com/edgehook/ithings/webserver/api/jwt"
	responce "github.com/edgehook/ithings/webserver/types"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var auth v1.Auth
	if err := c.Bind(&auth); err != nil {
		responce.FailWithMessage("Parameter error", c)
		return
	}

	user, err := model.GetUserByName(auth.Username)
	if err != nil {
		responce.FailWithMessage("User does not exist", c)
		return
	}

	type resp struct {
		AccessToken string `form:"accessToken" json:"accessToken"`
	}
	enPwd := utils.Md5V(auth.Password)

	if enPwd != user.Password {
		responce.FailWithMessage("Password error", c)
		return
	}

	token, err := jwt.GenerateToken(auth.Username)
	if err != nil {
		responce.FailWithMessage("Grnerate token error", c)
		return
	}
	responce.OkWithData(&resp{
		AccessToken: token,
	}, c)
}
