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
