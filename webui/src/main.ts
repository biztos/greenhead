/*

main.ts - Greenhead API *simple* UI code (SPA) entry point.

Ideas:

- separate *all* DOM logic into functions in another lib so here only bizlogic

*/

import "./style.css"; // weird-ass TS trick, required for vite packaging.

import { elem, hide, show, flash, disable, enable } from "./utils";
import { User } from "./user";
import { Agent, ToolCall } from "./agent";
import { API } from "./api";
import { marked } from "marked";
import sanitizeHtml from "sanitize-html";

function setupAll(): void {
  // Init objects and elements.
  const config = (window as any).__CONFIG__;
  if (config === null) {
    showError("Session not configured.", "Did you load from dist?");
    return;
  }
  const user = new User(
    config.user.api_key,
    config.user.name,
    config.user.agent_names,
  );
  const api = new API(user);

  let agent_select = elem("#newchat-agent-select") as HTMLSelectElement;
  api.user.agent_names.forEach((name) => {
    let opt = document.createElement("option");
    opt.value = name;
    opt.innerText = name;
    agent_select.appendChild(opt);
  });
  elem("#greeting").innerText = `Hello, ${api.user.name}!`;
  disable("#chat-cancel-button");

  // Add listeners.
  elem("#prompt-textarea").addEventListener("focus", function () {
    // NB: execCommand hack doesn't work cross-browser. Using textarea now.
    const ta = this as HTMLTextAreaElement;
    ta.select();
  });
  elem("#prompt-textarea").addEventListener("keydown", watchPrompt);
  elem("#prompt-send-button").addEventListener("click", runCompletion);
  elem("#error-dismiss-button").addEventListener("click", dismissError);
  elem("#newchat-button").addEventListener("click", newChat);
  elem("#reset-key-button").addEventListener("click", resetKey);
  elem("#chat-clear-button").addEventListener("click", clearHistory);
  elem("#chat-cancel-button").addEventListener("click", cancelCompletion);

  // Show our new-chat controls.
  show("#newchat");
}

async function runCompletion(): Promise<void> {
  // Make sure there's something useful to send.
  const ta = elem("#prompt-textarea") as HTMLTextAreaElement;
  const prompt = ta.value;
  if (!prompt.match(/\S/)) {
    flash("#prompt-textarea");
    return;
  }

  // Off you go!
  const api = API.getInstance();
  api.abort();

  // Prep UI.
  disable("#prompt-textarea");
  enable("#chat-cancel-button");
  await addUserPrompt(prompt);
  showProgress();

  // Send off the request!
  api.abortController = new AbortController();

  try {
    const response = await fetch(`/v1/agents/${api.agent!.id}/chat`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: "Bearer " + api.user.api_key,
      },
      body: JSON.stringify({ prompt: prompt }),
      signal: api.abortController.signal,
    });

    // Completed or failed, so reset UI elements:
    enable("#prompt-textarea");
    disable("#chat-cancel-button");
    hideProgress();
    api.abortController = undefined;

    if (!response.ok) {
      // A non-OK response goes into the error, but also we get a system msg.
      const errorText = await response.text();
      showError(`${response.status} ${response.statusText}`, errorText);
      addSystemMessage("Error running completion.");
      return;
    }

    // OK, right now we just have the response text, but this is crap.
    // TODO: get the tool calls in here so we can show you what happened!
    const data = await response.json();
    let calls: ToolCall[] = [];
    if (data.tool_calls != null) {
      for (const call of data.tool_calls) {
        calls.push(new ToolCall(call.id, call.name, call.args));
      }
    }
    await addCompletion(data.content, data.tool_calls);
    const ta = elem("#prompt-textarea") as HTMLTextAreaElement;
    ta.value = "";
    sizeTextArea(ta);
    ta.focus();
  } catch (err) {
    api.abortController = undefined;
    enable("#prompt-textarea");
    disable("#chat-cancel-button");
    hideProgress();
    // We may have been canceled!
    if (err instanceof Error && err.name === "AbortError") {
      addSystemMessage("User canceled completion.");
      return;
    }
    showError("Request failed.", err);
    addSystemMessage("Error running completion.");
  }
}

function addSystemMessage(message: string): void {
  // Ever care about HTML in a system message?
  addHistoryText(message, "system");
}

async function addUserPrompt(message: string): Promise<void> {
  addHistoryMarkdown(message, "user");
}

async function addCompletion(
  message: string,
  tool_calls: ToolCall[],
): Promise<void> {
  // Put in the tool calls first.  BYO Node for now.
  if (tool_calls.length > 0) {
    const tc = document.createElement("div");
    tc.classList.add("tool");
    const ol = document.createElement("ol");
    for (const call of tool_calls) {
      const li = document.createElement("li");
      li.innerText = call.name + " " + call.args; // args is string IRL
      ol.appendChild(li);
    }
    tc.appendChild(ol);
    elem("#history").appendChild(tc);
  }
  // TODO: treat agent completion as markdown *but safe*
  await addHistoryMarkdown(message, "agent");
}

async function addHistoryMarkdown(
  markdown: string,
  messageClass: string,
): Promise<void> {
  // NOTE: we actually do get code lang in the fenced block in the markdown
  // usually, but marked throws it away.  Look into fixing that and/or doing
  // a server render... but yucko, dislike including the html in the
  // response if we're KISS...
  const html = await marked.parse(markdown);
  const safe = sanitizeHtml(html);
  const div = document.createElement("div");
  div.classList.add(messageClass);
  div.innerHTML = safe;
  elem("#history").appendChild(div);
  scrollToBottom();
}

