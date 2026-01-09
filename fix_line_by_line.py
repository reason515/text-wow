#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import re

file_path = 'server/internal/test/runner/test_runner.go'

with open(file_path, 'rb') as f:
    content_bytes = f.read()

content = content_bytes.decode('utf-8', errors='replace')
original = content
lines = content.split('\n')

fixed_count = 0

# 逐行检查并修复
for i, line in enumerate(lines):
    original_line = line
    
    # 检查是否包含损坏字符模式
    if '计算最大生命' in line and ('?' in line or '\ufffd' in line):
        # 修复 "计算最大生命?) -> "计算最大生命值")
        line = re.sub(r'"计算最大生命[^"]*?\?"\)', '"计算最大生命值")', line)
        if line != original_line:
            print(f'修复了第{i+1}行: 计算最大生命')
            fixed_count += 1
    
    if '计算生命' in line and '?' in line and '计算最大生命' not in line:
        # 修复 "计算生命?) -> "计算生命值")
        line = re.sub(r'"计算生命[^"]*?\?"\)', '"计算生命值")', line)
        if line != original_line:
            print(f'修复了第{i+1}行: 计算生命')
            fixed_count += 1
    
    if '计算物理暴击' in line and '?' in line:
        line = re.sub(r'"计算物理暴击[^"]*?\?"\)', '"计算物理暴击率")', line)
        if line != original_line:
            print(f'修复了第{i+1}行: 计算物理暴击')
            fixed_count += 1
    
    if '计算法术暴击' in line and '?' in line:
        line = re.sub(r'"计算法术暴击[^"]*?\?"\)', '"计算法术暴击率")', line)
        if line != original_line:
            print(f'修复了第{i+1}行: 计算法术暴击')
            fixed_count += 1
    
    if '计算物理防御' in line and '?' in line:
        line = re.sub(r'"计算物理防御[^"]*?\?"\)', '"计算物理防御力")', line)
        if line != original_line:
            print(f'修复了第{i+1}行: 计算物理防御')
            fixed_count += 1
    
    if '计算魔法防御' in line and '?' in line:
        line = re.sub(r'"计算魔法防御[^"]*?\?"\)', '"计算魔法防御力")', line)
        if line != original_line:
            print(f'修复了第{i+1}行: 计算魔法防御')
            fixed_count += 1
    
    if '计算闪避' in line and '?' in line:
        line = re.sub(r'"计算闪避[^"]*?\?"\)', '"计算闪避率")', line)
        if line != original_line:
            print(f'修复了第{i+1}行: 计算闪避')
            fixed_count += 1
    
    if '次攻' in line and '?' in line:
        line = re.sub(r'"次攻[^"]*?\?"\)', '"次攻击")', line)
        if line != original_line:
            print(f'修复了第{i+1}行: 次攻')
            fixed_count += 1
    
    lines[i] = line

new_content = '\n'.join(lines)

if new_content != original:
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(new_content)
    print(f'修复完成！共修复了 {fixed_count} 处')
else:
    print('没有需要修复的内容')
