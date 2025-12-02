package api

import (
	"database/sql"
	"net/http"
	"strings"

	"text-wow/internal/auth"
	"text-wow/internal/models"
	"text-wow/internal/repository"

	"github.com/gin-gonic/gin"
)

// Handler API处理器
type Handler struct {
	userRepo *repository.UserRepository
	charRepo *repository.CharacterRepository
	gameRepo *repository.GameRepository
}

// NewHandler 创建处理器
func NewHandler() *Handler {
	return &Handler{
		userRepo: repository.NewUserRepository(),
		charRepo: repository.NewCharacterRepository(),
		gameRepo: repository.NewGameRepository(),
	}
}

// ═══════════════════════════════════════════════════════════
// 认证中间件
// ═══════════════════════════════════════════════════════════

// AuthMiddleware JWT认证中间件
func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Header获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Error:   "missing authorization header",
			})
			c.Abort()
			return
		}

		// 解析 Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Error:   "invalid authorization format",
			})
			c.Abort()
			return
		}

		// 验证token
		claims, err := auth.ValidateToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Error:   "invalid or expired token",
			})
			c.Abort()
			return
		}

		// 将用户信息存入context
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}

// ═══════════════════════════════════════════════════════════
// 认证相关API
// ═══════════════════════════════════════════════════════════

// Register 用户注册
func (h *Handler) Register(c *gin.Context) {
	var req models.UserRegister
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "invalid request: " + err.Error(),
		})
		return
	}

	// 检查用户名是否存在
	exists, err := h.userRepo.UsernameExists(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "database error",
		})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, models.APIResponse{
			Success: false,
			Error:   "username already exists",
		})
		return
	}

	// 加密密码
	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "failed to hash password",
		})
		return
	}

	// 创建用户
	user, err := h.userRepo.Create(req.Username, passwordHash, req.Email)
	if err != nil {
		// 输出错误日志以便调试
		println("Error creating user:", err.Error())
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "failed to create user: " + err.Error(),
		})
		return
	}

	// 生成token
	token, err := auth.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "registration successful",
		Data: models.AuthResponse{
			Token: token,
			User:  *user,
		},
	})
}

// Login 用户登录
func (h *Handler) Login(c *gin.Context) {
	var req models.UserCredentials
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "invalid request: " + err.Error(),
		})
		return
	}

	// 获取用户密码哈希
	userID, passwordHash, err := h.userRepo.GetPasswordHash(req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Error:   "invalid username or password",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "database error",
		})
		return
	}

	// 验证密码
	if !auth.CheckPassword(req.Password, passwordHash) {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "invalid username or password",
		})
		return
	}

	// 更新最后登录时间
	h.userRepo.UpdateLastLogin(userID)

	// 获取用户信息
	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "failed to get user info",
		})
		return
	}

	// 生成token
	token, err := auth.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "login successful",
		Data: models.AuthResponse{
			Token: token,
			User:  *user,
		},
	})
}

// GetCurrentUser 获取当前用户信息
func (h *Handler) GetCurrentUser(c *gin.Context) {
	userID := c.GetInt("userID")

	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "failed to get user info",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    user,
	})
}

// ═══════════════════════════════════════════════════════════
// 角色相关API
// ═══════════════════════════════════════════════════════════

