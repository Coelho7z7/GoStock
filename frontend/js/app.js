// Toggle de mostrar/ocultar senha (Login)

// ---------------------------------------------------------------------
// Navegação mobile: sidebar sanduíche, overlay, ESC e fechamento ao navegar.
// Também transforma cabeçalhos de tabela em data-labels para a visualização
// em cards no celular.
// ---------------------------------------------------------------------
function initMobileNavigation() {
    const menuButton = document.querySelector('.mobile-menu-button');
    const sidebar = document.querySelector('.sidebar');
    const overlay = document.querySelector('.mobile-sidebar-overlay');
    if (!menuButton || !sidebar || !overlay) return;

    const setOpen = function (open) {
        document.body.classList.toggle('mobile-menu-open', open);
        menuButton.setAttribute('aria-expanded', open ? 'true' : 'false');
        menuButton.setAttribute('aria-label', open ? 'Fechar menu' : 'Abrir menu');
    };

    menuButton.addEventListener('click', function () {
        setOpen(!document.body.classList.contains('mobile-menu-open'));
    });

    overlay.addEventListener('click', function () {
        setOpen(false);
    });

    sidebar.querySelectorAll('a').forEach(function (link) {
        link.addEventListener('click', function () {
            setOpen(false);
        });
    });

    document.addEventListener('keydown', function (event) {
        if (event.key === 'Escape') setOpen(false);
    });

    window.addEventListener('resize', function () {
        if (window.innerWidth > 760) setOpen(false);
    });
}

function initResponsiveTableLabels() {
    document.querySelectorAll('.content table').forEach(function (table) {
        const headers = Array.from(table.querySelectorAll('thead th')).map(function (th) {
            return th.textContent.trim();
        });
        if (!headers.length) return;

        table.querySelectorAll('tbody tr').forEach(function (row) {
            Array.from(row.children).forEach(function (cell, index) {
                if (cell.hasAttribute('colspan')) return;
                if (headers[index]) cell.setAttribute('data-label', headers[index]);
            });
        });
    });
}

document.addEventListener("DOMContentLoaded", function () {
    initMobileNavigation();
    initResponsiveTableLabels();
    const password = document.getElementById("password");
    const togglePassword = document.getElementById("togglePassword");
    const eyeIcon = document.getElementById("eyeIcon");

    if (password && togglePassword && eyeIcon) {
        togglePassword.addEventListener("click", function () {
            const passwordVisible = password.type === "text";

            password.type = passwordVisible ? "password" : "text";

            if (passwordVisible) {
                eyeIcon.innerHTML = `
                    <path d="M2 12s3.5-7 10-7 10 7 10 7-3.5 7-10 7S2 12 2 12Z"></path>
                    <circle cx="12" cy="12" r="3"></circle>
                `;
                togglePassword.setAttribute("aria-label", "Mostrar senha");
                return;
            }

            eyeIcon.innerHTML = `
                <path d="M3 3l18 18"></path>
                <path d="M10.6 10.6a2 2 0 0 0 2.8 2.8"></path>
                <path d="M9.9 4.2A10.8 10.8 0 0 1 12 4c6.5 0 10 8 10 8a17.2 17.2 0 0 1-3.1 4.2"></path>
                <path d="M6.6 6.6C3.8 8.5 2 12 2 12s3.5 8 10 8a9.8 9.8 0 0 0 4.2-.9"></path>
            `;
            togglePassword.setAttribute("aria-label", "Ocultar senha");
        });
    }

    // Modal de cadastro de produto (Produtos)
    const createModal = document.getElementById("create-product-modal");
    const openCreateModal = document.getElementById("open-create-product");
    const closeCreateModal = document.getElementById("close-create-product");

    if (createModal && openCreateModal && closeCreateModal) {
        const closeModal = function () {
            createModal.hidden = true;
        };

        openCreateModal.addEventListener("click", function () {
            createModal.hidden = false;
            createModal.querySelector("input")?.focus();
        });

        closeCreateModal.addEventListener("click", closeModal);

        createModal.addEventListener("click", function (event) {
            if (event.target === createModal) {
                closeModal();
            }
        });

        document.addEventListener("keydown", function (event) {
            if (event.key === "Escape" && !createModal.hidden) {
                closeModal();
            }
        });
    }

    initAnimations();
});

