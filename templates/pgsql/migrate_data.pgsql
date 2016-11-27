{{define "migrateData"}}
BEGIN;
INSERT INTO admins(id, username, password, name, status)
VALUES(1,'admin','$2a$10$4F5Hpu0NM8Uy4bI/XQWKDO552uK77WwNpi3zIforzLngziZVszk06', 'Administrator', 'active');
COMMIT;
{{end}}
