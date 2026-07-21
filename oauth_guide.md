# Google OAuth Implementation Guide

This guide explains the implementation details, technical stack, source code, and API endpoints for the newly integrated Google OAuth authentication flow in LinkPulse.

---

## 1. Technical Architecture & Flow

The Google OAuth system is designed with two distinct authentication flows to support both standard web/mobile clients and API-only testing clients (such as Postman).

### Flow A: Standard Redirect-based Authentication (Browser)
1. **Consent Initialization**: The client requests `/auth/google/login`.
2. **Redirect to Google**: The server responds with an HTTP `302 Found` redirecting the client to Google's OAuth 2.0 Consent Screen with standard scopes (`email`, `profile`).
3. **Google Callback**: Upon user consent, Google redirects the user back to `/auth/google/callback?code=AUTH_CODE&state=STATE`.
4. **Token Exchange**: The server exchanges the authorization code for a Google Access Token and requests the user profile (`id`, `email`, `name`) from Google User Info API.
5. **Session Resolution**: 
   - If the user exists with this `google_id`, the system signs them in and generates a JWT.
   - If the user exists with the same `email` but has no `google_id` linked, the system links their Google ID and signs them in.
   - If no user exists, a new user is created with a random username (derived from name/email prefix) and an unusable password hash. A default profile is created, and the user is signed in with a generated JWT.

### Flow B: Direct Token Exchange (API / Postman / Mobile App)
1. **Token Retrieval**: The client retrieves an access token from Google directly (e.g. using Postman's built-in OAuth helper or a frontend Google SDK).
2. **Exchange Request**: The client sends a `POST` request to `/auth/google/token` containing the `access_token` in the JSON body.
3. **Validation & Resolution**: The server validates the token directly with Google User Info, then performs registration/login and returns a JWT session token.

---

## 2. Technology Stack & Dependencies

- **Core Language**: Go 1.26
- **Web Framework**: Gin Gonic (`github.com/gin-gonic/gin`)
- **OAuth Library**: `golang.org/x/oauth2` & `golang.org/x/oauth2/google` (Official Google OAuth client)
- **Session Tokens**: JWT (`github.com/golang-jwt/jwt/v5`)
- **Database Driver**: pgx (`github.com/jackc/pgx/v5`)

---

## 3. Implemented Code

### A. Database Column
We added a nullable, unique `google_id` field to the `users` table:
```sql
ALTER TABLE users ADD COLUMN IF NOT EXISTS google_id VARCHAR(255) UNIQUE;
```

### B. User Model
**File:** [internal/models/user.go](file:///d:/LinkPulse/internal/models/user.go)
```go
type User struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // "-" excludes from JSON
	Plan         string    `json:"plan"`
	GoogleID     *string   `json:"google_id,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
```

### C. User Repository
**File:** [internal/repository/user_repository.go](file:///d:/LinkPulse/internal/repository/user_repository.go)
```go
// Extends UserRepository interface:
type UserRepository interface {
    // ... other methods ...
    GetByGoogleID(ctx context.Context, googleID string) (*models.User, error)
    LinkGoogleAccount(ctx context.Context, id uuid.UUID, googleID string) error
}

// Implementation:
func (r *userRepository) GetByGoogleID(ctx context.Context, googleID string) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, plan, google_id, created_at, updated_at
		FROM users
		WHERE google_id = $1
	`
	var user models.User
	err := r.db.QueryRow(ctx, query, googleID).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Plan, &user.GoogleID, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by google id: %w", err)
	}
	return &user, nil
}

func (r *userRepository) LinkGoogleAccount(ctx context.Context, id uuid.UUID, googleID string) error {
	query := `UPDATE users SET google_id = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(ctx, query, googleID, time.Now(), id)
	return err
}
```

### D. Authentication Service
**File:** [internal/auth/service.go](file:///d:/LinkPulse/internal/auth/service.go)
```go
func (s *authService) LoginOrRegisterWithGoogle(ctx context.Context, googleID, email, name string) (string, error) {
	// 1. Check if user exists by Google ID
	user, err := s.userRepo.GetByGoogleID(ctx, googleID)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve user by Google ID: %w", err)
	}
	if user != nil {
		return GenerateToken(user.ID, user.Email)
	}

	// 2. Check if user exists by Email
	user, err = s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve user by email: %w", err)
	}
	if user != nil {
		// Link Google ID to existing account
		err = s.userRepo.LinkGoogleAccount(ctx, user.ID, googleID)
		if err != nil {
			return "", fmt.Errorf("failed to link Google account: %w", err)
		}
		return GenerateToken(user.ID, user.Email)
	}

	// 3. User does not exist, register new user
	username, err := s.generateUniqueUsername(ctx, email, name)
	if err != nil {
		return "", fmt.Errorf("failed to generate unique username: %w", err)
	}

	// Generate safe, random unusable password hash
	tempPassBytes := make([]byte, 32)
	if _, err := rand.Read(tempPassBytes); err != nil {
		return "", err
	}
	hashedPassword, err := utils.HashPassword(hex.EncodeToString(tempPassBytes))
	if err != nil {
		return "", err
	}

	newUser := &models.User{
		ID:           uuid.New(),
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword,
		Plan:         "free",
		GoogleID:     &googleID,
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	// Create default Profile
	displayName := name
	if displayName == "" {
		displayName = username
	}
	profile := &models.Profile{
		UserID:      newUser.ID,
		DisplayName: displayName,
		Bio:         "Welcome to my LinkPulse profile!",
	}
	if err := s.profileRepo.CreateProfile(profile); err != nil {
		return "", fmt.Errorf("failed to create default user profile: %w", err)
	}

	return GenerateToken(newUser.ID, newUser.Email)
}
```

### E. Handler & Helper Functions
**File:** [internal/handler/auth_handler.go](file:///d:/LinkPulse/internal/handler/auth_handler.go)
```go
// Configures OAuth client details from environment
func getGoogleOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}