// ---------------------------------------------------------------------
// Animações globais (entrada suave de conteúdo, linhas de tabela em
// cascata, fechamento animado de modais e destaque de mensagens).
// Ficam num único lugar para funcionar em todas as páginas sem precisar
// mexer em cada template.
// ---------------------------------------------------------------------
function initAnimations() {
    injectAnimationStyles();
    enableFormLoadingState();
    enableDeleteConfirmation();
    enableToasts();
    enableSearchHighlightFromUrl();

    if (window.matchMedia && window.matchMedia("(prefers-reduced-motion: reduce)").matches) {
        return;
    }

    enableClickRipple();
    enablePageTransition();
    animateContentEntrance();
    animateRowsCascade();
    animateCards();
    animateCounters();
    animateModalClosing();
    animateLoginError();
}

function injectAnimationStyles() {
    if (document.getElementById("gs-animations")) return;

    const style = document.createElement("style");
    style.id = "gs-animations";
    style.textContent = `
        .gs-enter {
            opacity: 0;
            transform: translateY(14px);
        }
        .gs-enter.gs-visible {
            opacity: 1;
            transform: none;
            transition: opacity 420ms ease, transform 420ms cubic-bezier(.22,.8,.2,1);
        }
        .gs-row {
            opacity: 0;
            transform: translateY(6px);
            animation: gs-row-enter 380ms ease forwards;
            animation-delay: calc(var(--gs-i, 0) * 45ms);
        }
        @keyframes gs-row-enter {
            to { opacity: 1; transform: none; }
        }
        .gs-banner {
            animation: gs-banner-enter 320ms cubic-bezier(.22,.8,.2,1);
        }
        @keyframes gs-banner-enter {
            from { opacity: 0; transform: translateY(-6px); }
            to { opacity: 1; transform: none; }
        }
        .modal-closing {
            animation: gs-modal-fade-out 160ms ease forwards !important;
        }
        .modal-closing .modal-create-content,
        .modal-closing .modal-edit-content {
            animation: gs-modal-content-fade-out 160ms ease forwards !important;
        }
        @keyframes gs-modal-fade-out {
            to { opacity: 0; }
        }
        @keyframes gs-modal-content-fade-out {
            to { opacity: 0; transform: translateY(10px) scale(.98); }
        }
        button, .edit-product-button, .pagination a, .sidebar-logout {
            transition: transform 120ms ease, background 0.2s, opacity 0.2s;
        }
        button:active, .edit-product-button:active, .pagination a:active {
            transform: scale(0.96);
        }
        .card {
            transition: transform 220ms cubic-bezier(.22,.8,.2,1), box-shadow 220ms ease;
        }
        .card:hover {
            transform: translateY(-5px);
            box-shadow: 0 16px 32px rgba(0, 0, 0, 0.22);
        }
        .gs-ripple {
            position: absolute;
            border-radius: 50%;
            background: rgba(255, 255, 255, 0.45);
            transform: scale(0);
            opacity: 0.55;
            pointer-events: none;
            animation: gs-ripple-anim 550ms ease-out forwards;
        }
        @keyframes gs-ripple-anim {
            to { transform: scale(1.8); opacity: 0; }
        }
        .gs-loading {
            position: relative;
            color: transparent !important;
            pointer-events: none;
        }
        .gs-loading::after {
            content: "";
            position: absolute;
            top: 50%;
            left: 50%;
            width: 16px;
            height: 16px;
            margin: -8px 0 0 -8px;
            border: 2px solid rgba(255, 255, 255, 0.35);
            border-top-color: #fff;
            border-radius: 50%;
            animation: gs-spin 700ms linear infinite;
        }
        @keyframes gs-spin {
            to { transform: rotate(360deg); }
        }
        body.gs-leaving {
            animation: gs-page-leave 180ms ease forwards;
        }
        @keyframes gs-page-leave {
            to { opacity: 0; }
        }
        .gs-shake {
            animation: gs-shake-anim 480ms ease;
        }
        @keyframes gs-shake-anim {
            10%, 90% { transform: translateX(-2px); }
            20%, 80% { transform: translateX(4px); }
            30%, 50%, 70% { transform: translateX(-7px); }
            40%, 60% { transform: translateX(7px); }
        }
        .gs-toasts {
            position: fixed;
            top: 20px;
            right: 20px;
            z-index: 200;
            display: flex;
            flex-direction: column;
            gap: 10px;
            max-width: min(360px, calc(100vw - 40px));
            pointer-events: none;
        }
        .gs-toast {
            display: flex;
            align-items: flex-start;
            gap: 10px;
            padding: 14px 14px 14px 16px;
            border-radius: 10px;
            box-shadow: 0 16px 32px rgba(0, 0, 0, 0.35);
            opacity: 0;
            transform: translateX(24px);
            transition: opacity 220ms ease, transform 220ms cubic-bezier(.22,.8,.2,1);
            pointer-events: auto;
        }
        .gs-toast-visible {
            opacity: 1;
            transform: translateX(0);
        }
        .gs-toast-leaving {
            opacity: 0;
            transform: translateX(24px) scale(0.97);
        }
        .gs-toast-success {
            color: #bbf7d0;
            background: #0f2e1c;
            border: 1px solid #16a34a;
        }
        .gs-toast-error {
            color: #fecaca;
            background: #3b1218;
            border: 1px solid #7f1d1d;
        }
        .gs-toast-icon svg {
            width: 18px;
            height: 18px;
            fill: none;
            stroke: currentColor;
            stroke-width: 2;
            stroke-linecap: round;
            stroke-linejoin: round;
            flex-shrink: 0;
            margin-top: 1px;
        }
        .gs-toast-text {
            flex: 1;
            font-size: 13.5px;
            line-height: 1.4;
        }
        .gs-toast-close {
            background: none;
            border: 0;
            color: inherit;
            opacity: 0.6;
            font-size: 18px;
            line-height: 1;
            cursor: pointer;
            padding: 0;
        }
        .gs-toast-close:hover {
            opacity: 1;
        }
        @media (max-width: 480px) {
            .gs-toasts {
                left: 16px;
                right: 16px;
                top: 12px;
                max-width: none;
            }
        }
        .empty-table {
            display: flex;
            flex-direction: column;
            align-items: center;
            gap: 10px;
            padding: 36px 12px;
            color: #8a93a6;
        }
        .empty-table svg {
            width: 34px;
            height: 34px;
            fill: none;
            stroke: currentColor;
            stroke-width: 1.6;
            stroke-linecap: round;
            stroke-linejoin: round;
            opacity: 0.6;
        }
        .empty-table span {
            font-size: 13.5px;
        }
        .gs-highlight {
            background: rgba(96, 165, 250, 0.35);
            color: #fff;
            border-radius: 3px;
            padding: 0 2px;
        }
        .gs-pulse {
            animation: gs-pulse-anim 380ms cubic-bezier(.22,.8,.2,1);
        }
        @keyframes gs-pulse-anim {
            0% { transform: scale(1); }
            40% { transform: scale(1.18); }
            100% { transform: scale(1); }
        }
        @media (prefers-reduced-motion: reduce) {
            .gs-toast, .gs-pulse, .gs-highlight {
                animation: none !important;
                transition: none !important;
            }
            .gs-toast {
                opacity: 1;
                transform: none;
            }
        }
    `;
    document.head.appendChild(style);
}

