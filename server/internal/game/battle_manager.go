package game

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"text-wow/internal/models"
	"text-wow/internal/repository"
)

// BattleManager ææç®¡çå?- ç®¡çææç¨æ·çææç¶æ?type BattleManager struct {
	mu                  sync.RWMutex
	sessions            map[int]*BattleSession // key: userID
	gameRepo            *repository.GameRepository
	charRepo            *repository.CharacterRepository
	skillManager        *SkillManager
	buffManager         *BuffManager
	passiveSkillManager *PassiveSkillManager
}

// BattleSession ç¨æ·ææä¼è¯
type BattleSession struct {
	UserID             int
	IsRunning          bool
	CurrentZone        *models.Zone
	CurrentEnemy       *models.Monster   // ä¿çç¨äºååå¼å®¹
	CurrentEnemies     []*models.Monster // å¤ä¸ªæäººæ¯æ
	BattleLogs         []models.BattleLog
	BattleCount        int
	SessionKills       int
	SessionGold        int
	SessionExp         int
	StartedAt          time.Time
	LastTick           time.Time
	IsResting          bool       // æ¯å¦å¨ä¼æ?	RestUntil          *time.Time // ä¼æ¯ç»ææ¶é´
	RestStartedAt      *time.Time // ä¼æ¯å¼å§æ¶é?	LastRestTick       *time.Time // ä¸æ¬¡æ¢å¤å¤ççæ¶é?	RestSpeed          float64    // æ¢å¤éåº¦åç
	CurrentBattleExp   int        // æ¬åºææè·å¾çç»éª?	CurrentBattleGold  int        // æ¬åºææè·å¾çéå¸?	CurrentBattleKills int        // æ¬åºææå»ææ?	CurrentTurnIndex   int        // ååæ§å¶ï¼?1=ç©å®¶ååï¼?=0=æäººç´¢å¼
	JustEncountered    bool       // åé­éæäººï¼éè¦ç­å¾?ä¸ªtickåå¼å§ææ?}

// NewBattleManager åå»ºææç®¡çå?func NewBattleManager() *BattleManager {
	return &BattleManager{
		sessions:            make(map[int]*BattleSession),
		gameRepo:            repository.NewGameRepository(),
		charRepo:            repository.NewCharacterRepository(),
		skillManager:        NewSkillManager(),
		buffManager:         NewBuffManager(),
		passiveSkillManager: NewPassiveSkillManager(),
	}
}

// å¨å±ææç®¡çå¨å®ä¾?var battleManager *BattleManager
var once sync.Once

// GetBattleManager è·åææç®¡çå¨åä¾?func GetBattleManager() *BattleManager {
	once.Do(func() {
		battleManager = NewBattleManager()
	})
	return battleManager
}

// GetOrCreateSession è·åæåå»ºææä¼è¯?func (m *BattleManager) GetOrCreateSession(userID int) *BattleSession {
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
		CurrentTurnIndex: -1, // åå§åä¸ºç©å®¶åå
		RestSpeed:        1.0,
	}
	m.sessions[userID] = session
	return session
}

// GetSession è·åææä¼è¯
func (m *BattleManager) GetSession(userID int) *BattleSession {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.sessions[userID]
}

// ToggleBattle åæ¢ææç¶æ?func (m *BattleManager) ToggleBattle(userID int) (bool, error) {
	session := m.GetOrCreateSession(userID)

	m.mu.Lock()
	defer m.mu.Unlock()

	session.IsRunning = !session.IsRunning
	session.LastTick = time.Now()

	if session.IsRunning {
		// å¦ææ²¡æè®¾ç½®åºåï¼è®¾ç½®é»è®¤åºå?		if session.CurrentZone == nil {
			zone, err := m.gameRepo.GetZoneByID("elwynn")
			if err == nil {
				session.CurrentZone = zone
			}
		}
		session.CurrentTurnIndex = -1 // éç½®ä¸ºç©å®¶åå?		m.addLog(session, "system", ">> å¼å§èªå¨ææ?..", "#33ff33")
	} else {
		m.addLog(session, "system", ">> æåèªå¨ææ", "#ffff00")
	}

	return session.IsRunning, nil
}

// StartBattle å¼å§ææ?func (m *BattleManager) StartBattle(userID int) (bool, error) {
	session := m.GetOrCreateSession(userID)

	m.mu.Lock()
	defer m.mu.Unlock()

	if session.IsRunning {
		return true, nil
	}

	session.IsRunning = true
	session.LastTick = time.Now()
	session.CurrentTurnIndex = -1 // éç½®ä¸ºç©å®¶åå?
	// è®¾ç½®é»è®¤åºå
	if session.CurrentZone == nil {
		zone, err := m.gameRepo.GetZoneByID("elwynn")
		if err == nil {
			session.CurrentZone = zone
		}
	}

	m.addLog(session, "system", ">> å¼å§èªå¨ææ?..", "#33ff33")
	return true, nil
}

