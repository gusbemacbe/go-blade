// factory.go
package blade

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Blade is the core struct for the template engine.
type Blade struct {
	viewPaths []string
	cachePath string
	mu        sync.Mutex
}

// New creates a new Blade engine instance.
func New(viewPaths []string, cachePath string) *Blade {
	return &Blade{
		viewPaths: viewPaths,
		cachePath: cachePath,
	}
}

// Run renders a template with the provided data.
func (b *Blade) Run(templateName string, data map[string]interface{}) (string, error) {
	// 1. Finding the template file.
	viewFile, err := b.findView(templateName)
	if err != nil {
		return "", err
	}

	// 2. Compiling the Blade template to a Go template.
	compiledPath, err := b.compile(viewFile)
	if err != nil {
		return "", err
	}

	// 3. Executing the compiled Go template with the data.
	tpl, err := template.ParseFiles(compiledPath)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// findView searches for the view file in the configured view paths.
func (b *Blade) findView(name string) (string, error) {
	name = strings.Replace(name, ".", "/", -1)
	for _, viewPath := range b.viewPaths {
		// The original library supported multiple extensions; we will simplify to .blade for now.
		viewFile := filepath.Join(viewPath, name+".blade")
		if _, err := os.Stat(viewFile); err == nil {
			return viewFile, nil
		}
	}
	return "", fmt.Errorf("view '%s' not found", name)
}

// compile translates the Blade syntax into Go's html/template syntax and caches it.
func (b *Blade) compile(path string) (string, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	cacheFile := filepath.Join(b.cachePath, fmt.Sprintf("%x", md5.Sum([]byte(path))))
	fileInfo, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	cacheInfo, err := os.Stat(cacheFile)
	if err == nil && cacheInfo.ModTime().After(fileInfo.ModTime()) {
		return cacheFile, nil // Returning the cached file if it's still fresh.
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	compiledContent := compileContent(string(content))

	if err := os.WriteFile(cacheFile, []byte(compiledContent), 0644); err != nil {
		return "", err
	}

	return cacheFile, nil
}
