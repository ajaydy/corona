package api

import (
	"context"
	"corona/helpers"
	"corona/models"
	"corona/util"
	"database/sql"
	"encoding/base64"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type (
	UserModule struct {
		db     *sql.DB
		logger *helpers.Logger
		name   string
	}

	UserDetailParam struct {
		Id uuid.UUID `json:"id"`
	}

	UserLoginParam struct {
		Email    string `json:"email" validate:"email,required"`
		Password string `json:"password" validate:"required"`
	}

	UserRegisterParam struct {
		SubscriptionId  uuid.UUID `json:"subscription_id" validate:"required"`
		Name            string    `json:"name" validate:"max=20,min=4,required"`
		Email           string    `json:"email" validate:"email,required"`
		Password        string    `json:"password" validate:"required"`
		ConfirmPassword string    `json:"confirm_password" validate:"required"`
	}

	UserWithToken struct {
		User  models.UserResponse `json:"user"`
		Token string              `json:"token"`
	}
)

func NewUserModule(db *sql.DB, logger *helpers.Logger) *UserModule {
	return &UserModule{
		db:     db,
		logger: logger,
		name:   "module/user",
	}
}

func (u UserModule) Register(ctx context.Context, param UserRegisterParam) (interface{}, *helpers.Error) {

	if param.Password != param.ConfirmPassword {
		return nil, helpers.ErrorWrap(errors.New("Password Does Not Match !"), u.name,
			"Register/Password", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	password, err := bcrypt.GenerateFromPassword([]byte(param.Password), 12)

	if err != nil {
		return nil, helpers.ErrorWrap(err, u.name, "Register/bcrypt.GenerateFromPassword",
			helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	user := models.UserModel{
		SubscriptionId: param.SubscriptionId,
		Name:           param.Name,
		Email:          param.Email,
		Password:       string(password),
		CreatedBy:      uuid.NewV4(),
	}

	err = user.Insert(ctx, u.db)

	if err != nil {
		return nil, helpers.ErrorWrap(err, u.name, "Register/user.Insert", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	userResponse, err := user.Response(ctx, u.db, u.logger)

	if err != nil {
		return nil, helpers.ErrorWrap(err, u.name, "Register/user.Response", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	tokenKey := util.GenerateToken()

	expiredAt := time.Now().AddDate(0, 1, 0)

	encToken := base64.StdEncoding.EncodeToString([]byte(tokenKey))

	token := models.TokenModel{
		UserId:    user.Id,
		TokenKey:  encToken,
		ExpiredAt: expiredAt,
		CreatedBy: uuid.NewV4(),
	}

	err = token.Insert(ctx, u.db)

	if err != nil {
		return nil, helpers.ErrorWrap(err, u.name, "Register/token.Insert", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	rateLimit := models.RateLimitModel{
		UserId:       user.Id,
		TotalRequest: 0,
		CreatedBy:    uuid.NewV4(),
	}

	err = rateLimit.Insert(ctx, u.db)

	userWithToken := UserWithToken{
		User:  userResponse,
		Token: tokenKey,
	}

	return userWithToken, nil
}

func (u UserModule) Login(ctx context.Context, param UserLoginParam) (interface{}, *helpers.Error) {

	user, err := models.GetOneUserByEmail(ctx, u.db, param.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, helpers.ErrorWrap(err, u.name, "Login/Email", helpers.IncorrectEmailMessage,
				http.StatusInternalServerError)
		}
		return nil, helpers.ErrorWrap(err, u.name, "Login/GetOneUserByEmail", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(param.Password))

	if err != nil {
		if err != nil {
			return nil, helpers.ErrorWrap(errors.New("Invalid Password"), u.name,
				"Login/CompareHashAndPassword",
				"Incorrect Password",
				http.StatusInternalServerError)
		}
	}

	userResponse, err := user.Response(ctx, u.db, u.logger)

	if err != nil {
		return nil, helpers.ErrorWrap(err, u.name, "Login/user.Response", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return userResponse, nil

}

func (u UserModule) List(ctx context.Context, filter helpers.Filter) (interface{}, *helpers.Error) {

	users, err := models.GetAllUser(ctx, u.db, filter)

	if err != nil {
		return nil, helpers.ErrorWrap(err, u.name, "List/GetAllUser",
			helpers.InternalServerError, http.StatusInternalServerError)
	}

	var userResponses []models.UserResponse
	for _, user := range users {
		response, err := user.Response(ctx, u.db, u.logger)
		if err != nil {
			return nil, helpers.ErrorWrap(err, u.name, "List/user.Response",
				helpers.InternalServerError, http.StatusInternalServerError)
		}
		userResponses = append(userResponses, response)
	}

	return userResponses, nil
}

func (u UserModule) Detail(ctx context.Context, param UserDetailParam) (interface{}, *helpers.Error) {

	user, err := models.GetOneUser(ctx, u.db, param.Id)

	if err != nil {
		return nil, helpers.ErrorWrap(err, u.name, "Detail/GetOneUser",
			helpers.InternalServerError, http.StatusInternalServerError)
	}

	response, err := user.Response(ctx, u.db, u.logger)

	if err != nil {
		return nil, helpers.ErrorWrap(err, u.name, "Detail/user.Response",
			helpers.InternalServerError, http.StatusInternalServerError)
	}

	return response, nil

}
