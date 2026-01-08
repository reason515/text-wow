-- 配置版本表 - 用于配置版本管理和热更新
CREATE TABLE IF NOT EXISTS config_versions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    config_type VARCHAR(32) NOT NULL,  -- monster/skill/item/economy/zone
    version INTEGER NOT NULL,
    config_data TEXT NOT NULL,         -- JSON格式的配置数据
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(32),
    description TEXT,
    UNIQUE(config_type, version)
);

CREATE INDEX IF NOT EXISTS idx_config_versions_type ON config_versions(config_type);
CREATE INDEX IF NOT EXISTS idx_config_versions_version ON config_versions(config_type, version DESC);
























































