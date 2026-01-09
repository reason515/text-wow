#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import re

file_path = 'server/internal/test/runner/test_runner.go'
backup_path = file_path + '.backup_normalize'

print('读取文件...')
with open(file_path, 'rb') as f:
    content_bytes = f.read()

print('处理换行符...')
content = content_bytes.decode('utf-8', errors='replace')
# 统一换行符
content = content.replace('\r\n', '\n').replace('\r', '\n')

# 移除多余的空行（连续3个以上空行保留2个）
lines = content.split('\n')
print(f'原始行数: {len(lines)}')

# 处理多余空行
new_lines = []
prev_empty_count = 0
for line in lines:
    is_empty = not line.strip()
    if is_empty:
        prev_empty_count += 1
        # 最多保留2个连续空行
        if prev_empty_count <= 2:
            new_lines.append('')
    else:
        prev_empty_count = 0
        new_lines.append(line)

print(f'处理后行数: {len(new_lines)}')
print(f'移除了 {len(lines) - len(new_lines)} 个多余空行')

# 创建备份
print('创建备份...')
with open(backup_path, 'wb') as f:
    f.write(content_bytes)

# 写入处理后的内容
print('写入处理后的文件...')
content_normalized = '\n'.join(new_lines)
with open(file_path, 'w', encoding='utf-8') as f:
    f.write(content_normalized)

print('完成！文件已规范化')
