package entity

import "github.com/golang-jwt/jwt/v5"

type JsonTokenType string

const (
	JsonWebTokenRefreshToken JsonTokenType = "refresh"
	JsonWebTokenAccessToken  JsonTokenType = "access"
)

type JsonWebToken struct {
	ID      string        `json:"id" gorm:"type:varchar(255);primarykey"`
	UserId  string        `json:"user_id" gorm:"type:varchar(255);index"`
	Type    JsonTokenType `gorm:"type:varchar(255)"`
	Revoked bool          `json:"revoked" gorm:"default:false"`
	Ref     string        `json:"ref" gorm:"type:varchar(255)"`
	Iat     int64         `json:"iat" gorm:"default:0"`
	Exp     int64         `json:"exp" gorm:"default:0"`
}

type JwtTokenResponse struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

type JwtToken struct {
	Sub   string `json:"sub"`
	Iat   int64  `json:"iat"`
	Exp   int64  `json:"exp"`
	Iss   string `json:"iss"`
	Aud   string `json:"aud"`
	Email string `json:"email"`
}

func (j *JwtToken) ToMapClaims() jwt.MapClaims {
	return jwt.MapClaims{
		"sub":   j.Sub,
		"iat":   j.Iat,
		"exp":   j.Exp,
		"iss":   j.Iss,
		"aud":   j.Aud,
		"email": j.Email,
	}
}
