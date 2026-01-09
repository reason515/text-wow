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
    # 直接查找并替换损坏字符模式
    # 查找: strings.Contains(instruction, "角色") && strings.Contains(instruction, "?))
    # 替换为: strings.Contains(instruction, "角色") && strings.Contains(instruction, "在"))
    
    # 使用更精确的匹配
    if 'strings.Contains(instruction, "角色") && strings.Contains(instruction, "' in line:
        # 查找从 "角色" 到 "?)) 的部分并替换
        line = re.sub(
            r'strings\.Contains\(instruction, "角色"\) && strings\.Contains\(instruction, "[^"]*\?"\)\)',
            'strings.Contains(instruction, "角色") && strings.Contains(instruction, "在"))',
            line
        )
        # 如果正则没有匹配，尝试直接字符串替换
        if '"?))' in line or '\ufffd?))' in line:
            # 找到损坏字符的位置并替换
            import re
            line = re.sub(r'("角色"\) && strings\.Contains\(instruction, ")[^\"]*\?\)\)', r'\1在"))', line)
    
    if line != original_line:
        lines[332] = line
        print(f'修复了第333行')

# 修复第342行（索引341）
if len(lines) > 341:
    line = lines[341]
    original_line = line
    # 直接查找并替换损坏字符模式
    if 'strings.Contains(instruction, "怪物") && strings.Contains(instruction, "' in line:
        # 查找从 "怪物" 到 "?)) 的部分并替换
        line = re.sub(
            r'strings\.Contains\(instruction, "怪物"\) && strings\.Contains\(instruction, "[^"]*\?"\)\)',
            'strings.Contains(instruction, "怪物") && strings.Contains(instruction, "在"))',
            line
        )
        # 如果正则没有匹配，尝试直接字符串替换
        if '"?))' in line or '\ufffd?))' in line:
            # 找到损坏字符的位置并替换
            line = re.sub(r'("怪物"\) && strings\.Contains\(instruction, ")[^\"]*\?\)\)', r'\1在"))', line)
    
    if line != original_line:
        lines[341] = line
        print(f'修复了第342行')

# 修复所有包含损坏字符的 strings.Contains 行
for i, line in enumerate(lines):
    original_line = line
    # 修复 "计算最大生命?") -> "计算最大生命值")
    line = re.sub(r'strings\.Contains\(instruction, "计算最大生命[^\"]*\?"\)', 
                  'strings.Contains(instruction, "计算最大生命值")', line)
    # 修复 "计算生命?") -> "计算生命值")
    line = re.sub(r'strings\.Contains\(instruction, "计算生命[^\"]*\?"\)', 
                  'strings.Contains(instruction, "计算生命值")', line)
    # 修复 "计算物理暴击?") -> "计算物理暴击率")
    line = re.sub(r'strings\.Contains\(instruction, "计算物理暴击[^\"]*\?"\)', 
                  'strings.Contains(instruction, "计算物理暴击率")', line)
    # 修复 "计算法术暴击?") -> "计算法术暴击率")
    line = re.sub(r'strings\.Contains\(instruction, "计算法术暴击[^\"]*\?"\)', 
                  'strings.Contains(instruction, "计算法术暴击率")', line)
    # 修复 "计算物理防御?") -> "计算物理防御力")
    line = re.sub(r'strings\.Contains\(instruction, "计算物理防御[^\"]*\?"\)', 
                  'strings.Contains(instruction, "计算物理防御力")', line)
    # 修复 "计算魔法防御?") -> "计算魔法防御力")
    line = re.sub(r'strings\.Contains\(instruction, "计算魔法防御[^\"]*\?"\)', 
                  'strings.Contains(instruction, "计算魔法防御力")', line)
    # 修复 "计算闪避?") -> "计算闪避率")
    line = re.sub(r'strings\.Contains\(instruction, "计算闪避[^\"]*\?"\)', 
                  'strings.Contains(instruction, "计算闪避率")', line)
    # 修复 "次攻?") -> "次攻击")
    line = re.sub(r'strings\.Contains\(instruction, "次攻[^\"]*\?"\)', 
                  'strings.Contains(instruction, "次攻击")', line)
    
    if line != original_line:
        lines[i] = line
        print(f'修复了第{i+1}行')

new_content = '\n'.join(lines)

if new_content != original:
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(new_content)
    print('修复完成！')
else:
    print('没有需要修复的内容')
