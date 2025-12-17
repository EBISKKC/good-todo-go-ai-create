# Good Todo Go - AI実装用設計図

このドキュメントは、AIがこのプロジェクトを0から構築するための設計図です。

---

## 1. プロジェクト概要

### 目的
PostgreSQLの **Row Level Security (RLS)** を用いたマルチテナント分離を学習・実践するためのTodoアプリケーション

### 技術的ゴール
- RLS (Row Level Security) によるテナント間のデータ分離
- マルチテナントアーキテクチャの実装パターン
- クリーンアーキテクチャによるGoバックエンドの設計
- OpenAPI (oapi-codegen) を用いた型安全なAPI開発

---

## 2. 技術スタック

### バックエンド
| カテゴリ | 技術 | バージョン |
|----------|------|-----------|
| 言語 | Go | 1.24 |
| Webフレームワーク | Echo | v4 |
| ORM | Ent | v0.14.5 |
| データベース | PostgreSQL | 17 |
| マイグレーション | Atlas | - |
| 認証 | JWT | golang-jwt/jwt/v5 |
| DI | Uber Dig | v1.19.0 |
| API仕様 | OpenAPI 3.0 + oapi-codegen | - |
| テスト | testify, testcontainers-go | - |
| モック | go.uber.org/mock | v0.6.0 |

### フロントエンド
| カテゴリ | 技術 | バージョン |
|----------|------|-----------|
| フレームワーク | Next.js (App Router) | 16 |
| 言語 | TypeScript | 5 |
| UIライブラリ | React | 19 |
| スタイリング | Tailwind CSS | v4 |
| 状態管理 | TanStack React Query | v5 |
| フォーム | React Hook Form + Zod | v7 / v4 |
| UIコンポーネント | Radix UI + shadcn/ui | - |
| APIクライアント | Orval (自動生成) | v7 |

### インフラ
- Docker & Docker Compose
- MailHog (開発用メールサーバー)
- Atlas (マイグレーション管理)

---

## 3. データベース設計

### ER図

```
┌─────────────────┐
│     tenants     │
├─────────────────┤
│ id (PK, UUID)   │
│ name            │
│ slug (UNIQUE)   │
│ created_at      │
│ updated_at      │
└────────┬────────┘
         │ 1:N
         ▼
┌─────────────────────────────────────┐
│              users                   │
├─────────────────────────────────────┤
│ id (PK, UUID)                       │
│ tenant_id (FK) ──────────────────── │ ← RLS適用
│ email                               │
│ password_hash                       │
│ name                                │
│ role (admin/member)                 │
│ email_verified                      │
│ verification_token                  │
│ verification_token_expires_at       │
│ created_at                          │
│ updated_at                          │
├─────────────────────────────────────┤
│ UNIQUE(tenant_id, email)            │
└────────┬────────────────────────────┘
         │ 1:N
         ▼
┌─────────────────────────────────────┐
│              todos                   │
├─────────────────────────────────────┤
│ id (PK, UUID)                       │
│ tenant_id (FK) ──────────────────── │ ← RLS適用
│ user_id (FK)                        │
│ title                               │
│ description                         │
│ completed                           │
│ is_public                           │
│ due_date                            │
│ completed_at                        │
│ created_at                          │
│ updated_at                          │
└─────────────────────────────────────┘
```

### Entスキーマ定義

#### Tenant
```go
field.String("id").NotEmpty().Immutable(),
field.String("name").NotEmpty(),
field.String("slug").NotEmpty().Unique(),
field.Time("created_at").Default(time.Now).Immutable(),
field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
```

#### User
```go
field.String("id").NotEmpty().Immutable(),
field.String("tenant_id").NotEmpty().Immutable(),
field.String("email").NotEmpty(),
field.String("password_hash").NotEmpty().Sensitive(),
field.String("name").Default(""),
field.Enum("role").Values("admin", "member").Default("member"),
field.Bool("email_verified").Default(false),
field.String("verification_token").Optional().Nillable(),
field.Time("verification_token_expires_at").Optional().Nillable(),
field.Time("created_at").Default(time.Now).Immutable(),
field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
// Index: UNIQUE(tenant_id, email)
```

