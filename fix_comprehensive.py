#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import re

file_path = 'server/internal/test/runner/test_runner.go'

with open(file_path, 'rb') as f:
    content_bytes = f.read()

content = content_bytes.decode('utf-8', errors='replace')
original = content

# 查找所有包含损坏字符的行
lines = content.split('\n')
fixed_lines = []

for i, line in enumerate(lines):
    original_line = line
    
    # 修复 "计算最大生命?) -> "计算最大生命值")
    # 检查是否包含损坏字符（Unicode替换字符或?）
    if '计算最大生命' in line and ('?' in line or '\ufffd' in line):
        # 使用更宽松的匹配，匹配任何字符直到?)
        line = re.sub(r'"计算最大生命.*?\?"\)', '"计算最大生命值")', line, flags=re.DOTALL)
        # 也尝试直接替换（包括Unicode替换字符）
        line = line.replace('"计算最大生命?)', '"计算最大生命值")')
        # 尝试匹配包含Unicode替换字符的模式
        if '\ufffd' in line:
            line = re.sub(r'"计算最大生命[^\"]*\ufffd\?"\)', '"计算最大生命值")', line)
    
    # 修复 "计算生命?) -> "计算生命值")
    if '计算生命' in line and '?' in line and '计算最大生命' not in line:
        # 使用更宽松的匹配，匹配任何字符直到?)
        line = re.sub(r'"计算生命.*?\?"\)', '"计算生命值")', line, flags=re.DOTALL)
        # 也尝试直接替换
        line = line.replace('"计算生命?)', '"计算生命值")')
    
    # 修复 "计算物理暴击?) -> "计算物理暴击率")
    if '计算物理暴击' in line and '?' in line:
        line = re.sub(r'"计算物理暴击[^"]*?\?"\)', '"计算物理暴击率")', line)
    
    # 修复 "计算法术暴击?) -> "计算法术暴击率")
    if '计算法术暴击' in line and '?' in line:
        line = re.sub(r'"计算法术暴击[^"]*?\?"\)', '"计算法术暴击率")', line)
    
    # 修复 "计算物理防御?) -> "计算物理防御力")
    if '计算物理防御' in line and '?' in line:
        line = re.sub(r'"计算物理防御[^"]*?\?"\)', '"计算物理防御力")', line)
    
    # 修复 "计算魔法防御?) -> "计算魔法防御力")
    if '计算魔法防御' in line and '?' in line:
        line = re.sub(r'"计算魔法防御[^"]*?\?"\)', '"计算魔法防御力")', line)
    
    # 修复 "计算闪避?) -> "计算闪避率")
    if '计算闪避' in line and '?' in line:
        line = re.sub(r'"计算闪避[^"]*?\?"\)', '"计算闪避率")', line)
    
    # 修复 "次攻?) -> "次攻击")
    if '次攻' in line and '?' in line:
        line = re.sub(r'"次攻[^"]*?\?"\)', '"次攻击")', line)
    
    if line != original_line:
        fixed_lines.append(i + 1)
        print(f'修复了第{i+1}行')
    
    lines[i] = line

new_content = '\n'.join(lines)

if new_content != original:
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(new_content)
    print(f'修复完成！共修复了 {len(fixed_lines)} 行')
else:
    print('没有需要修复的内容')
