package middleware

import (
	"context"
	"douyin-service/domain"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/dgrijalva/jwt-go"
	"github.com/hertz-contrib/keyauth"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

type DouyinMiddleware struct {
}

func validateToken(token string) (*domain.TokenClaims, error) {
	key := viper.GetString("jwt_key")
	var c domain.TokenClaims
	_, err := jwt.ParseWithClaims(token, &c, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	fmt.Println(c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (m *DouyinMiddleware) TokenAuth() app.HandlerFunc {
	return keyauth.New(
		keyauth.WithKeyLookUp("query:token", "Bearer"),
		keyauth.WithValidator(func(ctx context.Context, c *app.RequestContext, s string) (bool, error) {
			//log.Println("request uri: ", c.URI())
			claim, err := validateToken(s)
			if err != nil {
				fmt.Println("token解析失败")
				return false, nil
			}
			fmt.Println("claim:", claim)
			if claim.ExpiresAt < time.Now().Unix() {
				return false, nil
			}
			c.Set("uid", claim.Id)
			//c.Next(ctx)
			return true, nil
		}),
		keyauth.WithErrorHandler(func(ctx context.Context, c *app.RequestContext, err error) {
			c.AbortWithStatusJSON(http.StatusOK, domain.Response{
				StatusCode: 2,
				StatusMsg:  "认证失败",
			})
		}),
	)
}
