// Loads and initializes the WebAssembly module for EUI-64 calculations.
const go = new Go();
WebAssembly.instantiateStreaming(fetch("./main.wasm"), go.importObject)
  .then((result) => {
    go.run(result.instance);
    console.log("WebAssembly module initialized");

    // Ensure required WebAssembly functions are available.
    if (
      typeof window.validateMAC !== "function" ||
      typeof window.validateIPv6Prefix !== "function" ||
      typeof window.calculateEUI64 !== "function"
    ) {
      console.error("Required WebAssembly functions missing");
    }
  })
  .catch((err) => console.error("Failed to load WebAssembly module:", err));

// Copies the value of an input element to the clipboard and updates the button's tooltip to "Copied!" for 2 seconds.
function copyToClipboard(elementId, buttonId) {
  const input = document.getElementById(elementId);
  const button = document.getElementById(buttonId);

  // Validate input and button elements exist.
  if (!input || !button) {
    console.error(`Element not found: ${elementId} or ${buttonId}`);
    return;
  }

  // Copy input value to clipboard and update tooltip.
  navigator.clipboard
    .writeText(input.value)
    .then(() => {
      button.classList.add("copied");
      const tooltip = button.querySelector(".copy-tooltip");
      tooltip.textContent = "Copied!";
      setTimeout(() => {
        button.classList.remove("copied");
        tooltip.textContent = "Copy";
      }, 2000);
    })
    .catch((err) => console.error("Clipboard copy failed:", err));
}

// Sets up form event listeners for submission and clearing, handling input validation and EUI-64 calculation via WebAssembly.
document.addEventListener("DOMContentLoaded", () => {
  // Retrieve DOM elements for form interaction.
  const form = document.querySelector("form");
  const resultContainer = document.querySelector(".result-container");
  const formResults = document.querySelector(".form-results");
  const macInput = document.getElementById("mac");
  const prefixInput = document.getElementById("ip-start");
  const copyMac = document.getElementById("copy-mac");
  const copyPrefix = document.getElementById("copy-ip-start");

  // Validate all required DOM elements are present.
  if (
    !form ||
    !resultContainer ||
    !formResults ||
    !macInput ||
    !prefixInput ||
    !copyMac ||
    !copyPrefix
  ) {
    console.error("Required DOM elements missing");
    return;
  }

  console.log("Form event listeners set up");

  // Initialize tooltips for input copy buttons.
  [copyMac, copyPrefix].forEach((button) => {
    const tooltip = button.querySelector(".copy-tooltip");
    if (tooltip && !tooltip.textContent) {
      tooltip.textContent = "Copy";
    }
  });

  // Attach copy event listeners for MAC and IPv6 prefix inputs.
  copyMac.addEventListener("click", () => copyToClipboard("mac", "copy-mac"));
  copyPrefix.addEventListener("click", () =>
    copyToClipboard("ip-start", "copy-ip-start")
  );

  // Handle form submission for EUI-64 calculation.
  form.addEventListener("submit", (e) => {
    e.preventDefault(); // Prevent default form submission behavior.
    console.log("Form submitted");

    // Clear previous results and show result container.
    resultContainer.innerHTML = "";
    formResults.classList.remove("hidden");

    const mac = macInput.value;
    const prefix = prefixInput.value;

    // Ensure WebAssembly validation function is available.
    if (typeof window.validateMAC !== "function") {
      resultContainer.innerHTML = `<p class="error-message">Error: WebAssembly module not loaded</p>`;
      resultContainer.classList.remove("hidden");
      return;
    }

    // Validate MAC address.
    let macErr = window.validateMAC(mac);
    if (macErr) {
      resultContainer.innerHTML = `<p class="error-message">Invalid MAC address (e.g., 00-14-22-01-23-45): ${macErr}</p>`;
      resultContainer.classList.remove("hidden");
      return;
    }

    // Validate IPv6 prefix.
    let prefixErr = window.validateIPv6Prefix(prefix);
    if (prefixErr) {
      resultContainer.innerHTML = `<p class="error-message">Invalid IPv6 prefix (e.g., 2001:db8::): ${prefixErr}</p>`;
      resultContainer.classList.remove("hidden");
      return;
    }

    // Calculate EUI-64 address.
    let result = window.calculateEUI64(mac, prefix);
    if (typeof result === "string") {
      resultContainer.innerHTML = `<p class="error-message">EUI-64 calculation failed: ${result}</p>`;
      resultContainer.classList.remove("hidden");
      return;
    }

    // Render result HTML with interface ID and full IPv6 address.
    resultContainer.innerHTML = `
      <div class="form-field-container">
        <label class="form-label" for="interface-id">End of IPv6 Address</label>
        <div class="input-copy-container">
          <input type="text" class="form-field" id="interface-id" readonly value="${result.interfaceID}" aria-describedby="interface-id-copy"/>
          <button class="copy-button" id="copy-interface" aria-label="Copy Interface ID">
            <svg class="copy-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
              <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
            </svg>
            <span class="copy-tooltip">Copy</span>
          </button>
        </div>
      </div>
      <br/>
      <div class="form-field-container">
        <label class="form-label" for="ip-full">IPv6 Address</label>
        <div class="input-copy-container">
          <input type="text" class="form-field" id="ip-full" readonly value="${result.fullIP}" aria-describedby="ip-full-copy"/>
          <button class="copy-button" id="copy-ip-full" aria-label="Copy IPv6 Address">
            <svg class="copy-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
              <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
            </svg>
            <span class="copy-tooltip">Copy</span>
          </button>
        </div>
      </div>
    `;
    resultContainer.classList.remove("hidden");

    // Attach event listeners to result copy buttons, ensuring no duplicates.
    const copyInterface = document.getElementById("copy-interface");
    const copyIpFull = document.getElementById("copy-ip-full");
    if (copyInterface) {
      copyInterface.removeEventListener("click", copyToClipboard);
      copyInterface.addEventListener("click", () =>
        copyToClipboard("interface-id", "copy-interface")
      );
    } else {
      console.error("copy-interface button not found");
    }
    if (copyIpFull) {
      copyIpFull.removeEventListener("click", copyToClipboard);
      copyIpFull.addEventListener("click", () =>
        copyToClipboard("ip-full", "copy-ip-full")
      );
    } else {
      console.error("copy-ip-full button not found");
    }
  });

  // Clear form results and hide containers on clear button click.
  document.querySelector(".form-clear").addEventListener("click", () => {
    resultContainer.innerHTML = "";
    formResults.classList.add("hidden");
    resultContainer.classList.add("hidden");
  });
});
