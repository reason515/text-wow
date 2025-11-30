package repository

import (
	"testing"

	"text-wow/internal/auth"
	"text-wow/internal/database"
)

// ═══════════════════════════════════════════════════════════
// 测试辅助函数
// ═══════════════════════════════════════════════════════════

func setupUserRepoTest(t *testing.T) (*UserRepository, func()) {
	testDB, err := database.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}

	repo := NewUserRepository()
	cleanup := func() {
		database.TeardownTestDB(testDB)
	}

	return repo, cleanup
}

// ═══════════════════════════════════════════════════════════
// Create 测试
// ═══════════════════════════════════════════════════════════

func TestUserRepository_Create(t *testing.T) {
	repo, cleanup := setupUserRepoTest(t)
	defer cleanup()

	passwordHash, _ := auth.HashPassword("testpassword")

	user, err := repo.Create("testuser", passwordHash, "test@example.com")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if user == nil {
		t.Fatal("Create returned nil user")
	}

	if user.ID == 0 {
		t.Error("User ID should not be 0")
	}

	if user.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", user.Username)
	}

	if user.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", user.Email)
	}

	// 验证默认值 (根据 schema.sql: max_team_size=5, unlocked_slots=1)
	if user.MaxTeamSize != 5 {
		t.Errorf("Expected MaxTeamSize 5, got %d", user.MaxTeamSize)
	}

	if user.UnlockedSlots != 1 {
		t.Errorf("Expected UnlockedSlots 1, got %d", user.UnlockedSlots)
	}

	if user.Gold != 0 {
		t.Errorf("Expected Gold 0, got %d", user.Gold)
	}
}

func TestUserRepository_Create_WithEmptyEmail(t *testing.T) {
	repo, cleanup := setupUserRepoTest(t)
	defer cleanup()

	passwordHash, _ := auth.HashPassword("testpassword")

	user, err := repo.Create("testuser", passwordHash, "")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// 空邮箱应该被存储为空字符串（从NULL转换）
	if user.Email != "" {
		t.Errorf("Expected empty email, got '%s'", user.Email)
	}
}

func TestUserRepository_Create_DuplicateUsername(t *testing.T) {
	repo, cleanup := setupUserRepoTest(t)
	defer cleanup()

	passwordHash, _ := auth.HashPassword("testpassword")

	// 第一次创建
	_, err := repo.Create("duplicateuser", passwordHash, "first@example.com")
	if err != nil {
		t.Fatalf("First create failed: %v", err)
	}

	// 第二次创建相同用户名
	_, err = repo.Create("duplicateuser", passwordHash, "second@example.com")
	if err == nil {
		t.Error("Expected error for duplicate username, got nil")
	}
}

// ═══════════════════════════════════════════════════════════
// GetByID 测试
// ═══════════════════════════════════════════════════════════

func TestUserRepository_GetByID(t *testing.T) {
	repo, cleanup := setupUserRepoTest(t)
	defer cleanup()

	passwordHash, _ := auth.HashPassword("testpassword")
	created, _ := repo.Create("testuser", passwordHash, "test@example.com")

	user, err := repo.GetByID(created.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if user.ID != created.ID {
		t.Errorf("Expected ID %d, got %d", created.ID, user.ID)
	}

	if user.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", user.Username)
	}
}

func TestUserRepository_GetByID_NotFound(t *testing.T) {
	repo, cleanup := setupUserRepoTest(t)
	defer cleanup()

	_, err := repo.GetByID(99999)
	if err == nil {
		t.Error("Expected error for non-existent user, got nil")
	}
}

// ═══════════════════════════════════════════════════════════
// GetByUsername 测试
// ═══════════════════════════════════════════════════════════

func TestUserRepository_GetByUsername(t *testing.T) {
	repo, cleanup := setupUserRepoTest(t)
	defer cleanup()

	passwordHash, _ := auth.HashPassword("testpassword")
	repo.Create("findme", passwordHash, "findme@example.com")

	user, err := repo.GetByUsername("findme")
	if err != nil {
		t.Fatalf("GetByUsername failed: %v", err)
	}

	if user.Username != "findme" {
		t.Errorf("Expected username 'findme', got '%s'", user.Username)
	}
}

func TestUserRepository_GetByUsername_NotFound(t *testing.T) {
	repo, cleanup := setupUserRepoTest(t)
	defer cleanup()

	_, err := repo.GetByUsername("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent username, got nil")
	}
}

// ═══════════════════════════════════════════════════════════
// GetPasswordHash 测试
// ═══════════════════════════════════════════════════════════

func TestUserRepository_GetPasswordHash(t *testing.T) {
	repo, cleanup := setupUserRepoTest(t)
	defer cleanup()

	password := "mysecretpassword"
	passwordHash, _ := auth.HashPassword(password)
	created, _ := repo.Create("hashuser", passwordHash, "")

	userID, hash, err := repo.GetPasswordHash("hashuser")
	if err != nil {
		t.Fatalf("GetPasswordHash failed: %v", err)
	}

	if userID != created.ID {
		t.Errorf("Expected userID %d, got %d", created.ID, userID)
	}

	// 验证密码hash可以正确验证
	if !auth.CheckPassword(password, hash) {
		t.Error("Password hash should be valid")
	}
}

