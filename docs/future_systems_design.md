# 未来扩展系统设计文档（概要）

> 📌 **核心设计理念**: 轻量级联网功能，支持单机游戏体验，玩家精力开销较小

---

## 📋 目录

1. [系统概览](#系统概览)
2. [交易系统](#交易系统)
3. [排名系统](#排名系统)
4. [聊天系统](#聊天系统)
5. [成就系统](#成就系统)

---

## 系统概览

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          未来扩展系统架构                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   轻量级联网功能:                                                             │
│   ├─ 交易系统: 玩家间装备交易，异步交易                                       │
│   ├─ 排名系统: 排行榜展示，定期更新                                           │
│   ├─ 聊天系统: 世界聊天、公告板，轻量级                                       │
│   └─ 成就系统: 成就追踪，奖励发放                                             │
│                                                                             │
│   设计原则:                                                                  │
│   ├─ 单机优先: 游戏核心体验为单机，联网功能为辅助                             │
│   ├─ 轻量级: 玩家精力开销小，无需实时在线                                    │
│   ├─ 异步设计: 大部分功能支持异步，无需实时交互                              │
│   └─ 可选参与: 玩家可以选择是否参与联网功能                                  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 交易系统

### 系统概览

交易系统允许玩家间进行装备和金币交易，促进装备流通，增加游戏社交性。

### 核心功能

#### 1. 玩家间直接交易

- **面对面交易**: 两个玩家实时交易
- **交易确认**: 双方确认后完成交易
- **支持物品**: 装备、金币、材料

#### 2. 拍卖行系统

- **上架装备**: 玩家上架装备到拍卖行
- **浏览购买**: 其他玩家浏览和购买
- **价格排序**: 支持按价格、品质、等级排序
- **筛选功能**: 按职业、槽位、品质筛选

### 设计特点

- **异步交易**: 上架后无需在线，其他玩家可以购买
- **交易安全**: 交易确认机制，防止欺诈
- **交易记录**: 记录所有交易，便于查询
- **手续费机制**: 上架费和交易手续费，防止刷装备

### 数据库设计

```sql
-- 拍卖行表（已在装备系统设计中）
-- 交易记录表（已在装备系统设计中）
```

### 接口设计

```go
// TradingManager 交易管理器接口
type TradingManager interface {
    // 创建交易
    CreateTrade(sellerID, buyerID int, items []Item, price int) error
    
    // 上架装备
    ListItem(sellerID int, equipmentID int, price int) error
    
    // 购买装备
    BuyItem(buyerID int, auctionID int) error
    
    // 取消上架
    CancelListing(sellerID int, auctionID int) error
}
```

---

## 排名系统

### 系统概览

排名系统展示玩家排行榜，提供竞争和成就感，定期更新，无需实时同步。

### 核心功能

#### 1. 排行榜类型

- **等级排行榜**: 按角色等级排序
- **战力排行榜**: 按综合战力排序
- **击杀排行榜**: 按总击杀数排序
- **财富排行榜**: 按总金币数排序
- **深渊排行榜**: 按深渊层数排序

#### 2. 排名更新

- **定期更新**: 每小时更新一次
- **增量更新**: 只更新变化的排名
- **缓存机制**: 排名数据缓存，减少数据库压力

### 设计特点

- **轻量级**: 定期更新，无需实时同步
- **多维度**: 多种排行榜，满足不同玩家需求
- **奖励机制**: 排名奖励，激励玩家竞争

### 数据库设计

```sql
-- 排行榜表
CREATE TABLE leaderboards (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    category VARCHAR(32) NOT NULL,  -- level/power/kills/gold/abyss
    user_id INTEGER NOT NULL,
    character_id INTEGER,
    score INTEGER NOT NULL,
    rank INTEGER,
    last_updated DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (character_id) REFERENCES characters(id)
);

CREATE INDEX idx_leaderboard_category ON leaderboards(category, score DESC);
```

### 接口设计

```go
// RankingManager 排名管理器接口
type RankingManager interface {
    // 更新排名
    UpdateRanking(userID int, category string, score int) error
    
    // 获取排行榜
    GetLeaderboard(category string, limit int) ([]RankingEntry, error)
    
    // 获取玩家排名
    GetUserRank(userID int, category string) (int, error)
}
```

---

## 聊天系统

### 系统概览

聊天系统提供世界聊天和公告板功能，轻量级设计，支持玩家交流和官方公告。

### 核心功能

#### 1. 世界聊天

- **公共频道**: 所有玩家可见
- **消息限制**: 防止刷屏（如每分钟3条）
- **消息历史**: 保存最近100条消息
- **敏感词过滤**: 自动过滤敏感词

#### 2. 公告板

- **官方公告**: 游戏更新、活动通知
- **系统消息**: 重要系统消息
- **玩家公告**: 玩家可以发布交易、组队等信息

### 设计特点

- **轻量级**: 消息历史有限，减少存储压力
- **异步设计**: 消息发送和接收异步处理
- **可选参与**: 玩家可以选择关闭聊天

### 数据库设计

```sql
-- 聊天消息表
CREATE TABLE chat_messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    channel VARCHAR(32) NOT NULL,  -- world/trade/party
    message TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_chat_channel ON chat_messages(channel, created_at DESC);

-- 公告表
CREATE TABLE announcements (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    type VARCHAR(32) NOT NULL,  -- official/system/player
    title VARCHAR(128) NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME
);
```

### 接口设计

```go
// ChatManager 聊天管理器接口
type ChatManager interface {
    // 发送消息
    SendMessage(userID int, channel string, message string) error
    
    // 获取消息历史
    GetMessageHistory(channel string, limit int) ([]ChatMessage, error)
    
    // 发布公告
    CreateAnnouncement(announcement *Announcement) error
    
    // 获取公告
    GetAnnouncements(limit int) ([]Announcement, error)
}
```

---

## 成就系统

### 系统概览

成就系统追踪玩家游戏进度，提供成就和奖励，增加游戏长期目标。

### 核心功能

#### 1. 成就类型

- **战斗成就**: 击杀数、胜利数、连杀等
- **收集成就**: 装备收集、角色收集等
- **探索成就**: 区域探索、副本完成等
- **成长成就**: 等级、属性、技能等

#### 2. 成就奖励

- **金币奖励**: 完成成就获得金币
- **装备奖励**: 完成成就获得特殊装备
- **称号奖励**: 完成成就获得称号
- **经验奖励**: 完成成就获得经验

### 设计特点

- **进度追踪**: 自动追踪成就进度
- **奖励机制**: 完成成就获得奖励
- **成就展示**: 展示已完成的成就

### 数据库设计

```sql
-- 成就配置表
CREATE TABLE achievements (
    id VARCHAR(32) PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    description TEXT,
    category VARCHAR(32) NOT NULL,  -- battle/collection/exploration/growth
    condition_type VARCHAR(32) NOT NULL,  -- kill_count/level/item_collect
    condition_value INTEGER NOT NULL,
    reward_type VARCHAR(32),  -- gold/item/title/exp
    reward_value INTEGER,
    icon VARCHAR(64)
);

-- 玩家成就表
CREATE TABLE user_achievements (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    achievement_id VARCHAR(32) NOT NULL,
    progress INTEGER DEFAULT 0,
    is_completed INTEGER DEFAULT 0,
    completed_at DATETIME,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (achievement_id) REFERENCES achievements(id),
    UNIQUE(user_id, achievement_id)
);
```

### 接口设计

```go
// AchievementManager 成就管理器接口
type AchievementManager interface {
    // 检查成就进度
    CheckProgress(userID int, achievementID string) error
    
    // 完成成就
    CompleteAchievement(userID int, achievementID string) error
    
    // 获取成就列表
    GetAchievements(userID int) ([]UserAchievement, error)
    
    // 获取成就进度
    GetProgress(userID int, achievementID string) (int, error)
}
```

---

## 系统集成

### 依赖关系

```
交易系统 → 装备系统、经济系统
排名系统 → 角色系统、战斗系统、战斗统计系统
聊天系统 → 用户系统
成就系统 → 战斗系统、角色系统、战斗统计系统
```

### 数据流

```
玩家操作
    ↓
各系统处理
    ↓
数据更新
    ↓
排行榜/成就更新（异步）
    ↓
前端展示
```

---

## 总结

### 设计亮点

1. **轻量级设计**: 所有系统都设计为轻量级，不影响单机体验
2. **异步处理**: 大部分功能支持异步，无需实时交互
3. **可选参与**: 玩家可以选择是否参与联网功能
4. **单机优先**: 游戏核心体验为单机，联网功能为辅助

### 实施优先级

1. **交易系统**: 高优先级（装备系统核心）
2. **排名系统**: 中优先级（增加竞争性）
3. **聊天系统**: 中优先级（增加社交性）
4. **成就系统**: 低优先级（长期目标）

---

**文档版本**: v1.0  
**最后更新**: 2025年


