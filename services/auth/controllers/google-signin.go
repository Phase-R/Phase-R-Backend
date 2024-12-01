package controllers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
	"github.com/gorilla/sessions"
	"github.com/nrednav/cuid2"
	"github.com/Phase-R/Phase-R-Backend/services/auth/db"
	"github.com/Phase-R/Phase-R-Backend/db/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
	"github.com/golang-jwt/jwt/v4"
)

// Register models.User type with gob to allow serialization
func init() {
	gob.Register(models.User{})
}

var (
	oauthConfGl = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  "http://localhost:8080/user/google/auth", // Update with your correct redirect URL
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
	oauthStateStringGl = ""
)

var store = sessions.NewCookieStore([]byte("secret"))

// InitializeSessionMiddleware initializes the session middleware.
func InitializeSessionMiddleware(r *gin.Engine) {
	// Using the gorilla cookie store
	r.Use(func(c *gin.Context) {
		session, err := store.Get(c.Request, "mysession")
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to get session: %v", err)
			c.Abort()
			return
		}
		c.Set("session", session)
		c.Next()
	})
}

// InitializeOAuthGoogle sets up OAuth configuration using environment variables.
func InitializeOAuthGoogle() {
	oauthConfGl.ClientID = os.Getenv("CLIENT_ID")
	oauthConfGl.ClientSecret = os.Getenv("CLIENT_SECRET")
	oauthStateStringGl = generateStateOauthCookie()
}

// generateStateOauthCookie creates a secure random state parameter.
func generateStateOauthCookie() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// HandleGoogleLogin redirects users to Google's OAuth 2.0 authentication page.
func HandleGoogleLogin(c *gin.Context) {
	URL, err := url.Parse(oauthConfGl.Endpoint.AuthURL)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to parse AuthURL: %v", err)
		return
	}

	parameters := url.Values{}
	parameters.Add("client_id", oauthConfGl.ClientID)
	parameters.Add("scope", strings.Join(oauthConfGl.Scopes, " "))
	parameters.Add("redirect_uri", oauthConfGl.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", oauthStateStringGl)
	URL.RawQuery = parameters.Encode()

	c.Redirect(http.StatusTemporaryRedirect, URL.String())
}

// CallBackFromGoogle handles the OAuth 2.0 callback from Google.
func CallBackFromGoogle(c *gin.Context) {
    state := c.Query("state")
    if state != oauthStateStringGl {
        c.String(http.StatusBadRequest, "Invalid state parameter")
        return
    }

    code := c.Query("code")
    if code == "" {
        c.String(http.StatusBadRequest, "Authorization code not found")
        return
    }

    token, err := oauthConfGl.Exchange(context.Background(), code)
    if err != nil {
        c.String(http.StatusInternalServerError, "Failed to exchange token: %v", err)
        return
    }

    client := oauthConfGl.Client(context.Background(), token)
    resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
    if err != nil {
        c.String(http.StatusInternalServerError, "Failed to get user info: %v", err)
        return
    }
    defer resp.Body.Close()

    var userInfo map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
        c.String(http.StatusInternalServerError, "Failed to parse user info: %v", err)
        return
    }

    email, _ := userInfo["email"].(string)
    fname, _ := userInfo["given_name"].(string)
    lname, _ := userInfo["family_name"].(string)

    if email == "" {
        c.String(http.StatusBadRequest, "Email not found in user info")
        return
    }

    id := cuid2.Generate()
    if id == "" {
        c.String(http.StatusInternalServerError, "Failed to generate user ID")
        return
    }

    // Check if the user exists in the database
    var user models.User
    result := db.DB.Where("email = ?", email).First(&user)
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            // User does not exist, create a new user
            user = models.User{
                ID:       id,
                Username: fname + lname,
                Fname:    fname,
                Lname:    lname,
                Email:    email,
                Password: "",
                Age:      0,
                Access:   "free",
                Verified: true,
            }
            if err := db.DB.Create(&user).Error; err != nil {
                c.String(http.StatusInternalServerError, "Failed to create user: %v", err)
                return
            }
        } else {
            c.String(http.StatusInternalServerError, "Database error: %v", result.Error)
            return
        }
    }

    // Generate JWT token
    claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "iss": user.Email,
        "exp": time.Now().Add(24 * time.Hour).Unix(),
    })

    tokenString, err := claims.SignedString([]byte(os.Getenv("SECRET_KEY")))
    if err != nil {
        c.String(http.StatusInternalServerError, "Token generation failed: %v", err)
        return
    }

    // Set the JWT token as a cookie (use HTTPS in production)
    c.SetSameSite(http.SameSiteLaxMode)
    c.SetCookie("Auth", tokenString, 3600*24*30, "", "", false, false)

    // Save user info in session
    session, err := store.Get(c.Request, "mysession")
    if err != nil {
        c.String(http.StatusInternalServerError, "Failed to get session: %v", err)
        return
    }
    session.Values["user"] = user
    if err := session.Save(c.Request, c.Writer); err != nil {
        c.String(http.StatusInternalServerError, "Failed to save session: %v", err)
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message":      "Login successful",
        "isAuthenticated": true,  
    })
}