#### Todo
```go
field.String("id").NotEmpty().Immutable(),
field.String("tenant_id").NotEmpty().Immutable(),
field.String("user_id").NotEmpty().Immutable(),
field.String("title").NotEmpty(),
field.Text("description").Optional().Default(""),
field.Bool("completed").Default(false),
field.Bool("is_public").Default(false),
field.Time("due_date").Optional().Nillable(),
field.Time("completed_at").Optional().Nillable(),
field.Time("created_at").Default(time.Now).Immutable(),
field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
```

### RLS (Row Level Security) 設定

```sql
-- RLSの有効化
ALTER TABLE "users" ENABLE ROW LEVEL SECURITY;
ALTER TABLE "todos" ENABLE ROW LEVEL SECURITY;

-- Usersテーブルのポリシー (認証トークン検証を考慮)
CREATE POLICY "users_tenant_isolation" ON "users"
    FOR ALL
    USING (
        "tenant_id" = current_setting('app.current_tenant_id', true)
        OR
        current_setting('app.current_tenant_id', true) = ''
    )
    WITH CHECK (
        "tenant_id" = current_setting('app.current_tenant_id', true)
    );

-- Todosテーブルのポリシー
CREATE POLICY "todos_tenant_isolation" ON "todos"
    FOR ALL
    USING ("tenant_id" = current_setting('app.current_tenant_id', true))
    WITH CHECK ("tenant_id" = current_setting('app.current_tenant_id', true));
```

### PostgreSQLユーザー設計

```sql
-- 管理者ユーザー (マイグレーション用)
-- postgres / postgres

-- アプリケーションユーザー (RLS適用)
CREATE USER app_user WITH PASSWORD 'app_password';
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO app_user;
```

---

## 4. バックエンドアーキテクチャ

### ディレクトリ構成

```
backend/
├── cmd/
│   ├── api/
│   │   └── main.go                    # エントリーポイント
│   └── test_rls/
│       └── main.go                    # RLSテスト用
├── internal/
│   ├── domain/
│   │   ├── model/
│   │   │   ├── user.go                # ドメインモデル
│   │   │   └── todo.go
│   │   └── repository/
│   │       ├── user.go                # リポジトリインターフェース
│   │       ├── todo.go
│   │       ├── auth.go
│   │       └── mock/                  # mockgen生成モック
│   │           ├── user.go
│   │           ├── todo.go
│   │           └── auth.go
│   ├── ent/
│   │   ├── schema/
│   │   │   ├── tenant.go              # Entスキーマ
│   │   │   ├── user.go
│   │   │   └── todo.go
│   │   ├── migrate/
│   │   │   └── migrations/            # Atlasマイグレーション
│   │   └── generate.go                # go generate設定
│   ├── infrastructure/
│   │   ├── database/
│   │   │   ├── db.go                  # DB接続
│   │   │   └── tenant_context.go      # テナントコンテキスト設定
│   │   ├── repository/
│   │   │   ├── user.go                # リポジトリ実装
│   │   │   └── todo.go
│   │   └── environment/
│   │       └── environment.go         # 環境変数
│   ├── usecase/
│   │   ├── auth.go                    # 認証ユースケース
│   │   ├── user.go                    # ユーザーユースケース
│   │   ├── todo.go                    # Todoユースケース
│   │   ├── todo_test.go               # ユニットテスト
│   │   ├── user_test.go
│   │   ├── input/
│   │   │   ├── auth.go                # 入力DTO
│   │   │   ├── user.go
│   │   │   └── todo.go
│   │   └── output/
│   │       ├── auth.go                # 出力DTO
│   │       └── todo.go
│   ├── presentation/
│   │   ├── public/                    # Public API
│   │   │   ├── api/
│   │   │   │   └── api.go             # oapi-codegen生成
│   │   │   ├── controller/
│   │   │   │   ├── auth.go
│   │   │   │   ├── user.go
│   │   │   │   └── todo.go
│   │   │   ├── presenter/
│   │   │   │   ├── auth.go
│   │   │   │   ├── user.go
│   │   │   │   └── todo.go
│   │   │   └── router/
│   │   │       ├── auth.go            # Serverメソッド実装
│   │   │       ├── user.go
│   │   │       ├── todo.go
│   │   │       ├── health.go
│   │   │       ├── context_keys/
│   │   │       │   └── context_keys.go
│   │   │       ├── middleware/
│   │   │       │   └── jwt_auth.go    # JWT認証ミドルウェア
│   │   │       └── dependency/
│   │   │           └── dependency.go  # DIコンテナ設定
│   ├── integration_test/
│   │   ├── common/
│   │   │   ├── setup.go               # Testcontainers設定
│   │   │   └── testdata.go            # テストデータヘルパー
│   │   ├── core/
│   │   │   ├── auth_test.go           # 認証統合テスト
│   │   │   ├── todo_test.go           # Todo統合テスト
│   │   │   ├── user_test.go           # ユーザー統合テスト
│   │   │   └── helper_test.go
│   │   └── rls_test.go                # RLS分離テスト
│   └── pkg/
│       ├── jwt.go                     # JWT処理
│       ├── password.go                # パスワードハッシュ
│       ├── uuid.go                    # UUID生成
│       └── mock/
│           └── uuid.go                # UUIDモック
├── openapi/
│   ├── openapi-public.yaml            # Public API定義
│   ├── config-public.yaml             # oapi-codegen設定
│   ├── paths/
│   │   └── public/
│   │       ├── health.yaml
│   │       ├── auth.yaml
│   │       ├── me.yaml
│   │       └── todo.yaml
│   └── components/
│       └── schemas/
│           ├── auth.yaml
│           ├── user.yaml
│           ├── todo.yaml
│           └── error.yaml
├── docker-compose.yml
├── Dockerfile.api.local
├── Dockerfile.migrate
├── Makefile
├── atlas.hcl
├── go.mod
├── go.sum
├── .env
└── .env.example
```

