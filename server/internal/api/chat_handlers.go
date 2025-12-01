package api

import (
	"net/http"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"text-wow/internal/repository"

	"github.com/gin-gonic/gin"
)

// ChatHandler 聊天处理器
type ChatHandler struct {
	chatRepo *repository.ChatRepository
	charRepo *repository.CharacterRepository
	userRepo *repository.UserRepository
	
	// 简单的刷屏检测 (生产环境应使用Redis)
	lastMessages map[int]time.Time
}

// NewChatHandler 创建聊天处理器
func NewChatHandler() *ChatHandler {
	return &ChatHandler{
		chatRepo:     repository.NewChatRepository(),
		charRepo:     repository.NewCharacterRepository(),
		userRepo:     repository.NewUserRepository(),
		lastMessages: make(map[int]time.Time),
	}
}

// 敏感词列表 (简化版，生产环境应使用更完整的词库)
var sensitiveWords = []string{
	// 这里添加敏感词...
}

// 敏感词正则
var sensitiveRegex *regexp.Regexp

func init() {
	if len(sensitiveWords) > 0 {
		pattern := strings.Join(sensitiveWords, "|")
		sensitiveRegex = regexp.MustCompile("(?i)" + pattern)
	}
}

// 消息请求
type SendMessageRequest struct {
	Channel  string `json:"channel" binding:"required"`
	Content  string `json:"content" binding:"required,min=1,max=200"`
	ZoneID   string `json:"zoneId,omitempty"`
	Receiver string `json:"receiver,omitempty"` // 私聊目标角色名
}

// ═══════════════════════════════════════════════════════════
// 聊天 API
// ═══════════════════════════════════════════════════════════

// SendMessage 发送消息
func (h *ChatHandler) SendMessage(c *gin.Context) {
	userID := c.GetInt("userID")

	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid request: " + err.Error(),
		})
		return
	}

	// 验证频道
	validChannels := map[string]bool{
		"world": true, "zone": true, "trade": true,
		"lfg": true, "whisper": true,
	}
	if !validChannels[req.Channel] {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid channel",
		})
		return
	}

	// 刷屏检测 (3秒内不能发送相同消息)
	if lastTime, ok := h.lastMessages[userID]; ok {
		if time.Since(lastTime) < 3*time.Second {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error":   "please wait before sending another message",
			})
			return
		}
	}

	// 获取发送者信息
	chars, err := h.charRepo.GetActiveByUserID(userID)
	if err != nil || len(chars) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "no active character",
		})
		return
	}
	sender := chars[0] // 使用第一个活跃角色

	// 内容过滤
	content := filterContent(req.Content)
	if utf8.RuneCountInString(content) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "message content is empty after filtering",
		})
		return
	}

	msg := &repository.ChatMessage{
		Channel:     req.Channel,
		Faction:     sender.Faction,
		ZoneID:      req.ZoneID,
		SenderID:    userID,
		SenderName:  sender.Name,
		SenderClass: sender.ClassID,
		Content:     content,
	}

	// 私聊处理
	if req.Channel == "whisper" {
		if req.Receiver == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "receiver is required for whisper",
			})
			return
		}

		// 查找接收者
		receiver, err := h.chatRepo.GetUserByName(req.Receiver)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "player not found",
			})
			return
		}

		// 检查阵营 (不能跨阵营私聊)
		if receiver.Faction != sender.Faction {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "cannot whisper to enemy faction",
			})
			return
		}

		// 检查是否被屏蔽
		blocked, _ := h.chatRepo.IsBlocked(receiver.UserID, userID)
		if blocked {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "this player has blocked you",
			})
			return
		}

		msg.ReceiverID = receiver.UserID
	}

	// 保存消息
	savedMsg, err := h.chatRepo.SendMessage(msg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "failed to send message",
		})
		return
	}

	// 更新刷屏检测
	h.lastMessages[userID] = time.Now()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    savedMsg,
	})
}

// GetMessages 获取消息
func (h *ChatHandler) GetMessages(c *gin.Context) {
	userID := c.GetInt("userID")

	channel := c.Query("channel")
	if channel == "" {
		channel = "world"
	}

	// 获取用户阵营
	chars, err := h.charRepo.GetActiveByUserID(userID)
	if err != nil || len(chars) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "no active character",
		})
		return
	}
	faction := chars[0].Faction

	// 获取屏蔽列表
	blockedUsers, _ := h.chatRepo.GetBlockedUsers(userID)
	blockedMap := make(map[int]bool)
	for _, id := range blockedUsers {
		blockedMap[id] = true
	}

	var messages []repository.ChatMessage

	switch channel {
	case "world", "trade", "lfg":
		messages, err = h.chatRepo.GetChannelMessages(channel, faction, "", 50, 0)
	case "zone":
		zoneID := c.Query("zoneId")
		messages, err = h.chatRepo.GetChannelMessages(channel, faction, zoneID, 50, 0)
	case "whisper":
		otherName := c.Query("with")
		if otherName != "" {
			other, err := h.chatRepo.GetUserByName(otherName)
			if err == nil {
				messages, _ = h.chatRepo.GetWhisperMessages(userID, other.UserID, 50, 0)
			}
		}
	case "recent":
		messages, err = h.chatRepo.GetRecentMessages(faction, 100)
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid channel",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "failed to get messages",
		})
		return
	}

	// 过滤屏蔽用户的消息
	filtered := make([]repository.ChatMessage, 0, len(messages))
	for _, msg := range messages {
		if !blockedMap[msg.SenderID] {
			filtered = append(filtered, msg)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    filtered,
	})
}

