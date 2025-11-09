package models

type ChapterData struct {
	ID     int    `yaml:"id" json:"id"`
	Title  string `yaml:"title" json:"title"`
	Folder string `yaml:"folder" json:"folder"`
}
type CourseData struct {
	Name         string   `yaml:"name" json:"name"`
	Description  string   `yaml:"description" json:"description"`
	Slug         string   `yaml:"slug" json:"slug"`
	Tags         []string `yaml:"tags" json:"tags"`
	Folder       string   `yaml:"folder" json:"folder"`
	ChapterCount int      `yaml:"chapter_count" json:"chapter_count"`
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
