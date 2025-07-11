// Script pour la page principale (index.html)
document.addEventListener("DOMContentLoaded", function () {
  const shortenForm = document.getElementById("shorten-form");

  if (shortenForm) {
    // Animation d'entrée au chargement de la page
    const container = document.querySelector(".container");
    if (container) {
      container.style.opacity = "0";
      container.style.transform = "translateY(30px)";

      setTimeout(() => {
        container.style.transition = "all 0.5s ease-out";
        container.style.opacity = "1";
        container.style.transform = "translateY(0)";
      }, 100);
    }

    // Gestion du formulaire de raccourcissement
    shortenForm.addEventListener("submit", async (e) => {
      e.preventDefault();

      const url = document.getElementById("url").value.trim();
      const alias = document.getElementById("alias").value.trim();
      const expires = document.getElementById("expires").value.trim();
      const submitButton = document.querySelector("button[type='submit']");
      const messageDiv = document.getElementById("message");

      // Préparation de la requête
      const payload = { url };
      if (alias) payload.alias = alias;
      if (expires) payload.expiration_minutes = parseInt(expires, 10);

      // Affichage du loading
      submitButton.disabled = true;
      submitButton.innerHTML =
        '<span class="loading"></span>Raccourcissement...';

      messageDiv.style.display = "block";
      messageDiv.className = "result";
      messageDiv.innerHTML =
        '<span class="loading"></span>Traitement en cours...';

      try {
        const headers = { "Content-Type": "application/json" };
        if (localStorage.token)
          headers.Authorization = "Bearer " + localStorage.token;

        const res = await fetch("/api/private/shorten", {
          method: "POST",
          headers,
          body: JSON.stringify(payload),
        });

        const data = await res.json();

        if (!res.ok) {
          throw new Error(data.message || "Erreur inconnue");
        }

        // Affichage du succès avec animation
        let html = `
          <div style="text-align: center; margin-bottom: 1rem;">
            <h3 style="margin: 0.5rem 0; color: #2d3748;">URL raccourcie avec succès !</h3>
          </div>
          <div style="background: #fff; padding: 1rem; border-radius: 8px; margin: 1rem 0;">
            <strong>URL courte :</strong><br>
            <a href="${
              data.short_url
            }" target="_blank" style="font-size: 1.1rem;">${data.short_url}</a>
          </div>
        <div style="text-align:center;">
          <img src="/qr/${data.short_url
            .split("/")
            .pop()}" alt="QR Code" style="max-width:150px; margin-top:1rem;" />
          <div style="font-size:0.8rem;color:#718096;margin-top:0.5rem;">Scannez pour ouvrir</div>
        </div>
        `;

        if (data.expires_at) {
          const d = new Date(data.expires_at);
          html += `
            <div style="text-align: center; color: #718096; font-size: 0.9rem;">
              Expire le : ${d.toLocaleString()}
            </div>
          `;
        }

        messageDiv.innerHTML = html;

        // Réinitialisation du formulaire après succès
        setTimeout(() => {
          shortenForm.reset();
        }, 1000);
      } catch (err) {
        // Affichage de l'erreur
        messageDiv.className = "result error";
        messageDiv.innerHTML = `
          <div style="text-align: center;">
            <h3 style="margin: 0.5rem 0;">Erreur</h3>
            <p>${err.message}</p>
          </div>
        `;
      } finally {
        // Restauration du bouton
        submitButton.disabled = false;
        submitButton.innerHTML = "Raccourcir l'URL";
      }
    });

    // Amélioration de l'expérience utilisateur avec les inputs
    document.querySelectorAll("input").forEach((input) => {
      input.addEventListener("focus", function () {
        this.parentElement.style.transform = "scale(1.02)";
      });

      input.addEventListener("blur", function () {
        this.parentElement.style.transform = "scale(1)";
      });
    });
  }
});

