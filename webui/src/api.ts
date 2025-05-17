import { User } from "./user";
import { Agent } from "./agent";

export class API {
  private static _instance: API;

  readonly user: User;
  agent: Agent | undefined;
  abortController: AbortController | undefined;

  static initFromDOM(): API {
    const user = User.initFromDOM();
    return new API(user);
  }

  constructor(user: User) {
    if (API._instance) {
      throw new Error("Error: API singleton already initialized.");
    }
    this.user = user;
    API._instance = this;
  }

  static getInstance(): API {
    if (!API._instance) {
      throw new Error("Error: API singleton not initialized.");
    }
    return API._instance;
  }

  abort(): void {
    this.abortController?.abort();
    this.abortController = undefined;
  }
}
