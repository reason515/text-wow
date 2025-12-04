package game

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"text-wow/internal/models"
	"text-wow/internal/repository"
)

// BattleManager 战斗管理器 - 管理所有用户的战斗状态
type BattleManager struct {
	mu                  sync.RWMutex
	sessions            map[int]*BattleSession // key: userID
	gameRepo            *repository.GameRepository
	charRepo            *repository.CharacterRepository
	skillManager        *SkillManager
	buffManager         *BuffManager
	passiveSkillManager *PassiveSkillManager
}

// BattleSession 用户战斗会话
type BattleSession struct {
	UserID             int
	IsRunning          bool
	CurrentZone        *models.Zone
	CurrentEnemy       *models.Monster   // 保留用于向后兼容
	CurrentEnemies     []*models.Monster // 多个敌人支持
	BattleLogs         []models.BattleLog
	BattleCount        int
	SessionKills       int
	SessionGold        int
	SessionExp         int
	StartedAt          time.Time
	LastTick           time.Time
	IsResting          bool       // 是否在休息
	RestUntil          *time.Time // 休息结束时间
	RestStartedAt      *time.Time // 休息开始时间
	LastRestTick       *time.Time // 上次恢复处理的时间
	RestSpeed          float64    // 恢复速度倍率
	CurrentBattleExp   int        // 本场战斗获得的经验
	CurrentBattleGold  int        // 本场战斗获得的金币
	CurrentBattleKills int        // 本场战斗击杀数
	CurrentTurnIndex   int        // 回合控制：-1=玩家回合，>=0=敌人索引
	JustEncountered    bool       // 刚遭遇敌人，需要等待1个tick再开始战斗
}

// NewBattleManager 创建战斗管理器
func NewBattleManager() *BattleManager {
	return &BattleManager{
		sessions:            make(map[int]*BattleSession),
		gameRepo:            repository.NewGameRepository(),
		charRepo:            repository.NewCharacterRepository(),
		skillManager:        NewSkillManager(),
		buffManager:         NewBuffManager(),
		passiveSkillManager: NewPassiveSkillManager(),
	}
}

// 全局战斗管理器实例
var battleManager *BattleManager
var once sync.Once

// GetBattleManager 获取战斗管理器单例
func GetBattleManager() *BattleManager {
	once.Do(func() {
		battleManager = NewBattleManager()
	})
	return battleManager
}

// GetOrCreateSession 获取或创建战斗会话
func (m *BattleManager) GetOrCreateSession(userID int) *BattleSession {
	m.mu.Lock()
	defer m.mu.Unlock()

	if session, exists := m.sessions[userID]; exists {
		return session
	}

	session := &BattleSession{
		UserID:           userID,
		BattleLogs:       make([]models.BattleLog, 0),
		StartedAt:        time.Now(),
		CurrentEnemies:   make([]*models.Monster, 0),
		CurrentTurnIndex: -1, // 初始化为玩家回合
		RestSpeed:        1.0,
	}
	m.sessions[userID] = session
	return session
}

// GetSession 获取战斗会话
func (m *BattleManager) GetSession(userID int) *BattleSession {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.sessions[userID]
}

// ToggleBattle 切换战斗状态
func (m *BattleManager) ToggleBattle(userID int) (bool, error) {
	session := m.GetOrCreateSession(userID)

	m.mu.Lock()
	defer m.mu.Unlock()

	session.IsRunning = !session.IsRunning
	session.LastTick = time.Now()

	if session.IsRunning {
		// 如果没有设置区域，设置默认区域
		if session.CurrentZone == nil {
			zone, err := m.gameRepo.GetZoneByID("elwynn")
			if err == nil {
				session.CurrentZone = zone
			}
		}
		session.CurrentTurnIndex = -1 // 重置为玩家回合
		m.addLog(session, "system", ">> 开始自动战斗...", "#33ff33")
	} else {
		m.addLog(session, "system", ">> 暂停自动战斗", "#ffff00")
	}

	return session.IsRunning, nil
}

// StartBattle 开始战斗
func (m *BattleManager) StartBattle(userID int) (bool, error) {
	session := m.GetOrCreateSession(userID)

	m.mu.Lock()
	defer m.mu.Unlock()

	if session.IsRunning {
		return true, nil
	}

	session.IsRunning = true
	session.LastTick = time.Now()
	session.CurrentTurnIndex = -1 // 重置为玩家回合

	// 设置默认区域
	if session.CurrentZone == nil {
		zone, err := m.gameRepo.GetZoneByID("elwynn")
		if err == nil {
			session.CurrentZone = zone
		}
	}

	m.addLog(session, "system", ">> 开始自动战斗...", "#33ff33")
	return true, nil
}