// StopBattle åæ­¢ææ
func (m *BattleManager) StopBattle(userID int) error {
	session := m.GetSession(userID)
	if session == nil {
		return nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	session.IsRunning = false
	m.addLog(session, "system", ">> æåèªå¨ææ", "#ffff00")
	return nil
}

// ExecuteBattleTick æ§è¡ææååï¼ååå¶ï¼æ¯tickåªæ§è¡ä¸ä¸ªå¨ä½ï¼
func (m *BattleManager) ExecuteBattleTick(userID int, characters []*models.Character) (*BattleTickResult, error) {
	session := m.GetOrCreateSession(userID)

	m.mu.Lock()
	defer m.mu.Unlock()

	// å¦ææ²¡æè§è²ï¼è¿ånil
	if len(characters) == 0 {
		return nil, nil
	}

	// ä½¿ç¨ç¬¬ä¸ä¸ªè§è²è¿è¡ææ?	char := characters[0]

	// ç¡®ä¿æå£«çææ°ä¸éä¸?00ï¼æ¯æ¬¡tické½æ£æ¥ï¼é²æ­¢è¢«è¦çï¼
	if char.ResourceType == "rage" {
		char.MaxResource = 100
	}

	// å è½½è§è²çæè½ï¼å¦æè¿æ²¡æå è½½ï¼
	if m.skillManager != nil {
		if err := m.skillManager.LoadCharacterSkills(char.ID); err != nil {
			// å¦æå è½½å¤±è´¥ï¼è®°å½æ¥å¿ä½ä¸ä¸­æ­ææ?			m.addLog(session, "system", fmt.Sprintf("è­¦åï¼æ æ³å è½½è§è²æè? %v", err), "#ffaa00")
		}
	}

	// å è½½è§è²çè¢«å¨æè½ï¼å¦æè¿æ²¡æå è½½ï¼
	if m.passiveSkillManager != nil {
		if err := m.passiveSkillManager.LoadCharacterPassiveSkills(char.ID); err != nil {
			// å¦æå è½½å¤±è´¥ï¼è®°å½æ¥å¿ä½ä¸ä¸­æ­ææ?			m.addLog(session, "system", fmt.Sprintf("è­¦åï¼æ æ³å è½½è§è²è¢«å¨æè? %v", err), "#ffaa00")
		}
	}

	// å¦ææææªè¿è¡ä¸ä¸å¨ä¼æ¯ç¶æï¼æ£æ¥æ¯å¦éè¦è¿åè§è²æ°æ?	// å¦æè§è²åå¤æ´»ï¼ä¹åæ­»äº¡ä½ç°å¨ä¸æ­»äº¡ï¼ï¼éè¦è¿åä¸æ¬¡æ°æ®è®©åç«¯æ´æ°
	if !session.IsRunning && !session.IsResting {
		// ä»æ°æ®åºéæ°å è½½è§è²æ°æ®ä»¥ç¡®ä¿ç¶ææ­£ç¡?		updatedChar, err := m.charRepo.GetByID(char.ID)
		if err == nil && updatedChar != nil {
			char = updatedChar
			// ç¡®ä¿æå£«çææ°ä¸éä¸?00
			if char.ResourceType == "rage" {
				char.MaxResource = 100
			}
			// å¦æè§è²å·²ç»å¤æ´»ï¼ä¹åæ­»äº¡ä½ç°å¨ä¸æ­»äº¡ï¼ï¼è¿åè§è²æ°æ?			if !char.IsDead {
				// è¿åè§è²æ°æ®ï¼è®©åç«¯ç¥éè§è²å·²ç»å¤æ´»
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
		// å¦ææ æ³è·åè§è²æ°æ®æè§è²ä»ç¶æ­»äº¡ï¼è¿ånil
		return nil, nil
	}

	session.LastTick = time.Now()
	logs := make([]models.BattleLog, 0)

	// æ£æ¥è§è²æ¯å¦æ­»äº¡ä¸è¿æ²¡å°å¤æ´»æ¶é?	now := time.Now()
	if char.IsDead && char.ReviveAt != nil && now.Before(*char.ReviveAt) {
		// è§è²æ­»äº¡ä½è¿æ²¡å°å¤æ´»æ¶é´ï¼è¿å¥ä¼æ¯ç¶æ?		if !session.IsResting {
			// è®¡ç®ä¼æ¯æ¶é´ï¼å¤æ´»æ¶é?+ æ¢å¤æ¶é´ï¼?			reviveRemaining := char.ReviveAt.Sub(now)
			recoveryTime := 25 * time.Second // æ¢å¤ä¸åHPéè¦çæ¶é´
			restDuration := reviveRemaining + recoveryTime
			restUntil := now.Add(restDuration)
			session.IsResting = true
			session.RestUntil = &restUntil
			session.RestStartedAt = &now
			session.LastRestTick = &now
			session.RestSpeed = 1.0
			// ä¿æ isRunning = trueï¼è¿æ ·æé®ä¼æ¾ç¤º"åæ­¢ææº"ï¼ä¼æ¯ç¶æå¯ä»¥èªå¨å¤ç?
			remainingSeconds := int(reviveRemaining.Seconds()) + 1
			m.addLog(session, "death", fmt.Sprintf("%s æ­£å¨å¤æ´»ä¸?.. (å©ä½ %d ç§?", char.Name, remainingSeconds), "#ff0000")
			logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
		}
	}

	// å¦ææ­£å¨ä¼æ¯ï¼å¤çä¼æ?	if session.IsResting && session.RestUntil != nil {
		initialHP := char.HP
		initialMP := char.Resource
		now := time.Now()
		m.processRest(session, char)

		// æ´æ°LastTickï¼ç¨äºä¸æ¬¡è®¡ç®æ¶é´å·®
		session.LastTick = now

		if !session.IsResting {
			// ä¼æ¯ç»æï¼ä¿å­è§è²æ°æ?			m.charRepo.UpdateAfterBattle(char.ID, char.HP, char.Resource, char.Exp, char.Level,
				char.ExpToNext, char.MaxHP, char.MaxResource, char.PhysicalAttack, char.MagicAttack, char.PhysicalDefense, char.MagicDefense,
				char.Strength, char.Agility, char.Stamina, char.TotalKills)

			// ä¼æ¯ç»æåï¼ç¡®ä¿è¿åè§è²æ°æ®ï¼è®©åç«¯ç¥éä¼æ¯å·²ç»æ?			// ä»æ°æ®åºéæ°å è½½è§è²æ°æ®ä»¥ç¡®ä¿ç¶ææ­£ç¡?			updatedChar, err := m.charRepo.GetByID(char.ID)
			if err == nil && updatedChar != nil {
				char = updatedChar
				// ç¡®ä¿æå£«çææ°ä¸éä¸?00
				if char.ResourceType == "rage" {
					char.MaxResource = 100
				}
			}

			// å¦æè§è²å·²ç»å¤æ´»ï¼ä¸åæ­»äº¡ï¼ï¼èªå¨æ¢å¤ææ?			if !char.IsDead {
				session.IsRunning = true
				m.addLog(session, "system", ">> ä¼æ¯ç»æï¼èªå¨æ¢å¤ææ?, "#33ff33")
			} else {
				m.addLog(session, "system", ">> ä¼æ¯ç»æï¼åå¤ä¸ä¸åºææ?, "#33ff33")
			}
			logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
		} else {
			// ä»å¨ä¼æ¯ä¸?			remaining := session.RestUntil.Sub(time.Now())
			if remaining > 0 {
				m.addLog(session, "system", fmt.Sprintf(">> ä¼æ¯ä¸?.. (å©ä½ %d ç§?", int(remaining.Seconds())+1), "#888888")
				logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
			}
		}

		// ä¿å­è§è²æ°æ®æ´æ°
		if char.HP != initialHP || char.Resource != initialMP {
			m.charRepo.UpdateAfterBattle(char.ID, char.HP, char.Resource, char.Exp, char.Level,
				char.ExpToNext, char.MaxHP, char.MaxResource, char.PhysicalAttack, char.MagicAttack, char.PhysicalDefense, char.MagicDefense,
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

	// è·åå­æ´»çæäº?	aliveEnemies := make([]*models.Monster, 0)
	for _, enemy := range session.CurrentEnemies {
		if enemy != nil && enemy.HP > 0 {
			aliveEnemies = append(aliveEnemies, enemy)
		}
	}

	// å¦ææ²¡ææäººï¼çææ°ç?	if len(aliveEnemies) == 0 {
		// éç½®æ¬åºææç»è®¡
		session.CurrentBattleExp = 0
		session.CurrentBattleGold = 0
		session.CurrentBattleKills = 0
		session.CurrentTurnIndex = -1 // ç©å®¶åå

		// ææå¼å§æ¶ï¼ç¡®ä¿æå£«çææ°ä¸?ï¼æå¤§ææ°ä¸?00
		if char.ResourceType == "rage" {
			char.Resource = 0
			char.MaxResource = 100
		}

		err := m.spawnEnemies(session, char.Level)
		if err != nil {
			return nil, err
		}
		logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])

		// æ è®°åé­éæäººï¼éè¦ç­å¾?ä¸ªtickåå¼å§ææ?		session.JustEncountered = true

		// æ´æ°å­æ´»æäººåè¡¨
		aliveEnemies = session.CurrentEnemies

		// åé­éæäººï¼è¿ä¸ªtickåªæ¾ç¤ºä¿¡æ¯ï¼ä¸æ§è¡ææ?		return &BattleTickResult{
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

	// å¦æåé­éæäººï¼è¿ä¸ªtickåªæ¾ç¤ºä¿¡æ¯ï¼ä¸æ§è¡ææ?	if session.JustEncountered {
		session.JustEncountered = false // æ¸é¤æ å¿ï¼ä¸ä¸ä¸ªtickå¼å§ææ?		return &BattleTickResult{
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

	// ååå¶é»è¾ï¼CurrentTurnIndex == -1 è¡¨ç¤ºç©å®¶ååï¼?=0 è¡¨ç¤ºæäººç´¢å¼
	if session.CurrentTurnIndex == -1 {
		// ç©å®¶ååï¼æ»å»ç¬¬ä¸ä¸ªå­æ´»çæäºº
		if len(aliveEnemies) > 0 {
			target := aliveEnemies[0]
			targetHPPercent := float64(target.HP) / float64(target.MaxHP)
			hasMultipleEnemies := len(aliveEnemies) > 1

			// ä½¿ç¨æè½ç®¡çå¨éæ©æè?			var skillState *CharacterSkillState
			if m.skillManager != nil {
				skillState = m.skillManager.SelectBestSkill(char.ID, char.Resource, targetHPPercent, hasMultipleEnemies, m.buffManager)
			}

			var skillName string
			var playerDamage int
			var resourceCost int
			var resourceGain int
			var usedSkill bool
			var skillEffects map[string]interface{}
			var isCrit bool
			var damageDetails *DamageCalculationDetails
			var shouldDealDamage bool // æ¯å¦åºè¯¥é æä¼¤å®³ï¼åªæattackç±»åçæè½æé æä¼¤å®³ï¼?
			if skillState != nil && skillState.Skill != nil {
				// ä½¿ç¨æè?				skillName = skillState.Skill.Name
				resourceCost = m.skillManager.GetSkillResourceCost(skillState)

				// å¤æ­æè½æ¯å¦åºè¯¥é æä¼¤å®³ï¼åªæattackç±»åçæè½æé æä¼¤å®³ï¼?				shouldDealDamage = skillState.Skill.Type == "attack"

				// æ£æ¥èµæºæ¯å¦è¶³å¤?				if resourceCost <= char.Resource {
					
					var baseDamage int
					// playerDamage, isCrit, and damageDetails are already declared in outer scope
					// Do not redeclare them here to avoid shadowing outer scope variables
					if shouldDealDamage {
						// è®¡ç®æè½ä¼¤å®³ï¼åºç¡ä¼¤å®³ï¼æ´å»å¨åé¢å¤çï¼?						baseDamage = m.skillManager.CalculateSkillDamage(skillState, char, target, m.passiveSkillManager, m.buffManager)
						
						// åå»ºæè½ä¼¤å®³è¯¦æï¼ç®åçï¼?						damageDetails = &DamageCalculationDetails{
							BaseAttack:      char.PhysicalAttack,
							BaseDefense:     target.PhysicalDefense,
							BaseDamage:      float64(baseDamage),
							AttackModifiers: []string{fmt.Sprintf("æè½åç: %.1f", skillState.Skill.ScalingRatio)},
							DefenseModifiers: []string{},
							ActualCritRate:  -1, // -1 è¡¨ç¤ºæªè®¾ç½?							RandomRoll:      -1, // -1 è¡¨ç¤ºæªè®¾ç½?						}

						// è®¡ç®æ´å»ï¼æè½ä¹å¯ä»¥æ´å»ï¼åºç¨è¢«å¨æè½åBuffå æï¼?						actualCritRate := char.CritRate
						damageDetails.BaseCritRate = char.CritRate
						damageDetails.CritModifiers = []string{}
						
						if m.passiveSkillManager != nil {
							critModifier := m.passiveSkillManager.GetPassiveModifier(char.ID, "crit_rate")
							if critModifier > 0 {
								actualCritRate = char.CritRate + critModifier/100.0
								damageDetails.CritModifiers = append(damageDetails.CritModifiers, 
									fmt.Sprintf("è¢«å¨æ´å»+%.0f%%", critModifier))
							}
						}
						// åºç¨Buffçæ´å»çå æï¼é²è½ç­ï¼?						if m.buffManager != nil {
							critBuffValue := m.buffManager.GetBuffValue(char.ID, "crit_rate")
							if critBuffValue > 0 {
								actualCritRate = actualCritRate + critBuffValue/100.0
								damageDetails.CritModifiers = append(damageDetails.CritModifiers, 
									fmt.Sprintf("Buffæ´å»+%.0f%%", critBuffValue))
							}
						}
						if actualCritRate > 1.0 {
							actualCritRate = 1.0
						}
						damageDetails.ActualCritRate = actualCritRate
						randomRoll := rand.Float64()
						damageDetails.RandomRoll = randomRoll
						isCrit = randomRoll < actualCritRate
						damageDetails.IsCrit = isCrit
						damageDetails.CritMultiplier = char.CritDamage
						
						if isCrit {
							playerDamage = int(float64(baseDamage) * char.CritDamage)
						} else {
							playerDamage = baseDamage
						}
						damageDetails.FinalDamage = playerDamage
					}

					// åºç¨æè½ææ?					skillEffects = m.skillManager.ApplySkillEffects(skillState, char, target)

					// åºç¨Buff/Debuffææ
					m.applySkillBuffs(skillState, char, target, skillEffects)

					// åºç¨Debuffå°æäººï¼æ«å¿æå¼ãæé£æ©ç­ï¼
					m.applySkillDebuffs(skillState, char, target, aliveEnemies, skillEffects)

					// æ¶èèµæº?					char.Resource -= resourceCost
					if char.Resource < 0 {
						char.Resource = 0
					}

					// ä½¿ç¨æè½ï¼è®¾ç½®å·å´ï¼?					m.skillManager.UseSkill(char.ID, skillState.SkillID)
					usedSkill = true

					// å¤çæè½ç¹æ®ææï¼ææ°è·å¾ç­ï¼
					if rageGain, ok := skillEffects["rageGain"].(int); ok {
						// åºç¨è¢«å¨æè½çææ°è·å¾å æï¼æ¤æææ¡ç­ï¼?						actualRageGain := m.applyRageGenerationModifiers(char.ID, rageGain)
						char.Resource += actualRageGain
						resourceGain = actualRageGain
						if char.Resource > char.MaxResource {
							char.Resource = char.MaxResource
						}
					}

					// åªæattackç±»åçæè½æé æä¼¤å®³
					if shouldDealDamage {
						// å¤çAOEæè½ï¼æé£æ©ç­ï¼?						if skillState.Skill.TargetType == "enemy_all" {
							// å¯¹æææäººé æä¼¤å®³
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
							// playerDamageç¨äºæ¥å¿æ¾ç¤ºï¼ä¸»ç®æ ä¼¤å®³ï¼?						} else if skillState.SkillID == "warrior_cleave" {
							// é¡ºåæ©ï¼ä¸»ç®æ ?ç¸é»ç®æ 
							target.HP -= playerDamage

							// å¯¹ç¸é»ç®æ é æä¼¤å®³ï¼æå¤?ä¸ªï¼
							adjacentCount := 0
							for _, enemy := range aliveEnemies {
								if enemy != target && enemy.HP > 0 && adjacentCount < 2 {
									// è®¡ç®ç¸é»ç®æ ä¼¤å®³
									if effect, ok := skillState.Effect["adjacentMultiplier"].(float64); ok {
										adjacentDamage := int(float64(char.PhysicalAttack) * effect)
										// åºç¡ä¼¤å®³ = å®éæ»å»å?- ç®æ é²å¾¡åï¼ä¸åé¤ä»¥2ï¼?										adjacentDamage = adjacentDamage - enemy.PhysicalDefense
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
										m.addLog(session, "combat", fmt.Sprintf("%s çé¡ºåæ©æ³¢åå?%sï¼é æ %d ç¹ä¼¤å®?, char.Name, enemy.Name, adjacentDamage), "#ffaa00")
										logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
									}
								}
							}
						} else {
							// åä½æè?							target.HP -= playerDamage
						}
					}
				} else {
					// èµæºä¸è¶³ï¼ä½¿ç¨æ®éæ»å?					skillState = nil
				}
			}

			// å¦ææ²¡æä½¿ç¨æè½æèµæºä¸è¶³ï¼ä½¿ç¨æ®éæ»å?			if skillState == nil {
				skillName = "æ®éæ»å?
				// è®¡ç®å®éç©çæ»å»åï¼åºç¨è¢«å¨æè½å æï¼
				actualAttack := float64(char.PhysicalAttack)
				damageDetails = &DamageCalculationDetails{
					BaseAttack:      char.PhysicalAttack,
					BaseDefense:     target.PhysicalDefense,
					AttackModifiers: []string{},
					DefenseModifiers: []string{},
					ActualCritRate:  -1, // -1 è¡¨ç¤ºæªè®¾ç½?					RandomRoll:      -1, // -1 è¡¨ç¤ºæªè®¾ç½?				}
				
				if m.passiveSkillManager != nil {
					attackModifier := m.passiveSkillManager.GetPassiveModifier(char.ID, "attack")
					if attackModifier > 0 {
						actualAttack = actualAttack * (1.0 + attackModifier/100.0)
						damageDetails.AttackModifiers = append(damageDetails.AttackModifiers, 
							fmt.Sprintf("è¢«å¨æ»å»+%.0f%%", attackModifier))
					}
					// åºç¨è¢«å¨æè½çä¼¤å®³å æ
					damageModifier := m.passiveSkillManager.GetPassiveModifier(char.ID, "damage")
					if damageModifier > 0 {
						actualAttack = actualAttack * (1.0 + damageModifier/100.0)
						damageDetails.AttackModifiers = append(damageDetails.AttackModifiers, 
							fmt.Sprintf("è¢«å¨ä¼¤å®³+%.0f%%", damageModifier))
					}

					// å¤çä½è¡éæ¶çæ»å»åå æï¼çæ´ä¹å¿ï¼
					hpPercent := float64(char.HP) / float64(char.MaxHP)
					passives := m.passiveSkillManager.GetPassiveSkills(char.ID)
					for _, passive := range passives {
						if passive.Passive.EffectType == "stat_mod" && passive.Passive.ID == "warrior_passive_berserker_heart" {
							// æ ¹æ®ç­çº§è®¡ç®è§¦åéå¼ï¼1çº?0%ï¼?çº?0%ï¼?							threshold := 0.50 - float64(passive.Level-1)*0.05
							if hpPercent < threshold {
								// æ ¹æ®ç­çº§è®¡ç®æ»å»åå æï¼1çº?0%ï¼?çº?0%ï¼?								attackBonus := 20.0 + float64(passive.Level-1)*10.0
								actualAttack = actualAttack * (1.0 + attackBonus/100.0)
								damageDetails.AttackModifiers = append(damageDetails.AttackModifiers, 
									fmt.Sprintf("çæ´ä¹å¿+%.0f%%", attackBonus))
							}
						}
					}
				}
				// åºç¨Buffçæ»å»åå æï¼æææå¼ãçæ´ä¹æãå¤©ç¥ä¸å¡ç­ï¼?				if m.buffManager != nil {
					attackBuffValue := m.buffManager.GetBuffValue(char.ID, "attack")
					if attackBuffValue > 0 {
						actualAttack = actualAttack * (1.0 + attackBuffValue/100.0)
						damageDetails.AttackModifiers = append(damageDetails.AttackModifiers, 
							fmt.Sprintf("Buffæ»å»+%.0f%%", attackBuffValue))
					}
				}
				
				damageDetails.ActualAttack = actualAttack
				damageDetails.ActualDefense = float64(target.PhysicalDefense)
				
				baseDamage, calcDetails := m.calculatePhysicalDamageWithDetails(int(actualAttack), target.PhysicalDefense)
				damageDetails.BaseDamage = calcDetails.BaseDamage
				damageDetails.Variance = calcDetails.Variance
				
				// è®¡ç®æ´å»çï¼åºç¨è¢«å¨æè½åBuffå æï¼?				actualCritRate := char.CritRate
				damageDetails.BaseCritRate = char.CritRate
				damageDetails.CritModifiers = []string{}
				
				if m.passiveSkillManager != nil {
					critModifier := m.passiveSkillManager.GetPassiveModifier(char.ID, "crit_rate")
					if critModifier > 0 {
						actualCritRate = char.CritRate + critModifier/100.0
						damageDetails.CritModifiers = append(damageDetails.CritModifiers, 
							fmt.Sprintf("è¢«å¨æ´å»+%.0f%%", critModifier))
					}
				}
				// åºç¨Buffçæ´å»çå æï¼é²è½ç­ï¼?				if m.buffManager != nil {
					critBuffValue := m.buffManager.GetBuffValue(char.ID, "crit_rate")
					if critBuffValue > 0 {
						actualCritRate = actualCritRate + critBuffValue/100.0
						damageDetails.CritModifiers = append(damageDetails.CritModifiers, 
							fmt.Sprintf("Buffæ´å»+%.0f%%", critBuffValue))
					}
				}
				if actualCritRate > 1.0 {
					actualCritRate = 1.0
				}
				damageDetails.ActualCritRate = actualCritRate
				randomRoll := rand.Float64()
				damageDetails.RandomRoll = randomRoll
				isCrit = randomRoll < actualCritRate
				damageDetails.IsCrit = isCrit
				damageDetails.CritMultiplier = char.CritDamage
				
				if isCrit {
					playerDamage = int(float64(baseDamage) * char.CritDamage)
				} else {
					playerDamage = baseDamage
				}
				damageDetails.FinalDamage = playerDamage
				target.HP -= playerDamage
				resourceCost = 0
				usedSkill = false
			}
			// å¦æä½¿ç¨äºæè½ï¼isCritå·²ç»å¨ä¸é¢è®¡ç®äº

			// æ®éæ»å»è·å¾ææ°ï¼åªææ®éæ»å»æè·å¾ææ°ï¼ä½¿ç¨æè½æ¶ä¸è·å¾ï¼
			if char.ResourceType == "rage" && !usedSkill {
				var baseRageGain int
				if isCrit {
					baseRageGain = 10 // æ´å»è·å¾10ç¹ææ°
				} else {
					baseRageGain = 5 // æ®éæ»å»è·å¾?ç¹ææ°
				}

				// åºç¨è¢«å¨æè½çææ°è·å¾å æï¼æ¤æææ¡ç­ï¼?				rageGain := m.applyRageGenerationModifiers(char.ID, baseRageGain)

				char.Resource += rageGain
				resourceGain = rageGain
				// ç¡®ä¿ä¸è¶è¿æå¤§å?				if char.Resource > char.MaxResource {
					char.Resource = char.MaxResource
				}
			}

			// å¤çè¢«å¨æè½çç¹æ®ææï¼æ»å»æ¶è§¦åï¼?			m.handlePassiveOnHitEffects(char, playerDamage, usedSkill, session, &logs)

			// æå»ºæææ¥å¿æ¶æ¯ï¼åå«èµæºååï¼å¸¦é¢è²ï¼
			resourceChangeText := m.formatResourceChange(char.ResourceType, resourceCost, resourceGain)
			
			// æ ¼å¼åä¼¤å®³å¬å¼?			formulaText := ""
			if damageDetails != nil {
				formulaText = m.formatDamageFormula(damageDetails)
			}

			// å¤çæè½ç¹æ®æææ¥å¿?			if skillEffects != nil {
				if stun, ok := skillEffects["stun"].(bool); ok && stun {
					m.addLog(session, "combat", fmt.Sprintf("%s è¢«ç©æäºï¼?, target.Name), "#ff00ff")
					logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
				}
				// å¤çåºäºä¼¤å®³çæ¢å¤ï¼åè¡ç­ï¼
				if healPercent, ok := skillEffects["healPercent"].(float64); ok && usedSkill {
					healAmount := int(float64(playerDamage) * healPercent / 100.0)
					char.HP += healAmount
					if char.HP > char.MaxHP {
						char.HP = char.MaxHP
					}
					m.addLog(session, "heal", fmt.Sprintf("%s æ¢å¤äº?%d ç¹çå½å?, char.Name, healAmount), "#00ff00")
					logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
				}
				// å¤çç ´éæ²èçç«å³æ¢å¤ï¼åºäºæå¤§HPï¼?				if healMaxHpPercent, ok := skillEffects["healMaxHpPercent"].(float64); ok && usedSkill {
					healAmount := int(float64(char.MaxHP) * healMaxHpPercent / 100.0)
					char.HP += healAmount
					if char.HP > char.MaxHP {
						char.HP = char.MaxHP
					}
					m.addLog(session, "heal", fmt.Sprintf("%s çç ´éæ²èæ¢å¤äº %d ç¹çå½å?, char.Name, healAmount), "#00ff00")
					logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
				}
			}

			// è®°å½æè½ä½¿ç¨æ¥å¿?			if shouldDealDamage {
				// æ»å»ç±»æè½ï¼è®°å½ä¼¤å®³
				if isCrit {
					m.addLog(session, "combat", fmt.Sprintf("%s ä½¿ç¨ [%s] ð¥æ´å»ï¼å¯¹ %s é æ %d ç¹ä¼¤å®?s%s", char.Name, skillName, target.Name, playerDamage, formulaText, resourceChangeText), "#ff6b6b")
				} else {
					m.addLog(session, "combat", fmt.Sprintf("%s ä½¿ç¨ [%s] å¯?%s é æ %d ç¹ä¼¤å®?s%s", char.Name, skillName, target.Name, playerDamage, formulaText, resourceChangeText), "#ffaa00")
				}
			} else {
				// éæ»å»ç±»æè½ï¼buff/debuff/controlç­ï¼ï¼åªè®°å½ä½¿ç¨ï¼ä¸è®°å½ä¼¤å®³
				m.addLog(session, "combat", fmt.Sprintf("%s ä½¿ç¨ [%s]%s", char.Name, skillName, resourceChangeText), "#8888ff")
			}
			logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])

			// åå°æè½å·å´æ¶é?			m.skillManager.TickCooldowns(char.ID)

			// åå°Buff/Debuffæç»­æ¶é´
			expiredBuffs := m.buffManager.TickBuffs(char.ID)
			for _, effectID := range expiredBuffs {
				m.addLog(session, "buff", fmt.Sprintf("%s ç?%s æææ¶å¤±äº?, char.Name, effectID), "#888888")
				logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
			}

			// æ£æ¥ç®æ æ¯å¦æ­»äº?			if target.HP <= 0 {
				// ç¡®ä¿HPå½é¶
				target.HP = 0

				// å¤çæäºæºå¨çå»æåæææ?				m.handleWarMachineRageGain(char, session, &logs)

				// æäººæ­»äº¡
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

				// æ£æ¥åçº?				for char.Exp >= char.ExpToNext {
					char.Exp -= char.ExpToNext
					char.Level++
					char.ExpToNext = int(float64(char.ExpToNext) * 1.5)

					// åçº§å±æ§æå?					char.MaxHP += 15
					char.HP = char.MaxHP

					// æå£«çææ°æå¤§å¼åºå®ä¸º100ï¼ä¸éåçº§æ¹å?					if char.ResourceType == "rage" {
						char.MaxResource = 100
						// åçº§æ¶ææ°ä¿æä¸åï¼ä¸éç½®ä¸ºæå¤§å?					} else {
						char.MaxResource += 8
						char.Resource = char.MaxResource
					}

					char.Strength += 2
					char.Agility += 1
					char.Stamina += 2
					char.Intellect += 1
					char.Spirit += 1
					char.PhysicalAttack = char.Strength / 2
					char.MagicAttack = char.Intellect / 2
					char.PhysicalDefense = char.Stamina / 3
					char.MagicDefense = (char.Intellect + char.Spirit) / 4

					m.addLog(session, "levelup", fmt.Sprintf("ðãåçº§ãæ­åï¼%s åå°äº?%d çº§ï¼", char.Name, char.Level), "#ffd700")
					logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
				}
			}

			// ç§»å¨å°ä¸ä¸ä¸ªæäººåå?			session.CurrentTurnIndex = 0
		}
	} else {
		// æäººååï¼å½åç´¢å¼çæäººæ»å»ç©å®¶
		if session.CurrentTurnIndex < len(aliveEnemies) {
			enemy := aliveEnemies[session.CurrentTurnIndex]
			// æäººé»è®¤ä½¿ç¨ç©çæ»å»
			baseEnemyDamage, enemyDamageDetails := m.calculatePhysicalDamageWithDetails(enemy.PhysicalAttack, char.PhysicalDefense)
			enemyDamageDetails.BaseAttack = enemy.PhysicalAttack
			enemyDamageDetails.BaseDefense = char.PhysicalDefense
			enemyDamageDetails.AttackModifiers = []string{}
			enemyDamageDetails.DefenseModifiers = []string{}

			// åºç¨buff/debuffææï¼å¦ç¾çæ ¼æ¡çåä¼¤ç­ï¼?			originalDamage := baseEnemyDamage
			enemyDamage := m.buffManager.CalculateDamageTakenWithBuffs(baseEnemyDamage, char.ID, true)
			if enemyDamage != originalDamage {
				reduction := float64(originalDamage-enemyDamage) / float64(originalDamage) * 100.0
				enemyDamageDetails.DefenseModifiers = append(enemyDamageDetails.DefenseModifiers, 
					fmt.Sprintf("åä¼¤Buff -%.0f%%", reduction))
			}

			// å¤çè¢«å¨æè½çåä¼¤ææï¼ä¸ç­æå¿ç­ï¼?			originalDamage2 := enemyDamage
			enemyDamage = m.handlePassiveDamageReduction(char, enemyDamage)
			if enemyDamage != originalDamage2 {
				reduction := float64(originalDamage2-enemyDamage) / float64(originalDamage2) * 100.0
				enemyDamageDetails.DefenseModifiers = append(enemyDamageDetails.DefenseModifiers, 
					fmt.Sprintf("è¢«å¨åä¼¤ -%.0f%%", reduction))
			}
			enemyDamageDetails.FinalDamage = enemyDamage

			// å¤çæ¤ç¾ææï¼ä¸ç­å£åç­ï¼?			shieldAmount := m.buffManager.GetBuffValue(char.ID, "shield")
			if shieldAmount > 0 {
				// ææ¤ç¾ï¼åæ¶èæ¤ç?				shieldInt := int(shieldAmount)
				if enemyDamage <= shieldInt {
					// ä¼¤å®³å®å¨è¢«æ¤ç¾å¸æ?					shieldInt -= enemyDamage
					absorbedDamage := enemyDamage
					enemyDamage = 0
					m.addLog(session, "shield", fmt.Sprintf("%s çæ¤ç¾å¸æ¶äº %d ç¹ä¼¤å®?, char.Name, absorbedDamage), "#00ffff")
					logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
					// æ´æ°æ¤ç¾å¼ï¼éè¿æ´æ°Buffçvalueï¼?					m.updateShieldValue(char.ID, float64(shieldInt))
				} else {
					// æ¤ç¾è¢«å»ç ´ï¼å©ä½ä¼¤å®³ç»§ç»­
					absorbedDamage := shieldInt
					enemyDamage -= shieldInt
					m.addLog(session, "shield", fmt.Sprintf("%s çæ¤ç¾å¸æ¶äº %d ç¹ä¼¤å®³åè¢«å»ç ?, char.Name, absorbedDamage), "#00ffff")
					logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
					m.updateShieldValue(char.ID, 0)
				}
			}

			// å¤çè¢«å¨æè½ççå­ææï¼åé§ä¸æç­ï¼? å¨åå°ä¼¤å®³åæ£æ?			originalHP := char.HP
			char.HP -= enemyDamage

			// å¦æåå°è´å½ä¼¤å®³ï¼æ£æ¥åé§ä¸æææ?			if originalHP > 0 && char.HP <= 0 {
				if m.passiveSkillManager != nil {
					passives := m.passiveSkillManager.GetPassiveSkills(char.ID)
					for _, passive := range passives {
						if passive.Passive.EffectType == "survival" && passive.Passive.ID == "warrior_passive_unbreakable" {
							// åé§ä¸æï¼åå°è´å½ä¼¤å®³æ¶ä¿ç1ç¹HP
							char.HP = 1
							m.addLog(session, "survival", fmt.Sprintf("%s çåé§ä¸æææè§¦åï¼ä¿çäº?ç¹çå½å¼ï¼", char.Name), "#ff00ff")
							logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])
							break // åªè§¦åä¸æ¬?						}
					}
				}
			}

			// å¤çåå»ææï¼åå»é£æ´ãå¤ä»è¢«å¨ç­ï¼?			m.handleCounterAttacks(char, enemy, enemyDamage, session, &logs)

			// å¤çè¢«å¨æè½çåå°ææï¼ç¾çåå°è¢«å¨ç­ï¼?			m.handlePassiveReflectEffects(char, enemy, enemyDamage, session, &logs)

			// å¤çä¸»å¨æè½çåå°ææï¼ç¾çåå°æè½ç­ï¼?			m.handleActiveReflectEffects(char, enemy, enemyDamage, session, &logs)

			// æå£«åå°ä¼¤å®³æ¶è·å¾ææ°
			resourceGain := 0
			if char.ResourceType == "rage" && enemyDamage > 0 {
				// åå°ä¼¤å®³è·å¾ææ°: ä¼¤å®³/æå¤§HP Ã 50ï¼è³å°?ç?				baseRageGain := int(float64(enemyDamage) / float64(char.MaxHP) * 50)
				if baseRageGain < 1 {
					baseRageGain = 1
				}

				// åºç¨è¢«å¨æè½çææ°è·å¾å æï¼æ¤æææ¡ç­ï¼?				rageGain := m.applyRageGenerationModifiers(char.ID, baseRageGain)

				char.Resource += rageGain
				resourceGain = rageGain
				if char.Resource > char.MaxResource {
					char.Resource = char.MaxResource
				}
			}

			// æå»ºæææ¥å¿æ¶æ¯ï¼åå«èµæºååï¼å¸¦é¢è²ï¼
			resourceChangeText := m.formatResourceChange(char.ResourceType, 0, resourceGain)
			
			// æ ¼å¼åä¼¤å®³å¬å¼?			enemyFormulaText := ""
			if enemyDamageDetails != nil {
				enemyFormulaText = m.formatDamageFormula(enemyDamageDetails)
			}

			m.addLog(session, "combat", fmt.Sprintf("%s æ»å»äº?%sï¼é æ %d ç¹ä¼¤å®?s%s", enemy.Name, char.Name, enemyDamage, enemyFormulaText, resourceChangeText), "#ff4444")
			logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])

			// æ£æ¥ç©å®¶æ¯å¦æ­»äº?			if char.HP <= 0 {
				char.TotalDeaths++
				// è§è²æ­»äº¡æ¶ä¸åæ­¢ææï¼ä¿æ?isRunning = trueï¼è¿æ ·ä¼æ¯ç¶æå¯ä»¥èªå¨å¤ç?				// ç¨æ·å·²ç»å¼å¯äºèªå¨ææï¼æ­»äº¡åªæ¯ææ¶è¿å¥ä¼æ¯ç¶æï¼ä¼æ¯ç»æååºè¯¥èªå¨æ¢å¤ææ?				session.CurrentEnemies = nil
				session.CurrentEnemy = nil
				session.CurrentTurnIndex = -1

				// è§è²æ­»äº¡æ¶ï¼æå£«çææ°å½?
				if char.ResourceType == "rage" {
					char.Resource = 0
				}

				// è®¡ç®å¤æ´»æ¶é´
				reviveDuration := m.calculateReviveTime(userID)
				now := time.Now()
				reviveAt := now.Add(reviveDuration)

				// è®¾ç½®è§è²HPä¸?ï¼æ­»äº¡ç¶æï¼
				char.HP = 0
				char.IsDead = true
				char.ReviveAt = &reviveAt

				// è§è²æ­»äº¡æ¶ï¼ç«å³æ¸é¤ææbuffådebuff
				if m.buffManager != nil {
					m.buffManager.ClearBuffs(char.ID)
				}

				m.addLog(session, "death", fmt.Sprintf("%s è¢«å»è´¥äº... éè¦?%d ç§å¤æ´?, char.Name, int(reviveDuration.Seconds())), "#ff0000")
				logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])

				// ææå¤±è´¥æ»ç»
				m.addBattleSummary(session, false, &logs)

				// ä¿å­æ­»äº¡æ°æ®ï¼åæ¬æ­»äº¡æ è®°ãå¤æ´»æ¶é´åææ°å½?ï¼?				m.charRepo.UpdateAfterDeath(char.ID, char.HP, char.Resource, char.TotalDeaths, &reviveAt)

				// è¿å¥ä¼æ¯ç¶æï¼ä¼æ¯æ¶é´ = å¤æ´»æ¶é´ + æ¢å¤æ¶é´ï¼æ¢å¤ä¸åHPéè¦çæ¶é´ï¼?				// æ¢å¤æ¶é´ï¼ä»0æ¢å¤å?0% HPï¼æ¯ç§æ¢å¤?%ï¼éè¦?5ç§?				recoveryTime := 25 * time.Second
				restDuration := reviveDuration + recoveryTime
				restUntil := now.Add(restDuration)
				session.IsResting = true
				session.RestUntil = &restUntil
				session.RestStartedAt = &now
				session.LastRestTick = &now
				session.RestSpeed = 1.0

				m.addLog(session, "system", fmt.Sprintf(">> è¿å¥ä¼æ¯æ¢å¤ç¶æ?(é¢è®¡ %d ç§?", int(restDuration.Seconds())+1), "#33ff33")
				logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])

				// éç½®æ¬åºææç»è®¡
				session.CurrentBattleExp = 0
				session.CurrentBattleGold = 0
				session.CurrentBattleKills = 0
				session.CurrentTurnIndex = -1

				// è§è²æ­»äº¡æ¶ï¼ç«å³è¿åï¼ç¡®ä¿åç«¯æ¸é¤æäººæ¾ç¤?				// ä¿æ isRunning = trueï¼è¿æ ·æé®ä¼æ¾ç¤º"åæ­¢ææº"ï¼ä¼æ¯ç¶æå¯ä»¥èªå¨å¤ç?				return &BattleTickResult{
					Character:    char,
					Enemy:        nil,
					Enemies:      nil, // æç¡®è¿å nilï¼è®©åç«¯æ¸é¤æäººæ¾ç¤º
					Logs:         logs,
					IsRunning:    session.IsRunning, // ä¿æè¿è¡ç¶æï¼ä¸åæ­?					IsResting:    session.IsResting,
					RestUntil:    session.RestUntil,
					SessionKills: session.SessionKills,
					SessionGold:  session.SessionGold,
					SessionExp:   session.SessionExp,
					BattleCount:  session.BattleCount,
				}, nil
			} else {
				// ç§»å¨å°ä¸ä¸ä¸ªæäººæåå°ç©å®¶åå
				session.CurrentTurnIndex++
				if session.CurrentTurnIndex >= len(aliveEnemies) {
					session.CurrentTurnIndex = -1 // åå°ç©å®¶åå
				}
			}
		} else {
			// ç´¢å¼è¶åºèå´ï¼åå°ç©å®¶åå?			session.CurrentTurnIndex = -1
		}
	}

	// æ´æ°å­æ´»æäººåè¡¨
	aliveEnemies = make([]*models.Monster, 0)
	for _, enemy := range session.CurrentEnemies {
		if enemy != nil && enemy.HP > 0 {
			aliveEnemies = append(aliveEnemies, enemy)
		}
	}

	// å¦ææææäººé½è¢«å»è´¥ï¼å¤çææç»æ
	if len(aliveEnemies) == 0 && len(session.CurrentEnemies) > 0 {
		// ç¡®ä¿æææäººçHPé½å½é?		for _, enemy := range session.CurrentEnemies {
			if enemy != nil && enemy.HP <= 0 {
				enemy.HP = 0
			}
		}

		// ææèå©æ»ç»
		m.addBattleSummary(session, true, &logs)

		// ææç»æåï¼æææå£«è§è²çææ°é½å½0
		for _, c := range characters {
			if c.ResourceType == "rage" {
				c.Resource = 0
			}
			// ä¿å­ææè§è²çæ°æ®ï¼åæ¬æå£«çææ°å½?ï¼?			m.charRepo.UpdateAfterBattle(c.ID, c.HP, c.Resource, c.Exp, c.Level,
				c.ExpToNext, c.MaxHP, c.MaxResource, c.PhysicalAttack, c.MagicAttack, c.PhysicalDefense, c.MagicDefense,
				c.Strength, c.Agility, c.Stamina, c.TotalKills)
		}

		// è®¡ç®å¹¶å¼å§ä¼æ?		restDuration := m.calculateRestTime(char)
		now := time.Now()
		restUntil := now.Add(restDuration)
		session.IsResting = true
		session.RestUntil = &restUntil
		session.RestStartedAt = &now
		session.LastRestTick = &now
		session.RestSpeed = 1.0 // é»è®¤æ¢å¤éåº¦

		m.addLog(session, "system", fmt.Sprintf(">> å¼å§ä¼æ¯æ¢å¤?(é¢è®¡ %d ç§?", int(restDuration.Seconds())+1), "#33ff33")
		logs = append(logs, session.BattleLogs[len(session.BattleLogs)-1])

		// éç½®æ¬åºææç»è®¡
		session.CurrentBattleExp = 0
		session.CurrentBattleGold = 0
		session.CurrentBattleKills = 0
		session.CurrentTurnIndex = -1

		// åè¿åä¸æ¬¡å¸¦æHP=0çæäººç¶æï¼è®©åç«¯æ´æ°æ¾ç¤?		// åå»ºæäººå¯æ¬ï¼ç¡®ä¿HPä¸?
		defeatedEnemies := make([]*models.Monster, len(session.CurrentEnemies))
		for i, enemy := range session.CurrentEnemies {
			if enemy != nil {
				defeatedEnemy := *enemy
				defeatedEnemy.HP = 0
				defeatedEnemies[i] = &defeatedEnemy
			}
		}

		// æ¸é¤æäººï¼å¨è¿åç»æä¹åï¼?		session.CurrentEnemies = nil
		session.CurrentEnemy = nil

		// è¿åå¸¦æHP=0çæäººç¶æ?		return &BattleTickResult{
			Character:    char,
			Enemy:        nil,
			Enemies:      defeatedEnemies, // è¿åHP=0çæäººå¯æ?			Logs:         logs,
			IsRunning:    session.IsRunning,
			IsResting:    session.IsResting,
			RestUntil:    session.RestUntil,
			SessionKills: session.SessionKills,
			SessionGold:  session.SessionGold,
			SessionExp:   session.SessionExp,
			BattleCount:  session.BattleCount,
		}, nil
	}

	// ä¿å­è§è²æ°æ®æ´æ°
	m.charRepo.UpdateAfterBattle(char.ID, char.HP, char.Resource, char.Exp, char.Level,
		char.ExpToNext, char.MaxHP, char.MaxResource, char.PhysicalAttack, char.MagicAttack, char.PhysicalDefense, char.MagicDefense,
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

// spawnEnemy çææäººï¼ååå¼å®¹ï¼
func (m *BattleManager) spawnEnemy(session *BattleSession, playerLevel int) error {
	return m.spawnEnemies(session, playerLevel)
}

// spawnEnemies çæå¤ä¸ªæäººï¼?-3ä¸ªéæºï¼
func (m *BattleManager) spawnEnemies(session *BattleSession, playerLevel int) error {
	if session.CurrentZone == nil {
		// å è½½é»è®¤åºå
		zone, err := m.gameRepo.GetZoneByID("elwynn")
		if err != nil {
			fmt.Printf("[ERROR] Failed to get zone: %v\n", err)
			return fmt.Errorf("failed to get zone: %v", err)
		}
		session.CurrentZone = zone
		fmt.Printf("[DEBUG] Loaded zone: %s\n", zone.Name)
	}

	// è·ååºåæªç©
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

	// éæºçæ1-3ä¸ªæäº?	enemyCount := 1 + rand.Intn(3) // 1-3ä¸?	session.CurrentEnemies = make([]*models.Monster, 0, enemyCount)

	var enemyNames []string
	for i := 0; i < enemyCount; i++ {
		// éæºéæ©ä¸ä¸ªæªç©æ¨¡æ¿
		template := monsters[rand.Intn(len(monsters))]

		enemy := &models.Monster{
			ID:              template.ID,
			ZoneID:          template.ZoneID,
			Name:            template.Name,
			Level:           template.Level,
			Type:            template.Type,
			HP:              template.HP,
			MaxHP:           template.HP,
			PhysicalAttack:  template.PhysicalAttack,
			MagicAttack:     template.MagicAttack,
			PhysicalDefense: template.PhysicalDefense,
			MagicDefense:    template.MagicDefense,
			ExpReward:       template.ExpReward,
			GoldMin:         template.GoldMin,
			GoldMax:         template.GoldMax,
		}
		session.CurrentEnemies = append(session.CurrentEnemies, enemy)
		enemyNames = append(enemyNames, fmt.Sprintf("%s (Lv.%d)", enemy.Name, enemy.Level))
	}

	// ä¿ç CurrentEnemy ç¨äºååå¼å®¹ï¼æåç¬¬ä¸ä¸ªæäººï¼
	if len(session.CurrentEnemies) > 0 {
		session.CurrentEnemy = session.CurrentEnemies[0]
	}

	session.BattleCount++
	if len(enemyNames) == 0 {
		return fmt.Errorf("failed to generate enemies")
	}
	enemyList := fmt.Sprintf("%s", enemyNames[0])
	if len(enemyNames) > 1 {
		enemyList = fmt.Sprintf("%s ç­?%d ä¸ªæäº?, enemyNames[0], len(enemyNames))
	}
	m.addLog(session, "encounter", fmt.Sprintf("âââ?ææ #%d âââ?é­é: %s", session.BattleCount, enemyList), "#ffff00")

	return nil
}

// ChangeZone åæ¢åºå
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

	m.addLog(session, "zone", fmt.Sprintf(">> ä½ æ¥å°äº [%s]", zone.Name), "#00ffff")
	m.addLog(session, "zone", zone.Description, "#888888")

	return nil
}

// GetBattleStatus è·åææç¶æ?func (m *BattleManager) GetBattleStatus(userID int) *models.BattleStatus {
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

// GetCharacterBuffs è·åè§è²çææBuff/Debuffä¿¡æ¯ï¼ç¨äºAPIè¿åï¼?func (m *BattleManager) GetCharacterBuffs(characterID int) []*models.BuffInfo {
	if m.buffManager == nil {
		return []*models.BuffInfo{}
	}
	
	buffInstances := m.buffManager.GetBuffs(characterID)
	buffs := make([]*models.BuffInfo, 0, len(buffInstances))
	
	for _, buff := range buffInstances {
		description := m.getBuffDescription(buff)
		buffInfo := &models.BuffInfo{
			EffectID:     buff.EffectID,
			Name:         buff.Name,
			Type:         buff.Type,
			IsBuff:       buff.IsBuff,
			Duration:     buff.Duration,
			Value:        buff.Value,
			StatAffected: buff.StatAffected,
			Description:  description,
		}
		buffs = append(buffs, buffInfo)
	}
	
	return buffs
}

