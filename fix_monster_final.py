#!/usr/bin/env python3
# -*- coding: utf-8 -*-

def fix_monster_final():
    filepath = 'server/internal/test/runner/monster.go'
    
    with open(filepath, 'rb') as f:
        lines = f.readlines()
    
    # 修复第57-59行的问题：将 '")"[0]' 替换为 '")"[0]'
    # 实际上应该是 strings.Split(..., ")")[0]
    for i in range(56, 60):  # 第57-60行，索引从0开始
        if i < len(lines):
            line = lines[i]
            # 查找 '")"[0]' 模式并替换为 '")"[0]'
            if b'")"[0]' in line and b'")"[0]' not in line:
                # 查找 strings.Split(..., ")"[0] 模式
                if b'strings.Split' in line and b'")"[0]' in line:
                    # 替换为 strings.Split(..., ")")[0]
                    new_line = line.replace(b'")"[0]', b'")"[0]')
                    lines[i] = new_line
                    print(f"修复了第 {i+1} 行")
    
    # 修复第97行的问题
    if len(lines) > 96:
        line = lines[96]
        if b'")"[0]' in line and b'strings.Split' in line:
            new_line = line.replace(b'")"[0]', b'")"[0]')
            lines[96] = new_line
            print(f"修复了第 97 行")
    
    with open(filepath, 'wb') as f:
        f.writelines(lines)

if __name__ == '__main__':
    fix_monster_final()
