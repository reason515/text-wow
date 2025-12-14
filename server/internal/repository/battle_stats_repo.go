package repository

import (
	"database/sql"
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
		       COALESCE(s.name, bsb.skill_id) as skill_name
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
