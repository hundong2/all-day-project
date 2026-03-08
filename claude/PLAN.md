# MAAO (Multi-Agent AI Orchestrator) - 최종 개발 계획서

> **프로젝트 코드명**: MAAO  
> **타겟 OS**: macOS, Linux, WSL  
> **실행 방식**: 터미널 Foreground  
> **상태**: 별도 독립 프로젝트 (OneMind-CLI과 분리)  
> **GitHub Actions CI 연동**: 포함  
> **문서 버전**: v0.2.0 | 2026-03-06

---

## 1. 프로젝트 개요

### 1.1 문제 정의
Claude, Gemini, Codex, Copilot CLI를 모두 구독하고 있지만, 각 서비스의 토큰/요청 한도를 충분히 소진하지 못하는 비효율 발생.

### 1.2 솔루션
GitHub 레포지토리를 협업 허브로 활용하여, 4개 AI CLI 에이전트가 자동으로 계획 수립 → 토론 → 개발 → 리뷰 → 머지까지 수행하는 오케스트레이션 시스템.

### 1.3 핵심 원칙
- **GitHub-Native**: 모든 커뮤니케이션은 Issue, Comment, PR을 통해 이루어짐
- **PM 중심 조율**: 지정된 PM 에이전트가 토큰 예산 기반으로 작업 분배
- **Git Worktree 병렬성**: 각 에이전트가 독립 워크트리에서 동시 작업
- **완전 자동화**: plan.md 생성부터 완료 메일까지 무인 운영

---

## 2. 에이전트 특화 전략 (Agent Specialization)

리서치 결과, 2026년 현재 각 AI CLI 도구의 강점이 뚜렷하게 분화되었습니다. 이를 기반으로 역할 특화 전략을 수립합니다.

### 2.1 에이전트별 강점 분석

```
┌─────────────┬──────────────────────────────────────────────────┐
│  에이전트     │  핵심 강점                                        │
├─────────────┼──────────────────────────────────────────────────┤
│ Claude Code │ • 자율적 멀티파일 아키텍처 변경에 가장 뛰어남        │
│             │ • Agent Teams로 서브에이전트 병렬 실행 가능          │
│             │ • 복잡한 리팩토링, 의존성 업그레이드에 최적            │
│             │ • SWE-bench 최고 수준 (80.84%)                     │
│             │ • 세션 지속 우수 (--session-id, --resume)           │
│             │ • 토큰 효율성 가장 좋음 (자동 컴팩션)                │
├─────────────┼──────────────────────────────────────────────────┤
│ Gemini CLI  │ • 1M 토큰 컨텍스트 윈도우 (경쟁사 대비 5~8배)       │
│             │ • 비용 효율성 최고 (무료 티어: 일 1,000건)           │
│             │ • Google Search 그라운딩으로 실시간 정보 활용         │
│             │ • 계획 수립, 스펙 작성에 강점                        │
│             │ • JSON 출력에 토큰 stats 자동 포함                   │
│             │ • 이미지/UI 스케치 → 코드 변환 우수                  │
├─────────────┼──────────────────────────────────────────────────┤
│ Codex CLI   │ • 터미널 워크플로우 최강 (Terminal-Bench 77.3%)     │
│             │ • Rust 기반으로 빠른 실행 속도                       │
│             │ • 샌드박스 우선 접근 (실행 안전성)                    │
│             │ • 디버깅, 트러블슈팅에 보완적 강점                   │
│             │ • 내장 /review 명령으로 코드 리뷰 특화               │
│             │ • exec 모드로 CI/CD 통합 용이                       │
├─────────────┼──────────────────────────────────────────────────┤
│ Copilot CLI │ • GitHub 네이티브 통합 (Issue, PR, Search 내장)     │
│             │ • 멀티 모델 지원 (Claude, GPT, Gemini 전환 가능)     │
│             │ • 전문화 에이전트 (Explore, Task, Review, Plan)      │
│             │ • MCP 서버 내장 (GitHub MCP)                        │
│             │ • 백그라운드 클라우드 위임 (&로 비동기 실행)           │
│             │ • LSP 지원으로 코드 인텔리전스 제공                   │
└─────────────┴──────────────────────────────────────────────────┘
```

### 2.2 역할 기반 태스크 매핑

PM이 이슈를 생성할 때, 각 이슈의 특성에 따라 최적의 에이전트에 자동 배정합니다.

