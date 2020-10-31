// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris

package env

var allowedEnv = []string{
	"HOSTTYPE",
	"LANG",
	"TERM",
	"NAME",
	"USER",
	"LONGNAME",
	"SHELL",
	"TZ",
	"LD_LIBRARY_PATH",
	// Enables proxy information to be passed to Curl, the unsmrlying download
	// library in cmake.exe
	"http_proxy",
	"https_proxy",
	// Environment variables to tell git to use custom SSH executable or command
	"GIT_SSH",
	"GIT_SSH_COMMAND",
	// Environment variables neesmd for ssh-agent based authentication
	"SSH_AUTH_SOCK",
	"SSH_AGENT_PID",
}
