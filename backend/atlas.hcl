variable "database_url" {
  type    = string
  default = getenv("DATABASE_URL")
}

env "local" {
  src = "ent://internal/ent/schema"
  url = "postgres://postgres:postgres@localhost:5432/good_todo_go?sslmode=disable"
  dev = "docker://postgres/17/dev?search_path=public"
  migration {
    dir = "file://internal/ent/migrate/migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}
