document.addEventListener("DOMContentLoaded", function () {
    const senha = document.getElementById("senha");
    const mostrarSenha = document.getElementById("mostrarSenha");
    const iconeOlho = document.getElementById("iconeOlho");

    if (senha && mostrarSenha && iconeOlho) {
        mostrarSenha.addEventListener("click", function () {
        const senhaVisivel = senha.type === "text";

        senha.type = senhaVisivel ? "password" : "text";

        if (senhaVisivel) {
            iconeOlho.innerHTML = `
                <path d="M2 12s3.5-7 10-7 10 7 10 7-3.5 7-10 7S2 12 2 12Z"></path>
                <circle cx="12" cy="12" r="3"></circle>
            `;

            mostrarSenha.setAttribute("aria-label", "Mostrar senha");
            mostrarSenha.setAttribute("title", "Mostrar senha");
            return;
        }

        iconeOlho.innerHTML = `
            <path d="M3 3l18 18"></path>
            <path d="M10.6 10.6a2 2 0 0 0 2.8 2.8"></path>
            <path d="M9.9 4.2A10.8 10.8 0 0 1 12 4c6.5 0 10 8 10 8a17.2 17.2 0 0 1-3.1 4.2"></path>
            <path d="M6.6 6.6C3.8 8.5 2 12 2 12s3.5 8 10 8a9.8 9.8 0 0 0 4.2-.9"></path>
        `;

        mostrarSenha.setAttribute("aria-label", "Ocultar senha");
        mostrarSenha.setAttribute("title", "Ocultar senha");
        });
    }

    const modalCadastro = document.getElementById("modal-cadastro-produto");
    const abrirCadastro = document.getElementById("abrir-cadastro-produto");
    const fecharCadastro = document.getElementById("fechar-cadastro-produto");

    if (!modalCadastro || !abrirCadastro || !fecharCadastro) {
        return;
    }

    const fecharModal = function () {
        modalCadastro.hidden = true;
    };

    abrirCadastro.addEventListener("click", function () {
        modalCadastro.hidden = false;
        modalCadastro.querySelector("input")?.focus();
    });

    fecharCadastro.addEventListener("click", fecharModal);

    modalCadastro.addEventListener("click", function (evento) {
        if (evento.target === modalCadastro) {
            fecharModal();
        }
    });

    document.addEventListener("keydown", function (evento) {
        if (evento.key === "Escape" && !modalCadastro.hidden) {
            fecharModal();
        }
    });
});
