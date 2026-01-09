#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
修复剩余的编码问题
"""
import re
import os
from datetime import datetime

def fix_remaining_encoding(file_path):
    """修复剩余的编码问题"""
    # 读取文件
    with open(file_path, 'r', encoding='utf-8', errors='replace') as f:
        content = f.read()
    
    original_content = content
    
    # 修复字符串中的损坏字符（在 strings.Contains 中，带引号和括号结束）
    fixes = [
        # 修复 strings.Contains 中的损坏字符
        (r'strings\.Contains\(instruction, "创建.*角色.*\?\s*\)\)', lambda m: m.group(0).replace('?))', '"在"))')),
        (r'strings\.Contains\(instruction, "创建.*怪物.*\?\s*\)\)', lambda m: m.group(0).replace('?))', '"在"))')),
        (r'strings\.Contains\(instruction, "计算最大生命\?\s*\)', 'strings.Contains(instruction, "计算最大生命值")'),
        (r'strings\.Contains\(instruction, "计算生命\?\s*\)', 'strings.Contains(instruction, "计算生命值")'),
        (r'strings\.Contains\(instruction, "计算物理暴击\?\s*\)', 'strings.Contains(instruction, "计算物理暴击率")'),
        (r'strings\.Contains\(instruction, "计算法术暴击\?\s*\)', 'strings.Contains(instruction, "计算法术暴击率")'),
        (r'strings\.Contains\(instruction, "计算物理防御\?\s*\)', 'strings.Contains(instruction, "计算物理防御力")'),
        (r'strings\.Contains\(instruction, "计算魔法防御\?\s*\)', 'strings.Contains(instruction, "计算魔法防御力")'),
        (r'strings\.Contains\(instruction, "计算闪避\?\s*\)', 'strings.Contains(instruction, "计算闪避率")'),
        (r'strings\.Contains\(instruction, "次攻\?\s*\)', 'strings.Contains(instruction, "次攻击")'),
        (r'strings\.Contains\(instruction, "计算队伍总生命\?\s*\)', 'strings.Contains(instruction, "计算队伍总生命值")'),
        (r'strings\.Contains\(instruction, "计算减伤后伤\?\s*\)', 'strings.Contains(instruction, "计算减伤后伤害")'),
        (r'strings\.Contains\(instruction, "计算该区\?\s*\)', 'strings.Contains(instruction, "计算该区域")'),
        # 修复注释中合并的行
        (r'//.*同\?\s*tr\.updateAssertionContext\(\)', '// 在setup执行后立即更新断言上下文，确保所有计算属性都被正确同步\n\ttr.updateAssertionContext()'),
        (r'//.*数据\?\s*tr\.updateAssertionContext\(\)', '// 更新断言上下文（同步测试数据）\n\ttr.updateAssertionContext()'),
    ]
    
    for pattern, replacement in fixes:
        if callable(replacement):
            content = re.sub(pattern, replacement, content)
        else:
            content = re.sub(pattern, replacement, content)
    
    # 直接替换损坏字符
    content = content.replace('?))', '"在"))')
    content = content.replace('计算最大生命?)', '计算最大生命值")')
    content = content.replace('计算生命?)', '计算生命值")')
    content = content.replace('计算物理暴击?)', '计算物理暴击率")')
    content = content.replace('计算法术暴击?)', '计算法术暴击率")')
    content = content.replace('计算物理防御?)', '计算物理防御力")')
    content = content.replace('计算魔法防御?)', '计算魔法防御力")')
    content = content.replace('计算闪避?)', '计算闪避率")')
    content = content.replace('次攻?)', '次攻击")')
    content = content.replace('计算队伍总生命?)', '计算队伍总生命值")')
    content = content.replace('计算减伤后伤?)', '计算减伤后伤害")')
    content = content.replace('计算该区?)', '计算该区域")')
    content = content.replace('同?	tr.updateAssertionContext()', '同步\n\ttr.updateAssertionContext()')
    content = content.replace('数据?	tr.updateAssertionContext()', '数据）\n\ttr.updateAssertionContext()')
    
    if content != original_content:
        # 创建备份
        backup_path = f"{file_path}.backup_{datetime.now().strftime('%Y%m%d_%H%M%S')}"
        with open(backup_path, 'w', encoding='utf-8') as f:
            f.write(original_content)
        print(f"已创建备份文件：{backup_path}")
        
        # 写入修复后的内容
        with open(file_path, 'w', encoding='utf-8') as f:
            f.write(content)
        print(f"成功修复编码问题：{file_path}")
        return True
    else:
        print("没有发现需要修复的问题")
        return False

if __name__ == '__main__':
    file_path = 'server/internal/test/runner/test_runner.go'
    if os.path.exists(file_path):
        fix_remaining_encoding(file_path)
        print("编码修复完成！")
    else:
        print(f"文件不存在：{file_path}")
