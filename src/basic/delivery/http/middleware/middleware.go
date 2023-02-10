package middleware

import (
	"context"
	"douyin-service/domain"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/jwt"
	"github.com/spf13/viper"
	"strconv"
)

type DouyinMiddleware struct {
	JWTMid *jwt.HertzJWTMiddleware
}

//func validateToken(token string) (*domain.TokenClaims, error) {
//	key := viper.GetString("jwt_key")
//	var c domain.TokenClaims
//	_, err := jwt.ParseWithClaims(token, &c, func(token *jwt.Token) (interface{}, error) {
//		return []byte(key), nil
//	})
//	fmt.Println(c)
//	if err != nil {
//		return nil, err
//	}
//	return &c, nil
//}

func New() *DouyinMiddleware {
	key := viper.GetString("jwt_key")
	mid, _ := jwt.New(&jwt.HertzJWTMiddleware{
		Realm:       "douyin",
		Key:         []byte(key),
		IdentityKey: "uid",
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			return jwt.MapClaims{"uid": strconv.FormatInt(data.(*domain.TokenClaims).Uid, 10)}
		},
		IdentityHandler: func(ctx context.Context, c *app.RequestContext) interface{} {
			claims := jwt.ExtractClaims(ctx, c)
			i, _ := strconv.ParseInt(claims["uid"].(string), 10, 64)
			return i
		},
		TokenLookup: "form:token, query: token",
		Unauthorized: func(ctx context.Context, c *app.RequestContext, code int, message string) {
			c.JSON(code, map[string]interface{}{
				"code":    code,
				"message": message,
			})
		},
	})
	return &DouyinMiddleware{mid}
}

func (m *DouyinMiddleware) TokenAuth() app.HandlerFunc {
	return m.JWTMid.MiddlewareFunc()
}

//func (m *DouyinMiddleware) TokenAuthPublishAction() app.HandlerFunc {
//	return keyauth.New(
//		keyauth.WithKeyLookUp("form:token", ""),
//		keyauth.WithValidator(func(ctx context.Context, c *app.RequestContext, s string) (bool, error) {
//			//log.Println("request uri: ", c.URI())
//			claim, err := validateToken(s)
//			if err != nil {
//				fmt.Println("token解析失败")
//				return false, nil
//			}
//			fmt.Println("claim:", claim)
//			if claim.ExpiresAt < time.Now().Unix() {
//				return false, nil
//			}
//			c.Set("uid", claim.Id)
//			//c.Next(ctx)
//			return true, nil
//		}),
//		keyauth.WithErrorHandler(func(ctx context.Context, c *app.RequestContext, err error) {
//			c.AbortWithStatusJSON(http.StatusOK, domain.Response{
//				StatusCode: 2,
//				StatusMsg:  "认证失败",
//			})
//		}),
//	)
//}
