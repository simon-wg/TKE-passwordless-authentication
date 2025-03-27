/**
 * Logs out the user by sending a http request to the server. The server then deletes
 * the session on the server then deletes the cookie in the browser
 * @throws {Error} - Throws an error if the HTTP request fails.
 */
function LogoutButton() {
  const handleClick = async () => {
    try {
      const response = await fetch("/api/logout", {
        method: "POST",
        credentials: "include",
      });

      if (response.ok) {
        window.location.reload();
      } else {
        const errorData = await response.json();
        console.error(
          "Logout failed:",
          errorData.message || response.statusText
        );
      }
    } catch (error) {
      console.error("Unable to logout user:", error);
    }
  };

  return (
    <button
      className={"logout-button"}
      onClick={handleClick}
      style={{ width: "fit-content" }}
    >
      Logout
    </button>
  );
}

export default LogoutButton;
