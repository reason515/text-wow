#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
修复 test_runner.go 文件中的编码问题
"""
import re
import os
import sys
import shutil
from datetime import datetime

def backup_file(file_path):
    """备份文件"""
    try:
        # 生成备份文件名（包含时间戳）
        timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
        backup_path = f"{file_path}.backup_{timestamp}"
        
        # 复制文件
        shutil.copy2(file_path, backup_path)
        print(f"已创建备份文件：{backup_path}")
        return backup_path
    except Exception as e:
        print(f"警告：创建备份文件时出错 - {e}")
        return None

def fix_encoding_issues(file_path):
    """修复文件中的编码问题"""
    
    # 先备份文件
    backup_path = backup_file(file_path)
    if backup_path is None:
        response = input("备份失败，是否继续修复？(y/n): ")
        if response.lower() != 'y':
            print("已取消修复操作")
            return False
    
    # 读取文件
    try:
        with open(file_path, 'rb') as f:
            content = f.read().decode('utf-8', errors='replace')
    except FileNotFoundError:
        print(f"错误：找不到文件 {file_path}")
        return False
    except Exception as e:
        print(f"错误：读取文件时出错 - {e}")
        return False
    
    original_content = content
    
    # 定义修复规则 - 使用字符串替换和正则表达式结合
    # 先尝试直接字符串替换（匹配包含损坏字符的完整模式）
    fixes_direct = [
        # 第24行：控制?var -> 控制）\nvar
        ('控制', 'var', '控制)\nvar'),
        
        # 第33行：测试运行?type -> 测试运行器\ntype
        ('测试运行', 'type', '测试运行器\ntype'),
        
        # 第42行：测试上下?type -> 测试上下文\ntype
        ('测试上下', 'type', '测试上下文\ntype'),
        
        # 第49行：创建测试运行?func -> 创建测试运行器\nfunc
        ('创建测试运行', 'func', '创建测试运行器\nfunc'),
        
        # 第83行：//?	MaxRounds -> // 秒\n	MaxRounds
        ('//', '\tMaxRounds', '// 秒\n\tMaxRounds'),
        
        # 第99行：期望?	Tolerance -> 期望值\n	Tolerance
        ('期望', '\tTolerance', '期望值\n\tTolerance'),
        
        # 第99行：approximately?	Message -> approximately）\n	Message
        ('approximately', '\tMessage', 'approximately)\n\tMessage'),
        
        # 第191行：影响?	tr.context -> 影响）\n	tr.context
        ('影响', '\ttr.context', '影响)\n\ttr.context'),
        
        # 第272行：使?	if step.MaxRounds -> 使用\n	if step.MaxRounds
        ('使', '\tif step.MaxRounds', '使用\n\tif step.MaxRounds'),
        
        # 第305行：检查词缀数?) -> 检查词缀数值")
        ('检查词缀数', ')', '检查词缀数值")'),
        
        # 第306行：处?		return -> 处理\n		return
        ('处', '\t\treturn', '处理\n\t\treturn'),
        
        # 第312行：创建一?) -> 创建一个")
        ('创建一', ')', '创建一个")'),
    ]
    
    # 使用正则表达式匹配损坏字符（可能是任何非正常字符）
    fixes_regex = [
        # 匹配：中文字符 + 非中文字母数字空白字符 + 目标字符
        (r'控制[^\u4e00-\u9fa5a-zA-Z0-9\s\)]var', '控制)\nvar'),
        (r'测试运行[^\u4e00-\u9fa5a-zA-Z0-9\s]type', '测试运行器\ntype'),
        (r'测试上下[^\u4e00-\u9fa5a-zA-Z0-9\s]type', '测试上下文\ntype'),
        (r'创建测试运行[^\u4e00-\u9fa5a-zA-Z0-9\s]func', '创建测试运行器\nfunc'),
        (r'//[^\u4e00-\u9fa5a-zA-Z0-9\s秒]\s*MaxRounds', '// 秒\n\tMaxRounds'),
        (r'期望[^\u4e00-\u9fa5a-zA-Z0-9\s值]\s*Tolerance', '期望值\n\tTolerance'),
        (r'approximately[^\u4e00-\u9fa5a-zA-Z0-9\s\)]\s*Message', 'approximately)\n\tMessage'),
        (r'影响[^\u4e00-\u9fa5a-zA-Z0-9\s\)]\s*tr\.context', '影响)\n\ttr.context'),
        (r'使[^\u4e00-\u9fa5a-zA-Z0-9\s用]\s*if\s+step\.MaxRounds', '使用\n\tif step.MaxRounds'),
        (r'检查词缀数[^\u4e00-\u9fa5a-zA-Z0-9\s值]\)', '检查词缀数值")'),
        (r'处[^\u4e00-\u9fa5a-zA-Z0-9\s理]\s*return', '处理\n\t\treturn'),
        (r'创建一[^\u4e00-\u9fa5a-zA-Z0-9\s个]\)', '创建一个")'),
    ]
    
    # 先检测损坏字符（调试用）
    print("正在检测损坏字符...")
    lines = content.split('\n')
    for i, line in enumerate(lines[:50], 1):  # 检查前50行
        if '控制' in line and 'var' in line and len(line) > 50:
            idx = line.find('控制')
            if idx >= 0:
                chars = line[idx:idx+25]
                print(f"第{i}行 '控制' 后的字符: {repr(chars)}")
                print(f"  Unicode值: {[ord(c) for c in chars[5:15]]}")
        if '测试运行' in line and 'type' in line:
            idx = line.find('测试运行')
            if idx >= 0:
                chars = line[idx:idx+20]
                print(f"第{i}行 '测试运行' 后的字符: {repr(chars)}")
                print(f"  Unicode值: {[ord(c) for c in chars[4:14]]}")
    
    # 使用更直接的方法：匹配包含损坏字符的特定行模式
    changes_made = False
    
    # 定义修复规则：匹配模式 -> 替换内容
    # 使用 . 匹配任何字符（包括损坏字符）
    fixes = [
        # 匹配：控制 + 任何字符（包括损坏字符）+ var
        (r'控制.{1,5}var', '控制)\nvar'),
        # 匹配：测试运行 + 任何字符 + type
        (r'测试运行.{1,5}type', '测试运行器\ntype'),
        # 匹配：测试上下 + 任何字符 + type
        (r'测试上下.{1,5}type', '测试上下文\ntype'),
        # 匹配：创建测试运行 + 任何字符 + func
        (r'创建测试运行.{1,5}func', '创建测试运行器\nfunc'),
        # 匹配：// + 任何字符 + MaxRounds（但排除包含"秒"的行）
        (r'//(?!.*秒).{1,5}MaxRounds', '// 秒\n\tMaxRounds'),
        # 匹配：期望 + 任何字符 + Tolerance（但排除包含"期望值"的行）
        (r'期望(?!.*期望值).{1,5}Tolerance', '期望值\n\tTolerance'),
        # 匹配：approximately + 任何字符 + Message（但排除包含")"的行）
        (r'approximately(?!.*\)).{1,5}Message', 'approximately)\n\tMessage'),
        # 匹配：影响 + 任何字符 + tr.context（但排除包含"影响)"的行）
        (r'影响(?!.*影响\)).{1,5}tr\.context', '影响)\n\ttr.context'),
        # 匹配：使 + 任何字符 + if step.MaxRounds（但排除包含"使用"的行）
        (r'使(?!.*使用).{1,5}if\s+step\.MaxRounds', '使用\n\tif step.MaxRounds'),
        # 匹配：检查词缀数 + 任何字符 + )（但排除包含"检查词缀数值"的行）
        (r'检查词缀数(?!.*检查词缀数值).{1,5}\)', '检查词缀数值")'),
        # 匹配：处 + 任何字符 + return（但排除包含"处理"的行）
        (r'处(?!.*处理).{1,5}return', '处理\n\t\treturn'),
        # 匹配：创建一 + 任何字符 + )（但排除包含"创建一个"的行）
        (r'创建一(?!.*创建一个).{1,5}\)', '创建一个")'),
    ]
    
    # 应用修复
    for pattern, replacement in fixes:
        new_content = re.sub(pattern, replacement, content)
        if new_content != content:
            changes_made = True
            content = new_content
            print(f"修复匹配：{pattern[:20]}...")
    
    # 修复其他常见的编码问题（在字符串和注释中）
    # 使用更宽泛的匹配，匹配任何包含损坏字符的模式
    additional_fixes = [
        # 修复字符串中的损坏字符：个角? -> 个角色（在字符串中）
        (r'"个角\ufffd\?"', '"个角色"'),
        (r'"个角.{1,5}\?"', '"个角色"'),
        (r'个角\ufffd\?\)', '个角色")'),
        (r'个角.{1,5}\?\)', '个角色")'),
        # 修复：? -> 在（在字符串中）
        (r'"\ufffd\?"', '"在"'),
        (r'创建一\ufffd\?\)', '创建一个")'),
        (r'创建一.{1,5}\?\)', '创建一个")'),
        # 修复：计算物理攻击? -> 计算物理攻击力")
        (r'"计算物理攻击\ufffd\?"', '"计算物理攻击力"'),
        (r'"计算物理攻击.{1,5}\?"', '"计算物理攻击力"'),
        (r'计算物理攻击\ufffd\?\)', '计算物理攻击力")'),
        (r'计算物理攻击.{1,5}\?\)', '计算物理攻击力")'),
        # 修复：计算法术攻击? -> 计算法术攻击力")
        (r'"计算法术攻击\ufffd\?"', '"计算法术攻击力"'),
        (r'"计算法术攻击.{1,5}\?"', '"计算法术攻击力"'),
        (r'计算法术攻击\ufffd\?\)', '计算法术攻击力")'),
        (r'计算法术攻击.{1,5}\?\)', '计算法术攻击力")'),
        # 修复注释中的损坏字符
        (r'敏\ufffd\?', '敏捷'),
        (r'敏.{1,5}\?', '敏捷'),
        (r'指\ufffd\?', '指令'),
        (r'指.{1,5}\?', '指令'),
        (r'排\ufffd\?', '排除'),
        (r'排.{1,5}\?', '排除'),
        (r'包\ufffd\?', '包含'),
        (r'包.{1,5}\?', '包含'),
        (r'角\ufffd\?', '角色'),
        (r'角.{1,5}\?', '角色'),
        # 修复字符串字面量中的损坏字符（在 strings.Contains 中）
        (r'strings\.Contains\(instruction, "个角\ufffd\?"\)', 'strings.Contains(instruction, "个角色")'),
        (r'strings\.Contains\(instruction, "个角.{1,5}\?"\)', 'strings.Contains(instruction, "个角色")'),
        (r'strings\.Contains\(instruction, "\ufffd\?"\)', 'strings.Contains(instruction, "在")'),
        (r'strings\.Contains\(instruction, "计算物理攻击\ufffd\?"\)', 'strings.Contains(instruction, "计算物理攻击力")'),
        (r'strings\.Contains\(instruction, "计算物理攻击.{1,5}\?"\)', 'strings.Contains(instruction, "计算物理攻击力")'),
        (r'strings\.Contains\(instruction, "计算法术攻击\ufffd\?"\)', 'strings.Contains(instruction, "计算法术攻击力")'),
        (r'strings\.Contains\(instruction, "计算法术攻击.{1,5}\?"\)', 'strings.Contains(instruction, "计算法术攻击力")'),
        (r'strings\.Contains\(instruction, "计算最大生命\ufffd\?"\)', 'strings.Contains(instruction, "计算最大生命值")'),
        (r'strings\.Contains\(instruction, "计算最大生命.{1,5}\?"\)', 'strings.Contains(instruction, "计算最大生命值")'),
        (r'strings\.Contains\(instruction, "计算生命\ufffd\?"\)', 'strings.Contains(instruction, "计算生命值")'),
        (r'strings\.Contains\(instruction, "计算生命.{1,5}\?"\)', 'strings.Contains(instruction, "计算生命值")'),
        (r'strings\.Contains\(instruction, "计算物理暴击\ufffd\?"\)', 'strings.Contains(instruction, "计算物理暴击率")'),
        (r'strings\.Contains\(instruction, "计算物理暴击.{1,5}\?"\)', 'strings.Contains(instruction, "计算物理暴击率")'),
        (r'strings\.Contains\(instruction, "计算法术暴击\ufffd\?"\)', 'strings.Contains(instruction, "计算法术暴击率")'),
        (r'strings\.Contains\(instruction, "计算法术暴击.{1,5}\?"\)', 'strings.Contains(instruction, "计算法术暴击率")'),
        (r'strings\.Contains\(instruction, "计算物理防御\ufffd\?"\)', 'strings.Contains(instruction, "计算物理防御力")'),
        (r'strings\.Contains\(instruction, "计算物理防御.{1,5}\?"\)', 'strings.Contains(instruction, "计算物理防御力")'),
        (r'strings\.Contains\(instruction, "计算魔法防御\ufffd\?"\)', 'strings.Contains(instruction, "计算魔法防御力")'),
        (r'strings\.Contains\(instruction, "计算魔法防御.{1,5}\?"\)', 'strings.Contains(instruction, "计算魔法防御力")'),
        (r'strings\.Contains\(instruction, "计算闪避\ufffd\?"\)', 'strings.Contains(instruction, "计算闪避率")'),
        (r'strings\.Contains\(instruction, "计算闪避.{1,5}\?"\)', 'strings.Contains(instruction, "计算闪避率")'),
        (r'strings\.Contains\(instruction, "次攻\ufffd\?"\)', 'strings.Contains(instruction, "次攻击")'),
        (r'strings\.Contains\(instruction, "次攻.{1,5}\?"\)', 'strings.Contains(instruction, "次攻击")'),
        # 修复字符串中的损坏字符（在 strings.Contains 中，带引号结束）
        (r'strings\.Contains\(instruction, "计算最大生命\ufffd\?"\)', 'strings.Contains(instruction, "计算最大生命值")'),
        (r'strings\.Contains\(instruction, "计算最大生命.{1,5}\?"\)', 'strings.Contains(instruction, "计算最大生命值")'),
        (r'strings\.Contains\(instruction, "计算生命\ufffd\?"\)', 'strings.Contains(instruction, "计算生命值")'),
        (r'strings\.Contains\(instruction, "计算生命.{1,5}\?"\)', 'strings.Contains(instruction, "计算生命值")'),
        (r'strings\.Contains\(instruction, "计算物理暴击\ufffd\?"\)', 'strings.Contains(instruction, "计算物理暴击率")'),
        (r'strings\.Contains\(instruction, "计算物理暴击.{1,5}\?"\)', 'strings.Contains(instruction, "计算物理暴击率")'),
        (r'strings\.Contains\(instruction, "计算法术暴击\ufffd\?"\)', 'strings.Contains(instruction, "计算法术暴击率")'),
        (r'strings\.Contains\(instruction, "计算法术暴击.{1,5}\?"\)', 'strings.Contains(instruction, "计算法术暴击率")'),
        (r'strings\.Contains\(instruction, "计算物理防御\ufffd\?"\)', 'strings.Contains(instruction, "计算物理防御力")'),
        (r'strings\.Contains\(instruction, "计算物理防御.{1,5}\?"\)', 'strings.Contains(instruction, "计算物理防御力")'),
        (r'strings\.Contains\(instruction, "计算魔法防御\ufffd\?"\)', 'strings.Contains(instruction, "计算魔法防御力")'),
        (r'strings\.Contains\(instruction, "计算魔法防御.{1,5}\?"\)', 'strings.Contains(instruction, "计算魔法防御力")'),
        (r'strings\.Contains\(instruction, "计算闪避\ufffd\?"\)', 'strings.Contains(instruction, "计算闪避率")'),
        (r'strings\.Contains\(instruction, "计算闪避.{1,5}\?"\)', 'strings.Contains(instruction, "计算闪避率")'),
        (r'strings\.Contains\(instruction, "计算队伍总生命\ufffd\?"\)', 'strings.Contains(instruction, "计算队伍总生命值")'),
        (r'strings\.Contains\(instruction, "计算队伍总生命.{1,5}\?"\)', 'strings.Contains(instruction, "计算队伍总生命值")'),
        (r'strings\.Contains\(instruction, "有队伍生命\ufffd\?"\)', 'strings.Contains(instruction, "有队伍生命值")'),
        (r'strings\.Contains\(instruction, "有队伍生命.{1,5}\?"\)', 'strings.Contains(instruction, "有队伍生命值")'),
        (r'strings\.Contains\(instruction, "队伍攻击\ufffd\?"\)', 'strings.Contains(instruction, "队伍攻击力")'),
        (r'strings\.Contains\(instruction, "队伍攻击.{1,5}\?"\)', 'strings.Contains(instruction, "队伍攻击力")'),
        (r'strings\.Contains\(instruction, "队伍生命\ufffd\?"\)', 'strings.Contains(instruction, "队伍生命值")'),
        (r'strings\.Contains\(instruction, "队伍生命.{1,5}\?"\)', 'strings.Contains(instruction, "队伍生命值")'),
        # 修复字符串中的损坏字符（在 strings.Contains 中，带引号和括号结束）
        (r'strings\.Contains\(instruction, "\ufffd\?"\)\)', 'strings.Contains(instruction, "在"))'),
        (r'strings\.Contains\(instruction, ".{1,5}\?"\)\)', lambda m: 'strings.Contains(instruction, "在"))' if '?' in m.group(0) and len(m.group(0)) < 50 else m.group(0)),
        # 修复字符串中的损坏字符（在条件表达式中，带引号和括号结束）
        (r'\(strings\.Contains\(instruction, "角色"\) && strings\.Contains\(instruction, "\ufffd\?"\)\)', '(strings.Contains(instruction, "角色") && strings.Contains(instruction, "在"))'),
        (r'\(strings\.Contains\(instruction, "怪物"\) && strings\.Contains\(instruction, "\ufffd\?"\)\)', '(strings.Contains(instruction, "怪物") && strings.Contains(instruction, "在"))'),
        # 修复字符串中的损坏字符（在条件表达式中，带引号和括号结束）
        (r'\(strings\.Contains\(instruction, "角色"\) && strings\.Contains\(instruction, "\ufffd\?"\)\)', '(strings.Contains(instruction, "角色") && strings.Contains(instruction, "在"))'),
        (r'\(strings\.Contains\(instruction, "怪物"\) && strings\.Contains\(instruction, "\ufffd\?"\)\)', '(strings.Contains(instruction, "怪物") && strings.Contains(instruction, "在"))'),
        # 修复注释中的损坏字符
        (r'必须\ufffd\?创建', '必须在创建'),
        (r'必须.{1,5}\?创建', '必须在创建'),
        # 修复注释中的损坏字符
        (r'上下文状\ufffd\?', '上下文状态'),
        (r'上下文状.{1,5}\?', '上下文状态'),
        (r'Variables中的\ufffd\?', 'Variables中的值'),
        (r'Variables中的.{1,5}\?', 'Variables中的值'),
    ]
    
    # 修复注释和代码中合并的行（损坏字符导致换行丢失）
    # 匹配：注释 + 损坏字符 + 代码（在同一行）
    content = re.sub(r'//.*状\ufffd\?\s*debugPrint', r'// 调试：检查setup后的上下文状态\n\tdebugPrint', content)
    content = re.sub(r'//.*状.{1,5}\?\s*debugPrint', r'// 调试：检查setup后的上下文状态\n\tdebugPrint', content)
    content = re.sub(r'//.*Variables中的\ufffd\?\s*if', r'// 也检查Variables中的值\n\t\tif', content)
    content = re.sub(r'//.*Variables中的.{1,5}\?\s*if', r'// 也检查Variables中的值\n\t\tif', content)
    content = re.sub(r'//.*速度=80\ufffd\?.*指令\s*return', r'// 处理"创建3个怪物：怪物1（速度=40），怪物2（速度=80）"这样的指令\n\t\treturn', content)
    content = re.sub(r'//.*速度=80.{1,5}.*指令\s*return', r'// 处理"创建3个怪物：怪物1（速度=40），怪物2（速度=80）"这样的指令\n\t\treturn', content)
    # 修复字符串中的损坏字符（在条件表达式中，带引号和括号结束，在同一行）
    content = re.sub(r'strings\.Contains\(instruction, "\ufffd\?"\)\)', 'strings.Contains(instruction, "在"))', content)
    # 修复字符串中的损坏字符（在条件表达式中，带引号和括号结束，在同一行，使用Unicode替换字符）
    content = re.sub(r'strings\.Contains\(instruction, "角色"\) && strings\.Contains\(instruction, "\ufffd\?"\)\)', 'strings.Contains(instruction, "角色") && strings.Contains(instruction, "在"))', content)
    content = re.sub(r'strings\.Contains\(instruction, "怪物"\) && strings\.Contains\(instruction, "\ufffd\?"\)\)', 'strings.Contains(instruction, "怪物") && strings.Contains(instruction, "在"))', content)
    # 更通用的修复：匹配任何包含?的损坏字符模式
    content = re.sub(r'strings\.Contains\(instruction, "角色"\) && strings\.Contains\(instruction, "[^"]*\?"\)\)', 'strings.Contains(instruction, "角色") && strings.Contains(instruction, "在"))', content)
    content = re.sub(r'strings\.Contains\(instruction, "怪物"\) && strings\.Contains\(instruction, "[^"]*\?"\)\)', 'strings.Contains(instruction, "怪物") && strings.Contains(instruction, "在"))', content)
    # 修复括号问题：在 "创建一个" 后面添加右括号
    content = re.sub(r'\(strings\.Contains\(instruction, "创建"\) && strings\.Contains\(instruction, "个角色"\) && !strings\.Contains\(instruction, "创建一个"\) \|\|', '(strings.Contains(instruction, "创建") && strings.Contains(instruction, "个角色") && !strings.Contains(instruction, "创建一个")) ||', content)
    # 修复其他合并的行
    content = re.sub(r'//.*同\ufffd\?\s*tr\.updateAssertionContext\(\)', r'// 在setup执行后立即更新断言上下文，确保所有计算属性都被正确同步\n\ttr.updateAssertionContext()', content)
    content = re.sub(r'//.*同.{1,5}\?\s*tr\.updateAssertionContext\(\)', r'// 在setup执行后立即更新断言上下文，确保所有计算属性都被正确同步\n\ttr.updateAssertionContext()', content)
    content = re.sub(r'//.*数据\ufffd\?\s*tr\.updateAssertionContext\(\)', r'// 更新断言上下文（同步测试数据）\n\ttr.updateAssertionContext()', content)
    content = re.sub(r'//.*数据.{1,5}\?\s*tr\.updateAssertionContext\(\)', r'// 更新断言上下文（同步测试数据）\n\ttr.updateAssertionContext()', content)
    content = re.sub(r'//.*上下\ufffd\?\s*tr\.updateAssertionContext\(\)', r'// 更新断言上下文\n\ttr.updateAssertionContext()', content)
    content = re.sub(r'//.*上下.{1,5}\?\s*tr\.updateAssertionContext\(\)', r'// 更新断言上下文\n\ttr.updateAssertionContext()', content)
    
    for pattern, replacement in additional_fixes:
        new_content = re.sub(pattern, replacement, content)
        if new_content != content:
            changes_made = True
            content = new_content
            print(f"修复额外匹配：{pattern[:30]}...")
    
    # 如果还是没有匹配到，尝试匹配 Unicode 替换字符 (U+FFFD)
    if not changes_made:
        print("尝试匹配 Unicode 替换字符 (U+FFFD)...")
        # Unicode 替换字符的匹配
        fixes_fffd = [
            (r'控制\ufffd\s*var', '控制)\nvar'),
            (r'测试运行\ufffd\s*type', '测试运行器\ntype'),
            (r'测试上下\ufffd\s*type', '测试上下文\ntype'),
            (r'创建测试运行\ufffd\s*func', '创建测试运行器\nfunc'),
            (r'//\ufffd\s*MaxRounds', '// 秒\n\tMaxRounds'),
            (r'期望\ufffd\s*Tolerance', '期望值\n\tTolerance'),
            (r'approximately\ufffd\s*Message', 'approximately)\n\tMessage'),
            (r'影响\ufffd\s*tr\.context', '影响)\n\ttr.context'),
            (r'使\ufffd\s*if\s+step\.MaxRounds', '使用\n\tif step.MaxRounds'),
            (r'检查词缀数\ufffd\s*\)', '检查词缀数值")'),
            (r'处\ufffd\s*return', '处理\n\t\treturn'),
            (r'创建一\ufffd\s*\)', '创建一个")'),
        ]
        for pattern, replacement in fixes_fffd:
            new_content = re.sub(pattern, replacement, content)
            if new_content != content:
                changes_made = True
                content = new_content
                print(f"修复匹配（U+FFFD）：{pattern[:20]}...")
                break
    
    if not changes_made:
        print("仍然未发现需要修复的编码问题")
        return True
    
    # 写回文件
    try:
        with open(file_path, 'wb') as f:
            f.write(content.encode('utf-8'))
        print(f"成功修复编码问题：{file_path}")
        return True
    except Exception as e:
        print(f"错误：写入文件时出错 - {e}")
        return False

def main():
    # 文件路径
    file_path = os.path.join('server', 'internal', 'test', 'runner', 'test_runner.go')
    
    # 检查文件是否存在
    if not os.path.exists(file_path):
        print(f"错误：文件不存在 - {file_path}")
        print(f"当前工作目录：{os.getcwd()}")
        sys.exit(1)
    
    # 修复编码问题
    success = fix_encoding_issues(file_path)
    
    if success:
        print("编码修复完成！")
        sys.exit(0)
    else:
        print("编码修复失败！")
        sys.exit(1)

if __name__ == '__main__':
    main()
