package service

import (
	"errors"
	"fmt"

	"text-wow/internal/models"
	"text-wow/internal/repository"
)

// SkillService 技能服务
type SkillService struct {
	skillRepo     *repository.SkillRepository
	characterRepo *repository.CharacterRepository
}

// NewSkillService 创建技能服务
func NewSkillService(skillRepo *repository.SkillRepository, characterRepo *repository.CharacterRepository) *SkillService {
	return &SkillService{
		skillRepo:     skillRepo,
		characterRepo: characterRepo,
	}
}

// GetInitialSkillSelection 获取初始技能选择机会（角色创建时）
func (s *SkillService) GetInitialSkillSelection(characterID int) (*models.SkillSelection, error) {
	character, err := s.characterRepo.GetByID(characterID)
	if err != nil {
		return nil, err
	}

	// 检查是否已经选择过初始技能
	existingSkills, err := s.skillRepo.GetCharacterSkills(characterID)
	if err == nil && len(existingSkills) > 0 {
		return nil, errors.New("初始技能已选择")
	}

	// 获取初始技能池
	initialSkills, err := s.skillRepo.GetInitialSkills(character.ClassID)
	if err != nil {
		return nil, err
	}

	return &models.SkillSelection{
		CharacterID:   characterID,
		Level:         1,
		SelectionType: "initial_active",
		CanUpgrade:    false,
		NewSkills:     initialSkills,
	}, nil
}

// SelectInitialSkills 选择初始技能（必须选择2个）
func (s *SkillService) SelectInitialSkills(req *models.InitialSkillSelectionRequest) error {
	if len(req.SkillIDs) != 2 {
		return errors.New("必须选择2个初始技能")
	}

	character, err := s.characterRepo.GetByID(req.CharacterID)
	if err != nil {
		return err
	}

	// 检查是否已经选择过初始技能
	existingSkills, err := s.skillRepo.GetCharacterSkills(req.CharacterID)
	if err == nil && len(existingSkills) > 0 {
		return errors.New("初始技能已选择")
	}

	// 验证技能是否在初始技能池中
	initialSkills, err := s.skillRepo.GetInitialSkills(character.ClassID)
	if err != nil {
		return err
	}

	initialSkillMap := make(map[string]bool)
	for _, skill := range initialSkills {
		initialSkillMap[skill.ID] = true
	}

	for _, skillID := range req.SkillIDs {
		if !initialSkillMap[skillID] {
			return fmt.Errorf("技能 %s 不在初始技能池中", skillID)
		}

		// 添加技能
		err := s.skillRepo.AddCharacterSkill(req.CharacterID, skillID, 1)
		if err != nil {
			return err
		}
	}

	return nil
}

// CheckSkillSelectionOpportunity 检查是否有技能选择机会
func (s *SkillService) CheckSkillSelectionOpportunity(characterID int) (*models.SkillSelection, error) {
	character, err := s.characterRepo.GetByID(characterID)
	if err != nil {
		return nil, err
	}

	level := character.Level

	// 检查被动技能选择机会（3的倍数）
	if level%3 == 0 && level >= 3 {
		return s.GetPassiveSkillSelection(characterID)
	}

	// 检查主动技能选择机会（5的倍数）
	if level%5 == 0 && level >= 5 {
		return s.GetActiveSkillSelection(characterID)
	}

	return nil, nil // 没有选择机会
}

// GetActiveSkillSelection 获取主动技能选择机会
func (s *SkillService) GetActiveSkillSelection(characterID int) (*models.SkillSelection, error) {
	character, err := s.characterRepo.GetByID(characterID)
	if err != nil {
		return nil, err
	}

	// 获取已学会的技能
	learnedSkills, err := s.skillRepo.GetCharacterSkills(characterID)
	if err != nil {
		return nil, err
	}

	// 筛选可升级的技能（等级<5）
	upgradeSkills := make([]*models.CharacterSkill, 0)
	learnedSkillIDs := make([]string, 0)
	for _, cs := range learnedSkills {
		learnedSkillIDs = append(learnedSkillIDs, cs.SkillID)
		if cs.SkillLevel < 5 {
			upgradeSkills = append(upgradeSkills, cs)
		}
	}

	// 随机获取4个新技能选项
	newSkills, err := s.skillRepo.GetRandomActiveSkills(character.ClassID, 4, learnedSkillIDs)
	if err != nil {
		return nil, err
	}

	return &models.SkillSelection{
		CharacterID:   characterID,
		Level:         character.Level,
		SelectionType: "active",
		CanUpgrade:    len(upgradeSkills) > 0,
		UpgradeSkills: upgradeSkills,
		NewSkills:     newSkills,
	}, nil
}

