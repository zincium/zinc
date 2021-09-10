//go:build windows
// +build windows

package env

var allowedEnv = []string{
	"ALLUSERSPROFILE",
	"APPDATA",
	"CommonProgramFiles",
	"CommonProgramFiles(x86)",
	"CommonProgramW6432",
	"COMPUTERNAME",
	"ComSpec",
	"HOMEDRIVE",
	"HOMEPATH",
	"LOCALAPPDATA",
	"LOGONSERVER",
	"NUMBER_OF_PROCESSORS",
	"OS",
	"PATHEXT",
	"PROCESSOR_ARCHITECTURE",
	"PROCESSOR_ARCHITEW6432",
	"PROCESSOR_IDENTIFIER",
	"PROCESSOR_LEVEL",
	"PROCESSOR_REVISION",
	"ProgramData",
	"ProgramFiles",
	"ProgramFiles(x86)",
	"ProgramW6432",
	"PROMPT",
	"PATH",
	"PSModulePath",
	"PUBLIC",
	"SystemDrive",
	"SystemRoot",
	"TEMP",
	"TMP",
	"USERDNSDOMAIN",
	"USERDOMAIN",
	"USERDOMAIN_ROAMINGPROFILE",
	"USERNAME",
	"USERPROFILE",
	"windir",
	// Windows Terminal
	"SESSIONNAME",
	"WT_SESSION",
	"WSLENV",
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
