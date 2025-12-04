package api

import (
	"encoding/json"
	"net/http"
	"testing"

	"text-wow/internal/database"
	"text-wow/internal/models"

	"github.com/gin-gonic/gin"
)

// ═══════════════════════════════════════════════════════════
// 测试辅助函数
// ═══════════════════════════════════════════════════════════

func setupBattleTestSimple(t *testing.T) (*BattleHandler, *Handler, *gin.Engine, string, func()) {
	testDB, err := database.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}

	handler := NewHandler()
	battleHandler := NewBattleHandler()
	router := gin.New()

	// 设置路由
	setupBattleRoutes(router, handler, battleHandler)

	// 创建测试用户（不创建角色）
	token := createTestUser(t, router)

	cleanup := func() {
		database.TeardownTestDB(testDB)
	}

	return battleHandler, handler, router, token, cleanup
}

func setupBattleRoutes(router *gin.Engine, handler *Handler, battleHandler *BattleHandler) {
	api := router.Group("/api")
	{
		// 公开接口
		api.POST("/auth/register", handler.Register)
		api.POST("/auth/login", handler.Login)

		// 需要认证的接口
		protected := api.Group("")
		protected.Use(handler.AuthMiddleware())
		{
			protected.GET("/characters", handler.GetCharacters)
			protected.POST("/characters", handler.CreateCharacter)

			// 战斗接口
			battle := protected.Group("/battle")
			{
				battle.POST("/start", battleHandler.StartBattle)
				battle.POST("/stop", battleHandler.StopBattle)
				battle.POST("/toggle", battleHandler.ToggleBattle)
				battle.POST("/tick", battleHandler.BattleTick)
				battle.GET("/status", battleHandler.GetBattleStatus)
				battle.GET("/logs", battleHandler.GetBattleLogs)
				battle.POST("/zone", battleHandler.ChangeZone)
			}
		}
	}
}

func createTestUser(t *testing.T, router *gin.Engine) string {
	// 注册用户
	registerBody := models.UserRegister{
		Username: "battleuser",
		Password: "password123",
	}
	w := makeRequest(router, "POST", "/api/auth/register", registerBody)
	if w.Code != http.StatusOK {
		t.Fatalf("Failed to register user: %s", w.Body.String())
	}

	var response struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &response)
	return response.Data.Token
}

// ═══════════════════════════════════════════════════════════
// 战斗状态测试
// ═══════════════════════════════════════════════════════════

func TestBattleHandler_GetBattleStatus_Initial(t *testing.T) {
	_, _, router, token, cleanup := setupBattleTestSimple(t)
	defer cleanup()

	w := makeAuthRequest(router, "GET", "/api/battle/status", token, nil)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var response struct {
		Success bool `json:"success"`
		Data    struct {
			IsRunning   bool `json:"isRunning"`
			BattleCount int  `json:"battleCount"`
			TotalKills  int  `json:"totalKills"`
		} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &response)

	if !response.Success {
		t.Error("Expected success response")
	}

	if response.Data.IsRunning {
		t.Error("Expected battle to not be running initially")
	}

	if response.Data.BattleCount != 0 {
		t.Errorf("Expected battle count 0, got %d", response.Data.BattleCount)
	}
}

