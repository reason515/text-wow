#!/usr/bin/env python3
# -*- coding: utf-8 -*-

def fix_instruction_comment():
    filepath = 'server/internal/test/runner/instruction.go'
    
    with open(filepath, 'rb') as f:
        content = f.read()
    
    # 修复第640行的注释问题
    # 查找 "检查金币是否足" 后面跟着 tab 或空格然后 if 的模式
    pattern1 = b'// \xe6\xa3\x80\xe6\x9f\xa5\xe9\x87\x91\xe5\xb8\x81\xe6\x98\xaf\xe5\x90\xa6\xe8\xb6\xb3\xef\xbf\xbd\tif'
    replacement1 = b'// \xe6\xa3\x80\xe6\x9f\xa5\xe9\x87\x91\xe5\xb8\x81\xe6\x98\xaf\xe5\x90\xa6\xe8\xb6\xb3\xe5\xa4\x9f\n\tif'
    
    if pattern1 in content:
        content = content.replace(pattern1, replacement1)
        print("修复了第640行的注释问题")
    else:
        # 尝试其他可能的编码
        pattern2 = b'// \xe6\xa3\x80\xe6\x9f\xa5\xe9\x87\x91\xe5\xb8\x81\xe6\x98\xaf\xe5\x90\xa6\xe8\xb6\xb3'
        if pattern2 in content:
            # 查找这个模式后面跟着 tab 和 if 的位置
            pos = content.find(pattern2)
            if pos != -1:
                # 检查后面是否有 tab 和 if
                after_comment = content[pos+len(pattern2):pos+len(pattern2)+10]
                if b'\tif' in after_comment or b' if' in after_comment:
                    # 替换为正确的格式
                    new_comment = b'// \xe6\xa3\x80\xe6\x9f\xa5\xe9\x87\x91\xe5\xb8\x81\xe6\x98\xaf\xe5\x90\xa6\xe8\xb6\xb3\xe5\xa4\x9f\n'
                    # 找到 if 的位置
                    if_pos = content.find(b'if', pos+len(pattern2))
                    if if_pos != -1:
                        # 替换从注释开始到 if 之前的部分
                        content = content[:pos] + new_comment + content[if_pos:]
                        print("修复了第640行的注释问题（方法2）")
    
    with open(filepath, 'wb') as f:
        f.write(content)

if __name__ == '__main__':
    fix_instruction_comment()