```
┌──────────────────────┬─────────────┬──────────────────────────────┐
│  태스크 유형           │ 우선 배정    │  이유                         │
├──────────────────────┼─────────────┼──────────────────────────────┤
│ 아키텍처 설계/리팩토링  │ Claude Code │ 멀티파일 자율 변경 최강         │
│ 복잡한 비즈니스 로직    │ Claude Code │ 깊은 추론, 높은 정확도          │
│ 대규모 의존성 업그레이드 │ Claude Code │ 의존성 체인 분석력             │
├──────────────────────┼─────────────┼──────────────────────────────┤
│ 프로젝트 계획/PM 역할   │ Gemini CLI  │ 대규모 컨텍스트, 계획 수립 강점  │
│ 스펙 문서 작성          │ Gemini CLI  │ 비용 효율적, 문서화 강점         │
│ UI/디자인 관련 코드     │ Gemini CLI  │ 이미지 인식/변환 우수           │
│ 리서치/정보 수집 필요   │ Gemini CLI  │ Google Search 그라운딩          │
├──────────────────────┼─────────────┼──────────────────────────────┤
│ 버그 수정/디버깅        │ Codex CLI   │ 디버깅 특화, 빠른 실행          │
│ 테스트 코드 작성        │ Codex CLI   │ 샌드박스로 안전한 테스트 실행     │
│ CI/CD 파이프라인 설정   │ Codex CLI   │ exec 모드로 자동화 용이         │
│ 셸 스크립트/인프라      │ Codex CLI   │ Terminal-Bench 최고 점수        │
├──────────────────────┼─────────────┼──────────────────────────────┤
│ GitHub 연동 작업       │ Copilot CLI │ GitHub MCP 내장, 네이티브 통합   │
│ 코드 리뷰              │ Copilot CLI │ 전문화 Review 에이전트 내장     │
│ 코드베이스 탐색/분석    │ Copilot CLI │ Explore 에이전트 특화           │
│ PR 작성/이슈 관리       │ Copilot CLI │ GitHub 기능 직접 접근           │
└──────────────────────┴─────────────┴──────────────────────────────┘
```

### 2.3 기본 PM 선정 전략

사용자가 PM을 지정하지 않은 경우, **Gemini CLI를 기본 PM으로 선정**합니다.

근거:
- 1M 토큰 컨텍스트로 전체 plan.md + 모든 에이전트 의견을 한 번에 수용 가능
- 비용 효율적 (무료 티어 활용 가능, PM 역할은 토큰 소비가 많음)
- 계획 수립, 태스크 분해에 강점
- Google Search 그라운딩으로 기술 조사 시 실시간 정보 활용

대안 PM 옵션:
- **Claude Code**: 복잡한 대규모 프로젝트에서 더 정확한 태스크 분해 필요 시
- **Copilot CLI**: GitHub 중심 워크플로우가 매우 많은 경우

---

## 3. 기술 스택

### 3.1 핵심 스택

| 영역 | 기술 | 선택 이유 |
|------|------|----------|
| **런타임** | Go 1.22+ | 크로스 플랫폼 단일 바이너리, goroutine 병렬성, CLI 최적 |
| **CLI 프레임워크** | cobra + viper | 서브커맨드, 플래그, 설정 파일 통합 관리 |
| **GitHub API** | go-github/v68 | Issue, Comment, PR, Actions 전체 커버 |
| **Git 조작** | go-git + exec(git) | Worktree 관리는 git CLI 직접 호출, 나머지는 go-git |
| **상태 저장** | SQLite (modernc.org/sqlite) | CGo 없는 순수 Go SQLite, 단일 파일 DB |
| **설정** | YAML (gopkg.in/yaml.v3) | 에이전트/레포 설정에 YAML |
| **TUI 대시보드** | bubbletea + lipgloss | 터미널 실시간 모니터링 UI |
| **로깅** | log/slog (표준 라이브러리) | Go 1.21+ 내장 구조화 로깅 |
| **이메일** | net/smtp (표준 라이브러리) | 완료 알림 발송 |
| **테스트** | testify + gomock | 에이전트 어댑터 모킹 |

### 3.2 각 CLI 도구의 Headless 실행 사양

```go
// Claude Code
cmd := exec.Command("claude", "-p", prompt,
    "--output-format", "json",
    "--allowedTools", "Read,Write,Edit,Bash",
    "--dangerously-skip-permissions",
    "--max-turns", "30")
cmd.Dir = worktreePath

// Gemini CLI
cmd := exec.Command("gemini", "-p", prompt,
    "--output-format", "json")
cmd.Dir = worktreePath

// Codex CLI
cmd := exec.Command("codex", "exec", prompt,
    "--yolo",
    "--jsonl",
    "-m", "gpt-5.3-codex")
cmd.Dir = worktreePath

// Copilot CLI
cmd := exec.Command("copilot",
    "--allow-all-tools",
    "--deny-tool", "shell(rm -rf)",
    "--deny-tool", "shell(git push --force)",
    "-p", prompt)
cmd.Dir = worktreePath
```

---

## 4. 시스템 아키텍처

