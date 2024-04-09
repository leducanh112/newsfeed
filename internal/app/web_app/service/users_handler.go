package service

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/leducanh112/newsfeed/internal/pkg/types"
	"github.com/leducanh112/newsfeed/pkg/types/proto/pb/authen_and_post"
)

// CheckUserNamePassword godoc
//
//	@Summary		check user authentication
//	@Description	check user user_name and password
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			request body types.LoginRequest true "login param"
//	@Success		200	{object} types.MessageResponse
//	@Failure		400	{object} types.MessageResponse
//	@Failure		500	{object} types.MessageResponse
//	@Router			/users/login [post]
func (svc *WebService) CheckUserNamePassword(ctx *gin.Context) {
	start := time.Now()
	status := http.StatusOK
	countExporter.WithLabelValues("check_user_login", "total").Inc()
	defer func() {
		latencyExporter.WithLabelValues("check_user_login", strconv.Itoa(status)).Observe(float64(start.UnixMilli()))
	}()

	// parse request
	var jsonRequest types.LoginRequest
	err := ctx.ShouldBindJSON(&jsonRequest)
	if err != nil {
		countExporter.WithLabelValues("check_user_login", "bad_request").Inc()
		status = http.StatusBadRequest
		ctx.JSON(status, &types.MessageResponse{Message: err.Error()})
		return
	}

	// call AAP gRPC service to check auth
	authentication, err := svc.authenticateAndPostClient.CheckUserAuthentication(ctx, &authen_and_post.CheckUserAuthenticationRequest{
		UserName:     jsonRequest.UserName,
		UserPassword: jsonRequest.Password,
	})
	if err != nil {
		countExporter.WithLabelValues("check_user_login", "call_api_failed").Inc()
		status = http.StatusInternalServerError
		ctx.JSON(status, &types.MessageResponse{Message: err.Error()})
		return
	}

	// process result
	if authentication.GetStatus() == authen_and_post.CheckUserAuthenticationResponse_NOT_FOUND {
		countExporter.WithLabelValues("check_user_login", "not_found").Inc()
		ctx.JSON(status, &types.MessageResponse{Message: "not found"})
		return
	} else if authentication.GetStatus() == authen_and_post.CheckUserAuthenticationResponse_WRONG_PASSWORD {
		countExporter.WithLabelValues("check_user_login", "wrong_password").Inc()
		ctx.JSON(status, &types.MessageResponse{Message: "wrong password"})
		return
	}

	countExporter.WithLabelValues("check_user_login", "success").Inc()
	ctx.JSON(status, &types.MessageResponse{Message: "ok"})
	ctx.SetCookie("session_id", fmt.Sprintf("%d", authentication.UserId), 0, "", "", false, false)
}

// CreateUser godoc
//
//	@Summary		create user
//	@Description	create new user using user provided information
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			request body types.CreateUserRequest true "create user param"
//	@Success		200	{object} types.MessageResponse
//	@Failure		400	{object} types.MessageResponse
//	@Failure		500	{object} types.MessageResponse
//	@Router			/users [post]
func (svc *WebService) CreateUser(ctx *gin.Context) {
	start := time.Now()
	status := http.StatusOK
	countExporter.WithLabelValues("create_user", "total").Inc()
	defer func() {
		latencyExporter.WithLabelValues("create_user", strconv.Itoa(status)).Observe(float64(start.UnixMilli()))
		svc.log.Debug("processed create_user API")
	}()

	// parse request
	var jsonRequest types.CreateUserRequest
	err := ctx.ShouldBindJSON(&jsonRequest)
	if err != nil {
		countExporter.WithLabelValues("create_user", "bad_request").Inc()

		ctx.JSON(http.StatusBadRequest, &types.MessageResponse{Message: err.Error()})
		svc.log.Debug("processed create_user API with req=%+v status=%d", zap.String("req", fmt.Sprintf("%+v", jsonRequest)), zap.Int("status", status))
		return
	}

	// parse DOB string to time struct
	dob, err := time.Parse("2006-01-02", jsonRequest.Dob)
	if err != nil {
		countExporter.WithLabelValues("create_user", "bad_request.wrong_dob").Inc()
		ctx.JSON(http.StatusBadRequest, &types.MessageResponse{Message: err.Error()})
		return
	}

	// call AAP grpc service to create user
	createdUser, err := svc.authenticateAndPostClient.CreateUser(ctx, &authen_and_post.UserDetailInfo{
		FirstName:    jsonRequest.FirstName,
		LastName:     jsonRequest.LastName,
		Dob:          timestamppb.New(dob),
		UserName:     jsonRequest.UserName,
		UserPassword: jsonRequest.Password,
		Email:        jsonRequest.Email,
	})
	if err != nil {
		countExporter.WithLabelValues("create_user", "call_api_failed").Inc()
		ctx.JSON(http.StatusInternalServerError, &types.MessageResponse{Message: err.Error()})
		return
	}

	countExporter.WithLabelValues("create_user", "success").Inc()
	ctx.JSON(http.StatusOK, &types.MessageResponse{Message: fmt.Sprintf("Successfully created user with id: %d", createdUser.Info.UserId)})
}

