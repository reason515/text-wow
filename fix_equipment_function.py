#!/usr/bin/env python3
# -*- coding: utf-8 -*-

def fix_equipment_function():
    filepath = 'server/internal/test/runner/equipment.go'
    print(f"处理文件: {filepath}")
    
    with open(filepath, 'rb') as f:
        content = f.read()
    
    original_content = content
    
    # 查找函数定义
    func_start = content.find(b'func (tr *TestRunner) generateMultipleEquipments')
    if func_start == -1:
        print("  未找到函数定义")
        return
    
    func_line = content[:func_start].count(b'\n') + 1
    print(f"  函数从第 {func_line} 行开始")
    
    # 计算大括号，找到函数结束位置
    brace_count = 0
    lines = content[func_start:].split(b'\n')
    
    for i, line in enumerate(lines):
        brace_count += line.count(b'{') - line.count(b'}')
        current_line = func_line + i
        
        # 检查第164行
        if current_line == 164:
            print(f"  第164行: brace_count={brace_count}")
            print(f"    内容: {repr(line)}")
            if brace_count == 0:
                print(f"    警告: 大括号计数为0，函数可能已结束")
            elif brace_count < 0:
                print(f"    错误: 大括号计数为负")
            else:
                print(f"    正常: 大括号计数为正，函数仍在进行中")
        
        # 如果大括号计数为0，说明函数结束
        if brace_count == 0 and i > 0:
            print(f"  函数在第 {current_line} 行结束")
            if current_line < 164:
                print(f"    警告: 函数在目标行之前结束！")
                print(f"    需要检查第 {current_line} 行附近的代码")
                # 显示前后几行
                for j in range(max(0, i-3), min(len(lines), i+4)):
                    marker = ">>>" if j == i else "   "
                    line_num = func_line + j
                    print(f"    {marker} Line {line_num}: {repr(lines[j][:100])}")
            break
    
    # 如果函数在第164行之前结束，可能是编码问题
    # 检查第127行附近是否有编码问题
    if brace_count == 0 and current_line < 164:
        print(f"\n  检查第 {current_line} 行附近的编码问题...")
        # 查找可能的编码问题
        problem_line_idx = current_line - func_line
        if problem_line_idx < len(lines):
            problem_line = lines[problem_line_idx]
            print(f"  问题行内容: {repr(problem_line)}")
            
            # 检查是否有UTF-8替换字符
            if b'\xef\xbf\xbd' in problem_line:
                print("  发现UTF-8替换字符，尝试修复...")
                fixed_line = problem_line.replace(b'\xef\xbf\xbd', b'')
                # 替换原内容
                old_line_pos = content.find(problem_line, func_start)
                if old_line_pos != -1:
                    content = content[:old_line_pos] + fixed_line + content[old_line_pos + len(problem_line):]
                    print("  已修复编码问题")
    
    if content != original_content:
        with open(filepath, 'wb') as f:
            f.write(content)
        print(f"  已保存更改 ({len(content)} 字节)")
    else:
        print(f"  无需更改")

if __name__ == '__main__':
    fix_equipment_function()
