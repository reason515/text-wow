#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import re

file_path = 'server/internal/test/runner/test_runner.go'

with open(file_path, 'rb') as f:
    content_bytes = f.read()

content = content_bytes.decode('utf-8', errors='replace')
original = content

# 直接替换所有损坏字符模式（包括Unicode替换字符）
replacements = [
    ('"计算最大生命�?)', '"计算最大生命值")'),
    ('"计算生命�?)', '"计算生命值")'),
    ('"计算物理暴击�?)', '"计算物理暴击率")'),
    ('"计算法术暴击�?)', '"计算法术暴击率")'),
    ('"计算物理防御�?)', '"计算物理防御力")'),
    ('"计算魔法防御�?)', '"计算魔法防御力")'),
    ('"计算闪避�?)', '"计算闪避率")'),
    ('"次攻�?)', '"次攻击")'),
    ('"计算队伍总生命�?)', '"计算队伍总生命值")'),
    ('"计算减伤后伤�?)', '"计算减伤后伤害")'),
    ('"计算该区�?)', '"计算该区域")'),
    ('"有队伍生命�?)', '"有队伍生命值")'),
    # 也尝试没有Unicode替换字符的版本（以防万一）
    ('"计算最大生命?)', '"计算最大生命值")'),
    ('"计算生命?)', '"计算生命值")'),
    ('"计算物理暴击?)', '"计算物理暴击率")'),
    ('"计算法术暴击?)', '"计算法术暴击率")'),
    ('"计算物理防御?)', '"计算物理防御力")'),
    ('"计算魔法防御?)', '"计算魔法防御力")'),
    ('"计算闪避?)', '"计算闪避率")'),
    ('"次攻?)', '"次攻击")'),
    ('"计算队伍总生命?)', '"计算队伍总生命值")'),
    ('"计算减伤后伤?)', '"计算减伤后伤害")'),
    ('"计算该区?)', '"计算该区域")'),
    ('"有队伍生命?)', '"有队伍生命值")'),
]

fixed_count = 0
for old, new in replacements:
    if old in content:
        count = content.count(old)
        content = content.replace(old, new)
        fixed_count += count
        print(f'修复了 {count} 处: {old} -> {new}')

# 也使用正则表达式匹配包含Unicode替换字符的模式
patterns = [
    (r'"计算最大生命[^\"]*\ufffd\?"\)', '"计算最大生命值")'),
    (r'"计算生命[^\"]*\ufffd\?"\)', '"计算生命值")'),
    (r'"计算物理暴击[^\"]*\ufffd\?"\)', '"计算物理暴击率")'),
    (r'"计算法术暴击[^\"]*\ufffd\?"\)', '"计算法术暴击率")'),
    (r'"计算物理防御[^\"]*\ufffd\?"\)', '"计算物理防御力")'),
    (r'"计算魔法防御[^\"]*\ufffd\?"\)', '"计算魔法防御力")'),
    (r'"计算闪避[^\"]*\ufffd\?"\)', '"计算闪避率")'),
    (r'"次攻[^\"]*\ufffd\?"\)', '"次攻击")'),
    (r'"计算队伍总生命[^\"]*\ufffd\?"\)', '"计算队伍总生命值")'),
    (r'"计算减伤后伤[^\"]*\ufffd\?"\)', '"计算减伤后伤害")'),
    (r'"计算该区[^\"]*\ufffd\?"\)', '"计算该区域")'),
    (r'"有队伍生命[^\"]*\ufffd\?"\)', '"有队伍生命值")'),
]

for pattern, replacement in patterns:
    new_content = re.sub(pattern, replacement, content)
    if new_content != content:
        print(f'修复了（正则）: {pattern[:30]}...')
        content = new_content

if content != original:
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(content)
    print(f'修复完成！共修复了 {fixed_count} 处')
else:
    print('没有需要修复的内容')