// StopBattle 停止战斗
func (m *BattleManager) StopBattle(userID int) error {
	session := m.GetSession(userID)
	if session == nil {
		return nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	session.IsRunning = false
	m.addLog(session, "system", ">> 暂停自动战斗", "#ffff00")
	return nil
}

// ExecuteBattleTick 执行战斗回合（回合制：每tick只执行一个动作）
func (m *BattleManager) ExecuteBattleTick(userID int, characters []*models.Character) (*BattleTickResult, error) {
	session := m.GetOrCreateSession(userID)

	m.mu.Lock()
	defer m.mu.Unlock()

	// 如果没有角色，返回nil
	if len(characters) == 0 {
		return nil, nil
	}

	// 使用第一个角色进行战斗
	char := characters[0]

	// 确保战士的怒气上限为100（每次tick都检查，防止被覆盖）
	if char.ResourceType == "rage" {
		char.MaxResource = 100
	}

	// 加载角色的技能（如果还没有加载）
	if err := m.skillManager.LoadCharacterSkills(char.ID); err != nil {
		// 如果加载失败，记录日志但不中断战斗
		m.addLog(session, "system", fmt.Sprintf("警告：无法加载角色技能: %v", err), "#ffaa00")
	}

	// 加载角色的被动技能（如果还没有加载）
	if err := m.passiveSkillManager.LoadCharacterPassiveSkills(char.ID); err != nil {
		// 如果加载失败，记录日志但不中断战斗
		m.addLog(session, "system", fmt.Sprintf("警告：无法加载角色被动技能: %v", err), "#ffaa00")
	}

	// 如果战斗未运行且不在休息状态，检查是否需要返回角色数据
	// 如果角色刚复活（之前死亡但现在不死亡），需要返回一次数据让前端更新
	if !session.IsRunning && !session.IsResting {
		// 从数据库重新加载角色数据以确保状态正确
		updatedChar, err := m.charRepo.GetByID(char.ID)
		if err == nil && updatedChar != nil {
			char = updatedChar
			// 确保战士的怒气上限为100
			if char.ResourceType == "rage" {
				char.MaxResource = 100
			}
			// 如果角色已经复活（之前死亡但现在不死亡），返回角色数据
			if !char.IsDead {
				// 返回角色数据，让前端知道角色已经复活
				return &BattleTickResult{
					Character:    char,
					Enemy:        nil,
					Enemies:      nil,
					Logs:         []models.BattleLog{},
					IsRunning:    false,
					IsResting:    false,
					RestUntil:    nil,
					SessionKills: session.SessionKills,
					SessionGold:  session.SessionGold,
					SessionExp:   session.SessionExp,
					BattleCount:  session.BattleCount,
				}, nil
			}
		}
		// 如果无法获取角色数据或角色仍然死亡，返回nil
		return nil, nil
	}

	session.LastTick = time.Now()
	logs := make([]models.BattleLog, 0)

	// 检查角色是否死亡且还没到复活时间
	now := time.Now()
	if char.IsDead && char.ReviveAt != nil && now.Before(*char.ReviveAt) {
		// 角色死亡但还没到复活时间，进入休息状态
		if !session.IsResting {
			// 计算休息时间（复活时间 + 恢复时间）
			reviveRemaining := char.ReviveAt.Sub(now)
			recoveryTime := 25 * time.Second // 恢复一半HP需要的时间
			restDuration := reviveRemaining + recoveryTime
			restUntil := now.Add(restDuration)
			session.IsResting = true
			session.RestUntil = &restUntil
			session.RestStartedAt = &now
			session.LastRestTick = &now
			session.RestSpeed = 1.0
			// 保持 isRunning = true，这样按钮会显示"停止挂机"，休息状态可以自动处理

			remainingSeconds := int(reviveRemaining.Seconds()) + 1
			m.addLog(session, "death", fmt.Sprintf("%s 正在复活中... (剩余 %d 秒)", char.Name, remainingSeconds), "#ff0000")
			logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
		}
	}

	// 如果正在休息，处理休息
	if session.IsResting && session.RestUntil != nil {
		initialHP := char.HP
		initialMP := char.Resource
		now := time.Now()
		m.processRest(session, char)

		// 更新LastTick，用于下次计算时间差
		session.LastTick = now

		if !session.IsResting {
			// 休息结束，保存角色数据
			m.charRepo.UpdateAfterBattle(char.ID, char.HP, char.Resource, char.Exp, char.Level,
				char.ExpToNext, char.MaxHP, char.MaxResource, char.Attack, char.Defense,
				char.Strength, char.Agility, char.Stamina, char.TotalKills)

			// 休息结束后，确保返回角色数据，让前端知道休息已结束
			// 从数据库重新加载角色数据以确保状态正确
			updatedChar, err := m.charRepo.GetByID(char.ID)
			if err == nil && updatedChar != nil {
				char = updatedChar
				// 确保战士的怒气上限为100
				if char.ResourceType == "rage" {
					char.MaxResource = 100
				}
			}

			// 如果角色已经复活（不再死亡），自动恢复战斗
			if !char.IsDead {
				session.IsRunning = true
				m.addLog(session, "system", ">> 休息结束，自动恢复战斗", "#33ff33")
			} else {
				m.addLog(session, "system", ">> 休息结束，准备下一场战斗", "#33ff33")
			}
			logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
		} else {
			// 仍在休息中
			remaining := session.RestUntil.Sub(time.Now())
			if remaining > 0 {
				m.addLog(session, "system", fmt.Sprintf(">> 休息中... (剩余 %d 秒)", int(remaining.Seconds())+1), "#888888")
				logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
			}
		}

		// 保存角色数据更新
		if char.HP != initialHP || char.Resource != initialMP {
			m.charRepo.UpdateAfterBattle(char.ID, char.HP, char.Resource, char.Exp, char.Level,
				char.ExpToNext, char.MaxHP, char.MaxResource, char.Attack, char.Defense,
				char.Strength, char.Agility, char.Stamina, char.TotalKills)
		}

		return &BattleTickResult{
			Character:    char,
			Enemy:        nil,
			Enemies:      session.CurrentEnemies,
			Logs:         logs,
			IsRunning:    session.IsRunning,
			IsResting:    session.IsResting,
			RestUntil:    session.RestUntil,
			SessionKills: session.SessionKills,
			SessionGold:  session.SessionGold,
			SessionExp:   session.SessionExp,
			BattleCount:  session.BattleCount,
		}, nil
	}

	// 获取存活的敌人
	aliveEnemies := make([]*models.Monster, 0)
	for _, enemy := range session.CurrentEnemies {
		if enemy != nil && enemy.HP > 0 {
			aliveEnemies = append(aliveEnemies, enemy)
		}
	}

	// 如果没有敌人，生成新的
	if len(aliveEnemies) == 0 {
		// 重置本场战斗统计
		session.CurrentBattleExp = 0
		session.CurrentBattleGold = 0
		session.CurrentBattleKills = 0
		session.CurrentTurnIndex = -1 // 玩家回合

		// 战斗开始时，确保战士的怒气为0，最大怒气为100
		if char.ResourceType == "rage" {
			char.Resource = 0
			char.MaxResource = 100
		}

		err := m.spawnEnemies(session, char.Level)
		if err != nil {
			return nil, err
		}
		logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])

		// 标记刚遭遇敌人，需要等待1个tick再开始战斗
		session.JustEncountered = true

		// 更新存活敌人列表
		aliveEnemies = session.CurrentEnemies

		// 刚遭遇敌人，这个tick只显示信息，不执行战斗
		return &BattleTickResult{
			Character:    char,
			Enemy:        session.CurrentEnemy,
			Enemies:      session.CurrentEnemies,
			Logs:         logs,
			IsRunning:    session.IsRunning,
			IsResting:    session.IsResting,
			RestUntil:    session.RestUntil,
			SessionKills: session.SessionKills,
			SessionGold:  session.SessionGold,
			SessionExp:   session.SessionExp,
			BattleCount:  session.BattleCount,
		}, nil
	}

	// 如果刚遭遇敌人，这个tick只显示信息，不执行战斗
	if session.JustEncountered {
		session.JustEncountered = false // 清除标志，下一个tick开始战斗
		return &BattleTickResult{
			Character:    char,
			Enemy:        session.CurrentEnemy,
			Enemies:      session.CurrentEnemies,
			Logs:         logs,
			IsRunning:    session.IsRunning,
			IsResting:    session.IsResting,
			RestUntil:    session.RestUntil,
			SessionKills: session.SessionKills,
			SessionGold:  session.SessionGold,
			SessionExp:   session.SessionExp,
			BattleCount:  session.BattleCount,
		}, nil
	}

	// 回合制逻辑：CurrentTurnIndex == -1 表示玩家回合，>=0 表示敌人索引
	if session.CurrentTurnIndex == -1 {
		// 玩家回合：攻击第一个存活的敌人
		if len(aliveEnemies) > 0 {
			target := aliveEnemies[0]
			targetHPPercent := float64(target.HP) / float64(target.MaxHP)
			hasMultipleEnemies := len(aliveEnemies) > 1

			// 使用技能管理器选择技能
			skillState := m.skillManager.SelectBestSkill(char.ID, char.Resource, targetHPPercent, hasMultipleEnemies)

			var skillName string
			var playerDamage int
			var resourceCost int
			var resourceGain int
			var usedSkill bool
			var skillEffects map[string]interface{}
			var isCrit bool

			if skillState != nil {
				// 使用技能
				skillName = skillState.Skill.Name
				resourceCost = m.skillManager.GetSkillResourceCost(skillState)

				// 检查资源是否足够
				if resourceCost <= char.Resource {
					// 计算技能伤害（基础伤害，暴击在后面处理）
					baseDamage := m.skillManager.CalculateSkillDamage(skillState, char, target, m.passiveSkillManager, m.buffManager)

					// 计算暴击（技能也可以暴击，应用被动技能和Buff加成）
					actualCritRate := char.CritRate
					if m.passiveSkillManager != nil {
						critModifier := m.passiveSkillManager.GetPassiveModifier(char.ID, "crit_rate")
						actualCritRate = char.CritRate + critModifier/100.0
					}
					// 应用Buff的暴击率加成（鲁莽等）
					if m.buffManager != nil {
						critBuffValue := m.buffManager.GetBuffValue(char.ID, "crit_rate")
						if critBuffValue > 0 {
							actualCritRate = actualCritRate + critBuffValue/100.0
						}
					}
					if actualCritRate > 1.0 {
						actualCritRate = 1.0
					}
					isCrit = rand.Float64() < actualCritRate
					if isCrit {
						playerDamage = int(float64(baseDamage) * char.CritDamage)
					} else {
						playerDamage = baseDamage
					}

					// 应用技能效果
					skillEffects = m.skillManager.ApplySkillEffects(skillState, char, target)

					// 应用Buff/Debuff效果
					m.applySkillBuffs(skillState, char, target, skillEffects)

					// 应用Debuff到敌人（挫志怒吼、旋风斩等）
					m.applySkillDebuffs(skillState, char, target, aliveEnemies, skillEffects)

					// 消耗资源
					char.Resource -= resourceCost
					if char.Resource < 0 {
						char.Resource = 0
					}

					// 使用技能（设置冷却）
					m.skillManager.UseSkill(char.ID, skillState.SkillID)
					usedSkill = true

					// 处理技能特殊效果（怒气获得等）
					if rageGain, ok := skillEffects["rageGain"].(int); ok {
						// 应用被动技能的怒气获得加成（愤怒掌握等）
						actualRageGain := m.applyRageGenerationModifiers(char.ID, rageGain)
						char.Resource += actualRageGain
						resourceGain = actualRageGain
						if char.Resource > char.MaxResource {
							char.Resource = char.MaxResource
						}
					}

					// 处理AOE技能（旋风斩等）
					if skillState.Skill.TargetType == "enemy_all" {
						// 对所有敌人造成伤害
						for _, enemy := range aliveEnemies {
							if enemy.HP > 0 {
								damage := m.skillManager.CalculateSkillDamage(skillState, char, enemy, m.passiveSkillManager, m.buffManager)
								if isCrit {
									damage = int(float64(damage) * char.CritDamage)
								}
								enemy.HP -= damage
								if enemy.HP < 0 {
									enemy.HP = 0
								}
							}
						}
						// playerDamage用于日志显示（主目标伤害）
					} else if skillState.SkillID == "warrior_cleave" {
						// 顺劈斩：主目标+相邻目标
						target.HP -= playerDamage

						// 对相邻目标造成伤害（最多2个）
						adjacentCount := 0
						for _, enemy := range aliveEnemies {
							if enemy != target && enemy.HP > 0 && adjacentCount < 2 {
								// 计算相邻目标伤害
								if effect, ok := skillState.Effect["adjacentMultiplier"].(float64); ok {
									adjacentDamage := int(float64(char.Attack) * effect)
									adjacentDamage = adjacentDamage - enemy.Defense/2
									if adjacentDamage < 1 {
										adjacentDamage = 1
									}
									if isCrit {
										adjacentDamage = int(float64(adjacentDamage) * char.CritDamage)
									}
									enemy.HP -= adjacentDamage
									if enemy.HP < 0 {
										enemy.HP = 0
									}
									adjacentCount++
									m.addLog(session, "combat", fmt.Sprintf("%s 的顺劈斩波及到 %s，造成 %d 点伤害", char.Name, enemy.Name, adjacentDamage), "#ffaa00")
									logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
								}
							}
						}
					} else {
						// 单体技能
						target.HP -= playerDamage
					}
				} else {
					// 资源不足，使用普通攻击
					skillState = nil
				}
			}

			// 如果没有使用技能或资源不足，使用普通攻击
			if skillState == nil {
				skillName = "普通攻击"
				// 计算实际攻击力（应用被动技能加成）
				actualAttack := float64(char.Attack)
				if m.passiveSkillManager != nil {
					attackModifier := m.passiveSkillManager.GetPassiveModifier(char.ID, "attack")
					actualAttack = actualAttack * (1.0 + attackModifier/100.0)
					// 应用被动技能的伤害加成
					damageModifier := m.passiveSkillManager.GetPassiveModifier(char.ID, "damage")
					actualAttack = actualAttack * (1.0 + damageModifier/100.0)

					// 处理低血量时的攻击力加成（狂暴之心）
					hpPercent := float64(char.HP) / float64(char.MaxHP)
					passives := m.passiveSkillManager.GetPassiveSkills(char.ID)
					for _, passive := range passives {
						if passive.Passive.EffectType == "stat_mod" && passive.Passive.ID == "warrior_passive_berserker_heart" {
							// 根据等级计算触发阈值（1级50%，5级30%）
							threshold := 0.50 - float64(passive.Level-1)*0.05
							if hpPercent < threshold {
								// 根据等级计算攻击力加成（1级20%，5级60%）
								attackBonus := 20.0 + float64(passive.Level-1)*10.0
								actualAttack = actualAttack * (1.0 + attackBonus/100.0)
							}
						}
					}
				}
				// 应用Buff的攻击力加成（战斗怒吼、狂暴之怒、天神下凡等）
				if m.buffManager != nil {
					attackBuffValue := m.buffManager.GetBuffValue(char.ID, "attack")
					if attackBuffValue > 0 {
						actualAttack = actualAttack * (1.0 + attackBuffValue/100.0)
					}
				}
				baseDamage := m.calculateDamage(int(actualAttack), target.Defense)
				// 计算暴击率（应用被动技能和Buff加成）
				actualCritRate := char.CritRate
				if m.passiveSkillManager != nil {
					critModifier := m.passiveSkillManager.GetPassiveModifier(char.ID, "crit_rate")
					actualCritRate = char.CritRate + critModifier/100.0
				}
				// 应用Buff的暴击率加成（鲁莽等）
				if m.buffManager != nil {
					critBuffValue := m.buffManager.GetBuffValue(char.ID, "crit_rate")
					if critBuffValue > 0 {
						actualCritRate = actualCritRate + critBuffValue/100.0
					}
				}
				if actualCritRate > 1.0 {
					actualCritRate = 1.0
				}
				isCrit = rand.Float64() < actualCritRate
				if isCrit {
					playerDamage = int(float64(baseDamage) * char.CritDamage)
				} else {
					playerDamage = baseDamage
				}
				target.HP -= playerDamage
				resourceCost = 0
				usedSkill = false
			}
			// 如果使用了技能，isCrit已经在上面计算了

			// 普通攻击获得怒气（只有普通攻击才获得怒气，使用技能时不获得）
			if char.ResourceType == "rage" && !usedSkill {
				var baseRageGain int
				if isCrit {
					baseRageGain = 10 // 暴击获得10点怒气
				} else {
					baseRageGain = 5 // 普通攻击获得5点怒气
				}

				// 应用被动技能的怒气获得加成（愤怒掌握等）
				rageGain := m.applyRageGenerationModifiers(char.ID, baseRageGain)

				char.Resource += rageGain
				resourceGain = rageGain
				// 确保不超过最大值
				if char.Resource > char.MaxResource {
					char.Resource = char.MaxResource
				}
			}

			// 处理被动技能的特殊效果（攻击时触发）
			m.handlePassiveOnHitEffects(char, playerDamage, usedSkill, session, &logs)

			// 构建战斗日志消息，包含资源变化（带颜色）
			resourceChangeText := m.formatResourceChange(char.ResourceType, resourceCost, resourceGain)

			// 处理技能特殊效果日志
			if skillEffects != nil {
				if stun, ok := skillEffects["stun"].(bool); ok && stun {
					m.addLog(session, "combat", fmt.Sprintf("%s 被眩晕了！", target.Name), "#ff00ff")
					logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
				}
				// 处理基于伤害的恢复（嗜血等）
				if healPercent, ok := skillEffects["healPercent"].(float64); ok && usedSkill {
					healAmount := int(float64(playerDamage) * healPercent / 100.0)
					char.HP += healAmount
					if char.HP > char.MaxHP {
						char.HP = char.MaxHP
					}
					m.addLog(session, "heal", fmt.Sprintf("%s 恢复了 %d 点生命值", char.Name, healAmount), "#00ff00")
					logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
				}
				// 处理破釜沉舟的立即恢复（基于最大HP）
				if healMaxHpPercent, ok := skillEffects["healMaxHpPercent"].(float64); ok && usedSkill {
					healAmount := int(float64(char.MaxHP) * healMaxHpPercent / 100.0)
					char.HP += healAmount
					if char.HP > char.MaxHP {
						char.HP = char.MaxHP
					}
					m.addLog(session, "heal", fmt.Sprintf("%s 的破釜沉舟恢复了 %d 点生命值", char.Name, healAmount), "#00ff00")
					logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
				}
			}

			if isCrit {
				m.addLog(session, "combat", fmt.Sprintf("%s 使用 [%s] 💥暴击！对 %s 造成 %d 点伤害%s", char.Name, skillName, target.Name, playerDamage, resourceChangeText), "#ff6b6b")
			} else {
				m.addLog(session, "combat", fmt.Sprintf("%s 使用 [%s] 对 %s 造成 %d 点伤害%s", char.Name, skillName, target.Name, playerDamage, resourceChangeText), "#ffaa00")
			}
			logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])

			// 减少技能冷却时间
			m.skillManager.TickCooldowns(char.ID)

			// 减少Buff/Debuff持续时间
			expiredBuffs := m.buffManager.TickBuffs(char.ID)
			for _, effectID := range expiredBuffs {
				m.addLog(session, "buff", fmt.Sprintf("%s 的 %s 效果消失了", char.Name, effectID), "#888888")
				logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
			}

			// 检查目标是否死亡
			if target.HP <= 0 {
				// 确保HP归零
				target.HP = 0

				// 处理战争机器的击杀回怒效果
				m.handleWarMachineRageGain(char, session, &logs)

				// 敌人死亡
				expGain := target.ExpReward
				goldGain := target.GoldMin + rand.Intn(target.GoldMax-target.GoldMin+1)

				session.CurrentBattleExp += expGain
				session.CurrentBattleGold += goldGain
				session.CurrentBattleKills++
				session.SessionExp += expGain
				session.SessionGold += goldGain
				session.SessionKills++

				char.Exp += expGain
				char.TotalKills++

				// 检查升级
				for char.Exp >= char.ExpToNext {
					char.Exp -= char.ExpToNext
					char.Level++
					char.ExpToNext = int(float64(char.ExpToNext) * 1.5)

					// 升级属性提升
					char.MaxHP += 15
					char.HP = char.MaxHP

					// 战士的怒气最大值固定为100，不随升级改变
					if char.ResourceType == "rage" {
						char.MaxResource = 100
						// 升级时怒气保持不变，不重置为最大值
					} else {
						char.MaxResource += 8
						char.Resource = char.MaxResource
					}

					char.Strength += 2
					char.Agility += 1
					char.Stamina += 2
					char.Attack = char.Strength / 2
					char.Defense = char.Stamina / 3

					m.addLog(session, "levelup", fmt.Sprintf("🎉【升级】恭喜！%s 升到了 %d 级！", char.Name, char.Level), "#ffd700")
					logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
				}
			}

			// 移动到下一个敌人回合
			session.CurrentTurnIndex = 0
		}
	} else {
		// 敌人回合：当前索引的敌人攻击玩家
		if session.CurrentTurnIndex < len(aliveEnemies) {
			enemy := aliveEnemies[session.CurrentTurnIndex]
			enemyDamage := m.calculateDamage(enemy.Attack, char.Defense)

			// 应用buff/debuff效果（如盾牌格挡的减伤等）
			enemyDamage = m.buffManager.CalculateDamageTakenWithBuffs(enemyDamage, char.ID, true)

			// 处理被动技能的减伤效果（不灭意志等）
			enemyDamage = m.handlePassiveDamageReduction(char, enemyDamage)

			// 处理护盾效果（不灭壁垒等）
			shieldAmount := m.buffManager.GetBuffValue(char.ID, "shield")
			if shieldAmount > 0 {
				// 有护盾，先消耗护盾
				shieldInt := int(shieldAmount)
				if enemyDamage <= shieldInt {
					// 伤害完全被护盾吸收
					shieldInt -= enemyDamage
					absorbedDamage := enemyDamage
					enemyDamage = 0
					m.addLog(session, "shield", fmt.Sprintf("%s 的护盾吸收了 %d 点伤害", char.Name, absorbedDamage), "#00ffff")
					logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
					// 更新护盾值（通过更新Buff的value）
					m.updateShieldValue(char.ID, float64(shieldInt))
				} else {
					// 护盾被击破，剩余伤害继续
					absorbedDamage := shieldInt
					enemyDamage -= shieldInt
					m.addLog(session, "shield", fmt.Sprintf("%s 的护盾吸收了 %d 点伤害后被击破", char.Name, absorbedDamage), "#00ffff")
					logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
					m.updateShieldValue(char.ID, 0)
				}
			}

			// 处理被动技能的生存效果（坚韧不拔等）- 在受到伤害前检查
			originalHP := char.HP
			char.HP -= enemyDamage

			// 如果受到致命伤害，检查坚韧不拔效果
			if originalHP > 0 && char.HP <= 0 {
				if m.passiveSkillManager != nil {
					passives := m.passiveSkillManager.GetPassiveSkills(char.ID)
					for _, passive := range passives {
						if passive.Passive.EffectType == "survival" && passive.Passive.ID == "warrior_passive_unbreakable" {
							// 坚韧不拔：受到致命伤害时保留1点HP
							char.HP = 1
							m.addLog(session, "survival", fmt.Sprintf("%s 的坚韧不拔效果触发，保留了1点生命值！", char.Name), "#ff00ff")
							logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
							break // 只触发一次
						}
					}
				}
			}

			// 处理反击效果（反击风暴、复仇被动等）
			m.handleCounterAttacks(char, enemy, enemyDamage, session, &logs)

			// 处理被动技能的反射效果（盾牌反射被动等）
			m.handlePassiveReflectEffects(char, enemy, enemyDamage, session, &logs)

			// 处理主动技能的反射效果（盾牌反射技能等）
			m.handleActiveReflectEffects(char, enemy, enemyDamage, session, &logs)

			// 战士受到伤害时获得怒气
			resourceGain := 0
			if char.ResourceType == "rage" && enemyDamage > 0 {
				// 受到伤害获得怒气: 伤害/最大HP × 50，至少1点
				baseRageGain := int(float64(enemyDamage) / float64(char.MaxHP) * 50)
				if baseRageGain < 1 {
					baseRageGain = 1
				}

				// 应用被动技能的怒气获得加成（愤怒掌握等）
				rageGain := m.applyRageGenerationModifiers(char.ID, baseRageGain)

				char.Resource += rageGain
				resourceGain = rageGain
				if char.Resource > char.MaxResource {
					char.Resource = char.MaxResource
				}
			}

			// 构建战斗日志消息，包含资源变化（带颜色）
			resourceChangeText := m.formatResourceChange(char.ResourceType, 0, resourceGain)

			m.addLog(session, "combat", fmt.Sprintf("%s 攻击了 %s，造成 %d 点伤害%s", enemy.Name, char.Name, enemyDamage, resourceChangeText), "#ff4444")
			logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])

			// 检查玩家是否死亡
			if char.HP <= 0 {
				char.TotalDeaths++
				// 角色死亡时不停止战斗，保持 isRunning = true，这样休息状态可以自动处理
				// 用户已经开启了自动战斗，死亡只是暂时进入休息状态，休息结束后应该自动恢复战斗
				session.CurrentEnemies = nil
				session.CurrentEnemy = nil
				session.CurrentTurnIndex = -1

				// 角色死亡时，战士的怒气归0
				if char.ResourceType == "rage" {
					char.Resource = 0
				}

				// 计算复活时间
				reviveDuration := m.calculateReviveTime(userID)
				now := time.Now()
				reviveAt := now.Add(reviveDuration)

				// 设置角色HP为0（死亡状态）
				char.HP = 0
				char.IsDead = true
				char.ReviveAt = &reviveAt

				m.addLog(session, "death", fmt.Sprintf("%s 被击败了... 需要 %d 秒复活", char.Name, int(reviveDuration.Seconds())), "#ff0000")
				logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])

				// 战斗失败总结
				m.addBattleSummary(session, false, &logs)

				// 保存死亡数据（包括死亡标记、复活时间和怒气归0）
				m.charRepo.UpdateAfterDeath(char.ID, char.HP, char.Resource, char.TotalDeaths, &reviveAt)

				// 进入休息状态，休息时间 = 复活时间 + 恢复时间（恢复一半HP需要的时间）
				// 恢复时间：从0恢复到50% HP，每秒恢复2%，需要25秒
				recoveryTime := 25 * time.Second
				restDuration := reviveDuration + recoveryTime
				restUntil := now.Add(restDuration)
				session.IsResting = true
				session.RestUntil = &restUntil
				session.RestStartedAt = &now
				session.LastRestTick = &now
				session.RestSpeed = 1.0

				m.addLog(session, "system", fmt.Sprintf(">> 进入休息恢复状态 (预计 %d 秒)", int(restDuration.Seconds())+1), "#33ff33")
				logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])

				// 重置本场战斗统计
				session.CurrentBattleExp = 0
				session.CurrentBattleGold = 0
				session.CurrentBattleKills = 0
				session.CurrentTurnIndex = -1

				// 角色死亡时，立即返回，确保前端清除敌人显示
				// 保持 isRunning = true，这样按钮会显示"停止挂机"，休息状态可以自动处理
				return &BattleTickResult{
					Character:    char,
					Enemy:        nil,
					Enemies:      nil, // 明确返回 nil，让前端清除敌人显示
					Logs:         logs,
					IsRunning:    session.IsRunning, // 保持运行状态，不停止
					IsResting:    session.IsResting,
					RestUntil:    session.RestUntil,
					SessionKills: session.SessionKills,
					SessionGold:  session.SessionGold,
					SessionExp:   session.SessionExp,
					BattleCount:  session.BattleCount,
				}, nil
			} else {
				// 移动到下一个敌人或回到玩家回合
				session.CurrentTurnIndex++
				if session.CurrentTurnIndex >= len(aliveEnemies) {
					session.CurrentTurnIndex = -1 // 回到玩家回合
				}
			}
		} else {
			// 索引超出范围，回到玩家回合
			session.CurrentTurnIndex = -1
		}
	}

	// 更新存活敌人列表
	aliveEnemies = make([]*models.Monster, 0)
	for _, enemy := range session.CurrentEnemies {
		if enemy != nil && enemy.HP > 0 {
			aliveEnemies = append(aliveEnemies, enemy)
		}
	}

	// 如果所有敌人都被击败，处理战斗结束
	if len(aliveEnemies) == 0 && len(session.CurrentEnemies) > 0 {
		// 确保所有敌人的HP都归零
		for _, enemy := range session.CurrentEnemies {
			if enemy != nil && enemy.HP <= 0 {
				enemy.HP = 0
			}
		}

		// 战斗胜利总结
		m.addBattleSummary(session, true, &logs)

		// 战斗结束后，所有战士角色的怒气都归0
		for _, c := range characters {
			if c.ResourceType == "rage" {
				c.Resource = 0
			}
			// 保存所有角色的数据（包括战士的怒气归0）
			m.charRepo.UpdateAfterBattle(c.ID, c.HP, c.Resource, c.Exp, c.Level,
				c.ExpToNext, c.MaxHP, c.MaxResource, c.Attack, c.Defense,
				c.Strength, c.Agility, c.Stamina, c.TotalKills)
		}

		// 计算并开始休息
		restDuration := m.calculateRestTime(char)
		now := time.Now()
		restUntil := now.Add(restDuration)
		session.IsResting = true
		session.RestUntil = &restUntil
		session.RestStartedAt = &now
		session.LastRestTick = &now
		session.RestSpeed = 1.0 // 默认恢复速度

		m.addLog(session, "system", fmt.Sprintf(">> 开始休息恢复 (预计 %d 秒)", int(restDuration.Seconds())+1), "#33ff33")
		logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])

		// 重置本场战斗统计
		session.CurrentBattleExp = 0
		session.CurrentBattleGold = 0
		session.CurrentBattleKills = 0
		session.CurrentTurnIndex = -1

		// 先返回一次带有HP=0的敌人状态，让前端更新显示
		// 创建敌人副本，确保HP为0
		defeatedEnemies := make([]*models.Monster, len(session.CurrentEnemies))
		for i, enemy := range session.CurrentEnemies {
			if enemy != nil {
				defeatedEnemy := *enemy
				defeatedEnemy.HP = 0
				defeatedEnemies[i] = &defeatedEnemy
			}
		}

		// 清除敌人（在返回结果之后）
		session.CurrentEnemies = nil
		session.CurrentEnemy = nil

		// 返回带有HP=0的敌人状态
		return &BattleTickResult{
			Character:    char,
			Enemy:        nil,
			Enemies:      defeatedEnemies, // 返回HP=0的敌人副本
			Logs:         logs,
			IsRunning:    session.IsRunning,
			IsResting:    session.IsResting,
			RestUntil:    session.RestUntil,
			SessionKills: session.SessionKills,
			SessionGold:  session.SessionGold,
			SessionExp:   session.SessionExp,
			BattleCount:  session.BattleCount,
		}, nil
	}

	// 保存角色数据更新
	m.charRepo.UpdateAfterBattle(char.ID, char.HP, char.Resource, char.Exp, char.Level,
		char.ExpToNext, char.MaxHP, char.MaxResource, char.Attack, char.Defense,
		char.Strength, char.Agility, char.Stamina, char.TotalKills)

	return &BattleTickResult{
		Character:    char,
		Enemy:        session.CurrentEnemy,
		Enemies:      session.CurrentEnemies,
		Logs:         logs,
		IsRunning:    session.IsRunning,
		IsResting:    session.IsResting,
		RestUntil:    session.RestUntil,
		SessionKills: session.SessionKills,
		SessionGold:  session.SessionGold,
		SessionExp:   session.SessionExp,
		BattleCount:  session.BattleCount,
	}, nil
}

