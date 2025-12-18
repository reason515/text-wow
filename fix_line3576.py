#!/usr/bin/env python3
# -*- coding: utf-8 -*-
with open('server/internal/game/battle_manager.go', 'r', encoding='utf-8', errors='ignore') as f:
    lines = f.readlines()

# 修复第3576行
if len(lines) > 3575:
    line = lines[3575]
    print(f'修复前: {repr(line[:200])}')
    
    # 检查是否有字符串格式问题
    # 可能是 "å\x9c£å\x85\x89æ\x9c?" 这样的字符串没有正确闭合
    # 需要检查是否有 ?" 后面直接跟逗号的情况
    if 'paladin' in line:
        # 检查字符串数组中的字符串是否都正确闭合
        # 如果看到 ?" 后面直接是逗号，可能需要修复
        if '?"' in line and ', "' not in line[line.find('?"')+2:line.find('?"')+10]:
            # 可能需要添加引号
            pass
    
    # 更简单的方法：检查是否有未闭合的字符串
    # 如果看到 ?" 后面直接是 ", 可能是问题
    if '?"' in line and 'paladin' in line:
        # 尝试修复：?" 后面如果是 ", 可能需要改为 ?",
        # 但这里看起来字符串应该是完整的
        # 让我检查实际的字符串内容
        print('需要检查字符串格式')

with open('server/internal/game/battle_manager.go', 'w', encoding='utf-8') as f:
    f.writelines(lines)































