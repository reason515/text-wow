#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import re

file_path = 'server/internal/test/runner/test_runner.go'

with open(file_path, 'rb') as f:
    content_bytes = f.read()

content = content_bytes.decode('utf-8', errors='replace')
lines = content.split('\n')

fixed_count = 0

# 修复第35621行的多余括号
if len(lines) > 35620:
    line = lines[35620]
    if 'strings.TrimSuffix(singleCharInstruction, ")")))' in line:
        lines[35620] = line.replace('strings.TrimSuffix(singleCharInstruction, ")")))', 'strings.TrimSuffix(singleCharInstruction, ")"))')
        fixed_count += 1
        print(f'修复了第35621行的多余括号')

content = '\n'.join(lines)

if fixed_count > 0:
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(content)
    print(f'修复完成！共修复了 {fixed_count} 处')
else:
    print('没有需要修复的内容')