// Script pour la page des liens (links.html)
document.addEventListener("DOMContentLoaded", async function () {
  const linksBody = document.getElementById("links-body");

  if (linksBody) {
    // Animation d'entrée au chargement de la page
    const container = document.querySelector(".links-container");
    if (container) {
      container.style.opacity = "0";
      container.style.transform = "translateY(30px)";

      setTimeout(() => {
        container.style.transition = "all 0.5s ease-out";
        container.style.opacity = "1";
        container.style.transform = "translateY(0)";
      }, 100);
    }

    // Chargement des liens
    try {
      const headers = { "Content-Type": "application/json" };
      if (localStorage.token) {
        headers.Authorization = "Bearer " + localStorage.token;
      } else {
        window.location.href = "/login.html";
      }
      const res = await fetch("/api/private/links", { headers });
      const data = await res.json();

      if (!res.ok) {
        throw new Error(data.message || "Erreur inconnue");
      }

      if (data.length === 0) {
        linksBody.innerHTML = `
          <tr>
            <td colspan="3" class="empty-state">
              Aucun lien raccourci pour le moment
            </td>
          </tr>
        `;
        return;
      }

      data.forEach((link) => {
        const tr = document.createElement("tr");
        tr.innerHTML = `
          <td class="url-cell">
            <a href="${link.target_url}" target="_blank" title="${link.target_url}">
              ${link.target_url}
            </a>
          </td>
          <td class="short-url-cell">
            <a href="/${link.alias}" target="_blank">
              ${location.origin}/${link.alias}
            </a>
          </td>
          <td class="center">
            ${link.click_count}
          </td>
          <td class="center">
            <button class="qr-btn auth-btn" data-alias="${link.alias}">QR</button>
          </td>
        `;
        linksBody.appendChild(tr);
      });

      // Délégation d'événement pour boutons QR
      linksBody.addEventListener("click", function (e) {
        const btn = e.target.closest(".qr-btn");
        if (!btn) return;
        const alias = btn.getAttribute("data-alias");
        window.open(`/qr/${alias}`, "_blank");
      });
    } catch (err) {
      linksBody.innerHTML = `
        <tr>
          <td colspan="3" class="error-state">
            ${err.message}
          </td>
        </tr>
      `;
    }
  }
});

// Script pour la page de connexion (login.html)
document.addEventListener("DOMContentLoaded", function () {
  const loginForm = document.getElementById("login-form");
  if (loginForm) {
    loginForm.addEventListener("submit", async (e) => {
      e.preventDefault();
      const username = document.getElementById("login-username").value.trim();
      const password = document.getElementById("login-password").value;
      const messageDiv = document.getElementById("login-message");

      messageDiv.style.display = "block";
      messageDiv.className = "result";
      messageDiv.textContent = "Connexion en cours...";

      try {
        const res = await fetch("/api/login", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ username, password }),
        });
        const data = await res.json();
        if (!res.ok) throw new Error(data.message || "Erreur inconnue");

        // Stocke le token dans le stockage local
        localStorage.setItem("token", data.token);
        messageDiv.textContent = "Connexion réussie ! Redirection...";
        setTimeout(() => {
          window.location.href = "/"; // redirige vers la page principale
        }, 800);
      } catch (err) {
        messageDiv.className = "result error";
        messageDiv.textContent = err.message;
      }
    });
  }
});

// Script pour la page d'inscription (register.html)
document.addEventListener("DOMContentLoaded", function () {
  const regForm = document.getElementById("register-form");
  if (!regForm) return; // on n'est pas sur register.html

  regForm.addEventListener("submit", async (e) => {
    e.preventDefault();
    const username = document.getElementById("reg-username").value.trim();
    const password = document.getElementById("reg-password").value;
    const msg = document.getElementById("register-message");

    msg.style.display = "block";
    msg.className = "result";
    msg.textContent = "Création du compte…";

    try {
      const res = await fetch("/api/register", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password }),
      });

      if (!res.ok) {
        const data = await res.json();
        throw new Error(data.message || "Erreur inconnue");
      }

      msg.textContent = "Compte créé ! Redirection vers la connexion…";
      setTimeout(() => (window.location.href = "/login.html"), 800);
    } catch (err) {
      msg.className = "result error";
      msg.textContent = err.message;
    }
  });
});

