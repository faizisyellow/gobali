function decodeJWTClaims(token) {
  if (!token || typeof token !== "string") {
    throw new Error("Token must be a valid string");
  }

  const parts = token.split(".");

  if (parts.length !== 3) {
    throw new Error("Invalid JWT format");
  }

  try {
    let payload = parts[1];

    payload = payload.replace(/-/g, "+").replace(/_/g, "/");

    while (payload.length % 4) {
      payload += "=";
    }

    const decoded = JSON.parse(atob(payload));

    return decoded;
  } catch (error) {
    throw new Error("Failed to decode JWT token: " + error.message);
  }
}
export { decodeJWTClaims };