// spawnEnemy 生成敌人（向后兼容）
func (m *BattleManager) spawnEnemy(session *BattleSession, playerLevel int) error {
	return m.spawnEnemies(session, playerLevel)
}

// spawnEnemies 生成多个敌人（1-3个随机）
func (m *BattleManager) spawnEnemies(session *BattleSession, playerLevel int) error {
	if session.CurrentZone == nil {
		// 加载默认区域
		zone, err := m.gameRepo.GetZoneByID("elwynn")
		if err != nil {
			fmt.Printf("[ERROR] Failed to get zone: %v\n", err)
			return fmt.Errorf("failed to get zone: %v", err)
		}
		session.CurrentZone = zone
		fmt.Printf("[DEBUG] Loaded zone: %s\n", zone.Name)
	}

	// 获取区域怪物
	monsters, err := m.gameRepo.GetMonstersByZone(session.CurrentZone.ID)
	if err != nil {
		fmt.Printf("[ERROR] Failed to get monsters: %v\n", err)
		return fmt.Errorf("failed to get monsters: %v", err)
	}
	if len(monsters) == 0 {
		fmt.Printf("[ERROR] No monsters in zone %s\n", session.CurrentZone.ID)
		return fmt.Errorf("no monsters in zone %s", session.CurrentZone.ID)
	}
	fmt.Printf("[DEBUG] Found %d monsters in zone\n", len(monsters))

	// 随机生成1-3个敌人
	enemyCount := 1 + rand.Intn(3) // 1-3个
	session.CurrentEnemies = make([]*models.Monster, 0, enemyCount)

	var enemyNames []string
	for i := 0; i < enemyCount; i++ {
		// 随机选择一个怪物模板
		template := monsters[rand.Intn(len(monsters))]

		enemy := &models.Monster{
			ID:        template.ID,
			ZoneID:    template.ZoneID,
			Name:      template.Name,
			Level:     template.Level,
			Type:      template.Type,
			HP:        template.HP,
			MaxHP:     template.HP,
			Attack:    template.Attack,
			Defense:   template.Defense,
			ExpReward: template.ExpReward,
			GoldMin:   template.GoldMin,
			GoldMax:   template.GoldMax,
		}
		session.CurrentEnemies = append(session.CurrentEnemies, enemy)
		enemyNames = append(enemyNames, fmt.Sprintf("%s (Lv.%d)", enemy.Name, enemy.Level))
	}

	// 保留 CurrentEnemy 用于向后兼容（指向第一个敌人）
	if len(session.CurrentEnemies) > 0 {
		session.CurrentEnemy = session.CurrentEnemies[0]
	}

	session.BattleCount++
	enemyList := fmt.Sprintf("%s", enemyNames[0])
	if len(enemyNames) > 1 {
		enemyList = fmt.Sprintf("%s 等 %d 个敌人", enemyNames[0], len(enemyNames))
	}
	m.addLog(session, "encounter", fmt.Sprintf("━━━ 战斗 #%d ━━━ 遭遇: %s", session.BattleCount, enemyList), "#ffff00")

	return nil
}

