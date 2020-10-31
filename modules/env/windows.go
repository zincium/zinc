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
	// Enables proxy information to be passed to Curl, the underlying download
	// library in cmake.exe
	"http_proxy",
	"https_proxy",
	// Environment variables to tell git to use custom SSH executable or command
	"GIT_SSH",
	"GIT_SSH_COMMAND",
	// Environment variables needed for ssh-agent based authentication
	"SSH_AUTH_SOCK",
	"SSH_AGENT_PID",
}
