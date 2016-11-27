{{define "get_admin_id"}}
SELECT id, password
FROM admins
WHERE username=$1
{{end}}