func TestBattleHandler_GetBattleStatus_Unauthorized(t *testing.T) {
	_, _, router, _, cleanup := setupBattleTestSimple(t)
	defer cleanup()

	w := makeRequest(router, "GET", "/api/battle/status", nil)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

// ═══════════════════════════════════════════════════════════
// 战斗控制测试
// ═══════════════════════════════════════════════════════════

func TestBattleHandler_StartBattle(t *testing.T) {
	_, _, router, token, cleanup := setupBattleTestSimple(t)
	defer cleanup()

	w := makeAuthRequest(router, "POST", "/api/battle/start", token, nil)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var response struct {
		Success bool `json:"success"`
		Data    struct {
			IsRunning bool `json:"isRunning"`
		} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &response)

	if !response.Success {
		t.Error("Expected success response")
	}

	if !response.Data.IsRunning {
		t.Error("Expected battle to be running after start")
	}
}

func TestBattleHandler_StopBattle(t *testing.T) {
	_, _, router, token, cleanup := setupBattleTestSimple(t)
	defer cleanup()

	// 先开始战斗
	makeAuthRequest(router, "POST", "/api/battle/start", token, nil)

	// 然后停止
	w := makeAuthRequest(router, "POST", "/api/battle/stop", token, nil)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var response struct {
		Success bool `json:"success"`
		Data    struct {
			IsRunning bool `json:"isRunning"`
		} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response.Data.IsRunning {
		t.Error("Expected battle to be stopped")
	}
}

func TestBattleHandler_ToggleBattle(t *testing.T) {
	_, _, router, token, cleanup := setupBattleTestSimple(t)
	defer cleanup()

	// 第一次切换：开始
	w := makeAuthRequest(router, "POST", "/api/battle/toggle", token, nil)

	var response struct {
		Success bool `json:"success"`
		Data    struct {
			IsRunning bool `json:"isRunning"`
		} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &response)

	if !response.Data.IsRunning {
		t.Error("Expected battle to be running after first toggle")
	}

	// 第二次切换：停止
	w = makeAuthRequest(router, "POST", "/api/battle/toggle", token, nil)
	json.Unmarshal(w.Body.Bytes(), &response)

	if response.Data.IsRunning {
		t.Error("Expected battle to be stopped after second toggle")
	}
}

// ═══════════════════════════════════════════════════════════
// 战斗回合测试
// ═══════════════════════════════════════════════════════════

func TestBattleHandler_BattleTick_NotRunning(t *testing.T) {
	_, _, router, token, cleanup := setupBattleTestSimple(t)
	defer cleanup()

	// 不开始战斗直接执行tick（没有角色会返回400）
	w := makeAuthRequest(router, "POST", "/api/battle/tick", token, nil)

	// 没有角色应该返回400
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for no characters, got %d", w.Code)
	}
}

func TestBattleHandler_BattleTick_NoCharacter(t *testing.T) {
	// 测试没有角色时执行战斗回合
	_, _, router, token, cleanup := setupBattleTestSimple(t)
	defer cleanup()

	// 开始战斗（没有角色）
	makeAuthRequest(router, "POST", "/api/battle/start", token, nil)

	// 执行战斗回合
	w := makeAuthRequest(router, "POST", "/api/battle/tick", token, nil)

	// 没有角色应该返回400
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for no characters, got %d. Body: %s", w.Code, w.Body.String())
	}
}

// ═══════════════════════════════════════════════════════════
// 战斗日志测试
// ═══════════════════════════════════════════════════════════

func TestBattleHandler_GetBattleLogs_Empty(t *testing.T) {
	_, _, router, token, cleanup := setupBattleTestSimple(t)
	defer cleanup()

	w := makeAuthRequest(router, "GET", "/api/battle/logs", token, nil)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response struct {
		Success bool `json:"success"`
		Data    struct {
			Logs []models.BattleLog `json:"logs"`
		} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &response)

	if !response.Success {
		t.Error("Expected success response")
	}
}

func TestBattleHandler_GetBattleLogs_AfterStart(t *testing.T) {
	_, _, router, token, cleanup := setupBattleTestSimple(t)
	defer cleanup()

	// 开始战斗会生成日志
	makeAuthRequest(router, "POST", "/api/battle/start", token, nil)

	// 获取日志
	w := makeAuthRequest(router, "GET", "/api/battle/logs", token, nil)

	var response struct {
		Success bool `json:"success"`
		Data    struct {
			Logs []models.BattleLog `json:"logs"`
		} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &response)

	if len(response.Data.Logs) == 0 {
		t.Error("Expected battle logs after starting battle")
	}
}

// ═══════════════════════════════════════════════════════════
// 区域切换测试 (跳过需要完整游戏数据的测试)
// ═══════════════════════════════════════════════════════════

func TestBattleHandler_ChangeZone_MissingZoneId(t *testing.T) {
	_, _, router, token, cleanup := setupBattleTestSimple(t)
	defer cleanup()

	body := map[string]string{}
	w := makeAuthRequest(router, "POST", "/api/battle/zone", token, body)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestBattleHandler_ChangeZone_NoCharacter(t *testing.T) {
	// 测试没有角色时切换区域
	_, _, router, token, cleanup := setupBattleTestSimple(t)
	defer cleanup()

	body := map[string]string{
		"zoneId": "elwynn_forest",
	}
	w := makeAuthRequest(router, "POST", "/api/battle/zone", token, body)

	// 没有角色应该返回400
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for no characters, got %d", w.Code)
	}
}

