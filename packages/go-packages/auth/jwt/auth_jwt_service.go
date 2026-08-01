package jwt

import (
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	signingMethod = jwt.SigningMethodHS256
)

type authJWTService struct {
	issuer                     string
	audience                   string
	secretKey                  []byte
	accessTokenExpirationTime  time.Duration
	refreshTokenExpirationTime time.Duration
	leeway                     time.Duration
}

// New returns an AuthJWTService configured with the given issuer, audience,
// secret key, token expiration times, and validation leeway. It returns
// ErrEmptyIssuer if issuer is empty, ErrEmptyAudience if audience is empty, or
// ErrWeakSecretKey if secretKey is empty or shorter than 32 bytes. A non-empty
// issuer and audience are required so that issuer and audience validation are
// actually enforced during token verification.
func New(
	issuer string,
	audience string,
	secretKey string,
	accessTokenExpirationTime time.Duration,
	refreshTokenExpirationTime time.Duration,
	leeway time.Duration,
) (AuthJWTService, error) {
	if strings.TrimSpace(issuer) == "" {
		return nil, ErrEmptyIssuer
	}
	if strings.TrimSpace(audience) == "" {
		return nil, ErrEmptyAudience
	}

	trimmedKey := strings.TrimSpace(secretKey)
	if trimmedKey == "" || len(trimmedKey) < 32 {
		return nil, ErrWeakSecretKey
	}

	return &authJWTService{
		issuer:                     issuer,
		audience:                   audience,
		secretKey:                  []byte(trimmedKey),
		accessTokenExpirationTime:  accessTokenExpirationTime,
		refreshTokenExpirationTime: refreshTokenExpirationTime,
		leeway:                     leeway,
	}, nil
}

func (s *authJWTService) GenerateAccessToken(userID uuid.UUID, sessionID uuid.UUID) (string, error) {
	now := time.Now().UTC()

	claims := AccessTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Audience:  []string{s.audience},
			Subject:   userID.String(),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTokenExpirationTime)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
		TokenType: TokenTypeAccessToken,
		SessionID: sessionID.String(),
	}

	token := jwt.NewWithClaims(signingMethod, claims)
	return token.SignedString(s.secretKey)
}

func (s *authJWTService) GenerateRefreshToken(userID uuid.UUID, sessionID uuid.UUID) (string, error) {
	now := time.Now().UTC()

	claims := RefreshTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Audience:  []string{s.audience},
			Subject:   userID.String(),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshTokenExpirationTime)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
		TokenType: TokenTypeRefreshToken,
		SessionID: sessionID.String(),
	}

	token := jwt.NewWithClaims(signingMethod, claims)
	return token.SignedString(s.secretKey)
}

func (s *authJWTService) ValidateAccessToken(tokenString string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errUnexpectedSigningMethod
		}
		return s.secretKey, nil
	},
		jwt.WithValidMethods([]string{signingMethod.Name}),
		jwt.WithIssuer(s.issuer),
		jwt.WithIssuedAt(),
		jwt.WithAudience(s.audience),
		jwt.WithExpirationRequired(),
		jwt.WithLeeway(s.leeway),
	)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			slog.Debug("token expired",
				"component", "auth_jwt_service",
			)
			return nil, ErrExpiredToken
		}

		slog.Debug("token validation failed",
			"component", "auth_jwt_service",
			"error", err.Error(),
		)
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*AccessTokenClaims); ok && token.Valid {
		if claims.TokenType != TokenTypeAccessToken {
			return nil, ErrInvalidTokenType
		}
		return claims, nil
	}

	return nil, ErrInvalidToken
}

func (s *authJWTService) ValidateRefreshToken(tokenString string) (*RefreshTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshTokenClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errUnexpectedSigningMethod
		}
		return s.secretKey, nil
	},
		jwt.WithValidMethods([]string{signingMethod.Name}),
		jwt.WithIssuer(s.issuer),
		jwt.WithIssuedAt(),
		jwt.WithAudience(s.audience),
		jwt.WithExpirationRequired(),
		jwt.WithLeeway(s.leeway),
	)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			slog.Debug("refresh token expired",
				"component", "auth_jwt_service",
			)
			return nil, ErrExpiredToken
		}

		slog.Debug("refresh token validation failed",
			"component", "auth_jwt_service",
			"error", err.Error(),
		)
		return nil, ErrInvalidRefreshToken
	}

	if claims, ok := token.Claims.(*RefreshTokenClaims); ok && token.Valid {
		if claims.TokenType != TokenTypeRefreshToken {
			return nil, ErrInvalidTokenType
		}
		return claims, nil
	}

	return nil, ErrInvalidRefreshToken
}

func (s *authJWTService) AccessTokenTTL() time.Duration {
	return s.accessTokenExpirationTime
}

func (s *authJWTService) RefreshTokenTTL() time.Duration {
	return s.refreshTokenExpirationTime
}
