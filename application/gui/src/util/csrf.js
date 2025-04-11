/* Function that fetches and saves CSRF Tokens from server
 *
 * @returns a CSRF token
 *
 */

let csrfTokenCache = null;

export const fetchCsrfToken = async () => {
  if (csrfTokenCache) return csrfTokenCache;

  const res = await fetch("/api/csrf-token", {
    credentials: "include",
  });

  if (!res.ok) throw new Error("Failed to fetch CSRF token");

  const token = res.headers.get("X-CSRF-Token");
  if (!token) throw new Error("CSRF token not found in response headers");

  csrfTokenCache = token;
  console.log("Fetched csrf token:", { csrfTokenCache });
  return csrfTokenCache;
};

export const clearCsrfToken = () => {
  csrfTokenCache = null;
};
