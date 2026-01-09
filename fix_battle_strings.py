#!/usr/bin/env python3
# -*- coding: utf-8 -*-

file_path = 'server/internal/test/runner/battle.go'

with open(file_path, 'rb') as f:
    content_bytes = f.read()

content = content_bytes.decode('utf-8', errors='replace')

# 查找所有fmt.Sprintf调用，检查是否有未闭合的字符串
import re

# 修复fmt.Sprintf中的字符串问题
# 查找所有fmt.Sprintf("...伤...")的模式
lines = content.split('\n')

fixed_lines = []
for i, line in enumerate(lines):
    # 检查是否有fmt.Sprintf且字符串未正确闭合
    if 'fmt.Sprintf' in line and '伤害' not in line and '伤' in line:
        # 尝试修复
        if '伤' in line or '伤?' in line:
            # 替换为伤害
            line = line.replace('伤', '伤害').replace('伤?', '伤害')
            print(f'Line {i+1}: Fixed fmt.Sprintf string')
    
    # 检查是否有未闭合的字符串（包含fmt.Sprintf但缺少闭合引号）
    if 'fmt.Sprintf' in line and line.count('"') % 2 != 0:
        # 检查下一行是否继续字符串
        if i + 1 < len(lines):
            next_line = lines[i + 1]
            # 如果下一行以引号开始，可能是字符串继续
            if next_line.strip().startswith('"'):
                # 合并这两行
                line = line.rstrip() + next_line.strip()
                lines[i + 1] = ''  # 标记下一行为空
                print(f'Line {i+1}: Merged with next line')
    
    fixed_lines.append(line)

# 重新组合内容
content = '\n'.join(fixed_lines)

# 再次修复fmt.Sprintf中的乱码
content = re.sub(
    r'fmt\.Sprintf\("([^"]*?)伤([^"]*?)"',
    r'fmt.Sprintf("\1伤害\2"',
    content
)

# 修复第91行的注释问题
content = re.sub(
    r'// 收集所有参与者（角色和怪物[^\n]*?\)\s*type participant struct',
    '// 收集所有参与者（角色和怪物）\n\ttype participant struct',
    content
)

# 写入文件
with open(file_path, 'w', encoding='utf-8') as f:
    f.write(content)

print('修复完成')
