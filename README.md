# Supabase Admin API Server

To run on KPS and administer configs for services:
- Kong
- Gotrue
- Realtime
- Postgrest
- Pg-listen

## API Interface
### Configs
GET `/config/postgrest` - returns current config `{ raw_contents: <string-of-file-contents> }`

POST `/config/postgrest` - sets new config - params: `{ raw_contents: <string-of-file-contents>, restart_services : <bool> }`

GET `/config/kong` - returns current config `{ raw_contents: <string-of-file-contents> }`

POST `/config/kong` - sets new config - params: `{ raw_contents: <string-of-file-contents>, restart_services : <bool> }`

GET `/config/pglisten` - returns current config `{ raw_contents: <string-of-file-contents> }`

POST `/config/pglisten` - sets new config - params: `{ raw_contents: <string-of-file-contents>, restart_services : <bool> }`

GET `/config/realtime` - returns current config `{ raw_contents: <string-of-file-contents> }`

POST `/config/realtime` - sets new config - params: `{ raw_contents: <string-of-file-contents>, restart_services : <bool> }`

GET `/config/goauth` - returns current config as `{ raw_contents: <string-of-file-contents> }`

POST `/config/goauth` - sets new config - params: `{ raw_contents: <string-of-file-contents>, restart_services : <bool> }`

### Restarting
GET `/services/restart` - re-reads all configs and restarts all services

GET `/services/reboot` - reboot the server

### Logs
GET `/logs/<application>/<head|tail>/<max_lines>` - get logs for a given application (postgrest,kong,admin,gotrue,syslog,pglisten)

