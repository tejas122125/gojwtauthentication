package helpers

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func ChechUserType(c *gin.Context, role string) (err error){

userType := c.GetString("user_type")
err = nil
if userType != role{
	err = errors.New("unathorised access to this resources")
	return err
}
return err



}
func MatchTypeToUid(c *gin.Context, userId string) error {
	userType := c.GetString("user_type")
	uid := c.GetString("uid")
	var err error = nil
	if userType == "USER" && uid != userId {
		err = errors.New(" unauthorised to access this route")
		return err
	}
	err = ChechUserType(c, userType)
	return err
}
