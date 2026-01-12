package runner



import (
	"fmt"
	"strings"

	"text-wow/internal/repository"
)

// Context 相关函数



func (tr *TestRunner) updateAssertionContext() {

	// 同步角色信息

	if char, ok := tr.context.Characters["character"]; ok && char != nil {

		// 确保char不是nil指针

		tr.safeSetContext("character.hp", char.HP)

		tr.safeSetContext("character.max_hp", char.MaxHP)

		tr.safeSetContext("character.level", char.Level)

		tr.safeSetContext("character.resource", char.Resource)

		tr.safeSetContext("character.max_resource", char.MaxResource)

		tr.safeSetContext("character.physical_attack", char.PhysicalAttack)

		tr.safeSetContext("character.magic_attack", char.MagicAttack)

		tr.safeSetContext("character.physical_defense", char.PhysicalDefense)

		tr.safeSetContext("character.magic_defense", char.MagicDefense)

		tr.safeSetContext("character.phys_crit_rate", char.PhysCritRate)

		tr.safeSetContext("character.phys_crit_damage", char.PhysCritDamage)

		tr.safeSetContext("character.spell_crit_rate", char.SpellCritRate)

		tr.safeSetContext("character.spell_crit_damage", char.SpellCritDamage)

		tr.safeSetContext("character.dodge_rate", char.DodgeRate)

		tr.safeSetContext("character.id", char.ID)

		tr.safeSetContext("character.strength", char.Strength)

		tr.safeSetContext("character.agility", char.Agility)

		tr.safeSetContext("character.intellect", char.Intellect)

		tr.safeSetContext("character.stamina", char.Stamina)

		tr.safeSetContext("character.spirit", char.Spirit)

		

		// 获取用户金币（Gold在User模型中，不在Character模型中）

		userRepo := repository.NewUserRepository()

		user, err := userRepo.GetByID(char.UserID)

		if err == nil && user != nil {

			tr.safeSetContext("character.gold", user.Gold)

			tr.context.Variables["character.gold"] = user.Gold

			tr.safeSetContext("gold", user.Gold)

			tr.context.Variables["gold"] = user.Gold

		} else {

			// 如果获取失败，从Variables中获取（可能在setup中设置了）
			if goldVal, exists := tr.context.Variables["character.gold"]; exists {

				tr.safeSetContext("character.gold", goldVal)

				tr.safeSetContext("gold", goldVal)

				tr.context.Variables["gold"] = goldVal

			} else {

				tr.safeSetContext("character.gold", 0)

				tr.context.Variables["character.gold"] = 0

				tr.safeSetContext("gold", 0)

				tr.context.Variables["gold"] = 0

			}

		}



		// 同时设置简化路径（不带character.前缀），以支持测试用例中的直接访问
		tr.safeSetContext("hp", char.HP)

		tr.safeSetContext("max_hp", char.MaxHP)

		tr.safeSetContext("level", char.Level)

		tr.safeSetContext("resource", char.Resource)

		tr.safeSetContext("max_resource", char.MaxResource)

		tr.safeSetContext("physical_attack", char.PhysicalAttack)

		tr.safeSetContext("magic_attack", char.MagicAttack)

		tr.safeSetContext("physical_defense", char.PhysicalDefense)

		tr.safeSetContext("magic_defense", char.MagicDefense)

		tr.safeSetContext("phys_crit_rate", char.PhysCritRate)

		tr.safeSetContext("phys_crit_damage", char.PhysCritDamage)

		tr.safeSetContext("spell_crit_rate", char.SpellCritRate)

		tr.safeSetContext("spell_crit_damage", char.SpellCritDamage)

		tr.safeSetContext("dodge_rate", char.DodgeRate)

		tr.safeSetContext("strength", char.Strength)

		tr.safeSetContext("agility", char.Agility)

		tr.safeSetContext("intellect", char.Intellect)

		tr.safeSetContext("stamina", char.Stamina)

		tr.safeSetContext("spirit", char.Spirit)



		// 计算并同步速度（speed = agility）
		speed := tr.calculator.CalculateSpeed(char)

		tr.safeSetContext("character.speed", speed)

		tr.safeSetContext("speed", speed)



		// 同步从Variables中存储的计算属性（如果存在，优先使用）

		// 这些值可能是通过"计算物理攻击"等步骤计算出来的

		if physicalAttack, exists := tr.context.Variables["physical_attack"]; exists {

			tr.safeSetContext("physical_attack", physicalAttack)

		}

		if magicAttack, exists := tr.context.Variables["magic_attack"]; exists {

			tr.safeSetContext("magic_attack", magicAttack)

		}

		if maxHP, exists := tr.context.Variables["max_hp"]; exists {

			tr.safeSetContext("max_hp", maxHP)

		}

		if physCritRate, exists := tr.context.Variables["phys_crit_rate"]; exists {

			tr.safeSetContext("phys_crit_rate", physCritRate)

		}

		if spellCritRate, exists := tr.context.Variables["spell_crit_rate"]; exists {

			tr.safeSetContext("spell_crit_rate", spellCritRate)

		}

		if dodgeRate, exists := tr.context.Variables["dodge_rate"]; exists {

			tr.safeSetContext("dodge_rate", dodgeRate)

		}

		if physCritDamage, exists := tr.context.Variables["phys_crit_damage"]; exists {

			tr.safeSetContext("phys_crit_damage", physCritDamage)

			tr.safeSetContext("character.phys_crit_damage", physCritDamage)

		}

		if spellCritDamage, exists := tr.context.Variables["spell_crit_damage"]; exists {

			tr.safeSetContext("spell_crit_damage", spellCritDamage)

			tr.safeSetContext("character.spell_crit_damage", spellCritDamage)

		}

		if speedVal, exists := tr.context.Variables["speed"]; exists {

			tr.safeSetContext("speed", speedVal)

			tr.safeSetContext("character.speed", speedVal)

		}

		if manaRegen, exists := tr.context.Variables["mana_regen"]; exists {

			tr.safeSetContext("mana_regen", manaRegen)

		}

		if rageGain, exists := tr.context.Variables["rage_gain"]; exists {

			tr.safeSetContext("rage_gain", rageGain)

		}

		if energyRegen, exists := tr.context.Variables["energy_regen"]; exists {

			tr.safeSetContext("energy_regen", energyRegen)

		}

		if physicalDefense, exists := tr.context.Variables["physical_defense"]; exists {

			tr.safeSetContext("physical_defense", physicalDefense)

			tr.safeSetContext("character.physical_defense", physicalDefense)

		}

		if magicDefense, exists := tr.context.Variables["magic_defense"]; exists {

			tr.safeSetContext("magic_defense", magicDefense)

			tr.safeSetContext("character.magic_defense", magicDefense)

		}



		// 同步Buff信息（从上下文获取）

		if buffModifier, exists := tr.context.Variables["character_buff_attack_modifier"]; exists {

			tr.safeSetContext("character.buff_attack_modifier", buffModifier)

		}

		if buffDuration, exists := tr.context.Variables["character_buff_duration"]; exists {

			tr.safeSetContext("character.buff_duration", buffDuration)

		}

	}



	// 同步所有角色信息（character, character_1, character_2等）

	for key, char := range tr.context.Characters {

		if char != nil {

			// 设置角色的基本属性
			tr.safeSetContext(fmt.Sprintf("%s.hp", key), char.HP)

			tr.safeSetContext(fmt.Sprintf("%s.max_hp", key), char.MaxHP)

			tr.safeSetContext(fmt.Sprintf("%s.level", key), char.Level)

			tr.safeSetContext(fmt.Sprintf("%s.resource", key), char.Resource)

			tr.safeSetContext(fmt.Sprintf("%s.max_resource", key), char.MaxResource)

			tr.safeSetContext(fmt.Sprintf("%s.physical_attack", key), char.PhysicalAttack)

			tr.safeSetContext(fmt.Sprintf("%s.magic_attack", key), char.MagicAttack)

			tr.safeSetContext(fmt.Sprintf("%s.id", key), char.ID)

			tr.safeSetContext(fmt.Sprintf("%s.name", key), char.Name)

			

			// 如果key是职业名称（如warrior, mage, priest），也设置
			// 这需要从角色名称或ClassID推断

			if strings.Contains(strings.ToLower(char.Name), "战士") || char.ClassID == "warrior" {

				tr.safeSetContext("warrior.hp", char.HP)

				tr.safeSetContext("warrior.max_hp", char.MaxHP)

				tr.safeSetContext("warrior.id", char.ID)

			}

			if strings.Contains(strings.ToLower(char.Name), "法师") || char.ClassID == "mage" {

				tr.safeSetContext("mage.hp", char.HP)

				tr.safeSetContext("mage.max_hp", char.MaxHP)

				tr.safeSetContext("mage.id", char.ID)

			}

			if strings.Contains(strings.ToLower(char.Name), "牧师") || char.ClassID == "priest" {

				tr.safeSetContext("priest.hp", char.HP)

				tr.safeSetContext("priest.max_hp", char.MaxHP)

				tr.safeSetContext("priest.id", char.ID)

			}

		}

	}



	// 同步怪物信息

	for key, monster := range tr.context.Monsters {

		if monster != nil {

			tr.safeSetContext(fmt.Sprintf("%s.hp", key), monster.HP)

			tr.safeSetContext(fmt.Sprintf("%s.max_hp", key), monster.MaxHP)

		}

	}



	// 同步所有monster_X.hp_damage值（从Variables中读取，只同步可序列化的值）

	for i := 1; i <= 10; i++ {

		damageKey := fmt.Sprintf("monster_%d.hp_damage", i)

		if hpDamage, exists := tr.context.Variables[damageKey]; exists {

			tr.safeSetContext(damageKey, hpDamage)

		}

	}



	// 同步技能伤害值（只同步可序列化的值）

	if skillDamage, exists := tr.context.Variables["skill_damage_dealt"]; exists {

		tr.safeSetContext("skill_damage_dealt", skillDamage)

	}



	// 同步治疗相关值（只同步可序列化的值）

	if overhealing, exists := tr.context.Variables["overhealing"]; exists {

		tr.safeSetContext("overhealing", overhealing)

	}

	if skillHealing, exists := tr.context.Variables["skill_healing_done"]; exists {

		tr.safeSetContext("skill_healing_done", skillHealing)

	}



	// 同步怪物技能相关值（只同步可序列化的值）

	if monsterSkillDamage, exists := tr.context.Variables["monster_skill_damage_dealt"]; exists {

		tr.safeSetContext("monster_skill_damage_dealt", monsterSkillDamage)

	}

	if monsterHealing, exists := tr.context.Variables["monster_healing_dealt"]; exists {

		tr.safeSetContext("monster_healing_dealt", monsterHealing)

	}

	if monsterResource, exists := tr.context.Variables["monster.resource"]; exists {

		tr.safeSetContext("monster.resource", monsterResource)

	}

	if monsterSkillResourceCost, exists := tr.context.Variables["monster_skill_resource_cost"]; exists {

		tr.safeSetContext("monster_skill_resource_cost", monsterSkillResourceCost)

	}

	if monsterSkillIsCrit, exists := tr.context.Variables["monster_skill_is_crit"]; exists {

		tr.safeSetContext("monster_skill_is_crit", monsterSkillIsCrit)

	}

	if monsterSkillCritDamage, exists := tr.context.Variables["monster_skill_crit_damage"]; exists {

		tr.safeSetContext("monster_skill_crit_damage", monsterSkillCritDamage)

	}

	if monsterDebuffDuration, exists := tr.context.Variables["character_debuff_duration"]; exists {

		tr.safeSetContext("character_debuff_duration", monsterDebuffDuration)

	}



	// 同步装备信息（从 Equipments map Variables 中的 equipment_id 获取）
	if eqID, ok := tr.context.Variables["equipment_id"].(int); ok {

		if eq, exists := tr.context.Equipments[fmt.Sprintf("%d", eqID)]; exists {

			tr.syncEquipmentToContext("equipment", eq)

		}

	}

	if weaponID, ok := tr.context.Variables["weapon_id"].(int); ok {

		if eq, exists := tr.context.Equipments[fmt.Sprintf("%d", weaponID)]; exists {

			tr.syncEquipmentToContext("weapon", eq)

		}

	}

	if oldWeaponID, ok := tr.context.Variables["old_weapon_id"].(int); ok {

		if eq, exists := tr.context.Equipments[fmt.Sprintf("%d", oldWeaponID)]; exists {

			tr.syncEquipmentToContext("old_weapon", eq)

		}

	}

	if oldEquipmentID, ok := tr.context.Variables["old_equipment_id"].(int); ok {

		if eq, exists := tr.context.Equipments[fmt.Sprintf("%d", oldEquipmentID)]; exists {

			tr.syncEquipmentToContext("old_equipment", eq)

		}

	}

	if newWeaponID, ok := tr.context.Variables["new_weapon_id"].(int); ok {

		if eq, exists := tr.context.Equipments[fmt.Sprintf("%d", newWeaponID)]; exists {

			tr.syncEquipmentToContext("new_weapon", eq)

		}

	}

	if newEquipmentID, ok := tr.context.Variables["new_equipment_id"].(int); ok {

		if eq, exists := tr.context.Equipments[fmt.Sprintf("%d", newEquipmentID)]; exists {

			tr.syncEquipmentToContext("new_equipment", eq)

		}

	}



	// 同步装备槽位计数（用于测试槽位冲突）

	if char, ok := tr.context.Characters["character"]; ok && char != nil {

		equipmentRepo := repository.NewEquipmentRepository()

		mainHandCount := 0

		equippedEquipments, _ := equipmentRepo.GetByCharacterID(char.ID)

		for _, eq := range equippedEquipments {

			if eq.Slot == "main_hand" {

				mainHandCount++

			}

		}

		tr.safeSetContext("equipped_main_hand_count", mainHandCount)

	}



	// 同步战斗状态相关变量（只同步可序列化的值）

	if battleState, exists := tr.context.Variables["battle_state"]; exists {

		tr.safeSetContext("battle_state", battleState)

	}

	if isResting, exists := tr.context.Variables["is_resting"]; exists {

		tr.safeSetContext("is_resting", isResting)

	}

	if restUntil, exists := tr.context.Variables["rest_until"]; exists {

		tr.safeSetContext("rest_until", restUntil)

	}

	if restSpeed, exists := tr.context.Variables["rest_speed"]; exists {

		tr.safeSetContext("rest_speed", restSpeed)

	}

	if turnOrder, exists := tr.context.Variables["turn_order"]; exists {

		tr.safeSetContext("turn_order", turnOrder)

	} else {

		debugPrint("[DEBUG] updateAssertionContext: turn_order is not serializable, skipping\n")

	}

	if turnOrderLength, exists := tr.context.Variables["turn_order_length"]; exists {

		tr.safeSetContext("turn_order_length", turnOrderLength)

	}

	if enemyCount, exists := tr.context.Variables["enemy_count"]; exists {

		tr.safeSetContext("enemy_count", enemyCount)

	}

	if enemyAliveCount, exists := tr.context.Variables["enemy_alive_count"]; exists {

		tr.safeSetContext("enemy_alive_count", enemyAliveCount)

			// 同时设置别名 enemies_alive_count（复数形式）

			tr.safeSetContext("enemies_alive_count", enemyAliveCount)

		}

	if currentRound, exists := tr.context.Variables["current_round"]; exists {

		tr.safeSetContext("current_round", currentRound)

	}

	// 同步战斗日志

	if battleLogs, exists := tr.context.Variables["battle_logs"]; exists {

		tr.safeSetContext("battle_logs", battleLogs)

	}

	// 同步战斗结果

	if battleResultVictory, exists := tr.context.Variables["battle_result.is_victory"]; exists {

		tr.safeSetContext("battle_result.is_victory", battleResultVictory)

	}

	if battleResultDuration, exists := tr.context.Variables["battle_result.duration_seconds"]; exists {

		tr.safeSetContext("battle_result.duration_seconds", battleResultDuration)

	}

	// 同步角色状态
	if isDead, exists := tr.context.Variables["character.is_dead"]; exists {

		tr.safeSetContext("character.is_dead", isDead)

	}

	if expGained, exists := tr.context.Variables["character.exp_gained"]; exists {

		tr.safeSetContext("character.exp_gained", expGained)

	}

	if goldGained, exists := tr.context.Variables["character.gold_gained"]; exists {

		tr.safeSetContext("character.gold_gained", goldGained)

	}

	if battleRounds, exists := tr.context.Variables["battle_rounds"]; exists {

		tr.safeSetContext("battle_rounds", battleRounds)

	}

	// 同步队伍信息

	tr.syncTeamToContext()



	// 同步所有变量（包括上面已经同步的，确保覆盖）
	// 只复制可序列化的基本类型，避免序列化错误

	for key, value := range tr.context.Variables {

		tr.safeSetContext(key, value)

	}

}



