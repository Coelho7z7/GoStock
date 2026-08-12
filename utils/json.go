package utils

import (
	"encoding/json"
	"errors"
	"os"

	"gostock/models"
)

func SalvarProdutos(produtos []models.Produto) error {
	dados, err := json.MarshalIndent(produtos, "", "    ")

	if err != nil {
		return err
	}

	return os.WriteFile("data/produtos.json", dados, 0644)
}

func CarregarProdutos() ([]models.Produto, error) {
	dados, err := os.ReadFile("data/produtos.json")

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []models.Produto{}, nil
		}

		return nil, err
	}

	var produtos []models.Produto

	err = json.Unmarshal(dados, &produtos)

	if err != nil {
		return nil, err
	}

	return produtos, nil
}
