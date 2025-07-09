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
        const res = await fetch("/api/shorten", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
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
            <a href="${data.short_url}" target="_blank" style="font-size: 1.1rem;">${data.short_url}</a>
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
      const res = await fetch("/api/links");
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
        `;
        linksBody.appendChild(tr);
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
