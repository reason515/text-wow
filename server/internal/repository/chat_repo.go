package repository

import (
	"database/sql"
	"time"

	"text-wow/internal/database"
)

// ChatMessage 聊天消息
type ChatMessage struct {
	ID          int       `json:"id"`
	Channel     string    `json:"channel"`
	Faction     string    `json:"faction,omitempty"`
	ZoneID      string    `json:"zoneId,omitempty"`
	SenderID    int       `json:"senderId"`
	SenderName  string    `json:"senderName"`
	SenderClass string    `json:"senderClass,omitempty"`
	ReceiverID  int       `json:"receiverId,omitempty"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"createdAt"`
}

// OnlineUser 在线用户
type OnlineUser struct {
	UserID        int       `json:"userId"`
	CharacterID   int       `json:"characterId,omitempty"`
	CharacterName string    `json:"characterName"`
	Faction       string    `json:"faction"`
	ZoneID        string    `json:"zoneId,omitempty"`
	LastActive    time.Time `json:"lastActive"`
	IsOnline      bool      `json:"isOnline"`
}

// ChatRepository 聊天数据仓库
type ChatRepository struct{}

// NewChatRepository 创建聊天仓库
func NewChatRepository() *ChatRepository {
	return &ChatRepository{}
}

// SendMessage 发送消息
func (r *ChatRepository) SendMessage(msg *ChatMessage) (*ChatMessage, error) {
	result, err := database.DB.Exec(`
		INSERT INTO chat_messages (channel, faction, zone_id, sender_id, sender_name, sender_class, receiver_id, content, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		msg.Channel, nullString(msg.Faction), nullString(msg.ZoneID),
		msg.SenderID, msg.SenderName, nullString(msg.SenderClass),
		nullInt(msg.ReceiverID), msg.Content, time.Now(),
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	msg.ID = int(id)
	msg.CreatedAt = time.Now()
	return msg, nil
}

// GetChannelMessages 获取频道消息
func (r *ChatRepository) GetChannelMessages(channel, faction, zoneID string, limit, offset int) ([]ChatMessage, error) {
	query := `
		SELECT id, channel, COALESCE(faction, ''), COALESCE(zone_id, ''),
		       sender_id, sender_name, COALESCE(sender_class, ''),
		       COALESCE(receiver_id, 0), content, created_at
		FROM chat_messages
		WHERE channel = ?`
	
	args := []interface{}{channel}

	if faction != "" {
		query += " AND faction = ?"
		args = append(args, faction)
	}

	if zoneID != "" {
		query += " AND zone_id = ?"
		args = append(args, zoneID)
	}

	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanMessages(rows)
}

// GetWhisperMessages 获取私聊消息
func (r *ChatRepository) GetWhisperMessages(userID, otherUserID int, limit, offset int) ([]ChatMessage, error) {
	rows, err := database.DB.Query(`
		SELECT id, channel, COALESCE(faction, ''), COALESCE(zone_id, ''),
		       sender_id, sender_name, COALESCE(sender_class, ''),
		       COALESCE(receiver_id, 0), content, created_at
		FROM chat_messages
		WHERE channel = 'whisper'
		  AND ((sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?))
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`,
		userID, otherUserID, otherUserID, userID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanMessages(rows)
}

// GetRecentMessages 获取最近消息 (用于刚上线时加载)
func (r *ChatRepository) GetRecentMessages(faction string, limit int) ([]ChatMessage, error) {
	rows, err := database.DB.Query(`
		SELECT id, channel, COALESCE(faction, ''), COALESCE(zone_id, ''),
		       sender_id, sender_name, COALESCE(sender_class, ''),
		       COALESCE(receiver_id, 0), content, created_at
		FROM chat_messages
		WHERE (faction = ? OR faction IS NULL OR channel = 'system')
		  AND channel != 'whisper'
		ORDER BY created_at DESC
		LIMIT ?`,
		faction, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanMessages(rows)
}

// BlockUser 屏蔽用户
func (r *ChatRepository) BlockUser(userID, blockedID int) error {
	_, err := database.DB.Exec(`
		INSERT OR IGNORE INTO chat_blocks (user_id, blocked_id, created_at)
		VALUES (?, ?, ?)`,
		userID, blockedID, time.Now(),
	)
	return err
}

// UnblockUser 取消屏蔽
func (r *ChatRepository) UnblockUser(userID, blockedID int) error {
	_, err := database.DB.Exec(`
		DELETE FROM chat_blocks WHERE user_id = ? AND blocked_id = ?`,
		userID, blockedID,
	)
	return err
}

// GetBlockedUsers 获取屏蔽列表
func (r *ChatRepository) GetBlockedUsers(userID int) ([]int, error) {
	rows, err := database.DB.Query(`
		SELECT blocked_id FROM chat_blocks WHERE user_id = ?`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blocked []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		blocked = append(blocked, id)
	}
	return blocked, nil
}

// IsBlocked 检查是否被屏蔽
func (r *ChatRepository) IsBlocked(userID, blockedID int) (bool, error) {
	var count int
	err := database.DB.QueryRow(`
		SELECT COUNT(*) FROM chat_blocks WHERE user_id = ? AND blocked_id = ?`,
		userID, blockedID,
	).Scan(&count)
	return count > 0, err
}

// SetOnlineStatus 设置在线状态
func (r *ChatRepository) SetOnlineStatus(userID int, charID int, charName, faction, zoneID string, online bool) error {
	_, err := database.DB.Exec(`
		INSERT INTO user_online_status (user_id, character_id, character_name, faction, zone_id, last_active, is_online)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(user_id) DO UPDATE SET
			character_id = excluded.character_id,
			character_name = excluded.character_name,
			faction = excluded.faction,
			zone_id = excluded.zone_id,
			last_active = excluded.last_active,
			is_online = excluded.is_online`,
		userID, charID, charName, faction, zoneID, time.Now(), boolToInt(online),
	)
	return err
}

// UpdateLastActive 更新最后活跃时间
func (r *ChatRepository) UpdateLastActive(userID int) error {
	_, err := database.DB.Exec(`
		UPDATE user_online_status SET last_active = ? WHERE user_id = ?`,
		time.Now(), userID,
	)
	return err
}

// GetOnlineUsers 获取在线用户
func (r *ChatRepository) GetOnlineUsers(faction string) ([]OnlineUser, error) {
	query := `
		SELECT user_id, COALESCE(character_id, 0), COALESCE(character_name, ''),
		       COALESCE(faction, ''), COALESCE(zone_id, ''), last_active, is_online
		FROM user_online_status
		WHERE is_online = 1`
	
	args := []interface{}{}
	if faction != "" {
		query += " AND faction = ?"
		args = append(args, faction)
	}
	query += " ORDER BY character_name"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []OnlineUser
	for rows.Next() {
		var u OnlineUser
		var isOnline int
		err := rows.Scan(&u.UserID, &u.CharacterID, &u.CharacterName,
			&u.Faction, &u.ZoneID, &u.LastActive, &isOnline)
		if err != nil {
			return nil, err
		}
		u.IsOnline = isOnline == 1
		users = append(users, u)
	}
	return users, nil
}

// GetOnlineCount 获取在线人数
func (r *ChatRepository) GetOnlineCount(faction string) (int, error) {
	query := "SELECT COUNT(*) FROM user_online_status WHERE is_online = 1"
	args := []interface{}{}
	
	if faction != "" {
		query += " AND faction = ?"
		args = append(args, faction)
	}

	var count int
	err := database.DB.QueryRow(query, args...).Scan(&count)
	return count, err
}

// CleanupInactiveUsers 清理不活跃用户 (超过5分钟无活动)
func (r *ChatRepository) CleanupInactiveUsers() error {
	threshold := time.Now().Add(-5 * time.Minute)
	_, err := database.DB.Exec(`
		UPDATE user_online_status SET is_online = 0
		WHERE is_online = 1 AND last_active < ?`,
		threshold,
	)
	return err
}

// GetUserByName 根据角色名获取用户
func (r *ChatRepository) GetUserByName(name string) (*OnlineUser, error) {
	var u OnlineUser
	var isOnline int
	err := database.DB.QueryRow(`
		SELECT user_id, COALESCE(character_id, 0), COALESCE(character_name, ''),
		       COALESCE(faction, ''), COALESCE(zone_id, ''), last_active, is_online
		FROM user_online_status
		WHERE character_name = ?`, name,
	).Scan(&u.UserID, &u.CharacterID, &u.CharacterName,
		&u.Faction, &u.ZoneID, &u.LastActive, &isOnline)
	if err != nil {
		return nil, err
	}
	u.IsOnline = isOnline == 1
	return &u, nil
}

// 辅助函数
func scanMessages(rows *sql.Rows) ([]ChatMessage, error) {
	var messages []ChatMessage
	for rows.Next() {
		var msg ChatMessage
		err := rows.Scan(
			&msg.ID, &msg.Channel, &msg.Faction, &msg.ZoneID,
			&msg.SenderID, &msg.SenderName, &msg.SenderClass,
			&msg.ReceiverID, &msg.Content, &msg.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	// 反转顺序，使最新消息在最后
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}
	return messages, nil
}

func nullString(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

func nullInt(i int) interface{} {
	if i == 0 {
		return nil
	}
	return i
}

