package jwt_auth

import (
	"fmt"
	"iu7-2022-sd-labs/buisness/entities"
	"iu7-2022-sd-labs/buisness/ports/repositories"
	"iu7-2022-sd-labs/server/ports"
	"time"

	"github.com/golang-jwt/jwt"
)

const organizerSubject = "organizer"
const consumerSubject = "consumer"

type JWTAuth struct {
	signingKey    []byte
	organizerRepo repositories.OrganizerRepository
	consumerRepo  repositories.ConsumerRepository
}

func NewJWTAuth(
	signingKey string,
	organizerRepo repositories.OrganizerRepository,
	consumerRepo repositories.ConsumerRepository,
) JWTAuth {
	return JWTAuth{[]byte(signingKey), organizerRepo, consumerRepo}
}

func (a *JWTAuth) newToken(id string, subject string) (string, error) {
	claims := jwt.StandardClaims{
		Id:       id,
		IssuedAt: time.Now().Unix(),
		Subject:  subject,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(a.signingKey)
	return tokenString, Wrap(err, "token signed string")
}

func (a *JWTAuth) NewOrganizerToken(organizer *entities.Organizer) (string, error) {
	return a.newToken(organizer.ID(), organizerSubject)
}

func (a *JWTAuth) NewConsumerToken(consumer *entities.Consumer) (string, error) {
	return a.newToken(consumer.ID(), consumerSubject)
}

func (a *JWTAuth) parseTokenClaims(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrWrongSigningMethod
		}

		return a.signingKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed to convert to map claims")
	}

	return claims, nil
}

func (a *JWTAuth) parseToken(tokenString string, targetSubject string) (string, error) {
	claims, err := a.parseTokenClaims(tokenString)
	if err != nil {
		return "", err
	}

	subject, ok := claims["sub"].(string)
	if !ok {
		return "", fmt.Errorf("failed to convert sub")
	}

	if subject != targetSubject {
		return "", ports.ErrWrongSubject
	}

	id, ok := claims["jti"].(string)
	if !ok {
		return "", fmt.Errorf("failed to convert jti")
	}

	return id, nil
}

func (a *JWTAuth) ParseOrganizerToken(token string) (entities.Organizer, error) {
	id, err := a.parseToken(token, organizerSubject)
	if err != nil {
		return entities.Organizer{}, Wrap(err, "parse token")
	}

	organizer, err := a.organizerRepo.Get(id)
	return organizer, Wrap(err, "organizer repo get")
}

func (a *JWTAuth) ParseConsumerToken(token string) (entities.Consumer, error) {
	id, err := a.parseToken(token, consumerSubject)
	if err != nil {
		return entities.Consumer{}, Wrap(err, "parse token")
	}

	consumer, err := a.consumerRepo.Get(id)
	return consumer, Wrap(err, "consumer repo get")
}