// GetPassiveSkillSelection 获取被动技能选择机会
func (s *SkillService) GetPassiveSkillSelection(characterID int) (*models.SkillSelection, error) {
	character, err := s.characterRepo.GetByID(characterID)
	if err != nil {
		return nil, err
	}

	// 获取已学会的被动技能
	learnedPassives, err := s.skillRepo.GetCharacterPassiveSkills(characterID)
	if err != nil {
		return nil, err
	}

	// 筛选可升级的被动技能（等级<5）
	upgradePassives := make([]*models.CharacterPassiveSkill, 0)
	learnedPassiveIDs := make([]string, 0)
	for _, cps := range learnedPassives {
		learnedPassiveIDs = append(learnedPassiveIDs, cps.PassiveID)
		if cps.Level < 5 {
			upgradePassives = append(upgradePassives, cps)
		}
	}

	// 随机获取4个新被动技能选项
	newPassives, err := s.skillRepo.GetRandomPassiveSkills(character.ClassID, 4, learnedPassiveIDs)
	if err != nil {
		return nil, err
	}

	return &models.SkillSelection{
		CharacterID:     characterID,
		Level:           character.Level,
		SelectionType:   "passive",
		CanUpgrade:      len(upgradePassives) > 0,
		UpgradePassives: upgradePassives,
		NewPassives:     newPassives,
	}, nil
}

// SelectSkill 选择技能（新技能或升级）
func (s *SkillService) SelectSkill(req *models.SkillSelectionRequest) error {
	character, err := s.characterRepo.GetByID(req.CharacterID)
	if err != nil {
		return err
	}

	level := character.Level

	// 验证选择机会
	if req.SkillID != "" {
		// 主动技能选择
		if level%5 != 0 || level < 5 {
			return errors.New("当前等级没有主动技能选择机会")
		}

		if req.IsUpgrade {
			// 升级现有技能
			return s.skillRepo.UpgradeCharacterSkill(req.CharacterID, req.SkillID)
		} else {
			// 学习新技能
			return s.skillRepo.AddCharacterSkill(req.CharacterID, req.SkillID, 1)
		}
	} else if req.PassiveID != "" {
		// 被动技能选择
		if level%3 != 0 || level < 3 {
			return errors.New("当前等级没有被动技能选择机会")
		}

		if req.IsUpgrade {
			// 升级现有被动技能
			return s.skillRepo.UpgradeCharacterPassiveSkill(req.CharacterID, req.PassiveID)
		} else {
			// 学习新被动技能
			return s.skillRepo.AddCharacterPassiveSkill(req.CharacterID, req.PassiveID, 1)
		}
	}

	return errors.New("必须指定skillId或passiveId")
}

// GetCharacterAllSkills 获取角色的所有技能（主动+被动）
func (s *SkillService) GetCharacterAllSkills(characterID int) (map[string]interface{}, error) {
	activeSkills, err := s.skillRepo.GetCharacterSkills(characterID)
	if err != nil {
		return nil, err
	}

	passiveSkills, err := s.skillRepo.GetCharacterPassiveSkills(characterID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"activeSkills":  activeSkills,
		"passiveSkills": passiveSkills,
	}, nil
}

