document.addEventListener("DOMContentLoaded", () => {
    const form = document.getElementById("signupForm");
    if (form) {
        form.addEventListener("submit", async (event) => {
            event.preventDefault();

            const formData = new FormData(form);
            const payload = {
                name: formData.get("name"),
                email: formData.get("email"),
                phone: formData.get("phone"),
                password: formData.get("password"),
            };
            
            console.log("Form data:", payload);  // Log form data to check if password is being sent

            try {
                const response = await fetch("/api/v1/user/signup", {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify(payload),
                });

                if (response.ok) {
                    const result = await response.json();
                    alert("Signup successful: " + result.message);
                    window.location.href = "../login"
                } else {
                    alert("Signup failed: " + response.statusText);
                }
            } catch (error) {
                console.error("Error:", error);
                alert("An error occurred during signup.");
            }
        });
    }
});
