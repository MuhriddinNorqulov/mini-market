package services

import (
	"errors"
	"fmt"
	"mini-market/src/core/application/response"
	"mini-market/src/core/domain/entity"
	"mini-market/src/core/domain/entity/enum"
	"mini-market/src/core/domain/ports/config"
	"mini-market/src/core/domain/ports/security"
	"mini-market/src/core/utils"
	"strconv"
	"time"
)

type UserAuthTokenService struct {
	provider security.JwtTokenProvider
	cfg      config.ConfigProvider
}

// @inject
func NewUserAuthTokenService(provider security.JwtTokenProvider, cfg config.ConfigProvider) *UserAuthTokenService {
	return &UserAuthTokenService{provider: provider, cfg: cfg}
}

func (this *UserAuthTokenService) GenerateToken(user *entity.UserEntity, sessionID string) (access, refresh string, err error) {
	access, err = this.GenerateAccessToken(user, sessionID)
	if err != nil {
		return "", "", err
	}
	refresh, err = this.GenerateRefreshToken(user, sessionID)
	return access, refresh, err
}

func (this *UserAuthTokenService) GenerateAccessToken(user *entity.UserEntity, sessionID string) (string, error) {
	return this.generate(user, sessionID, enum.TokenTypeAccess, this.cfg.GetAccessTokenExpireMinutes())
}

func (this *UserAuthTokenService) GenerateRefreshToken(user *entity.UserEntity, sessionID string) (string, error) {
	return this.generate(&entity.UserEntity{ID: user.ID}, sessionID, enum.TokenTypeRefresh, this.cfg.GetRefreshTokenExpireMinutes())
}

func (this *UserAuthTokenService) generate(user *entity.UserEntity, sessionID string, tokenType enum.TokenType, expireMinutes int64) (string, error) {
	payload, err := utils.JsonMarshal(_TokenPayload{User: user, Type: tokenType, SessionID: sessionID})
	if err != nil {
		return "", err
	}
	return this.provider.Encode(&security.JwtToken{
		Exp:     time.Now().Add(time.Minute * time.Duration(expireMinutes)),
		Subject: strconv.Itoa(int(user.ID)),
		Payload: payload,
	}), nil
}

func (this *UserAuthTokenService) VerifyToken(tokenString string, tokenType enum.TokenType) (*entity.UserEntity, string, error) {
	token, err := this.provider.Decode(tokenString)
	if err != nil {
		return nil, "", err
	}

	payload, err := utils.JsonUnmarshal[_TokenPayload](token.Payload)
	if err != nil {
		return nil, "", response.NewSafeError(response.CodeInvalidToken, fmt.Errorf("[UserAuthTokenService][VerifyToken] invalid token payload %w", err), utils.CallerPath(1))
	}

	if payload.Type != tokenType {
		return nil, "", response.NewSafeError(response.CodeInvalidToken, fmt.Errorf("[UserAuthTokenService][VerifyToken] invalid token type %w", err), utils.CallerPath(1))
	}

	if payload.User == nil {
		return nil, "", response.NewSafeError(response.CodeInvalidToken, errors.New("[UserAuthTokenService][VerifyToken] token payload has no user"), utils.CallerPath(1))
	}

	return payload.User, payload.SessionID, nil
}

type _TokenPayload struct {
	User *entity.UserEntity `json:"user"`
	Type enum.TokenType     `json:"type"`

	SessionID string `json:"sid,omitempty"`
}
