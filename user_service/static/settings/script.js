document.addEventListener("DOMContentLoaded", function() {
    // Get the user ID from localStorage (assuming it's stored when the user logs in)
    const userId = localStorage.getItem('userId');
    const userName = localStorage.getItem("userName");
    const userEmail = localStorage.getItem("userEmail");

    if (!userId || !userName || !userEmail) {
        alert("No user data found. Please log in.");
        window.location.href = "../home"; // Redirect to login page if no data
    }

    // Trim and sanitize userId so it wont be 1:1 but 1 instead
    /* if (userId) {
      userId = userId.trim().replace(/[^0-9]/g, ''); // Keep only numbers
    } */

    const updateProfileForm = document.getElementById("updateProfileForm");
    const membershipTierElement = document.getElementById("membershipTier");
    
    // Fetch Membership Status
    const fetchMembershipStatus = async () => {
      const userId = localStorage.getItem("userId");
      if (!userId) {
          membershipTierElement.textContent = "No user logged in.";
          return;
      }

      try {
          const response = await fetch(`/api/v1/user/settings?user_id=${userId}`, {
              method: "GET",
          });
          
          if (response.ok) {
              let data = await response.text(); // Get the raw text response

              // Clean up any non-JSON characters if present (like parentheses)
              data = data.replace(/[^\{]*\{/, '{').replace(/\}[^}]*$/, '}');
              // Parse the cleaned-up JSON
              const jsonData = JSON.parse(data);
              membershipTierElement.textContent = `Membership Tier: ${jsonData.membership_tier || "Unknown"}`;
          } else {
              membershipTierElement.textContent = "Failed to load membership status.";
              console.error("Failed to fetch membership status:", response.statusText);
          }
      } catch (error) {
          console.error("Error fetching membership status:", error);
          membershipTierElement.textContent = "Error loading membership status.";
      }
  };

  // Fetch membership status if the element exists
  if (membershipTierElement) {
      fetchMembershipStatus();
  }

  // Handle profile updates
  if (updateProfileForm) {
      updateProfileForm.addEventListener("submit", async (event) => {
          event.preventDefault();

          const formData = new FormData(updateProfileForm);
          const payload = {
              name: formData.get("name"),
              email: formData.get("email"),
              phone: formData.get("phone"),
              password: formData.get("password"),
          };

          const userId = localStorage.getItem("userId");
          if (!userId) {
              alert("No user is logged in.");
              return;
          }

          try {
              const response = await fetch(`/api/v1/user/settings?user_id=${userId}`, {
                  method: "PUT",
                  headers: {
                      "Content-Type": "application/json",
                  },
                  body: JSON.stringify(payload),
              });

              if (response.ok) {
                  const data = await response.json();
                  alert(data.message || "Profile updated successfully!");
                  // Optionally update local storage with new details
                  localStorage.setItem("userName", payload.name);
                  localStorage.setItem("userEmail", payload.email);
                  localStorage.setItem("userPhone", payload.phone);
              } else {
                  alert("Failed to update profile: " + response.statusText);
              }
          } catch (error) {
              console.error("Error updating profile:", error);
              alert("An error occurred while updating the profile.");
          }
      });
  }

});

// Fetch and display the membership status
/* function fetchMembershipStatus() {
  console.log("Fetching membership status...");
  fetch(`/api/v1/user/settings?user_id=${userId}`, {
    method: 'GET',
  })
    .then((response) => {
        console.log("Response received:", response);
        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }
        return response.json();
    })
    .then((data) => {
        console.log("Response data:", data);
        if (data.membership_tier) {
            document.getElementById('membershipTier').textContent = `Membership Tier: ${data.membership_tier}`;
        } else {
            document.getElementById('membershipTier').textContent = 'Membership tier not found.';
        }
    })
    .catch((error) => {
        console.error("Error fetching membership status:", error);
        document.getElementById('membershipTier').textContent = 'Error fetching membership status.';
    });
}

document.getElementById("updateProfileForm").addEventListener("submit", async function(event) {
    event.preventDefault();

    const name = document.getElementById("name").value.trim();
    const email = document.getElementById("email").value.trim();
    const phone = document.getElementById("phone").value.trim();
    const password = document.getElementById("password").value.trim();

    if (!name || !email) {
        alert("Name and email are required.");
        return;
    }

    const formData = new FormData();
    formData.append("name", name);
    formData.append("email", email);
    formData.append("phone", phone);
    formData.append("password", password);

    fetch(`/api/v1/user/settings?user_id=${userId}`, {
        method: 'PUT',
        body: formData,
    })
    .then((response) => response.json())
    .then((data) => {
        console.log("Response data:", data);
        if (data.message) {
            document.getElementById('updateMessage').textContent = data.message;
            fetchMembershipStatus();
            localStorage.setItem("userName", userData.name);
            localStorage.setItem("userEmail", userData.email);
            localStorage.setItem("userPhone", userData.phone);
        } else {
            document.getElementById('updateMessage').textContent = 'Error updating profile.';
        }
    })
    .catch((error) => {
        console.error("Error updating profile:", error);
        document.getElementById('updateMessage').textContent = 'Error updating profile.';
    });
}); */