// Revela suavemente o bloco principal da página (conteúdo interno,
// tela de acesso negado ou o formulário de login).
function animateContentEntrance() {
    const target = document.querySelector(".content, .access-denied-card, .login-form-side");
    if (!target) return;

    target.classList.add("gs-enter");
    requestAnimationFrame(function () {
        requestAnimationFrame(function () {
            target.classList.add("gs-visible");
        });
    });
}

// Faz as linhas de tabela aparecerem em cascata, uma logo após a outra.
function animateRowsCascade() {
    document.querySelectorAll(".content table tbody tr").forEach(function (row, index) {
        row.style.setProperty("--gs-i", index);
        row.classList.add("gs-row");
    });
}

// Aplica o mesmo efeito de cascata aos cartões do dashboard.
function animateCards() {
    document.querySelectorAll(".cards .card").forEach(function (card, index) {
        card.style.setProperty("--gs-i", index);
        card.classList.add("gs-row");
    });
}

// Anima a saída dos modais (fade + leve deslocamento) antes de
// escondê-los de verdade, sem precisar alterar o código de cada modal.
function animateModalClosing() {
    document.querySelectorAll(".modal-create, .modal-edit, .modal-payment, .modal-receipt").forEach(function (modal) {
        let animating = false;
        let ignoreNextMutation = false;

        const observer = new MutationObserver(function () {
            if (ignoreNextMutation) {
                ignoreNextMutation = false;
                return;
            }
            if (!modal.hidden || animating) return;

            animating = true;
            ignoreNextMutation = true;
            modal.hidden = false;
            modal.classList.add("modal-closing");

            const finish = function (event) {
                if (event.target !== modal) return;
                modal.removeEventListener("animationend", finish);
                modal.classList.remove("modal-closing");
                animating = false;
                ignoreNextMutation = true;
                modal.hidden = true;
            };

            modal.addEventListener("animationend", finish);
        });

        observer.observe(modal, { attributes: true, attributeFilter: ["hidden"] });
    });
}

