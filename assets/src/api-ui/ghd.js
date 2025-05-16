/*

ghd.js

NOTES:

1. this is experimental, simple, and a little dumb -- by design!
2. (it might turn into IZK playground later)
3. no complex sessions/histories YET, just one chat in a window

*/

let LAST_XHR = null; // examine in dev tools.
let CURRENT_XHR = null; // in-flight XHR, move to above when done.
let AGENT_ID = null; // current-chat agent ID.
let API_KEY = null; // defined in the HTML by the server.
let AGENT_NAMES = []; // same.

document.addEventListener("DOMContentLoaded", function () {
  let agent_select = document.getElementById("agent-select");
  AGENT_NAMES.forEach((name) => {
    let opt = document.createElement("option");
    opt.value = name;
    opt.innerText = name;
    agent_select.appendChild(opt);
  });

  document.getElementById("prompt").addEventListener("focus", function (e) {
    // "deprecated" but: https://stackoverflow.com/a/70831583
    document.execCommand("selectAll", false, null);
  });
  document.getElementById("prompt").addEventListener("keydown", watchPrompt);
  document.getElementById("newchat-button").addEventListener("click", newChat);
  document
    .getElementById("reset-key-button")
    .addEventListener("click", resetKey);
});

let newChat = function (event) {

  // Deal with in-flight chat, should be OK to just nuke it.
  if (CURRENT_XHR != null) {
    CURRENT_XHR.abort();
  }

  // Hide current chat; we don't nuke it until we successfully created a new
  // chat, because in the failure case you might still want the old one.
  hide("#history");
  hide("#error");

  // Freeze new-chat form and show progress.
  element("#newchat-button").disabled = true;
  element("#agent-select").disabled = true;
  showProgress();
  
  // Set up the XHR.
  let xhr = new XMLHttpRequest();
  xhr.open("POST", "/v1/agents/new", true);
  xhr.setRequestHeader("Content-Type", "application/json");
  xhr.setRequestHeader("Authoriation", "Bearer " + API_KEY);
  xhr.onreadystatechange = function () {
  if (xhr.readyState === XMLHttpRequest.DONE) {

    element("#newchat-button").disabled = false;
    element("#agent-select").disabled = false;
    if (xhr.status === 200) {
      console.log("Response:", xhr.responseText);
      AGENT_ID = xhr.
    } else {
      if (AGENT_ID != null) {
        // had a chat before, put it back for now.
        show("#history");
        show("#prompt");
      }
      element("#error").innerText = xhr.status;
    }
  }
};

  xhr.send(JSON.stringify({ agent: element("#agent-select").value }));
  CURRENT_XHR = xhr;

  
};

let resetKey = function (event) {
  event.preventDefault();
  let do_reset = confirm("Clear chat history and reset API Key?");
  if (do_reset == true) {
    window.location = "/v1/ui"; // quite the stupid hack but OK for now.
  }
};

let watchPrompt = function (event) {
  if (event.defaultPrevented) {
    return;
  }
  if (event.key === "Enter") {
    if (!event.shiftKey) {
      event.preventDefault();
      // Grab the text and make sure it's got something in it.
      let p = event.target;
      let s = p.innerText;
      if (!s.match(/\S/)) {
        flashElem(p);
        return;
      }
      // Off you go!
      sendPrompt();
    }
  }
};

let sendPrompt = function () {
  if (IN_FLIGHT) {
    console.log("already in flight");
    return;
  }

  IN_FLIGHT = true;
  console.log("sending prompt");

  let p = document.getElementById("prompt");
  let s = p.innerText.trim();
  let h = document.getElementById("history");

  // Freeze and clear the prompt input.
  p.contentEditable = false;
  deselectAll(); // or it looks wonky.
  p.innerText = "";

  // Add the prompt text to the history.
  let div = document.createElement("div");
  div.classList.add("prompt");
  div.innerText = s;
  h.appendChild(div);

  // Show the progress indicator.
  showProgress();
};

let addPromptText = function (s) {};

// NB: total hack from one of the AIs here, at least at the beginning.
let showProgress = function () {
  unhide("#progress");
  let svg = element("#spinner");
  let duration = 3; // TBD, what looks nice?

  // Add CSS for the rotation animation
  svg.style.transformOrigin = "center";
  svg.style.animation = `spin ${duration}s linear infinite`;

  // Add the keyframes animation to the document
  const style = document.createElement("style");
  style.textContent =
    "@keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }";
  document.head.appendChild(style);
};

let hideProgress = function () {
  // Might be cool to freeze the spinner first but this'll do.
  document.getElementById("progress").classList.add("hidden");

  let svg = document.getElementById("spinner");

  const computedStyle = window.getComputedStyle(svg);
  const currentTransform = computedStyle.getPropertyValue("transform");

  // Remove the animation
  svg.style.animation = "none";

  // Keep the current rotation position
  if (currentTransform !== "none") {
    svg.style.transform = currentTransform;
  }
};

let flashElem = function (elem) {
  elem.classList.add("attn");

  setTimeout(() => {
    elem.classList.remove("attn");
  }, 100);
};

let deselectAll = function () {
  // https://stackoverflow.com/a/6562764
  if (window.getSelection) {
    window.getSelection().removeAllRanges();
  } else if (document.selection) {
    document.selection.empty();
  }
};

// can't believe we're still doing this shit in 2025...
let hide = function (elem) {
  element(elem).classList.add("hidden");
};

// and this!
let unhide = function (elem) {
  element(elem).classList.remove("hidden");
};

// and this! not even using $ because the debugger aliases that!
let element = function (elem) {
  if (typeof elem === "string") {
    let selected = document.querySelector(elem);
    if (selected == null) {
      throw "element not found by selector: " + elem;
    }
    return selected;
  } else {
    return elem;
  }
};
