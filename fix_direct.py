#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import re

file_path = 'server/internal/test/runner/test_runner.go'

with open(file_path, 'rb') as f:
    content = f.read().decode('utf-8', errors='replace')

original = content

# 修复第333行和342行的损坏字符
# 修复模式：strings.Contains(instruction, "角色") && strings.Contains(instruction, "?))
content = re.sub(
    r'strings\.Contains\(instruction, "角色"\) && strings\.Contains\(instruction, "\ufffd\?"\)\)',
    'strings.Contains(instruction, "角色") && strings.Contains(instruction, "在"))',
    content
)

content = re.sub(
    r'strings\.Contains\(instruction, "怪物"\) && strings\.Contains\(instruction, "\ufffd\?"\)\)',
    'strings.Contains(instruction, "怪物") && strings.Contains(instruction, "在"))',
    content
)

# 修复括号问题：在 "创建一个" 后面添加右括号
content = re.sub(
    r'\(strings\.Contains\(instruction, "创建"\) && strings\.Contains\(instruction, "个角色"\) && !strings\.Contains\(instruction, "创建一个"\) \|\|',
    '(strings.Contains(instruction, "创建") && strings.Contains(instruction, "个角色") && !strings.Contains(instruction, "创建一个")) ||',
    content
)

if content != original:
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(content)
    print('修复完成')
else:
    print('没有需要修复的内容')
