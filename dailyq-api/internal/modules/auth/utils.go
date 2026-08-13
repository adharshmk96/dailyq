package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Claims is the JWT payload. The session id is carried in the standard `jti`
// claim so every token maps to a revocable server-side session.
type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func newID() string { return uuid.NewString() }

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func hashPassword(password string, cost int) (string, error) {
	if cost <= 0 {
		cost = bcrypt.DefaultCost
	}
	b, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func checkPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// generateResetToken returns a URL-safe random token and its SHA-256 hash.
// Only the hash is persisted.
func generateResetToken() (token string, tokenHash string, err error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", "", err
	}
	token = base64.RawURLEncoding.EncodeToString(buf)
	return token, hashToken(token), nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// signToken issues a JWT bound to the given session.
func signToken(secret, issuer, userID, email, sessionID string, ttl time.Duration, now time.Time) (string, time.Time, error) {
	expiresAt := now.Add(ttl)

	claims := Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        sessionID,
			Subject:   userID,
			Issuer:    issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		return "", time.Time{}, err
	}
	return signed, expiresAt, nil
}

func parseToken(secret, tokenString string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, ErrInvalidToken
	}
	if claims.Subject == "" || claims.ID == "" {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

// bearerToken extracts the token from an `Authorization: Bearer <token>` header.
func bearerToken(header string) (string, bool) {
	parts := strings.SplitN(strings.TrimSpace(header), " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", false
	}
	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", false
	}
	return token, true
}