// getBuffDescription è·åBuffçæè¿°ææ?func (m *BattleManager) getBuffDescription(buff *BuffInstance) string {
	switch buff.StatAffected {
	case "attack":
		if buff.IsBuff {
			return fmt.Sprintf("æå%.0f%%ç©çæ»å»å?, buff.Value)
		}
		return fmt.Sprintf("éä½%.0f%%ç©çæ»å»å?, -buff.Value)
	case "defense":
		if buff.IsBuff {
			return fmt.Sprintf("æå%.0f%%ç©çé²å¾¡", buff.Value)
		}
		return fmt.Sprintf("éä½%.0f%%ç©çé²å¾¡", -buff.Value)
	case "physical_damage_taken":
		return fmt.Sprintf("åå°%.0f%%åå°çç©çä¼¤å®?, -buff.Value)
	case "damage_taken":
		return fmt.Sprintf("åå°%.0f%%åå°çä¼¤å®?, -buff.Value)
	case "crit_rate":
		if buff.IsBuff {
			return fmt.Sprintf("æå%.0f%%æ´å»ç?, buff.Value)
		}
		return fmt.Sprintf("éä½%.0f%%æ´å»ç?, -buff.Value)
	case "healing_received":
		return fmt.Sprintf("éä½%.0f%%æ²»çææ", buff.Value)
	case "shield":
		return fmt.Sprintf("è·å¾ç¸å½äºæå¤§HP %.0f%%çæ¤ç?, buff.Value/float64(100))
	case "reflect":
		return fmt.Sprintf("åå°%.0f%%åå°çä¼¤å®?, buff.Value)
	case "counter_attack":
		return fmt.Sprintf("åå°æ»å»æ¶åå»ï¼é æ%.0f%%ç©çæ»å»åä¼¤å®?, buff.Value)
	case "cc_immune":
		return "åç«æ§å¶ææ"
	default:
		// å¦ææ²¡æå¹éçç±»åï¼è¿åbuffåç§°
		return buff.Name
	}
}

// GetBattleLogs è·åæææ¥å¿
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

// DamageCalculationDetails ä¼¤å®³è®¡ç®è¯¦æ
type DamageCalculationDetails struct {
	BaseAttack      int     // åºç¡æ»å»å?	ActualAttack    float64 // å®éæ»å»åï¼åºç¨å æåï¼
	BaseDefense     int     // åºç¡é²å¾¡å?	ActualDefense   float64 // å®éé²å¾¡åï¼åºç¨Debuffåï¼
	BaseDamage      float64 // åºç¡ä¼¤å®³ï¼æ»å?é²å¾¡/2ï¼?	FinalDamage     int     // æç»ä¼¤å®³ï¼åºç¨éæºæ³¢å¨åï¼
	Variance        float64 // éæºæ³¢å¨å?	IsCrit          bool    // æ¯å¦æ´å»
	CritMultiplier  float64 // æ´å»åç
	BaseCritRate    float64 // åºç¡æ´å»ç?	ActualCritRate  float64 // å®éæ´å»çï¼åºç¨å æåï¼
	RandomRoll      float64 // éæºæ°ï¼ç¨äºæ´å»å¤å®ï¼?	AttackModifiers []string // æ»å»åå æè¯´æ?	DefenseModifiers []string // é²å¾¡åä¿®æ¹è¯´æ?	CritModifiers   []string // æ´å»çå æè¯´æ?}

// calculatePhysicalDamage è®¡ç®ç©çä¼¤å®³ï¼è¿åè¯¦æï¼
func (m *BattleManager) calculatePhysicalDamageWithDetails(attack, defense int) (int, *DamageCalculationDetails) {
	details := &DamageCalculationDetails{
		BaseAttack:      attack,
		ActualAttack:    float64(attack),
		BaseDefense:     defense,
		ActualDefense:   float64(defense),
		AttackModifiers: []string{},
		DefenseModifiers: []string{},
	}
	
	// åºç¡ä¼¤å®³ = å®éæ»å»å?- ç®æ é²å¾¡åï¼ä¸åé¤ä»¥2ï¼?	baseDamage := float64(attack) - float64(defense)
	if baseDamage < 1 {
		baseDamage = 1
	}
	details.BaseDamage = baseDamage
	details.Variance = 0 // ä¸åä½¿ç¨éæºæ³¢å¨ï¼æªæ¥éè¿è£å¤çæ»å»åä¸ä¸éå®ç?	details.FinalDamage = int(baseDamage)
	
	return int(baseDamage), details
}

// calculatePhysicalDamage è®¡ç®ç©çä¼¤å®³ï¼ä¿æååå¼å®¹ï¼
func (m *BattleManager) calculatePhysicalDamage(attack, defense int) int {
	damage, _ := m.calculatePhysicalDamageWithDetails(attack, defense)
	return damage
}

// calculateMagicDamage è®¡ç®é­æ³ä¼¤å®³ï¼è¿åè¯¦æï¼
func (m *BattleManager) calculateMagicDamageWithDetails(attack, defense int) (int, *DamageCalculationDetails) {
	details := &DamageCalculationDetails{
		BaseAttack:      attack,
		ActualAttack:    float64(attack),
		BaseDefense:     defense,
		ActualDefense:   float64(defense),
		AttackModifiers: []string{},
		DefenseModifiers: []string{},
	}
	
	// åºç¡ä¼¤å®³ = å®éæ»å»å?- ç®æ é²å¾¡åï¼ä¸åé¤ä»¥2ï¼?	baseDamage := float64(attack) - float64(defense)
	if baseDamage < 1 {
		baseDamage = 1
	}
	details.BaseDamage = baseDamage
	details.Variance = 0 // ä¸åä½¿ç¨éæºæ³¢å¨ï¼æªæ¥éè¿è£å¤çæ»å»åä¸ä¸éå®ç?	details.FinalDamage = int(baseDamage)
	
	return int(baseDamage), details
}

// calculateMagicDamage è®¡ç®é­æ³ä¼¤å®³ï¼ä¿æååå¼å®¹ï¼
func (m *BattleManager) calculateMagicDamage(attack, defense int) int {
	damage, _ := m.calculateMagicDamageWithDetails(attack, defense)
	return damage
}

// calculateDamage è®¡ç®ä¼¤å®³ï¼å¼å®¹æ§ä»£ç ï¼é»è®¤ä½¿ç¨ç©çï¼
func (m *BattleManager) calculateDamage(attack, defense int) int {
	return m.calculatePhysicalDamage(attack, defense)
}

// addLog æ·»å æ¥å¿
func (m *BattleManager) addLog(session *BattleSession, logType, message, color string) {
	log := models.BattleLog{
		Message:   message,
		LogType:   logType,
		CreatedAt: time.Now(),
	}
	session.BattleLogs = append(session.BattleLogs, log)

	// ä¿ææ¥å¿æ°éå¨åçèå?	if len(session.BattleLogs) > 200 {
		session.BattleLogs = session.BattleLogs[len(session.BattleLogs)-200:]
	}
}

// addBattleSummary æ·»å æææ»ç»ååå²çº¿
func (m *BattleManager) addBattleSummary(session *BattleSession, isVictory bool, logs *[]models.BattleLog) {
	// çææææ»ç»ï¼ä½¿ç¨ä¸åé¢è²æ è®°ä¸åææ ?	var summaryMsg string
	if isVictory {
		if session.CurrentBattleKills > 0 {
			// ä½¿ç¨HTMLæ ç­¾ä¸ºä¸åé¨åæ·»å é¢è?			// ç»æï¼éè?#ffd700ï¼å»æï¼çº¢è?#ff4444ï¼ç»éªï¼èè² #3d85c6ï¼éå¸ï¼éè² #ffd700
			summaryMsg = fmt.Sprintf("âââ?æææ»ç» âââ?ç»æ: <span style=\"color: #ffd700\">â?èå©</span> | å»æ: <span style=\"color: #ff4444\">%d</span> | ç»éª: <span style=\"color: #3d85c6\">%d</span> | éå¸: <span style=\"color: #ffd700\">%d</span>",
				session.CurrentBattleKills, session.CurrentBattleExp, session.CurrentBattleGold)
		} else {
			summaryMsg = "âââ?æææ»ç» âââ?ç»æ: <span style=\"color: #ffd700\">â?èå©</span>"
		}
		m.addLog(session, "battle_summary", summaryMsg, "#ffd700")
		*logs = append(*logs, session.BattleLogs[len(session.BattleLogs)-1])
	} else {
		// å¤±è´¥æ¶çæ»ç»
		if session.CurrentBattleKills > 0 {
			// ç»æï¼çº¢è?#ff6666ï¼å»æï¼æ©è?#ffaa00ï¼ç»éªï¼èè² #3d85c6ï¼éå¸ï¼éè² #ffd700
			summaryMsg = fmt.Sprintf("âââ?æææ»ç» âââ?ç»æ: <span style=\"color: #ff6666\">â?å¤±è´¥</span> | å»æ: <span style=\"color: #ffaa00\">%d</span> | ç»éª: <span style=\"color: #3d85c6\">%d</span> | éå¸: <span style=\"color: #ffd700\">%d</span>",
				session.CurrentBattleKills, session.CurrentBattleExp, session.CurrentBattleGold)
		} else {
			summaryMsg = "âââ?æææ»ç» âââ?ç»æ: <span style=\"color: #ff6666\">â?å¤±è´¥</span>"
		}
		m.addLog(session, "battle_summary", summaryMsg, "#ff6666")
		*logs = append(*logs, session.BattleLogs[len(session.BattleLogs)-1])
	}

	// æ·»å åå²çº?	m.addLog(session, "battle_separator", "ââââââââââââââââââââââââââââââââââââââââ", "#666666")
	*logs = append(*logs, session.BattleLogs[len(session.BattleLogs)-1])
}

// getResourceName è·åèµæºçä¸­æåç§?func (m *BattleManager) getResourceName(resourceType string) string {
	switch resourceType {
	case "rage":
		return "ææ°"
	case "mana":
		return "MP"
	case "energy":
		return "è½é"
	default:
		return "èµæº"
	}
}

// getResourceColor è·åèµæºçé¢è²ï¼åèé­å½ä¸çï¼
func (m *BattleManager) getResourceColor(resourceType string) string {
	switch resourceType {
	case "rage":
		return "#ff4444" // çº¢è² - ææ°
	case "mana":
		return "#3d85c6" // èè² - æ³å
	case "energy":
		return "#ffd700" // éè²/é»è² - è½é
	default:
		return "#ffffff" // ç½è² - é»è®¤
	}
}

// formatDamageFormula æ ¼å¼åä¼¤å®³è®¡ç®å¬å¼ææ?func (m *BattleManager) formatDamageFormula(details *DamageCalculationDetails) string {
	if details == nil {
		return ""
	}
	
	var formulaParts []string
	
	// åºç¡å¬å¼ï¼æ»å»å - é²å¾¡å?	baseFormula := fmt.Sprintf("%d - %d", details.BaseAttack, details.BaseDefense)
	if details.BaseDamage > 0 {
		baseFormula = fmt.Sprintf("%s = %.0f", baseFormula, details.BaseDamage)
	}
	formulaParts = append(formulaParts, baseFormula)
	
	// å¦æææ»å»åå æ
	if len(details.AttackModifiers) > 0 {
		modifierText := strings.Join(details.AttackModifiers, ", ")
		if details.ActualAttack > float64(details.BaseAttack) {
			formulaParts = append(formulaParts, fmt.Sprintf("æ»å»å æ: %s â?%.0f", modifierText, details.ActualAttack))
		} else {
			formulaParts = append(formulaParts, fmt.Sprintf("æ»å»å æ: %s", modifierText))
		}
	}
	
	// å¦ææé²å¾¡åä¿®æ¹
	if len(details.DefenseModifiers) > 0 {
		modifierText := strings.Join(details.DefenseModifiers, ", ")
		formulaParts = append(formulaParts, fmt.Sprintf("é²å¾¡ä¿®æ¹: %s", modifierText))
	}
	
	// æ¾ç¤ºæ´å»å¤å®è¿ç¨ï¼å¦æè¿è¡äºæ´å»å¤å®ï¼?	// æ£æ¥æ¯å¦è¿è¡äºæ´å»å¤å®ï¼å¦æ?ActualCritRate >= 0 ä¸?RandomRoll >= 0ï¼è¯´æè¿è¡äºå¤å®
	if details.ActualCritRate >= 0 && details.RandomRoll >= 0 {
		critInfo := fmt.Sprintf("æ´å»ç? %.1f%%", details.BaseCritRate*100)
		if len(details.CritModifiers) > 0 {
			critInfo += fmt.Sprintf(" + %s = %.1f%%", strings.Join(details.CritModifiers, " + "), details.ActualCritRate*100)
		} else if details.ActualCritRate != details.BaseCritRate {
			critInfo += fmt.Sprintf(" = %.1f%%", details.ActualCritRate*100)
		}
		critInfo += fmt.Sprintf(" | éæº: %.3f", details.RandomRoll)
		if details.IsCrit {
			critInfo += fmt.Sprintf(" < %.3f âæ´å?, details.ActualCritRate)
			formulaParts = append(formulaParts, fmt.Sprintf("ð¥%s | ä¼¤å®³: %.0f Ã %.1f = %d", 
				critInfo, details.BaseDamage, details.CritMultiplier, details.FinalDamage))
		} else {
			critInfo += fmt.Sprintf(" â?%.3f âæªæ´å»", details.ActualCritRate)
			formulaParts = append(formulaParts, critInfo)
		}
	}
	
	if len(formulaParts) == 0 {
		return ""
	}
	
	// ä½¿ç¨è¾äº®çç°è²æ¾ç¤ºå¬å¼ï¼æé«å¯è¯»æ?	// ä½¿ç¨åæ¬å·èä¸æ¯æ¹æ¬å·ï¼é¿åè¢«åç«¯çæè½åå¹éè§åå½±å
	formulaText := strings.Join(formulaParts, " | ")
	return fmt.Sprintf(" <span style=\"color: #bbbbbb !important; opacity: 0.95;\">(%s)</span>", formulaText)
}

// formatResourceChange æ ¼å¼åèµæºååææ¬ï¼å¸¦é¢è²ï¼
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

	// å°å¤ä¸ªé¨åç¨ç©ºæ ¼è¿æ¥
	changeText := ""
	for i, part := range parts {
		if i > 0 {
			changeText += " "
		}
		changeText += part
	}

	return fmt.Sprintf(" (<span style=\"color: %s\">%s</span> %s)", color, resourceName, changeText)
}

// getRandomSkillName è·åéæºæè½åç§?func (m *BattleManager) getRandomSkillName(classID string) string {
	skills := map[string][]string{
		"warrior": {"è±åæå»", "é·éä¸å?, "é¡ºåæ?, "è´æ­»æå»"},
		"paladin": {"å£åæ?, "åå­åæå?, "æ­£ä¹ä¹é¤", "å®¡å¤"},
		"hunter":  {"å¥¥æ¯å°å»", "å¤éå°å»", "çåå°å»", "ç¨³åºå°å»"},
		"rogue":   {"éªæ¶æ»å»", "åéª¨", "èåº", "æ¯å"},
		"priest":  {"æ©å»", "æè¨æ?ç?, "ç¥å£ä¹ç«", "å¿çµéç"},
		"mage":    {"ç«çæ?, "å¯å°ç®?, "å¥¥æ¯é£å¼¹", "ççæ?},
		"warlock": {"æå½±ç®?, "èèæ?, "ç®ç¥­", "æ··ä¹±ç®?},
		"druid":   {"æç«æ?, "æ¤æ?, "æ¥å»", "æ¨ªæ«"},
		"shaman":  {"éªçµç®?, "éªçµé?, "çå²©çè£", "çç°éå»"},
	}

	if classSkills, ok := skills[classID]; ok {
		return classSkills[rand.Intn(len(classSkills))]
	}
	return "æ®éæ»å?
}

// getSkillForAttack è·åæ»å»æè½åç§°åæ¶è?func (m *BattleManager) getSkillForAttack(char *models.Character) (string, int) {
	// æå£«æè½åå¶ææ°æ¶è?	warriorSkills := []struct {
		name string
		cost int
	}{
		{"è±åæå»", 10},
		{"é·éä¸å?, 15},
		{"é¡ºåæ?, 12},
		{"è´æ­»æå»", 20},
	}

	// å¦ææ¯æå£«ï¼è¿åéæºæè½åæ¶è?	if char.ResourceType == "rage" {
		skill := warriorSkills[rand.Intn(len(warriorSkills))]
		return skill.name, skill.cost
	}

	// å¶ä»èä¸ä½¿ç¨æ®éæè½ï¼ä¸æ¶èèµæºï¼ææ¶èæ³åï¼ä½è¿éç®åå¤çï¼
	skillName := m.getRandomSkillName(char.ClassID)
	return skillName, 0
}

// calculateReviveTime è®¡ç®å¤æ´»æ¶é´ï¼æ ¹æ®æ­»äº¡äººæ°ï¼
func (m *BattleManager) calculateReviveTime(userID int) time.Duration {
	// è·åææè§è²ï¼ææè§è²é½åä¸ææï¼?	characters, err := m.charRepo.GetByUserID(userID)
	if err != nil {
		return 30 * time.Second // é»è®¤30ç§?	}

	// ç»è®¡æ­»äº¡è§è²çæ°é?	deadCount := 0
	for _, char := range characters {
		if char.IsDead {
			deadCount++
		}
	}

	// å¦ææ²¡ææ­»äº¡è§è²ï¼è¿åé»è®¤å?	if deadCount == 0 {
		deadCount = 1 // è³å°æä¸ä¸ªè§è²æ­»äº¡æä¼è°ç¨è¿ä¸ªå½æ?	}

	// æ ¹æ®æ­»äº¡äººæ°è®¡ç®å¤æ´»æ¶é´ï¼ç§ï¼?	var reviveSeconds int
	switch deadCount {
	case 1:
		reviveSeconds = 30
	case 2:
		reviveSeconds = 60
	case 3:
		reviveSeconds = 90
	case 4:
		reviveSeconds = 120
	default: // 5äººææ´å¤
		reviveSeconds = 180
	}

	return time.Duration(reviveSeconds) * time.Second
}

// calculateRestTime è®¡ç®ä¼æ¯æ¶é´ï¼åºäºHP/MPæå¤±ï¼?// æ³¨æï¼æå£«çææ°ä¸éè¦æ¢å¤ï¼ææç»æåç´æ¥å½0ï¼æ¯åºææä»0å¼å§?func (m *BattleManager) calculateRestTime(char *models.Character) time.Duration {
	hpLoss := float64(char.MaxHP - char.HP)

	// æå£«çææ°ä¸éè¦æ¢å¤ï¼åªè®¡ç®HPæå¤±
	// å¶ä»èä¸éè¦è®¡ç®MPæå¤±
	var mpLoss float64
	if char.ResourceType != "rage" {
		mpLoss = float64(char.MaxResource - char.Resource)
	} else {
		// æå£«çææ°ä¸åä¸ä¼æ¯æ¶é´è®¡ç®?		mpLoss = 0
	}

	// å¦æå·²ç»æ»¡è¡æ»¡èï¼ææ»¡è¡ï¼ï¼ä¸éè¦ä¼æ?	if hpLoss <= 0 && mpLoss <= 0 {
		return 0
	}

	// åå«è®¡ç®HPåMPçæ¢å¤æ¶é?	// æ¯ç§æ¢å¤2%ï¼æä»¥éè¦çæ¶é´ = æå¤±ç¾åæ¯?/ 0.02 = æå¤±ç¾åæ¯?* 50
	hpLossPercent := hpLoss / float64(char.MaxHP)

	hpRestSeconds := hpLossPercent * 50.0
	var mpRestSeconds float64
	if char.ResourceType != "rage" && char.MaxResource > 0 {
		mpLossPercent := mpLoss / float64(char.MaxResource)
		mpRestSeconds = mpLossPercent * 50.0
	} else {
		mpRestSeconds = 0
	}

	// åä¸¤èä¸­çæå¤§å¼ï¼å ä¸ºHPåMPæ¯åæ¶æ¢å¤ç
	restSeconds := hpRestSeconds
	if mpRestSeconds > restSeconds {
		restSeconds = mpRestSeconds
	}

	// æå°?ç§?	if restSeconds < 1.0 {
		restSeconds = 1.0
	}

	// åºç¨æ¢å¤éåº¦åçï¼æªæ¥å¯ä»¥ä»è£å¤è·åï¼?	restSpeed := 1.0 // é»è®¤æ¢å¤éåº¦
	if restSpeed > 0 {
		restSeconds = restSeconds / restSpeed
	}

	return time.Duration(restSeconds) * time.Second
}

// processRest å¤çä¼æ¯æé´çæ¢å¤?func (m *BattleManager) processRest(session *BattleSession, char *models.Character) {
	if !session.IsResting || session.RestUntil == nil || session.RestStartedAt == nil {
		return
	}

	now := time.Now()

	// æ£æ¥è§è²æ¯å¦å·²ç»å¤æ´»ï¼å¦æè§è²æ­»äº¡ä¸æå¤æ´»æ¶é´ï¼?	if char.IsDead && char.ReviveAt != nil {
		if now.After(*char.ReviveAt) || now.Equal(*char.ReviveAt) {
			// å¤æ´»æ¶é´å°äºï¼æ¢å¤è§è²å°ä¸åHP
			char.HP = char.MaxHP / 2
			char.IsDead = false
			char.ReviveAt = nil

			// æ´æ°æ°æ®åºï¼æ¸é¤æ­»äº¡æ è®°
			m.charRepo.SetDead(char.ID, false, nil)

			// æ´æ°è§è²HP
			m.charRepo.UpdateAfterBattle(char.ID, char.HP, char.Resource, char.Exp, char.Level,
				char.ExpToNext, char.MaxHP, char.MaxResource, char.PhysicalAttack, char.MagicAttack, char.PhysicalDefense, char.MagicDefense,
				char.Strength, char.Agility, char.Stamina, char.TotalKills)

			// è®°å½å¤æ´»æ¥å¿
			m.addLog(session, "revive", fmt.Sprintf("%s å·²å¤æ´»ï¼HPæ¢å¤è?%d/%d", char.Name, char.HP, char.MaxHP), "#00ff00")
		}
	}

	// æ£æ¥æ¯å¦å·²ç»æ¢å¤æ»¡è¡æ»¡èï¼å¦ææ¯åæåç»æä¼æ?	if char.HP >= char.MaxHP && char.Resource >= char.MaxResource {
		session.IsResting = false
		session.RestUntil = nil
		session.RestStartedAt = nil
		session.LastRestTick = nil
		return
	}

	if now.Before(*session.RestUntil) {
		// è®¡ç®ä»ä¸æ¬¡æ¢å¤å°ç°å¨ç»è¿çæ¶é?		var elapsed time.Duration
		if session.LastRestTick == nil {
			// ç¬¬ä¸æ¬¡æ¢å¤ï¼ä»ä¼æ¯å¼å§æ¶é´è®¡ç®?			elapsed = now.Sub(*session.RestStartedAt)
		} else {
			// ä»ä¸æ¬¡æ¢å¤æ¶é´è®¡ç®?			elapsed = now.Sub(*session.LastRestTick)
		}

		// å¦ææ¶é´é´éå¤ªé¿ï¼è¶è¿?ç§ï¼ï¼éå¶ä¸º1ç§ï¼é¿åä¸æ¬¡æ§æ¢å¤è¿å¤?		if elapsed > time.Second {
			elapsed = time.Second
		}

		// å¦ææ¶é´é´éå¤ªç­ï¼å°äº?.1ç§ï¼ï¼è·³è¿æ¢å¤ï¼é¿åé¢ç¹è®¡ç®
		if elapsed < 100*time.Millisecond {
			return
		}

		// è®¡ç®æ¢å¤éåº¦ï¼æ¯ç§æ¢å¤æå¤§å¼ç2%ï¼?		restSpeed := session.RestSpeed
		if restSpeed <= 0 {
			restSpeed = 1.0
		}

		// è®¡ç®ç»è¿çç§æ?		elapsedSeconds := elapsed.Seconds()

		// å¦æè§è²å·²ç»æ­»äº¡ä½è¿æ²¡å°å¤æ´»æ¶é´ï¼ä¸æ¢å¤HP
		if char.IsDead && char.ReviveAt != nil && now.Before(*char.ReviveAt) {
			// åªæ¢å¤èµæºï¼å¦æéç¨ï¼ï¼ä¸æ¢å¤HP
		} else {
			// æ ¹æ®å®éç»è¿çæ¶é´è®¡ç®æ¢å¤é
			hpRegenPercent := 0.02 * restSpeed * elapsedSeconds // æ¯ç§2%

			hpRegen := int(float64(char.MaxHP) * hpRegenPercent)

			// ç¡®ä¿è³å°æ¢å¤1ç¹ï¼å¦æè¿æ²¡æ»¡ï¼
			if hpRegen < 1 && char.HP < char.MaxHP {
				hpRegen = 1
			}

			char.HP += hpRegen
			if char.HP > char.MaxHP {
				char.HP = char.MaxHP
			}
		}

		// æå£«çææ°ä¸å¨ä¼æ¯æ¶æ¢å¤ï¼åªå¨ææä¸­éè¿æ»å»/åå»è·å¾
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

		// æ´æ°ä¸æ¬¡æ¢å¤æ¶é´
		session.LastRestTick = &now

		// æ¢å¤ååæ¬¡æ£æ¥æ¯å¦å·²æ»¡ï¼å¦ææ»¡äºåæåç»æä¼æ?		if char.HP >= char.MaxHP && char.Resource >= char.MaxResource {
			session.IsResting = false
			session.RestUntil = nil
			session.RestStartedAt = nil
			session.LastRestTick = nil
		}
	} else {
		// ä¼æ¯æ¶é´å°äºï¼ç»æä¼æ?		// ç¡®ä¿HPå·²æ»¡
		if char.HP < char.MaxHP {
			char.HP = char.MaxHP
		}
		// æå£«çææ°ä¸å¨ä¼æ¯æ¶æ¢å¤ï¼åªå¨ææä¸­éè¿æ»å»/åå»è·å¾
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

// BattleTickResult ææååç»æ
type BattleTickResult struct {
	Character    *models.Character  `json:"character"`
	Enemy        *models.Monster    `json:"enemy,omitempty"`
	Enemies      []*models.Monster  `json:"enemies,omitempty"` // å¤ä¸ªæäººæ¯æ
	Logs         []models.BattleLog `json:"logs"`
	IsRunning    bool               `json:"isRunning"`
	IsResting    bool               `json:"isResting"`           // æ¯å¦å¨ä¼æ?	RestUntil    *time.Time         `json:"restUntil,omitempty"` // ä¼æ¯ç»ææ¶é´
	SessionKills int                `json:"sessionKills"`
	SessionGold  int                `json:"sessionGold"`
	SessionExp   int                `json:"sessionExp"`
	BattleCount  int                `json:"battleCount"`
}

// applySkillBuffs åºç¨æè½çBuff/Debuffææ
func (m *BattleManager) applySkillBuffs(skillState *CharacterSkillState, character *models.Character, target *models.Monster, skillEffects map[string]interface{}) {
	skill := skillState.Skill
	effect := skillState.Effect

	switch skill.ID {
	case "warrior_shield_block":
		// ç¾çæ ¼æ¡ï¼åå°åå°çç©çä¼¤å®³
		if damageReduction, ok := effect["damageReduction"].(float64); ok {
			duration := 2
			if d, ok := effect["duration"].(int); ok {
				duration = d
			}
			m.buffManager.ApplyBuff(character.ID, "shield_block", "ç¾çæ ¼æ¡", "buff", true, duration, -damageReduction, "physical_damage_taken", "")
		}
	case "warrior_battle_shout":
		// æææå¼ï¼æåæ»å»å
		if attackBonus, ok := effect["attackBonus"].(float64); ok {
			duration := 5
			if d, ok := effect["duration"].(int); ok {
				duration = d
			}
			m.buffManager.ApplyBuff(character.ID, "battle_shout", "æææå¼", "buff", true, duration, attackBonus, "attack", "")
		}
	case "warrior_demoralizing_shout":
		// æ«å¿æå¼ï¼éä½æææäººæ»å»åï¼å¨applySkillDebuffsä¸­å¤çï¼
	case "warrior_whirlwind":
		// æé£æ©ï¼éä½æææäººé²å¾¡ï¼å¨applySkillDebuffsä¸­å¤çï¼
	case "warrior_mortal_strike":
		// è´æ­»æå»ï¼éä½ç®æ æ²»çææ?		if healingReduction, ok := effect["healingReduction"].(float64); ok {
			duration := 3
			if d, ok := effect["debuffDuration"].(float64); ok {
				duration = int(d)
			}
			// åºç¨å°ç®æ æäº?			if target != nil {
				m.buffManager.ApplyEnemyDebuff(target.ID, "mortal_strike", "è´æ­»æå»", "debuff", duration, healingReduction, "healing_received", "")
			}
		}
	case "warrior_last_stand":
		// ç ´éæ²èï¼ç«å³æ¢å¤æå¤§HPçç¾åæ¯
		if healPercent, ok := effect["healPercent"].(float64); ok {
			// ç«å³æ¢å¤
			healAmount := int(float64(character.MaxHP) * healPercent / 100.0)
			character.HP += healAmount
			if character.HP > character.MaxHP {
				character.HP = character.MaxHP
			}
			// éè¿skillEffectsä¼ éï¼å¨æææ¥å¿ä¸­æ¾ç¤º
			skillEffects["healMaxHpPercent"] = healPercent
		}
	case "warrior_unbreakable_barrier":
		// ä¸ç­å£åï¼è·å¾æ¤ç?		if shieldPercent, ok := effect["shieldPercent"].(float64); ok {
			duration := 4
			if d, ok := effect["duration"].(float64); ok {
				duration = int(d)
			}
			shieldAmount := int(float64(character.MaxHP) * shieldPercent / 100.0)
			// ä½¿ç¨Buffå­å¨æ¤ç¾å¼ï¼statAffectedä¸?shield"ï¼valueä¸ºæ¤ç¾å?			m.buffManager.ApplyBuff(character.ID, "unbreakable_barrier", "ä¸ç­å£å", "buff", true, duration, float64(shieldAmount), "shield", "")
		}
	case "warrior_shield_reflection":
		// ç¾çåå°ï¼åå°åå°çä¼¤å®³
		if reflectPercent, ok := effect["reflectPercent"].(float64); ok {
			duration := 2
			if d, ok := effect["duration"].(float64); ok {
				duration = int(d)
			}
			// ä½¿ç¨Buffå­å¨åå°æ¯ä¾ï¼statAffectedä¸?reflect"ï¼valueä¸ºåå°ç¾åæ¯
			m.buffManager.ApplyBuff(character.ID, "shield_reflection", "ç¾çåå°", "buff", true, duration, reflectPercent, "reflect", "")
		}
	case "warrior_shield_wall":
		// ç¾å¢ï¼å¤§å¹åå°åå°çä¼¤å®³
		if damageReduction, ok := effect["damageReduction"].(float64); ok {
			duration := 2
			if d, ok := effect["duration"].(float64); ok {
				duration = int(d)
			}
			m.buffManager.ApplyBuff(character.ID, "shield_wall", "ç¾å¢", "buff", true, duration, -damageReduction, "damage_taken", "")
		}
	case "warrior_recklessness":
		// é²è½ï¼æåæ´å»çï¼ä½åå°ä¼¤å®³å¢å 
		if critBonus, ok := effect["critBonus"].(float64); ok {
			duration := 3
			if d, ok := effect["duration"].(float64); ok {
				duration = int(d)
			}
			m.buffManager.ApplyBuff(character.ID, "recklessness_crit", "é²è½", "buff", true, duration, critBonus, "crit_rate", "")
		}
		if damageIncrease, ok := effect["damageTakenIncrease"].(float64); ok {
			duration := 3
			if d, ok := effect["duration"].(float64); ok {
				duration = int(d)
			}
			m.buffManager.ApplyBuff(character.ID, "recklessness_damage", "é²è½", "debuff", false, duration, damageIncrease, "damage_taken", "")
		}
	case "warrior_retaliation":
		// åå»é£æ´ï¼åå°æ»å»æ¶åå»
		if counterDamage, ok := effect["counterDamagePercent"].(float64); ok {
			duration := 3
			if d, ok := effect["duration"].(float64); ok {
				duration = int(d)
			}
			m.buffManager.ApplyBuff(character.ID, "retaliation", "åå»é£æ´", "buff", true, duration, counterDamage, "counter_attack", "")
		}
	case "warrior_berserker_rage":
		// çæ´ä¹æï¼æåæ»å»ååææ°è·å
		if attackBonus, ok := effect["attackBonus"].(float64); ok {
			duration := 4
			if d, ok := effect["duration"].(float64); ok {
				duration = int(d)
			}
			m.buffManager.ApplyBuff(character.ID, "berserker_rage", "çæ´ä¹æ?, "buff", true, duration, attackBonus, "attack", "")
		}
	case "warrior_avatar":
		// å¤©ç¥ä¸å¡ï¼å¤§å¹æåæ»å»åï¼åç«æ§å?		if attackBonus, ok := effect["attackBonus"].(float64); ok {
			duration := 3
			if d, ok := effect["duration"].(float64); ok {
				duration = int(d)
			}
			m.buffManager.ApplyBuff(character.ID, "avatar", "å¤©ç¥ä¸å¡", "buff", true, duration, attackBonus, "attack", "")
		}
		if immuneCC, ok := effect["immuneCC"].(bool); ok && immuneCC {
			duration := 3
			if d, ok := effect["duration"].(float64); ok {
				duration = int(d)
			}
			m.buffManager.ApplyBuff(character.ID, "avatar_cc_immune", "å¤©ç¥ä¸å¡", "buff", true, duration, 1.0, "cc_immune", "")
		}
	}
}