// Efeito de "onda" (ripple) ao clicar em botões e links de ação —
// delegado no document, então funciona em qualquer botão da aplicação,
// mesmo os criados depois (dentro de modais, por exemplo).
function enableClickRipple() {
    document.addEventListener("click", function (event) {
        const target = event.target.closest(
            "button, .new-product-button, .save-button, .edit-product-button, .sidebar-logout, .pagination a"
        );
        if (!target || target.disabled) return;

        const rect = target.getBoundingClientRect();
        const size = Math.max(rect.width, rect.height);
        const ripple = document.createElement("span");
        ripple.className = "gs-ripple";
        ripple.style.width = ripple.style.height = size + "px";
        ripple.style.left = (event.clientX - rect.left - size / 2) + "px";
        ripple.style.top = (event.clientY - rect.top - size / 2) + "px";

        if (getComputedStyle(target).position === "static") {
            target.style.position = "relative";
        }
        target.style.overflow = "hidden";
        target.appendChild(ripple);
        ripple.addEventListener("animationend", function () { ripple.remove(); });
    });
}

// Conta os números dos cartões do dashboard (faturamento, produtos,
// estoque, vendas...) subindo de 0 até o valor real, em vez de já
// aparecerem prontos.
function animateCounters() {
    document.querySelectorAll(".card strong").forEach(function (el) {
        const originalText = el.textContent.trim();
        const prefix = (originalText.match(/^[^\d]*/) || [""])[0];
        const numberText = originalText.slice(prefix.length).trim();
        const finalValue = parseFloat(numberText.replace(/\./g, "").replace(",", "."));
        if (isNaN(finalValue)) return;

        const decimalPlaces = numberText.includes(",") ? 2 : 0;
        const duration = 900;
        const start = performance.now();

        function format(value) {
            return prefix + value.toLocaleString("pt-BR", {
                minimumFractionDigits: decimalPlaces,
                maximumFractionDigits: decimalPlaces,
            });
        }

        function step(now) {
            const progress = Math.min((now - start) / duration, 1);
            const eased = 1 - Math.pow(1 - progress, 3);
            el.textContent = format(finalValue * eased);
            if (progress < 1) requestAnimationFrame(step);
        }

        el.textContent = format(0);
        requestAnimationFrame(step);
    });
}

// Faz a página desaparecer suavemente antes de navegar para outra tela
// (menu lateral, paginação, sair), em vez de trocar de tela de golpe.
function enablePageTransition() {
    const linkSelectors = ".sidebar nav a, .pagination a, .sidebar-logout";

    document.addEventListener("click", function (event) {
        const link = event.target.closest(linkSelectors);
        if (!link) return;
        if (event.defaultPrevented || event.button !== 0) return;
        if (event.metaKey || event.ctrlKey || event.shiftKey || event.altKey) return;
        if (link.target === "_blank") return;

        event.preventDefault();
        document.body.classList.add("gs-leaving");
        setTimeout(function () {
            window.location.href = link.href;
        }, 160);
    });
}

// Mostra um pequeno spinner no botão de envio enquanto o formulário é
// processado, evitando cliques duplicados em ações mais demoradas.
function enableFormLoadingState() {
    document.addEventListener("submit", function (event) {
        if (event.defaultPrevented) return;

        const form = event.target;
        if (!(form instanceof HTMLFormElement)) return;

        const button = form.querySelector('button[type="submit"], button:not([type])');
        if (button) {
            button.classList.add("gs-loading");
        }
    });
}

