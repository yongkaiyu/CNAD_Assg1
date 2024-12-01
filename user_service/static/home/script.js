document.addEventListener("DOMContentLoaded", () => {
    const userId = localStorage.getItem("userId");
    const userName = localStorage.getItem("userName");
    const userEmail = localStorage.getItem("userEmail");

    if (userId && userName && userEmail) {
        const userInfoDiv = document.getElementById("userInfo");
        userInfoDiv.innerHTML = `
            <p>User ID: ${userId}</p>
            <p>Name: ${userName}</p>
            <p>Email: ${userEmail}</p>
        `;
    } else {
        alert("No user data found. Please log in.");
        window.location.href = "../../static/login"; // Redirect to login page if no data
        return
    }

    // Button redirects
    const settingsButton = document.getElementById("settingsButton");
    const historyButton = document.getElementById("historyButton");

    if (settingsButton) {
        settingsButton.addEventListener("click", () => {
            window.location.href = "../../static/settings";
        });
    }

    if (historyButton) {
        historyButton.addEventListener("click", () => {
            window.location.href = "../../static/history";
        });
    }
});