// handleCounterAttacks å¤çåå»ææ
func (m *BattleManager) handleCounterAttacks(character *models.Character, attacker *models.Monster, damageTaken int, session *BattleSession, logs *[]models.BattleLog) {
	// å¤çBuffçåå»ææï¼åå»é£æ´ï¼?	buffs := m.buffManager.GetBuffs(character.ID)
	for _, buff := range buffs {
		if buff.StatAffected == "counter_attack" && buff.IsBuff {
			// åå»é£æ´ï¼å¯¹æ»å»èé æåå»ä¼¤å®³
			counterDamage := int(float64(character.PhysicalAttack) * buff.Value / 100.0)
			attacker.HP -= counterDamage
			if attacker.HP < 0 {
				attacker.HP = 0
			}
			m.addLog(session, "combat", fmt.Sprintf("%s çåå»é£æ´å¯¹ %s é æ %d ç¹åå»ä¼¤å®³ï¼", character.Name, attacker.Name, counterDamage), "#ff8800")
			*logs = append(*logs, session.BattleLogs[len(session.BattleLogs)-1])
		}
	}

	// å¤çè¢«å¨æè½çåå»ææï¼å¤ä»ï¼
	if m.passiveSkillManager != nil {
		passives := m.passiveSkillManager.GetPassiveSkills(character.ID)
		for _, passive := range passives {
			if passive.Passive.EffectType == "counter_attack" {
				// å¤ä»ï¼åå°æ»å»æ¶æ¦çåå»
				// effectValueæ¯è§¦åæ¦çï¼ç¾åæ¯ï¼ï¼éè¦æ ¹æ®ç­çº§è®¡ç®å®éæ¦çåä¼¤å®³
				triggerChance := passive.EffectValue / 100.0
				if rand.Float64() < triggerChance {
					// è®¡ç®åå»ä¼¤å®³ï¼æ ¹æ®ç­çº§ï¼1çº?00%ï¼?çº?80%ï¼?					counterDamagePercent := 100.0 + float64(passive.Level-1)*20.0
					// è®¡ç®å®éæ»å»åï¼åºç¨è¢«å¨æè½åBuffå æï¼?					actualAttack := float64(character.PhysicalAttack)
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
					counterDamage = counterDamage - attacker.PhysicalDefense/2
					if counterDamage < 1 {
						counterDamage = 1
					}
					attacker.HP -= counterDamage
					if attacker.HP < 0 {
						attacker.HP = 0
					}
					m.addLog(session, "combat", fmt.Sprintf("%s çå¤ä»å¯¹ %s é æ %d ç¹åå»ä¼¤å®³ï¼", character.Name, attacker.Name, counterDamage), "#ff8800")
					*logs = append(*logs, session.BattleLogs[len(session.BattleLogs)-1])
				}
			}
		}
	}
}

