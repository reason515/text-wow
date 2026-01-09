#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import re

file_path = 'server/internal/test/runner/test_runner.go'

with open(file_path, 'rb') as f:
    content_bytes = f.read()

content = content_bytes.decode('utf-8', errors='replace')
original = content

# 修复所有剩余的损坏字符
replacements = [
    # 修复 strings.Split 中的损坏字符
    ('strings.Split(hpStr, "�?)[0]', 'strings.Split(hpStr, "=")[0]'),
    ('strings.Split(speedStr, "�?)[0]', 'strings.Split(speedStr, "=")[0]'),
    ('strings.Split(gainStr, "�?)[0]', 'strings.Split(gainStr, "=")[0]'),
    ('strings.Split(regenStr, "�?)[0]', 'strings.Split(regenStr, "=")[0]'),
    # 修复 strings.Contains 和 strings.Split 中的损坏字符（冒号分隔符）
    ('strings.Contains(instruction, "�?)', 'strings.Contains(instruction, ":")'),
    ('strings.Split(instruction, "�?)', 'strings.Split(instruction, ":")'),
    ('strings.Split(parts[1], "�?)', 'strings.Split(parts[1], ",")'),
    ('strings.Split(parts[1], "�?)', 'strings.Split(parts[1], ",")'),
    # 修复 strings.TrimPrefix 和 strings.TrimSuffix 中的损坏字符
    ('strings.TrimPrefix(singleCharInstruction, "�?))', 'strings.TrimPrefix(singleCharInstruction, "(")'),
    ('strings.TrimSuffix(singleCharInstruction, "�?))', 'strings.TrimSuffix(singleCharInstruction, ")")'),
]

fixed_count = 0
for old, new in replacements:
    if old in content:
        count = content.count(old)
        content = content.replace(old, new)
        fixed_count += count
        print(f'修复了 {count} 处: {old[:50]}... -> {new[:50]}...')

# 也使用正则表达式匹配包含Unicode替换字符的模式
patterns = [
    # 修复 strings.Split 中的损坏字符
    (r'strings\.Split\([^,]+,\s*"[^\"]*\ufffd\?"\)\[0\]', lambda m: m.group(0).replace('\ufffd?', '=')),
    # 修复 strings.Contains 中的损坏字符（冒号）
    (r'strings\.Contains\(instruction,\s*"[^\"]*\ufffd\?"\)', lambda m: m.group(0).replace('\ufffd?', ':')),
    # 修复 strings.Split 中的损坏字符（冒号）
    (r'strings\.Split\(instruction,\s*"[^\"]*\ufffd\?"\)', lambda m: m.group(0).replace('\ufffd?', ':')),
    # 修复 strings.Split 中的损坏字符（逗号）
    (r'strings\.Split\(parts\[1\],\s*"[^\"]*\ufffd\?"\)', lambda m: m.group(0).replace('\ufffd?', ',')),
    # 修复 strings.TrimPrefix 和 strings.TrimSuffix 中的损坏字符
    (r'strings\.TrimPrefix\([^,]+,\s*"[^\"]*\ufffd\?"\)', lambda m: m.group(0).replace('\ufffd?', '(')),
    (r'strings\.TrimSuffix\([^,]+,\s*"[^\"]*\ufffd\?"\)', lambda m: m.group(0).replace('\ufffd?', ')')),
]

for pattern, replacement in patterns:
    new_content = re.sub(pattern, replacement, content)
    if new_content != content:
        print(f'修复了（正则）: {pattern[:30]}...')
        content = new_content

if content != original:
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(content)
    print(f'修复完成！共修复了 {fixed_count} 处')
else:
    print('没有需要修复的内容')
