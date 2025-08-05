export function showLoader(loaderId = "loader") {
  const loader = document.getElementById(loaderId);
  if (loader) loader.classList.remove("hidden");
}

export function hideLoader(loaderId = "loader") {
  const loader = document.getElementById(loaderId);
  if (loader) loader.classList.add("hidden");
}AC