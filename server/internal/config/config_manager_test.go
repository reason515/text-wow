package config

import (
	"testing"

	"text-wow/internal/database"

	"github.com/stretchr/testify/assert"
)

func setupConfigTest(t *testing.T) (*ConfigManager, func()) {
	testDB, err := database.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}

	cleanup := func() {
		database.TeardownTestDB(testDB)
	}

	cm := NewConfigManager()
	return cm, cleanup
}

func TestNewConfigManager(t *testing.T) {
	cm := NewConfigManager()
	assert.NotNil(t, cm, "ConfigManager should not be nil")
	assert.NotNil(t, cm.configs, "configs map should be initialized")
	assert.NotNil(t, cm.listeners, "listeners slice should be initialized")
	assert.NotNil(t, cm.versionManager, "versionManager should be initialized")
}

func TestLoadConfig_InvalidType(t *testing.T) {
	cm, cleanup := setupConfigTest(t)
	defer cleanup()

	err := cm.LoadConfig("invalid_type")
	assert.Error(t, err, "Should return error for invalid config type")
	assert.Contains(t, err.Error(), "unknown config type", "Error message should mention unknown config type")
}

func TestGetConfigVersion_NotLoaded(t *testing.T) {
	cm, cleanup := setupConfigTest(t)
	defer cleanup()

	version := cm.getConfigVersion("monster")
	assert.Equal(t, 0, version, "Version should be 0 for not loaded config")
}

func TestRegisterConfigChangeListener(t *testing.T) {
	cm, cleanup := setupConfigTest(t)
	defer cleanup()

	listener := &testConfigListener{}
	cm.RegisterConfigChangeListener(listener)

	assert.Equal(t, 1, len(cm.listeners), "Should have one listener registered")
}

func TestReloadConfig_NotLoaded(t *testing.T) {
	cm, cleanup := setupConfigTest(t)
	defer cleanup()

	// 尝试重新加载未加载的配置
	err := cm.ReloadConfig("monster")
	// 这可能会失败，因为数据库可能没有数据，但应该不会panic
	_ = err
}

// testConfigListener 测试用的配置变更监听器
type testConfigListener struct {
	changeCount int
	lastType    string
	lastVersion int
}

func (l *testConfigListener) OnConfigChange(configType string, version int) {
	l.changeCount++
	l.lastType = configType
	l.lastVersion = version
}

func TestConfigChangeNotification(t *testing.T) {
	cm, cleanup := setupConfigTest(t)
	defer cleanup()

	listener := &testConfigListener{}
	cm.RegisterConfigChangeListener(listener)

	// 模拟配置变更通知
	cm.notifyConfigChange("monster", 1)

	assert.Equal(t, 1, listener.changeCount, "Listener should be notified once")
	assert.Equal(t, "monster", listener.lastType, "Last type should be monster")
	assert.Equal(t, 1, listener.lastVersion, "Last version should be 1")
}





















































