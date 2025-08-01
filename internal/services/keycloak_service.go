package services

import (
	"auth-service/internal/models"
	"context"
	"fmt"
	"os"

	"github.com/Nerzal/gocloak/v13"
)

var (
	KEYCLOAK_ADMIN_USERNAME = os.Getenv("KEYCLOAK_ADMIN_USERNAME")
	KEYCLOAK_ADMIN_PASSWORD = os.Getenv("KEYCLOAK_ADMIN_PASSWORD")
	KEYCLOAK_ADMIN_REALM    = os.Getenv("KEYCLOAK_ADMIN_REALM")

	hostname = "sdasd"
)

type KeycloakService struct {
	Gocloak      *gocloak.GoCloak
	ClientId     string
	ClientSecret string
	Realm        string
	Hostname     string
}

func NewKeycloakService(client_id string, client_secret string, realm string, hostname string) *KeycloakService {
	return &KeycloakService{
		Gocloak:      gocloak.NewClient(hostname),
		ClientId:     client_id,
		ClientSecret: client_secret,
		Realm:        realm,
		Hostname:     hostname,
	}
}

func (ks *KeycloakService) Login(login models.LoginParams) (*models.LoginResponse, error) {
	ctx := context.Background()
	token, err := ks.Gocloak.Login(ctx, login.Email, login.Password, ks.ClientId, ks.ClientSecret, ks.Realm)
	if err != nil {
		return nil, fmt.Errorf("login fail: %w", err)
	}

	// Response modelimize dönüştür
	response := &models.LoginResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    token.ExpiresIn,
		TokenType:    token.TokenType,
	}
	return response, nil
}

func (ks *KeycloakService) Register(register models.RegisterParams) error {
	ctx := context.Background()

	adminToken, err := ks.Gocloak.LoginAdmin(ctx, KEYCLOAK_ADMIN_USERNAME, KEYCLOAK_ADMIN_PASSWORD, KEYCLOAK_ADMIN_REALM)
	if err != nil {
		return fmt.Errorf("admin login failed: %w", err)
	}

	user := gocloak.User{
		FirstName: gocloak.StringP(register.Firstname),
		LastName:  gocloak.StringP(register.Lastname),
		Username:  gocloak.StringP(register.Username),
		Email:     gocloak.StringP(register.Login.Email),
		Enabled:   gocloak.BoolP(true),
	}

	userID, err := ks.Gocloak.CreateUser(ctx, adminToken.AccessToken, ks.Realm, user)
	if err != nil {
		fmt.Println("keycloak createUser error:", err)
		return fmt.Errorf("create user failed: %w", err)
	}

	cred := gocloak.CredentialRepresentation{
		Type:      gocloak.StringP("password"),
		Value:     gocloak.StringP(register.Login.Password),
		Temporary: gocloak.BoolP(false),
	}
	err = ks.Gocloak.SetPassword(ctx, adminToken.AccessToken, userID, ks.Realm, *cred.Value, false)
	if err != nil {
		return fmt.Errorf("setting password fail: %w", err)
	}
	return nil
}

func (ks *KeycloakService) GetUserByID(userID string) (*gocloak.User, error) {
	ctx := context.Background()
	adminToken, err := ks.Gocloak.LoginAdmin(ctx, KEYCLOAK_ADMIN_USERNAME, KEYCLOAK_ADMIN_PASSWORD, KEYCLOAK_ADMIN_REALM)
	if err != nil {
		return nil, fmt.Errorf("admin login failed: %w", err)
	}

	user, err := ks.Gocloak.GetUserByID(ctx, adminToken.AccessToken, ks.Realm, userID)
	if err != nil {
		return nil, fmt.Errorf("get user failed: %w", err)
	}
	return user, nil
}

func (ks *KeycloakService) UpdateUser(userID string, user gocloak.User) error {
	ctx := context.Background()
	adminToken, err := ks.Gocloak.LoginAdmin(ctx, KEYCLOAK_ADMIN_USERNAME, KEYCLOAK_ADMIN_PASSWORD, KEYCLOAK_ADMIN_REALM)
	if err != nil {
		return fmt.Errorf("admin login failed: %w", err)
	}

	user.ID = gocloak.StringP(userID)

	err = ks.Gocloak.UpdateUser(ctx, adminToken.AccessToken, ks.Realm, user)
	if err != nil {
		return fmt.Errorf("update user failed: %w", err)
	}
	return nil
}

func (ks *KeycloakService) DeleteUser(userID string) error {
	ctx := context.Background()
	adminToken, err := ks.Gocloak.LoginAdmin(ctx, KEYCLOAK_ADMIN_USERNAME, KEYCLOAK_ADMIN_PASSWORD, KEYCLOAK_ADMIN_REALM)
	if err != nil {
		return fmt.Errorf("admin login failed: %w", err)
	}

	err = ks.Gocloak.DeleteUser(ctx, adminToken.AccessToken, ks.Realm, userID)
	if err != nil {
		return fmt.Errorf("delete user failed: %w", err)
	}
	return nil
}
func (ks *KeycloakService) GetUserProfile(accessToken string) (*gocloak.User, error) {
	ctx := context.Background()
	
	userInfo, err := ks.Gocloak.GetUserInfo(ctx, accessToken, ks.Realm)
	if err != nil {
		return nil, fmt.Errorf("get user info failed: %w", err)
	}

	// Admin token ile kullanıcı detaylarını al
	adminToken, err := ks.Gocloak.LoginAdmin(ctx, KEYCLOAK_ADMIN_USERNAME, KEYCLOAK_ADMIN_PASSWORD, KEYCLOAK_ADMIN_REALM)
	if err != nil {
		return nil, fmt.Errorf("admin login failed: %w", err)
	}

	// UserInfo'dan gelen sub (subject) ID'sini kullanarak tam kullanıcı bilgisini al
	user, err := ks.Gocloak.GetUserByID(ctx, adminToken.AccessToken, ks.Realm, *userInfo.Sub)
	if err != nil {
		return nil, fmt.Errorf("get user failed: %w", err)
	}
	
	return user, nil
}