-- Create application user for RLS
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'app_user') THEN
        CREATE USER app_user WITH PASSWORD 'app_password';
    END IF;
END
$$;

-- Grant permissions to app_user
GRANT USAGE ON SCHEMA public TO app_user;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO app_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO app_user;

-- Enable RLS on users table
ALTER TABLE "users" ENABLE ROW LEVEL SECURITY;

-- Users table policy (allowing empty tenant_id for email verification)
CREATE POLICY "users_tenant_isolation" ON "users"
    FOR ALL
    USING (
        "tenant_id" = current_setting('app.current_tenant_id', true)
        OR
        current_setting('app.current_tenant_id', true) = ''
        OR
        current_setting('app.current_tenant_id', true) IS NULL
    )
    WITH CHECK (
        "tenant_id" = current_setting('app.current_tenant_id', true)
    );

-- Enable RLS on todos table
ALTER TABLE "todos" ENABLE ROW LEVEL SECURITY;

-- Todos table policy
CREATE POLICY "todos_tenant_isolation" ON "todos"
    FOR ALL
    USING ("tenant_id" = current_setting('app.current_tenant_id', true))
    WITH CHECK ("tenant_id" = current_setting('app.current_tenant_id', true));

-- Force RLS for app_user (bypass for table owner/superuser)
ALTER TABLE "users" FORCE ROW LEVEL SECURITY;
ALTER TABLE "todos" FORCE ROW LEVEL SECURITY;
