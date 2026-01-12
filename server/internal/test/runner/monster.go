package runner



import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"text-wow/internal/models"
)

// Monster 相关函数



func (tr *TestRunner) createMonster(instruction string) error {

	debugPrint("[DEBUG] createMonster: called with instruction: %s\n", instruction)

	// 解析数量（如"创建3个怪物"）
	count := 1

	if strings.Contains(instruction, "个") {

		parts := strings.Split(instruction, "个")

		if len(parts) > 0 {

			countStr := strings.TrimSpace(parts[0])

			// 提取数字

			for i, r := range countStr {

				if r >= '0' && r <= '9' {

					// 找到数字开始位置
					numStr := ""

					for j := i; j < len(countStr); j++ {

						if countStr[j] >= '0' && countStr[j] <= '9' {

							numStr += string(countStr[j])

						} else {

							break

						}

					}

					if c, err := strconv.Atoi(numStr); err == nil {

						count = c

					}

					break

				}

			}

		}

	}



	// 解析防御力（如"防御=10"）
	defense := 5 // 默认

	if strings.Contains(instruction, "防御=") {

		parts := strings.Split(instruction, "防御=")

		if len(parts) > 1 {

			defenseStr := strings.TrimSpace(strings.Split(parts[1], ")")[0])

			defenseStr = strings.TrimSpace(strings.Split(defenseStr, ")")[0])

			defenseStr = strings.TrimSpace(strings.Split(defenseStr, ")")[0])

			if d, err := strconv.Atoi(defenseStr); err == nil {

				defense = d

			}

		}

	}



	// 存储防御力到上下文（用于伤害计算）
	tr.context.Variables["monster_defense"] = defense



	// 创建指定数量的怪物

	for i := 1; i <= count; i++ {

		monster := &models.Monster{

			ID:              fmt.Sprintf("test_monster_%d", i),

			Name:            fmt.Sprintf("测试怪物%d", i),

			Type:            "normal",

			Level:           1,

			HP:              100, // 默认存活

			MaxHP:           100,

			PhysicalAttack:  10,

			MagicAttack:     5,

			PhysicalDefense: defense,

			MagicDefense:    3,

			DodgeRate:       0.05,

		}



		// 解析闪避率（如"闪避=10%"）
		if strings.Contains(instruction, "闪避=") {

			parts := strings.Split(instruction, "闪避=")

			if len(parts) > 1 {

				dodgeStr := strings.TrimSpace(strings.Split(parts[1], "%")[0])

				if dodge, err := strconv.ParseFloat(dodgeStr, 64); err == nil {

					monster.DodgeRate = dodge / 100.0

				}

			}

		}



		// 解析速度（如"速度=80"）
		if strings.Contains(instruction, "速度=") {

			parts := strings.Split(instruction, "速度=")

			if len(parts) > 1 {

				speedStr := strings.TrimSpace(strings.Split(parts[1], ")")[0])

				speedStr = strings.TrimSpace(strings.Split(speedStr, ")")[0])

				speedStr = strings.TrimSpace(strings.Split(speedStr, ")")[0])

				if speed, err := strconv.Atoi(speedStr); err == nil {

					monster.Speed = speed

				}

			}

		}



		// 解析攻击力（如"攻击=20"）
		if strings.Contains(instruction, "攻击=") {

			parts := strings.Split(instruction, "攻击=")

			if len(parts) > 1 {

				attackStr := strings.TrimSpace(strings.Split(parts[1], ")")[0])

				attackStr = strings.TrimSpace(strings.Split(attackStr, ")")[0])

				if attack, err := strconv.Atoi(attackStr); err == nil {

					monster.PhysicalAttack = attack

				}

			}

		}



		// 解析HP（如"HP=100"，"HP=50/100"）
		if strings.Contains(instruction, "HP=") {

			parts := strings.Split(instruction, "HP=")

			if len(parts) > 1 {

				hpStr := strings.TrimSpace(strings.Split(parts[1], ")")[0])

				if strings.Contains(hpStr, "/") {

					// 处理 "50/100" 格式

					hpParts := strings.Split(hpStr, "/")

					if len(hpParts) >= 1 {

						if hp, err := strconv.Atoi(strings.TrimSpace(hpParts[0])); err == nil {

							monster.HP = hp

						}

					}

					if len(hpParts) >= 2 {

						if maxHP, err := strconv.Atoi(strings.TrimSpace(hpParts[1])); err == nil {

							monster.MaxHP = maxHP

						}

					}

				} else {

					// 处理 "100" 格式

					if hp, err := strconv.Atoi(hpStr); err == nil {

						monster.HP = hp

						monster.MaxHP = hp

					}

				}

			}

		}



		// 解析资源（如"资源=100/100"）
		if strings.Contains(instruction, "资源=") {

			parts := strings.Split(instruction, "资源=")

			if len(parts) > 1 {

				resourceStr := strings.TrimSpace(strings.Split(parts[1], ")")[0])

				if strings.Contains(resourceStr, "/") {

					resourceParts := strings.Split(resourceStr, "/")

					if len(resourceParts) >= 1 {

						if resource, err := strconv.Atoi(strings.TrimSpace(resourceParts[0])); err == nil {

							tr.context.Variables["monster.resource"] = resource

						}

					}

				} else {

					if resource, err := strconv.Atoi(resourceStr); err == nil {

						tr.context.Variables["monster.resource"] = resource

					}

				}

			}

		}



		// 解析金币掉落（如"金币掉落=10-20"）
		if strings.Contains(instruction, "金币掉落=") {

			parts := strings.Split(instruction, "金币掉落=")

			if len(parts) > 1 {

				goldStr := strings.TrimSpace(strings.Split(parts[1], ")")[0])

				if strings.Contains(goldStr, "-") {

					// 解析范围，如"10-20"

					goldParts := strings.Split(goldStr, "-")

					if len(goldParts) >= 2 {

						if min, err := strconv.Atoi(strings.TrimSpace(goldParts[0])); err == nil {

							if max, err := strconv.Atoi(strings.TrimSpace(goldParts[1])); err == nil {

								monster.GoldMin = min

								monster.GoldMax = max

								tr.context.Variables["monster_gold_min"] = min

								tr.context.Variables["monster_gold_max"] = max

							}

						}

					}

				} else {

					// 单个值，�10"

					if gold, err := strconv.Atoi(goldStr); err == nil {

						monster.GoldMin = gold

						monster.GoldMax = gold

						tr.context.Variables["monster_gold_min"] = gold

						tr.context.Variables["monster_gold_max"] = gold

					}

				}

			}

		}



		// 存储怪物（monster_1, monster_2, monster_3等）

		// 注意：key用于context存储，monster.ID用于标识

		key := fmt.Sprintf("monster_%d", i)

		if count == 1 {

			key = "monster" // 单个怪物使用monster作为key

		}

		// 确保monster.ID格式正确（monster_1, monster_2等，而不是test_monster_1）
		monster.ID = fmt.Sprintf("monster_%d", i)

		tr.context.Monsters[key] = monster

		debugPrint("[DEBUG] createMonster: stored monster[%s] with PhysicalDefense=%d, HP=%d\n", key, monster.PhysicalDefense, monster.HP)

	}

	debugPrint("[DEBUG] createMonster: total monsters in context: %d\n", len(tr.context.Monsters))



	return nil

}



