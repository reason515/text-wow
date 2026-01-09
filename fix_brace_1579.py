#!/usr/bin/env python3
# -*- coding: utf-8 -*-

file_path = 'server/internal/test/runner/test_runner.go'

with open(file_path, 'r', encoding='utf-8', errors='replace') as f:
    lines = f.readlines()

print(f'文件总行数: {len(lines)}')

# 检查第1550-1580行之间的括号匹配
brace_count = 0
for i in range(1550, 1580):
    line = lines[i]
    brace_count += line.count('{') - line.count('}')
    if line.strip():
        print(f'{i+1:5d}: brace={brace_count:2d} {line[:80]}')

print(f'\n第1580行时的brace_count: {brace_count}')

# 如果brace_count < 0，需要在第1579行之前添加闭合大括号
if brace_count < 0:
    # 查找合适的位置插入 }
    # 在第1573行之后添加
    if len(lines) > 1573:
        lines.insert(1573, '\t}\n')
        print('在第1573行之后添加了闭合大括号')
        
        # 重新写入文件
        with open(file_path, 'w', encoding='utf-8') as f:
            f.writelines(lines)
        print('修复完成！')
    else:
        print('无法修复：行号超出范围')
else:
    print('没有发现未闭合的大括号')
