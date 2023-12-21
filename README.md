# Features
- Username to IPs
- IP to usernames
- Find alts (based on IP)

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
The admin user has full access to the tools. \
The staff user only has access to the `find-alts` tool, because it should not leak any IP addresses. \
Give your staff team access to the staff user, to not leak any IP addresses from the `server.log`
