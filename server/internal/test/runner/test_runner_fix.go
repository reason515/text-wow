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

// debugEnabled 控制是否输出调试信息（通过环境变量 TEST_DEBUG 控制）
var debugEnabled = os.Getenv("TEST_DEBUG") == "1" || os.Getenv("TEST_DEBUG") == "true"