### 4.1 전체 흐름

```
사용자 ─── plan.md 작성 & push ───→ GitHub Repository
                                          │
                                    MAAO (Foreground)
                                          │
              ┌───────────────────────────┤
              │                           │
         [Poller]                    [State Machine]
     (30s~5m 간격)                        │
     plan.md 감지?  ──── yes ────→  PHASE 1: 계획 토론
     issue 감지?    ──── yes ────→  PHASE 2: 개발 실행
     comment 감지?  ──── yes ────→  에이전트 라우팅
              │                           │
              │                    ┌──────┴──────┐
              │                    │             │
              │              [PM Agent]    [Token Tracker]
              │              (Gemini)      (사용량 추적)
              │                    │
              │         ┌──────────┼──────────┐──────────┐
              │         │          │          │          │
              │    [Worktree]  [Worktree] [Worktree] [Worktree]
              │    claude/     gemini/    codex/     copilot/
              │    issue-3     issue-1    issue-2    issue-4
              │         │          │          │          │
              │         └──────────┼──────────┘──────────┘
              │                    │
              │              [PR & Review Cycle]
              │              (최대 3회)
              │                    │
              │              [GitHub Actions CI]
              │              (자동 테스트 트리거)
              │                    │
              │              [Merge to main]
              │                    │
              │              [Email 알림]
              │                    │
              └────────────── [IDLE] ← 다시 폴링 시작
```

### 4.2 프로젝트 구조

```
maao/
├── cmd/
│   └── maao/
│       └── main.go                 # cobra root command
├── internal/
│   ├── agent/                      # 에이전트 어댑터 계층
│   │   ├── agent.go                # Agent 인터페이스 정의
│   │   ├── registry.go             # 에이전트 레지스트리 & 팩토리
│   │   ├── claude.go               # Claude Code 어댑터
│   │   ├── gemini.go               # Gemini CLI 어댑터
│   │   ├── codex.go                # Codex CLI 어댑터
│   │   ├── copilot.go              # Copilot CLI 어댑터
│   │   └── parser/                 # 에이전트별 출력 파서
│   │       ├── json_parser.go      # Claude, Gemini JSON 파싱
│   │       ├── jsonl_parser.go     # Codex JSONL 파싱
│   │       └── text_parser.go      # Copilot 텍스트 파싱
│   │
│   ├── orchestrator/               # 오케스트레이션 핵심
│   │   ├── pm.go                   # PM 로직 (작업 분배, 특화 전략)
│   │   ├── specializer.go          # 태스크 유형 → 에이전트 매핑
│   │   ├── discussion.go           # 3회 토론 사이클 관리
│   │   ├── workflow.go             # 상태 머신 (Phase 1~3)
│   │   └── reviewer.go             # PR 리뷰 사이클 관리
│   │
│   ├── github/                     # GitHub 연동 계층
│   │   ├── client.go               # go-github 래퍼
│   │   ├── poller.go               # 변경 감지 폴링
│   │   ├── issue.go                # Issue CRUD + 라벨링
│   │   ├── comment.go              # Comment 생성 (멘션 프로토콜)
│   │   ├── pr.go                   # PR 생성/머지
│   │   └── actions.go              # GitHub Actions 트리거 & 결과 확인
│   │
│   ├── workspace/                  # 작업 공간 관리
│   │   ├── worktree.go             # Git Worktree 생성/삭제/관리
│   │   ├── branch.go               # 브랜치 네이밍 전략
│   │   └── cleanup.go              # 머지 후 정리
│   │
│   ├── token/                      # 토큰 관리
│   │   ├── tracker.go              # 실행별 토큰 사용량 기록
│   │   ├── budget.go               # 일일 예산 계산 & 잔여량
│   │   └── estimator.go            # 태스크별 예상 토큰 추정
│   │
│   ├── notify/                     # 알림
│   │   └── email.go                # SMTP 메일 발송
│   │
│   ├── config/                     # 설정
│   │   ├── config.go               # 설정 구조체
│   │   ├── loader.go               # YAML 로더 + 환경 변수 치환
│   │   └── validator.go            # 설정 유효성 검증
│   │
│   ├── store/                      # 상태 저장
│   │   ├── sqlite.go               # SQLite 초기화 & 마이그레이션
│   │   ├── workflow_store.go       # 워크플로우 상태 저장
│   │   └── token_store.go          # 토큰 사용량 히스토리
│   │
│   └── tui/                        # 터미널 대시보드
│       ├── app.go                  # bubbletea 앱
│       ├── dashboard.go            # 메인 대시보드 뷰
│       ├── agent_panel.go          # 에이전트별 상태 패널
│       └── log_viewer.go           # 실시간 로그 스트림
│
├── prompts/                        # 프롬프트 템플릿
│   ├── planning/
│   │   ├── analyze_plan.md         # plan.md 분석 프롬프트
│   │   ├── discuss_round.md        # 토론 라운드 프롬프트
│   │   └── finalize_plan.md        # 최종 계획 확정 프롬프트
│   ├── development/
│   │   ├── implement_issue.md      # 이슈 구현 프롬프트
│   │   └── context_setup.md        # 작업 컨텍스트 설정
│   └── review/
│       ├── code_review.md          # 코드 리뷰 프롬프트
│       └── review_response.md      # 리뷰 피드백 반영 프롬프트
│
├── .maao/                          # 기본 설정 디렉토리
│   └── config.example.yaml
│
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

---

## 5. 상세 설계

### 5.1 설정 파일 (.maao/config.yaml)

```yaml
# ─── 에이전트 설정 ───
agents:
  claude:
    path: "/usr/local/bin/claude"
    enabled: true
    daily_token_budget: 500000
    specialties:                      # 특화 태스크 유형
      - architecture
      - refactoring
      - complex_logic
      - dependency_upgrade
    headless:
      flags: ["-p", "--output-format", "json", "--dangerously-skip-permissions"]
      max_turns: 30
      timeout: "30m"

  gemini:
    path: "/usr/local/bin/gemini"
    enabled: true
    daily_token_budget: 1000000
    specialties:
      - planning
      - documentation
      - spec_writing
      - ui_design
      - research
    headless:
      flags: ["-p", "--output-format", "json"]
      timeout: "20m"

  codex:
    path: "/usr/local/bin/codex"
    enabled: true
    daily_token_budget: 500000
    specialties:
      - debugging
      - testing
      - ci_cd
      - shell_scripting
      - infrastructure
    headless:
      flags: ["exec", "--yolo", "--jsonl"]
      model: "gpt-5.3-codex"
      timeout: "25m"

  copilot:
    path: "/usr/local/bin/copilot"
    enabled: true
    daily_token_budget: 300000
    specialties:
      - github_integration
      - code_review
      - codebase_exploration
      - pr_management
    headless:
      flags: ["--allow-all-tools", "-p"]
      deny_tools:                     # 안전 제한
        - "shell(rm -rf)"
        - "shell(git push --force)"
      timeout: "20m"

