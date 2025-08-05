document.addEventListener("DOMContentLoaded", () => {
    const form = document.getElementById("node-form");
  
    form?.addEventListener("submit", async (e) => {
      e.preventDefault();
  
      const projectName = document.getElementById("node-id").value;
      const cloudProvider = document.querySelector('input[name="cloud"]:checked')?.value;
  
      if (!cloudProvider || !projectName) {
        alert("Please fill out all fields.");
        return;
      }
  
      const payload = {
        user_id: "demo_user_123", // Replace this with Clerk user info if available
        project_name: projectName,
        cloud_provider: cloudProvider,
      };
  
      try {
        const res = await fetch("/api/scenario", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(payload),
        });
  
        const data = await res.json();
        if (res.ok) {
          alert("Scenario saved!");
          form.reset();
        } else {
          alert("Error: " + data.error);
        }
      } catch (err) {
        console.error(err);
        alert("Failed to save scenario.");
      }
    });
  });
  