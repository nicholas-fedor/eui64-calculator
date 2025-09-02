// Initializes WebAssembly module and handles form interactions for the EUI-64 calculator static site.
const go = new Go();
WebAssembly.instantiateStreaming(fetch("./main.wasm"), go.importObject)
  .then((result) => {
    go.run(result.instance);
    console.log("WebAssembly module loaded successfully");
    // Verify WebAssembly functions are available
    if (
      typeof window.validateMAC !== "function" ||
      typeof window.validateIPv6Prefix !== "function" ||
      typeof window.calculateEUI64 !== "function"
    ) {
      console.error("WebAssembly functions not available");
    }
  })
  .catch((err) => console.error("WASM load error:", err));

// copyToClipboard copies the value of an input element to the clipboard and updates the button's tooltip.
// It matches the behavior in ui/layout.templ, displaying "Copied!" for 2 seconds on success.
function copyToClipboard(elementId, buttonId) {
  const input = document.getElementById(elementId);
  const button = document.getElementById(buttonId);
  if (!input || !button) {
    console.error(
      `copyToClipboard: Element ${elementId} or ${buttonId} not found`
    );
    return;
  }
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
    .catch((err) => {
      console.error("Failed to copy: ", err);
    });
}

// Initializes form event listeners for submission and clearing, handling validation and calculation via WASM.
document.addEventListener("DOMContentLoaded", () => {
  const form = document.querySelector("form");
  const resultContainer = document.querySelector(".result-container");
  const formResults = document.querySelector(".form-results");
  const macInput = document.getElementById("mac");
  const prefixInput = document.getElementById("ip-start");
  const copyMac = document.getElementById("copy-mac");
  const copyPrefix = document.getElementById("copy-ip-start");

  if (
    !form ||
    !resultContainer ||
    !formResults ||
    !macInput ||
    !prefixInput ||
    !copyMac ||
    !copyPrefix
  ) {
    console.error("Required DOM elements not found");
    return;
  }

  console.log("Form event listeners initialized");

  copyMac.addEventListener("click", () => copyToClipboard("mac", "copy-mac"));
  copyPrefix.addEventListener("click", () =>
    copyToClipboard("ip-start", "copy-ip-start")
  );

  form.addEventListener("submit", (e) => {
    e.preventDefault(); // Prevent default form submission
    console.log("Form submitted");
    resultContainer.innerHTML = "";
    formResults.classList.remove("hidden");

    const mac = macInput.value;
    const prefix = prefixInput.value;

    // Check if WebAssembly functions are available
    if (typeof window.validateMAC !== "function") {
      resultContainer.innerHTML = `<p class="error-message">Error: WebAssembly module not loaded</p>`;
      resultContainer.classList.remove("hidden");
      return;
    }

    let macErr = window.validateMAC(mac);
    if (macErr) {
      resultContainer.innerHTML = `<p class="error-message">Please enter a valid MAC address (e.g., 00-14-22-01-23-45): ${macErr}</p>`;
      resultContainer.classList.remove("hidden");
      return;
    }
    let prefixErr = window.validateIPv6Prefix(prefix);
    if (prefixErr) {
      resultContainer.innerHTML = `<p class="error-message">Please enter a valid IPv6 prefix (e.g., 2001:db8::): ${prefixErr}</p>`;
      resultContainer.classList.remove("hidden");
      return;
    }

    let result = window.calculateEUI64(mac, prefix);
    if (typeof result === "string") {
      resultContainer.innerHTML = `<p class="error-message">Failed to calculate EUI-64 address: ${result}</p>`;
      resultContainer.classList.remove("hidden");
      return;
    }

    // Render result HTML, matching ui/result.templ structure and behavior.
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
                    <script>
                        document.getElementById("copy-interface").addEventListener("click", () => {
                            copyToClipboard("interface-id", "copy-interface");
                        });
                    </script>
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
                    <script>
                        document.getElementById("copy-ip-full").addEventListener("click", () => {
                            copyToClipboard("ip-full", "copy-ip-full");
                        });
                    </script>
                </div>
            </div>
        `;
    resultContainer.classList.remove("hidden");
  });

  document.querySelector(".form-clear").addEventListener("click", () => {
    resultContainer.innerHTML = "";
    formResults.classList.add("hidden");
    resultContainer.classList.add("hidden");
  });
});
