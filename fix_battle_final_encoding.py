#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
修复battle.go中剩余的编码问题
"""

file_path = 'server/internal/test/runner/battle.go'

with open(file_path, 'rb') as f:
    content_bytes = f.read()

# 使用UTF-8解码，将无法解码的字符替换为?
content = content_bytes.decode('utf-8', errors='replace')
lines = content.split('\n')

print(f'总行数: {len(lines)}')

# 查找并修复所有包含乱码字符的行
fixed_count = 0
for i, line in enumerate(lines):
    original_line = line
    
    # 修复 strings.Contains(instruction, "?) 为 strings.Contains(instruction, "在")
    if 'strings.Contains(instruction, "' in line and '\ufffd' in line:
        # 查找乱码字符的位置
        if '"' in line:
            # 尝试修复
            line = line.replace('\ufffd', '在')
            if line != original_line:
                print(f'Line {i+1}: 修复了 strings.Contains')
                fixed_count += 1
                lines[i] = line
    
    # 修复 strings.Split(instruction, "?) 为 strings.Split(instruction, "在")
    if 'strings.Split(instruction, "' in line and '\ufffd' in line:
        line = line.replace('\ufffd', '在')
        if line != original_line:
            print(f'Line {i+1}: 修复了 strings.Split instruction')
            fixed_count += 1
            lines[i] = line
    
    # 修复 strings.Split(parts[1], "?)[0] 为 strings.Split(parts[1], "个")[0]
    if 'strings.Split(parts[1], "' in line and '\ufffd' in line:
        line = line.replace('\ufffd', '个')
        if line != original_line:
            print(f'Line {i+1}: 修复了 strings.Split parts')
            fixed_count += 1
            lines[i] = line
    
    # 修复其他可能的乱码模式
    # 查找包含?的行（可能是乱码字符）
    if '\ufffd' in line and ('strings.Contains' in line or 'strings.Split' in line):
        # 尝试根据上下文修复
        if 'instruction' in line and 'Contains' in line:
            line = line.replace('\ufffd', '在')
            if line != original_line:
                print(f'Line {i+1}: 修复了 Contains 乱码')
                fixed_count += 1
                lines[i] = line
        elif 'parts[1]' in line and 'Split' in line:
            line = line.replace('\ufffd', '个')
            if line != original_line:
                print(f'Line {i+1}: 修复了 Split parts 乱码')
                fixed_count += 1
                lines[i] = line

content = '\n'.join(lines)

# 写入文件
with open(file_path, 'w', encoding='utf-8') as f:
    f.write(content)

print(f'\n修复完成，共修复了 {fixed_count} 处问题')
