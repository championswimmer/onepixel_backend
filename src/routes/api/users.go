package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"onepixel_backend/src/controllers"
	"onepixel_backend/src/dtos"
	"onepixel_backend/src/security"
	"onepixel_backend/src/server/parsers"
	"onepixel_backend/src/server/validators"
)

var usersController *controllers.UsersController

// UsersRoute defines the routes for /api/v1/users
func UsersRoute() func(router fiber.Router) {
	// initialize UsersController
	usersController = controllers.CreateUsersController()

	return func(router fiber.Router) {
		router.Post("/", security.MandatoryAdminApiKeyAuthMiddleware, registerUser)
		router.Post("/login", loginUser)
		router.Get("/:userid", security.MandatoryJwtAuthMiddleware, getUserInfo)
		router.Patch("/:userid", security.MandatoryJwtAuthMiddleware, updateUserInfo)
	}
}

// registerUser
//
//	@Summary		Register new user
//	@Description	Register new user
//	@ID				register-user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		dtos.CreateUserRequest	true	"User"
//	@Success		201		{object}	dtos.UserResponse
//	@Failure		400		{object}	dtos.ErrorResponse "The request body is not valid"
//	@Failure		422		{object}	dtos.ErrorResponse "email and password are required to create user"
//	@Failure		409		{object}	dtos.ErrorResponse "User with this email already exists"
//	@Router			/users [post]
//	@Security		APIKeyAuth
func registerUser(ctx *fiber.Ctx) error {

	u, parseError := parsers.ParseBody[dtos.CreateUserRequest](ctx)
	if parseError != nil {
		return parsers.SendParsingError(ctx, parseError)
	}

	validateErr := validators.ValidateCreateUserRequest(u)
	if validateErr != nil {
		return validators.SendValidationError(ctx, validateErr)
	}

	savedUser, token, err := usersController.Create(u.Email, u.Password)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ctx.Status(fiber.StatusConflict).JSON(dtos.CreateErrorResponse(fiber.StatusConflict, "User with this email already exists"))
		} else {
			return ctx.Status(fiber.StatusInternalServerError).JSON(dtos.CreateErrorResponse(fiber.StatusInternalServerError, "something went wrong"))
		}
	}

	return ctx.Status(fiber.StatusCreated).JSON(dtos.CreateUserResponseFromUser(savedUser, &token))
}

// loginUser
//
// @Summary		Login user
// @Description	Login user
// @ID				login-user
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			user	body		dtos.LoginUserRequest true	"User"
// @Success		200		{object}	dtos.UserResponse
// @Failure		401		{object}	dtos.ErrorResponse "Invalid email or password"
// @Router			/users/login [post]
func loginUser(ctx *fiber.Ctx) error {
	u, parseError := parsers.ParseBody[dtos.LoginUserRequest](ctx)
	if parseError != nil {
		return parsers.SendParsingError(ctx, parseError)
	}

	validateErr := validators.ValidateLoginUserRequest(u)
	if validateErr != nil {
		return validators.SendValidationError(ctx, validateErr)
	}

	user, err := usersController.VerifyEmailAndPassword(u.Email, u.Password)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(dtos.CreateErrorResponse(fiber.StatusUnauthorized, err.Error()))
	}
	token := security.CreateJWTFromUser(user)

	return ctx.Status(fiber.StatusOK).JSON(dtos.CreateUserResponseFromUser(user, &token))
}

// getUserInfo
//
// @Summary		Get user info
// @Description	Get user info
// @ID				get-user-info
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			userid	path	uint	true	"User ID"
// @Router			/users/{userid} [get]
func getUserInfo(ctx *fiber.Ctx) error {
	return ctx.SendString("GetUserInfo")
}

// updateUserInfo
//
// @Summary		Update user info
// @Description	Update user info
// @ID				update-user-info
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			userid	path	uint	true	"User ID"
// @Router			/users/{userid} [patch]
func updateUserInfo(ctx *fiber.Ctx) error {
	return ctx.SendString("UpdateUserInfo")
}
