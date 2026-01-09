#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import re

file_path = 'server/internal/test/runner/test_runner.go'

with open(file_path, 'rb') as f:
    content_bytes = f.read()

# 先尝试用UTF-8解码，如果有错误就用replace
try:
    content = content_bytes.decode('utf-8')
except:
    content = content_bytes.decode('utf-8', errors='replace')

original = content
lines = content.split('\n')

# 修复第333行（索引332）
if len(lines) > 332:
    line = lines[332]
    original_line = line
    # 修复：添加缺失的右括号，并修复损坏字符
    # 从: (strings.Contains(instruction, "创建") && strings.Contains(instruction, "个角色") && !strings.Contains(instruction, "创建一个") || (strings.Contains(instruction, "创建") && strings.Contains(instruction, "角色") && strings.Contains(instruction, "?))
    # 到: (strings.Contains(instruction, "创建") && strings.Contains(instruction, "个角色") && !strings.Contains(instruction, "创建一个")) || (strings.Contains(instruction, "创建") && strings.Contains(instruction, "角色") && strings.Contains(instruction, "在"))
    
    # 修复括号：在 "创建一个" 后面添加右括号
    line = re.sub(r'\(strings\.Contains\(instruction, "创建"\) && strings\.Contains\(instruction, "个角色"\) && !strings\.Contains\(instruction, "创建一个"\) \|\|',
                  '(strings.Contains(instruction, "创建") && strings.Contains(instruction, "个角色") && !strings.Contains(instruction, "创建一个")) ||', line)
    
    # 修复损坏字符：替换 "?)) 为 "在"))
    # 尝试多种匹配模式
    line = re.sub(r'strings\.Contains\(instruction, "角色"\) && strings\.Contains\(instruction, "[^"]*\?"\)\)',
                  'strings.Contains(instruction, "角色") && strings.Contains(instruction, "在"))', line)
    # 直接替换损坏字符模式（包括Unicode替换字符）
    if '\ufffd' in line or '' in line:
        # 替换包含Unicode替换字符的模式
        line = re.sub(r'strings\.Contains\(instruction, "角色"\) && strings\.Contains\(instruction, "[^"]*[\ufffd]\?"\)\)',
                      'strings.Contains(instruction, "角色") && strings.Contains(instruction, "在"))', line)
        # 直接替换字符串中的损坏字符
        line = line.replace('"?))', '"在"))')
        line = line.replace('"?))', '"在"))')
    
    if line != original_line:
        lines[332] = line
        print(f'修复了第333行')

# 修复第342行（索引341）
if len(lines) > 341:
    line = lines[341]
    original_line = line
    # 修复损坏字符：替换 "?)) 为 "在"))
    line = re.sub(r'strings\.Contains\(instruction, "怪物"\) && strings\.Contains\(instruction, "[^"]*\?"\)\)',
                  'strings.Contains(instruction, "怪物") && strings.Contains(instruction, "在"))', line)
    # 直接替换损坏字符模式（包括Unicode替换字符）
    if '\ufffd' in line or '' in line:
        # 替换包含Unicode替换字符的模式
        line = re.sub(r'strings\.Contains\(instruction, "怪物"\) && strings\.Contains\(instruction, "[^"]*[\ufffd]\?"\)\)',
                      'strings.Contains(instruction, "怪物") && strings.Contains(instruction, "在"))', line)
        # 直接替换字符串中的损坏字符
        line = line.replace('"?))', '"在"))')
        line = line.replace('"?))', '"在"))')
    
    if line != original_line:
        lines[341] = line
        print(f'修复了第342行')

new_content = '\n'.join(lines)

if new_content != original:
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(new_content)
    print('修复完成！')
else:
    print('没有需要修复的内容')