### クリーンアーキテクチャ層

```
┌─────────────────────────────────────────────────────────────┐
│                    Presentation Layer                        │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │  Controller │  │  Presenter  │  │    Router/Server    │  │
│  └──────┬──────┘  └──────▲──────┘  └──────────┬──────────┘  │
│         │                │                     │             │
│         │ uses           │ returns             │ implements  │
│         ▼                │                     ▼             │
│  ┌─────────────────────────────────────────────────────────┐│
│  │                   oapi-codegen ServerInterface          ││
│  └─────────────────────────────────────────────────────────┘│
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      Usecase Layer                           │
│  ┌─────────────────────────────────────────────────────────┐│
│  │     Interactor (ビジネスロジック)                         ││
│  │     - AuthInteractor                                    ││
│  │     - UserInteractor                                    ││
│  │     - TodoInteractor                                    ││
│  └─────────────────────────────────────────────────────────┘│
│  ┌────────────────────┐  ┌────────────────────┐             │
│  │   Input DTO        │  │   Output DTO       │             │
│  └────────────────────┘  └────────────────────┘             │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      Domain Layer                            │
│  ┌────────────────────┐  ┌────────────────────────────────┐ │
│  │     Model          │  │   Repository Interface         │ │
│  │   - User           │  │   - IUserRepository            │ │
│  │   - Todo           │  │   - ITodoRepository            │ │
│  │   - Tenant         │  │   - IAuthRepository            │ │
│  └────────────────────┘  └────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                   Infrastructure Layer                       │
│  ┌────────────────────┐  ┌────────────────────────────────┐ │
│  │   Repository       │  │     Database                   │ │
│  │   Implementation   │  │   - db.go (Ent Client)         │ │
│  │   - UserRepository │  │   - tenant_context.go (RLS)    │ │
│  │   - TodoRepository │  └────────────────────────────────┘ │
│  └────────────────────┘                                     │
└─────────────────────────────────────────────────────────────┘
```

### API実装パターン

新しいAPIエンドポイントを実装する際の手順:

