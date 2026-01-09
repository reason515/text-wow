#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import re

file_path = 'server/internal/test/runner/test_runner.go'

with open(file_path, 'rb') as f:
    content = f.read().decode('utf-8', errors='replace')

original = content

# 修复第333行：修复损坏字符和括号
content = re.sub(
    r'\(strings\.Contains\(instruction, "创建"\) && strings\.Contains\(instruction, "个角色"\) && !strings\.Contains\(instruction, "创建一个"\) \|\| \(strings\.Contains\(instruction, "创建"\) && strings\.Contains\(instruction, "角色"\) && strings\.Contains\(instruction, ".*\?"\)\)',
    '(strings.Contains(instruction, "创建") && strings.Contains(instruction, "个角色") && !strings.Contains(instruction, "创建一个")) || (strings.Contains(instruction, "创建") && strings.Contains(instruction, "角色") && strings.Contains(instruction, "在"))',
    content
)

# 修复第342行：修复损坏字符
content = re.sub(
    r'\(strings\.Contains\(instruction, "创建"\) && strings\.Contains\(instruction, "个怪物"\)\) \|\| \(strings\.Contains\(instruction, "创建"\) && strings\.Contains\(instruction, "怪物"\) && strings\.Contains\(instruction, ".*\?"\)\)',
    '(strings.Contains(instruction, "创建") && strings.Contains(instruction, "个怪物")) || (strings.Contains(instruction, "创建") && strings.Contains(instruction, "怪物") && strings.Contains(instruction, "在"))',
    content
)

# 更通用的修复：直接替换损坏字符模式
content = re.sub(r'strings\.Contains\(instruction, "角色"\) && strings\.Contains\(instruction, "[^\"]*\?"\)\)', 
                 'strings.Contains(instruction, "角色") && strings.Contains(instruction, "在"))', content)
content = re.sub(r'strings\.Contains\(instruction, "怪物"\) && strings\.Contains\(instruction, "[^\"]*\?"\)\)', 
                 'strings.Contains(instruction, "怪物") && strings.Contains(instruction, "在"))', content)

if content != original:
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(content)
    print('修复完成')
    print(f'修复了 {len([c for c in content if ord(c) == 0xFFFD])} 个Unicode替换字符')
else:
    print('没有需要修复的内容')
