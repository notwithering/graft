package preset

type ProjectConfig struct {
	Root     string
	Commands map[string]*CommandSpec
}
