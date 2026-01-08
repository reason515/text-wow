#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
怪物强度快速调整工具
用法：
  python adjust_monster_strength.py --list                    # 列出所有配置
  python adjust_monster_strength.py --set 1 10 1.5 1.4 1.4 0.02  # 设置配置
  python adjust_monster_strength.py --apply                  # 应用配置到怪物数据
  python adjust_monster_strength.py --reset                  # 重置所有配置为默认值
"""

import sqlite3
import argparse
import sys
from pathlib import Path

DB_PATH = Path("server/game.db")

# 默认配置（当前已应用的强度提升）
DEFAULT_CONFIGS = [
    (1, 10, 1.5, 1.4, 1.4, 0.02, "1-10级：HP +50%, 攻击 +40%, 防御 +40%, 暴击率 +2%"),
    (11, 20, 1.45, 1.35, 1.35, 0.02, "10-20级：HP +45%, 攻击 +35%, 防御 +35%, 暴击率 +2%"),
    (21, 30, 1.4, 1.35, 1.35, 0.02, "20-30级：HP +40%, 攻击 +35%, 防御 +35%, 暴击率 +2%"),
    (31, 40, 1.4, 1.35, 1.35, 0.03, "30-40级：HP +40%, 攻击 +35%, 防御 +35%, 暴击率 +3%"),
    (41, 50, 1.35, 1.3, 1.3, 0.03, "40-50级：HP +35%, 攻击 +30%, 防御 +30%, 暴击率 +3%"),
    (51, 60, 1.35, 1.3, 1.3, 0.03, "50-60级：HP +35%, 攻击 +30%, 防御 +30%, 暴击率 +3%"),
]


def ensure_config_table(db):
    """确保配置表存在"""
    db.execute("""
        CREATE TABLE IF NOT EXISTS monster_strength_config (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            level_min INTEGER NOT NULL,
            level_max INTEGER NOT NULL,
            hp_multiplier REAL DEFAULT 1.0,
            attack_multiplier REAL DEFAULT 1.0,
            defense_multiplier REAL DEFAULT 1.0,
            crit_rate_bonus REAL DEFAULT 0.0,
            description TEXT,
            is_active INTEGER DEFAULT 1,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            UNIQUE(level_min, level_max)
        )
    """)
    db.execute("CREATE INDEX IF NOT EXISTS idx_strength_config_level ON monster_strength_config(level_min, level_max)")
    db.execute("CREATE INDEX IF NOT EXISTS idx_strength_config_active ON monster_strength_config(is_active)")


def list_configs(db):
    """列出所有配置"""
    cursor = db.execute("""
        SELECT level_min, level_max, hp_multiplier, attack_multiplier, 
               defense_multiplier, crit_rate_bonus, description, is_active
        FROM monster_strength_config
        ORDER BY level_min
    """)
    
    configs = cursor.fetchall()
    if not configs:
        print("暂无配置")
        return
    
    print("当前强度配置：")
    print("=" * 60)
    for level_min, level_max, hp_mult, attack_mult, defense_mult, crit_bonus, desc, is_active in configs:
        status = "启用" if is_active else "禁用"
        print(f"等级 {level_min:2d}-{level_max:2d}: "
              f"HP×{hp_mult:.2f}, 攻击×{attack_mult:.2f}, 防御×{defense_mult:.2f}, "
              f"暴击+{crit_bonus*100:.2f}% [{status}]")
        if desc:
            print(f"          {desc}")
    print("=" * 60)


def set_config(db, level_min, level_max, hp_mult, attack_mult, defense_mult, crit_bonus, description=None):
    """设置配置"""
    if description is None:
        description = f"等级 {level_min}-{level_max}：HP ×{hp_mult:.2f}, 攻击 ×{attack_mult:.2f}, 防御 ×{defense_mult:.2f}, 暴击率 +{crit_bonus*100:.2f}%"
    
    db.execute("""
        INSERT OR REPLACE INTO monster_strength_config 
        (level_min, level_max, hp_multiplier, attack_multiplier, defense_multiplier, 
         crit_rate_bonus, description, updated_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
    """, (level_min, level_max, hp_mult, attack_mult, defense_mult, crit_bonus, description))
    
    print(f"✅ 配置已更新：等级 {level_min}-{level_max}")
    print(f"   HP倍数: {hp_mult:.2f}, 攻击倍数: {attack_mult:.2f}, 防御倍数: {defense_mult:.2f}, 暴击率加成: {crit_bonus*100:.2f}%")


def apply_config(db):
    """应用配置到怪物数据"""
    # 获取所有启用的配置
    cursor = db.execute("""
        SELECT level_min, level_max, hp_multiplier, attack_multiplier, 
               defense_multiplier, crit_rate_bonus
        FROM monster_strength_config
        WHERE is_active = 1
        ORDER BY level_min
    """)
    
    configs = cursor.fetchall()
    if not configs:
        print("⚠️  没有启用的配置")
        return
    
    print("正在应用配置到怪物数据...")
    updated_count = 0
    
    for level_min, level_max, hp_mult, attack_mult, defense_mult, crit_bonus in configs:
        cursor = db.execute("""
            UPDATE monsters
            SET 
                hp = ROUND(hp * ?),
                physical_attack = ROUND(physical_attack * ?),
                magic_attack = ROUND(magic_attack * ?),
                physical_defense = ROUND(physical_defense * ?),
                magic_defense = ROUND(magic_defense * ?),
                phys_crit_rate = MIN(0.4, phys_crit_rate + ?),
                spell_crit_rate = MIN(0.4, spell_crit_rate + ?)
            WHERE level >= ? AND level <= ?
        """, (hp_mult, attack_mult, attack_mult, defense_mult, defense_mult, 
              crit_bonus, crit_bonus, level_min, level_max))
        
        count = cursor.rowcount
        updated_count += count
        print(f"  等级 {level_min}-{level_max}: 更新了 {count} 个怪物")
    
    print(f"✅ 配置已应用到 {updated_count} 个怪物")


def reset_config(db):
    """重置为默认配置"""
    print("正在重置为默认配置...")
    db.execute("DELETE FROM monster_strength_config")
    
    for level_min, level_max, hp_mult, attack_mult, defense_mult, crit_bonus, desc in DEFAULT_CONFIGS:
        db.execute("""
            INSERT INTO monster_strength_config 
            (level_min, level_max, hp_multiplier, attack_multiplier, defense_multiplier, 
             crit_rate_bonus, description)
            VALUES (?, ?, ?, ?, ?, ?, ?)
        """, (level_min, level_max, hp_mult, attack_mult, defense_mult, crit_bonus, desc))
    
    print("✅ 已重置为默认配置")


def main():
    parser = argparse.ArgumentParser(description="怪物强度快速调整工具")
    parser.add_argument("--db", default="server/game.db", help="数据库路径")
    parser.add_argument("--list", action="store_true", help="列出所有配置")
    parser.add_argument("--set", nargs=6, metavar=("MIN", "MAX", "HP", "ATTACK", "DEFENSE", "CRIT"),
                       help="设置配置: 等级下限 等级上限 HP倍数 攻击倍数 防御倍数 暴击率加成")
    parser.add_argument("--apply", action="store_true", help="应用配置到怪物数据")
    parser.add_argument("--reset", action="store_true", help="重置为默认配置")
    
    args = parser.parse_args()
    
    if not Path(args.db).exists():
        print(f"❌ 数据库文件不存在: {args.db}")
        sys.exit(1)
    
    conn = sqlite3.connect(args.db)
    conn.execute("PRAGMA foreign_keys = ON")
    
    try:
        ensure_config_table(conn)
        
        if args.list:
            list_configs(conn)
        elif args.set:
            level_min, level_max = int(args.set[0]), int(args.set[1])
            hp_mult, attack_mult, defense_mult = float(args.set[2]), float(args.set[3]), float(args.set[4])
            crit_bonus = float(args.set[5])
            set_config(conn, level_min, level_max, hp_mult, attack_mult, defense_mult, crit_bonus)
            conn.commit()
        elif args.apply:
            apply_config(conn)
            conn.commit()
        elif args.reset:
            reset_config(conn)
            conn.commit()
        else:
            parser.print_help()
    finally:
        conn.close()


if __name__ == "__main__":
    main()














































































