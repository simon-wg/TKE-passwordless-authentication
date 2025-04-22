/* A function that makes API calls with a csrf token
 * Optional params that allow you to specify request type (POST,DELETE, etc)
 * will default to GET if nothing provided.
 * Always includes credentials
 * In practice works exactly as fetch() only with csrf tokens and cookies.
 */
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
