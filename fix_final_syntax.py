#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import re

file_path = 'server/internal/test/runner/test_runner.go'

with open(file_path, 'rb') as f:
    content_bytes = f.read()

content = content_bytes.decode('utf-8', errors='replace')
original = content

# 修复未闭合的括号
fixes = [
    # 修复 strings.TrimSuffix 中未闭合的括号
    (r'strings\.TrimSuffix\([^,]+,\s*"\)"\s*\)', lambda m: m.group(0).replace('")"', '")")')),
    # 修复 strings.TrimSuffix 中缺少闭合括号的情况
    (r'strings\.TrimSuffix\([^,]+,\s*"\)"\s*$', lambda m: m.group(0) + ')'),
    # 修复 strings.TrimSuffix(singleCharInstruction, ")") 缺少闭合括号
    (r'strings\.TrimSuffix\(singleCharInstruction,\s*"\)"\s*$', lambda m: m.group(0) + ')'),
]

fixed_count = 0
for pattern, replacement in fixes:
    matches = list(re.finditer(pattern, content, re.MULTILINE))
    if matches:
        # 从后往前替换，避免位置偏移
        for match in reversed(matches):
            old = match.group(0)
            new = replacement(match)
            if old != new:
                content = content[:match.start()] + new + content[match.end():]
                fixed_count += 1
                print(f'修复了: {old[:50]}... -> {new[:50]}...')

# 直接替换已知的问题
direct_fixes = [
    ('strings.TrimSuffix(singleCharInstruction, ")")', 'strings.TrimSuffix(singleCharInstruction, ")"))'),
]

for old, new in direct_fixes:
    if old in content and new not in content:
        count = content.count(old)
        content = content.replace(old, new)
        fixed_count += count
        print(f'修复了 {count} 处: {old[:50]}... -> {new[:50]}...')

if content != original:
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(content)
    print(f'修复完成！共修复了 {fixed_count} 处')
else:
    print('没有需要修复的内容')
