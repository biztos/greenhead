export class User {
  private static _instance: User;

  readonly api_key: string;
  readonly name: string;
  readonly agent_names: string[];

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