// CalculateSkillEffect 计算技能效果（根据等级）
func (s *SkillService) CalculateSkillEffect(skill *models.Skill, skillLevel int) map[string]interface{} {
	// 基础效果
	effect := map[string]interface{}{
		"baseDamage":     skill.BaseValue,
		"scalingRatio":   skill.ScalingRatio,
		"resourceCost":   skill.ResourceCost,
		"cooldown":       skill.Cooldown,
		"threatModifier": skill.ThreatModifier,
		"threatType":     skill.ThreatType,
		"type":           skill.Type,
		"targetType":     skill.TargetType,
		"damageType":     skill.DamageType,
	}

	// 根据技能类型和等级计算升级效果
	switch skill.ID {
	case "warrior_heroic_strike":
		// 英勇打击：每级+15%伤害（1级100%，5级160%）
		effect["damageMultiplier"] = 1.0 + float64(skillLevel-1)*0.15
	case "warrior_taunt":
		// 嘲讽：每级-0.5回合冷却（1级2回合，5级0回合）
		cooldownReduction := int(float64(skillLevel-1) * 0.5)
		effect["cooldown"] = max(0, skill.Cooldown-cooldownReduction)
	case "warrior_shield_block":
		// 盾牌格挡：每级+5%减伤（1级30%，5级50%），持续2回合
		effect["damageReduction"] = 30.0 + float64(skillLevel-1)*5.0
		effect["duration"] = 2
	case "warrior_cleave":
		// 顺劈斩：每级+10%伤害（主目标1级100%，5级140%；相邻1级80%，5级120%）
		effect["mainTargetMultiplier"] = 1.0 + float64(skillLevel-1)*0.10
		effect["adjacentMultiplier"] = 0.8 + float64(skillLevel-1)*0.10
	case "warrior_slam":
		// 重击：每级+20%伤害（1级150%，5级230%）
		effect["damageMultiplier"] = 1.5 + float64(skillLevel-1)*0.20
	case "warrior_battle_shout":
		// 战斗怒吼：每级+2%攻击力，+1回合持续时间（1级10%持续5回合，5级18%持续9回合）
		effect["attackBonus"] = 10.0 + float64(skillLevel-1)*2.0
		effect["duration"] = 5 + (skillLevel - 1)
	case "warrior_demoralizing_shout":
		// 挫志怒吼：每级+3%降低比例，+0.5回合持续时间（1级15%持续3回合，5级27%持续4回合）
		effect["attackReduction"] = 15.0 + float64(skillLevel-1)*3.0
		effect["duration"] = 3.0 + float64(skillLevel-1)*0.5
	case "warrior_last_stand":
		// 破釜沉舟：每级+5%恢复量，+0.5回合持续时间（1级30%持续3回合，5级50%持续4回合）
		effect["healPercent"] = 30.0 + float64(skillLevel-1)*5.0
		effect["duration"] = 3.0 + float64(skillLevel-1)*0.5
		// duration保持为float64，在使用时转换为int
	case "warrior_charge":
		// 冲锋：每级+10%伤害，+3怒气，+5%眩晕概率，-0.5回合冷却
		effect["damageMultiplier"] = 0.8 + float64(skillLevel-1)*0.10
		effect["rageGain"] = 15 + (skillLevel-1)*3
		effect["stunChance"] = 0.3 + float64(skillLevel-1)*0.05
		cooldownReduction := int(float64(skillLevel-1) * 0.5)
		effect["cooldown"] = max(1, skill.Cooldown-cooldownReduction)
	case "warrior_whirlwind":
		// 旋风斩：每级+15%伤害，+2%防御降低，+0.5回合持续时间（1级100%伤害，10%防御降低持续2回合，5级160%伤害，18%防御降低持续3回合）
		effect["damageMultiplier"] = 1.0 + float64(skillLevel-1)*0.15
		effect["defenseReduction"] = 10.0 + float64(skillLevel-1)*2.0
		effect["debuffDuration"] = 2.0 + float64(skillLevel-1)*0.5
	case "warrior_shield_slam":
		// 盾牌猛击：每级+15%攻击力加成，+10%防御力加成（1级120%攻击+50%防御，5级195%攻击+90%防御）
		effect["attackMultiplier"] = 1.2 + float64(skillLevel-1)*0.15
		effect["defenseMultiplier"] = 0.5 + float64(skillLevel-1)*0.10
	case "warrior_execute":
		// 斩杀：每级+30%伤害（1级200%，5级320%），仅对HP<20%的敌人
		effect["damageMultiplier"] = 2.0 + float64(skillLevel-1)*0.30
		effect["hpThreshold"] = 0.20 // 20%血量阈值
	case "warrior_shield_wall":
		// 盾墙：每级+5%减伤，+0.5回合持续时间（1级60%持续2回合，5级80%持续3回合）
		effect["damageReduction"] = 60.0 + float64(skillLevel-1)*5.0
		effect["duration"] = 2.0 + float64(skillLevel-1)*0.5
	case "warrior_unbreakable_barrier":
		// 不灭壁垒：每级+10%护盾量，+0.5回合持续时间（1级50%持续4回合，5级90%持续5回合）
		effect["shieldPercent"] = 50.0 + float64(skillLevel-1)*10.0
		effect["duration"] = 4.0 + float64(skillLevel-1)*0.5
	case "warrior_shield_reflection":
		// 盾牌反射：每级+10%反射比例，+0.5回合持续时间（1级50%持续2回合，5级90%持续3回合）
		effect["reflectPercent"] = 50.0 + float64(skillLevel-1)*10.0
		effect["duration"] = 2.0 + float64(skillLevel-1)*0.5
	case "warrior_challenging_shout":
		// 挑战怒吼：每级+0.5回合持续时间，-0.5回合冷却（1级持续1回合冷却5回合，5级持续3回合冷却3回合）
		effect["duration"] = 1.0 + float64(skillLevel-1)*0.5
		cooldownReduction := int(float64(skillLevel-1) * 0.5)
		effect["cooldown"] = max(3, skill.Cooldown-cooldownReduction)
	case "warrior_recklessness":
		// 鲁莽：每级+10%暴击率，-2%受到伤害增加，+0.5回合持续时间（1级50%暴击+20%受伤持续3回合，5级90%暴击+12%受伤持续4回合）
		effect["critBonus"] = 50.0 + float64(skillLevel-1)*10.0
		effect["damageTakenIncrease"] = 20.0 - float64(skillLevel-1)*2.0
		effect["duration"] = 3.0 + float64(skillLevel-1)*0.5
	case "warrior_retaliation":
		// 反击风暴：每级+10%反击伤害，+0.5回合持续时间，-1回合冷却（1级50%持续3回合冷却10回合，5级90%持续4回合冷却6回合）
		effect["counterDamagePercent"] = 50.0 + float64(skillLevel-1)*10.0
		effect["duration"] = 3.0 + float64(skillLevel-1)*0.5
		effect["cooldown"] = max(6, skill.Cooldown-(skillLevel-1))
	case "warrior_berserker_rage":
		// 狂暴之怒：每级+5%攻击力，+1点额外怒气，+0.5回合持续时间（1级30%攻击+5怒气持续4回合，5级50%攻击+9怒气持续5回合）
		effect["attackBonus"] = 30.0 + float64(skillLevel-1)*5.0
		effect["ragePerHit"] = 5 + (skillLevel - 1)
		effect["duration"] = 4.0 + float64(skillLevel-1)*0.5
	case "warrior_avatar":
		// 天神下凡：每级+10%攻击力，+0.5回合持续时间，-1回合冷却（1级50%持续3回合冷却12回合，5级90%持续4回合冷却8回合）
		effect["attackBonus"] = 50.0 + float64(skillLevel-1)*10.0
		effect["duration"] = 3.0 + float64(skillLevel-1)*0.5
		effect["cooldown"] = max(8, skill.Cooldown-(skillLevel-1))
		effect["immuneCC"] = true
	case "warrior_mortal_strike":
		// 致死打击：每级+25%伤害，+5%降低治疗效果，+0.5回合持续时间（1级180%伤害，50%降低治疗持续3回合，5级280%伤害，70%降低治疗持续4回合）
		effect["damageMultiplier"] = 1.8 + float64(skillLevel-1)*0.25
		effect["healingReduction"] = 50.0 + float64(skillLevel-1)*5.0
		effect["debuffDuration"] = 3.0 + float64(skillLevel-1)*0.5
	case "warrior_bloodthirst":
		// 嗜血：每级+15%伤害，+5%恢复比例（1级120%伤害，30%恢复，5级180%伤害，50%恢复）
		effect["damageMultiplier"] = 1.2 + float64(skillLevel-1)*0.15
		effect["healPercent"] = 30.0 + float64(skillLevel-1)*5.0
	default:
		// 默认：如果没有特殊升级效果，使用基础值
		effect["damageMultiplier"] = skill.ScalingRatio
	}

	return effect
}

// Helper function
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

