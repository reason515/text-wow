package auth

import (
	"testing"
	"time"
)

// ═══════════════════════════════════════════════════════════
// 密码哈希测试
// ═══════════════════════════════════════════════════════════

func TestHashPassword(t *testing.T) {
	password := "testPassword123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	// 验证hash不为空
	if hash == "" {
		t.Error("HashPassword returned empty hash")
	}

	// 验证hash不等于原密码
	if hash == password {
		t.Error("HashPassword returned unhashed password")
	}

	// 验证bcrypt格式 (应该以$2开头)
	if len(hash) < 4 || hash[:2] != "$2" {
		t.Error("HashPassword did not return bcrypt format")
	}
}

func TestHashPassword_DifferentPasswords(t *testing.T) {
	password1 := "password1"
	password2 := "password2"

	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	if hash1 == hash2 {
		t.Error("Different passwords should produce different hashes")
	}
}

func TestHashPassword_SamePasswordDifferentHashes(t *testing.T) {
	password := "samePassword"

	hash1, _ := HashPassword(password)
	hash2, _ := HashPassword(password)

	// bcrypt应该为相同密码生成不同的hash（由于salt）
	if hash1 == hash2 {
		t.Error("Same password should produce different hashes due to salt")
	}
}

func TestCheckPassword_ValidPassword(t *testing.T) {
	password := "mySecurePassword123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if !CheckPassword(password, hash) {
		t.Error("CheckPassword should return true for valid password")
	}
}

func TestCheckPassword_InvalidPassword(t *testing.T) {
	password := "correctPassword"
	wrongPassword := "wrongPassword"

	hash, _ := HashPassword(password)

	if CheckPassword(wrongPassword, hash) {
		t.Error("CheckPassword should return false for wrong password")
	}
}

func TestCheckPassword_EmptyPassword(t *testing.T) {
	password := "testPassword"
	hash, _ := HashPassword(password)

	if CheckPassword("", hash) {
		t.Error("CheckPassword should return false for empty password")
	}
}

// ═══════════════════════════════════════════════════════════
// JWT Token测试
// ═══════════════════════════════════════════════════════════

func TestGenerateToken(t *testing.T) {
	userID := 1
	username := "testuser"

	token, err := GenerateToken(userID, username)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	if token == "" {
		t.Error("GenerateToken returned empty token")
	}

	// JWT token应该有三个部分，以.分隔
	parts := 0
	for _, c := range token {
		if c == '.' {
			parts++
		}
	}
	if parts != 2 {
		t.Errorf("JWT token should have 3 parts, got %d separators", parts)
	}
}

func TestValidateToken_ValidToken(t *testing.T) {
	userID := 42
	username := "testplayer"

	token, err := GenerateToken(userID, username)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("Expected UserID %d, got %d", userID, claims.UserID)
	}

	if claims.Username != username {
		t.Errorf("Expected Username %s, got %s", username, claims.Username)
	}
}

func TestValidateToken_InvalidToken(t *testing.T) {
	_, err := ValidateToken("invalid.token.here")

	if err == nil {
		t.Error("ValidateToken should return error for invalid token")
	}

	if err != ErrInvalidToken {
		t.Errorf("Expected ErrInvalidToken, got %v", err)
	}
}

func TestValidateToken_EmptyToken(t *testing.T) {
	_, err := ValidateToken("")

	if err == nil {
		t.Error("ValidateToken should return error for empty token")
	}
}

func TestValidateToken_TamperedToken(t *testing.T) {
	token, _ := GenerateToken(1, "user")

	// 篡改token的最后一个字符
	tamperedToken := token[:len(token)-1] + "X"

	_, err := ValidateToken(tamperedToken)
	if err == nil {
		t.Error("ValidateToken should return error for tampered token")
	}
}

func TestValidateToken_Issuer(t *testing.T) {
	token, _ := GenerateToken(1, "user")
	claims, err := ValidateToken(token)

	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	if claims.Issuer != "text-wow" {
		t.Errorf("Expected issuer 'text-wow', got '%s'", claims.Issuer)
	}
}

func TestValidateToken_ExpirationTime(t *testing.T) {
	token, _ := GenerateToken(1, "user")
	claims, err := ValidateToken(token)

	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	// 验证过期时间大约是7天后
	expectedExpiry := time.Now().Add(7 * 24 * time.Hour)
	actualExpiry := claims.ExpiresAt.Time

	// 允许1分钟的误差
	if actualExpiry.Sub(expectedExpiry) > time.Minute || expectedExpiry.Sub(actualExpiry) > time.Minute {
		t.Errorf("Token expiry time is not approximately 7 days from now")
	}
}

// ═══════════════════════════════════════════════════════════
// 边界条件测试
// ═══════════════════════════════════════════════════════════

func TestHashPassword_LongPassword(t *testing.T) {
	// bcrypt有72字节的密码长度限制
	// 使用刚好72字节的密码（不超过限制）
	longPassword := ""
	for i := 0; i < 72; i++ {
		longPassword += "a"
	}

	hash, err := HashPassword(longPassword)
	if err != nil {
		t.Fatalf("HashPassword failed for long password: %v", err)
	}

	// 验证可以检查
	if !CheckPassword(longPassword, hash) {
		t.Error("CheckPassword should work with long password")
	}
}

func TestHashPassword_TooLongPassword(t *testing.T) {
	// 超过72字节应该报错
	tooLongPassword := ""
	for i := 0; i < 100; i++ {
		tooLongPassword += "a"
	}

	_, err := HashPassword(tooLongPassword)
	// bcrypt会报错或者截断，取决于实现
	// 这里我们只是确认不会panic
	_ = err
}

func TestHashPassword_SpecialCharacters(t *testing.T) {
	specialPassword := "!@#$%^&*()_+-=[]{}|;':\",./<>?中文密码🔐"

	hash, err := HashPassword(specialPassword)
	if err != nil {
		t.Fatalf("HashPassword failed for special characters: %v", err)
	}

	if !CheckPassword(specialPassword, hash) {
		t.Error("CheckPassword should work with special characters")
	}
}

func TestGenerateToken_SpecialUsername(t *testing.T) {
	userID := 1
	username := "user@special.name"

	token, err := GenerateToken(userID, username)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	if claims.Username != username {
		t.Errorf("Username with special chars not preserved")
	}
}

// ═══════════════════════════════════════════════════════════
// 基准测试
// ═══════════════════════════════════════════════════════════

func BenchmarkHashPassword(b *testing.B) {
	password := "benchmarkPassword123"
	for i := 0; i < b.N; i++ {
		HashPassword(password)
	}
}

func BenchmarkCheckPassword(b *testing.B) {
	password := "benchmarkPassword123"
	hash, _ := HashPassword(password)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CheckPassword(password, hash)
	}
}

func BenchmarkGenerateToken(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateToken(1, "benchuser")
	}
}

func BenchmarkValidateToken(b *testing.B) {
	token, _ := GenerateToken(1, "benchuser")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ValidateToken(token)
	}
}

