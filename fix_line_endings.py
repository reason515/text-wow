#!/usr/bin/env python3
# -*- coding: utf-8 -*-

file_path = 'server/internal/test/runner/test_runner.go'

print('读取文件...')
with open(file_path, 'rb') as f:
    content_bytes = f.read()

print('处理换行符...')
# 统一换行符为 \n
content = content_bytes.decode('utf-8', errors='replace')
# 替换所有可能的换行符组合
content = content.replace('\r\n', '\n').replace('\r', '\n')
# 移除多余的空行（连续3个以上空行保留2个）
lines = content.split('\n')
print(f'原始行数: {len(lines)}')

# 统计空行
empty_count = 0
for line in lines:
    if not line.strip():
        empty_count += 1

print(f'空行数: {empty_count}')
print(f'非空行数: {len(lines) - empty_count}')

# 检查错误行号对应的实际内容
error_lines = [8995, 10099, 10243, 18483, 21044, 22356, 36277, 36629, 38293, 41717]

print('\n检查错误行号对应的实际内容:')
for err_line in error_lines:
    if err_line <= len(lines):
        line = lines[err_line - 1]
        if line.strip():
            print(f'第 {err_line} 行: {repr(line[:100])}')
        else:
            # 查找附近有内容的行
            for offset in range(-10, 11):
                check_line = err_line - 1 + offset
                if 0 <= check_line < len(lines) and lines[check_line].strip():
                    print(f'第 {err_line} 行 (空行，附近第 {check_line+1} 行有内容): {repr(lines[check_line][:100])}')
                    break
