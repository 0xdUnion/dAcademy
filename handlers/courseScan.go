package handlers

import (
	"dAcademy/database"
	"dAcademy/models"
	"dAcademy/utils"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"

	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	courseScanMutex  sync.Mutex
	isCourseScanning bool
)

func CourseScanHandler(c *gin.Context) {
	courseScanMutex.Lock()
	if isCourseScanning {
		courseScanMutex.Unlock()
		c.JSON(http.StatusConflict, gin.H{"message": "A scan job is already in progress."})
		return
	}

	isCourseScanning = true
	courseScanMutex.Unlock()

	go func() {
		defer func() {
			courseScanMutex.Lock()
			isCourseScanning = false
			courseScanMutex.Unlock()
			log.Println("Courses scan process finished.")
		}()

		fmt.Println("Courses scan process starting.")
		if err := scanAll("./courses"); err != nil {
			log.Printf("❌ Scan failed: %v", err)
		} else {
			log.Println("✅ All courses indexed successfully.")
		}
	}()

	c.JSON(http.StatusAccepted, gin.H{"message": "Scan job has been accepted."})
}

// Scan all(courses and their chapters)
func scanAll(root string) error {
	var courses []models.CourseData

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if !d.IsDir() {
			return nil
		}

		folderName := filepath.Base(path)
		courseYaml := filepath.Join(path, "course.yaml")
		chaptersYaml := filepath.Join(path, "_chapters.yaml")

		if _, err := os.Stat(courseYaml); err == nil {
			fmt.Println("📘 構建章節索引:", path)

			chapters, err := buildChapters(path)
			if err != nil {
				return fmt.Errorf("構建章節索引失敗 %s: %v", path, err)
			}

			if len(chapters) > 0 {
				if err := utils.SaveYAML(chaptersYaml, chapters); err != nil {
					return fmt.Errorf("寫入 _chapters.yaml 失敗: %v", err)
				}
			}

			data, err := os.ReadFile(courseYaml)
			if err != nil {
				return fmt.Errorf("讀取 course.yaml 失敗 %s: %v", path, err)
			}

			var course models.CourseData
			if err := yaml.Unmarshal(data, &course); err != nil {
				return fmt.Errorf("解析 course.yaml 失敗 %s: %v", path, err)
			}
			course.ChapterCount = len(chapters)
			course.Folder = folderName

			db, err := database.Run()
			if err != nil {
				log.Fatal(err)
			}

			_, err = db.NamedExec(`
        INSERT INTO courses (slug, name, description, tags, folder, chapter_count)
        VALUES (:slug, :name, :description, :tags, :folder, :chapter_count)
        ON CONFLICT(slug) DO UPDATE SET
            slug = excluded.slug,
            name = excluded.name,
            description = excluded.description,
            tags = excluded.tags,
            folder = excluded.folder,
            chapter_count = excluded.chapter_count
    `, course)
			if err != nil {
				log.Fatalf("insert error: %v", err)
			}

		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("掃描課程目錄失敗: %v", err)
	}

	fmt.Println("✅ Scan finished, total: ", len(courses), " courses")
	return nil
}

// Scan chapter folders
func buildChapters(coursePath string) ([]models.ChapterData, error) {
	files, err := os.ReadDir(coursePath)
	if err != nil {
		return nil, err
	}

	var chapters []models.ChapterData
	for _, f := range files {
		if f.IsDir() {
			// Skip hidden folders
			if f.Name()[0] == '.' {
				continue
			}
			id, title, ok := parseChapterFolder(f.Name())
			if !ok {
				log.Fatalf("❌ 無效章節資料夾: %s", f.Name())
			}
			chapters = append(chapters, models.ChapterData{
				ID:     id,
				Title:  title,
				Folder: f.Name(),
			})
		}
	}

	// Sort by ID
	sort.Slice(chapters, func(i, j int) bool { return chapters[i].ID < chapters[j].ID })
	return chapters, nil
}

// Parse chapter folder name, like 1-Hello
func parseChapterFolder(name string) (int, string, bool) {
	parts := strings.SplitN(name, "-", 2)
	if len(parts) != 2 {
		return 0, "", false
	}
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", false
	}
	return id, parts[1], true
}
