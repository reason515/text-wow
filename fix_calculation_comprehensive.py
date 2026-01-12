#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
全面修复 calculation.go 文件中的编码问题和空行，注意不要产生多余空行
"""

def fix_calculation():
    file_path = 'server/internal/test/runner/calculation.go'
    
    with open(file_path, 'r', encoding='utf-8', errors='replace') as f:
        lines = f.readlines()
    
    fixed_lines = []
    in_import = False
    prev_was_blank = False
    
    for i, line in enumerate(lines):
        stripped = line.strip()
        is_blank = not stripped
        
        # 检测 import 块
        if stripped == 'import (':
            in_import = True
            fixed_lines.append(line)
            prev_was_blank = False
            continue
        
        if in_import and stripped == ')':
            in_import = False
            fixed_lines.append(line)
            prev_was_blank = False
            continue
        
        # import 块内：移除所有空行
        if in_import:
            if is_blank:
                continue
            fixed_lines.append(line)
            prev_was_blank = False
            continue
        
        # 修复编码问题：注释和代码混在一起的情况
        # 第32行：更新角色的属性
        if '// 更新角色的属' in line and 'char.PhysCritDamage = critDamage' in line:
            fixed_lines.append('	// 更新角色的属性\n')
            fixed_lines.append('	char.PhysCritDamage = critDamage\n')
            continue
        
        # 第50行：更新角色的属性
        if '// 更新角色的属' in line and 'char.SpellCritDamage = critDamage' in line:
            fixed_lines.append('	// 更新角色的属性\n')
            fixed_lines.append('	char.SpellCritDamage = critDamage\n')
            continue
        
        # 第257行：获取基础伤害
        if '// 获取基础伤害（如果已计算' in line and 'baseDamage := char.PhysicalAttack' in line:
            fixed_lines.append('	// 获取基础伤害（如果已计算）\n')
            fixed_lines.append('	baseDamage := char.PhysicalAttack\n')
            continue
        
        # 第267行：至少1点伤害
        if 'damageAfterDefense = 1 // 至少1点伤' in line:
            fixed_lines.append('		damageAfterDefense = 1 // 至少1点伤害\n')
            continue
        
        # 第271行：如果没有最终伤害
        if '// 如果没有最终伤害，使用减伤后伤害作为最终伤' in line and 'if _, exists' in line:
            fixed_lines.append('	// 如果没有最终伤害，使用减伤后伤害作为最终伤害\n')
            fixed_lines.append('	if _, exists := tr.context.Variables["final_damage"]; !exists {\n')
            continue
        
        # 第280行：从上下文中获取伤害
        if '// 从上下文中获取伤害' in line and 'var baseDamage int' in line:
            fixed_lines.append('	// 从上下文中获取伤害值\n')
            fixed_lines.append('	var baseDamage int\n')
            continue
        
        # 第301行：更新上下文
        if '// 更新上下' in line and 'tr.safeSetContext("damage_after_defense"' in line:
            fixed_lines.append('		// 更新上下文\n')
            fixed_lines.append('		tr.safeSetContext("damage_after_defense", baseDamage)\n')
            continue
        
        # 修复第306行的语法错误
        if i == 305 and 'if !ok ||' in line and 'char == nil {' not in line:
            # 检查下一行
            if i + 1 < len(lines) and 'return fmt.Errorf("character not found")' in lines[i + 1]:
                fixed_lines.append('	if !ok || char == nil {\n')
                continue
        
        # 处理空行：函数体内最多保留一个空行
        if is_blank:
            # 如果前面不是空行，且不是函数定义前，保留一个空行
            if not prev_was_blank:
                # 检查是否需要空行（函数定义前或注释后）
                if fixed_lines:
                    last_line = fixed_lines[-1].strip()
                    if last_line == '}' or (last_line.startswith('//') and not last_line.endswith(')')):
                        fixed_lines.append(line)
                        prev_was_blank = True
            continue
        
        # 非空行：如果前面有多个空行，减少到一个
        if prev_was_blank and fixed_lines:
            last_line = fixed_lines[-1].strip()
            if last_line and not last_line.endswith('{') and not last_line.endswith('('):
                # 保留空行
                pass
            else:
                # 移除多余的空行
                while fixed_lines and not fixed_lines[-1].strip():
                    fixed_lines.pop()
                if fixed_lines and fixed_lines[-1].strip():
                    fixed_lines.append('\n')
        
        fixed_lines.append(line)
        prev_was_blank = False
    
    # 移除文件末尾的空行
    while fixed_lines and not fixed_lines[-1].strip():
        fixed_lines.pop()
    
    # 写入文件
    with open(file_path, 'w', encoding='utf-8', errors='replace') as f:
        f.writelines(fixed_lines)
    
    original_count = len(lines)
    new_count = len(fixed_lines)
    removed = original_count - new_count
    print(f"修复完成！原文件 {original_count} 行，现在 {new_count} 行，移除了 {removed} 行空行")

if __name__ == '__main__':
    fix_calculation()
