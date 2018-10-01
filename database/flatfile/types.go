package flatfile

type Flatfile struct {
	Location string
	Data     map[string][]map[string]interface{}
}
