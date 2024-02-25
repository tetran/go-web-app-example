package auth

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/tetran/go-web-app-example/clock"
	"github.com/tetran/go-web-app-example/entity"
)

const (
	RoleKey     = "role"
	UserNameKey = "user_name"
)

//go:embed cert/secret.pem
var rawPrivKey []byte

//go:embed cert/public.pem
var rawPubKey []byte

//go:generate go run github.com/matryer/moq -out moq_test.go . Store
type Store interface {
	Save(ctx context.Context, key string, userID entity.UserID) error
	Load(ctx context.Context, key string) (entity.UserID, error)
}

type JWTer struct {
	PrivateKey, PulicKey jwk.Key
	Store                Store
	Clocker              clock.Clocker
}

func NewJWTer(s Store, clocker clock.Clocker) (*JWTer, error) {
	j := &JWTer{Store: s}
	privKey, err := parse(rawPrivKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	pubkey, err := parse(rawPubKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	j.PrivateKey = privKey
	j.PulicKey = pubkey
	j.Clocker = clocker
	return j, nil
}

func (j *JWTer) GenerateToken(ctx context.Context, u entity.User) ([]byte, error) {
	tok, err := jwt.NewBuilder().
		JwtID(uuid.New().String()).
		Issuer(`go-web-app-example`).
		Subject("access_token").
		IssuedAt(j.Clocker.Now()).
		Expiration(j.Clocker.Now().Add(30*time.Minute)).
		Claim(RoleKey, u.Role).
		Claim(UserNameKey, u.Name).
		Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build token: %w", err)
	}

	if err := j.Store.Save(ctx, tok.JwtID(), u.ID); err != nil {
		return nil, fmt.Errorf("failed to save token: %w", err)
	}

	signed, err := jwt.Sign(tok, jwt.WithKey(jwa.RS256, j.PrivateKey))
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %w", err)
	}

	return signed, nil
}

func (j *JWTer) GetToken(ctx context.Context, r *http.Request) (jwt.Token, error) {
	token, err := jwt.ParseRequest(r, jwt.WithKey(jwa.RS256, j.PulicKey), jwt.WithValidate(false))
	if err != nil {
		return nil, err
	}

	if err := jwt.Validate(token, jwt.WithClock(j.Clocker)); err != nil {
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}

	if _, err := j.Store.Load(ctx, token.JwtID()); err != nil {
		return nil, fmt.Errorf("token is expired: %w", err)
	}

	return token, nil
}

func (j *JWTer) FillContext(r *http.Request) (*http.Request, error) {
	token, err := j.GetToken(r.Context(), r)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	uid, err := j.Store.Load(r.Context(), token.JwtID())
	if err != nil {
		return nil, fmt.Errorf("failed to load token: %w", err)
	}

	ctx := SetUserID(r.Context(), uid)
	ctx = SetRole(ctx, token)
	clone := r.Clone(ctx)
	return clone, nil
}

func parse(raw []byte) (jwk.Key, error) {
	key, err := jwk.ParseKey(raw, jwk.WithPEM(true))
	if err != nil {
		return nil, err
	}

	return key, nil
}

type userIDKey struct{}
type roleKey struct{}

func SetUserID(ctx context.Context, uid entity.UserID) context.Context {
	return context.WithValue(ctx, userIDKey{}, uid)
}

func GetUserID(ctx context.Context) (entity.UserID, bool) {
	uid, ok := ctx.Value(userIDKey{}).(entity.UserID)
	return uid, ok
}

func SetRole(ctx context.Context, tok jwt.Token) context.Context {
	get, ok := tok.Get(RoleKey)
	if !ok {
		return context.WithValue(ctx, roleKey{}, "")
	}

	return context.WithValue(ctx, roleKey{}, get.(string))
}

func GetRole(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(roleKey{}).(string)
	return role, ok
}

func IsAdmin(ctx context.Context) bool {
	role, ok := GetRole(ctx)
	return ok && role == "admin"
}
