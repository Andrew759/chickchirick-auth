package service

import (
	"chickchirick-auth/cmd/service"
	"chickchirick-auth/internal/model/dto"
	"chickchirick-auth/pkg/chirick_config"
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

const ISS = "auth_service"

var (
	TokenInvalidErr = errors.New("token is invalid")
	TokenExpiredErr = errors.New("token has expired or revoked")
)

type Claims struct {
	UserUuid string `json:"user_uuid"`
	jwt.RegisteredClaims
}

// CreateTokens создает пару Access и Refresh, сохраняя JTI в Redis
func CreateTokens(ctx context.Context, redisDec service.RedisDecorator, userUuid string) (dto.AccessToken, dto.RefreshToken, error) {
	atStr, err := createAccessToken(userUuid)
	if err != nil {
		return dto.AccessToken{}, dto.RefreshToken{}, err
	}

	rtStr, jti, err := createRefreshToken(userUuid)
	if err != nil {
		return dto.AccessToken{}, dto.RefreshToken{}, err
	}

	//Сохраняет JTI в Redis. Ключ — JTI, значение — UserUUID.
	err = redisDec.Client.Set(ctx, jti, userUuid, viper.GetDuration(chirick_config.RefreshTokenLT)).Err()
	if err != nil {
		return dto.AccessToken{}, dto.RefreshToken{}, err
	}

	return dto.AccessToken{
			UserUuid: userUuid,
			Token:    atStr,
			Lt:       viper.GetDuration(chirick_config.AccessTokenLT),
		}, dto.RefreshToken{
			UserUuid: userUuid,
			Token:    rtStr,
			Lt:       viper.GetDuration(chirick_config.RefreshTokenLT),
		}, nil
}

// RefreshToken реализует ротацию: проверяет старый RT, удаляет его и выдает новые
func RefreshToken(ctx context.Context, redisDec service.RedisDecorator, refreshTokenStr string) (dto.AccessToken, dto.RefreshToken, error) {
	//Парсинг пришедшего токена
	token, err := jwt.ParseWithClaims(refreshTokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString(chirick_config.SecretKey)), nil
	})

	if err != nil || !token.Valid {
		return dto.AccessToken{}, dto.RefreshToken{}, TokenInvalidErr
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return dto.AccessToken{}, dto.RefreshToken{}, TokenInvalidErr
	}

	//Проверка наличия JTI в Redis и его удаление - ротация токена
	userUuid, err := redisDec.Client.GetDel(ctx, claims.ID).Result()
	if err != nil {
		//Если ключа нет — токен либо протух, либо уже был использован
		return dto.AccessToken{}, dto.RefreshToken{}, TokenExpiredErr
	}

	//Генерация новой пары токенов
	return CreateTokens(ctx, redisDec, userUuid)
}

func createAccessToken(userUuid string) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS512, Claims{
		UserUuid: userUuid,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    ISS,
			Subject:   userUuid,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(viper.GetDuration(chirick_config.AccessTokenLT))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}).SignedString([]byte(viper.GetString(chirick_config.SecretKey)))
}

// createRefreshToken возвращает сам токен и его внутренний ID (JTI)
func createRefreshToken(userUuid string) (string, string, error) {
	jti := uuid.New().String()
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserUuid: userUuid,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			Issuer:    ISS,
			Subject:   userUuid,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(viper.GetDuration(chirick_config.RefreshTokenLT))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}).SignedString([]byte(viper.GetString(chirick_config.SecretKey)))

	return token, jti, err
}

// Logout удаляет токен из Redis
func Logout(ctx context.Context, redisDec service.RedisDecorator, refreshTokenStr string) error {
	token, _ := jwt.ParseWithClaims(refreshTokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString(chirick_config.SecretKey)), nil
	})

	if claims, ok := token.Claims.(*Claims); ok {
		return redisDec.Client.Del(ctx, claims.ID).Err()
	}
	return nil
}