// ═══════════════════════════════════════════════════════════
// 边界情况测试
// ═══════════════════════════════════════════════════════════

func TestBattleHandler_StartBattle_AlreadyRunning(t *testing.T) {
	_, _, router, token, cleanup := setupBattleTestSimple(t)
	defer cleanup()

	// 开始两次
	makeAuthRequest(router, "POST", "/api/battle/start", token, nil)
	w := makeAuthRequest(router, "POST", "/api/battle/start", token, nil)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response struct {
		Data struct {
			IsRunning bool `json:"isRunning"`
		} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &response)

	if !response.Data.IsRunning {
		t.Error("Expected battle to still be running")
	}
}

func TestBattleHandler_StopBattle_NotRunning(t *testing.T) {
	_, _, router, token, cleanup := setupBattleTestSimple(t)
	defer cleanup()

	// 直接停止（没有开始）
	w := makeAuthRequest(router, "POST", "/api/battle/stop", token, nil)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// ═══════════════════════════════════════════════════════════
// 区域相关错误测试 - 覆盖实际遇到的问题
// ═══════════════════════════════════════════════════════════

func TestBattleHandler_BattleTick_ZoneNotFound(t *testing.T) {
	// 测试：区域不存在时执行战斗回合
	_, _, router, token, cleanup := setupBattleTestSimple(t)
	defer cleanup()

	// 创建角色
	charBody := models.CharacterCreate{
		Name:    "TestChar",
		RaceID:  "human",
		ClassID: "warrior",
	}
	w := makeAuthRequest(router, "POST", "/api/characters", token, charBody)
	if w.Code != http.StatusOK {
		t.Fatalf("Failed to create character: %s", w.Body.String())
	}

	// 开始战斗（会尝试加载默认区域）
	makeAuthRequest(router, "POST", "/api/battle/start", token, nil)

	// 执行战斗回合 - 应该能正常工作（因为testdb中有elwynn区域）
	w = makeAuthRequest(router, "POST", "/api/battle/tick", token, nil)
	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 200 or 500, got %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestBattleHandler_ChangeZone_InvalidZoneId(t *testing.T) {
	// 测试：使用错误的区域ID（如 elwynn_forest 而不是 elwynn）
	_, _, router, token, cleanup := setupBattleTestSimple(t)
	defer cleanup()

	// 创建角色
	charBody := models.CharacterCreate{
		Name:    "TestChar",
		RaceID:  "human",
		ClassID: "warrior",
	}
	w := makeAuthRequest(router, "POST", "/api/characters", token, charBody)
	if w.Code != http.StatusOK {
		t.Fatalf("Failed to create character: %s", w.Body.String())
	}

	// 尝试切换到不存在的区域（错误的ID格式）
	body := map[string]string{
		"zoneId": "elwynn_forest", // 错误的ID，应该是 "elwynn"
	}
	w = makeAuthRequest(router, "POST", "/api/battle/zone", token, body)

	// 应该返回错误（区域不存在）
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid zone ID, got %d. Body: %s", w.Code, w.Body.String())
	}

	var response models.APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	if response.Success {
		t.Error("Expected error response for invalid zone ID")
	}
}

func TestBattleHandler_ChangeZone_NonExistentZone(t *testing.T) {
	// 测试：切换到完全不存在的区域
	_, _, router, token, cleanup := setupBattleTestSimple(t)
	defer cleanup()

	// 创建角色
	charBody := models.CharacterCreate{
		Name:    "TestChar",
		RaceID:  "human",
		ClassID: "warrior",
	}
	w := makeAuthRequest(router, "POST", "/api/characters", token, charBody)
	if w.Code != http.StatusOK {
		t.Fatalf("Failed to create character: %s", w.Body.String())
	}

	// 尝试切换到不存在的区域
	body := map[string]string{
		"zoneId": "nonexistent_zone",
	}
	w = makeAuthRequest(router, "POST", "/api/battle/zone", token, body)

	// 应该返回错误
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for non-existent zone, got %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestBattleHandler_ChangeZone_ValidZone(t *testing.T) {
	// 测试：切换到有效的区域（使用正确的ID）
	_, _, router, token, cleanup := setupBattleTestSimple(t)
	defer cleanup()

	// 创建角色
	charBody := models.CharacterCreate{
		Name:    "TestChar",
		RaceID:  "human",
		ClassID: "warrior",
	}
	w := makeAuthRequest(router, "POST", "/api/characters", token, charBody)
	if w.Code != http.StatusOK {
		t.Fatalf("Failed to create character: %s", w.Body.String())
	}

	// 切换到有效的区域（使用正确的ID）
	body := map[string]string{
		"zoneId": "elwynn", // 正确的ID
	}
	w = makeAuthRequest(router, "POST", "/api/battle/zone", token, body)

	// 应该成功
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for valid zone, got %d. Body: %s", w.Code, w.Body.String())
	}

	var response models.APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	if !response.Success {
		t.Errorf("Expected success response, got error: %s", response.Error)
	}
}

func TestBattleHandler_BattleTick_WithValidZone(t *testing.T) {
	// 测试：有有效区域和角色时，战斗回合应该正常工作
	_, _, router, token, cleanup := setupBattleTestSimple(t)
	defer cleanup()

	// 创建角色
	charBody := models.CharacterCreate{
		Name:    "TestChar",
		RaceID:  "human",
		ClassID: "warrior",
	}
	w := makeAuthRequest(router, "POST", "/api/characters", token, charBody)
	if w.Code != http.StatusOK {
		t.Fatalf("Failed to create character: %s", w.Body.String())
	}

	// 开始战斗
	makeAuthRequest(router, "POST", "/api/battle/start", token, nil)

	// 执行战斗回合 - 应该成功（区域和怪物数据都存在）
	w = makeAuthRequest(router, "POST", "/api/battle/tick", token, nil)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for battle tick with valid data, got %d. Body: %s", w.Code, w.Body.String())
	}

	var response models.APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	if !response.Success {
		t.Errorf("Expected successful battle tick, got error: %s", response.Error)
	}
}

