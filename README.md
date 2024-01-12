# Features
- Username to IPs
- IP to usernames
- Find alts (based on IP)
- See your players total playtimes

# How to install
- Make sure you have [Go](https://go.dev) installed, and it's in your PATH environment variable
- Change directory into your ReIndev server. This program looks for `server.log` in the parent folder to read from
- `git clone https://github.com/kivattt/reindev-dashboard`
- `cd reindev-dashboard`
- `go build main.go readlog.go`
- Set secure passwords for admin and staff in `config.txt`
- Now launch it with `./main`
- Reach it at `localhost:8080`

# Important info
This is an HTTP server. Everything can be intercepted by your ISP, local starbucks squatter or big brother \
So use a unique password, and not one you use elsewhere. \
Alternatively, use a reverse proxy like [NGINX](https://docs.nginx.com/nginx/admin-guide/web-server/reverse-proxy/) or [caddy](https://caddyserver.com/)

The admin user has full access to the tools \
The staff user only has access to the `find-alts` and `playtimes` tool, because they should not leak any IP addresses \
Give your staff team access to the staff user, to not leak any IP addresses from the `server.log`
