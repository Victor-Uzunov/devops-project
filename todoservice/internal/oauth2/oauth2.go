package oauth2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Victor-Uzunov/devops-project/todoservice/internal/db"
	"github.com/Victor-Uzunov/devops-project/todoservice/internal/users"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/constants"
	token "github.com/Victor-Uzunov/devops-project/todoservice/pkg/jwt"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/log"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jmoiron/sqlx"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Handler struct {
	oauth2Config          *oauth2.Config
	oauth2State           string
	jwtKey                []byte
	jwtExpirationTime     time.Duration
	refreshExpirationTime time.Duration
	userService           users.UserService
	database              *sqlx.DB
}

type GitHubData struct {
	User struct {
		Email string `json:"email"`
	} `json:"user"`
	Role string `json:"role"`
}

type Tokens struct {
	AccessToken string `json:"access_token"`
}

func NewOAuth2(auth2 token.ConfigOAuth2, userService users.UserService, database *sqlx.DB) *Handler {
	return &Handler{
		oauth2Config: &oauth2.Config{
			ClientID:     auth2.ClientID,
			ClientSecret: auth2.ClientSecret,
			RedirectURL:  auth2.RedirectURL,
			Scopes:       auth2.Scopes,
			Endpoint:     github.Endpoint,
		},
		oauth2State:           auth2.OAuth2State,
		jwtKey:                []byte(auth2.JwtKey),
		jwtExpirationTime:     auth2.JWTExpirationTime,
		refreshExpirationTime: auth2.RefreshExpirationTime,
		userService:           userService,
		database:              database,
	}
}

func (h *Handler) RootHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, `<a href="/login/github">LOGIN</a>`)
	if err != nil {
		return
	}
}

func (h *Handler) GithubLoginHandler(w http.ResponseWriter, r *http.Request) {
	log.C(r.Context()).Info("github login handler")
	url := h.oauth2Config.AuthCodeURL(h.oauth2State)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler) GithubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	log.C(r.Context()).Info("github callback received")
	code := r.URL.Query().Get("code")
	tokenJWT, err := h.oauth2Config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.C(r.Context()).Errorf("failed to exchange code for token: %v", err)
		http.Error(w, "token exchange failed", http.StatusInternalServerError)
		return
	}
	log.C(r.Context()).Debugf("exchanged token:  %s", tokenJWT.AccessToken)

	accessToken := tokenJWT.AccessToken
	githubData, err := h.getGithubData(r.Context(), accessToken)
	if err != nil {
		log.C(r.Context()).Debugf("failed to fetch github data: %v", err)
		http.Error(w, "Failed to get github data", http.StatusInternalServerError)
		return
	}
	log.C(r.Context()).Debugf("successfully parsed github data: %v", githubData)

	h.loggedInHandler(w, r, githubData)
}

