package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"text-wow/internal/auth"
	"text-wow/internal/database"
	"text-wow/internal/models"

	"github.com/gin-gonic/gin"
)

// ═══════════════════════════════════════════════════════════
// 测试辅助函数
// ═══════════════════════════════════════════════════════════

func init() {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)
}

func setupHandlerTest(t *testing.T) (*Handler, *gin.Engine, func()) {
	testDB, err := database.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}

	handler := NewHandler()
	router := gin.New()

	// 设置路由
	setupRoutes(router, handler)

	cleanup := func() {
		database.TeardownTestDB(testDB)
	}

	return handler, router, cleanup
}

func setupRoutes(router *gin.Engine, handler *Handler) {
	api := router.Group("/api")
	{
		// 公开接口
		api.POST("/auth/register", handler.Register)
		api.POST("/auth/login", handler.Login)
		api.GET("/races", handler.GetRaces)
		api.GET("/classes", handler.GetClasses)

		// 需要认证的接口
		protected := api.Group("")
		protected.Use(handler.AuthMiddleware())
		{
			protected.GET("/user", handler.GetCurrentUser)
			protected.GET("/characters", handler.GetCharacters)
			protected.POST("/characters", handler.CreateCharacter)
		}
	}
}

func makeRequest(router *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBody)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req, _ := http.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func makeAuthRequest(router *gin.Engine, method, path, token string, body interface{}) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBody)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req, _ := http.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func parseResponse(w *httptest.ResponseRecorder) models.APIResponse {
	var response models.APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	return response
}

// ═══════════════════════════════════════════════════════════
// 注册测试
// ═══════════════════════════════════════════════════════════

func TestHandler_Register_Success(t *testing.T) {
	_, router, cleanup := setupHandlerTest(t)
	defer cleanup()

	body := models.UserRegister{
		Username: "newuser",
		Password: "password123",
		Email:    "test@example.com",
	}

	w := makeRequest(router, "POST", "/api/auth/register", body)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	response := parseResponse(w)
	if !response.Success {
		t.Errorf("Expected success, got error: %s", response.Error)
	}

	// 验证返回了token和用户信息
	if response.Data == nil {
		t.Error("Expected data in response")
	}
}

func TestHandler_Register_InvalidRequest(t *testing.T) {
	_, router, cleanup := setupHandlerTest(t)
	defer cleanup()

	// 缺少必填字段
	body := map[string]string{
		"username": "testuser",
		// missing password
	}

	w := makeRequest(router, "POST", "/api/auth/register", body)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestHandler_Register_DuplicateUsername(t *testing.T) {
	_, router, cleanup := setupHandlerTest(t)
	defer cleanup()

	body := models.UserRegister{
		Username: "duplicateuser",
		Password: "password123",
	}

	// 第一次注册
	w := makeRequest(router, "POST", "/api/auth/register", body)
	if w.Code != http.StatusOK {
		t.Fatalf("First registration failed: %s", w.Body.String())
	}

	// 第二次注册相同用户名
	w = makeRequest(router, "POST", "/api/auth/register", body)
	if w.Code != http.StatusConflict {
		t.Errorf("Expected status 409, got %d", w.Code)
	}
}

func TestHandler_Register_WithoutEmail(t *testing.T) {
	_, router, cleanup := setupHandlerTest(t)
	defer cleanup()

	body := models.UserRegister{
		Username: "noemailuser",
		Password: "password123",
	}

	w := makeRequest(router, "POST", "/api/auth/register", body)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}
}

// ═══════════════════════════════════════════════════════════
// 登录测试
// ═══════════════════════════════════════════════════════════

func TestHandler_Login_Success(t *testing.T) {
	_, router, cleanup := setupHandlerTest(t)
	defer cleanup()

	// 先注册
	registerBody := models.UserRegister{
		Username: "loginuser",
		Password: "password123",
	}
	makeRequest(router, "POST", "/api/auth/register", registerBody)

	// 然后登录
	loginBody := models.UserCredentials{
		Username: "loginuser",
		Password: "password123",
	}
	w := makeRequest(router, "POST", "/api/auth/login", loginBody)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	response := parseResponse(w)
	if !response.Success {
		t.Errorf("Expected success, got error: %s", response.Error)
	}
}