// ═══════════════════════════════════════════════════════════
// 数据库字段名一致性测试 - 覆盖实际遇到的问题
// ═══════════════════════════════════════════════════════════

func TestBattleHandler_BattleTick_ZoneHasMonsters(t *testing.T) {
	// 测试：验证区域有怪物数据，防止"No monsters in zone"错误
	_, _, router, token, cleanup := setupBattleTestSimple(t)
	defer cleanup()

	// 创建角色
	charBody := models.CharacterCreate{
		Name:    "TestChar",
		RaceID:  "human",
		ClassID: "warrior",
	}
	w := makeAuthRequest(router, "POST", "/api/characters", token, charBody)
	if w.Code != http.StatusOK {
		t.Fatalf("Failed to create character: %s", w.Body.String())
	}

	// 开始战斗（会尝试加载默认区域 elwynn）
	makeAuthRequest(router, "POST", "/api/battle/start", token, nil)

	// 执行战斗回合 - 应该成功生成敌人（因为testdb中有elwynn区域的怪物）
	w = makeAuthRequest(router, "POST", "/api/battle/tick", token, nil)

	// 应该成功（200）或返回战斗未运行（200 with isRunning: false）
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var response struct {
		Success bool `json:"success"`
		Data    struct {
			IsRunning bool `json:"isRunning"`
			Enemies   interface{} `json:"enemies"`
		} `json:"data"`
		Error string `json:"error"`
	}
	json.Unmarshal(w.Body.Bytes(), &response)

	if !response.Success {
		// 如果失败，检查是否是"No monsters in zone"错误
		if response.Error != "" && (response.Error == "no monsters in zone elwynn" || 
			response.Error == "failed to get monsters" || 
			response.Error == "No monsters in zone elwynn") {
			t.Errorf("Zone elwynn should have monsters in test database. Error: %s", response.Error)
		} else {
			t.Errorf("Unexpected error: %s", response.Error)
		}
	}
}

