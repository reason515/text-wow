#!/usr/bin/env python3
# -*- coding: utf-8 -*-

def fix_instruction_line81():
    filepath = 'server/internal/test/runner/instruction.go'
    print(f"处理文件: {filepath}")
    
    with open(filepath, 'rb') as f:
        lines = f.readlines()
    
    original_lines = lines[:]
    changed = False
    
    # 检查第81行（索引80）的问题
    if len(lines) > 80:
        line80 = lines[79]
        line81 = lines[80] if len(lines) > 80 else b''
        
        # 如果第80行包含 else if 和 {，且第81行是注释
        if b'} else if' in line80 and b'{' in line80 and line81.strip().startswith(b'//'):
            # 找到 { 的位置
            brace_pos = line80.rfind(b'{')
            if brace_pos != -1:
                # 分割：第一行是 else if ... {，第二行是注释
                before_brace = line80[:brace_pos+1].rstrip()
                after_brace = line80[brace_pos+1:].strip()
                
                # 如果 after_brace 是注释，需要移到下一行
                if after_brace.startswith(b'//'):
                    lines[79] = before_brace + b'\n'
                    # 将注释插入到第81行之前
                    lines.insert(80, after_brace + b'\n')
                    changed = True
                    print(f"  修复了第 80 行的 else if 语句")
    
    if changed:
        with open(filepath, 'wb') as f:
            f.writelines(lines)
        print(f"  已保存更改")
    else:
        print(f"  无需更改")

if __name__ == '__main__':
    fix_instruction_line81()
