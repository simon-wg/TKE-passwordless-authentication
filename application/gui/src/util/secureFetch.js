import { fetchCsrfToken } from "./csrf";

export const secureFetch = async (url, options = {}) => {
  const csrfToken = await fetchCsrfToken();
  console.log("Token is: ", { csrfToken });
  const headers = {
    "Content-Type": "application/json",
    "X-CSRF-Token": csrfToken,
    ...(options.headers || {}),
  };

  const finalOptions = {
    ...options,
    method: options.method || "GET",
    headers,
    credentials: "include",
  };

  return fetch(url, finalOptions);
};
