package main

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type Post struct {
	Title       string
	Content     string
	Date        time.Time
	Tags        []string
	Summary     string
	Slug        string
	PageData    // 嵌入页面数据
}


// parseFrontMatter 解析文章的 front matter（标题、日期、标签等）
func parseFrontMatter(content string) (map[string]string, string) {
	frontMatter := make(map[string]string)
	body := content
	
	// 检查是否有 front matter (以 --- 开头和结尾)
	if strings.HasPrefix(content, "---") {
		parts := strings.SplitN(content, "---", 3)
		if len(parts) >= 3 {
			frontMatterContent := parts[1]
			body = parts[2]
			
			// 解析 front matter
			lines := strings.Split(frontMatterContent, "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line == "" {
					continue
				}
				if colonIndex := strings.Index(line, ":"); colonIndex > 0 {
					key := strings.TrimSpace(line[:colonIndex])
					value := strings.TrimSpace(line[colonIndex+1:])
					frontMatter[key] = strings.Trim(value, "\"")
				}
			}
		}
	}
	
	return frontMatter, strings.TrimSpace(body)
}

// parseDate 解析日期字符串
func parseDate(dateStr string) time.Time {
	if dateStr == "" {
		return time.Now()
	}
	
	// 尝试多种日期格式
	formats := []string{
		"2006-01-02",
		"2006/01/02",
		"2006-01-02 15:04:05",
		"2006/01/02 15:04:05",
	}
	
	for _, format := range formats {
		if date, err := time.Parse(format, dateStr); err == nil {
			return date
		}
	}
	
	return time.Now()
}

// parseTags 解析标签字符串
func parseTags(tagsStr string) []string {
	if tagsStr == "" {
		return []string{}
	}
	
	tagsStr = tagsStr[2:len(tagsStr)-2]

	// 支持逗号分隔或空格分隔
	tags := strings.FieldsFunc(tagsStr, func(c rune) bool {
		return c == ',' || c == ' '
	})
	
	var cleanTags []string
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		tag = strings.Trim(tag, "\"")
		if tag != "" {
			cleanTags = append(cleanTags, tag)
		}
	}
	
	return cleanTags
}

// titleCase 将字符串转换为标题格式
func titleCase(s string) string {
	if s == "" {
		return s
	}
	words := strings.Fields(s)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		}
	}
	return strings.Join(words, " ")
}

// stripHTML 移除 HTML 标签
func stripHTML(s string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	return re.ReplaceAllString(s, "")
}

