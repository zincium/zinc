//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris
// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris

package env

var allowedEnv = []string{
	"HOME"
	"PATH"
	"LD_LIBRARY_PATH",
	"TZ",
	"HOSTTYPE",
	"LANG",
	"TERM",
	"NAME",
	"USER",
	"LONGNAME",
	"SHELL",

	// ENABLE trace
	"GIT_TRACE",
	"GIT_TRACE_PACK_ACCESS",
	"GIT_TRACE_PACKET",
	"GIT_TRACE_PERFORMANCE",
	"GIT_TRACE_SETUP",

	// Git HTTP proxy settings: https://git-scm.com/docs/git-config#git-config-httpproxy
	"all_proxy",
	"http_proxy",
	"HTTP_PROXY",
	"https_proxy",
	"HTTPS_PROXY",
	// libcurl settings: https://curl.haxx.se/libcurl/c/CURLOPT_NOPROXY.html
	"no_proxy",
	"NO_PROXY",
	// Environment variables to tell git to use custom SSH executable or command
	"GIT_SSH",
	"GIT_SSH_COMMAND",
	// Environment variables neesmd for ssh-agent based authentication
	"SSH_AUTH_SOCK",
	"SSH_AGENT_PID",
}
