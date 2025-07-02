package repository

import (
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"hunter-backend/di/config"
	"hunter-backend/entity"
	"hunter-backend/util"
	"time"
)

type JsonWebTokenRepository interface {
	createJsonWebToken(token *entity.JwtToken, tokenType entity.JsonTokenType, user *entity.Users, ref string) (*entity.JsonWebToken, error)
	GetTokenById(jwtId string) (*entity.JsonWebToken, error)
	GenerateToken(userEnt *entity.Users) (*entity.JwtTokenResponse, error)
	GenerateAccessToken(user *entity.Users, ref string) (string, error)

	generateRefreshToken(userEnt *entity.Users) (string, *entity.JsonWebToken, error)
}

type jsonWebTokenRepository struct {
	db                  *gorm.DB
	config              config.AppConfig
	encryptorRepository EncryptorRepository
}

func (j jsonWebTokenRepository) GetTokenById(jwtId string) (*entity.JsonWebToken, error) {
	var ent entity.JsonWebToken
	result := j.db.First(&ent, "id = ?", jwtId)
	if result.Error != nil {
		return &entity.JsonWebToken{}, result.Error
	}
	return &ent, nil
}

func (j jsonWebTokenRepository) createJsonWebToken(token *entity.JwtToken, tokenType entity.JsonTokenType, user *entity.Users, ref string) (*entity.JsonWebToken, error) {
	ent := &entity.JsonWebToken{
		ID:     token.Sub,
		UserId: user.ID,
		Type:   tokenType,
		Iat:    token.Iat,
		Exp:    token.Exp,
		Ref:    ref,
	}

	result := j.db.Create(ent)
	if result.Error != nil {
		return nil, result.Error
	}

	return ent, nil
}

func (j jsonWebTokenRepository) generateRefreshToken(user *entity.Users) (string, *entity.JsonWebToken, error) {
	id, err := j.encryptorRepository.GeneratePassphrase(20)
	if err != nil {
		return "", nil, err
	}

	jwtEnt := &entity.JwtToken{
		Sub:   id,
		Iat:   time.Now().Unix(),
		Exp:   time.Now().Add(time.Minute * 60 * 24 * 7).Unix(),
		Iss:   "Hunter",
		Aud:   "Hunter",
		Email: j.encryptorRepository.Decrypt(user.Email),
	}

	result, err := j.createJsonWebToken(jwtEnt, entity.JsonWebTokenRefreshToken, user, "")
	if err != nil {
		return "", nil, err
	}

	privateKey, err := util.EnsureRSAKeyPair()
	if err != nil {
		return "", nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwtEnt.ToMapClaims())
	token.Header["kid"] = "hunter app"
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", nil, err
	}

	return signedToken, result, nil
}

func (j jsonWebTokenRepository) GenerateAccessToken(user *entity.Users, ref string) (string, error) {
	id, err := j.encryptorRepository.GeneratePassphrase(20)
	if err != nil {
		return "", err
	}

	jwtEnt := &entity.JwtToken{
		Sub:   id,
		Iat:   time.Now().Unix(),
		Exp:   time.Now().Add(time.Minute * 15).Unix(),
		Iss:   "Hunter",
		Aud:   "Hunter",
		Email: j.encryptorRepository.Decrypt(user.Email),
	}

	_, err = j.createJsonWebToken(jwtEnt, entity.JsonWebTokenAccessToken, user, ref)
	if err != nil {
		return "", err
	}

	privateKey, err := util.EnsureRSAKeyPair()
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwtEnt.ToMapClaims())
	token.Header["kid"] = "hunter app"
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (j jsonWebTokenRepository) GenerateToken(userEnt *entity.Users) (*entity.JwtTokenResponse, error) {
	refreshToken, refreshTokenEnt, err := j.generateRefreshToken(userEnt)
	if err != nil {
		return nil, err
	}

	accessToken, err := j.GenerateAccessToken(userEnt, refreshTokenEnt.ID)
	if err != nil {
		return nil, err
	}

	return &entity.JwtTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func ProvideJsonWebTokenRepository(db *gorm.DB, config config.AppConfig) JsonWebTokenRepository {
	encryptorRepository := ProvideEncryptorRepository(db, config)
	return &jsonWebTokenRepository{
		db:                  db,
		config:              config,
		encryptorRepository: encryptorRepository,
	}
}