# ─── PM 설정 ───
pm:
  default_agent: "gemini"
  discussion_rounds: 3
  review_max_rounds: 3
  task_estimation:                    # 토큰 추정 기준
    small_loc_threshold: 100          # < 100 LOC → ~10K tokens
    medium_loc_threshold: 500         # 100-500 LOC → ~30K tokens
    large_token_estimate: 80000       # > 500 LOC → ~80K tokens

# ─── 레포지토리 ───
repositories:
  - url: "https://github.com/user/project"
    local_path: "/home/user/projects/project"
    poll_interval: "60s"
    default_branch: "main"

# ─── GitHub ───
github:
  token: "${GITHUB_TOKEN}"
  api_url: "https://api.github.com"

# ─── 알림 ───
notification:
  email:
    to: "${MAAO_EMAIL}"
    smtp_host: "smtp.gmail.com"
    smtp_port: 587
    username: "${SMTP_USERNAME}"
    password: "${SMTP_PASSWORD}"

# ─── 워크스페이스 ───
workspace:
  worktree_dir: ".worktrees"
  branch_prefix: "agent/"
  auto_cleanup: true

# ─── GitHub Actions ───
ci:
  enabled: true
  wait_for_completion: true
  timeout: "15m"
  required_checks:                    # 머지 전 필수 통과 체크
    - "test"
    - "lint"
    - "build"
```

### 5.2 에이전트 인터페이스 & 특화 시스템

```go
// internal/agent/agent.go

type TaskType string

const (
    TaskArchitecture      TaskType = "architecture"
    TaskRefactoring       TaskType = "refactoring"
    TaskComplexLogic      TaskType = "complex_logic"
    TaskPlanning          TaskType = "planning"
    TaskDocumentation     TaskType = "documentation"
    TaskDebugging         TaskType = "debugging"
    TaskTesting           TaskType = "testing"
    TaskCICD              TaskType = "ci_cd"
    TaskCodeReview        TaskType = "code_review"
    TaskGitHubIntegration TaskType = "github_integration"
    TaskGeneral           TaskType = "general"
)

type Agent interface {
    Name() string
    Execute(ctx context.Context, req ExecuteRequest) (*ExecuteResponse, error)
    IsAvailable() bool
    Specialties() []TaskType
    TokenBudget() TokenBudget
}

type ExecuteRequest struct {
    WorkDir   string
    Prompt    string
    Context   []string          // 추가 컨텍스트 파일
    SessionID string            // 세션 지속 (지원 시)
    Timeout   time.Duration
}

