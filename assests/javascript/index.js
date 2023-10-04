const sideMenu = document.querySelector("aside");
const menuBtn = document.querySelector("#menu-btn");
const closeBtn = document.querySelector("#close-btn");
const themeToggler = document.querySelector(".theme-toggler");


// show sidebar
menuBtn.addEventListener("click", () => {
  sideMenu.style.display = "block";
});

closeBtn.addEventListener("click", () => {
  sideMenu.style.display = "none";
});

// change theme
themeToggler.addEventListener("click", () => {
  document.body.classList.toggle("dark-theme-variables");

  themeToggler.querySelector('span:first-child').classList.toggle('active');
  themeToggler.querySelector('span:last-child').classList.toggle('active');
});

//For posting with date in index.html
var dateInput = document.getElementById("selected-date");

dateInput.addEventListener("change", function() {
  document.getElementById("search-by-date").submit();
});