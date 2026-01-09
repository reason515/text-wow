#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import re

file_path = 'server/internal/test/runner/test_runner.go'

with open(file_path, 'rb') as f:
    content = f.read().decode('utf-8', errors='replace')

original = content
lines = content.split('\n')

# 修复第333行（索引332）
if len(lines) > 332:
    line = lines[332]
    # 修复损坏字符和括号
    line = re.sub(r'strings\.Contains\(instruction, "角色"\) && strings\.Contains\(instruction, ".*\?"\)\)', 
                  'strings.Contains(instruction, "角色") && strings.Contains(instruction, "在"))', line)
    # 修复括号问题
    line = re.sub(r'\(strings\.Contains\(instruction, "创建"\) && strings\.Contains\(instruction, "个角色"\) && !strings\.Contains\(instruction, "创建一个"\) \|\|',
                  '(strings.Contains(instruction, "创建") && strings.Contains(instruction, "个角色") && !strings.Contains(instruction, "创建一个")) ||', line)
    lines[332] = line

# 修复第342行（索引341）
if len(lines) > 341:
    line = lines[341]
    # 修复损坏字符
    line = re.sub(r'strings\.Contains\(instruction, "怪物"\) && strings\.Contains\(instruction, ".*\?"\)\)',
                  'strings.Contains(instruction, "怪物") && strings.Contains(instruction, "在"))', line)
    lines[341] = line

new_content = '\n'.join(lines)

if new_content != original:
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(new_content)
    print('修复完成')
else:
    print('没有需要修复的内容')