type ExecuteResponse struct {
    Output     string
    Tokens     TokenUsage
    ExitCode   int
    Duration   time.Duration
    SessionID  string            // 다음 호출에서 세션 이어가기용
}

type TokenUsage struct {
    Input    int
    Output   int
    Total    int
    Cached   int                 // Gemini/Claude 캐시 토큰
}

type TokenBudget struct {
    DailyLimit int
    UsedToday  int
    Remaining  int
}
```

```go
// internal/orchestrator/specializer.go

type Specializer struct {
    agents map[string]agent.Agent
}

// AnalyzeTask: 이슈 제목/본문을 분석하여 TaskType 추론
func (s *Specializer) AnalyzeTask(issue *github.Issue) TaskType {
    // 키워드 기반 분류 + PM 에이전트에게 분류 요청 (하이브리드)
    // "refactor" → TaskRefactoring
    // "test", "spec" → TaskTesting
    // "CI", "pipeline", "workflow" → TaskCICD
    // "bug", "fix", "error" → TaskDebugging
    // 등등
}

// AssignAgent: TaskType + 토큰 잔여량 기반 최적 에이전트 선택
func (s *Specializer) AssignAgent(taskType TaskType) agent.Agent {
    // 1. 해당 TaskType을 specialty로 가진 에이전트 필터링
    // 2. 토큰 잔여량이 충분한 에이전트만 후보로
    // 3. 후보 중 잔여 토큰이 가장 많은 에이전트 선택
    // 4. 후보가 없으면 일반(general) 가능한 에이전트 중 선택
}
```

### 5.3 GitHub Comment 통신 프로토콜

```markdown
<!-- ─── 계획 토론 Comment 형식 ─── -->

### 🤖 @gemini (PM) → All Agents

**Phase**: Planning | **Round**: 1/3

plan.md를 분석했습니다. 다음 사항에 대해 각 에이전트의 의견을 요청합니다:

1. **아키텍처**: REST vs gRPC 선택에 대한 의견 (@claude)
2. **테스트 전략**: E2E 테스트 범위 (@codex)
3. **GitHub 워크플로우**: CI 파이프라인 구성 (@copilot)

각 에이전트는 이 이슈에 Comment로 의견을 남겨주세요.

---
📊 **Token Status**: claude=500K/500K | gemini=980K/1M | codex=500K/500K | copilot=300K/300K
🔄 **Auto-generated by MAAO v0.1.0** | 2026-03-06T14:30:00Z
```

```markdown
<!-- ─── 에이전트 응답 Comment 형식 ─── -->

### 🤖 @claude → @gemini (PM)

**Phase**: Planning | **Round**: 1/3 | **Response**

@gemini 아키텍처 관련 의견입니다:

REST API를 권장합니다. 이유:
- 클라이언트가 웹 브라우저 기반이므로 REST가 자연스러움
- gRPC는 서비스 간 통신에는 적합하나 이 프로젝트 규모에는 오버엔지니어링
- OpenAPI spec으로 자동 문서화 가능

추가로 미들웨어 구조에 대해 @codex의 의견이 필요합니다.

---
📊 **Tokens Used This Turn**: 2,340 | **Remaining**: 497,660/500K
🔄 **Auto-generated by MAAO v0.1.0** | 2026-03-06T14:32:15Z
```

```markdown
<!-- ─── 이슈 할당 Comment 형식 ─── -->

### 🤖 @gemini (PM) — Issue Assignment

최종 계획 기반으로 이슈를 생성하고 배정합니다:

| Issue | Title | Assignee | Type | Est. Tokens |
|-------|-------|----------|------|-------------|
| #1 | API 엔드포인트 구현 | @claude | architecture | ~45K |
| #2 | 단위/통합 테스트 작성 | @codex | testing | ~30K |
| #3 | GitHub Actions CI 설정 | @copilot | ci_cd | ~15K |
| #4 | API 문서 자동 생성 | @gemini | documentation | ~20K |

배정 근거: 각 에이전트의 특화 영역 + 잔여 토큰 기반

