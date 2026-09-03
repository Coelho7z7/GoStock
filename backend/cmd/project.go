package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// prepararDiretorioProjeto garante que o processo esteja rodando com o
// diretório de trabalho na raiz do projeto (onde ficam as pastas
// frontend/ e backend/), independente de onde o binário foi chamado.
func prepararDiretorioProjeto() error {
	diretorio, err := os.Getwd()
	if err != nil {
		return err
	}

	for {
		if _, err := os.Stat(filepath.Join(diretorio, "frontend", "html", "index.html")); err == nil {
			return os.Chdir(diretorio)
		}
		projetoAninhado := filepath.Join(diretorio, "GoStock")
		if _, err := os.Stat(filepath.Join(projetoAninhado, "frontend", "html", "index.html")); err == nil {
			return os.Chdir(projetoAninhado)
		}

		parent := filepath.Dir(diretorio)
		if parent == diretorio {
			return fmt.Errorf("frontend/html/index.html não encontrado")
		}
		diretorio = parent
	}
}
