#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
拆分test_runner.go文件为多个模块
"""

import re
import os

file_path = 'server/internal/test/runner/test_runner.go'

# 读取文件
with open(file_path, 'r', encoding='utf-8', errors='replace') as f:
    content = f.read()

lines = content.split('\n')

# 找到package声明和imports
package_line = None
import_start = None
import_end = None

for i, line in enumerate(lines):
    if line.startswith('package '):
        package_line = i
    elif line.strip() == 'import (':
        import_start = i
    elif import_start is not None and line.strip() == ')' and import_end is None:
        import_end = i
        break

# 提取package和imports
header_lines = []
if package_line is not None:
    header_lines.append(lines[package_line])
    header_lines.append('')
if import_start is not None and import_end is not None:
    header_lines.extend(lines[import_start:import_end+1])
    header_lines.append('')

header = '\n'.join(header_lines)

# 查找所有类型定义
type_definitions = []
type_pattern = re.compile(r'^type\s+(\w+)\s+')
type_start = None
type_name = None

for i, line in enumerate(lines):
    match = type_pattern.match(line)
    if match:
        if type_start is not None:
            # 保存前一个类型定义
            type_definitions.append((type_name, type_start, i))
        type_name = match.group(1)
        type_start = i
    elif type_start is not None and line.strip() == '}' and not line.strip().startswith('//'):
        # 类型定义结束
        type_definitions.append((type_name, type_start, i+1))
        type_start = None
        type_name = None

# 查找所有函数定义
function_definitions = []
func_pattern = re.compile(r'^func\s+(\(tr\s+\*TestRunner\)\s+)?(\w+)')
func_start = None
func_name = None
brace_count = 0

for i, line in enumerate(lines):
    match = func_pattern.match(line)
    if match:
        if func_start is not None:
            # 保存前一个函数定义
            function_definitions.append((func_name, func_start, i))
        func_name = match.group(2)
        func_start = i
        brace_count = line.count('{') - line.count('}')
    elif func_start is not None:
        brace_count += line.count('{') - line.count('}')
        if brace_count == 0 and line.strip() == '}':
            # 函数定义结束
            function_definitions.append((func_name, func_start, i+1))
            func_start = None
            func_name = None
            brace_count = 0

print(f'找到 {len(type_definitions)} 个类型定义')
print(f'找到 {len(function_definitions)} 个函数定义')

# 函数分类
character_funcs = ['createCharacter', 'createMultipleCharacters', 'createTestCharacter', 'executeGetCharacterData']
monster_funcs = ['createMonster', 'createMultipleMonsters', 'getFirstAliveMonster']
team_funcs = ['createTeam', 'executeCreateEmptyTeam', 'executeCreateTeamWithMembers', 
              'executeAddCharacterToTeamSlot', 'executeTryAddCharacterToTeamSlot',
              'executeUnlockTeamSlot', 'executeTryAddCharacterToUnlockedSlot']
equipment_funcs = ['generateMultipleEquipments', 'generateEquipmentFromMonster', 
                   'generateEquipmentWithAttributes', 'executeEquipItem', 'executeUnequipItem',
                   'executeEquipAllItems', 'executeTryEquipItem']
calculation_funcs = ['executeCalculatePhysCritDamage', 'executeCalculateSpellCritDamage',
                     'executeCalculateSpeed', 'executeCalculateResourceRegen',
                     'executeCalculateBaseDamage', 'executeCalculateDefenseReduction',
                     'executeApplyCrit', 'executeCalculatePhysicalAttack', 'executeCalculateMagicAttack',
                     'executeCalculateMaxHP', 'executeCalculatePhysCritRate', 'executeCalculateSpellCritRate',
                     'executeCalculatePhysicalDefense', 'executeCalculateMagicDefense', 'executeCalculateDodgeRate']
battle_funcs = ['executeBattleRound', 'executeAttackMonster', 'executeMonsterAttack',
                'executeRemainingMonstersAttack', 'executeBuildTurnOrder', 'buildTurnOrder',
                'setBattleResult', 'executeWaitRestRecovery', 'executeMultipleAttacks',
                'executeContinueBattleUntil']
instruction_funcs = ['executeInstruction', 'executeSetup', 'executeStep', 'executeTeardown',
                    'executeSetVariable', 'executeGainGold', 'executeGainExploration',
                    'executeCreateItem', 'executePurchaseItem', 'executeTryPurchaseItem',
                    'executeViewShopItems']
context_funcs = ['updateAssertionContext', 'safeSetContext', 'safeSetVariable']
core_funcs = ['RunTestSuite', 'RunTestCase', 'RunAllTests', 'NewTestRunner']

# 创建types.go
print('\n创建types.go...')
types_content = header
types_content += '\n// 类型定义\n\n'

for type_name, start, end in type_definitions:
    types_content += '\n'.join(lines[start:end]) + '\n\n'

with open('server/internal/test/runner/types.go', 'w', encoding='utf-8') as f:
    f.write(types_content)

print('types.go 创建完成')

# 创建各个功能模块文件
modules = {
    'character.go': character_funcs,
    'monster.go': monster_funcs,
    'team.go': team_funcs,
    'equipment.go': equipment_funcs,
    'calculation.go': calculation_funcs,
    'battle.go': battle_funcs,
    'instruction.go': instruction_funcs,
    'context.go': context_funcs,
}

for filename, func_names in modules.items():
    print(f'\n创建{filename}...')
    module_content = header
    module_content += f'// {filename.replace(".go", "").title()} 相关函数\n\n'
    
    found_funcs = []
    for func_name, start, end in function_definitions:
        if func_name in func_names:
            found_funcs.append((func_name, start, end))
            module_content += '\n'.join(lines[start:end]) + '\n\n'
    
    if found_funcs:
        with open(f'server/internal/test/runner/{filename}', 'w', encoding='utf-8') as f:
            f.write(module_content)
        print(f'{filename} 创建完成，包含 {len(found_funcs)} 个函数')
    else:
        print(f'{filename} 没有找到匹配的函数')

# 创建新的test_runner.go（只包含核心函数）
print('\n创建新的test_runner.go...')
core_content = header
core_content += '// TestRunner 核心结构和主要运行逻辑\n\n'

# 添加TestRunner结构体定义
for type_name, start, end in type_definitions:
    if type_name == 'TestRunner':
        core_content += '\n'.join(lines[start:end]) + '\n\n'
        break

# 添加NewTestRunner
for func_name, start, end in function_definitions:
    if func_name == 'NewTestRunner':
        core_content += '\n'.join(lines[start:end]) + '\n\n'
        break

# 添加核心运行函数
for func_name, start, end in function_definitions:
    if func_name in core_funcs:
        core_content += '\n'.join(lines[start:end]) + '\n\n'

# 备份原文件
os.rename(file_path, file_path + '.backup_split')

# 写入新的test_runner.go
with open(file_path, 'w', encoding='utf-8') as f:
    f.write(core_content)

print('新的test_runner.go 创建完成')
print(f'\n原文件已备份为: {file_path}.backup_split')