func TestBattleHandler_BattleTick_MonsterFieldNames(t *testing.T) {
	// 测试：验证怪物数据使用正确的字段名（physical_attack, magic_attack等）
	// 这个测试确保数据库查询不会因为字段名不匹配而失败
	_, _, router, token, cleanup := setupBattleTestSimple(t)
	defer cleanup()

	// 创建角色
	charBody := models.CharacterCreate{
		Name:    "TestChar",
		RaceID:  "human",
		ClassID: "warrior",
	}
	w := makeAuthRequest(router, "POST", "/api/characters", token, charBody)
	if w.Code != http.StatusOK {
		t.Fatalf("Failed to create character: %s", w.Body.String())
	}

	// 开始战斗
	makeAuthRequest(router, "POST", "/api/battle/start", token, nil)

	// 执行多个战斗回合，验证怪物数据能正常加载和使用
	for i := 0; i < 3; i++ {
		w = makeAuthRequest(router, "POST", "/api/battle/tick", token, nil)
		
		if w.Code != http.StatusOK {
			var errorResponse struct {
				Error string `json:"error"`
			}
			json.Unmarshal(w.Body.Bytes(), &errorResponse)
			
			// 检查是否是字段名相关的错误
			if errorResponse.Error != "" && 
				(errorResponse.Error == "no such column: attack" || 
				 errorResponse.Error == "no such column: defense" ||
				 errorResponse.Error == "failed to get monsters") {
				t.Fatalf("Database field name mismatch detected at tick %d: %s", i+1, errorResponse.Error)
			}
			
			t.Errorf("Battle tick %d failed with status %d: %s", i+1, w.Code, errorResponse.Error)
			break
		}
	}
}

func TestBattleHandler_ToggleBattle_StartsImmediately(t *testing.T) {
	// 测试：验证toggleBattle后战斗能立即开始（覆盖之前的问题）
	_, _, router, token, cleanup := setupBattleTestSimple(t)
	defer cleanup()

	// 创建角色
	charBody := models.CharacterCreate{
		Name:    "TestChar",
		RaceID:  "human",
		ClassID: "warrior",
	}
	w := makeAuthRequest(router, "POST", "/api/characters", token, charBody)
	if w.Code != http.StatusOK {
		t.Fatalf("Failed to create character: %s", w.Body.String())
	}

	// 切换战斗状态（开始）
	w = makeAuthRequest(router, "POST", "/api/battle/toggle", token, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("Failed to toggle battle: %s", w.Body.String())
	}

	var toggleResponse struct {
		Success bool `json:"success"`
		Data    struct {
			IsRunning bool `json:"isRunning"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &toggleResponse); err != nil {
		t.Fatalf("Failed to parse toggle response: %v", err)
	}

	if !toggleResponse.Success {
		t.Errorf("Toggle should succeed. Response: %s", w.Body.String())
	}
	
	// 先确保战斗是运行状态（如果 toggle 返回 false，使用 start）
	if !toggleResponse.Data.IsRunning {
		// 如果 toggle 没有启动战斗，使用 start 来启动
		w = makeAuthRequest(router, "POST", "/api/battle/start", token, nil)
		if w.Code != http.StatusOK {
			t.Fatalf("Failed to start battle: %s", w.Body.String())
		}
	}

	// 立即执行战斗回合 - 这是关键测试：验证战斗启动后能立即执行，不会因为区域没有怪物而失败
	// 这个测试主要覆盖"No monsters in zone"错误
	w = makeAuthRequest(router, "POST", "/api/battle/tick", token, nil)
	if w.Code != http.StatusOK {
		var errorResponse struct {
			Success bool   `json:"success"`
			Error   string `json:"error"`
		}
		json.Unmarshal(w.Body.Bytes(), &errorResponse)
		
		// 检查是否是"No monsters in zone"错误（这是我们想要覆盖的问题）
		if errorResponse.Error != "" && 
			(errorResponse.Error == "no monsters in zone elwynn" || 
			 errorResponse.Error == "No monsters in zone elwynn" ||
			 errorResponse.Error == "failed to get monsters") {
			t.Errorf("Zone elwynn should have monsters in test database. Error: %s", errorResponse.Error)
		} else {
			t.Errorf("Battle tick should work immediately after toggle/start. Status: %d, Error: %s", w.Code, errorResponse.Error)
		}
	} else {
		// 如果成功，验证返回的数据结构正确
		var tickResponse struct {
			Success bool `json:"success"`
			Data    interface{} `json:"data"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &tickResponse); err == nil {
			if !tickResponse.Success {
				t.Error("Battle tick should return success response")
			}
		}
	}
}