// GoogleLogin initiates OAuth redirection
func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	config := getGoogleOAuthConfig()
	if config.ClientID == "" || config.ClientSecret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Google OAuth is not configured on the server"})
		return
	}
	state := "state-token"
	url := config.AuthCodeURL(state, oauth2.AccessTypeOnline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleCallback processes OAuth auth code
func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	state := c.Query("state")
	if state != "state-token" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state parameter"})
		return
	}
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing code parameter"})
		return
	}

	config := getGoogleOAuthConfig()
	token, err := config.Exchange(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange authorization code: " + err.Error()})
		return
	}

	userInfo, err := fetchGoogleUserInfo(token.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user info: " + err.Error()})
		return
	}

	jwtToken, err := h.authService.LoginOrRegisterWithGoogle(c.Request.Context(), userInfo.ID, userInfo.Email, userInfo.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": jwtToken})
}

// GoogleTokenExchange validates access token directly (API-first client)
func (h *AuthHandler) GoogleTokenExchange(c *gin.Context) {
	var req GoogleTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: access_token is required"})
		return
	}

	userInfo, err := fetchGoogleUserInfo(req.AccessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to verify token with Google: " + err.Error()})
		return
	}

	jwtToken, err := h.authService.LoginOrRegisterWithGoogle(c.Request.Context(), userInfo.ID, userInfo.Email, userInfo.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": jwtToken})
}
```

---

## 4. Endpoint URLs & Postman Payload

Before calling redirect-based login, configure credentials in your `.env` file:
```env
GOOGLE_CLIENT_ID=your-google-client-id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your-client-secret
GOOGLE_REDIRECT_URL=http://localhost:8080/auth/google/callback
```

### Endpoint 1: Google Redirection Login (Browser)
* **Method**: `GET`
* **URL**: `http://localhost:8080/auth/google/login`
* **Description**: Initiates the consent flow and redirects users to Google's sign-in page.

### Endpoint 2: Google Callback (Redirect Receiver)
* **Method**: `GET`
* **URL**: `http://localhost:8080/auth/google/callback`
* **Parameters**: `code` (string), `state` (string)
* **Description**: The redirect URI Google calls. It returns a JWT key.

### Endpoint 3: Google Access Token Exchange (Postman / API First)
* **Method**: `POST`
* **URL**: `http://localhost:8080/auth/google/token`
* **Description**: Allows API clients (like Postman or mobile applications) to submit a Google Access Token directly to authenticate and receive a local JWT.
* **JSON Body**:
  ```json
  {
    "access_token": "ya29.a0ARWqs..."
  }
  ```
* **Success Response (200 OK)**:
  ```json
  {
    "message": "Login successful",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
  ```
