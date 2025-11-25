package models

type ChapterData struct {
	ID     int    `yaml:"id" json:"id"`
	Title  string `yaml:"title" json:"title"`
	Folder string `yaml:"folder" json:"folder"`
}
type CourseData struct {
	ID           int64       `db:"id" json:"id"`
	Slug         string      `yaml:"slug" json:"slug" db:"slug"`
	Name         string      `yaml:"name" json:"name" db:"name"`
	Description  string      `yaml:"description" json:"description" db:"description"`
	Tags         StringSlice `yaml:"tags" json:"tags" db:"tags"`
	Folder       string      `yaml:"folder" json:"folder" db:"folder"`
	ChapterCount int         `yaml:"chapter_count" json:"chapter_count" db:"chapter_count"`
}
type SectionData struct {
	Type string `yaml:"type" json:"type"`
	Text string `yaml:"text" json:"text"`
	Quiz []int  `yaml:"quiz" json:"quiz"`
}
type QuizData struct {
	ID      int      `yaml:"id" json:"id"`
	Type    string   `yaml:"type" json:"type"`
	Text    string   `yaml:"text" json:"text"`
	Options []string `yaml:"options,omitempty" json:"options,omitempty"`
	Answer  []string `yaml:"answer" json:"answer"`
}
