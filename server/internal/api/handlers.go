package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"text-wow/internal/auth"
	"text-wow/internal/game"
	"text-wow/internal/models"
	"text-wow/internal/repository"
	"text-wow/internal/service"

	"github.com/gin-gonic/gin"
)

// Handler API处理器
type Handler struct {
	userRepo     *repository.UserRepository
	charRepo     *repository.CharacterRepository
	gameRepo     *repository.GameRepository
	skillRepo    *repository.SkillRepository
	skillService *service.SkillService
}

// NewHandler 创建处理器
func NewHandler() *Handler {
	skillRepo := repository.NewSkillRepository()
	skillService := service.NewSkillService(skillRepo, repository.NewCharacterRepository())
	return &Handler{
		userRepo:     repository.NewUserRepository(),
		charRepo:     repository.NewCharacterRepository(),
		gameRepo:     repository.NewGameRepository(),
		skillRepo:    skillRepo,
		skillService: skillService,
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
		// 区分用户不存在和其他错误
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Error:   "user not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Error:   "failed to get user info",
			})
		}
		return
	}

	// 确保返回的用户数据有效
	if user == nil || user.ID == 0 {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "user not found",
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
	}

	// 计算基础属性 = 职业基础 + 种族加成
	char.Strength = class.BaseStrength + race.StrengthBase
	char.Agility = class.BaseAgility + race.AgilityBase
	char.Intellect = class.BaseIntellect + race.IntellectBase
	char.Stamina = class.BaseStamina + race.StaminaBase
	char.Spirit = class.BaseSpirit + race.SpiritBase

	// 计算HP和资源
	// 最大HP = 职业基础HP + 耐力×2
	char.MaxHP = class.BaseHP + char.Stamina*2
	char.HP = char.MaxHP

	// 战士的怒气最大值固定为100，初始值为0
	// 其他职业：最大MP = 职业基础MP + 精神×2
	if class.ResourceType == "rage" {
		char.MaxResource = 100
		char.Resource = 0
	} else if class.ResourceType == "energy" {
		// 能量职业（盗贼）固定100
		char.MaxResource = 100
		char.Resource = char.MaxResource
	} else {
		// 法力职业：最大MP = 基础MP + 精神×2
		char.MaxResource = class.BaseResource + char.Spirit*2
		char.Resource = char.MaxResource
	}

	// 计算物理和魔法攻击/防御
	// 物理攻击 = 力量×1.0 + 敏捷×0.2
	char.PhysicalAttack = int(float64(char.Strength)*1.0 + float64(char.Agility)*0.2)
	// 魔法攻击 = 智力×1.0 + 精神×0.2
	char.MagicAttack = int(float64(char.Intellect)*1.0 + float64(char.Spirit)*0.2)
	// 物理防御 = 力量×0.2 + 耐力×0.3
	char.PhysicalDefense = int(float64(char.Strength)*0.2 + float64(char.Stamina)*0.3)
	// 魔法防御 = 智力×0.2 + 精神×0.3
	char.MagicDefense = int(float64(char.Intellect)*0.2 + float64(char.Spirit)*0.3)

	// 计算暴击属性
	// 物理暴击率 = 基础5% + 敏捷/20
	char.PhysCritRate = 0.05 + float64(char.Agility)/20.0/100.0
	// 物理暴击伤害 = 150% + 力量×0.3%
	char.PhysCritDamage = 1.5 + float64(char.Strength)*0.3/100.0
	// 法术暴击率 = 基础5% + 精神/20
	char.SpellCritRate = 0.05 + float64(char.Spirit)/20.0/100.0
	// 法术暴击伤害 = 150% + 智力×0.3%
	char.SpellCritDamage = 1.5 + float64(char.Intellect)*0.3/100.0

	// 计算闪避率 = 基础5% + 敏捷/20
	char.DodgeRate = 0.05 + float64(char.Agility)/20.0/100.0

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
		char := characters[0]
		// 添加buff信息
		battleMgr := game.GetBattleManager()
		char.Buffs = battleMgr.GetCharacterBuffs(char.ID)

		c.JSON(http.StatusOK, models.APIResponse{
			Success: true,
			Data:    char,
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

// AllocateAttributePoints 分配属性点
func (h *Handler) AllocateAttributePoints(c *gin.Context) {
	userID := c.GetInt("userID")

	charID, err := strconv.Atoi(c.Param("characterId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "invalid character id",
		})
		return
	}

	// 载入角色并校验归属
	char, err := h.charRepo.GetByID(charID)
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

	var req struct {
		Strength  int `json:"strength"`
		Agility   int `json:"agility"`
		Intellect int `json:"intellect"`
		Stamina   int `json:"stamina"`
		Spirit    int `json:"spirit"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "invalid request body",
		})
		return
	}

	if req.Strength < 0 || req.Agility < 0 || req.Intellect < 0 || req.Stamina < 0 || req.Spirit < 0 {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "points must be non-negative",
		})
		return
	}

	// 校验点数
	toSpend := req.Strength + req.Agility + req.Intellect + req.Stamina + req.Spirit
	if toSpend == 0 {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "no points to allocate",
		})
		return
	}
	if toSpend > char.UnspentPoints {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "not enough unspent points",
		})
		return
	}

	// 应用加点
	char.Strength += req.Strength
	char.Agility += req.Agility
	char.Intellect += req.Intellect
	char.Stamina += req.Stamina
	char.Spirit += req.Spirit
	char.UnspentPoints -= toSpend

	// 重新计算派生属性（与创建时保持一致）
	class, err := h.gameRepo.GetClassByID(char.ClassID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "failed to load class config",
		})
		return
	}

	char.MaxHP = class.BaseHP + char.Stamina*2
	char.HP = char.MaxHP

	switch char.ResourceType {
	case "rage":
		char.MaxResource = 100
		// 怒气不回满
	case "energy":
		char.MaxResource = 100
		char.Resource = char.MaxResource
	default:
		char.MaxResource = class.BaseResource + char.Spirit*2
		char.Resource = char.MaxResource
	}

	char.PhysicalAttack = int(float64(char.Strength)*1.0 + float64(char.Agility)*0.2)
	char.MagicAttack = int(float64(char.Intellect)*1.0 + float64(char.Spirit)*0.2)
	char.PhysicalDefense = int(float64(char.Strength)*0.2 + float64(char.Stamina)*0.3)
	char.MagicDefense = int(float64(char.Intellect)*0.2 + float64(char.Spirit)*0.3)
	char.PhysCritRate = 0.05 + float64(char.Agility)/20.0/100.0
	char.PhysCritDamage = 1.5 + float64(char.Strength)*0.3/100.0
	char.SpellCritRate = 0.05 + float64(char.Spirit)/20.0/100.0
	char.SpellCritDamage = 1.5 + float64(char.Intellect)*0.3/100.0
	char.DodgeRate = 0.05 + float64(char.Agility)/20.0/100.0

	// 持久化
	if err := h.charRepo.Update(char); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "failed to save allocation",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: fmt.Sprintf("分配成功，剩余点数 %d", char.UnspentPoints),
		Data:    char,
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

// ═══════════════════════════════════════════════════════════
// 技能选择相关API
// ═══════════════════════════════════════════════════════════

// GetInitialSkillSelection 获取初始技能选择机会
func (h *Handler) GetInitialSkillSelection(c *gin.Context) {
	characterID := c.Param("characterId")
	if characterID == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "characterId is required",
		})
		return
	}

	var charID int
	if _, err := fmt.Sscanf(characterID, "%d", &charID); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "invalid characterId",
		})
		return
	}

	selection, err := h.skillService.GetInitialSkillSelection(charID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    selection,
	})
}

// SelectInitialSkills 选择初始技能
func (h *Handler) SelectInitialSkills(c *gin.Context) {
	var req models.InitialSkillSelectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 检查是否是长度验证错误
		if strings.Contains(err.Error(), "len") {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Error:   "必须选择2个初始技能",
			})
		} else {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Error:   err.Error(),
			})
		}
		return
	}

	// 验证角色所有权
	character, err := h.charRepo.GetByID(req.CharacterID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "character not found",
		})
		return
	}

	userID, _ := c.Get("userID")
	if character.UserID != userID {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Error:   "forbidden",
		})
		return
	}

	if err := h.skillService.SelectInitialSkills(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "初始技能选择成功",
	})
}

// GetSkillSelection 获取技能选择机会
func (h *Handler) GetSkillSelection(c *gin.Context) {
	characterID := c.Param("characterId")
	if characterID == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "characterId is required",
		})
		return
	}

	var charID int
	if _, err := fmt.Sscanf(characterID, "%d", &charID); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "invalid characterId",
		})
		return
	}

	// 验证角色所有权
	character, err := h.charRepo.GetByID(charID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "character not found",
		})
		return
	}

	userID, _ := c.Get("userID")
	if character.UserID != userID {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Error:   "forbidden",
		})
		return
	}

	selection, err := h.skillService.CheckSkillSelectionOpportunity(charID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if selection == nil {
		c.JSON(http.StatusOK, models.APIResponse{
			Success: true,
			Message: "当前没有技能选择机会",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    selection,
	})
}

// SelectSkill 选择技能（新技能或升级）
func (h *Handler) SelectSkill(c *gin.Context) {
	var req models.SkillSelectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// 验证角色所有权
	character, err := h.charRepo.GetByID(req.CharacterID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "character not found",
		})
		return
	}

	userID, _ := c.Get("userID")
	if character.UserID != userID {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Error:   "forbidden",
		})
		return
	}

	if err := h.skillService.SelectSkill(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "技能选择成功",
	})
}

// GetCharacterSkills 获取角色的所有技能
func (h *Handler) GetCharacterSkills(c *gin.Context) {
	characterID := c.Param("characterId")
	if characterID == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "characterId is required",
		})
		return
	}

	var charID int
	if _, err := fmt.Sscanf(characterID, "%d", &charID); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "invalid characterId",
		})
		return
	}

	// 验证角色所有权
	character, err := h.charRepo.GetByID(charID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "character not found",
		})
		return
	}

	userID, _ := c.Get("userID")
	if character.UserID != userID {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Error:   "forbidden",
		})
		return
	}

	skills, err := h.skillService.GetCharacterAllSkills(charID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    skills,
	})
}
