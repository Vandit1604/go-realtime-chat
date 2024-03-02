package user

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/vandit1604/go-realtime-chat/util"
)

const (
	secretKey = "secret"
)

type service struct {
	Repository
	timeout time.Duration
}

func NewService(repository Repository) Service {
	return &service{
		repository,
		time.Duration(2) * time.Second,
	}
}
func (s *service) CreateUser(c context.Context, req *CreateUserReq) (*CreateUserRes, error) {
	// we're using this so even if the db operation completes before the timeout the context will be automatically canceled
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		log.Fatal(err)
	}

	u := &User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	userCreatedInDb, err := s.Repository.CreateUser(ctx, u)
	if err != nil {
		log.Fatal(err)
	}

	res := &CreateUserRes{
		ID:       strconv.Itoa(int(userCreatedInDb.ID)),
		Username: userCreatedInDb.Username,
		Email:    userCreatedInDb.Email,
	}

	return res, nil
}

type MyJwtClaims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (s *service) Login(ctx context.Context, req *LoginUserReq) (*LoginUserRes, error) {
	// this context will be passed to the Repository layer for DB operations
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	user, err := s.Repository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return &LoginUserRes{}, err
	}

	err = util.CheckPassword(req.Password, user.Password)
	if err != nil {
		return &LoginUserRes{}, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyJwtClaims{
		ID:       strconv.Itoa(int(user.ID)),
		Username: user.Username,
		// learn more about this
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    strconv.Itoa(int(user.ID)),
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(24 * time.Hour)},
		},
	})

	ss, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return &LoginUserRes{}, err
	}

	return &LoginUserRes{
		accessToken: ss,
		ID:          strconv.Itoa(int(user.ID)),
		Username:    user.Username,
	}, nil
}
