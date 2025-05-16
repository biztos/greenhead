import { elem } from "./utils";

export class User {
  private static _instance: User;

  readonly api_key: string;
  readonly name: string;
  readonly agent_names: string[];

  static initFromDOM(): User {
    const apiKeyInput = document.querySelector(
      "#user-api-key",
    ) as HTMLInputElement;
    const nameInput = elem("#user-name") as HTMLInputElement;

    // Effing Typescript implements HTMLCollection wrong, can't iterate, WTF?!
    const agentNamesSelect = elem("#user-agent-name") as HTMLSelectElement;

    let agents: string[] = [];
    for (const opt of Array.from(agentNamesSelect.children)) {
      const option = opt as HTMLOptionElement;
      agents.push(option.value);
    }
    return new User(apiKeyInput.value, nameInput.value, agents);
  }

  constructor(api_key: string, name: string, agent_names: string[]) {
    if (User._instance) {
      throw new Error("Error: User singleton already initialized.");
    }
    this.api_key = api_key;
    this.name = name;
    this.agent_names = agent_names; // This was missing
    User._instance = this;
  }

  // Add a method to get the instance
  static getInstance(): User {
    if (!User._instance) {
      throw new Error("Error: User singleton not initialized.");
    }
    return User._instance;
  }
}