#### 1. Controller作成 (`internal/presentation/public/controller/`)
```go
type FooController struct {
    fooUsecase   usecase.IFooInteractor
    fooPresenter presenter.IFooPresenter
}

func NewFooController(usecase usecase.IFooInteractor, presenter presenter.IFooPresenter) *FooController {
    return &FooController{fooUsecase: usecase, fooPresenter: presenter}
}

func (fc *FooController) GetFoo(ctx echo.Context) error {
    userID, _ := ctx.Get(context_keys.UserIDContextKey).(string)
    result, err := fc.fooUsecase.Get(ctx.Request().Context(), userID)
    if err != nil { return err }
    return fc.fooPresenter.GetFoo(ctx, result)
}
```

#### 2. Presenter作成 (`internal/presentation/public/presenter/`)
```go
type IFooPresenter interface {
    GetFoo(ctx echo.Context, output *output.FooOutput) error
}

type FooPresenter struct{}

func NewFooPresenter() IFooPresenter {
    return &FooPresenter{}
}

func (p *FooPresenter) GetFoo(ctx echo.Context, output *output.FooOutput) error {
    response := api.FooResponse{...}
    return ctx.JSON(http.StatusOK, response)
}
```

#### 3. Router作成 (`internal/presentation/public/router/`)
```go
// router/foo.go - 個別ファイルにServerメソッドを実装
package router

func (s *Server) GetFoo(c echo.Context) error {
    return s.fooController.GetFoo(c)
}
```

#### 4. Server構造体に追加 (`router/dependency/dependency.go`)
```go
type Server struct {
    fooController *controller.FooController
    // ...
}

func NewServer(..., fooController *controller.FooController) *Server {
    return &Server{fooController: fooController, ...}
}
```

#### 5. DIコンテナ自動生成
```bash
make generate_di  # 自動的にNew*関数を検出してdependency.goを生成
```

---

## 5. フロントエンドアーキテクチャ

### ディレクトリ構成

```
frontend/
├── src/
│   ├── app/                           # Next.js App Router
│   │   ├── layout.tsx                 # ルートレイアウト
│   │   ├── page.tsx                   # ホームページ (リダイレクト)
│   │   ├── favicon.ico
│   │   ├── (auth)/                    # 認証グループ
│   │   │   ├── login/
│   │   │   │   └── page.tsx
│   │   │   └── register/
│   │   │       └── page.tsx
│   │   └── (main)/                    # メイングループ (認証必須)
│   │       └── todos/
│   │           └── page.tsx
│   ├── api/                           # Orval生成APIクライアント
│   │   ├── axios-instance.ts          # Axiosインスタンス設定
│   │   └── public/
│   │       ├── auth/
│   │       │   └── auth.ts
│   │       ├── todo/
│   │       │   └── todo.ts
│   │       ├── user/
│   │       │   └── user.ts
│   │       └── model/                 # 型定義
│   │           └── *.ts
│   ├── components/
│   │   ├── ui/                        # shadcn/uiコンポーネント
│   │   │   ├── button.tsx
│   │   │   ├── card.tsx
│   │   │   ├── checkbox.tsx
│   │   │   ├── dialog.tsx
│   │   │   ├── form.tsx
│   │   │   ├── input.tsx
│   │   │   ├── label.tsx
│   │   │   ├── sonner.tsx
│   │   │   └── tabs.tsx
│   │   ├── todo/
│   │   │   ├── todo-list.tsx
│   │   │   ├── todo-item.tsx
│   │   │   ├── create-todo-dialog.tsx
│   │   │   ├── edit-todo-dialog.tsx
│   │   │   └── view-todo-dialog.tsx
│   │   └── user/
│   │       └── profile-edit-dialog.tsx
│   ├── contexts/
│   │   └── auth-context.tsx           # 認証コンテキスト
│   ├── providers/
│   │   └── query-provider.tsx         # TanStack Query Provider
│   └── lib/
│       └── utils.ts                   # ユーティリティ (cn関数)
├── public/
│   └── *.svg
├── orval.config.ts                    # Orval設定
├── next.config.ts
├── tailwind.config.ts (v4はcss使用)
├── postcss.config.mjs
├── tsconfig.json
├── package.json
├── .env
└── .env.example
```

### 状態管理