---
🔄 **Auto-generated by MAAO v0.1.0**
```

### 5.4 워크플로우 상태 머신

```
┌────────┐  plan.md 감지   ┌──────────────┐
│  IDLE  │ ──────────────→ │ PLAN_ANALYZE │
└────────┘                 └──────┬───────┘
     ↑                           │ PM이 plan.md 읽고 분석
     │                           ▼
     │                    ┌──────────────┐
     │                    │ DISCUSS_R1   │ ← 각 에이전트 의견 수집
     │                    └──────┬───────┘
     │                           │
     │                           ▼
     │                    ┌──────────────┐
     │                    │ DISCUSS_R2   │ ← PM 종합 + 추가 질문
     │                    └──────┬───────┘
     │                           │
     │                           ▼
     │                    ┌──────────────┐
     │                    │ DISCUSS_R3   │ ← 최종 합의
     │                    └──────┬───────┘
     │                           │
     │                           ▼
     │                    ┌──────────────┐
     │                    │ PLAN_FINALIZE│ ← final-plan.md 커밋
     │                    └──────┬───────┘
     │                           │
     │                           ▼
     │                    ┌──────────────┐
     │                    │ ISSUE_CREATE │ ← PM이 이슈 생성 & 배정
     │                    └──────┬───────┘
     │                           │
     │                           ▼
     │                    ┌──────────────┐
     │                    │ WORKTREE_SETUP│ ← 에이전트별 worktree 생성
     │                    └──────┬───────┘
     │                           │
     │                           ▼
     │                    ┌──────────────┐
     │                    │ DEVELOPING   │ ← 에이전트들 병렬 코딩
     │                    └──────┬───────┘
     │                           │ 작업 완료
     │                           ▼
     │                    ┌──────────────┐
     │              ┌───→ │ PR_CREATED   │ ← PR 생성 + 리뷰 요청
     │              │     └──────┬───────┘
     │              │            │
     │              │            ▼
     │              │     ┌──────────────┐
     │              │     │ REVIEWING    │ ← 다른 에이전트가 리뷰
     │              │     └──────┬───────┘
     │              │            │
     │              │     수정 필요? ──yes──→ REVIEW_FIX (최대 3회) ─┐
     │              │            │ no                                │
     │              │            │         ←─────────────────────────┘
     │              │            ▼
     │              │     ┌──────────────┐
     │              │     │ CI_RUNNING   │ ← GitHub Actions 실행 대기
     │              │     └──────┬───────┘
     │              │            │
     │              │     CI 실패? ──yes──→ CI_FIX (에이전트에 전달) ─┘
     │              │            │ no
     │              │            ▼
     │              │     ┌──────────────┐
     │              │     │ MERGING      │ ← main에 순차 머지
     │              │     └──────┬───────┘
     │              │            │
     │              │     남은 이슈? ──yes──→ 다음 이슈 (DEVELOPING) ─┘
     │              │            │ no
     │                           ▼
     │                    ┌──────────────┐
     │                    │ CLEANUP      │ ← worktree 삭제, 브랜치 정리
     │                    └──────┬───────┘
     │                           │
     │                           ▼
     │                    ┌──────────────┐
     │                    │ NOTIFY       │ ← 완료 이메일 발송
     │                    └──────┬───────┘
     │                           │
     └───────────────────────────┘
```

### 5.5 Git Worktree 관리

```go
// internal/workspace/worktree.go

type WorktreeManager struct {
    repoPath    string
    worktreeDir string  // ".worktrees"
    prefix      string  // "agent/"
}

// CreateForAgent: 에이전트별 독립 작업 공간 생성
func (wm *WorktreeManager) CreateForAgent(agentName string, issueNum int) (string, error) {
    branch := fmt.Sprintf("%s%s/issue-%d", wm.prefix, agentName, issueNum)
    path := filepath.Join(wm.repoPath, wm.worktreeDir, fmt.Sprintf("%s-issue-%d", agentName, issueNum))

    // git worktree add <path> -b <branch> main
    cmd := exec.Command("git", "worktree", "add", path, "-b", branch, "main")
    cmd.Dir = wm.repoPath
    return path, cmd.Run()
}

// CleanupAgent: 머지 완료 후 worktree 정리
func (wm *WorktreeManager) CleanupAgent(agentName string, issueNum int) error {
    path := filepath.Join(wm.repoPath, wm.worktreeDir, fmt.Sprintf("%s-issue-%d", agentName, issueNum))
    branch := fmt.Sprintf("%s%s/issue-%d", wm.prefix, agentName, issueNum)

    // git worktree remove <path>
    exec.Command("git", "worktree", "remove", path, "--force").Run()
    // git branch -D <branch>
    return exec.Command("git", "branch", "-D", branch).Run()
}

// 결과 디렉토리 구조:
// project/
// ├── .worktrees/
// │   ├── claude-issue-3/       ← agent/claude/issue-3 브랜치
// │   ├── gemini-issue-1/       ← agent/gemini/issue-1 브랜치
// │   ├── codex-issue-2/        ← agent/codex/issue-2 브랜치
// │   └── copilot-issue-4/      ← agent/copilot/issue-4 브랜치
```

### 5.6 GitHub Actions CI 연동

```go
// internal/github/actions.go

type ActionsManager struct {
    client *github.Client
}

