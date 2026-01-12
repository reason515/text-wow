#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
修复 character.go 文件中剩余的问题
"""

def fix_remaining():
    file_path = 'server/internal/test/runner/character.go'
    
    with open(file_path, 'r', encoding='utf-8', errors='replace') as f:
        lines = f.readlines()
    
    fixed_lines = []
    
    for i, line in enumerate(lines):
        # 修复第4541行：编码问题
        if i == 4540 and '一个角' in line and ', 1)' in line and '"' not in line.split('"一个角')[1].split(',')[0] if '"一个角' in line else '' else True:
            # 找到 "一个角 的位置，替换为 "一个角色"
            if '"一个角' in line:
                line = line.replace('"一个角, 1)', '"一个角色", 1)')
            else:
                # 尝试其他方式
                line = line.replace('一个角, 1)', '一个角色", 1)')
                if ', 1)' in line and '"' not in line.split(', 1)')[0][-5:]:
                    # 在前面添加引号
                    idx = line.rfind('一个角')
                    if idx > 0:
                        line = line[:idx] + '"一个角色", 1)'
        
        fixed_lines.append(line)
    
    # 写入文件
    with open(file_path, 'w', encoding='utf-8', errors='replace') as f:
        f.writelines(fixed_lines)
    
    print("修复完成！")

if __name__ == '__main__':
    fix_remaining()
