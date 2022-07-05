package controllers

import (
	"boilerplate/apps/user/services"
	"boilerplate/core/infrastructures"
	"boilerplate/core/interfaces"
	"boilerplate/core/models"
	"boilerplate/core/responses"
	"boilerplate/core/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController struct {
	logger      interfaces.Logger
	env         *infrastructures.Env
	userService *services.UserService
}

func NewUserController(logger *infrastructures.Logger,
	env *infrastructures.Env,
	userService *services.UserService,
) *UserController {
	return &UserController{
		logger:      logger,
		env:         env,
		userService: userService,
	}
}

// @Summary get users list
// @Schemes
// @Description list of paginated response , authentication required
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} swagger.UsersListResponse
// @failure 401 {object} swagger.UnauthenticatedResponse
// @failure 403 {object} swagger.AccessForbiddenResponse
// @Router /users [get]
func (uc UserController) ListUser(c *gin.Context) {
	uc.paginateUserList(c, "")
}

// @Summary create users
// @Schemes
// @Description create user and admin , admin only
// @Tags admin
// @Accept json
// @Produce json
// @Param email formData string true "unique email"
// @Param password formData string true "password that have at least 8 length and contain an alphabet and number "
// @Param repeatPassword formData string true "repeatPassword that have at least 8 length and contain an alphabet and number "
// @Param firstName formData string true "firstName"
// @Param lastName formData string true "lastName"
// @Param isAdmin formData bool true "isAdmin"
// @Success 200 {object} swagger.UsersListResponse
// @failure 401 {object} swagger.UnauthenticatedResponse
// @failure 403 {object} swagger.AccessForbiddenResponse
// @Router /admin/users [post]
func (uc UserController) CreateUser(c *gin.Context) {
	var userData models.CreateUserRequestAdmin
	if err := c.ShouldBindJSON(&userData); err != nil {
		fieldErrors := make(map[string]string, 0)
		if !utils.IsGoodPassword(userData.Password) {
			fieldErrors["password"] = "Password must contain at least one alphabet and one number and its length must be 8 characters or more"

		}
		responses.ValidationErrorsJSON(c, err, "", fieldErrors)
		return
	}
	if !utils.IsGoodPassword(userData.Password) {
		fieldErrors := map[string]string{
			"password": "Password must contain at least one alphabet and one number and its length must be 8 characters or more",
		}
		responses.ManualValidationErrorsJSON(c, fieldErrors, "")
		return
	}

	err := uc.userService.CreateUser(userData)
	if err != nil {
		responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "Sorry an error occurred!")
	}
	uc.paginateUserList(c, "User created successfully.")
}

func (uc *UserController) paginateUserList(c *gin.Context, message string) {
	pagination := utils.BuildPagination(c)
	users, count, err := uc.userService.GetAllUsers(pagination)
	if err != nil {
		responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "Sorry an error occurred 😢")
		return
	}
	responses.JSON(c, http.StatusOK, gin.H{
		"users": users,
		"count": count,
	}, message)
}
