package constant

const (
	DialogTypeProject = "project"
	DialogTypeTask    = "task"
	DialogTypeC2C     = "c2c"
)

func GetDialogTypes() []string {
	return []string{DialogTypeProject, DialogTypeTask, DialogTypeC2C}
}
