package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mr-emerald-wolf/21BCE0665_Backend/database"
	"github.com/mr-emerald-wolf/21BCE0665_Backend/internal/db"
	"github.com/mr-emerald-wolf/21BCE0665_Backend/internal/models"
	"github.com/mr-emerald-wolf/21BCE0665_Backend/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(newUser models.CreateUserRequest) error {

	_, err := database.DB.GetUserByEmail(context.Background(), newUser.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	} else if err == nil {
		return fmt.Errorf("user already exists: %s", newUser.Email)
	}

	// Hash Password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(newUser.Password), 10)

	user := db.CreateUserParams{
		Uuid: pgtype.UUID{
			Bytes: uuid.New(),
			Valid: true,
		},
		Email:        newUser.Email,
		PasswordHash: string(hashedPassword),
	}

	// Create New User
	_, err = database.DB.CreateUser(context.Background(), user)
	if err != nil {
		return err
	}
	return nil
}

func LoginUser(loginRequest models.LoginUserRequest) (string, error) {

	user, err := database.DB.GetUserByEmail(context.Background(), loginRequest.Email)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("user does not exist: %s", loginRequest.Email)
	} else if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginRequest.Password))
	if err != nil {
		return "", fmt.Errorf("password does not match")
	}

	// Generate New Access Token
	access_token, err := utils.CreateToken(user.Email, utils.ACCESS_TOKEN)

	if err != nil {
		return "", err
	}

	return access_token, nil
}

func FindUserByEmail(email string) (db.User, error) {
	user, err := database.DB.GetUserByEmail(context.Background(), email)
	if err != nil {
		return db.User{}, err
	}

	return user, nil
}
