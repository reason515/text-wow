#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import re

file_path = 'server/internal/test/runner/test_runner.go'

with open(file_path, 'rb') as f:
    content_bytes = f.read()

content = content_bytes.decode('utf-8', errors='replace')
original = content

# 直接修复包含损坏字符的行
# 修复模式：strings.Contains(instruction, "计算最大生命?) || strings.Contains(instruction, "计算生命?)
content = re.sub(
    r'strings\.Contains\(instruction, "计算最大生命[^\"]*[\ufffd]?\?"\)',
    'strings.Contains(instruction, "计算最大生命值")',
    content
)

content = re.sub(
    r'strings\.Contains\(instruction, "计算生命[^\"]*[\ufffd]?\?"\)',
    'strings.Contains(instruction, "计算生命值")',
    content
)

# 也尝试直接字符串替换（包括Unicode替换字符）
content = content.replace('"计算最大生命?)', '"计算最大生命值")')
content = content.replace('"计算生命?)', '"计算生命值")')
# 尝试匹配包含Unicode替换字符的模式
content = re.sub(r'"计算最大生命[^\"]*[\ufffd]?\?"\)', '"计算最大生命值")', content)
content = re.sub(r'"计算生命[^\"]*[\ufffd]?\?"\)', '"计算生命值")', content)

if content != original:
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(content)
    print('修复完成！')
else:
    print('没有需要修复的内容')
