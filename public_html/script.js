let htmlElement;

window.addEventListener("DOMContentLoaded", function () {
  htmlElement = document.querySelector("html");

  if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
    htmlElement.setAttribute("data-bs-theme", "dark");
  } else {
    htmlElement.setAttribute("data-bs-theme", "light");
  }

});

window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', event => {
  const colorScheme = event.matches ? "dark" : "light";
  if (htmlElement) {
    htmlElement.setAttribute("data-bs-theme", colorScheme);
  }
});

function changeTheme() {
  const currentTheme = htmlElement.getAttribute("data-bs-theme");
  if (currentTheme === "dark") {
    document.getElementById("mode-icon").className = "bi bi-moon-stars-fill";
    if (htmlElement) {
      htmlElement.setAttribute("data-bs-theme", "light");
    }
  } else {
    document.getElementById("mode-icon").className = "bi bi-sun-fill";
    if (htmlElement) {
      htmlElement.setAttribute("data-bs-theme", "dark");
    }
  }
}
