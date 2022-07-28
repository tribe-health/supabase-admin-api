package data

const (
	PostgrestConf = `
# postgrest.conf

# The standard connection URI format, documented at
# https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING
db-uri       = "postgres://user:pass@host:5432/dbname"

# The name of which database schema to expose to REST clients
db-schemas   = "api"

# The database role to use when no client authentication is provided.
# Can (and should) differ from user in db-uri
db-anon-role = "anon"
db-admin-role = "admin"
`

	KongConf = `
_format_version: "1.1"

services:
- name: my-supa-service
  url: https://example.com
  plugins:
  - name: key-auth
  routes:
  - name: my-route
    paths:
    - /

consumers:
- username: my-supa-user
  keyauth_credentials:
  - key: my-supa-key
`

	PglistenConf = `
conn = "postgres://localhost/postgres"
notifyName = pglisten
callbackCommand = foo
`

	Pgbouncer = `
[databases]
supabase = host=localhost dbname=supabase auth_user=supauser

[pgbouncer]
pool_mode = session
listen_port = 6432
listen_addr = localhost
auth_type = md5
auth_file = users.txt
logfile = pgbouncer.log
pidfile = pgbouncer.pid
admin_users = someuser
stats_users = stat_collector
`

	PgbouncerNew = `
[databases]
supabase = host=localhost dbname=supabase auth_user=supauser

[pgbouncer]
pool_mode = session
listen_port = 6432
listen_addr = localhost
auth_type = md5
auth_file = users.txt
logfile = pgbouncer.log
pidfile = pgbouncer.pid
admin_users = someuser
stats_users = stat_collector
new_users = new_user
foo = bar
`
)
