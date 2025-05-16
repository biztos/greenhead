/*

main.ts - Greenhead API *simple* UI code (SPA) entry point.

*/

import "./style.css";

import { elem, hide, show, flash, disable, enable } from "./utils";
// import { User } from "./user";
import { Agent } from "./agent";
import { API } from "./api";

function setupAll(): void {
  // Bail out if we're not in the normal session context.
  if (document.querySelector("#user-session") == null) {
    showError("Session not initiated.", "Did you load from dist?");
    hide("#error-dismiss-button");
    return;
  }

  // Init objects and elements.
  const api = API.initFromDOM();
  let agent_select = elem("#newchat-agent-select") as HTMLSelectElement;
  api.user.agent_names.forEach((name) => {
    let opt = document.createElement("option");
    opt.value = name;
    opt.innerText = name;
    agent_select.appendChild(opt);
  });
  elem("#greeting").innerText = `Hello, ${api.user.name}!`;

  // Add listeners.

  elem("#prompt").addEventListener("focus", function () {
    // NB: execCommand hack doesn't work cross-browser. Using textarea now.
    const ta = this as HTMLTextAreaElement;
    ta.select();
  });
  elem("#prompt").addEventListener("keydown", watchPrompt);
  elem("#error-dismiss-button").addEventListener("click", dismissError);
  elem("#newchat-button").addEventListener("click", newChat);
  elem("#reset-key-button").addEventListener("click", resetKey);

  // Show our new-chat controls.
  show("#newchat");
}

async function newChat(): Promise<void> {
  // Bail on any currently in-flight activity, don't care for now.
  const api = API.getInstance();
  api.xhr?.abort();
  api.xhr = undefined;

  // Hide current chat and freeze our own controls.
  hide("#chat");
  disable("#newchat-agent-select");
  disable("#newchat-button");

  // With whom shall we be chatting?
  const sel = elem("#newchat-agent-select") as HTMLSelectElement;
  const agent_name = sel.value;

  try {
    showProgress();
    const response = await fetch("/v1/agents/new", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer " + api.user.api_key,
      },
      body: JSON.stringify({ agent: agent_name }),
    });

    enable("#newchat-agent-select");
    enable("#newchat-button");
    hideProgress();

    if (!response.ok) {
      if (api.agent != null) {
        show("#chat");
      }
      const errorText = await response.text();
      showError(`${response.status} ${response.statusText}`, errorText);
      return;
    }

    const val = await response.json();
    console.log(typeof(val));
    console.log(val);
    const agent = new Agent(val.id,val.name,val.description);
    api.agent = agent;
    startChat(agent);
  } catch (err) {
    enable("#newchat-agent-select");
    enable("#newchat-button");
    hideProgress();
    showError("Request failed", err instanceof Error ? err.message : String(err));
  }

}


function startChat(agent: Agent): void {

  // Delete previous chat. (Maybe save eventually but not yet.)
  // TODO: it would be nice to change agents without nuking the chat!
  // But it's not at all clear you can clone the context with tool calls,
  // for openai you can not, it knows what was called and what not.
  // Anyway for now just nuke it.
  elem("#history").innerHTML = "";
  const msg = document.createElement("div");
  msg.classList.add("system");
  if (agent.description == "") {
    msg.innerText = `You are chatting with ${agent.name}.`
  } else {
    msg.innerText = `You are chatting with ${agent.name}: ${agent.description}.`
  }
  elem("#history").appendChild(msg);

  // NB: do NOT reset the prompt, user might want to reuse it!

  // Show the chat area.
  show("#chat");

}

function resetKey(event: Event): void {
  event.preventDefault();
  let do_reset = confirm("Clear chat history and reset API Key?");
  if (do_reset == true) {
    window.location.href = "/v1/ui"; // quite the stupid hack but OK for now.
  }
}

function showError(message: string, detail: string): void {
  elem("#error-message").innerText = message;
  elem("#error-detail").innerText = detail;
  show("#error");
}

function dismissError(): void {
  // reset the values too, just to be thorough.
  elem("#error-message").innerText = "Error.";
  elem("#error-detail").innerText = "An error occurred.";
  hide("#error");
}

function sizeTextArea(ta: HTMLTextAreaElement): void {
  // Wonky-ass size adjustment for the TA.
  ta.style.height = "auto";

  const computed = window.getComputedStyle(ta);
  const paddingTop = parseFloat(computed.paddingTop);
  const paddingBottom = parseFloat(computed.paddingBottom);

  const totalPadding = paddingTop + paddingBottom;
  ta.style.height = ta.scrollHeight - totalPadding + "px";

  ta.style.height = ta.scrollHeight + 10 + "px";
}

function watchPrompt(event: KeyboardEvent): void {
  if (event.defaultPrevented) {
    return;
  }

  let ta = event.target! as HTMLTextAreaElement;

  sizeTextArea(ta);

  if (event.key === "Enter") {
    if (!event.shiftKey) {
      event.preventDefault();
      // Grab the text and make sure it's got something in it.
      let s = ta.value;
      if (!s.match(/\S/)) {
        flash(ta);
        return;
      }
      // Off you go!
      sendPrompt(s);
    }
  }
}

function sendPrompt(prompt: string): void {
  console.log("SEND PROMPT");
  console.log(prompt);
}

function showProgress(): void {
  show("#progress");
  let svg = elem("#spinner");
  let duration = 3; // TBD, what looks nice?

  // Add CSS for the rotation animation
  svg.style.transformOrigin = "center";
  svg.style.animation = `spin ${duration}s linear infinite`;
}

function hideProgress(): void {
  // Might be cool to freeze the spinner first but this'll do.
  hide("#progress");

  let svg = elem("#spinner");

  const computedStyle = window.getComputedStyle(svg);
  const currentTransform = computedStyle.getPropertyValue("transform");

  // Remove the animation
  svg.style.animation = "none";

  // Keep the current rotation position
  if (currentTransform !== "none") {
    svg.style.transform = currentTransform;
  }
}

document.addEventListener("DOMContentLoaded", setupAll);