```
┌─────────────────────────────────────────────────────────────┐
│                      App Provider Tree                       │
│                                                             │
│  <QueryClientProvider>                                      │
│    <AuthProvider>                                           │
│      <Toaster />                                            │
│      {children}                                             │
│    </AuthProvider>                                          │
│  </QueryClientProvider>                                     │
│                                                             │
└─────────────────────────────────────────────────────────────┘

状態管理:
- サーバー状態: TanStack React Query
- 認証状態: AuthContext (React Context)
- フォーム状態: React Hook Form
```

### 認証フロー

```
┌────────────────┐     ┌────────────────┐     ┌────────────────┐
│   Login Page   │────▶│   AuthContext  │────▶│  Local Storage │
│                │     │                │     │  accessToken   │
│  ユーザー入力   │     │  login()       │     │  refreshToken  │
└────────────────┘     │  logout()      │     └────────────────┘
                       │  user state    │
                       └────────┬───────┘
                                │
                                ▼
                       ┌────────────────┐
                       │ Axios Instance │
                       │                │
                       │ Authorization  │
                       │ Bearer token   │
                       │                │
                       │ 401 → Refresh  │
                       └────────────────┘
```

### Orval設定

```typescript
// orval.config.ts
export default {
  public: {
    input: '../backend/openapi-public.yaml',
    output: {
      target: './src/api/public',
      schemas: './src/api/public/model',
      client: 'react-query',
      mode: 'tags-split',
      override: {
        mutator: {
          path: './src/api/axios-instance.ts',
          name: 'customInstance',
        },
      },
    },
  },
};
```

---

## 6. 認証・認可設計

### JWTトークン構造

```go
type Claims struct {
    UserID   string `json:"user_id"`
    TenantID string `json:"tenant_id"`
    Email    string `json:"email"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}
```

### トークン有効期限
- アクセストークン: 15分
- リフレッシュトークン: 7日

### 認証フロー

```
┌──────────────┐                   ┌──────────────┐                   ┌──────────────┐
│    Client    │                   │   API Server │                   │   Database   │
└──────┬───────┘                   └──────┬───────┘                   └──────┬───────┘
       │                                  │                                  │
       │ POST /auth/register              │                                  │
       │ {email, password}                │                                  │
       │────────────────────────────────▶│                                  │
       │                                  │ Create Tenant                    │
       │                                  │─────────────────────────────────▶│
       │                                  │                                  │
       │                                  │ Create User (email_verified=false)
       │                                  │─────────────────────────────────▶│
       │                                  │                                  │
       │                                  │ Send Verification Email          │
       │◀────────────────────────────────│                                  │
       │ 201 Created                      │                                  │
       │                                  │                                  │
       │ GET /auth/verify-email?token=xxx │                                  │
       │────────────────────────────────▶│                                  │
       │                                  │ Verify Token                     │
       │                                  │─────────────────────────────────▶│
       │                                  │                                  │
       │                                  │ Update email_verified=true       │
       │                                  │─────────────────────────────────▶│
       │◀────────────────────────────────│                                  │
       │ 200 OK                           │                                  │
       │                                  │                                  │
       │ POST /auth/login                 │                                  │
       │ {email, password}                │                                  │
       │────────────────────────────────▶│                                  │
       │                                  │ Verify Credentials               │
       │                                  │─────────────────────────────────▶│
       │                                  │                                  │
       │◀────────────────────────────────│                                  │
       │ {accessToken, refreshToken}      │                                  │
       │                                  │                                  │