// handlePassiveOnHitEffects å¤çè¢«å¨æè½çæ»å»æ¶ææ?func (m *BattleManager) handlePassiveOnHitEffects(character *models.Character, damageDealt int, usedSkill bool, session *BattleSession, logs *[]models.BattleLog) {
	if m.passiveSkillManager == nil {
		return
	}

	passives := m.passiveSkillManager.GetPassiveSkills(character.ID)
	for _, passive := range passives {
		switch passive.Passive.EffectType {
		case "on_hit_heal":
			// è¡ä¹çç­ï¼æ¯æ¬¡æ»å»æ¢å¤çå½å?			healPercent := passive.EffectValue // ç¾åæ¯å¼ï¼å¦?.0è¡¨ç¤º1%ï¼?			healAmount := int(float64(character.MaxHP) * healPercent / 100.0)
			if healAmount > 0 {
				character.HP += healAmount
				if character.HP > character.MaxHP {
					character.HP = character.MaxHP
				}
				m.addLog(session, "heal", fmt.Sprintf("%s çè¡ä¹çç­æ¢å¤äº %d ç¹çå½å?, character.Name, healAmount), "#00ff00")
				*logs = append(*logs, session.BattleLogs[len(session.BattleLogs)-1])
			}
		}
	}
}

// handlePassiveDamageReduction å¤çè¢«å¨æè½çåä¼¤ææ
func (m *BattleManager) handlePassiveDamageReduction(character *models.Character, damage int) int {
	if m.passiveSkillManager == nil {
		return damage
	}

	passives := m.passiveSkillManager.GetPassiveSkills(character.ID)
	for _, passive := range passives {
		if passive.Passive.EffectType == "survival" && passive.Passive.ID == "warrior_passive_unbreakable_will" {
			// ä¸ç­æå¿ï¼HPä½äºéå¼æ¶åä¼¤
			hpPercent := float64(character.HP) / float64(character.MaxHP)
			// æ ¹æ®ç­çº§è®¡ç®è§¦åéå¼ï¼1çº?0%ï¼?çº?0%ï¼?			threshold := 0.30 - float64(passive.Level-1)*0.05
			if hpPercent < threshold {
				// æ ¹æ®ç­çº§è®¡ç®åä¼¤æ¯ä¾ï¼?çº?5%ï¼?çº?5%ï¼?				reductionPercent := 25.0 + float64(passive.Level-1)*10.0
				damage = int(float64(damage) * (1.0 - reductionPercent/100.0))
				if damage < 1 {
					damage = 1
				}
			}
		}
	}

	return damage
}

// handleActiveReflectEffects å¤çä¸»å¨æè½çåå°ææ
func (m *BattleManager) handleActiveReflectEffects(character *models.Character, attacker *models.Monster, damageTaken int, session *BattleSession, logs *[]models.BattleLog) {
	if m.buffManager == nil {
		return
	}

	buffs := m.buffManager.GetBuffs(character.ID)
	for _, buff := range buffs {
		if buff.StatAffected == "reflect" && buff.IsBuff && buff.EffectID == "shield_reflection" {
			// ç¾çåå°ï¼ä¸»å¨æè½ï¼ï¼åå°åå°çä¼¤å®³
			reflectPercent := buff.Value // ç¾åæ¯å¼ï¼å¦?0.0è¡¨ç¤º50%ï¼?			reflectDamage := int(float64(damageTaken) * reflectPercent / 100.0)
			if reflectDamage > 0 {
				attacker.HP -= reflectDamage
				if attacker.HP < 0 {
					attacker.HP = 0
				}
				m.addLog(session, "combat", fmt.Sprintf("%s çç¾çåå°å¯¹ %s é æ %d ç¹åå°ä¼¤å®³ï¼", character.Name, attacker.Name, reflectDamage), "#ff8800")
				*logs = append(*logs, session.BattleLogs[len(session.BattleLogs)-1])
			}
		}
	}
}

// updateShieldValue æ´æ°æ¤ç¾å?func (m *BattleManager) updateShieldValue(characterID int, newShieldValue float64) {
	if m.buffManager == nil {
		return
	}

	buffs := m.buffManager.GetBuffs(characterID)
	if buff, exists := buffs["unbreakable_barrier"]; exists {
		buff.Value = newShieldValue
	}
}

