import { useState, useEffect } from "react";

function useVerifySignature(publicKeyPem, message, signatureBase64) {
  const [isValid, setIsValid] = useState(null);

  useEffect(() => {
    async function verify() {
      try {
        if (!publicKeyPem || !message || !signatureBase64) {
          setIsValid(false);
          return;
        }

        // Import the public key
        const publicKey = await importPublicKey(publicKeyPem);

        // Hash the original message
        const hashedMessage = await hashMessage(message);

        // Decrypt the signature (recover original hash)
        const decryptedHash = await decryptSignature(
          publicKey,
          signatureBase64
        );

        // Compare the hashes
        setIsValid(hashedMessage === decryptedHash);
      } catch (error) {
        console.error("Error verifying signature:", error);
        setIsValid(false);
      }
    }

    verify();
  }, [publicKeyPem, message, signatureBase64]);

  return isValid;
}

// 🔹 Convert PEM Public Key to CryptoKey
async function importPublicKey() {
  const publicKeyPem = import.meta.env.REACT_APP_PUBLIC_KEY; // Fetch from .env

  if (!publicKeyPem) {
    throw new Error("Public key is missing from environment variables.");
  }

  return await crypto.subtle.importKey(
    "spki",
    binaryKey.buffer,
    { name: "RSASSA-PKCS1-v1_5", hash: "SHA-256" },
    true,
    ["verify"]
  );
}

// 🔹 Hash Message with SHA-256
async function hashMessage(message) {
  const encoder = new TextEncoder();
  const messageBuffer = encoder.encode(message);
  const hashBuffer = await crypto.subtle.digest("SHA-256", messageBuffer);
  return bufferToHex(hashBuffer);
}

// 🔹 Decrypt the Signature to Recover the Hash
async function decryptSignature(publicKey, signatureBase64) {
  const signature = Uint8Array.from(atob(signatureBase64), (c) =>
    c.charCodeAt(0)
  );

  try {
    // Verify signature instead of decrypting
    const isValid = await crypto.subtle.verify(
      { name: "RSASSA-PKCS1-v1_5", hash: "SHA-256" },
      publicKey,
      signature,
      new Uint8Array(signature.length).buffer
    );

    return isValid ? bufferToHex(signature) : null;
  } catch (error) {
    console.error("Signature verification failed:", error);
    return null;
  }
}

// 🔹 Convert ArrayBuffer to Hex
function bufferToHex(buffer) {
  return [...new Uint8Array(buffer)]
    .map((b) => b.toString(16).padStart(2, "0"))
    .join("");
}

export default useVerifySignature;
