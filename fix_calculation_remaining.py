#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import re

file_path = 'server/internal/test/runner/calculation.go'
with open(file_path, 'r', encoding='utf-8', errors='replace') as f:
    content = f.read()

# 修复所有注释和代码混在一起的情况
fixes = [
    (r'// 默认基础获得[^\n]*\tif baseGain == 0', '// 默认基础获得值\n\t\tif baseGain == 0'),
    (r'// 默认基础获得[^\n]*\tregen := tr.calculator.CalculateRageGain', '// 默认基础获得值\n\t\t\tregen := tr.calculator.CalculateRageGain'),
    (r'// 使用基础获得值和加成百分[^\n]*\tregen := tr.calculator.CalculateRageGain', '// 使用基础获得值和加成百分比\n\t\t\tregen := tr.calculator.CalculateRageGain'),
    (r'// 从Variables获取基础获得值和加成百分[^\n]*\trageBaseGain := 10', '// 从Variables获取基础获得值和加成百分比\n\t\t\trageBaseGain := 10'),
    (r'// 获取基础伤害（如果已计算[^\n]*\tbaseDamage := char.PhysicalAttack', '// 获取基础伤害（如果已计算）\n\tbaseDamage := char.PhysicalAttack'),
]

for pattern, replacement in fixes:
    content = re.sub(pattern, replacement, content)

with open(file_path, 'w', encoding='utf-8', errors='replace') as f:
    f.write(content)

print('修复完成')
