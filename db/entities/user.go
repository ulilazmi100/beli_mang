package entities

import (
	"errors"
	"net/mail"
	"regexp"
	"strconv"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type User struct {
	Id        string    `db:"id" json:"id"`
	Username  string    `db:"username" json:"username"`
	Email     string    `db:"email" json:"email"`
	Role      string    `db:"role" json:"role"`
	Password  string    `db:"password" json:"password,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

type RegistrationPayload struct {
	Username string `json:"username" form:"username"`
	Email    string `json:"email,omitempty" form:"email"`
	Password string `json:"password" form:"password"`
}

type Credential struct {
	Username string `json:"nip"`
	Password string `json:"password"`
}

type JWTPayload struct {
	Id       string
	Username string
	Role     string
}

type JWTClaims struct {
	Id       string
	Username string
	Role     string
	jwt.RegisteredClaims
}

func NewUser(nip, name, password string) *User {
	u := &User{
		Username: nip,
		Email:    name,
		Password: password,
	}

	return u
}

func (u *Credential) Validate() error {
	err := validation.ValidateStruct(u,
		validation.Field(&u.Username,
			validation.Required.Error("username is required"),
		),
		validation.Field(&u.Password,
			validation.Required.Error("password is required"),
			validation.Length(5, 33).Error("password must be between 5 and 33 characters"),
		),
	)

	return err
}

func (u *RegistrationPayload) Validate() error {
	err := validation.ValidateStruct(u,
		validation.Field(&u.Username,
			validation.Required.Error("username is required"),
			validation.Length(5, 30).Error("username must be between 5 and 30 characters"),
		),
		validation.Field(&u.Password,
			validation.Required.Error("password is required"),
			validation.Length(5, 30).Error("password must be between 5 and 30 characters"),
		),
		validation.Field(&u.Email,
			validation.Required.Error("email is required"),
			validation.By(ValidateEmail),
		),
	)

	return err
}

func ValidateEmail(value any) error {
	email, ok := value.(string)
	if !ok {
		return errors.New("parse error")
	}
	_, err := mail.ParseAddress(email)
	return err
}

func ValidateImageURL(value any) error {
	url, ok := value.(string)
	if !ok {
		return errors.New("parse error")
	}
	pattern := `^(?:http(s)?:\/\/)?[\w.-]+(?:\.[\w\.-]+)+[\w\-\._~:/?#[\]@!\$&'\(\)\*\+,;=.]+(?:.jpg|.jpeg|.png)+$`
	rgx := regexp.MustCompile(pattern)
	if !rgx.MatchString(url) {
		return errors.New("invalid image URL format")
	}

	return nil
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

// Function to convert int64 to string
func Int64ToString(num int64) string {
	return strconv.FormatInt(num, 10)
}

// Function to convert string to int64
func StringToInt64(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}