```

### RLS連携

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              Request Flow                                    │
│                                                                             │
│  1. Client sends request with JWT                                           │
│     Authorization: Bearer <token>                                           │
│                                                                             │
│  2. JWT Middleware extracts claims                                          │
│     - UserID                                                                │
│     - TenantID                                                              │
│     - Role                                                                  │
│                                                                             │
│  3. Middleware sets PostgreSQL session variable                             │
│     SET app.current_tenant_id = '<tenant_id>'                              │
│                                                                             │
│  4. RLS policy automatically filters data                                   │
│     WHERE tenant_id = current_setting('app.current_tenant_id')             │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 7. APIエンドポイント

### Public API

#### 認証
| メソッド | パス | 説明 | 認証 |
|---------|------|------|------|
| POST | `/api/v1/auth/register` | ユーザー登録 | 不要 |
| POST | `/api/v1/auth/login` | ログイン | 不要 |
| POST | `/api/v1/auth/verify-email` | メール認証 | 不要 |
| POST | `/api/v1/auth/refresh` | トークンリフレッシュ | 不要 |

#### ユーザー
| メソッド | パス | 説明 | 認証 |
|---------|------|------|------|
| GET | `/api/v1/me` | 現在のユーザー情報取得 | 必要 |
| PUT | `/api/v1/me` | プロフィール更新 | 必要 |

#### Todo
| メソッド | パス | 説明 | 認証 |
|---------|------|------|------|
| GET | `/api/v1/todos` | 自分のTodo一覧取得 | 必要 |
| GET | `/api/v1/todos-public` | 公開Todo一覧取得 (テナント内) | 必要 |
| POST | `/api/v1/todos` | Todo作成 | 必要 |
| PUT | `/api/v1/todos/:id` | Todo更新 | 必要 |
| DELETE | `/api/v1/todos/:id` | Todo削除 | 必要 |

#### ヘルスチェック
| メソッド | パス | 説明 | 認証 |
|---------|------|------|------|
| GET | `/health` | ヘルスチェック | 不要 |

---

## 8. 実装手順 (ステップバイステップ)

### Phase 1: プロジェクト初期化

1. **バックエンド初期化**
   ```bash
   mkdir backend && cd backend
   go mod init good-todo-go
   ```

2. **フロントエンド初期化**
   ```bash
   npx create-next-app@latest frontend --typescript --tailwind --app --src-dir
   ```

3. **Docker Compose設定**
   - PostgreSQL
   - MailHog
   - Atlas Migration

### Phase 2: データベース層

1. **Entスキーマ定義**
   - tenant.go
   - user.go
   - todo.go

2. **Atlasマイグレーション設定**
   - atlas.hcl
   - 初期マイグレーション生成

3. **RLS設定マイグレーション**
   - ポリシー作成
   - app_user作成

### Phase 3: インフラ層

1. **DB接続設定**
   - db.go
   - environment.go
   - tenant_context.go (RLS用)

2. **リポジトリ実装**
   - user.go
   - todo.go

### Phase 4: ドメイン層

1. **モデル定義**
   - user.go
   - todo.go

2. **リポジトリインターフェース**
   - IUserRepository
   - ITodoRepository
   - IAuthRepository

### Phase 5: ユースケース層

1. **認証ユースケース**
   - Register
   - Login
   - VerifyEmail
   - RefreshToken

2. **ユーザーユースケース**
   - GetMe
   - UpdateMe

3. **Todoユースケース**
   - List (自分のTodo)
   - ListPublic (公開Todo)
   - Create
   - Update
   - Delete

### Phase 6: プレゼンテーション層

1. **OpenAPI定義** (所与)
   - 各パス定義
   - スキーマ定義

2. **oapi-codegen生成**
   ```bash
   make oapi-gen
   ```

3. **Controller実装**
   - auth.go
   - user.go
   - todo.go

4. **Presenter実装**
   - auth.go
   - user.go
   - todo.go

5. **Router実装**
   - Server構造体
   - 各エンドポイントメソッド

6. **Middleware実装**
   - JWT認証ミドルウェア
   - テナントコンテキスト設定

### Phase 7: DI設定

1. **Uber Dig設定**
   ```bash
   make generate_di
   ```

### Phase 8: フロントエンド実装

1. **shadcn/ui設定**
   ```bash
   npx shadcn@latest init
   npx shadcn@latest add button card input label form dialog checkbox tabs
   ```

2. **Orval設定 & API生成**
   ```bash
   npm run generate:api
   ```

3. **認証コンテキスト**
   - auth-context.tsx
   - axios-instance.ts (インターセプター)

4. **ページ実装**
   - login/page.tsx
   - register/page.tsx
   - todos/page.tsx

5. **コンポーネント実装**
   - todo-list.tsx
   - todo-item.tsx
   - create-todo-dialog.tsx
   - edit-todo-dialog.tsx
   - profile-edit-dialog.tsx

### Phase 9: テスト実装

1. **ユニットテスト**
   - mockgen設定
   - usecase/todo_test.go
   - usecase/user_test.go

2. **統合テスト**
   - testcontainers設定
   - RLSテスト
   - 認証テスト
   - Todoテスト

---

## 9. 開発コマンド

### バックエンド (Makefile)

```bash
# Docker
make run                 # Dockerサービス起動
make stop                # Dockerサービス停止
make init-db             # DB初期化 (データ削除 + マイグレーション)

