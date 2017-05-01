{{define "dropSchema"}}
BEGIN;

DROP TABLE IF EXISTS admins CASCADE;
DROP TYPE IF EXISTS enum_admin_status CASCADE;

COMMIT;
{{end}}
