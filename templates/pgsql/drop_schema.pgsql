{{define "dropSchema"}}
BEGIN;

DROP TABLE admins CASCADE;
DROP TYPE enum_admin_status CASCADE;

COMMIT;
{{end}}