// Dá um leve "balanço" no card de login quando a página recarrega com
// uma mensagem de erro de autenticação. A classe não é removida depois:
// tirá-la faria a animação de entrada do card (definida no mesmo seletor)
// disparar de novo, já que o navegador reinicia a propriedade "animation"
// quando ela muda de novo para o valor original.
function animateLoginError() {
    const alert = document.querySelector(".login-alert");
    const card = document.querySelector(".login-card");
    if (!alert || !card) return;

    setTimeout(function () {
        card.classList.add("gs-shake");
    }, 900);
}

// Troca os confirm() nativos do navegador (feios e sem estilo) por um
// modal único, injetado uma vez e reaproveitado em qualquer formulário
// que tenha um botão com "data-confirm" — funciona em qualquer página,
// sem precisar duplicar HTML/JS de modal em cada tela (produtos,
// usuários, etc.).
function enableDeleteConfirmation() {
    if (!document.querySelector("form [data-confirm]")) return;

    let modal = document.getElementById("gs-confirm-modal");
    if (!modal) {
        modal = document.createElement("div");
        modal.id = "gs-confirm-modal";
        modal.className = "modal-create";
        modal.hidden = true;
        modal.innerHTML =
            '<div class="modal-confirm-content" role="alertdialog" aria-modal="true" ' +
            'aria-labelledby="gs-confirm-title" aria-describedby="gs-confirm-text">' +
            '<div class="modal-confirm-icon" aria-hidden="true">' +
            '<svg viewBox="0 0 24 24"><path d="M3 6h18"></path>' +
            '<path d="M8 6V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>' +
            '<path d="M19 6l-1 14a2 2 0 0 1-2 2H8a2 2 0 0 1-2-2L5 6"></path></svg>' +
            "</div>" +
            '<h3 id="gs-confirm-title">Confirmar remoção</h3>' +
            '<p id="gs-confirm-text"></p>' +
            '<div class="modal-confirm-actions">' +
            '<button type="button" class="cancel-delete-button" id="gs-confirm-cancel">Cancelar</button>' +
            '<button type="button" class="confirm-delete-button" id="gs-confirm-ok">Remover</button>' +
            "</div></div>";
        document.body.appendChild(modal);
    }

    const textEl = modal.querySelector("#gs-confirm-text");
    const confirmButton = modal.querySelector("#gs-confirm-ok");
    const cancelButton = modal.querySelector("#gs-confirm-cancel");
    let pendingForm = null;

    const close = function () {
        modal.hidden = true;
        pendingForm = null;
    };

    document.addEventListener("submit", function (event) {
        const form = event.target;
        if (!(form instanceof HTMLFormElement) || form.dataset.confirmed === "true") return;

        const button = form.querySelector("[data-confirm]");
        if (!button) return;

        event.preventDefault();
        pendingForm = form;
        textEl.textContent = button.dataset.confirm;
        modal.hidden = false;
    });

    confirmButton.addEventListener("click", function () {
        if (!pendingForm) return;
        const form = pendingForm;
        const originalButton = form.querySelector("[data-confirm]");
        close();
        form.dataset.confirmed = "true";
        if (originalButton) originalButton.classList.add("gs-loading");
        form.submit();
    });

    cancelButton.addEventListener("click", close);
    modal.addEventListener("click", function (event) {
        if (event.target === modal) close();
    });
    document.addEventListener("keydown", function (event) {
        if (event.key === "Escape" && !modal.hidden) close();
    });
}

