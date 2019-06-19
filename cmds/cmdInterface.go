package cmds

type Cmds interface {
	Execute() bool
	SubDispatch() bool
}
