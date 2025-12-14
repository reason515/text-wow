package game

import (
	"fmt"

	"text-wow/internal/models"
)

// applyRageGenerationModifiers 应用怒气获得加成（被动技能：愤怒掌握等）
func (m *BattleManager) applyRageGenerationModifiers(characterID int, baseRageGain int) int {
	if m.passiveSkillManager == nil {
		return baseRageGain
	}

	passives := m.passiveSkillManager.GetPassiveSkills(characterID)
	for _, passive := range passives {
		if passive.Passive.EffectType == "rage_generation" && passive.Passive.ID == "warrior_passive_anger_management" {
			// 愤怒掌握：提升怒气获得量（百分比加成）
			rageBonusPercent := passive.EffectValue // 百分比值（如10.0表示10%）
			bonusRage := int(float64(baseRageGain) * rageBonusPercent / 100.0)
			baseRageGain += bonusRage
		}
	}

	return baseRageGain
}

// handleWarMachineRageGain 处理战争机器的击杀回怒效果
func (m *BattleManager) handleWarMachineRageGain(character *models.Character, session *BattleSession, logs *[]models.BattleLog) {
	if m.passiveSkillManager == nil || character.ResourceType != "rage" {
		return
	}

	passives := m.passiveSkillManager.GetPassiveSkills(character.ID)
	for _, passive := range passives {
		if passive.Passive.EffectType == "rage_generation" && passive.Passive.ID == "warrior_passive_war_machine" {
			// 战争机器：击杀敌人时获得额外怒气
			rageGain := int(passive.EffectValue) // 效果值就是怒气数量（1级30，5级70）
			character.Resource += rageGain
			if character.Resource > character.MaxResource {
				character.Resource = character.MaxResource
			}
			m.addLog(session, "rage", fmt.Sprintf("%s 的战争机器效果触发，获得了 %d 点额外怒气！", character.Name, rageGain), "#ffaa00")
			*logs = append(*logs, session.BattleLogs[len(session.BattleLogs)-1])
		}
	}
}


