// Transforma os banners de sucesso/erro renderizados pelo servidor
// (.form-success / .form-error) em toasts flutuantes que somem
// sozinhos, em vez de ficarem ocupando espaço fixo no topo da página.
function enableToasts() {
    const messages = document.querySelectorAll(".form-success, .form-error");
    if (!messages.length) return;

    let container = document.getElementById("gs-toasts");
    if (!container) {
        container = document.createElement("div");
        container.id = "gs-toasts";
        container.className = "gs-toasts";
        document.body.appendChild(container);
    }

    const successIcon = '<svg viewBox="0 0 24 24"><path d="M20 6 9 17l-5-5"></path></svg>';
    const errorIcon = '<svg viewBox="0 0 24 24"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="8" x2="12" y2="12"></line><line x1="12" y1="16" x2="12.01" y2="16"></line></svg>';

    messages.forEach(function (original) {
        const text = original.textContent.trim();
        if (!text) return;

        const isError = original.classList.contains("form-error");
        original.hidden = true;

        const toast = document.createElement("div");
        toast.className = "gs-toast " + (isError ? "gs-toast-error" : "gs-toast-success");
        toast.setAttribute("role", isError ? "alert" : "status");
        toast.innerHTML =
            '<span class="gs-toast-icon" aria-hidden="true">' + (isError ? errorIcon : successIcon) + "</span>" +
            '<span class="gs-toast-text"></span>' +
            '<button type="button" class="gs-toast-close" aria-label="Fechar">&times;</button>';
        toast.querySelector(".gs-toast-text").textContent = text;

        container.appendChild(toast);
        requestAnimationFrame(function () { toast.classList.add("gs-toast-visible"); });

        let timer;
        const remove = function () {
            clearTimeout(timer);
            toast.classList.remove("gs-toast-visible");
            toast.classList.add("gs-toast-leaving");
            setTimeout(function () { toast.remove(); }, 240);
        };

        toast.querySelector(".gs-toast-close").addEventListener("click", remove);
        toast.addEventListener("mouseenter", function () { clearTimeout(timer); });
        toast.addEventListener("mouseleave", function () { timer = setTimeout(remove, 2500); });
        timer = setTimeout(remove, isError ? 6000 : 4200);
    });
}

// Destaca (em <mark>) o trecho de texto que bateu com a busca — usado
// tanto para buscas feitas pelo servidor (via ?busca= na URL) quanto para
// filtros que rodam só no navegador (como o do PDV). Só mexe em texto
// puro: células com botão/input/link ficam intactas.
function highlightSearch(element, term) {
    if (!element) return;
    clearHighlight(element);

    const cleanTerm = (term || "").trim();
    if (!cleanTerm) return;
    if (element.querySelector("button, input, select, form, a")) return;

    const termLower = cleanTerm.toLowerCase();
    const walker = document.createTreeWalker(element, NodeFilter.SHOW_TEXT);
    const nodes = [];
    let node;
    while ((node = walker.nextNode())) nodes.push(node);

    nodes.forEach(function (textNode) {
        const text = textNode.textContent;
        const index = text.toLowerCase().indexOf(termLower);
        if (index === -1) return;

        const span = document.createElement("span");
        span.className = "gs-highlight-wrap";
        span.appendChild(document.createTextNode(text.slice(0, index)));
        const mark = document.createElement("mark");
        mark.className = "gs-highlight";
        mark.textContent = text.slice(index, index + cleanTerm.length);
        span.appendChild(mark);
        span.appendChild(document.createTextNode(text.slice(index + cleanTerm.length)));
        textNode.replaceWith(span);
    });
}

// Desfaz o destaque aplicado por highlightSearch, devolvendo o texto puro.
function clearHighlight(element) {
    if (!element) return;
    element.querySelectorAll(".gs-highlight-wrap").forEach(function (span) {
        span.replaceWith(document.createTextNode(span.textContent));
    });
    element.normalize();
}

// Ao carregar uma página cujo resultado veio de uma busca feita pelo
// servidor (produtos, usuários...), destaca o termo buscado nas células
// de texto das tabelas.
function enableSearchHighlightFromUrl() {
    const term = new URLSearchParams(window.location.search).get("busca");
    if (!term || !term.trim()) return;

    document.querySelectorAll(".content table tbody td").forEach(function (cell) {
        highlightSearch(cell, term);
    });
}

// Dá um pequeno "pulso" em um elemento — usado no contador do carrinho do
// PDV sempre que a quantidade de itens muda, como feedback de que algo
// foi adicionado/removido.
function pulse(element) {
    if (!element) return;
    element.classList.remove("gs-pulse");
    void element.offsetWidth;
    element.classList.add("gs-pulse");
}

window.gsSearch = { highlight: highlightSearch, clear: clearHighlight };
window.gsPulse = pulse;