// CreateCharacter 创建角色
func (h *Handler) CreateCharacter(c *gin.Context) {
	userID := c.GetInt("userID")

	var req models.CharacterCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "invalid request: " + err.Error(),
		})
		return
	}

	// 获取用户信息（检查槽位）
	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "failed to get user info",
		})
		return
	}

	// 检查角色数量
	count, err := h.charRepo.CountByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "failed to check character count",
		})
		return
	}
	if count >= user.UnlockedSlots {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "no available character slots",
		})
		return
	}

	// 检查角色名是否存在
	exists, err := h.charRepo.NameExists(req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "database error",
		})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, models.APIResponse{
			Success: false,
			Error:   "character name already exists",
		})
		return
	}

	// 获取种族和职业信息
	race, err := h.gameRepo.GetRaceByID(req.RaceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "invalid race",
		})
		return
	}

	class, err := h.gameRepo.GetClassByID(req.ClassID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "invalid class",
		})
		return
	}

	// 获取下一个可用槽位
	slot, err := h.charRepo.GetNextSlot(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "failed to get slot",
		})
		return
	}

	// 计算初始属性
	char := &models.Character{
		UserID:       userID,
		Name:         req.Name,
		RaceID:       req.RaceID,
		ClassID:      req.ClassID,
		Faction:      race.Faction,
		TeamSlot:     slot,
		IsActive:     true,
		IsDead:       false,
		Level:        1,
		Exp:          0,
		ExpToNext:    100,
		ResourceType: class.ResourceType,
		CritRate:     0.05,
		CritDamage:   1.5,
	}

	// 计算基础属性 = 职业基础 + 种族加成
	char.Strength = class.BaseStrength + race.StrengthBase
	char.Agility = class.BaseAgility + race.AgilityBase
	char.Intellect = class.BaseIntellect + race.IntellectBase
	char.Stamina = class.BaseStamina + race.StaminaBase
	char.Spirit = class.BaseSpirit + race.SpiritBase

	// 计算HP和资源
	char.MaxHP = class.BaseHP + char.Stamina*2
	char.HP = char.MaxHP
	
	// 战士的怒气最大值固定为100，初始值为0
	if class.ResourceType == "rage" {
		char.MaxResource = 100
		char.Resource = 0
	} else {
		char.MaxResource = class.BaseResource
		char.Resource = char.MaxResource
	}

	// 计算攻击和防御
	char.Attack = char.Strength / 2
	char.Defense = char.Stamina / 3

	// 创建角色
	char, err = h.charRepo.Create(char)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "failed to create character: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "character created successfully",
		Data:    char,
	})
}

// GetCharacters 获取用户的所有角色
func (h *Handler) GetCharacters(c *gin.Context) {
	userID := c.GetInt("userID")

	characters, err := h.charRepo.GetByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "failed to get characters",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    characters,
	})
}

// GetCharacter 获取第一个角色（用于游戏界面显示）
// 即使角色死亡也会返回，以便显示复活倒计时
func (h *Handler) GetCharacter(c *gin.Context) {
	userID := c.GetInt("userID")

	// 获取所有角色（所有角色都参与战斗）
	characters, err := h.charRepo.GetByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "failed to get character",
		})
		return
	}

	// 返回第一个角色（即使死亡也会返回，以便显示复活状态）
	if len(characters) > 0 {
		c.JSON(http.StatusOK, models.APIResponse{
			Success: true,
			Data:    characters[0],
		})
		return
	}
	
	c.JSON(http.StatusNotFound, models.APIResponse{
		Success: false,
		Error:   "no character found",
	})
}

// GetTeam 获取小队信息
func (h *Handler) GetTeam(c *gin.Context) {
	userID := c.GetInt("userID")

	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "failed to get user info",
		})
		return
	}

	// 获取所有角色（所有角色都参与战斗）
	characters, err := h.charRepo.GetByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "failed to get team",
		})
		return
	}

	team := &models.Team{
		UserID:        userID,
		MaxSize:       user.MaxTeamSize,
		UnlockedSlots: user.UnlockedSlots,
		Characters:    characters,
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    team,
	})
}

// SetCharacterActive 设置角色激活状态
func (h *Handler) SetCharacterActive(c *gin.Context) {
	userID := c.GetInt("userID")

	var req struct {
		CharacterID int  `json:"characterId" binding:"required"`
		Active      bool `json:"active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "invalid request",
		})
		return
	}

	// 验证角色属于该用户
	char, err := h.charRepo.GetByID(req.CharacterID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "character not found",
		})
		return
	}
	if char.UserID != userID {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Error:   "not your character",
		})
		return
	}

	// 更新状态
	if err := h.charRepo.SetActive(req.CharacterID, req.Active); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "failed to update character",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "character updated",
	})
}

// ═══════════════════════════════════════════════════════════
// 游戏配置API
// ═══════════════════════════════════════════════════════════

// GetRaces 获取种族列表
func (h *Handler) GetRaces(c *gin.Context) {
	races, err := h.gameRepo.GetRaces()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "failed to get races",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    races,
	})
}

// GetClasses 获取职业列表
func (h *Handler) GetClasses(c *gin.Context) {
	classes, err := h.gameRepo.GetClasses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "failed to get classes",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    classes,
	})
}

// GetZones 获取区域列表
func (h *Handler) GetZones(c *gin.Context) {
	zones, err := h.gameRepo.GetZones()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "failed to get zones",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    zones,
	})
}
