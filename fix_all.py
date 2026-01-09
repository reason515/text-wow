#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import re

file_path = 'server/internal/test/runner/test_runner.go'

with open(file_path, 'rb') as f:
    content_bytes = f.read()

content = content_bytes.decode('utf-8', errors='replace')
original = content

# 修复所有包含损坏字符的 strings.Contains 行
# 使用更通用的匹配模式，匹配任何包含?的损坏字符
fixes = [
    # 匹配包含Unicode替换字符的模式
    (r'strings\.Contains\(instruction, "计算最大生命[^\"]*[\ufffd]?\?"\)', 'strings.Contains(instruction, "计算最大生命值")'),
    (r'strings\.Contains\(instruction, "计算生命[^\"]*[\ufffd]?\?"\)', 'strings.Contains(instruction, "计算生命值")'),
    # 也匹配不包含Unicode替换字符的模式
    (r'strings\.Contains\(instruction, "计算最大生命[^\"]*\?"\)', 'strings.Contains(instruction, "计算最大生命值")'),
    (r'strings\.Contains\(instruction, "计算生命[^\"]*\?"\)', 'strings.Contains(instruction, "计算生命值")'),
    (r'strings\.Contains\(instruction, "计算物理暴击[^\"]*\?"\)', 'strings.Contains(instruction, "计算物理暴击率")'),
    (r'strings\.Contains\(instruction, "计算法术暴击[^\"]*\?"\)', 'strings.Contains(instruction, "计算法术暴击率")'),
    (r'strings\.Contains\(instruction, "计算物理防御[^\"]*\?"\)', 'strings.Contains(instruction, "计算物理防御力")'),
    (r'strings\.Contains\(instruction, "计算魔法防御[^\"]*\?"\)', 'strings.Contains(instruction, "计算魔法防御力")'),
    (r'strings\.Contains\(instruction, "计算闪避[^\"]*\?"\)', 'strings.Contains(instruction, "计算闪避率")'),
    (r'strings\.Contains\(instruction, "次攻[^\"]*\?"\)', 'strings.Contains(instruction, "次攻击")'),
    (r'strings\.Contains\(instruction, "计算队伍总生命[^\"]*\?"\)', 'strings.Contains(instruction, "计算队伍总生命值")'),
    (r'strings\.Contains\(instruction, "计算减伤后伤[^\"]*\?"\)', 'strings.Contains(instruction, "计算减伤后伤害")'),
    (r'strings\.Contains\(instruction, "计算该区[^\"]*\?"\)', 'strings.Contains(instruction, "计算该区域")'),
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
