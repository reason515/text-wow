#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import re

file_path = 'server/internal/test/runner/test_runner.go'

with open(file_path, 'rb') as f:
    content = f.read().decode('utf-8', errors='replace')

original = content

# 修复字符串中的损坏字符（在条件表达式中，带引号和括号结束）
content = re.sub(r'strings\.Contains\(instruction, "角色"\) && strings\.Contains\(instruction, "\ufffd\?"\)\)', 'strings.Contains(instruction, "角色") && strings.Contains(instruction, "在"))', content)
content = re.sub(r'strings\.Contains\(instruction, "怪物"\) && strings\.Contains\(instruction, "\ufffd\?"\)\)', 'strings.Contains(instruction, "怪物") && strings.Contains(instruction, "在"))', content)

if content != original:
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(content)
    print('修复完成')
else:
    print('没有需要修复的内容')