// EditUser godoc
//
//	@Summary		edit user
//	@Description	edit user using user provided information
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			request body types.EditUserRequest true "edit user param"
//	@Success		200	{object} types.MessageResponse
//	@Failure		400	{object} types.MessageResponse
//	@Failure		500	{object} types.MessageResponse
//	@Router			/users [put]
func (svc *WebService) EditUser(ctx *gin.Context) {
	start := time.Now()
	status := http.StatusOK
	countExporter.WithLabelValues("edit_user", "total").Inc()
	defer func() {
		latencyExporter.WithLabelValues("edit_user", strconv.Itoa(status)).Observe(float64(start.UnixMilli()))
	}()

	// parse request
	var jsonRequest types.EditUserRequest
	err := ctx.ShouldBindJSON(&jsonRequest)
	if err != nil {
		countExporter.WithLabelValues("edit_user", "bad_request").Inc()
		ctx.JSON(http.StatusBadRequest, &types.MessageResponse{Message: err.Error()})
		return
	}

	var (
		userId    int64
		firstName *string
		lastName  *string
		dob       *timestamppb.Timestamp
		password  *string
	)

	userId = jsonRequest.UserId
	if userId == 0 {
		countExporter.WithLabelValues("edit_user", "bad_request.empty_user_id").Inc()
		ctx.JSON(http.StatusBadRequest, &types.MessageResponse{Message: "User id is required"})
		return
	}

	if jsonRequest.FirstName != "" {
		*firstName = jsonRequest.FirstName
	}

	if jsonRequest.LastName != "" {
		*lastName = jsonRequest.LastName
	}

	if jsonRequest.Dob != "" {
		parsedDob, err := time.Parse("2006-01-02", jsonRequest.Dob)
		if err != nil {
			countExporter.WithLabelValues("edit_user", "bad_request.invalid_dob").Inc()
			ctx.JSON(http.StatusBadRequest, &types.MessageResponse{Message: err.Error()})
			return
		}
		dob = timestamppb.New(parsedDob)
	}

	if jsonRequest.Password != "" {
		password = &jsonRequest.Password
	}

	// call AAP grpc service to edit user
	resp, err := svc.authenticateAndPostClient.EditUser(ctx, &authen_and_post.EditUserRequest{
		UserId:       userId,
		FirstName:    firstName,
		LastName:     lastName,
		Dob:          dob,
		UserPassword: password,
	})
	if err != nil {
		countExporter.WithLabelValues("edit_user", "call_api_failed").Inc()
		ctx.JSON(http.StatusInternalServerError, &types.MessageResponse{Message: err.Error()})
		return
	}

	countExporter.WithLabelValues("edit_user", "success").Inc()
	ctx.JSON(http.StatusOK, &types.MessageResponse{Message: "Successfully edited user with id: " + fmt.Sprintf("%d", resp.UserId)})
}
