$(document).ready(function () {
    $('#onlinebotable').DataTable({
        responsive: true,
        "iDisplayLength": 5,
        "aLengthMenu": [[5, 10, 25], ["5 Per Page", "10 Per Page", "25 Per Page"]]
    });
    $('#offlinebotable').DataTable({
        responsive: true,
        "iDisplayLength": 5,
        "aLengthMenu": [[5, 10, 25], ["5 Per Page", "10 Per Page", "25 Per Page"]]
    });

    const imageContainer = document.querySelector(".image-container");
    const closeButton = document.getElementById("close-button");

    imageContainer.addEventListener("click", () => {
        if (!imageContainer.classList.contains("expanded")) {
            imageContainer.classList.add("expanded");
            closeButton.classList.remove("hidden");
        }
    });

    closeButton.addEventListener("click", (event) => {
        event.stopPropagation(); // Prevent triggering the image click event
        imageContainer.classList.remove("expanded");
        closeButton.classList.add("hidden");
    });
});