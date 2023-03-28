package streams

type category struct {
	ID       string
	Name     string
	ParentID string
}

type CategoryRepository interface {
	// Save category and subcategories
	Save(string, []string) (string, error)

	// List all categories
	List(string) ([]Category, error)
}
