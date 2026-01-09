#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import re

file_path = 'server/internal/test/runner/test_runner.go'

with open(file_path, 'rb') as f:
    content_bytes = f.read()

content = content_bytes.decode('utf-8', errors='replace')
original = content

# 修复 strings.Split 中的损坏字符
replacements = [
    # 修复 strings.Split(value, "�?)[0] -> strings.Split(value, "=")[0]
    ('strings.Split(value, "�?)[0]', 'strings.Split(value, "=")[0]'),
    ('strings.Split(parts[1], "�?)[0]', 'strings.Split(parts[1], "=")[0]'),
    ('strings.Split(defenseStr, "�?)[0]', 'strings.Split(defenseStr, "=")[0]'),
    ('strings.Split(goldStr, "�?)[0]', 'strings.Split(goldStr, "=")[0]'),
    ('strings.Split(costStr, "�?)[0]', 'strings.Split(costStr, "=")[0]'),
    ('strings.Split(multiplierStr, "�?)[0]', 'strings.Split(multiplierStr, "=")[0]'),
    ('strings.Split(cooldownStr, "�?)[0]', 'strings.Split(cooldownStr, "=")[0]'),
    # 修复 strings.Index(value, "�?) -> strings.Index(value, "=")
    ('strings.Index(value, "�?)', 'strings.Index(value, "=")'),
    # 修复 strings.Split(instruction, "消�?) -> strings.Split(instruction, "消耗") 或 strings.Split(instruction, "消耗=")
    ('strings.Split(instruction, "消�?)', 'strings.Split(instruction, "消耗=")'),
    ('strings.Contains(instruction, "消�?)', 'strings.Contains(instruction, "消耗")'),
    # 修复 strings.Split(instruction, "治疗�?) -> strings.Split(instruction, "治疗=")
    ('strings.Split(instruction, "治疗�?)', 'strings.Split(instruction, "治疗=")'),
    # 修复 strings.Split(instruction, "攻击�?) -> strings.Split(instruction, "攻击=")
    ('strings.Split(instruction, "攻击�?)', 'strings.Split(instruction, "攻击=")'),
    ('strings.Contains(instruction, "攻击�?)', 'strings.Contains(instruction, "攻击")'),
    # 修复 strings.Contains(instruction, "效果�?) -> strings.Contains(instruction, "效果")
    ('strings.Contains(instruction, "效果�?)', 'strings.Contains(instruction, "效果")'),
    # 修复 strings.Split(instruction, "防御�?) -> strings.Split(instruction, "防御=")
    ('strings.Split(instruction, "防御�?)', 'strings.Split(instruction, "防御=")'),
    # 修复 strings.Split(instruction, "物理暴击�?) -> strings.Split(instruction, "物理暴击率=")
    ('strings.Split(instruction, "物理暴击�?)', 'strings.Split(instruction, "物理暴击率=")'),
    # 修复其他损坏字符
    ('strings.Split(multiplierStr, "�?)', 'strings.Split(multiplierStr, "=")'),
    ('strings.Contains(multiplierStr, "�?)', 'strings.Contains(multiplierStr, "=")'),
    ('strings.Split(cooldownStr, "�?)', 'strings.Split(cooldownStr, "=")'),
    ('strings.Contains(cooldownStr, "�?)', 'strings.Contains(cooldownStr, "=")'),
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
    (r'strings\.Split\([^,]+,\s*"[^\"]*\ufffd\?"\)', lambda m: m.group(0).replace('\ufffd?', '=')),
    (r'strings\.Index\([^,]+,\s*"[^\"]*\ufffd\?"\)', lambda m: m.group(0).replace('\ufffd?', '=')),
    (r'strings\.Contains\([^,]+,\s*"[^\"]*\ufffd\?"\)', lambda m: m.group(0).replace('\ufffd?', '') if '效果' in m.group(0) or '攻击' in m.group(0) or '消耗' in m.group(0) else m.group(0).replace('\ufffd?', '=')),
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