// GetOnlineUsers 获取在线用户
func (h *ChatHandler) GetOnlineUsers(c *gin.Context) {
	userID := c.GetInt("userID")

	// 获取用户阵营
	chars, err := h.charRepo.GetActiveByUserID(userID)
	if err != nil || len(chars) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "no active character",
		})
		return
	}
	faction := chars[0].Faction

	// 只能看到同阵营的在线玩家
	users, err := h.chatRepo.GetOnlineUsers(faction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "failed to get online users",
		})
		return
	}

	// 获取在线人数
	allianceCount, _ := h.chatRepo.GetOnlineCount("alliance")
	hordeCount, _ := h.chatRepo.GetOnlineCount("horde")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"users":         users,
			"factionCount":  len(users),
			"allianceCount": allianceCount,
			"hordeCount":    hordeCount,
		},
	})
}

// BlockUser 屏蔽用户
func (h *ChatHandler) BlockUser(c *gin.Context) {
	userID := c.GetInt("userID")

	var req struct {
		PlayerName string `json:"playerName" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid request",
		})
		return
	}

	// 查找目标玩家
	target, err := h.chatRepo.GetUserByName(req.PlayerName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "player not found",
		})
		return
	}

	if target.UserID == userID {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "cannot block yourself",
		})
		return
	}

	if err := h.chatRepo.BlockUser(userID, target.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "failed to block user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "player blocked",
	})
}

// UnblockUser 取消屏蔽
func (h *ChatHandler) UnblockUser(c *gin.Context) {
	userID := c.GetInt("userID")

	var req struct {
		PlayerName string `json:"playerName" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid request",
		})
		return
	}

	// 查找目标玩家
	target, err := h.chatRepo.GetUserByName(req.PlayerName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "player not found",
		})
		return
	}

	if err := h.chatRepo.UnblockUser(userID, target.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "failed to unblock user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "player unblocked",
	})
}

// SetOnline 设置在线状态
func (h *ChatHandler) SetOnline(c *gin.Context) {
	userID := c.GetInt("userID")

	var req struct {
		CharacterID int    `json:"characterId" binding:"required"`
		ZoneID      string `json:"zoneId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid request",
		})
		return
	}

	// 获取角色信息
	char, err := h.charRepo.GetByID(req.CharacterID)
	if err != nil || char.UserID != userID {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "character not found",
		})
		return
	}

	zoneID := req.ZoneID
	if zoneID == "" {
		// 从用户获取当前区域
		user, _ := h.userRepo.GetByID(userID)
		if user != nil {
			zoneID = user.CurrentZoneID
		}
	}

	if err := h.chatRepo.SetOnlineStatus(userID, char.ID, char.Name, char.Faction, zoneID, true); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "failed to set online status",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

// SetOffline 设置离线状态
func (h *ChatHandler) SetOffline(c *gin.Context) {
	userID := c.GetInt("userID")

	if err := h.chatRepo.SetOnlineStatus(userID, 0, "", "", "", false); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "failed to set offline status",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

// Heartbeat 心跳 (保持在线状态)
func (h *ChatHandler) Heartbeat(c *gin.Context) {
	userID := c.GetInt("userID")
	h.chatRepo.UpdateLastActive(userID)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// ═══════════════════════════════════════════════════════════
// 辅助函数
// ═══════════════════════════════════════════════════════════

// filterContent 过滤消息内容
func filterContent(content string) string {
	// 去除首尾空白
	content = strings.TrimSpace(content)
	
	// 敏感词替换
	if sensitiveRegex != nil {
		content = sensitiveRegex.ReplaceAllString(content, "***")
	}

	return content
}

// GarbleMessage 将消息转换为"敌方语言" (跨阵营时使用)
func GarbleMessage(content string, fromFaction string) string {
	// 兽人语音节
	orcSyllables := []string{"lok", "tar", "ogar", "zug", "kek", "gul", "mok", "throm", "ka"}
	// 通用语音节  
	commonSyllables := []string{"bur", "gol", "len", "mos", "nud", "ras", "ver", "ash", "thu"}

	syllables := orcSyllables
	if fromFaction == "horde" {
		syllables = commonSyllables
	}

	// 简单的乱码化：根据原文长度生成相应长度的"语言"
	words := strings.Fields(content)
	var result []string
	
	for _, word := range words {
		runeCount := utf8.RuneCountInString(word)
		// 每2-3个字符一个音节
		syllableCount := (runeCount + 2) / 3
		if syllableCount < 1 {
			syllableCount = 1
		}
		
		var newWord []string
		for i := 0; i < syllableCount; i++ {
			idx := (len(word) + i) % len(syllables)
			newWord = append(newWord, syllables[idx])
		}
		result = append(result, strings.Join(newWord, "'"))
	}

	return strings.Join(result, " ")
}