func (h *Handler) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.C(ctx).Info("refresh token received handler")

	cookie, err := r.Cookie("refresh_token")
	if err != nil || cookie.Value == "" {
		http.Error(w, "Refresh token not found", http.StatusUnauthorized)
		return
	}

	tx, err := h.database.BeginTxx(ctx, nil)
	if err != nil {
		log.C(ctx).Errorf("GenerateJWT refresh token transaction failed: %v", err)
		http.Error(w, "failed to start transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	ctx = db.SaveToContext(ctx, tx)

	user, err := h.userService.FindByRefreshToken(ctx, cookie.Value)
	if err != nil {
		log.C(ctx).Errorf("invalid refresh token: %v", err)
		http.Error(w, "invalid refresh token", http.StatusUnauthorized)
		return
	}
	log.C(ctx).Debugf("found user: %v", user)

	newAccessToken, err := h.GenerateJWT(ctx, h.jwtExpirationTime, user.Email, string(user.Role))
	if err != nil {
		log.C(ctx).Errorf("failed to generate new access token: %v", err)
		http.Error(w, "failed to generate new access token", http.StatusInternalServerError)
		return
	}

	newRefreshToken, err := h.generateRefreshToken(ctx, user.Email, string(user.Role))
	if err != nil {
		log.C(r.Context()).Errorf("failed to generate new refresh token: %v", err)
		http.Error(w, "failed to generate new refresh token", http.StatusInternalServerError)
		return
	}

	if err := h.userService.SaveRefreshToken(ctx, user.Email, newRefreshToken, time.Now().Add(h.refreshExpirationTime)); err != nil {
		log.C(r.Context()).Errorf("failed to save new refresh token: %v", err)
		http.Error(w, "failed to save new refresh token", http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		log.C(ctx).Errorf("JWTMiddleware transaction failed to commit: %v", err)
		http.Error(w, "failed to commit transaction", http.StatusInternalServerError)
		return
	}

	tokens := Tokens{
		AccessToken: newAccessToken,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(tokens); err != nil {
		log.C(r.Context()).Errorf("failed to write response: %v", err)
		http.Error(w, "failed to write response", http.StatusInternalServerError)
	}
}

func (h *Handler) getGithubData(context context.Context, accessToken string) (string, error) {
	log.C(context).Info("getting github data")
	userData, err := h.getGithubUserData(context, accessToken)
	if err != nil {
		return "", err
	}
	log.C(context).Debugf("successfully fetched github data: %v", userData)
	userOrg, err := h.getGithubUserOrg(context, accessToken)
	if err != nil {
		return "", err
	}
	log.C(context).Debugf("successfully fetched github organizations for a user: %v", userOrg)
	role, err := pkg.DetermineUserRole(userOrg)
	if err != nil {
		return "", err
	}
	log.C(context).Debugf("successfully determined user role: %v", role)

	fullData := fmt.Sprintf(`{
		"user": %s,
		"role": "%s",
		"organizations": %s
	}`, userData, role, userOrg)

	return fullData, nil
}

func (h *Handler) getGithubUserData(context context.Context, accessToken string) (string, error) {
	log.C(context).Info("getting github user data")
	client := h.oauth2Config.Client(oauth2.NoContext, &oauth2.Token{AccessToken: accessToken})
	resp, err := client.Get("https://api.github.com/user")
	log.C(context).Info("making get request to github for a user data")

	if err != nil {
		return "", err
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	log.C(context).Debugf("got github user data: %s", string(respBody))
	return string(respBody), nil
}

func (h *Handler) getGithubUserOrg(context context.Context, accessToken string) (string, error) {
	log.C(context).Info("getting github user organizations")
	client := h.oauth2Config.Client(oauth2.NoContext, &oauth2.Token{AccessToken: accessToken})
	resp, err := client.Get("https://api.github.com/user/orgs")
	log.C(context).Info("making get request to github for a user organizations")
	if err != nil {
		return "", err
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	log.C(context).Debugf("got github user organizations: %s", string(respBody))
	return string(respBody), nil
}

func (h *Handler) loggedInHandler(w http.ResponseWriter, r *http.Request, githubData string) {
	log.C(r.Context()).Info("logged in handler")

	if githubData == "" {
		log.C(r.Context()).Error("unauthorized access")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(githubData), "", "\t"); err != nil {
		log.C(r.Context()).Errorf("JSON parse error: %v", err)
		http.Error(w, "failed to parse JSON", http.StatusInternalServerError)
		return
	}

	var data GitHubData
	if err := json.Unmarshal(prettyJSON.Bytes(), &data); err != nil {
		log.C(r.Context()).Errorf("JSON parse error: %v", err)
		http.Error(w, "failed to parse JSON", http.StatusInternalServerError)
		return
	}
	tokenJWT, err := h.GenerateJWT(r.Context(), h.jwtExpirationTime, data.User.Email, data.Role)
	if err != nil {
		log.C(r.Context()).Errorf("JWT generation error: %v", err)
		http.Error(w, "failed to generate JWT", http.StatusInternalServerError)
		return
	}
	refreshToken, err := h.generateRefreshToken(r.Context(), data.User.Email, data.Role)
	if err != nil {
		log.C(r.Context()).Errorf("JWT refresh generation error: %v", err)
		http.Error(w, "failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   constants.CookieAge,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "user_role",
		Value:    data.Role,
		HttpOnly: false,
		Path:     "/",
		MaxAge:   constants.CookieAge,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    tokenJWT,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   constants.CookieAge,
	})

	redirectURL := "http://localhost:8000/"
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func (h *Handler) saveRefreshToken(ctx context.Context, email string, refreshToken string, expirationTime time.Time) error {
	log.C(ctx).Infof("saving refresh token for user: %v", email)
	err := h.userService.SaveRefreshToken(ctx, email, refreshToken, expirationTime)
	if err != nil {
		log.C(ctx).Errorf("failed to save refresh token: %v", err)
		return err
	}
	return nil
}

func (h *Handler) GenerateJWT(ctx context.Context, expTime time.Duration, email, role string) (string, error) {
	log.C(ctx).Info("generating JWT token")
	expirationTime := time.Now().Add(24 * expTime * 7)

	tx, err := h.database.BeginTxx(ctx, nil)
	if err != nil {
		log.C(ctx).Errorf("GenerateJWT transaction failed: %v", err)
		return "", err
	}
	defer tx.Rollback()

	ctx = db.SaveToContext(ctx, tx)
	_, err = h.userService.GetUserByEmail(ctx, email)
	if err != nil {
		if !strings.Contains(err.Error(), "user not found") {
			log.C(ctx).Errorf("error while getting user from the database: %v", err)
			return "", err
		}

		u := models.User{
			Email:    email,
			GithubID: email + "Git",
			Role:     pkg.StringToRole(role),
		}
		_, err = h.userService.CreateUser(ctx, u)
		if err != nil {
			log.C(ctx).Errorf("error while creating user in generate JWT token: %v", err)
			return "", err
		}
	}
	user, err := h.userService.GetUserByEmail(ctx, email)
	if err != nil {
		log.C(ctx).Errorf("error while getting user from the database: %v", err)
		return "", err
	}
	err = tx.Commit()
	if err != nil {
		log.C(ctx).Errorf("JWTMiddleware transaction failed to commit: %v", err)
		return "", err
	}

	claims := &token.Claims{
		ID:    user.ID,
		Email: email,
		Role:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	log.C(ctx).Debugf("claim for the token is: %v", claims)

	tokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := tokenJWT.SignedString(h.jwtKey)
	if err != nil {
		log.C(ctx).Errorf("failed to sign token: %v", err)
		return "", err
	}
	log.C(ctx).Debugf("generated JWT token: %s", tokenString)

	return tokenString, nil
}

func (h *Handler) generateRefreshToken(ctx context.Context, email string, role string) (string, error) {
	log.C(ctx).Info("generating refresh token")
	expirationTime := time.Now().Add(h.refreshExpirationTime)

	claims := &token.Claims{
		Email: email,
		Role:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	log.C(ctx).Debugf("claim for the token is: %v", claims)

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := refreshToken.SignedString(h.jwtKey)
	if err != nil {
		log.C(ctx).Errorf("failed to sign refresh token: %v", err)
		return "", err
	}

	tx, err := h.database.BeginTxx(ctx, nil)
	if err != nil {
		log.C(ctx).Errorf("GenerateJWT refresh token transaction failed: %v", err)
		return "", err
	}
	defer tx.Rollback()

	ctx = db.SaveToContext(ctx, tx)

	err = h.saveRefreshToken(ctx, email, tokenString, expirationTime)
	if err != nil {
		log.C(ctx).Errorf("failed to save refresh token: %v", err)
		return "", err
	}

	err = tx.Commit()
	if err != nil {
		log.C(ctx).Errorf("Refresh token transaction failed to commit: %v", err)
		return "", err
	}

	log.C(ctx).Debugf("generated refresh token: %s", tokenString)

	return tokenString, nil
}
