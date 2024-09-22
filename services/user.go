package services

import (
	"time"

	"github.com/google/uuid"
	"github.com/kiplikipli/technical-test-fm-tahap-2/database"
	"github.com/kiplikipli/technical-test-fm-tahap-2/entity"
	"golang.org/x/crypto/bcrypt"
)

type User entity.User

func GetUserByID(id uuid.UUID) (*User, error) {
	db := database.DB
	found := User{}
	query := User{
		ID: id,
	}

	err := db.First(&found, &query).Error
	return &found, err
}

func GetUserByPhoneNumber(phoneNumber string) (*User, error) {
	db := database.DB
	found := User{}
	query := User{
		PhoneNumber: phoneNumber,
	}

	err := db.First(&found, &query).Error
	return &found, err
}

func CreateUser(user *User) (*User, error) {
	db := database.DB
	createRequest := &User{
		ID:          uuid.New(),
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		Pin:         hashAndSalt([]byte(user.Pin)),
		Address:     user.Address,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := db.Create(createRequest).Error
	if err != nil {
		return nil, err
	}

	newUser := User{
		ID:          createRequest.ID,
		FirstName:   createRequest.FirstName,
		LastName:    createRequest.LastName,
		PhoneNumber: createRequest.PhoneNumber,
		Address:     createRequest.Address,
		CreatedAt:   createRequest.CreatedAt,
	}

	return &newUser, nil
}

func UpdateUser(userRequest *User) (*User, error) {
	db := database.DB
	var user User
	query := &User{
		ID: userRequest.ID,
	}

	err := db.First(&user, &query).Error
	if err != nil {
		return nil, err
	}

	if userRequest.FirstName != "" {
		user.FirstName = userRequest.FirstName
	}
	if userRequest.LastName != "" {
		user.LastName = userRequest.LastName
	}
	if userRequest.Address != "" {
		user.Address = userRequest.Address
	}

	err = db.Save(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func CompareHash(hashedPin string, plainPin string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPin), []byte(plainPin))
}

func hashAndSalt(pwd []byte) string {
	hash, _ := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	return string(hash)
}
