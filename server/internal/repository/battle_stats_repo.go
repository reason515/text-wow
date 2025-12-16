package repository

import (
	"database/sql"
	"fmt"
	"time"

	"text-wow/internal/database"
	"text-wow/internal/models"
)

// BattleStatsRepository 战斗统计数据仓库
type BattleStatsRepository struct{}

// NewBattleStatsRepository 创建战斗统计仓库
func NewBattleStatsRepository() *BattleStatsRepository {
	return &BattleStatsRepository{}
}

// ═══════════════════════════════════════════════════════════
// 战斗记录 (battle_records)
// ═══════════════════════════════════════════════════════════

// CreateBattleRecord 创建战斗记录
func (r *BattleStatsRepository) CreateBattleRecord(record *models.BattleRecord) (int64, error) {
	result, err := database.DB.Exec(`
		INSERT INTO battle_records (
			user_id, zone_id, battle_type, monster_id, opponent_user_id,
			total_rounds, duration_seconds, result,
			team_damage_dealt, team_damage_taken, team_healing_done,
			exp_gained, gold_gained, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		record.UserID, record.ZoneID, record.BattleType, record.MonsterID, record.OpponentUserID,
		record.TotalRounds, record.DurationSeconds, record.Result,
		record.TeamDamageDealt, record.TeamDamageTaken, record.TeamHealingDone,
		record.ExpGained, record.GoldGained, time.Now(),
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// GetBattleRecordByID 根据ID获取战斗记录
func (r *BattleStatsRepository) GetBattleRecordByID(id int) (*models.BattleRecord, error) {
	record := &models.BattleRecord{}
	var monsterID sql.NullString
	var opponentUserID sql.NullInt64

	err := database.DB.QueryRow(`
		SELECT id, user_id, zone_id, battle_type, monster_id, opponent_user_id,
		       total_rounds, duration_seconds, result,
		       team_damage_dealt, team_damage_taken, team_healing_done,
		       exp_gained, gold_gained, created_at
		FROM battle_records WHERE id = ?`, id,
	).Scan(
		&record.ID, &record.UserID, &record.ZoneID, &record.BattleType, &monsterID, &opponentUserID,
		&record.TotalRounds, &record.DurationSeconds, &record.Result,
		&record.TeamDamageDealt, &record.TeamDamageTaken, &record.TeamHealingDone,
		&record.ExpGained, &record.GoldGained, &record.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if monsterID.Valid {
		record.MonsterID = monsterID.String
	}
	if opponentUserID.Valid {
		id := int(opponentUserID.Int64)
		record.OpponentUserID = &id
	}

	return record, nil
}

// GetRecentBattleRecords 获取用户最近的战斗记录
func (r *BattleStatsRepository) GetRecentBattleRecords(userID int, limit int) ([]*models.BattleRecord, error) {
	rows, err := database.DB.Query(`
		SELECT id, user_id, zone_id, battle_type, monster_id, opponent_user_id,
		       total_rounds, duration_seconds, result,
		       team_damage_dealt, team_damage_taken, team_healing_done,
		       exp_gained, gold_gained, created_at
		FROM battle_records 
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT ?`, userID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*models.BattleRecord
	for rows.Next() {
		record := &models.BattleRecord{}
		var monsterID sql.NullString
		var opponentUserID sql.NullInt64

		err := rows.Scan(
			&record.ID, &record.UserID, &record.ZoneID, &record.BattleType, &monsterID, &opponentUserID,
			&record.TotalRounds, &record.DurationSeconds, &record.Result,
			&record.TeamDamageDealt, &record.TeamDamageTaken, &record.TeamHealingDone,
			&record.ExpGained, &record.GoldGained, &record.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if monsterID.Valid {
			record.MonsterID = monsterID.String
		}
		if opponentUserID.Valid {
			id := int(opponentUserID.Int64)
			record.OpponentUserID = &id
		}

		records = append(records, record)
	}

	return records, nil
}

// ═══════════════════════════════════════════════════════════
// 战斗角色统计 (battle_character_stats)
// ═══════════════════════════════════════════════════════════

// CreateBattleCharacterStats 创建单场战斗角色统计
func (r *BattleStatsRepository) CreateBattleCharacterStats(stats *models.BattleCharacterStats) (int64, error) {
	result, err := database.DB.Exec(`
		INSERT INTO battle_character_stats (
			battle_id, character_id, team_slot,
			damage_dealt, physical_damage, magic_damage,
			fire_damage, frost_damage, shadow_damage, holy_damage, nature_damage, dot_damage,
			crit_count, crit_damage, max_crit,
			damage_taken, physical_taken, magic_taken, damage_blocked, damage_absorbed,
			dodge_count, block_count, hit_count,
			healing_done, healing_received, overhealing, self_healing, hot_healing,
			skill_uses, skill_hits, skill_misses,
			cc_applied, cc_received, dispels, interrupts,
			kills, deaths, resurrects, resource_used, resource_generated
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		stats.BattleID, stats.CharacterID, stats.TeamSlot,
		stats.DamageDealt, stats.PhysicalDamage, stats.MagicDamage,
		stats.FireDamage, stats.FrostDamage, stats.ShadowDamage, stats.HolyDamage, stats.NatureDamage, stats.DotDamage,
		stats.CritCount, stats.CritDamage, stats.MaxCrit,
		stats.DamageTaken, stats.PhysicalTaken, stats.MagicTaken, stats.DamageBlocked, stats.DamageAbsorbed,
		stats.DodgeCount, stats.BlockCount, stats.HitCount,
		stats.HealingDone, stats.HealingReceived, stats.Overhealing, stats.SelfHealing, stats.HotHealing,
		stats.SkillUses, stats.SkillHits, stats.SkillMisses,
		stats.CcApplied, stats.CcReceived, stats.Dispels, stats.Interrupts,
		stats.Kills, stats.Deaths, stats.Resurrects, stats.ResourceUsed, stats.ResourceGenerated,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// GetBattleCharacterStats 获取战斗的角色统计
func (r *BattleStatsRepository) GetBattleCharacterStats(battleID int) ([]*models.BattleCharacterStats, error) {
	rows, err := database.DB.Query(`
		SELECT id, battle_id, character_id, team_slot,
		       damage_dealt, physical_damage, magic_damage,
		       fire_damage, frost_damage, shadow_damage, holy_damage, nature_damage, dot_damage,
		       crit_count, crit_damage, max_crit,
		       damage_taken, physical_taken, magic_taken, damage_blocked, damage_absorbed,
		       dodge_count, block_count, hit_count,
		       healing_done, healing_received, overhealing, self_healing, hot_healing,
		       skill_uses, skill_hits, skill_misses,
		       cc_applied, cc_received, dispels, interrupts,
		       kills, deaths, resurrects, resource_used, resource_generated
		FROM battle_character_stats
		WHERE battle_id = ?
		ORDER BY team_slot`, battleID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statsList []*models.BattleCharacterStats
	for rows.Next() {
		stats := &models.BattleCharacterStats{}
		err := rows.Scan(
			&stats.ID, &stats.BattleID, &stats.CharacterID, &stats.TeamSlot,
			&stats.DamageDealt, &stats.PhysicalDamage, &stats.MagicDamage,
			&stats.FireDamage, &stats.FrostDamage, &stats.ShadowDamage, &stats.HolyDamage, &stats.NatureDamage, &stats.DotDamage,
			&stats.CritCount, &stats.CritDamage, &stats.MaxCrit,
			&stats.DamageTaken, &stats.PhysicalTaken, &stats.MagicTaken, &stats.DamageBlocked, &stats.DamageAbsorbed,
			&stats.DodgeCount, &stats.BlockCount, &stats.HitCount,
			&stats.HealingDone, &stats.HealingReceived, &stats.Overhealing, &stats.SelfHealing, &stats.HotHealing,
			&stats.SkillUses, &stats.SkillHits, &stats.SkillMisses,
			&stats.CcApplied, &stats.CcReceived, &stats.Dispels, &stats.Interrupts,
			&stats.Kills, &stats.Deaths, &stats.Resurrects, &stats.ResourceUsed, &stats.ResourceGenerated,
		)
		if err != nil {
			return nil, err
		}
		statsList = append(statsList, stats)
	}

	return statsList, nil
}

// ═══════════════════════════════════════════════════════════
// 角色生涯统计 (character_lifetime_stats)
// ═══════════════════════════════════════════════════════════

// GetOrCreateLifetimeStats 获取或创建角色生涯统计
func (r *BattleStatsRepository) GetOrCreateLifetimeStats(characterID int) (*models.CharacterLifetimeStats, error) {
	stats, err := r.GetLifetimeStats(characterID)
	if err == sql.ErrNoRows {
		// 创建新的生涯统计记录
		_, err = database.DB.Exec(`
			INSERT INTO character_lifetime_stats (character_id, updated_at)
			VALUES (?, ?)`, characterID, time.Now(),
		)
		if err != nil {
			return nil, err
		}
		return r.GetLifetimeStats(characterID)
	}
	return stats, err
}

// GetLifetimeStats 获取角色生涯统计
func (r *BattleStatsRepository) GetLifetimeStats(characterID int) (*models.CharacterLifetimeStats, error) {
	stats := &models.CharacterLifetimeStats{}
	var lastBattleAt sql.NullTime

	err := database.DB.QueryRow(`
		SELECT character_id,
		       total_battles, victories, defeats, pve_battles, pvp_battles, boss_kills,
		       total_damage_dealt, total_physical_damage, total_magic_damage,
		       total_crit_damage, total_crit_count, highest_damage_single, highest_damage_battle,
		       total_damage_taken, total_damage_blocked, total_damage_absorbed, total_dodge_count,
		       total_healing_done, total_healing_received, total_overhealing,
		       highest_healing_single, highest_healing_battle,
		       total_kills, total_deaths, kill_streak_best, current_kill_streak,
		       total_skill_uses, total_skill_hits,
		       total_resource_used, total_rounds, total_battle_time,
		       last_battle_at, updated_at
		FROM character_lifetime_stats
		WHERE character_id = ?`, characterID,
	).Scan(
		&stats.CharacterID,
		&stats.TotalBattles, &stats.Victories, &stats.Defeats, &stats.PveBattles, &stats.PvpBattles, &stats.BossKills,
		&stats.TotalDamageDealt, &stats.TotalPhysicalDamage, &stats.TotalMagicDamage,
		&stats.TotalCritDamage, &stats.TotalCritCount, &stats.HighestDamageSingle, &stats.HighestDamageBattle,
		&stats.TotalDamageTaken, &stats.TotalDamageBlocked, &stats.TotalDamageAbsorbed, &stats.TotalDodgeCount,
		&stats.TotalHealingDone, &stats.TotalHealingReceived, &stats.TotalOverhealing,
		&stats.HighestHealingSingle, &stats.HighestHealingBattle,
		&stats.TotalKills, &stats.TotalDeaths, &stats.KillStreakBest, &stats.CurrentKillStrk,
		&stats.TotalSkillUses, &stats.TotalSkillHits,
		&stats.TotalResourceUsed, &stats.TotalRounds, &stats.TotalBattleTime,
		&lastBattleAt, &stats.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if lastBattleAt.Valid {
		stats.LastBattleAt = &lastBattleAt.Time
	}

	return stats, nil
}

// UpdateLifetimeStats 更新角色生涯统计（增量更新）
func (r *BattleStatsRepository) UpdateLifetimeStats(characterID int, battleStats *models.BattleCharacterStats, isVictory bool, battleType string, totalRounds int) error {
	// 先获取或创建
	current, err := r.GetOrCreateLifetimeStats(characterID)
	if err != nil {
		return err
	}

	// 计算新的统计值
	newTotalBattles := current.TotalBattles + 1
	newVictories := current.Victories
	newDefeats := current.Defeats
	if isVictory {
		newVictories++
	} else {
		newDefeats++
	}

	newPveBattles := current.PveBattles
	newPvpBattles := current.PvpBattles
	if battleType == "pvp" {
		newPvpBattles++
	} else {
		newPveBattles++
	}

	// 更新最高伤害记录
	newHighestDamageSingle := current.HighestDamageSingle
	if battleStats.MaxCrit > newHighestDamageSingle {
		newHighestDamageSingle = battleStats.MaxCrit
	}

	newHighestDamageBattle := current.HighestDamageBattle
	if battleStats.DamageDealt > newHighestDamageBattle {
		newHighestDamageBattle = battleStats.DamageDealt
	}

	// 更新连杀
	newCurrentKillStreak := current.CurrentKillStrk
	newKillStreakBest := current.KillStreakBest
	if battleStats.Deaths > 0 {
		newCurrentKillStreak = 0
	} else {
		newCurrentKillStreak += battleStats.Kills
	}
	if newCurrentKillStreak > newKillStreakBest {
		newKillStreakBest = newCurrentKillStreak
	}

	now := time.Now()
	_, err = database.DB.Exec(`
		UPDATE character_lifetime_stats SET
			total_battles = ?,
			victories = ?,
			defeats = ?,
			pve_battles = ?,
			pvp_battles = ?,
			total_damage_dealt = total_damage_dealt + ?,
			total_physical_damage = total_physical_damage + ?,
			total_magic_damage = total_magic_damage + ?,
			total_crit_damage = total_crit_damage + ?,
			total_crit_count = total_crit_count + ?,
			highest_damage_single = ?,
			highest_damage_battle = ?,
			total_damage_taken = total_damage_taken + ?,
			total_damage_blocked = total_damage_blocked + ?,
			total_damage_absorbed = total_damage_absorbed + ?,
			total_dodge_count = total_dodge_count + ?,
			total_healing_done = total_healing_done + ?,
			total_healing_received = total_healing_received + ?,
			total_overhealing = total_overhealing + ?,
			total_kills = total_kills + ?,
			total_deaths = total_deaths + ?,
			kill_streak_best = ?,
			current_kill_streak = ?,
			total_skill_uses = total_skill_uses + ?,
			total_skill_hits = total_skill_hits + ?,
			total_resource_used = total_resource_used + ?,
			total_rounds = total_rounds + ?,
			last_battle_at = ?,
			updated_at = ?
		WHERE character_id = ?`,
		newTotalBattles, newVictories, newDefeats, newPveBattles, newPvpBattles,
		battleStats.DamageDealt, battleStats.PhysicalDamage, battleStats.MagicDamage,
		battleStats.CritDamage, battleStats.CritCount,
		newHighestDamageSingle, newHighestDamageBattle,
		battleStats.DamageTaken, battleStats.DamageBlocked, battleStats.DamageAbsorbed, battleStats.DodgeCount,
		battleStats.HealingDone, battleStats.HealingReceived, battleStats.Overhealing,
		battleStats.Kills, battleStats.Deaths,
		newKillStreakBest, newCurrentKillStreak,
		battleStats.SkillUses, battleStats.SkillHits,
		battleStats.ResourceUsed, totalRounds,
		now, now, characterID,
	)

	return err
}

// GetLifetimeStatsByUserID 获取用户所有角色的生涯统计
func (r *BattleStatsRepository) GetLifetimeStatsByUserID(userID int) ([]*models.CharacterLifetimeStats, error) {
	rows, err := database.DB.Query(`
		SELECT cls.character_id,
		       cls.total_battles, cls.victories, cls.defeats, cls.pve_battles, cls.pvp_battles, cls.boss_kills,
		       cls.total_damage_dealt, cls.total_physical_damage, cls.total_magic_damage,
		       cls.total_crit_damage, cls.total_crit_count, cls.highest_damage_single, cls.highest_damage_battle,
		       cls.total_damage_taken, cls.total_damage_blocked, cls.total_damage_absorbed, cls.total_dodge_count,
		       cls.total_healing_done, cls.total_healing_received, cls.total_overhealing,
		       cls.highest_healing_single, cls.highest_healing_battle,
		       cls.total_kills, cls.total_deaths, cls.kill_streak_best, cls.current_kill_streak,
		       cls.total_skill_uses, cls.total_skill_hits,
		       cls.total_resource_used, cls.total_rounds, cls.total_battle_time,
		       cls.last_battle_at, cls.updated_at
		FROM character_lifetime_stats cls
		JOIN characters c ON cls.character_id = c.id
		WHERE c.user_id = ?
		ORDER BY c.team_slot`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statsList []*models.CharacterLifetimeStats
	for rows.Next() {
		stats := &models.CharacterLifetimeStats{}
		var lastBattleAt sql.NullTime

		err := rows.Scan(
			&stats.CharacterID,
			&stats.TotalBattles, &stats.Victories, &stats.Defeats, &stats.PveBattles, &stats.PvpBattles, &stats.BossKills,
			&stats.TotalDamageDealt, &stats.TotalPhysicalDamage, &stats.TotalMagicDamage,
			&stats.TotalCritDamage, &stats.TotalCritCount, &stats.HighestDamageSingle, &stats.HighestDamageBattle,
			&stats.TotalDamageTaken, &stats.TotalDamageBlocked, &stats.TotalDamageAbsorbed, &stats.TotalDodgeCount,
			&stats.TotalHealingDone, &stats.TotalHealingReceived, &stats.TotalOverhealing,
			&stats.HighestHealingSingle, &stats.HighestHealingBattle,
			&stats.TotalKills, &stats.TotalDeaths, &stats.KillStreakBest, &stats.CurrentKillStrk,
			&stats.TotalSkillUses, &stats.TotalSkillHits,
			&stats.TotalResourceUsed, &stats.TotalRounds, &stats.TotalBattleTime,
			&lastBattleAt, &stats.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if lastBattleAt.Valid {
			stats.LastBattleAt = &lastBattleAt.Time
		}

		statsList = append(statsList, stats)
	}

	return statsList, nil
}

// ═══════════════════════════════════════════════════════════
// 战斗技能明细 (battle_skill_breakdown)
// ═══════════════════════════════════════════════════════════

// CreateBattleSkillBreakdown 创建战斗技能明细
func (r *BattleStatsRepository) CreateBattleSkillBreakdown(breakdown *models.BattleSkillBreakdown) (int64, error) {
	result, err := database.DB.Exec(`
		INSERT INTO battle_skill_breakdown (
			battle_id, character_id, skill_id,
			use_count, hit_count, crit_count,
			total_damage, total_healing, resource_cost
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		breakdown.BattleID, breakdown.CharacterID, breakdown.SkillID,
		breakdown.UseCount, breakdown.HitCount, breakdown.CritCount,
		breakdown.TotalDamage, breakdown.TotalHealing, breakdown.ResourceCost,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// GetBattleSkillBreakdown 获取战斗的技能明细
func (r *BattleStatsRepository) GetBattleSkillBreakdown(battleID int, characterID int) ([]*models.BattleSkillBreakdown, error) {
	rows, err := database.DB.Query(`
		SELECT bsb.id, bsb.battle_id, bsb.character_id, bsb.skill_id,
		       bsb.use_count, bsb.hit_count, bsb.crit_count,
		       bsb.total_damage, bsb.total_healing, bsb.resource_cost,
		       COALESCE(s.name, 
		         CASE 
		           WHEN bsb.skill_id = '' OR bsb.skill_id IS NULL THEN '普通攻击'
		           ELSE bsb.skill_id 
		         END
		       ) as skill_name
		FROM battle_skill_breakdown bsb
		LEFT JOIN skills s ON bsb.skill_id = s.id
		WHERE bsb.battle_id = ? AND bsb.character_id = ?
		ORDER BY bsb.total_damage DESC`, battleID, characterID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var breakdowns []*models.BattleSkillBreakdown
	for rows.Next() {
		b := &models.BattleSkillBreakdown{}
		err := rows.Scan(
			&b.ID, &b.BattleID, &b.CharacterID, &b.SkillID,
			&b.UseCount, &b.HitCount, &b.CritCount,
			&b.TotalDamage, &b.TotalHealing, &b.ResourceCost,
			&b.SkillName,
		)
		if err != nil {
			return nil, err
		}

		// 如果skillName仍然为空或等于skillID（说明技能表中没有该技能），且skillID为空，设置为"普通攻击"
		if b.SkillName == "" || (b.SkillName == b.SkillID && (b.SkillID == "" || b.SkillID == "normal_attack")) {
			b.SkillName = "普通攻击"
		}

		breakdowns = append(breakdowns, b)
	}

	return breakdowns, nil
}

// ═══════════════════════════════════════════════════════════
// 每日统计 (daily_statistics)
// ═══════════════════════════════════════════════════════════

// GetOrCreateDailyStats 获取或创建今日统计
func (r *BattleStatsRepository) GetOrCreateDailyStats(userID int, date string) (*models.DailyStatistics, error) {
	stats, err := r.GetDailyStats(userID, date)
	if err == sql.ErrNoRows {
		// 创建新的每日统计记录
		_, err = database.DB.Exec(`
			INSERT INTO daily_statistics (user_id, stat_date)
			VALUES (?, ?)`, userID, date,
		)
		if err != nil {
			return nil, err
		}
		return r.GetDailyStats(userID, date)
	}
	return stats, err
}

// GetDailyStats 获取指定日期的统计
func (r *BattleStatsRepository) GetDailyStats(userID int, date string) (*models.DailyStatistics, error) {
	stats := &models.DailyStatistics{}

	err := database.DB.QueryRow(`
		SELECT id, user_id, stat_date, battles_count, victories, defeats,
		       total_damage, total_healing, total_damage_taken,
		       exp_gained, gold_gained, play_time, kills, deaths
		FROM daily_statistics
		WHERE user_id = ? AND stat_date = ?`, userID, date,
	).Scan(
		&stats.ID, &stats.UserID, &stats.StatDate, &stats.BattlesCount, &stats.Victories, &stats.Defeats,
		&stats.TotalDamage, &stats.TotalHealing, &stats.TotalDamageTaken,
		&stats.ExpGained, &stats.GoldGained, &stats.PlayTime, &stats.Kills, &stats.Deaths,
	)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

// UpdateDailyStats 更新每日统计（增量更新）
func (r *BattleStatsRepository) UpdateDailyStats(userID int, date string, isVictory bool, damage, healing, damageTaken, exp, gold, kills, deaths int) error {
	// 确保记录存在
	_, err := r.GetOrCreateDailyStats(userID, date)
	if err != nil {
		return err
	}

	victories := 0
	defeats := 0
	if isVictory {
		victories = 1
	} else {
		defeats = 1
	}

	_, err = database.DB.Exec(`
		UPDATE daily_statistics SET
			battles_count = battles_count + 1,
			victories = victories + ?,
			defeats = defeats + ?,
			total_damage = total_damage + ?,
			total_healing = total_healing + ?,
			total_damage_taken = total_damage_taken + ?,
			exp_gained = exp_gained + ?,
			gold_gained = gold_gained + ?,
			kills = kills + ?,
			deaths = deaths + ?
		WHERE user_id = ? AND stat_date = ?`,
		victories, defeats, damage, healing, damageTaken, exp, gold, kills, deaths, userID, date,
	)

	return err
}

// GetDailyStatsRange 获取日期范围内的统计
func (r *BattleStatsRepository) GetDailyStatsRange(userID int, startDate, endDate string) ([]*models.DailyStatistics, error) {
	rows, err := database.DB.Query(`
		SELECT id, user_id, stat_date, battles_count, victories, defeats,
		       total_damage, total_healing, total_damage_taken,
		       exp_gained, gold_gained, play_time, kills, deaths
		FROM daily_statistics
		WHERE user_id = ? AND stat_date >= ? AND stat_date <= ?
		ORDER BY stat_date DESC`, userID, startDate, endDate,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statsList []*models.DailyStatistics
	for rows.Next() {
		stats := &models.DailyStatistics{}
		err := rows.Scan(
			&stats.ID, &stats.UserID, &stats.StatDate, &stats.BattlesCount, &stats.Victories, &stats.Defeats,
			&stats.TotalDamage, &stats.TotalHealing, &stats.TotalDamageTaken,
			&stats.ExpGained, &stats.GoldGained, &stats.PlayTime, &stats.Kills, &stats.Deaths,
		)
		if err != nil {
			return nil, err
		}
		statsList = append(statsList, stats)
	}

	return statsList, nil
}

// ═══════════════════════════════════════════════════════════
// 统计汇总查询
// ═══════════════════════════════════════════════════════════

// GetCharacterBattleSummary 获取角色战斗摘要
func (r *BattleStatsRepository) GetCharacterBattleSummary(characterID int) (*models.CharacterBattleSummary, error) {
	stats, err := r.GetLifetimeStats(characterID)
	if err != nil {
		if err == sql.ErrNoRows {
			// 没有统计数据，返回空摘要
			return &models.CharacterBattleSummary{CharacterID: characterID}, nil
		}
		return nil, err
	}

	summary := &models.CharacterBattleSummary{
		CharacterID:  characterID,
		TotalBattles: stats.TotalBattles,
		Victories:    stats.Victories,
		TotalDamage:  stats.TotalDamageDealt,
		TotalHealing: stats.TotalHealingDone,
		TotalKills:   stats.TotalKills,
		TotalDeaths:  stats.TotalDeaths,
	}

	// 计算胜率
	if stats.TotalBattles > 0 {
		summary.WinRate = float64(stats.Victories) / float64(stats.TotalBattles) * 100
	}

	// 计算 K/D 比
	if stats.TotalDeaths > 0 {
		summary.KDRatio = float64(stats.TotalKills) / float64(stats.TotalDeaths)
	} else if stats.TotalKills > 0 {
		summary.KDRatio = float64(stats.TotalKills)
	}

	// 计算平均 DPS 和 HPS
	if stats.TotalRounds > 0 {
		summary.AvgDPS = float64(stats.TotalDamageDealt) / float64(stats.TotalRounds)
		summary.AvgHPS = float64(stats.TotalHealingDone) / float64(stats.TotalRounds)
	}

	return summary, nil
}

// GetBattleStatsOverview 获取战斗统计概览
func (r *BattleStatsRepository) GetBattleStatsOverview(userID int) (*models.BattleStatsOverview, error) {
	overview := &models.BattleStatsOverview{}

	// 获取角色生涯统计
	lifetimeStats, err := r.GetLifetimeStatsByUserID(userID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	overview.LifetimeStats = lifetimeStats

	// 获取今日统计
	today := time.Now().Format("2006-01-02")
	todayStats, err := r.GetDailyStats(userID, today)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	overview.TodayStats = todayStats

	// 获取最近10场战斗
	recentBattles, err := r.GetRecentBattleRecords(userID, 10)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	overview.RecentBattles = recentBattles

	return overview, nil
}

// ═══════════════════════════════════════════════════════════
// DPS分析
// ═══════════════════════════════════════════════════════════

// GetBattleDPSAnalysis 获取战斗DPS分析
func (r *BattleStatsRepository) GetBattleDPSAnalysis(battleID int) (*models.BattleDPSAnalysis, error) {
	// 获取战斗记录
	battle, err := r.GetBattleRecordByID(battleID)
	if err != nil {
		return nil, err
	}

	// 获取角色统计
	charStats, err := r.GetBattleCharacterStats(battleID)
	if err != nil {
		return nil, err
	}

	analysis := &models.BattleDPSAnalysis{
		BattleID:    battleID,
		Duration:    battle.DurationSeconds,
		TotalRounds: battle.TotalRounds,
		Characters:  make([]*models.CharacterDPSAnalysis, 0),
	}

	// 计算队伍总DPS和HPS
	if battle.DurationSeconds > 0 {
		analysis.TeamDPS = float64(battle.TeamDamageDealt) / float64(battle.DurationSeconds)
		analysis.TeamHPS = float64(battle.TeamHealingDone) / float64(battle.DurationSeconds)
	}

	// 队伍伤害构成
	teamComposition := &models.DamageComposition{
		Percentages: make(map[string]float64),
	}

	// 处理每个角色的DPS分析
	for _, charStat := range charStats {
		// 查询角色名称
		var characterName string
		err := database.DB.QueryRow(`SELECT name FROM characters WHERE id = ?`, charStat.CharacterID).Scan(&characterName)
		if err != nil {
			// 如果查询失败，使用默认名称
			characterName = fmt.Sprintf("角色 #%d", charStat.CharacterID)
		}

		charAnalysis := &models.CharacterDPSAnalysis{
			CharacterID:    charStat.CharacterID,
			CharacterName:  characterName,
			TotalDamage:    charStat.DamageDealt,
			TotalHealing:   charStat.HealingDone,
			Duration:       battle.DurationSeconds,
			SkillBreakdown: make([]*models.SkillDPSAnalysis, 0),
		}

		// 计算角色DPS和HPS
		if battle.DurationSeconds > 0 {
			charAnalysis.TotalDPS = float64(charStat.DamageDealt) / float64(battle.DurationSeconds)
			charAnalysis.TotalHPS = float64(charStat.HealingDone) / float64(battle.DurationSeconds)
		}

		// 获取技能明细
		skillBreakdowns, err := r.GetBattleSkillBreakdown(battleID, charStat.CharacterID)
		if err == nil {
			for _, skill := range skillBreakdowns {
				skillAnalysis := &models.SkillDPSAnalysis{
					SkillID:      skill.SkillID,
					SkillName:    skill.SkillName,
					TotalDamage:  skill.TotalDamage,
					UseCount:     skill.UseCount,
					HitCount:     skill.HitCount,
					CritCount:    skill.CritCount,
					ResourceCost: skill.ResourceCost,
				}

				// 计算平均伤害
				if skill.HitCount > 0 {
					skillAnalysis.AvgDamage = float64(skill.TotalDamage) / float64(skill.HitCount)
				}

				// 计算DPS
				if battle.DurationSeconds > 0 {
					skillAnalysis.DPS = float64(skill.TotalDamage) / float64(battle.DurationSeconds)
				}

				// 计算伤害占比
				if charStat.DamageDealt > 0 {
					skillAnalysis.DamagePercent = float64(skill.TotalDamage) / float64(charStat.DamageDealt) * 100
				}

				// 计算每点能量伤害
				if skill.ResourceCost > 0 {
					skillAnalysis.DamagePerResource = float64(skill.TotalDamage) / float64(skill.ResourceCost)
				}

				// 计算命中率
				if skill.UseCount > 0 {
					skillAnalysis.HitRate = float64(skill.HitCount) / float64(skill.UseCount) * 100
				}

				// 计算暴击率
				if skill.HitCount > 0 {
					skillAnalysis.CritRate = float64(skill.CritCount) / float64(skill.HitCount) * 100
				}

				charAnalysis.SkillBreakdown = append(charAnalysis.SkillBreakdown, skillAnalysis)
			}
		}

		// 角色伤害构成
		charComposition := &models.DamageComposition{
			Physical:    charStat.PhysicalDamage,
			Magic:       charStat.MagicDamage,
			Fire:        charStat.FireDamage,
			Frost:       charStat.FrostDamage,
			Shadow:      charStat.ShadowDamage,
			Holy:        charStat.HolyDamage,
			Nature:      charStat.NatureDamage,
			Dot:         charStat.DotDamage,
			Total:       charStat.DamageDealt,
			Percentages: make(map[string]float64),
		}

		// 计算各类型占比（确保Percentages已初始化）
		charComposition.Percentages = make(map[string]float64)
		if charComposition.Total > 0 {
			charComposition.Percentages["physical"] = float64(charComposition.Physical) / float64(charComposition.Total) * 100
			charComposition.Percentages["magic"] = float64(charComposition.Magic) / float64(charComposition.Total) * 100
			charComposition.Percentages["fire"] = float64(charComposition.Fire) / float64(charComposition.Total) * 100
			charComposition.Percentages["frost"] = float64(charComposition.Frost) / float64(charComposition.Total) * 100
			charComposition.Percentages["shadow"] = float64(charComposition.Shadow) / float64(charComposition.Total) * 100
			charComposition.Percentages["holy"] = float64(charComposition.Holy) / float64(charComposition.Total) * 100
			charComposition.Percentages["nature"] = float64(charComposition.Nature) / float64(charComposition.Total) * 100
			charComposition.Percentages["dot"] = float64(charComposition.Dot) / float64(charComposition.Total) * 100
		}

		charAnalysis.DamageComposition = charComposition

		// 累加队伍伤害构成
		teamComposition.Physical += charStat.PhysicalDamage
		teamComposition.Magic += charStat.MagicDamage
		teamComposition.Fire += charStat.FireDamage
		teamComposition.Frost += charStat.FrostDamage
		teamComposition.Shadow += charStat.ShadowDamage
		teamComposition.Holy += charStat.HolyDamage
		teamComposition.Nature += charStat.NatureDamage
		teamComposition.Dot += charStat.DotDamage
		teamComposition.Total += charStat.DamageDealt

		analysis.Characters = append(analysis.Characters, charAnalysis)
	}

	// 计算队伍伤害构成占比
	teamComposition.Percentages = make(map[string]float64)
	if teamComposition.Total > 0 {
		teamComposition.Percentages["physical"] = float64(teamComposition.Physical) / float64(teamComposition.Total) * 100
		teamComposition.Percentages["magic"] = float64(teamComposition.Magic) / float64(teamComposition.Total) * 100
		teamComposition.Percentages["fire"] = float64(teamComposition.Fire) / float64(teamComposition.Total) * 100
		teamComposition.Percentages["frost"] = float64(teamComposition.Frost) / float64(teamComposition.Total) * 100
		teamComposition.Percentages["shadow"] = float64(teamComposition.Shadow) / float64(teamComposition.Total) * 100
		teamComposition.Percentages["holy"] = float64(teamComposition.Holy) / float64(teamComposition.Total) * 100
		teamComposition.Percentages["nature"] = float64(teamComposition.Nature) / float64(teamComposition.Total) * 100
		teamComposition.Percentages["dot"] = float64(teamComposition.Dot) / float64(teamComposition.Total) * 100
	}

	analysis.TeamDamageComposition = teamComposition

	return analysis, nil
}

// GetCumulativeDPSAnalysis 获取累计DPS分析（按时间范围）
func (r *BattleStatsRepository) GetCumulativeDPSAnalysis(userID int, startTime time.Time) (*models.BattleDPSAnalysis, error) {
	// 获取时间范围内的所有战斗记录
	rows, err := database.DB.Query(`
		SELECT id, user_id, zone_id, battle_type, monster_id, opponent_user_id,
		       total_rounds, duration_seconds, result,
		       team_damage_dealt, team_damage_taken, team_healing_done,
		       exp_gained, gold_gained, created_at
		FROM battle_records 
		WHERE user_id = ? AND created_at >= ?
		ORDER BY created_at ASC`, userID, startTime,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var battles []*models.BattleRecord
	var totalDuration int
	var totalRounds int
	var totalDamage int
	var totalHealing int

	for rows.Next() {
		battle := &models.BattleRecord{}
		var monsterID sql.NullString
		var opponentUserID sql.NullInt64

		err := rows.Scan(
			&battle.ID, &battle.UserID, &battle.ZoneID, &battle.BattleType, &monsterID, &opponentUserID,
			&battle.TotalRounds, &battle.DurationSeconds, &battle.Result,
			&battle.TeamDamageDealt, &battle.TeamDamageTaken, &battle.TeamHealingDone,
			&battle.ExpGained, &battle.GoldGained, &battle.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if monsterID.Valid {
			battle.MonsterID = monsterID.String
		}
		if opponentUserID.Valid {
			id := int(opponentUserID.Int64)
			battle.OpponentUserID = &id
		}

		battles = append(battles, battle)
		totalDuration += battle.DurationSeconds
		totalRounds += battle.TotalRounds
		totalDamage += battle.TeamDamageDealt
		totalHealing += battle.TeamHealingDone
	}

	battleCount := len(battles)

	if battleCount == 0 {
		// 没有战斗记录，返回空分析
		return &models.BattleDPSAnalysis{
			BattleID:    0,
			Duration:    0,
			TotalRounds: 0,
			TeamDPS:     0,
			TeamHPS:     0,
			Characters:  make([]*models.CharacterDPSAnalysis, 0),
			TeamDamageComposition: &models.DamageComposition{
				Percentages: make(map[string]float64),
			},
		}, nil
	}

	// 获取所有战斗的角色统计（按角色ID聚合）
	charStatsMap := make(map[int]*models.BattleCharacterStats)
	characterIDs := make(map[int]bool)

	for _, battle := range battles {
		charStats, err := r.GetBattleCharacterStats(battle.ID)
		if err != nil {
			continue
		}

		for _, charStat := range charStats {
			characterIDs[charStat.CharacterID] = true
			if existing, exists := charStatsMap[charStat.CharacterID]; exists {
				// 累加统计数据
				existing.DamageDealt += charStat.DamageDealt
				existing.PhysicalDamage += charStat.PhysicalDamage
				existing.MagicDamage += charStat.MagicDamage
				existing.FireDamage += charStat.FireDamage
				existing.FrostDamage += charStat.FrostDamage
				existing.ShadowDamage += charStat.ShadowDamage
				existing.HolyDamage += charStat.HolyDamage
				existing.NatureDamage += charStat.NatureDamage
				existing.DotDamage += charStat.DotDamage
				existing.CritCount += charStat.CritCount
				existing.CritDamage += charStat.CritDamage
				if charStat.MaxCrit > existing.MaxCrit {
					existing.MaxCrit = charStat.MaxCrit
				}
				existing.DamageTaken += charStat.DamageTaken
				existing.PhysicalTaken += charStat.PhysicalTaken
				existing.MagicTaken += charStat.MagicTaken
				existing.DamageBlocked += charStat.DamageBlocked
				existing.DamageAbsorbed += charStat.DamageAbsorbed
				existing.DodgeCount += charStat.DodgeCount
				existing.BlockCount += charStat.BlockCount
				existing.HitCount += charStat.HitCount
				existing.HealingDone += charStat.HealingDone
				existing.HealingReceived += charStat.HealingReceived
				existing.Overhealing += charStat.Overhealing
				existing.SelfHealing += charStat.SelfHealing
				existing.HotHealing += charStat.HotHealing
				existing.SkillUses += charStat.SkillUses
				existing.SkillHits += charStat.SkillHits
				existing.SkillMisses += charStat.SkillMisses
				existing.CcApplied += charStat.CcApplied
				existing.CcReceived += charStat.CcReceived
				existing.Dispels += charStat.Dispels
				existing.Interrupts += charStat.Interrupts
				existing.Kills += charStat.Kills
				existing.Deaths += charStat.Deaths
				existing.Resurrects += charStat.Resurrects
				existing.ResourceUsed += charStat.ResourceUsed
				existing.ResourceGenerated += charStat.ResourceGenerated
			} else {
				// 创建新的统计记录
				newStat := &models.BattleCharacterStats{
					CharacterID:       charStat.CharacterID,
					TeamSlot:          charStat.TeamSlot,
					DamageDealt:       charStat.DamageDealt,
					PhysicalDamage:    charStat.PhysicalDamage,
					MagicDamage:       charStat.MagicDamage,
					FireDamage:        charStat.FireDamage,
					FrostDamage:       charStat.FrostDamage,
					ShadowDamage:      charStat.ShadowDamage,
					HolyDamage:        charStat.HolyDamage,
					NatureDamage:      charStat.NatureDamage,
					DotDamage:         charStat.DotDamage,
					CritCount:         charStat.CritCount,
					CritDamage:        charStat.CritDamage,
					MaxCrit:           charStat.MaxCrit,
					DamageTaken:       charStat.DamageTaken,
					PhysicalTaken:     charStat.PhysicalTaken,
					MagicTaken:        charStat.MagicTaken,
					DamageBlocked:     charStat.DamageBlocked,
					DamageAbsorbed:    charStat.DamageAbsorbed,
					DodgeCount:        charStat.DodgeCount,
					BlockCount:        charStat.BlockCount,
					HitCount:          charStat.HitCount,
					HealingDone:       charStat.HealingDone,
					HealingReceived:   charStat.HealingReceived,
					Overhealing:       charStat.Overhealing,
					SelfHealing:       charStat.SelfHealing,
					HotHealing:        charStat.HotHealing,
					SkillUses:         charStat.SkillUses,
					SkillHits:         charStat.SkillHits,
					SkillMisses:       charStat.SkillMisses,
					CcApplied:         charStat.CcApplied,
					CcReceived:        charStat.CcReceived,
					Dispels:           charStat.Dispels,
					Interrupts:        charStat.Interrupts,
					Kills:             charStat.Kills,
					Deaths:            charStat.Deaths,
					Resurrects:        charStat.Resurrects,
					ResourceUsed:      charStat.ResourceUsed,
					ResourceGenerated: charStat.ResourceGenerated,
				}
				charStatsMap[charStat.CharacterID] = newStat
			}
		}
	}

	// 构建累计DPS分析
	analysis := &models.BattleDPSAnalysis{
		BattleID:    0, // 累计分析没有特定的battleID
		Duration:    totalDuration,
		TotalRounds: totalRounds,
		BattleCount: battleCount, // 累计战斗场次
		Characters:  make([]*models.CharacterDPSAnalysis, 0),
	}

	// 计算队伍总DPS和HPS
	if totalDuration > 0 {
		analysis.TeamDPS = float64(totalDamage) / float64(totalDuration)
		analysis.TeamHPS = float64(totalHealing) / float64(totalDuration)
	}

	// 队伍伤害构成
	teamComposition := &models.DamageComposition{
		Percentages: make(map[string]float64),
	}

	// 处理每个角色的累计DPS分析
	for charID := range characterIDs {
		charStat := charStatsMap[charID]

		// 查询角色名称
		var characterName string
		err := database.DB.QueryRow(`SELECT name FROM characters WHERE id = ?`, charID).Scan(&characterName)
		if err != nil {
			characterName = fmt.Sprintf("角色 #%d", charID)
		}

		charAnalysis := &models.CharacterDPSAnalysis{
			CharacterID:    charID,
			CharacterName:  characterName,
			TotalDamage:    charStat.DamageDealt,
			TotalHealing:   charStat.HealingDone,
			Duration:       totalDuration,
			SkillBreakdown: make([]*models.SkillDPSAnalysis, 0),
		}

		// 计算角色DPS和HPS
		if totalDuration > 0 {
			charAnalysis.TotalDPS = float64(charStat.DamageDealt) / float64(totalDuration)
			charAnalysis.TotalHPS = float64(charStat.HealingDone) / float64(totalDuration)
		}

		// 获取累计技能明细（聚合所有战斗的技能数据）
		skillBreakdownMap := make(map[string]*models.SkillDPSAnalysis)

		for _, battle := range battles {
			skillBreakdowns, err := r.GetBattleSkillBreakdown(battle.ID, charID)
			if err != nil {
				continue
			}

			for _, skill := range skillBreakdowns {
				if existing, exists := skillBreakdownMap[skill.SkillID]; exists {
					// 累加技能统计
					existing.TotalDamage += skill.TotalDamage
					existing.UseCount += skill.UseCount
					existing.HitCount += skill.HitCount
					existing.CritCount += skill.CritCount
					existing.ResourceCost += skill.ResourceCost
				} else {
					// 创建新的技能统计
					skillAnalysis := &models.SkillDPSAnalysis{
						SkillID:      skill.SkillID,
						SkillName:    skill.SkillName,
						TotalDamage:  skill.TotalDamage,
						UseCount:     skill.UseCount,
						HitCount:     skill.HitCount,
						CritCount:    skill.CritCount,
						ResourceCost: skill.ResourceCost,
					}
					skillBreakdownMap[skill.SkillID] = skillAnalysis
				}
			}
		}

		// 计算技能DPS和其他指标
		for _, skillAnalysis := range skillBreakdownMap {
			// 计算平均伤害
			if skillAnalysis.HitCount > 0 {
				skillAnalysis.AvgDamage = float64(skillAnalysis.TotalDamage) / float64(skillAnalysis.HitCount)
			}

			// 计算DPS
			if totalDuration > 0 {
				skillAnalysis.DPS = float64(skillAnalysis.TotalDamage) / float64(totalDuration)
			}

			// 计算伤害占比
			if charStat.DamageDealt > 0 {
				skillAnalysis.DamagePercent = float64(skillAnalysis.TotalDamage) / float64(charStat.DamageDealt) * 100
			}

			// 计算每点能量伤害
			if skillAnalysis.ResourceCost > 0 {
				skillAnalysis.DamagePerResource = float64(skillAnalysis.TotalDamage) / float64(skillAnalysis.ResourceCost)
			}

			// 计算命中率
			if skillAnalysis.UseCount > 0 {
				skillAnalysis.HitRate = float64(skillAnalysis.HitCount) / float64(skillAnalysis.UseCount) * 100
			}

			// 计算暴击率
			if skillAnalysis.HitCount > 0 {
				skillAnalysis.CritRate = float64(skillAnalysis.CritCount) / float64(skillAnalysis.HitCount) * 100
			}

			charAnalysis.SkillBreakdown = append(charAnalysis.SkillBreakdown, skillAnalysis)
		}

		// 角色伤害构成
		charComposition := &models.DamageComposition{
			Physical:    charStat.PhysicalDamage,
			Magic:       charStat.MagicDamage,
			Fire:        charStat.FireDamage,
			Frost:       charStat.FrostDamage,
			Shadow:      charStat.ShadowDamage,
			Holy:        charStat.HolyDamage,
			Nature:      charStat.NatureDamage,
			Dot:         charStat.DotDamage,
			Total:       charStat.DamageDealt,
			Percentages: make(map[string]float64),
		}

		// 计算各类型占比
		if charComposition.Total > 0 {
			charComposition.Percentages["physical"] = float64(charComposition.Physical) / float64(charComposition.Total) * 100
			charComposition.Percentages["magic"] = float64(charComposition.Magic) / float64(charComposition.Total) * 100
			charComposition.Percentages["fire"] = float64(charComposition.Fire) / float64(charComposition.Total) * 100
			charComposition.Percentages["frost"] = float64(charComposition.Frost) / float64(charComposition.Total) * 100
			charComposition.Percentages["shadow"] = float64(charComposition.Shadow) / float64(charComposition.Total) * 100
			charComposition.Percentages["holy"] = float64(charComposition.Holy) / float64(charComposition.Total) * 100
			charComposition.Percentages["nature"] = float64(charComposition.Nature) / float64(charComposition.Total) * 100
			charComposition.Percentages["dot"] = float64(charComposition.Dot) / float64(charComposition.Total) * 100
		}

		charAnalysis.DamageComposition = charComposition

		// 累加队伍伤害构成
		teamComposition.Physical += charStat.PhysicalDamage
		teamComposition.Magic += charStat.MagicDamage
		teamComposition.Fire += charStat.FireDamage
		teamComposition.Frost += charStat.FrostDamage
		teamComposition.Shadow += charStat.ShadowDamage
		teamComposition.Holy += charStat.HolyDamage
		teamComposition.Nature += charStat.NatureDamage
		teamComposition.Dot += charStat.DotDamage
		teamComposition.Total += charStat.DamageDealt

		analysis.Characters = append(analysis.Characters, charAnalysis)
	}

	// 计算队伍伤害构成占比
	teamComposition.Percentages = make(map[string]float64)
	if teamComposition.Total > 0 {
		teamComposition.Percentages["physical"] = float64(teamComposition.Physical) / float64(teamComposition.Total) * 100
		teamComposition.Percentages["magic"] = float64(teamComposition.Magic) / float64(teamComposition.Total) * 100
		teamComposition.Percentages["fire"] = float64(teamComposition.Fire) / float64(teamComposition.Total) * 100
		teamComposition.Percentages["frost"] = float64(teamComposition.Frost) / float64(teamComposition.Total) * 100
		teamComposition.Percentages["shadow"] = float64(teamComposition.Shadow) / float64(teamComposition.Total) * 100
		teamComposition.Percentages["holy"] = float64(teamComposition.Holy) / float64(teamComposition.Total) * 100
		teamComposition.Percentages["nature"] = float64(teamComposition.Nature) / float64(teamComposition.Total) * 100
		teamComposition.Percentages["dot"] = float64(teamComposition.Dot) / float64(teamComposition.Total) * 100
	}

	analysis.TeamDamageComposition = teamComposition

	return analysis, nil
}
