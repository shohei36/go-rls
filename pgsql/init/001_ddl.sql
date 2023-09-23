CREATE SCHEMA myschema;

ALTER SCHEMA public OWNER TO current_user;
DROP SCHEMA IF EXISTS public CASCADE;
ALTER ROLE current_user IN DATABASE mydb SET search_path='myschema';

SET search_path TO myschema;

CREATE USER tenant WITH LOGIN PASSWORD 'password';
CREATE USER admin WITH BYPASSRLS LOGIN PASSWORD 'password';
GRANT USAGE ON SCHEMA myschema TO tenant;
GRANT USAGE ON SCHEMA myschema TO admin;

CREATE TABLE users (
    tenant_id TEXT NOT NULL,
    id        TEXT NOT NULL,
    name      TEXT NOT NULL,
    gender    TEXT NOT NULL,
    age       TEXT NOT NULL
);

ALTER TABLE users ADD PRIMARY KEY(tenant_id, id);

ALTER TABLE users enable ROW level security;

CREATE POLICY user_tenant_policy ON users TO tenant USING (tenant_id = current_setting('app.current_tenant_id'::text));

GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA myschema TO tenant;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA myschema TO tenant;
GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA myschema TO tenant;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA myschema TO admin;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA myschema TO admin;
GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA myschema TO admin;
