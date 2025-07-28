package blade

import (
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gusbemacbe/go-blade/lexers"
)

// The Compiler struct is responsible for transpiling Blade files into Go template files.
type Compiler struct {
	compiledFilePath string
	lexers           []Lexer
}

// The Compile method checks if a Blade template is expired and recompiles it if necessary.
func (compiler *Compiler) Compile(file string) ([]byte, error) {
	// Checking if the compiled template is expired.
	isExpired, err := compiler.IsExpired(file)
	if err != nil {
		return nil, err
	}

	if isExpired {
		// Reading the raw Blade template content.
		bytes, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}

		// Applying each lexer to transpile the content.
		for _, lexer := range compiler.lexers {
			bytes = lexer.Parse(bytes)
		}

		// Writing the compiled Go template to the cache.
		err = compiler.WriteCompiled(compiler.CompiledPath(file), bytes)
		return bytes, err
	}

	// Loading the fresh, compiled template from the cache.
	compiled := compiler.CompiledPath(file)
	return os.ReadFile(compiled)
}

// The CompiledPath method generates a unique, hashed path for the cached template file.
func (compiler *Compiler) CompiledPath(file string) string {
	// Instantiating the hasher within the method for safety.
	hasher := sha1.New()
	hasher.Write([]byte(file))
	hashed := fmt.Sprintf("%x", hasher.Sum(nil))

	return filepath.Join(compiler.compiledFilePath, hashed+".blade.html")
}

// The IsExpired method determines if a compiled template is out of date.
func (compiler *Compiler) IsExpired(file string) (bool, error) {
	compiled := compiler.CompiledPath(file)

	// Checking if a compiled version exists.
	compiledInfo, err := os.Stat(compiled)

	if os.IsNotExist(err) {
		return true, nil // The file is "expired" if it hasn't been compiled yet.
	}
	if err != nil {
		return true, err // Returning other errors (e.g., permissions).
	}

	// Getting the file info for the original Blade template.
	fileInfo, err := os.Stat(file)
	if err != nil {
		return true, err
	}

	// The template is expired if the original file's modification time is after the compiled file's.
	return fileInfo.ModTime().After(compiledInfo.ModTime()), nil
}

// The WriteCompiled method writes the transpiled content to the cache directory.
func (compiler *Compiler) WriteCompiled(filename string, contents []byte) error {
	// Using the modern `os.WriteFile` function.
	return os.WriteFile(filename, contents, 0644)
}

// The applyLexers method initializes the sequence of transpilation steps.
func (compiler *Compiler) applyLexers() {
	compiler.lexers = []Lexer{
		new(lexers.Echo),
		new(lexers.If),
		new(lexers.Else),
		new(lexers.EndIf),
	}
}

// The NewCompiler function creates and initializes a new Compiler instance.
func NewCompiler() *Compiler {
	compiler := new(Compiler)
	compiler.applyLexers()
	return compiler
}
