package executions

type PathRepository interface {
	Current(string, string) (string, error)
}
