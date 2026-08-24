const senha = document.getElementById("senha");
const mostrarSenha = document.getElementById("mostrarSenha");
const iconeOlho = document.getElementById("iconeOlho");

mostrarSenha.addEventListener("click", function () {

    if (senha.type === "password") {

        senha.type = "text";

        iconeOlho.innerHTML = `
            <path d="M3 3l18 18"/>
            <path d="M10.6 10.6a2 2 0 0 0 2.8 2.8"/>
            <path d="M9.9 4.2A10.8 10.8 0 0 1 12 4c6.5 0 10 8 10 8a17.2 17.2 0 0 1-3.1 4.2"/>
            <path d="M6.6 6.6C3.8 8.5 2 12 2 12s3.5 8 10 8a9.8 9.8 0 0 0 4.2-.9"/>
        `;

        mostrarSenha.setAttribute("aria-label", "Ocultar senha");

    } else {

        senha.type = "password";

        iconeOlho.innerHTML = `
            <path d="M2 12s3.5-7 10-7 10 7 10 7-3.5 7-10 7S2 12 2 12z"/>
            <circle cx="12" cy="12" r="3"/>
        `;

        mostrarSenha.setAttribute("aria-label", "Mostrar senha");
    }
});