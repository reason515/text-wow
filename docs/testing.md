# TEXT-WOW 自动化测试指南

## 概述

本项目使用以下测试框架构建了完整的自动化测试体系：

- **后端 (Go)**: 使用 Go 内置测试框架 `testing`
- **前端 (Vue)**: 使用 Vitest + Vue Test Utils

## 快速开始

### 运行所有测试

```bash
# Windows
test.bat

# 或分别运行
test.bat backend    # 仅后端测试
test.bat frontend   # 仅前端测试
```

### 监视模式（开发时推荐）

```bash
test.bat watch
```

### 生成覆盖率报告

```bash
test.bat coverage
```

## 测试结构

### 后端测试

```
server/
├── internal/
│   ├── auth/
│   │   └── auth_test.go        # 认证模块测试
│   ├── api/
│   │   └── handlers_test.go    # API处理器测试
│   ├── repository/
│   │   └── user_repo_test.go   # 用户仓库测试
│   └── database/
│       └── testdb.go           # 测试数据库工具
```

### 前端测试

```
client/
├── src/
│   ├── test/
│   │   └── setup.ts            # 测试全局配置
│   ├── api/
│   │   └── client.test.ts      # API客户端测试
│   ├── stores/
│   │   └── auth.test.ts        # Auth Store测试
│   └── components/
│       └── AuthScreen.test.ts  # 登录组件测试
├── vitest.config.ts            # Vitest配置
```

## 测试类型

### 1. 单元测试

测试独立的函数和模块。

**后端示例** - 密码哈希测试：
```go
func TestHashPassword(t *testing.T) {
    password := "testPassword123"
    hash, err := HashPassword(password)
    
    if err != nil {
        t.Fatalf("HashPassword failed: %v", err)
    }
    
    if !CheckPassword(password, hash) {
        t.Error("CheckPassword should return true for valid password")
    }
}
```

**前端示例** - Token管理测试：
```typescript
describe('Token Management', () => {
  it('should store token in localStorage', () => {
    setToken('test-token-123')
    expect(localStorage.setItem).toHaveBeenCalledWith('token', 'test-token-123')
  })
})
```

### 2. 集成测试

测试多个模块的协作。

**后端示例** - API端点测试：
```go
func TestHandler_Login_Success(t *testing.T) {
    _, router, cleanup := setupHandlerTest(t)
    defer cleanup()

    // 先注册
    registerBody := models.UserRegister{
        Username: "loginuser",
        Password: "password123",
    }
    makeRequest(router, "POST", "/api/auth/register", registerBody)

    // 然后登录
    loginBody := models.UserCredentials{
        Username: "loginuser",
        Password: "password123",
    }
    w := makeRequest(router, "POST", "/api/auth/login", loginBody)

    if w.Code != http.StatusOK {
        t.Errorf("Expected status 200, got %d", w.Code)
    }
}
```

### 3. 组件测试

测试 Vue 组件的渲染和交互。

```typescript
describe('AuthScreen Component', () => {
  it('should call login on submit', async () => {
    const wrapper = mount(AuthScreen)
    const authStore = useAuthStore()
    const loginSpy = vi.spyOn(authStore, 'login').mockResolvedValue(true)
    
    await wrapper.find('input[type="text"]').setValue('testuser')
    await wrapper.find('input[type="password"]').setValue('password')
    await wrapper.find('form').trigger('submit')
    
    expect(loginSpy).toHaveBeenCalledWith({
      username: 'testuser',
      password: 'password',
    })
  })
})
```

## 测试工具

### Mock 数据

前端测试使用预定义的 mock 数据：

```typescript
import { createMockUser, createMockAuthResponse } from '@/test/setup'

const mockUser = createMockUser({ username: 'testuser' })
const mockAuth = createMockAuthResponse(mockUser)
```

### 测试数据库

后端测试使用内存 SQLite 数据库：

```go
func setupTest(t *testing.T) func() {
    testDB, err := database.SetupTestDB()
    if err != nil {
        t.Fatalf("Failed to setup test database: %v", err)
    }
    
    return func() {
        database.TeardownTestDB(testDB)
    }
}
```

## 测试覆盖的功能

### 认证功能 ✅
- [x] 密码哈希和验证
- [x] JWT Token 生成和验证
- [x] 用户注册
- [x] 用户登录
- [x] 认证中间件
- [x] Token 过期处理

### 用户仓库 ✅
- [x] 创建用户
- [x] 根据ID/用户名查询
- [x] 用户名存在检查
- [x] 更新登录时间
- [x] 金币和击杀数更新

### 前端 Store ✅
- [x] 登录状态管理
- [x] 错误处理
- [x] Token 持久化
- [x] 并发操作

### 组件 ✅
- [x] 表单渲染
- [x] 表单验证
- [x] 模式切换（登录/注册）
- [x] 加载状态
- [x] 错误显示

## 最佳实践

### 1. 测试命名

使用描述性的测试名称：
```go
func TestUserRepository_Create_DuplicateUsername(t *testing.T)
func TestHandler_Login_WrongPassword(t *testing.T)
```

### 2. 隔离测试

每个测试应该独立运行，不依赖其他测试：
```go
func setupTest(t *testing.T) func() {
    // 设置测试环境
    return func() {
        // 清理
    }
}
```

### 3. Mock 外部依赖

```typescript
beforeEach(() => {
  mockFetch.mockReset()
  vi.mocked(localStorage.getItem).mockReturnValue(null)
})
```

### 4. 测试边界条件

```go
func TestHashPassword_SpecialCharacters(t *testing.T) {
    password := "!@#$%^&*()中文🔐"
    hash, err := HashPassword(password)
    // ...
}
```

## 持续集成

建议在 CI/CD 流水线中添加测试步骤：

```yaml
# GitHub Actions 示例
test:
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v4
    
    - name: Setup Go
      uses: actions/setup-go@v5
      with:
        go-version: '1.21'
    
    - name: Setup Node.js
      uses: actions/setup-node@v4
      with:
        node-version: '20'
    
    - name: Run Backend Tests
      run: |
        cd server
        go test ./... -v
    
    - name: Run Frontend Tests
      run: |
        cd client
        npm ci
        npm run test:run
```

## 扩展测试

### 添加新的后端测试

1. 创建 `xxx_test.go` 文件
2. 导入 `testing` 包
3. 使用 `TestXxx` 命名函数

### 添加新的前端测试

1. 创建 `xxx.test.ts` 文件
2. 导入 `describe`, `it`, `expect` from `vitest`
3. 使用 `@vue/test-utils` 测试组件

## 常见问题

### Q: 测试失败提示 "database is locked"
A: 确保每个测试都正确清理数据库连接。使用 `defer cleanup()` 确保清理执行。

### Q: 前端测试报 "Cannot find module"
A: 确保已安装依赖 `npm install`，并检查 `vitest.config.ts` 的路径别名配置。

### Q: 如何跳过某个测试?
后端: `t.Skip("reason")`
前端: `it.skip('test name', ...)`





