// Gestion du widget d'authentification
document.addEventListener("DOMContentLoaded", function () {
  const authWidget = document.getElementById("auth-widget");
  if (!authWidget) return; // pas sur la page d'index

  const authLoggedOut = document.getElementById("auth-logged-out");
  const authLoggedIn = document.getElementById("auth-logged-in");
  const usernameDisplay = document.getElementById("username-display");
  const loginBtn = document.getElementById("login-btn");
  const registerBtn = document.getElementById("register-btn");
  const logoutBtn = document.getElementById("logout-btn");

  // Fonction pour décoder le token JWT (partie payload)
  function decodeJWT(token) {
    try {
      const base64Url = token.split(".")[1];
      const base64 = base64Url.replace(/-/g, "+").replace(/_/g, "/");
      const jsonPayload = decodeURIComponent(
        atob(base64)
          .split("")
          .map(function (c) {
            return "%" + ("00" + c.charCodeAt(0).toString(16)).slice(-2);
          })
          .join("")
      );
      return JSON.parse(jsonPayload);
    } catch (error) {
      console.error("Erreur lors du décodage du token:", error);
      return null;
    }
  }

  // Fonction pour vérifier si le token est expiré
  function isTokenExpired(token) {
    const payload = decodeJWT(token);
    if (!payload || !payload.exp) return true;
    return Date.now() >= payload.exp * 1000;
  }

  // Fonction pour mettre à jour l'affichage du widget
  function updateAuthWidget() {
    const token = localStorage.getItem("token");

    if (!token || isTokenExpired(token)) {
      // Utilisateur non connecté ou token expiré
      if (token && isTokenExpired(token)) {
        localStorage.removeItem("token");
      }
      authLoggedOut.style.display = "flex";
      authLoggedIn.style.display = "none";
      const adminLink = document.getElementById("admin-link");
      if (adminLink) adminLink.style.display = "none";
    } else {
      // Utilisateur connecté
      const payload = decodeJWT(token);
      if (payload && payload.sub) {
        usernameDisplay.textContent = `👤 ${payload.sub}`;
        if (payload.adm) {
          usernameDisplay.textContent += " (Admin)";
          const adminLink = document.getElementById("admin-link");
          if (adminLink) adminLink.style.display = "inline-block";
        }
      }
      authLoggedOut.style.display = "none";
      authLoggedIn.style.display = "flex";
    }
  }

  // Gestionnaires d'événements pour les boutons
  loginBtn.addEventListener("click", function () {
    window.location.href = "/login.html";
  });

  registerBtn.addEventListener("click", function () {
    window.location.href = "/register.html";
  });

  logoutBtn.addEventListener("click", function () {
    localStorage.removeItem("token");
    updateAuthWidget();

    // Optionnel : rafraîchir la page ou afficher un message
    const messageDiv = document.getElementById("message");
    if (messageDiv) {
      messageDiv.style.display = "block";
      messageDiv.className = "result";
      messageDiv.innerHTML =
        '<div style="text-align: center;">Déconnexion réussie !</div>';
      setTimeout(() => {
        messageDiv.style.display = "none";
      }, 3000);
    }
  });

  // Mettre à jour l'affichage au chargement
  updateAuthWidget();

  // Vérifier périodiquement si le token a expiré
  setInterval(updateAuthWidget, 60000); // Vérifier chaque minute
});

// Script pour la page admin (admin.html)
document.addEventListener("DOMContentLoaded", async function () {
  const usersBody = document.getElementById("users-body");
  if (!usersBody) return; // pas sur admin.html

  const msgDiv = document.getElementById("admin-message");

  function showMessage(text, isError = false) {
    msgDiv.style.display = "block";
    msgDiv.className = isError ? "result error" : "result";
    msgDiv.textContent = text;
    setTimeout(() => (msgDiv.style.display = "none"), 3000);
  }

  async function loadUsers() {
    usersBody.innerHTML = "<tr><td colspan='3'>Chargement...</td></tr>";
    try {
      const headers = { "Content-Type": "application/json" };
      if (!localStorage.token) {
        window.location.href = "/login.html";
        return;
      }
      headers.Authorization = "Bearer " + localStorage.token;
      const res = await fetch("/api/private/admin/users", { headers });
      const data = await res.json();
      if (!res.ok) throw new Error(data.message || "Erreur inconnue");

      if (data.length === 0) {
        usersBody.innerHTML = `<tr><td colspan="3" class="empty-state">Aucun utilisateur</td></tr>`;
        return;
      }

      usersBody.innerHTML = "";
      data.forEach((user) => {
        const tr = document.createElement("tr");
        tr.innerHTML = `
          <td>${user.username}</td>
          <td class="center">${user.is_admin ? "✔️" : ""}</td>
          <td class="center">
            <button class="auth-btn logout delete-user" data-username="${
              user.username
            }">Supprimer</button>
          </td>
        `;
        usersBody.appendChild(tr);
      });
    } catch (err) {
      usersBody.innerHTML = `<tr><td colspan="3" class="error-state">${err.message}</td></tr>`;
    }
  }

  // Chargement initial
  await loadUsers();

  // Délégation d'événement pour les boutons supprimer
  usersBody.addEventListener("click", async function (e) {
    const btn = e.target.closest(".delete-user");
    if (!btn) return;
    const username = btn.getAttribute("data-username");
    if (!confirm(`Supprimer l'utilisateur ${username} ?`)) return;

    try {
      const headers = { "Content-Type": "application/json" };
      headers.Authorization = "Bearer " + localStorage.token;
      const res = await fetch(`/api/private/admin/users/${username}`, {
        method: "DELETE",
        headers,
      });
      if (!res.ok) {
        const data = await res.json();
        throw new Error(data.message || "Erreur inconnue");
      }
      showMessage("Utilisateur supprimé");
      await loadUsers();
    } catch (err) {
      showMessage(err.message, true);
    }
  });
});
