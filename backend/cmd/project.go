package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// prepareProjectDirectory garante que o processo esteja rodando com o
// diretório de trabalho na raiz do projeto (onde ficam as pastas
// frontend/ e backend/), independente de onde o binário foi chamado.
func prepareProjectDirectory() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "frontend", "html", "index.html")); err == nil {
			return os.Chdir(dir)
		}
		nestedProject := filepath.Join(dir, "GoStock")
		if _, err := os.Stat(filepath.Join(nestedProject, "frontend", "html", "index.html")); err == nil {
			return os.Chdir(nestedProject)
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return fmt.Errorf("frontend/html/index.html não encontrado")
		}
		dir = parent
	}
}
