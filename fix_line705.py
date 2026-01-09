#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import re

file_path = 'server/internal/test/runner/test_runner.go'

with open(file_path, 'rb') as f:
    content_bytes = f.read()

content = content_bytes.decode('utf-8', errors='replace')
original = content

# 修复所有包含损坏字符的 strings.Contains 行
# 使用更宽松的匹配模式，匹配任何字符直到?)（包括Unicode替换字符）
fixes = [
    # 修复 "计算最大生命?) -> "计算最大生命值")
    # 匹配包含Unicode替换字符的模式
    (r'strings\.Contains\(instruction, "计算最大生命[^\"]*[\ufffd]?\?"\)', 'strings.Contains(instruction, "计算最大生命值")'),
    (r'strings\.Contains\(instruction, "计算最大生命.*?\?"\)', 'strings.Contains(instruction, "计算最大生命值")'),
    # 修复 "计算生命?) -> "计算生命值")
    (r'strings\.Contains\(instruction, "计算生命[^\"]*[\ufffd]?\?"\)', 'strings.Contains(instruction, "计算生命值")'),
    (r'strings\.Contains\(instruction, "计算生命.*?\?"\)', 'strings.Contains(instruction, "计算生命值")'),
    # 修复 "计算物理暴击?) -> "计算物理暴击率")
    (r'strings\.Contains\(instruction, "计算物理暴击[^\"]*[\ufffd]?\?"\)', 'strings.Contains(instruction, "计算物理暴击率")'),
    (r'strings\.Contains\(instruction, "计算物理暴击.*?\?"\)', 'strings.Contains(instruction, "计算物理暴击率")'),
    # 修复 "计算法术暴击?) -> "计算法术暴击率")
    (r'strings\.Contains\(instruction, "计算法术暴击[^\"]*[\ufffd]?\?"\)', 'strings.Contains(instruction, "计算法术暴击率")'),
    (r'strings\.Contains\(instruction, "计算法术暴击.*?\?"\)', 'strings.Contains(instruction, "计算法术暴击率")'),
    # 修复 "计算物理防御?) -> "计算物理防御力")
    (r'strings\.Contains\(instruction, "计算物理防御[^\"]*[\ufffd]?\?"\)', 'strings.Contains(instruction, "计算物理防御力")'),
    (r'strings\.Contains\(instruction, "计算物理防御.*?\?"\)', 'strings.Contains(instruction, "计算物理防御力")'),
    # 修复 "计算魔法防御?) -> "计算魔法防御力")
    (r'strings\.Contains\(instruction, "计算魔法防御[^\"]*[\ufffd]?\?"\)', 'strings.Contains(instruction, "计算魔法防御力")'),
    (r'strings\.Contains\(instruction, "计算魔法防御.*?\?"\)', 'strings.Contains(instruction, "计算魔法防御力")'),
    # 修复 "计算闪避?) -> "计算闪避率")
    (r'strings\.Contains\(instruction, "计算闪避[^\"]*[\ufffd]?\?"\)', 'strings.Contains(instruction, "计算闪避率")'),
    (r'strings\.Contains\(instruction, "计算闪避.*?\?"\)', 'strings.Contains(instruction, "计算闪避率")'),
    # 修复 "次攻?) -> "次攻击")
    (r'strings\.Contains\(instruction, "次攻[^\"]*[\ufffd]?\?"\)', 'strings.Contains(instruction, "次攻击")'),
    (r'strings\.Contains\(instruction, "次攻.*?\?"\)', 'strings.Contains(instruction, "次攻击")'),
]

for pattern, replacement in fixes:
    new_content = re.sub(pattern, replacement, content)
    if new_content != content:
        print(f'修复了: {pattern[:50]}...')
        content = new_content

if content != original:
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(content)
    print('修复完成！')
else:
    print('没有需要修复的内容')