# 開発
make dev                 # 開発サーバー起動

# コード生成
make generate_ent        # Ent ORMコード生成
make mockgen             # モック生成
make oapi-gen            # OpenAPIコード生成
make generate_di         # DIコンテナ生成

# マイグレーション
make migrate_diff        # マイグレーション作成
make migrate_apply       # マイグレーション適用
make migrate_status      # マイグレーション状態確認
make migrate_down n=1    # マイグレーションロールバック

# テスト
make test_unit           # ユニットテスト
make test_integration    # 統合テスト
make test                # 全テスト実行

# コード品質
make fmt                 # フォーマット
make lint                # リント
make vet                 # vet
```

### フロントエンド

```bash
npm run dev              # 開発サーバー起動
npm run build            # ビルド
npm run generate:api     # APIクライアント生成
npm run lint             # リント
```

---

## 10. 環境変数

### バックエンド (.env)
```env
# PostgreSQL (管理者)
POSTGRES_DB_USER=postgres
POSTGRES_DB_PASSWORD=postgres
POSTGRES_DB_NAME=good_todo_go
POSTGRES_DB_PORT=5432
POSTGRES_DB_HOST=localhost

# PostgreSQL (アプリケーション用 - RLS適用)
POSTGRES_APP_USER=app_user
POSTGRES_APP_PASSWORD=app_password

# JWT
JWT_SECRET=your-super-secret-jwt-key-change-in-production

# Server
PUBLIC_API_PORT=8000

# Mail (MailHog)
SMTP_HOST=localhost
SMTP_PORT=1025
```

### フロントエンド (.env)
```env
NEXT_PUBLIC_API_BASE_URL=http://localhost:8000/api/v1
```

---

## 11. 重要な実装ポイント

### RLSの実装

1. **テナントコンテキスト設定は必須**
   - 認証済みリクエストでは必ず `SET app.current_tenant_id` を実行
   - 未設定の場合、RLSポリシーによりデータアクセス不可

2. **管理者接続とアプリ接続の分離**
   - マイグレーション: postgres (RLSバイパス)
   - アプリケーション: app_user (RLS適用)

3. **メール認証時の特殊処理**
   - 認証トークン検証時はテナントIDが不明
   - RLSポリシーで `current_tenant_id = ''` の場合を許可

### Ent ORMの使用

1. **スキーマ変更後は必ず再生成**
   ```bash
   make generate_ent
   make migrate_diff
   make migrate_apply
   ```

2. **IDフィールドはStringで定義**
   - UUID形式を使用
   - `field.String("id").NotEmpty().Immutable()`

### oapi-codegenの使用

1. **OpenAPI定義変更後は再生成**
   ```bash
   make oapi-gen
   ```

2. **ServerInterface実装**
   - 全メソッドを実装しないとコンパイルエラー
   - `router/` 配下に機能別ファイルで分割実装

---

## 12. テスト戦略

### ユニットテスト
- 対象: Usecase層
- モック: mockgenで生成されたリポジトリモック
- 実行: `make test_unit`

### 統合テスト
- 対象: Controller → Usecase → Repository
- 環境: Testcontainers (PostgreSQL)
- RLS: 実際のRLS適用状態でテスト
- 実行: `make test_integration`

### テスト観点
1. 正常系 (Happy Path)
2. 異常系 (エラーハンドリング)
3. RLSテナント分離 (他テナントデータへのアクセス不可)
4. 認証・認可 (未認証、認証済み、権限不足)

---

この設計図に従って実装を進めることで、RLSを活用したマルチテナントTodoアプリケーションを構築できます。
