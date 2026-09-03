package services

import (
	"errors"
	"strings"

	database "gostock/backend/Database"
	"gostock/backend/models"
	"gostock/backend/utils"

	"golang.org/x/crypto/bcrypt"
)

// ListarUsuariosPaginados retorna os usuários cadastrados, filtrados por
// nome/email e paginados, além do total de registros encontrados.
func ListarUsuariosPaginados(busca string, pagina, porPagina int) ([]models.Usuario, int, error) {
	busca = strings.TrimSpace(busca)
	filtro := "%" + busca + "%"
	offset := (pagina - 1) * porPagina

	var total int
	if err := database.DB.QueryRow(`
		SELECT COUNT(*) FROM usuarios
		WHERE nome LIKE ? OR email LIKE ?
	`, filtro, filtro).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := database.DB.Query(`
		SELECT id, nome, email, role
		FROM usuarios
		WHERE nome LIKE ? OR email LIKE ?
		ORDER BY nome
		LIMIT ? OFFSET ?
	`, filtro, filtro, porPagina, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	usuarios := make([]models.Usuario, 0, porPagina)
	for rows.Next() {
		var usuario models.Usuario
		if err := rows.Scan(&usuario.ID, &usuario.Nome, &usuario.Email, &usuario.Role); err != nil {
			return nil, 0, err
		}
		usuarios = append(usuarios, usuario)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return usuarios, total, nil
}

// CriarUsuarioWeb cadastra um novo usuário a partir do painel de
// administração, com a permissão escolhida pelo administrador.
func CriarUsuarioWeb(nome, email, senha, role string) error {
	nome = strings.TrimSpace(nome)
	email = strings.ToLower(strings.TrimSpace(email))
	role = strings.TrimSpace(role)

	if !utils.ValidarNome(nome) {
		return errors.New("Informe o nome do usuário.")
	}
	if !utils.ValidarEmail(email) {
		return errors.New("Email inválido. Use um endereço @gmail.com.")
	}
	// O endereço do CEO é reservado e não pode ser reutilizado por outra conta.
	if strings.EqualFold(email, "admin@gmail.com") {
		return errors.New("O email admin@gmail.com é reservado ao CEO e não pode ser cadastrado.")
	}
	if !utils.ValidarSenha(senha) {
		return errors.New("A senha deve ter no mínimo 6 caracteres e 1 caractere especial.")
	}
	if strings.EqualFold(role, "ceo") {
		return errors.New("O cargo CEO é reservado exclusivamente para admin@gmail.com.")
	}
	if role != "admin" && role != "basico" {
		return errors.New("Permissão inválida.")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("Erro ao proteger a senha.")
	}

	if _, err := database.DB.Exec(`
		INSERT INTO usuarios (nome, email, senha, role)
		VALUES (?, ?, ?, ?)
	`, nome, email, string(hash), role); err != nil {
		return errors.New("Esse email já está cadastrado.")
	}

	return nil
}

// AtualizarRoleUsuarioWeb altera a permissão de um usuário.
// O cargo CEO é exclusivo do admin@gmail.com e nunca pode ser alterado.
func AtualizarRoleUsuarioWeb(alvoID int, novaRole string) error {
	novaRole = strings.TrimSpace(novaRole)
	if novaRole != "admin" && novaRole != "basico" {
		return errors.New("Permissão inválida.")
	}

	alvo, err := BuscarUsuarioPorID(alvoID)
	if err != nil {
		return errors.New("Usuário não encontrado.")
	}
	if strings.EqualFold(strings.TrimSpace(alvo.Role), "ceo") ||
		strings.EqualFold(strings.TrimSpace(alvo.Email), "admin@gmail.com") {
		return errors.New("O cargo CEO é protegido e não pode ser alterado.")
	}

	if _, err := database.DB.Exec(`UPDATE usuarios SET role = ? WHERE id = ?`, novaRole, alvoID); err != nil {
		return errors.New("Erro ao atualizar a permissão.")
	}

	return nil
}

// RemoverUsuarioWeb apaga um usuário e suas sessões ativas.
// O CEO nunca pode ser removido por este caminho.
func RemoverUsuarioWeb(alvoID int) error {
	alvo, err := BuscarUsuarioPorID(alvoID)
	if err != nil {
		return errors.New("Usuário não encontrado.")
	}
	if strings.EqualFold(strings.TrimSpace(alvo.Role), "ceo") ||
		strings.EqualFold(strings.TrimSpace(alvo.Email), "admin@gmail.com") {
		return errors.New("O CEO é protegido e não pode ser removido.")
	}

	if _, err := database.DB.Exec(`DELETE FROM sessoes WHERE usuario_id = ?`, alvoID); err != nil {
		return errors.New("Erro ao remover usuário.")
	}

	if _, err := database.DB.Exec(`DELETE FROM usuarios WHERE id = ?`, alvoID); err != nil {
		return errors.New("Erro ao remover usuário.")
	}

	return nil
}