// TODO: addHistoryHTML
function addHistoryText(messageText: string, messageClass: string): void {
  const div = document.createElement("div");
  div.classList.add(messageClass);
  div.innerText = messageText;
  elem("#history").appendChild(div);
  scrollToBottom();
}

function clearHistory(): void {
  elem("#history").innerHTML = "";
}

function cancelCompletion(): void {
  // Abort will cause an error in the fetch, which we catch.
  API.getInstance().abort();
}

// end the existing agent, throwing error on any non-OK non-404
async function endAgent(): Promise<void> {
  const api = API.getInstance();
  if (api.agent == undefined) {
    return;
  }
  const response = await fetch(`/v1/agents/${api.agent.id}/end`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: "Bearer " + api.user.api_key,
    },
    body: JSON.stringify({}), // no payload at this time!
  });
  if (response.ok || response.status == 404) {
    api.agent = undefined;
    return;
  }
  throw new Error(`${response.status} ${response.statusText}`);
}

async function newChat(): Promise<void> {
  // Bit awkward here, but don't let them easily just nuke the history.
  const user_hist = document.querySelectorAll("#history .user");
  if (user_hist.length > 0) {
    let do_clear = confirm("Clear chat history and start fresh?");
    if (do_clear == false) {
      return;
    }
  }
  cancelCompletion();

  const api = API.getInstance();

  // Hide current chat and freeze our own controls.
  hide("#chat");
  disable("#newchat-agent-select");
  disable("#newchat-button");

  // With whom shall we be chatting?
  const sel = elem("#newchat-agent-select") as HTMLSelectElement;
  const agent_name = sel.value;

  try {
    showProgress();
    await endAgent();
    const response = await fetch("/v1/agents/new", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: "Bearer " + api.user.api_key,
      },
      body: JSON.stringify({ agent: agent_name }),
    });

    enable("#newchat-agent-select");
    enable("#newchat-button");
    hideProgress();

    if (!response.ok) {
      if (api.agent) {
        show("#chat");
      }
      const errorText = await response.text();
      showError(`${response.status} ${response.statusText}`, errorText);
      return;
    }

    const data = await response.json();
    const agent = new Agent(data.id, data.name, data.description);
    api.agent = agent;
    startChat(agent);
  } catch (err) {
    enable("#newchat-agent-select");
    enable("#newchat-button");
    if (api.agent) {
      show("#chat");
    }
    hideProgress();
    showError("Request failed.", err);
  }
}

function startChat(agent: Agent): void {
  // Delete previous chat. (Maybe save eventually but not yet.)
  // TODO: it would be nice to change agents without nuking the chat!
  // But it's not at all clear you can clone the context with tool calls,
  // for openai you can not, it knows what was called and what not.
  // Anyway for now just nuke it.
  clearHistory();
  let msg = `You are chatting with ${agent.name}`;
  if (agent.description == "") {
    msg += ".";
  } else {
    // TODO: safe markdown in description!
    msg += `: ${agent.description}`;
  }
  addSystemMessage(msg);

  // NB: do NOT reset the prompt, user might want to reuse it!
  // TODO: figure out what IN THE HAYOW is causing the ta to shrink on newChat
  // (it's not the hide/show, already checked that... sth about adding to DOM
  // maybe...)
  const ta = elem("#prompt-textarea") as HTMLTextAreaElement;
  setTimeout(() => {
    sizeTextArea(ta);
  }, 1);

  // Show the chat area.
  show("#chat");
}

async function resetKey(event: Event): Promise<void> {
  event.preventDefault();
  disable("#prompt-textarea");
  disable("#newchat-button");
  let do_reset = confirm("Clear chat history and reset API Key?");
  if (do_reset == false) {
    enable("#prompt-textarea");
    enable("#newchat-button");
    return;
  }
  try {
    await endAgent(); // let's not dangle!
    window.location.href = "/v1/ui"; // quite the stupid hack but OK for now.
  } catch (err) {
    enable("#prompt-textarea");
    enable("#newchat-button");
    showError("Failed to end current agent.", err);
  }
}

function showError(message: string, err: any): void {
  const detail = err instanceof Error ? err.message : String(err);
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
  const prev_height = ta.style.height;
  ta.style.height = "auto";

  const computed = window.getComputedStyle(ta);
  const paddingTop = parseFloat(computed.paddingTop);
  const paddingBottom = parseFloat(computed.paddingBottom);

  const totalPadding = paddingTop + paddingBottom;

  const new_height = ta.scrollHeight - totalPadding + "px";
  // Does this make sense?!
  ta.style.height = new_height;
  ta.style.height = ta.scrollHeight + 10 + "px";

  // Sizing/scrolling still a total mess, TODO: figure that out!
  // (in SwiftUI I'd do a view stack of some kind...)
  if (prev_height != new_height) {
    // TODO: something with scrolling but so far it's a shit-show...
  }
}

function scrollToBottom(): void {
  window.scrollTo({
    top: document.documentElement.scrollHeight,
    behavior: "smooth",
  });
}

function watchPrompt(event: KeyboardEvent): void {
  if (event.defaultPrevented) {
    return;
  }

  let ta = event.target! as HTMLTextAreaElement;

  // Size the text area in any case, 'tis annoying but somehow we have to.
  sizeTextArea(ta);

  // Ignore if not set.
  const chk = elem("#prompt-return-sends-checkbox") as HTMLInputElement;
  if (!chk.checked) {
    return;
  }

  if (event.key === "Enter") {
    if (!event.shiftKey) {
      event.preventDefault();
      runCompletion();
    }
  }
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