func (tr *TestRunner) createMultipleMonsters(instruction string) error {

	// 解析怪物列表（通过冒号分隔）
	var monsterDescs []string

	if strings.Contains(instruction, "个") {

		parts := strings.Split(instruction, "个")

		if len(parts) > 1 {

			monsterDescs = strings.Split(parts[1], ",")

		}

	} else if strings.Contains(instruction, ":") {

		parts := strings.Split(instruction, ":")

		if len(parts) > 1 {

			monsterDescs = strings.Split(parts[1], ",")

		}

	}



	for _, monsterDesc := range monsterDescs {

		monsterDesc = strings.TrimSpace(monsterDesc)

		if monsterDesc == "" {

			continue

		}



		// 解析怪物索引（如"怪物1"�怪物2"等）

		monsterIndex := 1

		if strings.Contains(monsterDesc, "怪物") {

			// 提取数字

			re := regexp.MustCompile(`怪物(\d+)`)

			matches := re.FindStringSubmatch(monsterDesc)

			if len(matches) > 1 {

				if idx, err := strconv.Atoi(matches[1]); err == nil {

					monsterIndex = idx

				}

			}

		}



		// 解析速度

		speed := 0

		if strings.Contains(monsterDesc, "速度=") {

			parts := strings.Split(monsterDesc, "速度=")

			if len(parts) > 1 {

				speedStr := strings.TrimSpace(strings.Split(parts[1], ")")[0])

				speedStr = strings.TrimSpace(strings.Split(speedStr, ")")[0])

				speedStr = strings.TrimSpace(strings.Split(speedStr, ")")[0])

				if s, err := strconv.Atoi(speedStr); err == nil {

					speed = s

				}

			}

		}



		// 创建怪物

		monster := &models.Monster{

			ID:              fmt.Sprintf("monster_%d", monsterIndex),

			Name:            fmt.Sprintf("测试怪物%d", monsterIndex),

			Type:            "normal",

			Level:           1,

			HP:              100,

			MaxHP:           100,

			PhysicalAttack:  10,

			MagicAttack:     5,

			PhysicalDefense: 5,

			MagicDefense:    3,

			Speed:           speed,

			DodgeRate:       0.05,

		}



		// 存储怪物（使用monster_1, monster_2等作为key）
		key := fmt.Sprintf("monster_%d", monsterIndex)

		tr.context.Monsters[key] = monster

		debugPrint("[DEBUG] createMultipleMonsters: created monster[%s] with Speed=%d\n", key, speed)

	}



	return nil

}



func (tr *TestRunner) getFirstAliveMonster() *models.Monster {

	// 按key排序，获取第一个存活的怪物

	monsterKeys := []string{}

	for key := range tr.context.Monsters {

		if tr.context.Monsters[key] != nil && tr.context.Monsters[key].HP > 0 {

			monsterKeys = append(monsterKeys, key)

		}

	}



	if len(monsterKeys) == 0 {

		return nil

	}



	sort.Strings(monsterKeys)

	return tr.context.Monsters[monsterKeys[0]]

}



