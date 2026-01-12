#!/usr/bin/env python3
# -*- coding: utf-8 -*-

def fix_line146():
    filepath = 'server/internal/test/runner/instruction.go'
    print(f"处理文件: {filepath}")
    
    with open(filepath, 'rb') as f:
        lines = f.readlines()
    
    original_lines = lines[:]
    changed = False
    
    # 修复第146行（索引145）
    if len(lines) > 145:
        line = lines[145]
        # 检查注释和代码是否在同一行
        if b'\t\t//' in line and b'\t\treturn' in line:
            # 找到注释结束位置和return开始位置
            comment_pos = line.find(b'\t\t//')
            return_pos = line.find(b'\t\treturn', comment_pos)
            if return_pos != -1:
                # 分割为两行
                comment_line = line[comment_pos:return_pos].rstrip()
                return_line = line[return_pos:]
                # 确保注释以引号结束
                if not comment_line.endswith(b'"'):
                    comment_line = comment_line.rstrip() + b'")'
                lines[145] = comment_line + b'\n'
                lines.insert(146, return_line)
                changed = True
                print(f"  修复了第 146 行的格式问题")
    
    if changed:
        with open(filepath, 'wb') as f:
            f.writelines(lines)
        print(f"  已保存更改")
    else:
        print(f"  无需更改")

if __name__ == '__main__':
    fix_line146()