// applySkillDebuffs åºç¨æè½çDebuffææå°æäº?func (m *BattleManager) applySkillDebuffs(skillState *CharacterSkillState, character *models.Character, target *models.Monster, allEnemies []*models.Monster, skillEffects map[string]interface{}) {
	skill := skillState.Skill
	effect := skillState.Effect

	switch skill.ID {
	case "warrior_demoralizing_shout":
		// æ«å¿æå¼ï¼éä½æææäººæ»å»å
		if attackReduction, ok := effect["attackReduction"].(float64); ok {
			duration := 3
			if d, ok := effect["duration"].(float64); ok {
				duration = int(d)
			}
			// åºç¨å°ææå­æ´»çæäºº
			for _, enemy := range allEnemies {
				if enemy.HP > 0 {
					m.buffManager.ApplyEnemyDebuff(enemy.ID, "demoralizing_shout", "æ«å¿æå¼", "debuff", duration, attackReduction, "attack", "")
				}
			}
		}
	case "warrior_whirlwind":
		// æé£æ©ï¼éä½æææäººé²å¾?		if defenseReduction, ok := effect["defenseReduction"].(float64); ok {
			duration := 2
			if d, ok := effect["debuffDuration"].(float64); ok {
				duration = int(d)
			}
			// åºç¨å°ææå­æ´»çæäºº
			for _, enemy := range allEnemies {
				if enemy.HP > 0 {
					m.buffManager.ApplyEnemyDebuff(enemy.ID, "whirlwind", "æé£æ?, "debuff", duration, defenseReduction, "defense", "")
				}
			}
		}
	case "warrior_mortal_strike":
		// è´æ­»æå»ï¼éä½ç®æ æ²»çææ?		if healingReduction, ok := effect["healingReduction"].(float64); ok {
			duration := 3
			if d, ok := effect["debuffDuration"].(float64); ok {
				duration = int(d)
			}
			// åºç¨å°ç®æ æäº?			if target != nil && target.HP > 0 {
				m.buffManager.ApplyEnemyDebuff(target.ID, "mortal_strike", "è´æ­»æå»", "debuff", duration, healingReduction, "healing_received", "")
			}
		}
	}
}

