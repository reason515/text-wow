package runner

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"text-wow/internal/database"
	"text-wow/internal/game"
	"text-wow/internal/models"
	"text-wow/internal/repository"

	"gopkg.in/yaml.v3"
)

// 类型定义

type TestRunner struct {
	parser           *YAMLParser
	assertion        *AssertionExecutor
	reporter         *Reporter
	calculator       *game.Calculator
	equipmentManager *game.EquipmentManager
	context          *TestContext
}

type TestContext struct {
	Characters map[string]*models.Character         // key: character_id
	Monsters   map[string]*models.Monster           // key: monster_id
	Equipments map[string]*models.EquipmentInstance // key: equipment_id
	Variables  map[string]interface{}               // 其他测试变量
}

type TestSuite struct {
	TestSuite   string     `yaml:"test_suite"`
	Description string     `yaml:"description"`
	Version     string     `yaml:"version"`
	Tests       []TestCase `yaml:"tests"`
}

type TestCase struct {
	Name        string      `yaml:"name"`
	Description string      `yaml:"description"`
	Category    string      `yaml:"category"` // unit/integration/e2e
	Priority    string      `yaml:"priority"` // high/medium/low
	Setup       []string    `yaml:"setup"`
	Steps       []TestStep  `yaml:"steps"`
	Assertions  []Assertion `yaml:"assertions"`
	Teardown    []string    `yaml:"teardown"`
	Timeout     int         `yaml:"timeout"`    // �	MaxRounds   int         `yaml:"max_rounds"` // 最大回合数
}

type TestStep struct {
	Action     string   `yaml:"action"`
	Expected   string   `yaml:"expected"`
	Timeout    int      `yaml:"timeout"`
	MaxRounds  int      `yaml:"max_rounds"` // 最大回合数（用�继续战斗直到"等指令）
	Assertions []string `yaml:"assertions"`
}

type Assertion struct {
	Type      string  `yaml:"type"`      // equals/greater_than/less_than/contains/approximately/range
	Target    string  `yaml:"target"`    // 目标路径，如 "character.hp"
	Expected  string  `yaml:"expected"`  // 期望�	Tolerance float64 `yaml:"tolerance"` // 容差（用于approximately�	Message   string  `yaml:"message"`   // 错误消息
}

type TestResult struct {
	TestName   string
	Status     string // passed/failed/skipped
	Duration   time.Duration
	Error      string
	Assertions []AssertionResult
}

type AssertionResult struct {
	Type     string
	Target   string
	Expected string
	Actual   interface{}
	Status   string // passed/failed
	Message  string
	Error    string // 错误信息
}

type TestSuiteResult struct {
	TestSuite    string
	TotalTests   int
	PassedTests  int
	FailedTests  int
	SkippedTests int
	Duration     time.Duration
	Results      []TestResult
}

