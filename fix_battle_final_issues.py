#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
修复 battle.go 文件中剩余的编码问题和格式问题
"""

def fix_battle_file():
    file_path = 'server/internal/test/runner/battle.go'
    
    with open(file_path, 'r', encoding='utf-8', errors='replace') as f:
        lines = f.readlines()
    
    fixed_lines = []
    
    for i, line in enumerate(lines):
        # 修复编码问题：注释和代码混在一起的情况
        if '// 更新上下' in line and 'tr.context.Characters["character"]' in line:
            # 分离注释和代码
            fixed_lines.append('	// 更新上下文\n')
            fixed_lines.append('	tr.context.Characters["character"] = char\n')
            continue
        
        if '// 更新数据' in line and 'charRepo := repository.NewCharacterRepository()' in line:
            # 分离注释和代码
            fixed_lines.append('			// 更新数据库\n')
            fixed_lines.append('			charRepo := repository.NewCharacterRepository()\n')
            continue
        
        if '// 如果角色死亡，不再获得怒气，直接返' in line and 'tr.context.Characters["character"]' in line:
            # 分离注释和代码
            fixed_lines.append('		// 如果角色死亡，不再获得怒气，直接返回\n')
            fixed_lines.append('		tr.context.Characters["character"] = char\n')
            continue
        
        if '// 解析回合数（执行回合"执行一个回' in line and 'roundNum := 1' in line:
            # 分离注释和代码
            fixed_lines.append('	// 解析回合数（如"执行X回合"或"执行一个回合"）\n')
            fixed_lines.append('	roundNum := 1\n')
            continue
        
        # 修复其他编码问题
        line = line.replace('', '')
        line = line.replace('更新上下', '更新上下文')
        line = line.replace('更新数据', '更新数据库')
        line = line.replace('直接返', '直接返回')
        line = line.replace('执行回合"执行一个回', '如"执行X回合"或"执行一个回合"')
        
        # 在函数定义前添加空行（如果前面没有空行且不是第一个函数）
        if line.strip().startswith('func ') and fixed_lines:
            last_line = fixed_lines[-1].strip()
            if last_line and last_line != '' and not last_line.startswith('//'):
                # 检查是否需要添加空行
                if last_line != '}' and not last_line.startswith('package '):
                    fixed_lines.append('\n')
        
        fixed_lines.append(line)
    
    # 写入文件
    with open(file_path, 'w', encoding='utf-8', errors='replace') as f:
        f.writelines(fixed_lines)
    
    print(f"修复完成！处理了编码问题和格式问题")

if __name__ == '__main__':
    fix_battle_file()