func TestHandler_Login_WrongPassword(t *testing.T) {
	_, router, cleanup := setupHandlerTest(t)
	defer cleanup()

	// 先注册
	registerBody := models.UserRegister{
		Username: "wrongpassuser",
		Password: "correctpassword",
	}
	makeRequest(router, "POST", "/api/auth/register", registerBody)

	// 用错误密码登录
	loginBody := models.UserCredentials{
		Username: "wrongpassuser",
		Password: "wrongpassword",
	}
	w := makeRequest(router, "POST", "/api/auth/login", loginBody)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestHandler_Login_NonexistentUser(t *testing.T) {
	_, router, cleanup := setupHandlerTest(t)
	defer cleanup()

	loginBody := models.UserCredentials{
		Username: "nonexistentuser",
		Password: "anypassword",
	}
	w := makeRequest(router, "POST", "/api/auth/login", loginBody)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

// ═══════════════════════════════════════════════════════════
// 认证中间件测试
// ═══════════════════════════════════════════════════════════

func TestHandler_AuthMiddleware_NoToken(t *testing.T) {
	_, router, cleanup := setupHandlerTest(t)
	defer cleanup()

	w := makeRequest(router, "GET", "/api/user", nil)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestHandler_AuthMiddleware_InvalidToken(t *testing.T) {
	_, router, cleanup := setupHandlerTest(t)
	defer cleanup()

	w := makeAuthRequest(router, "GET", "/api/user", "invalid.token.here", nil)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestHandler_AuthMiddleware_ValidToken(t *testing.T) {
	_, router, cleanup := setupHandlerTest(t)
	defer cleanup()

	// 先注册获取token
	registerBody := models.UserRegister{
		Username: "authuser",
		Password: "password123",
	}
	w := makeRequest(router, "POST", "/api/auth/register", registerBody)

	var response struct {
		Success bool `json:"success"`
		Data    struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &response)

	// 使用token访问受保护接口
	w = makeAuthRequest(router, "GET", "/api/user", response.Data.Token, nil)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestHandler_AuthMiddleware_InvalidFormat(t *testing.T) {
	_, router, cleanup := setupHandlerTest(t)
	defer cleanup()

	// 不使用 Bearer 前缀
	req, _ := http.NewRequest("GET", "/api/user", nil)
	req.Header.Set("Authorization", "InvalidFormat token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

// ═══════════════════════════════════════════════════════════
// GetCurrentUser 测试
// ═══════════════════════════════════════════════════════════

func TestHandler_GetCurrentUser(t *testing.T) {
	_, router, cleanup := setupHandlerTest(t)
	defer cleanup()

	// 注册用户
	registerBody := models.UserRegister{
		Username: "currentuser",
		Password: "password123",
		Email:    "current@example.com",
	}
	w := makeRequest(router, "POST", "/api/auth/register", registerBody)

	var response struct {
		Success bool `json:"success"`
		Data    struct {
			Token string      `json:"token"`
			User  models.User `json:"user"`
		} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &response)

	// 获取当前用户
	w = makeAuthRequest(router, "GET", "/api/user", response.Data.Token, nil)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var userResponse struct {
		Success bool        `json:"success"`
		Data    models.User `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &userResponse)

	if userResponse.Data.Username != "currentuser" {
		t.Errorf("Expected username 'currentuser', got '%s'", userResponse.Data.Username)
	}
}

// ═══════════════════════════════════════════════════════════
// 游戏配置API测试
// ═══════════════════════════════════════════════════════════

func TestHandler_GetRaces(t *testing.T) {
	_, router, cleanup := setupHandlerTest(t)
	defer cleanup()

	w := makeRequest(router, "GET", "/api/races", nil)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	response := parseResponse(w)
	if !response.Success {
		t.Errorf("Expected success, got error: %s", response.Error)
	}
}

func TestHandler_GetClasses(t *testing.T) {
	_, router, cleanup := setupHandlerTest(t)
	defer cleanup()

	w := makeRequest(router, "GET", "/api/classes", nil)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	response := parseResponse(w)
	if !response.Success {
		t.Errorf("Expected success, got error: %s", response.Error)
	}
}

// ═══════════════════════════════════════════════════════════
// 角色创建测试
// ═══════════════════════════════════════════════════════════

func TestHandler_CreateCharacter_Success(t *testing.T) {
	t.Skip("Skipping: requires full game data (races, classes) in test database")
	// This test would require the full seed.sql data
}

func TestHandler_CreateCharacter_DuplicateName(t *testing.T) {
	t.Skip("Skipping: requires full game data (races, classes) in test database")
	// This test would require the full seed.sql data
}

func TestHandler_CreateCharacter_InvalidRace(t *testing.T) {
	_, router, cleanup := setupHandlerTest(t)
	defer cleanup()

	// 注册用户
	registerBody := models.UserRegister{
		Username: "invalidraceuser",
		Password: "password123",
	}
	w := makeRequest(router, "POST", "/api/auth/register", registerBody)

	var response struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &response)

	// 使用无效种族
	charBody := models.CharacterCreate{
		Name:    "InvalidRaceHero",
		RaceID:  "invalid_race",
		ClassID: "warrior",
	}
	w = makeAuthRequest(router, "POST", "/api/characters", response.Data.Token, charBody)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

// ═══════════════════════════════════════════════════════════
// Token安全性测试
// ═══════════════════════════════════════════════════════════

func TestHandler_TokenSecurity_DifferentUsersGetDifferentTokens(t *testing.T) {
	_, router, cleanup := setupHandlerTest(t)
	defer cleanup()

	// 注册第一个用户
	registerBody1 := models.UserRegister{
		Username: "user1",
		Password: "password123",
	}
	w1 := makeRequest(router, "POST", "/api/auth/register", registerBody1)

	var response1 struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	json.Unmarshal(w1.Body.Bytes(), &response1)

	// 注册第二个用户
	registerBody2 := models.UserRegister{
		Username: "user2",
		Password: "password123",
	}
	w2 := makeRequest(router, "POST", "/api/auth/register", registerBody2)

	var response2 struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	json.Unmarshal(w2.Body.Bytes(), &response2)

	// 验证token不同
	if response1.Data.Token == response2.Data.Token {
		t.Error("Different users should get different tokens")
	}
}

func TestHandler_TokenContainsCorrectClaims(t *testing.T) {
	_, router, cleanup := setupHandlerTest(t)
	defer cleanup()

	// 注册用户
	registerBody := models.UserRegister{
		Username: "claimsuser",
		Password: "password123",
	}
	w := makeRequest(router, "POST", "/api/auth/register", registerBody)

	var response struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &response)

	// 验证token中的claims
	claims, err := auth.ValidateToken(response.Data.Token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if claims.Username != "claimsuser" {
		t.Errorf("Expected username 'claimsuser', got '%s'", claims.Username)
	}
}