// ChangeZone 切换区域
func (m *BattleManager) ChangeZone(userID int, zoneID string, playerLevel int) error {
	session := m.GetOrCreateSession(userID)

	zone, err := m.gameRepo.GetZoneByID(zoneID)
	if err != nil {
		return fmt.Errorf("zone not found: %s", zoneID)
	}

	if playerLevel < zone.MinLevel {
		return fmt.Errorf("level too low, need level %d", zone.MinLevel)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	session.CurrentZone = zone
	session.CurrentEnemy = nil

	m.addLog(session, "zone", fmt.Sprintf(">> 你来到了 [%s]", zone.Name), "#00ffff")
	m.addLog(session, "zone", zone.Description, "#888888")

	return nil
}

// GetBattleStatus 获取战斗状态
func (m *BattleManager) GetBattleStatus(userID int) *models.BattleStatus {
	session := m.GetSession(userID)
	if session == nil {
		return &models.BattleStatus{
			IsRunning: false,
		}
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	status := &models.BattleStatus{
		IsRunning:      session.IsRunning,
		CurrentMonster: session.CurrentEnemy,
		CurrentEnemies: session.CurrentEnemies,
		BattleCount:    session.BattleCount,
		TotalKills:     session.SessionKills,
		TotalExp:       session.SessionExp,
		TotalGold:      session.SessionGold,
		IsResting:      session.IsResting,
		RestUntil:      session.RestUntil,
	}

	if session.CurrentZone != nil {
		status.CurrentZoneID = session.CurrentZone.ID
	}

	return status
}

// GetBattleLogs 获取战斗日志
func (m *BattleManager) GetBattleLogs(userID int, limit int) []models.BattleLog {
	session := m.GetSession(userID)
	if session == nil {
		return []models.BattleLog{}
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	logs := session.BattleLogs
	if limit > 0 && len(logs) > limit {
		logs = logs[len(logs)-limit:]
	}
	return logs
}

// calculateDamage 计算伤害
func (m *BattleManager) calculateDamage(attack, defense int) int {
	baseDamage := attack - defense/2
	if baseDamage < 1 {
		baseDamage = 1
	}
	// 添加随机波动 ±20%
	variance := float64(baseDamage) * 0.2
	damage := float64(baseDamage) + (rand.Float64()*2-1)*variance
	return int(damage)
}

// addLog 添加日志
func (m *BattleManager) addLog(session *BattleSession, logType, message, color string) {
	log := models.BattleLog{
		Message:   message,
		LogType:   logType,
		CreatedAt: time.Now(),
	}
	session.BattleLogs = append(session.BattleLogs, log)

	// 保持日志数量在合理范围
	if len(session.BattleLogs) > 200 {
		session.BattleLogs = session.BattleLogs[len(session.BattleLogs)-200:]
	}
}

// addBattleSummary 添加战斗总结和分割线
func (m *BattleManager) addBattleSummary(session *BattleSession, isVictory bool, logs *[]models.BattleLog) {
	// 生成战斗总结，使用不同颜色标记不同指标
	var summaryMsg string
	if isVictory {
		if session.CurrentBattleKills > 0 {
			// 使用HTML标签为不同部分添加颜色
			// 结果：金色 #ffd700，击杀：红色 #ff4444，经验：蓝色 #3d85c6，金币：金色 #ffd700
			summaryMsg = fmt.Sprintf("━━━ 战斗总结 ━━━ 结果: <span style=\"color: #ffd700\">✓ 胜利</span> | 击杀: <span style=\"color: #ff4444\">%d</span> | 经验: <span style=\"color: #3d85c6\">%d</span> | 金币: <span style=\"color: #ffd700\">%d</span>",
				session.CurrentBattleKills, session.CurrentBattleExp, session.CurrentBattleGold)
		} else {
			summaryMsg = "━━━ 战斗总结 ━━━ 结果: <span style=\"color: #ffd700\">✓ 胜利</span>"
		}
		m.addLog(session, "battle_summary", summaryMsg, "#ffd700")
		*logs = append(*logs, session.BattleLogs[len(session.BattleLogs)-1])
	} else {
		// 失败时的总结
		if session.CurrentBattleKills > 0 {
			// 结果：红色 #ff6666，击杀：橙色 #ffaa00，经验：蓝色 #3d85c6，金币：金色 #ffd700
			summaryMsg = fmt.Sprintf("━━━ 战斗总结 ━━━ 结果: <span style=\"color: #ff6666\">✗ 失败</span> | 击杀: <span style=\"color: #ffaa00\">%d</span> | 经验: <span style=\"color: #3d85c6\">%d</span> | 金币: <span style=\"color: #ffd700\">%d</span>",
				session.CurrentBattleKills, session.CurrentBattleExp, session.CurrentBattleGold)
		} else {
			summaryMsg = "━━━ 战斗总结 ━━━ 结果: <span style=\"color: #ff6666\">✗ 失败</span>"
		}
		m.addLog(session, "battle_summary", summaryMsg, "#ff6666")
		*logs = append(*logs, session.BattleLogs[len(session.BattleLogs)-1])
	}

	// 添加分割线
	m.addLog(session, "battle_separator", "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━", "#666666")
	*logs = append(*logs, session.BattleLogs[len(session.BattleLogs)-1])
}

// getResourceName 获取资源的中文名称
func (m *BattleManager) getResourceName(resourceType string) string {
	switch resourceType {
	case "rage":
		return "怒气"
	case "mana":
		return "MP"
	case "energy":
		return "能量"
	default:
		return "资源"
	}
}

// getResourceColor 获取资源的颜色（参考魔兽世界）
func (m *BattleManager) getResourceColor(resourceType string) string {
	switch resourceType {
	case "rage":
		return "#ff4444" // 红色 - 怒气
	case "mana":
		return "#3d85c6" // 蓝色 - 法力
	case "energy":
		return "#ffd700" // 金色/黄色 - 能量
	default:
		return "#ffffff" // 白色 - 默认
	}
}

// formatResourceChange 格式化资源变化文本（带颜色）
func (m *BattleManager) formatResourceChange(resourceType string, cost int, gain int) string {
	if cost == 0 && gain == 0 {
		return ""
	}

	resourceName := m.getResourceName(resourceType)
	color := m.getResourceColor(resourceType)

	var parts []string
	if cost > 0 {
		parts = append(parts, fmt.Sprintf("<span style=\"color: %s\">-%d</span>", color, cost))
	}
	if gain > 0 {
		parts = append(parts, fmt.Sprintf("<span style=\"color: %s\">+%d</span>", color, gain))
	}

	if len(parts) == 0 {
		return ""
	}

	// 将多个部分用空格连接
	changeText := ""
	for i, part := range parts {
		if i > 0 {
			changeText += " "
		}
		changeText += part
	}

	return fmt.Sprintf(" (<span style=\"color: %s\">%s</span> %s)", color, resourceName, changeText)
}

// getRandomSkillName 获取随机技能名称
func (m *BattleManager) getRandomSkillName(classID string) string {
	skills := map[string][]string{
		"warrior": {"英勇打击", "雷霆一击", "顺劈斩", "致死打击"},
		"paladin": {"圣光术", "十字军打击", "正义之锤", "审判"},
		"hunter":  {"奥术射击", "多重射击", "瞄准射击", "稳固射击"},
		"rogue":   {"邪恶攻击", "剔骨", "背刺", "毒刃"},
		"priest":  {"惩击", "暗言术:痛", "神圣之火", "心灵震爆"},
		"mage":    {"火球术", "寒冰箭", "奥术飞弹", "炎爆术"},
		"warlock": {"暗影箭", "腐蚀术", "献祭", "混乱箭"},
		"druid":   {"月火术", "愤怒", "挥击", "横扫"},
		"shaman":  {"闪电箭", "闪电链", "熔岩爆裂", "烈焰震击"},
	}

	if classSkills, ok := skills[classID]; ok {
		return classSkills[rand.Intn(len(classSkills))]
	}
	return "普通攻击"
}

// getSkillForAttack 获取攻击技能名称和消耗
func (m *BattleManager) getSkillForAttack(char *models.Character) (string, int) {
	// 战士技能及其怒气消耗
	warriorSkills := []struct {
		name string
		cost int
	}{
		{"英勇打击", 10},
		{"雷霆一击", 15},
		{"顺劈斩", 12},
		{"致死打击", 20},
	}

	// 如果是战士，返回随机技能和消耗
	if char.ResourceType == "rage" {
		skill := warriorSkills[rand.Intn(len(warriorSkills))]
		return skill.name, skill.cost
	}

	// 其他职业使用普通技能，不消耗资源（或消耗法力，但这里简化处理）
	skillName := m.getRandomSkillName(char.ClassID)
	return skillName, 0
}

// calculateReviveTime 计算复活时间（根据死亡人数）
func (m *BattleManager) calculateReviveTime(userID int) time.Duration {
	// 获取所有角色（所有角色都参与战斗）
	characters, err := m.charRepo.GetByUserID(userID)
	if err != nil {
		return 30 * time.Second // 默认30秒
	}

	// 统计死亡角色的数量
	deadCount := 0
	for _, char := range characters {
		if char.IsDead {
			deadCount++
		}
	}

	// 如果没有死亡角色，返回默认值
	if deadCount == 0 {
		deadCount = 1 // 至少有一个角色死亡才会调用这个函数
	}

	// 根据死亡人数计算复活时间（秒）
	var reviveSeconds int
	switch deadCount {
	case 1:
		reviveSeconds = 30
	case 2:
		reviveSeconds = 60
	case 3:
		reviveSeconds = 90
	case 4:
		reviveSeconds = 120
	default: // 5人或更多
		reviveSeconds = 180
	}

	return time.Duration(reviveSeconds) * time.Second
}

// calculateRestTime 计算休息时间（基于HP/MP损失）
// 注意：战士的怒气不需要恢复，战斗结束后直接归0，每场战斗从0开始
func (m *BattleManager) calculateRestTime(char *models.Character) time.Duration {
	hpLoss := float64(char.MaxHP - char.HP)

	// 战士的怒气不需要恢复，只计算HP损失
	// 其他职业需要计算MP损失
	var mpLoss float64
	if char.ResourceType != "rage" {
		mpLoss = float64(char.MaxResource - char.Resource)
	} else {
		// 战士的怒气不参与休息时间计算
		mpLoss = 0
	}

	// 如果已经满血满蓝（或满血），不需要休息
	if hpLoss <= 0 && mpLoss <= 0 {
		return 0
	}

	// 分别计算HP和MP的恢复时间
	// 每秒恢复2%，所以需要的时间 = 损失百分比 / 0.02 = 损失百分比 * 50
	hpLossPercent := hpLoss / float64(char.MaxHP)

	hpRestSeconds := hpLossPercent * 50.0
	var mpRestSeconds float64
	if char.ResourceType != "rage" && char.MaxResource > 0 {
		mpLossPercent := mpLoss / float64(char.MaxResource)
		mpRestSeconds = mpLossPercent * 50.0
	} else {
		mpRestSeconds = 0
	}

	// 取两者中的最大值，因为HP和MP是同时恢复的
	restSeconds := hpRestSeconds
	if mpRestSeconds > restSeconds {
		restSeconds = mpRestSeconds
	}

	// 最少1秒
	if restSeconds < 1.0 {
		restSeconds = 1.0
	}

	// 应用恢复速度倍率（未来可以从装备获取）
	restSpeed := 1.0 // 默认恢复速度
	if restSpeed > 0 {
		restSeconds = restSeconds / restSpeed
	}

	return time.Duration(restSeconds) * time.Second
}

// processRest 处理休息期间的恢复
func (m *BattleManager) processRest(session *BattleSession, char *models.Character) {
	if !session.IsResting || session.RestUntil == nil || session.RestStartedAt == nil {
		return
	}

	now := time.Now()

	// 检查角色是否已经复活（如果角色死亡且有复活时间）
	if char.IsDead && char.ReviveAt != nil {
		if now.After(*char.ReviveAt) || now.Equal(*char.ReviveAt) {
			// 复活时间到了，恢复角色到一半HP
			char.HP = char.MaxHP / 2
			char.IsDead = false
			char.ReviveAt = nil

			// 更新数据库，清除死亡标记
			m.charRepo.SetDead(char.ID, false, nil)

			// 更新角色HP
			m.charRepo.UpdateAfterBattle(char.ID, char.HP, char.Resource, char.Exp, char.Level,
				char.ExpToNext, char.MaxHP, char.MaxResource, char.Attack, char.Defense,
				char.Strength, char.Agility, char.Stamina, char.TotalKills)

			// 记录复活日志
			m.addLog(session, "revive", fmt.Sprintf("%s 已复活，HP恢复至 %d/%d", char.Name, char.HP, char.MaxHP), "#00ff00")
		}
	}

	// 检查是否已经恢复满血满蓝，如果是则提前结束休息
	if char.HP >= char.MaxHP && char.Resource >= char.MaxResource {
		session.IsResting = false
		session.RestUntil = nil
		session.RestStartedAt = nil
		session.LastRestTick = nil
		return
	}

	if now.Before(*session.RestUntil) {
		// 计算从上次恢复到现在经过的时间
		var elapsed time.Duration
		if session.LastRestTick == nil {
			// 第一次恢复，从休息开始时间计算
			elapsed = now.Sub(*session.RestStartedAt)
		} else {
			// 从上次恢复时间计算
			elapsed = now.Sub(*session.LastRestTick)
		}

		// 如果时间间隔太长（超过1秒），限制为1秒，避免一次性恢复过多
		if elapsed > time.Second {
			elapsed = time.Second
		}

		// 如果时间间隔太短（小于0.1秒），跳过恢复，避免频繁计算
		if elapsed < 100*time.Millisecond {
			return
		}

		// 计算恢复速度（每秒恢复最大值的2%）
		restSpeed := session.RestSpeed
		if restSpeed <= 0 {
			restSpeed = 1.0
		}

		// 计算经过的秒数
		elapsedSeconds := elapsed.Seconds()

		// 如果角色已经死亡但还没到复活时间，不恢复HP
		if char.IsDead && char.ReviveAt != nil && now.Before(*char.ReviveAt) {
			// 只恢复资源（如果适用），不恢复HP
		} else {
			// 根据实际经过的时间计算恢复量
			hpRegenPercent := 0.02 * restSpeed * elapsedSeconds // 每秒2%

			hpRegen := int(float64(char.MaxHP) * hpRegenPercent)

			// 确保至少恢复1点（如果还没满）
			if hpRegen < 1 && char.HP < char.MaxHP {
				hpRegen = 1
			}

			char.HP += hpRegen
			if char.HP > char.MaxHP {
				char.HP = char.MaxHP
			}
		}

		// 战士的怒气不在休息时恢复，只在战斗中通过攻击/受击获得
		if char.ResourceType != "rage" {
			mpRegenPercent := 0.02 * restSpeed * elapsedSeconds
			mpRegen := int(float64(char.MaxResource) * mpRegenPercent)

			if mpRegen < 1 && char.Resource < char.MaxResource {
				mpRegen = 1
			}

			char.Resource += mpRegen
			if char.Resource > char.MaxResource {
				char.Resource = char.MaxResource
			}
		}

		// 更新上次恢复时间
		session.LastRestTick = &now

		// 恢复后再次检查是否已满，如果满了则提前结束休息
		if char.HP >= char.MaxHP && char.Resource >= char.MaxResource {
			session.IsResting = false
			session.RestUntil = nil
			session.RestStartedAt = nil
			session.LastRestTick = nil
		}
	} else {
		// 休息时间到了，结束休息
		// 确保HP已满
		if char.HP < char.MaxHP {
			char.HP = char.MaxHP
		}
		// 战士的怒气不在休息时恢复，只在战斗中通过攻击/受击获得
		if char.ResourceType != "rage" {
			if char.Resource < char.MaxResource {
				char.Resource = char.MaxResource
			}
		}
		session.IsResting = false
		session.RestUntil = nil
		session.RestStartedAt = nil
		session.LastRestTick = nil
	}
}

// BattleTickResult 战斗回合结果
type BattleTickResult struct {
	Character    *models.Character  `json:"character"`
	Enemy        *models.Monster    `json:"enemy,omitempty"`
	Enemies      []*models.Monster  `json:"enemies,omitempty"` // 多个敌人支持
	Logs         []models.BattleLog `json:"logs"`
	IsRunning    bool               `json:"isRunning"`
	IsResting    bool               `json:"isResting"`           // 是否在休息
	RestUntil    *time.Time         `json:"restUntil,omitempty"` // 休息结束时间
	SessionKills int                `json:"sessionKills"`
	SessionGold  int                `json:"sessionGold"`
	SessionExp   int                `json:"sessionExp"`
	BattleCount  int                `json:"battleCount"`
}

// applySkillBuffs 应用技能的Buff/Debuff效果
func (m *BattleManager) applySkillBuffs(skillState *CharacterSkillState, character *models.Character, target *models.Monster, skillEffects map[string]interface{}) {
	skill := skillState.Skill
	effect := skillState.Effect

	switch skill.ID {
	case "warrior_shield_block":
		// 盾牌格挡：减少受到的物理伤害
		if damageReduction, ok := effect["damageReduction"].(float64); ok {
			duration := 2
			if d, ok := effect["duration"].(int); ok {
				duration = d
			}
			m.buffManager.ApplyBuff(character.ID, "shield_block", "盾牌格挡", "buff", true, duration, -damageReduction, "physical_damage_taken", "")
		}
	case "warrior_battle_shout":
		// 战斗怒吼：提升攻击力
		if attackBonus, ok := effect["attackBonus"].(float64); ok {
			duration := 5
			if d, ok := effect["duration"].(int); ok {
				duration = d
			}
			m.buffManager.ApplyBuff(character.ID, "battle_shout", "战斗怒吼", "buff", true, duration, attackBonus, "attack", "")
		}
	case "warrior_demoralizing_shout":
		// 挫志怒吼：降低所有敌人攻击力（在applySkillDebuffs中处理）
	case "warrior_whirlwind":
		// 旋风斩：降低所有敌人防御（在applySkillDebuffs中处理）
	case "warrior_mortal_strike":
		// 致死打击：降低目标治疗效果
		if healingReduction, ok := effect["healingReduction"].(float64); ok {
			duration := 3
			if d, ok := effect["debuffDuration"].(float64); ok {
				duration = int(d)
			}
			// 应用到目标敌人
			if target != nil {
				m.buffManager.ApplyEnemyDebuff(target.ID, "mortal_strike", "致死打击", "debuff", duration, healingReduction, "healing_received", "")
			}
		}
	case "warrior_last_stand":
		// 破釜沉舟：立即恢复最大HP的百分比
		if healPercent, ok := effect["healPercent"].(float64); ok {
			// 立即恢复
			healAmount := int(float64(character.MaxHP) * healPercent / 100.0)
			character.HP += healAmount
			if character.HP > character.MaxHP {
				character.HP = character.MaxHP
			}
			// 通过skillEffects传递，在战斗日志中显示
			skillEffects["healMaxHpPercent"] = healPercent
		}
	case "warrior_unbreakable_barrier":
		// 不灭壁垒：获得护盾
		if shieldPercent, ok := effect["shieldPercent"].(float64); ok {
			duration := 4
			if d, ok := effect["duration"].(float64); ok {
				duration = int(d)
			}
			shieldAmount := int(float64(character.MaxHP) * shieldPercent / 100.0)
			// 使用Buff存储护盾值，statAffected为"shield"，value为护盾值
			m.buffManager.ApplyBuff(character.ID, "unbreakable_barrier", "不灭壁垒", "buff", true, duration, float64(shieldAmount), "shield", "")
		}
	case "warrior_shield_reflection":
		// 盾牌反射：反射受到的伤害
		if reflectPercent, ok := effect["reflectPercent"].(float64); ok {
			duration := 2
			if d, ok := effect["duration"].(float64); ok {
				duration = int(d)
			}
			// 使用Buff存储反射比例，statAffected为"reflect"，value为反射百分比
			m.buffManager.ApplyBuff(character.ID, "shield_reflection", "盾牌反射", "buff", true, duration, reflectPercent, "reflect", "")
		}
	case "warrior_shield_wall":
		// 盾墙：大幅减少受到的伤害
		if damageReduction, ok := effect["damageReduction"].(float64); ok {
			duration := 2
			if d, ok := effect["duration"].(float64); ok {
				duration = int(d)
			}
			m.buffManager.ApplyBuff(character.ID, "shield_wall", "盾墙", "buff", true, duration, -damageReduction, "damage_taken", "")
		}
	case "warrior_recklessness":
		// 鲁莽：提升暴击率，但受到伤害增加
		if critBonus, ok := effect["critBonus"].(float64); ok {
			duration := 3
			if d, ok := effect["duration"].(float64); ok {
				duration = int(d)
			}
			m.buffManager.ApplyBuff(character.ID, "recklessness_crit", "鲁莽", "buff", true, duration, critBonus, "crit_rate", "")
		}
		if damageIncrease, ok := effect["damageTakenIncrease"].(float64); ok {
			duration := 3
			if d, ok := effect["duration"].(float64); ok {
				duration = int(d)
			}
			m.buffManager.ApplyBuff(character.ID, "recklessness_damage", "鲁莽", "debuff", false, duration, damageIncrease, "damage_taken", "")
		}
	case "warrior_retaliation":
		// 反击风暴：受到攻击时反击
		if counterDamage, ok := effect["counterDamagePercent"].(float64); ok {
			duration := 3
			if d, ok := effect["duration"].(float64); ok {
				duration = int(d)
			}
			m.buffManager.ApplyBuff(character.ID, "retaliation", "反击风暴", "buff", true, duration, counterDamage, "counter_attack", "")
		}
	case "warrior_berserker_rage":
		// 狂暴之怒：提升攻击力和怒气获取
		if attackBonus, ok := effect["attackBonus"].(float64); ok {
			duration := 4
			if d, ok := effect["duration"].(float64); ok {
				duration = int(d)
			}
			m.buffManager.ApplyBuff(character.ID, "berserker_rage", "狂暴之怒", "buff", true, duration, attackBonus, "attack", "")
		}
	case "warrior_avatar":
		// 天神下凡：大幅提升攻击力，免疫控制
		if attackBonus, ok := effect["attackBonus"].(float64); ok {
			duration := 3
			if d, ok := effect["duration"].(float64); ok {
				duration = int(d)
			}
			m.buffManager.ApplyBuff(character.ID, "avatar", "天神下凡", "buff", true, duration, attackBonus, "attack", "")
		}
		if immuneCC, ok := effect["immuneCC"].(bool); ok && immuneCC {
			duration := 3
			if d, ok := effect["duration"].(float64); ok {
				duration = int(d)
			}
			m.buffManager.ApplyBuff(character.ID, "avatar_cc_immune", "天神下凡", "buff", true, duration, 1.0, "cc_immune", "")
		}
	}
}

// handleCounterAttacks 处理反击效果
func (m *BattleManager) handleCounterAttacks(character *models.Character, attacker *models.Monster, damageTaken int, session *BattleSession, logs *[]models.BattleLog) {
	// 处理Buff的反击效果（反击风暴）
	buffs := m.buffManager.GetBuffs(character.ID)
	for _, buff := range buffs {
		if buff.StatAffected == "counter_attack" && buff.IsBuff {
			// 反击风暴：对攻击者造成反击伤害
			counterDamage := int(float64(character.Attack) * buff.Value / 100.0)
			attacker.HP -= counterDamage
			if attacker.HP < 0 {
				attacker.HP = 0
			}
			m.addLog(session, "combat", fmt.Sprintf("%s 的反击风暴对 %s 造成 %d 点反击伤害！", character.Name, attacker.Name, counterDamage), "#ff8800")
			*logs = append(*logs, session.BattleLogs[len(session.BattleLogs)-1])
		}
	}

	// 处理被动技能的反击效果（复仇）
	if m.passiveSkillManager != nil {
		passives := m.passiveSkillManager.GetPassiveSkills(character.ID)
		for _, passive := range passives {
			if passive.Passive.EffectType == "counter_attack" {
				// 复仇：受到攻击时概率反击
				// effectValue是触发概率（百分比），需要根据等级计算实际概率和伤害
				triggerChance := passive.EffectValue / 100.0
				if rand.Float64() < triggerChance {
					// 计算反击伤害（根据等级：1级100%，5级180%）
					counterDamagePercent := 100.0 + float64(passive.Level-1)*20.0
					// 计算实际攻击力（应用被动技能和Buff加成）
					actualAttack := float64(character.Attack)
					if m.passiveSkillManager != nil {
						attackModifier := m.passiveSkillManager.GetPassiveModifier(character.ID, "attack")
						actualAttack = actualAttack * (1.0 + attackModifier/100.0)
					}
					if m.buffManager != nil {
						attackBuffValue := m.buffManager.GetBuffValue(character.ID, "attack")
						if attackBuffValue > 0 {
							actualAttack = actualAttack * (1.0 + attackBuffValue/100.0)
						}
					}
					counterDamage := int(actualAttack * counterDamagePercent / 100.0)
					counterDamage = counterDamage - attacker.Defense/2
					if counterDamage < 1 {
						counterDamage = 1
					}
					attacker.HP -= counterDamage
					if attacker.HP < 0 {
						attacker.HP = 0
					}
					m.addLog(session, "combat", fmt.Sprintf("%s 的复仇对 %s 造成 %d 点反击伤害！", character.Name, attacker.Name, counterDamage), "#ff8800")
					*logs = append(*logs, session.BattleLogs[len(session.BattleLogs)-1])
				}
			}
		}
	}
}

// handlePassiveOnHitEffects 处理被动技能的攻击时效果
func (m *BattleManager) handlePassiveOnHitEffects(character *models.Character, damageDealt int, usedSkill bool, session *BattleSession, logs *[]models.BattleLog) {
	if m.passiveSkillManager == nil {
		return
	}

	passives := m.passiveSkillManager.GetPassiveSkills(character.ID)
	for _, passive := range passives {
		switch passive.Passive.EffectType {
		case "on_hit_heal":
			// 血之狂热：每次攻击恢复生命值
			healPercent := passive.EffectValue // 百分比值（如1.0表示1%）
			healAmount := int(float64(character.MaxHP) * healPercent / 100.0)
			if healAmount > 0 {
				character.HP += healAmount
				if character.HP > character.MaxHP {
					character.HP = character.MaxHP
				}
				m.addLog(session, "heal", fmt.Sprintf("%s 的血之狂热恢复了 %d 点生命值", character.Name, healAmount), "#00ff00")
				*logs = append(*logs, session.BattleLogs[len(session.BattleLogs)-1])
			}
		}
	}
}

// handlePassiveDamageReduction 处理被动技能的减伤效果
func (m *BattleManager) handlePassiveDamageReduction(character *models.Character, damage int) int {
	if m.passiveSkillManager == nil {
		return damage
	}

	passives := m.passiveSkillManager.GetPassiveSkills(character.ID)
	for _, passive := range passives {
		if passive.Passive.EffectType == "survival" && passive.Passive.ID == "warrior_passive_unbreakable_will" {
			// 不灭意志：HP低于阈值时减伤
			hpPercent := float64(character.HP) / float64(character.MaxHP)
			// 根据等级计算触发阈值（1级30%，5级10%）
			threshold := 0.30 - float64(passive.Level-1)*0.05
			if hpPercent < threshold {
				// 根据等级计算减伤比例（1级25%，5级65%）
				reductionPercent := 25.0 + float64(passive.Level-1)*10.0
				damage = int(float64(damage) * (1.0 - reductionPercent/100.0))
				if damage < 1 {
					damage = 1
				}
			}
		}
	}

	return damage
}

// handleActiveReflectEffects 处理主动技能的反射效果
func (m *BattleManager) handleActiveReflectEffects(character *models.Character, attacker *models.Monster, damageTaken int, session *BattleSession, logs *[]models.BattleLog) {
	if m.buffManager == nil {
		return
	}

	buffs := m.buffManager.GetBuffs(character.ID)
	for _, buff := range buffs {
		if buff.StatAffected == "reflect" && buff.IsBuff && buff.EffectID == "shield_reflection" {
			// 盾牌反射（主动技能）：反射受到的伤害
			reflectPercent := buff.Value // 百分比值（如50.0表示50%）
			reflectDamage := int(float64(damageTaken) * reflectPercent / 100.0)
			if reflectDamage > 0 {
				attacker.HP -= reflectDamage
				if attacker.HP < 0 {
					attacker.HP = 0
				}
				m.addLog(session, "combat", fmt.Sprintf("%s 的盾牌反射对 %s 造成 %d 点反射伤害！", character.Name, attacker.Name, reflectDamage), "#ff8800")
				*logs = append(*logs, session.BattleLogs[len(session.BattleLogs)-1])
			}
		}
	}
}

// updateShieldValue 更新护盾值
func (m *BattleManager) updateShieldValue(characterID int, newShieldValue float64) {
	if m.buffManager == nil {
		return
	}

	buffs := m.buffManager.GetBuffs(characterID)
	if buff, exists := buffs["unbreakable_barrier"]; exists {
		buff.Value = newShieldValue
	}
}

// applySkillDebuffs 应用技能的Debuff效果到敌人
func (m *BattleManager) applySkillDebuffs(skillState *CharacterSkillState, character *models.Character, target *models.Monster, allEnemies []*models.Monster, skillEffects map[string]interface{}) {
	skill := skillState.Skill
	effect := skillState.Effect

	switch skill.ID {
	case "warrior_demoralizing_shout":
		// 挫志怒吼：降低所有敌人攻击力
		if attackReduction, ok := effect["attackReduction"].(float64); ok {
			duration := 3
			if d, ok := effect["duration"].(float64); ok {
				duration = int(d)
			}
			// 应用到所有存活的敌人
			for _, enemy := range allEnemies {
				if enemy.HP > 0 {
					m.buffManager.ApplyEnemyDebuff(enemy.ID, "demoralizing_shout", "挫志怒吼", "debuff", duration, attackReduction, "attack", "")
				}
			}
		}
	case "warrior_whirlwind":
		// 旋风斩：降低所有敌人防御
		if defenseReduction, ok := effect["defenseReduction"].(float64); ok {
			duration := 2
			if d, ok := effect["debuffDuration"].(float64); ok {
				duration = int(d)
			}
			// 应用到所有存活的敌人
			for _, enemy := range allEnemies {
				if enemy.HP > 0 {
					m.buffManager.ApplyEnemyDebuff(enemy.ID, "whirlwind", "旋风斩", "debuff", duration, defenseReduction, "defense", "")
				}
			}
		}
	case "warrior_mortal_strike":
		// 致死打击：降低目标治疗效果
		if healingReduction, ok := effect["healingReduction"].(float64); ok {
			duration := 3
			if d, ok := effect["debuffDuration"].(float64); ok {
				duration = int(d)
			}
			// 应用到目标敌人
			if target != nil && target.HP > 0 {
				m.buffManager.ApplyEnemyDebuff(target.ID, "mortal_strike", "致死打击", "debuff", duration, healingReduction, "healing_received", "")
			}
		}
	}
}

// handlePassiveReflectEffects 处理被动技能的反射效果
func (m *BattleManager) handlePassiveReflectEffects(character *models.Character, attacker *models.Monster, damageTaken int, session *BattleSession, logs *[]models.BattleLog) {
	if m.passiveSkillManager == nil {
		return
	}

	passives := m.passiveSkillManager.GetPassiveSkills(character.ID)
	for _, passive := range passives {
		if passive.Passive.EffectType == "reflect" && passive.Passive.ID == "warrior_passive_shield_reflection" {
			// 盾牌反射（被动）：受到物理攻击时反射伤害
			reflectPercent := passive.EffectValue // 百分比值（如10.0表示10%）
			reflectDamage := int(float64(damageTaken) * reflectPercent / 100.0)
			if reflectDamage > 0 {
				attacker.HP -= reflectDamage
				if attacker.HP < 0 {
					attacker.HP = 0
				}
				m.addLog(session, "combat", fmt.Sprintf("%s 的盾牌反射对 %s 造成 %d 点反射伤害！", character.Name, attacker.Name, reflectDamage), "#ff8800")
				*logs = append(*logs, session.BattleLogs[len(session.BattleLogs)-1])
			}
		}
	}
}
