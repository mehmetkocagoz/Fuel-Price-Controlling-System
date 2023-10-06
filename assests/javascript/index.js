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

// Create a new Date object representing today's date
var today = new Date();

// Format the date as "YYYY-MM-DD" for the input field
var formattedDate = today.toISOString().split("T")[0];

// Set the date input field's value to today's date
dateInput.value = formattedDate;

dateInput.addEventListener("change", function() {
  document.getElementById("search-by-date").submit();
});