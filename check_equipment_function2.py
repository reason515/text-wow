#!/usr/bin/env python3
# -*- coding: utf-8 -*-

def check_equipment_function2():
    filepath = 'server/internal/test/runner/equipment.go'
    print(f"检查文件: {filepath}")
    
    with open(filepath, 'rb') as f:
        content = f.read()
    
    # 查找 generateEquipmentFromMonster 函数
    func_start = content.find(b'func (tr *TestRunner) generateEquipmentFromMonster')
    if func_start == -1:
        print("  未找到函数定义")
        return
    
    func_line = content[:func_start].count(b'\n') + 1
    print(f"  函数从第 {func_line} 行开始")
    
    # 计算大括号
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
            
            # 显示前后几行
            print(f"    前后上下文:")
            for j in range(max(0, i-3), min(len(lines), i+4)):
                marker = ">>>" if j == i else "   "
                line_num = func_line + j
                print(f"    {marker} Line {line_num}: {repr(lines[j][:100])}")
            break
        
        # 如果大括号计数为0，说明函数结束
        if brace_count == 0 and i > 0:
            print(f"  函数在第 {current_line} 行结束")
            if current_line < 164:
                print(f"    警告: 函数在目标行之前结束！")
            break

if __name__ == '__main__':
    check_equipment_function2()
