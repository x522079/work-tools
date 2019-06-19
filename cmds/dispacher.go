package cmds

import (
	"math/rand"
)

type Commands struct {
	Cmd         string
	SubCmd      string
	Path        string
	RemotePath  string
	RemoteAlias string
}

type Dispacher struct {
	// 随机全局ID
	Uid int
	Cmd Commands
}

func NewDispacher(cmds Commands) {
	Dispatcher = &Dispacher{Uid: rand.Int(), Cmd: cmds}
}

func (d *Dispacher) Dispach() {
	switch d.Cmd.Cmd {
	case "scp":
		(&Scp{}).SubDispatch()
		break
	default:
		panic("暂不支持的命令")
	}
}
