document.addEventListener('DOMContentLoaded', () => {
  console.log("Certifications JS loaded");
  fetch('/api/certifications')
    .then(res => res.json())
    .then(data => {
      console.log("Certifications data:", data);
      // If data is an object with a 'certifications' key, use that
      if (!Array.isArray(data)) {
        if (Array.isArray(data.certifications)) {
          data = data.certifications;
        } else {
          console.error("Certifications data is not an array:", data);
          return;
        }
      }
      const tbody = document.getElementById('certTableBody');
      if (!tbody) {
        console.error("Table body not found");
        return;
      }
      tbody.innerHTML = '';
      data.forEach(cert => {
        const tr = document.createElement('tr');
        tr.innerHTML = `
          <td>${cert.Name || ""}</td>
          <td>${cert.Issuer || ""}</td>
          <td>${cert.Platform || ""}</td>
          <td>${cert.Filename || ""}</td>
        `;
        tbody.appendChild(tr);
      });
    })
    .catch(err => {
      console.error("Error loading certifications:", err);
    });
});