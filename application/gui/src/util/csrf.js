let csrfTokenCache = null;

export const fetchCsrfToken = async () => {
  if (csrfTokenCache) return csrfTokenCache;

  const res = await fetch("/api/csrf-token", {
    credentials: "include",
  });

  if (!res.ok) throw new Error("Failed to fetch CSRF token");

  const data = await res.json();
  csrfTokenCache = data.csrfToken;
  console.log("Fetched csrf token:", { csrfTokenCache });
  return csrfTokenCache;
};

export const clearCsrfToken = () => {
  csrfTokenCache = null;
};