func TestUserRepository_GetPasswordHash_NotFound(t *testing.T) {
	repo, cleanup := setupUserRepoTest(t)
	defer cleanup()

	_, _, err := repo.GetPasswordHash("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent user, got nil")
	}
}

// ═══════════════════════════════════════════════════════════
// UsernameExists 测试
// ═══════════════════════════════════════════════════════════

func TestUserRepository_UsernameExists(t *testing.T) {
	repo, cleanup := setupUserRepoTest(t)
	defer cleanup()

	passwordHash, _ := auth.HashPassword("testpassword")
	repo.Create("existinguser", passwordHash, "")

	// 测试存在的用户名
	exists, err := repo.UsernameExists("existinguser")
	if err != nil {
		t.Fatalf("UsernameExists failed: %v", err)
	}
	if !exists {
		t.Error("Expected username to exist")
	}

	// 测试不存在的用户名
	exists, err = repo.UsernameExists("nonexistent")
	if err != nil {
		t.Fatalf("UsernameExists failed: %v", err)
	}
	if exists {
		t.Error("Expected username to not exist")
	}
}

// ═══════════════════════════════════════════════════════════
// UpdateLastLogin 测试
// ═══════════════════════════════════════════════════════════

func TestUserRepository_UpdateLastLogin(t *testing.T) {
	repo, cleanup := setupUserRepoTest(t)
	defer cleanup()

	passwordHash, _ := auth.HashPassword("testpassword")
	created, _ := repo.Create("loginuser", passwordHash, "")

	// 初始时last_login_at应该为空
	user, _ := repo.GetByID(created.ID)
	if user.LastLoginAt != nil {
		t.Error("Expected LastLoginAt to be nil initially")
	}

	// 更新登录时间
	err := repo.UpdateLastLogin(created.ID)
	if err != nil {
		t.Fatalf("UpdateLastLogin failed: %v", err)
	}

	// 验证登录时间已更新
	user, _ = repo.GetByID(created.ID)
	if user.LastLoginAt == nil {
		t.Error("Expected LastLoginAt to be set after update")
	}
}

// ═══════════════════════════════════════════════════════════
// UpdateGold 测试
// ═══════════════════════════════════════════════════════════

func TestUserRepository_UpdateGold(t *testing.T) {
	repo, cleanup := setupUserRepoTest(t)
	defer cleanup()

	passwordHash, _ := auth.HashPassword("testpassword")
	created, _ := repo.Create("golduser", passwordHash, "")

	// 初始金币为0
	user, _ := repo.GetByID(created.ID)
	if user.Gold != 0 {
		t.Errorf("Expected initial gold 0, got %d", user.Gold)
	}

	// 增加金币
	err := repo.UpdateGold(created.ID, 100)
	if err != nil {
		t.Fatalf("UpdateGold failed: %v", err)
	}

	user, _ = repo.GetByID(created.ID)
	if user.Gold != 100 {
		t.Errorf("Expected gold 100, got %d", user.Gold)
	}

	// 验证total_gold_gained也增加了
	if user.TotalGoldGained != 100 {
		t.Errorf("Expected TotalGoldGained 100, got %d", user.TotalGoldGained)
	}
}

// ═══════════════════════════════════════════════════════════
// UpdateZone 测试
// ═══════════════════════════════════════════════════════════

func TestUserRepository_UpdateZone(t *testing.T) {
	repo, cleanup := setupUserRepoTest(t)
	defer cleanup()

	passwordHash, _ := auth.HashPassword("testpassword")
	created, _ := repo.Create("zoneuser", passwordHash, "")

	// 更新区域
	err := repo.UpdateZone(created.ID, "durotar")
	if err != nil {
		t.Fatalf("UpdateZone failed: %v", err)
	}

	user, _ := repo.GetByID(created.ID)
	if user.CurrentZoneID != "durotar" {
		t.Errorf("Expected zone 'durotar', got '%s'", user.CurrentZoneID)
	}
}

// ═══════════════════════════════════════════════════════════
// IncrementKills 测试
// ═══════════════════════════════════════════════════════════

func TestUserRepository_IncrementKills(t *testing.T) {
	repo, cleanup := setupUserRepoTest(t)
	defer cleanup()

	passwordHash, _ := auth.HashPassword("testpassword")
	created, _ := repo.Create("killuser", passwordHash, "")

	// 初始击杀数为0
	user, _ := repo.GetByID(created.ID)
	if user.TotalKills != 0 {
		t.Errorf("Expected initial kills 0, got %d", user.TotalKills)
	}

	// 增加击杀数
	for i := 0; i < 5; i++ {
		err := repo.IncrementKills(created.ID)
		if err != nil {
			t.Fatalf("IncrementKills failed: %v", err)
		}
	}

	user, _ = repo.GetByID(created.ID)
	if user.TotalKills != 5 {
		t.Errorf("Expected kills 5, got %d", user.TotalKills)
	}
}

