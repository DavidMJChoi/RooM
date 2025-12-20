package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// PageData 包含页面共用数据
type PageData struct {
	CurrentPage string
}

// Post 表示一篇博客文章
type Post struct {
	Title       string
	Content     string
	Date        time.Time
	Tags        []string
	Summary     string
	Slug        string
	PageData    // 嵌入页面数据
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




func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/posts/", postsHandler)
	http.HandleFunc("/about/", aboutHandler)

	fmt.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// 如果是根路径，显示首页
	// if r.URL.Path != "/" {
	// 	http.NotFound(w, r)
	// 	return
	// }
	
	// 解析所有模板文件
	tmpl, err := template.ParseFiles("templates/index.html", "templates/nav.html", "templates/footer.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// 传递页面数据
	data := PageData{
		CurrentPage: "home",
	}
	
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// postsHandler 处理文章列表页面
func postsHandler(w http.ResponseWriter, r *http.Request) {
	// 清理路径
	path := strings.TrimPrefix(r.URL.Path, "/posts/")
	
	// 如果路径为空或只是 "/"，显示文章列表
	if path == "" || path == "/" {
		posts, err := loadPosts()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		// 解析所有模板文件
		tmpl, err := template.ParseFiles("templates/posts.html", "templates/nav.html", "templates/footer.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		// 创建包含页面数据和文章列表的数据结构
		data := struct {
			PageData
			Posts []Post
		}{
			PageData: PageData{CurrentPage: "posts"},
			Posts:    posts,
		}
		
		err = tmpl.Execute(w, data)
		if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
		return
	}
	
	// 否则显示单篇文章
	postHandler(w, r)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	// 提取文章 slug
	slug := strings.TrimPrefix(r.URL.Path, "/posts/")
	slug = strings.TrimSuffix(slug, "/")
	
	// 使用 loadPost 函数加载文章
	post, err := loadPost(slug)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	
	// 创建模板并添加自定义函数，解析所有模板文件
	tmpl := template.Must(template.New("post.html").Funcs(template.FuncMap{
		"toHTML": simpleMarkdownToHTML,
	}).ParseFiles("templates/post.html", "templates/nav.html", "templates/footer.html"))
	
	err = tmpl.Execute(w, post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/about.html", "templates/nav.html", "templates/footer.html"))
	
	data := PageData{
		CurrentPage: "about",
	}
	
	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}