// handlePassiveReflectEffects å¤çè¢«å¨æè½çåå°ææ
func (m *BattleManager) handlePassiveReflectEffects(character *models.Character, attacker *models.Monster, damageTaken int, session *BattleSession, logs *[]models.BattleLog) {
	if m.passiveSkillManager == nil {
		return
	}

	passives := m.passiveSkillManager.GetPassiveSkills(character.ID)
	for _, passive := range passives {
		if passive.Passive.EffectType == "reflect" && passive.Passive.ID == "warrior_passive_shield_reflection" {
			// ç¾çåå°ï¼è¢«å¨ï¼ï¼åå°ç©çæ»å»æ¶åå°ä¼¤å®³
			reflectPercent := passive.EffectValue // ç¾åæ¯å¼ï¼å¦?0.0è¡¨ç¤º10%ï¼?			reflectDamage := int(float64(damageTaken) * reflectPercent / 100.0)
			if reflectDamage > 0 {
				attacker.HP -= reflectDamage
				if attacker.HP < 0 {
					attacker.HP = 0
				}
				m.addLog(session, "combat", fmt.Sprintf("%s çç¾çåå°å¯¹ %s é æ %d ç¹åå°ä¼¤å®³ï¼", character.Name, attacker.Name, reflectDamage), "#ff8800")
				*logs = append(*logs, session.BattleLogs[len(session.BattleLogs)-1])
			}
		}
	}
}
