package grpc

import (
	authGRPC "chickchirick-auth/internal/gen/auth"
	"context"
	"fmt"

	"chickchirick-auth/internal/service"
	"chickchirick-auth/pkg/chirick_config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type AuthGRPCController struct {
	authGRPC.UnimplementedAuthServiceServer
}

func (s *AuthGRPCController) ValidateToken(ctx context.Context, req *authGRPC.ValidateRequest) (*authGRPC.ValidateResponse, error) {
	tokenStr := req.GetToken()
	if tokenStr == "" {
		return &authGRPC.ValidateResponse{Valid: false}, nil
	}

	claims := &service.Claims{}
	t, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(viper.GetString(chirick_config.SecretKey)), nil
	})

	if err != nil || !t.Valid {
		return &authGRPC.ValidateResponse{
			Valid: false,
		}, nil
	}

	return &authGRPC.ValidateResponse{
		Valid:    true,
		UserUuid: claims.UserUuid,
	}, nil
}
