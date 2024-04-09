package authen_and_post_svc

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/leducanh112/newsfeed/internal/pkg/types"
	"github.com/leducanh112/newsfeed/pkg/types/proto/pb/authen_and_post"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func (a *AuthenticateAndPostServer) CheckUserAuthentication(ctx context.Context, info *authen_and_post.CheckUserAuthenticationRequest) (*authen_and_post.CheckUserAuthenticationResponse, error) {
	// query from db
	var user types.User
	result := a.db.Where(&types.User{UserName: info.GetUserName()}).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &authen_and_post.CheckUserAuthenticationResponse{
			Status: authen_and_post.CheckUserAuthenticationResponse_NOT_FOUND,
		}, nil
	} else if result.Error != nil {
		return nil, fmt.Errorf("failed to query user %+v: %s", info.GetUserName(), result.Error)
	}

	// compare hashed password & hash(password+salt)
	passwordWithSalt := []byte(info.UserPassword + user.Salt)
	err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), passwordWithSalt)
	if err != nil {
		return &authen_and_post.CheckUserAuthenticationResponse{
			Status: authen_and_post.CheckUserAuthenticationResponse_WRONG_PASSWORD,
		}, nil
	}

	return &authen_and_post.CheckUserAuthenticationResponse{
		Status: authen_and_post.CheckUserAuthenticationResponse_OK,
	}, nil
}

func (a *AuthenticateAndPostServer) CreateUser(ctx context.Context, info *authen_and_post.UserDetailInfo) (*authen_and_post.UserResult, error) {
	// TODO: validate if user is existed in db

	salt := generateAlphabetSalt(16)
	hashedPassword, err := hashPassword(info.GetUserPassword(), salt)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %s", err)
	}

	newUser := types.User{
		HashedPassword: hashedPassword,
		Salt:           string(salt),
		FirstName:      info.GetFirstName(),
		LastName:       info.GetLastName(),
		DateOfBirth:    info.Dob.AsTime(),
		Email:          info.GetEmail(),
		UserName:       info.GetUserName(),
	}

	res := a.db.Create(&newUser)
	if res.Error != nil {
		return nil, fmt.Errorf("failed to create user: %s", res.Error)
	}

	return &authen_and_post.UserResult{
		Status: authen_and_post.UserStatus_OK,
		Info: &authen_and_post.UserDetailInfo{
			UserId:   int64(newUser.ID),
			UserName: newUser.UserName,
		},
	}, nil
}

// EditUser edit user info by looking up user id in mysql database and update it
func (a *AuthenticateAndPostServer) EditUser(ctx context.Context, info *authen_and_post.EditUserRequest) (*authen_and_post.EditUserResponse, error) {
	var user types.User
	a.db.Where(&types.User{ID: uint(info.UserId)}).First(&user)
	if user.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}
	if info.FirstName != nil {
		user.FirstName = info.GetFirstName()
	}
	if info.LastName != nil {
		user.LastName = info.GetLastName()
	}
	if info.UserPassword != nil {
		salt := generateAlphabetSalt(16)
		hashedPassword, err := hashPassword(info.GetUserPassword(), salt)
		if err != nil {
			return nil, err
		}
		user.HashedPassword = hashedPassword
		user.Salt = string(salt)
	}
	if info.Dob != nil {
		user.DateOfBirth = info.Dob.AsTime()
	}
	a.db.Save(&user)

	return &authen_and_post.EditUserResponse{
		UserId: int64(user.ID),
	}, nil
}

func generateAlphabetSalt(length int) []byte {
	rand.Seed(time.Now().UnixNano())

	salt := make([]byte, length)
	for i := 0; i < length; i++ {
		salt[i] = alphabet[rand.Intn(len(alphabet))]
	}

	return salt
}

func hashPassword(password string, salt []byte) (string, error) {
	// Append the salt to the password
	passwordWithSalt := []byte(password + string(salt))

	// Generate the bcrypt hash
	hash, err := bcrypt.GenerateFromPassword(passwordWithSalt, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (a *AuthenticateAndPostServer) GetUserFollower(ctx context.Context, info *authen_and_post.UserInfo) (*authen_and_post.UserFollower, error) {
	// TODO implement me
	panic("implement me")
}

func (a *AuthenticateAndPostServer) GetPostDetail(ctx context.Context, request *authen_and_post.GetPostRequest) (*authen_and_post.Post, error) {
	// TODO implement me
	panic("implement me")
}