// TriggerCI: PR에 대한 CI 워크플로우 트리거 확인
func (am *ActionsManager) WaitForChecks(ctx context.Context, owner, repo string, prNum int, requiredChecks []string, timeout time.Duration) error {
    deadline := time.After(timeout)
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-deadline:
            return fmt.Errorf("CI timeout after %v", timeout)
        case <-ticker.C:
            pr, _, _ := am.client.PullRequests.Get(ctx, owner, repo, prNum)
            // 모든 required checks가 success인지 확인
            // success → return nil
            // failure → return error (에이전트에 수정 요청)
        case <-ctx.Done():
            return ctx.Err()
        }
    }
}
```

---

## 6. CLI 사용자 인터페이스

### 6.1 명령어 구조

```bash
# 설정 초기화
maao init                           # .maao/config.yaml 생성

# 레포 등록
maao register <repo-url>            # 모니터링 대상 레포 추가
maao unregister <repo-url>          # 레포 제거

# 에이전트 관리
maao agents check                   # 모든 에이전트 실행 가능 여부 확인
maao agents status                  # 에이전트별 토큰 사용량 조회

# 실행
maao run                            # foreground 실행 (메인 명령어)
maao run --verbose                  # 상세 로그 출력
maao run --dashboard                # TUI 대시보드 모드

# 상태 조회
maao status                         # 현재 워크플로우 상태
maao status --repo <url>            # 특정 레포 상태
maao logs                           # 최근 로그 조회
maao logs --agent claude            # 특정 에이전트 로그