// simpleMarkdownToHTML 简单的 Markdown 到 HTML 转换
func simpleMarkdownToHTML(s string) template.HTML {

	// 转换标题
	s = regexp.MustCompile(`(?m)^###\s+(.+)\n`).ReplaceAllString(s, "<h3 class='text-xl font-semibold mb-4'>$1</h3>")
	s = regexp.MustCompile(`(?m)^##\s+(.+)\n`).ReplaceAllString(s, "<h2 class='text-2xl font-bold mb-6'>$1</h2>")
	s = regexp.MustCompile(`(?m)^#\s+(.+)\n`).ReplaceAllString(s, "<h1 class='text-3xl font-bold mb-8'>$1</h1>")
	
	// 转换粗体文本
	s = regexp.MustCompile(`\*\*(.+?)\*\*`).ReplaceAllString(s, "<strong>$1</strong>")
	s = regexp.MustCompile(`\*(.+?)\*`).ReplaceAllString(s, "<em>$1</em>")
	
	// 转换代码块
	s = regexp.MustCompile("```(\\w+)?\\n([\\s\\S]*?)```").ReplaceAllString(s, "<pre class='bg-gray-100 dark:bg-gray-800 p-4 rounded-lg mb-6 overflow-x-auto'><code>$2</code></pre>")
	// convert inline code
	s = regexp.MustCompile("`(.*?)`").ReplaceAllString(s, "<code class='bg-gray-100 dark:bg-gray-800 px-1 py-0.5 rounded text-sm'>$1</code>")
	
	// 转换列表
	lines := strings.Split(s, "\n")
	var result []string
	inList := false
	
	for _, v := range lines {
		fmt.Println(v)
	}

	for i := 0; i < len(lines); i++ {

		if strings.HasPrefix(lines[i], "<pre class='bg-gray-100 dark:bg-gray-800 p-4 rounded-lg mb-6 overflow-x-auto'><code>") {
			for !strings.HasPrefix(lines[i], "</code></pre>") {
				result = append(result, lines[i])
				i++
			}
			result = append(result, lines[i])
		} else {
			if strings.HasPrefix(lines[i], "- ") {
			if !inList {
				result = append(result, "<ul class='list-disc list-inside mb-6 space-y-2'>")
				inList = true
			}
			content := strings.TrimPrefix(lines[i], "- ")
			result = append(result, "<li>"+content+"</li>")
		} else {
			if inList {
				result = append(result, "</ul>")
				inList = false
			}
			if lines[i] != "" {
				// result = append(result,lines[i])//, "<p class='mb-4'>"+lines[i]+"</p>")
				hasTag := strings.HasPrefix(lines[i], "<h") || 
                 strings.HasPrefix(lines[i], "<li>") || 
                 strings.HasPrefix(lines[i], "<hr")
        
				if hasTag {
					result = append(result, lines[i])
				} else {
					result = append(result, fmt.Sprintf("<p class='mb-4'>%s</p>", lines[i]))
				}
			}
		}
		}
	}
	
	if inList {
		result = append(result, "</ul>")
	}
	
	for _, v := range result {
		fmt.Println(v)
	}

	return template.HTML(strings.Join(result, "\n"))
}

// loadPost 加载单篇文章
func loadPost(slug string) (*Post, error) {
	// 查找对应的文章文件
	postPath := filepath.Join("posts", slug+".md")
	if _, err := os.Stat(postPath); os.IsNotExist(err) {
		return nil, err
	}
	
	content, err := os.ReadFile(postPath)
	if err != nil {
		return nil, err
	}
	
	contentStr := string(content)
	frontMatter, body := parseFrontMatter(contentStr)
	
	post := &Post{
		Title:   frontMatter["title"],
		Date:    parseDate(frontMatter["date"]),
		Tags:    parseTags(frontMatter["tags"]),
		Content: body,
		Slug:    slug,
		PageData: PageData{CurrentPage: "posts"}, // 单篇文章也在posts页面
	}
	
	// 如果没有标题，使用文件名
	if post.Title == "" {
		post.Title = titleCase(strings.ReplaceAll(slug, "-", " "))
	}
	
	// 生成摘要（取前100个字符）
	plainText := stripHTML(body)
	if len(plainText) > 100 {
		post.Summary = plainText[:100] + "..."
	} else {
		post.Summary = plainText
	}
	
	return post, nil
}

// loadPosts 加载所有文章
func loadPosts() ([]Post, error) {
	postsDir := "posts"
	var posts []Post
	
	// 遍历 posts 目录
	err := filepath.Walk(postsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// 只处理 .md 文件，排除目录
		if !info.IsDir() && strings.HasSuffix(path, ".md") {
			// 提取文件名作为 slug
			fileName := strings.TrimSuffix(filepath.Base(path), ".md")
			
			fmt.Println(fileName)

			// 使用 loadPost 函数加载单篇文章
			post, err := loadPost(fileName)
			if err != nil {
				return err
			}
			
			posts = append(posts, *post)
		}
		
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	
	// 按日期排序（最新的在前）
	for i := 0; i < len(posts)-1; i++ {
		for j := i + 1; j < len(posts); j++ {
			if posts[i].Date.Before(posts[j].Date) {
				posts[i], posts[j] = posts[j], posts[i]
			}
		}
	}
	
	return posts, nil
}

