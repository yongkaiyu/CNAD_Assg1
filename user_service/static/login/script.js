document.addEventListener("DOMContentLoaded", () => {
    // Login form handler
    const loginForm = document.getElementById("loginForm");
    if (loginForm) {
        loginForm.addEventListener("submit", async (event) => {
            event.preventDefault();

            const formData = new FormData(loginForm);
            const payload = {
                email: formData.get("email"),
                password: formData.get("password"),
            };

            try {
                const response = await fetch("/api/v1/user/login", {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify(payload),
                });

                if (response.ok) {
                    const userData = await response.json(); // Assuming login response is JSON
                    localStorage.setItem("userId", userData.user_id);
                    localStorage.setItem("userName", userData.name);
                    localStorage.setItem("userEmail", userData.email);
                    localStorage.setItem("userPhone", userData.phone);
                    alert("Login successful!");
                    window.location.href = "../../static/home"
                } else {
                    alert("Login failed: " + response.statusText);
                }
            } catch (error) {
                console.error("Error:", error);
                alert("An error occurred during login.");
            }
        });
    }
});
