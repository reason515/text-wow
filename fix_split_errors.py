#!/usr/bin/env python3
# -*- coding: utf-8 -*-

# 修复battle.go中的乱码
battle_file = 'server/internal/test/runner/battle.go'
with open(battle_file, 'r', encoding='utf-8', errors='replace') as f:
    content = f.read()

# 修复乱码
content = content.replace('怪物?', '怪物）')
content = content.replace('伤?', '伤害')
content = content.replace('伤?', '伤害')

with open(battle_file, 'w', encoding='utf-8') as f:
    f.write(content)

print('修复了battle.go中的乱码')

# 添加NewTestRunner到test_runner.go
test_runner_file = 'server/internal/test/runner/test_runner.go'
backup_file = 'server/internal/test/runner/test_runner.go.backup_split'

# 从备份文件中读取NewTestRunner
with open(backup_file, 'r', encoding='utf-8', errors='replace') as f:
    backup_content = f.read()

# 查找NewTestRunner函数
import re
newtestrunner_match = re.search(r'func NewTestRunner\(\) \*TestRunner \{.*?\n\}', backup_content, re.DOTALL)
if newtestrunner_match:
    newtestrunner_func = newtestrunner_match.group(0)
    
    # 读取当前的test_runner.go
    with open(test_runner_file, 'r', encoding='utf-8', errors='replace') as f:
        current_content = f.read()
    
    # 在TestRunner结构体定义后添加NewTestRunner
    # 查找TestRunner结构体定义结束的位置
    testrunner_end = current_content.find('type TestRunner struct {')
    if testrunner_end != -1:
        # 找到结构体结束的}
        brace_count = 0
        start_pos = testrunner_end
        for i in range(start_pos, len(current_content)):
            if current_content[i] == '{':
                brace_count += 1
            elif current_content[i] == '}':
                brace_count -= 1
                if brace_count == 0:
                    # 在}后添加NewTestRunner
                    insert_pos = i + 1
                    # 找到下一个非空行
                    while insert_pos < len(current_content) and current_content[insert_pos] in ['\n', '\r', ' ']:
                        insert_pos += 1
                    current_content = current_content[:insert_pos] + '\n\n' + newtestrunner_func + '\n\n' + current_content[insert_pos:]
                    break
        
        with open(test_runner_file, 'w', encoding='utf-8') as f:
            f.write(current_content)
        print('添加了NewTestRunner函数到test_runner.go')
    else:
        print('未找到TestRunner结构体定义')
else:
    print('未找到NewTestRunner函数')