# 설정
maao config show                    # 현재 설정 출력
maao config set pm.default_agent gemini
```

### 6.2 TUI 대시보드 레이아웃

```
┌─ MAAO Dashboard ─────────────────────────────────────────────────┐
│                                                                   │
│  📋 Workflow: user/project | Phase: DEVELOPING | Elapsed: 12m 34s │
│                                                                   │
│  ┌─ Agents ──────────────────────────────────────────────────┐   │
│  │                                                            │   │
│  │  🟢 Claude   [████████░░] 78%  Issue #3  architecture      │   │
│  │  🟢 Gemini   [██████░░░░] 56%  Issue #1  documentation     │   │
│  │  🟡 Codex    [███░░░░░░░] 32%  Issue #2  testing          │   │
│  │  ⏳ Copilot  [░░░░░░░░░░]  0%  Waiting   (queued: #4)     │   │
│  │                                                            │   │
│  └────────────────────────────────────────────────────────────┘   │
│                                                                   │
│  ┌─ Token Budget (Today) ────────────────────────────────────┐   │
│  │  Claude:  ████████░░  423K / 500K remaining               │   │
│  │  Gemini:  █████████░  812K / 1M   remaining               │   │
│  │  Codex:   ████████░░  398K / 500K remaining               │   │
│  │  Copilot: ██████████  300K / 300K remaining               │   │
│  └────────────────────────────────────────────────────────────┘   │
│                                                                   │
│  ┌─ Activity Log ────────────────────────────────────────────┐   │
│  │  14:32:15  [claude]  Implementing REST endpoints...        │   │
│  │  14:31:58  [gemini]  Writing API documentation section 3   │   │
│  │  14:31:42  [codex]   Running test suite: 12/18 passed     │   │
│  │  14:31:20  [pm]      Assigned issue #4 to copilot         │   │
│  │  14:30:55  [codex]   Generated 6 test files               │   │
│  └────────────────────────────────────────────────────────────┘   │
│                                                                   │
│  [q] Quit  [p] Pause  [r] Resume  [l] Full Logs  [s] Settings   │
└───────────────────────────────────────────────────────────────────┘
```

---

## 7. 개발 로드맵

### Phase 1: 기반 & 에이전트 계층 (2주)

**Week 1**
- Go 프로젝트 초기화 (go mod, cobra CLI, Makefile)
- 설정 파일 로더 구현 (YAML + 환경 변수 치환)
- Agent 인터페이스 정의
- Claude Code 어댑터 구현 (headless 실행 + JSON 파싱)
- Gemini CLI 어댑터 구현

**Week 2**
- Codex CLI 어댑터 구현 (JSONL 파싱)
- Copilot CLI 어댑터 구현 (텍스트 파싱)
- `maao agents check` 명령어 (실행 파일 존재 + 인증 확인)
- Git Worktree 관리 모듈
- 단위 테스트

### Phase 2: GitHub 연동 (2주)

**Week 3**
- go-github 클라이언트 래퍼
- Issue CRUD (생성, 라벨링, 할당, 클로즈)
- Comment 생성 (멘션 프로토콜 적용)
- Poller 구현 (plan.md 변경, 새 이슈, 새 코멘트 감지)

**Week 4**
- PR 생성 & 머지
- GitHub Actions 체크 결과 대기
- `maao register` / `maao unregister` 명령어
- 통합 테스트 (실제 GitHub 테스트 레포 대상)

### Phase 3: 오케스트레이션 핵심 (3주)

**Week 5**
- PM 로직 (에이전트 선택, 특화 기반 배정)
- Specializer (태스크 유형 분석 → 에이전트 매핑)
- 토큰 트래커 + SQLite 저장소

**Week 6**
- 3회 토론 사이클 구현 (Discussion Manager)
- 워크플로우 상태 머신 (Phase 1: 계획 토론 전체)
- final-plan.md 생성 & 커밋

**Week 7**
- 상태 머신 (Phase 2: 병렬 개발)
- 에이전트별 worktree 자동 생성 + 병렬 실행 (goroutine)
- 리뷰 사이클 관리 (PR 생성 → 리뷰 → 수정 최대 3회)
- CI 체크 대기 & 실패 시 에이전트에 수정 요청

### Phase 4: 통합 & 완성도 (2주)

**Week 8**
- 상태 머신 (Phase 3: 머지, 정리, 알림)
- 이메일 알림 시스템
- 에러 복구 / 재시도 로직 (에이전트 실패, 네트워크 오류)
- `maao run` foreground 실행 모드

**Week 9**
- TUI 대시보드 (bubbletea)
- `maao status`, `maao logs` 명령어
- Makefile 빌드 (macOS amd64/arm64, Linux amd64)
- E2E 테스트 (전체 워크플로우 시나리오)

### Phase 5: 문서화 & 릴리스 (1주)

**Week 10**
- README.md (설치, 설정, 사용법, 트러블슈팅)
- plan.md 작성 가이드 + 예제
- 설정 파일 레퍼런스 문서
- 릴리스 바이너리 빌드 (goreleaser)
- GitHub 레포 공개

**총 기간: 10주**

---

## 8. 추가 개선 제안

### 8.1 즉시 적용 권장

**프롬프트 캐싱**: Claude Code와 Gemini CLI 모두 캐싱을 지원합니다. 반복되는 시스템 프롬프트(역할 설명, 프로젝트 컨텍스트)를 캐시하면 토큰 절약과 속도 개선이 가능합니다.

**Conflict Resolution**: 두 에이전트가 같은 파일을 다른 이슈에서 수정한 경우, 머지 시 충돌이 발생할 수 있습니다. PM이 머지 순서를 결정하고, 후속 에이전트에게 rebase 후 충돌 해결을 요청하는 로직이 필요합니다.

**Graceful Shutdown**: Ctrl+C 시 현재 진행 중인 에이전트 프로세스를 안전하게 종료하고, 상태를 SQLite에 저장하여 다음 실행 시 이어서 진행할 수 있도록 해야 합니다.

### 8.2 향후 확장

**Webhook 모드**: 폴링 대신 GitHub Webhook으로 실시간 이벤트 수신. smee.io 프록시 또는 내장 HTTP 서버로 구현.

**MCP 서버 모드**: MAAO 자체를 MCP 서버로 노출하여, 에이전트들이 직접 토큰 잔량 조회, 이슈 상태 확인 등을 MCP 도구로 호출.

**Web Dashboard**: TUI 외에 웹 기반 대시보드 제공. 브라우저에서 실시간 모니터링 (SSE 또는 WebSocket).

**멀티 레포 오케스트레이션**: 마이크로서비스처럼 여러 레포에 걸친 작업을 단일 plan.md로 관리.

**에이전트 성능 학습**: 에이전트별 이슈 해결 성공률, 리뷰 통과율을 기록하여 Specializer의 배정 정확도를 지속 개선.

---

## 9. 리스크 & 대응

| 리스크 | 영향 | 대응 |
|--------|------|------|
| CLI 도구 API/플래그 변경 | 에이전트 어댑터 동작 불능 | 어댑터 계층 분리로 격리. 버전 핀닝 + 호환성 테스트 |
| 에이전트 토큰 한도 초과 | 작업 중단 | 실시간 토큰 추적 + 한도 80%에서 경고 + 다른 에이전트로 재배정 |
| Git Worktree 충돌 | 머지 실패 | PM이 머지 순서 제어 + 충돌 시 에이전트에 rebase 요청 |
| GitHub API Rate Limit | 폴링/Comment 실패 | 지수 백오프 재시도 + 폴링 간격 자동 조절 |
| 에이전트 무한 루프 | 토큰 낭비 | --max-turns 제한 + timeout 설정 + 비정상 토큰 소비 감지 시 강제 중단 |
| Copilot CLI 출력 파싱 불안정 | 결과 해석 실패 | 정규식 + LLM 기반 이중 파싱. 파싱 실패 시 raw 텍스트로 폴백 |

---

