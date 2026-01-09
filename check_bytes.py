#!/usr/bin/env python3
# -*- coding: utf-8 -*-

file_path = 'server/internal/test/runner/test_runner.go'

with open(file_path, 'rb') as f:
    content_bytes = f.read()

# 统一换行符
content = content_bytes.decode('utf-8', errors='replace')
content = content.replace('\r\n', '\n').replace('\r', '\n')
lines = content.split('\n')

print(f'文件总行数: {len(lines)}')

# 检查第1579行的字节内容
if len(lines) > 1578:
    line_1579 = lines[1578]
    print(f'\n第1579行内容: {repr(line_1579)}')
    print(f'第1579行字节: {line_1579.encode("utf-8")}')
    print(f'长度: {len(line_1579)}')
    
    # 检查是否有特殊字符
    for i, char in enumerate(line_1579):
        if ord(char) > 127:
            print(f'位置 {i}: {char} (U+{ord(char):04X})')

# 检查第1579行前后50行的内容
print('\n=== 第1579行前后50行（有内容的行） ===')
for i in range(max(0, 1529), min(len(lines), 1630)):
    line = lines[i].strip()
    if line:
        marker = '>>>' if i == 1578 else '   '
        print(f'{marker} {i+1:5d}: {line[:100]}')
