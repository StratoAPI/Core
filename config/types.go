package config

type Config struct {
	Host string `json:"host"`
	Port uint32 `json:"port"`

	Database Database `json:"database"`
}

type Database struct {
	Type     string   `json:"type"`
	Flatfile Flatfile `json:"flatfile"`
}

type Flatfile struct {
	Location string `json:"location"`
}
