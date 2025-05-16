/**
 * Represents an Agent entity within the system.
 *
 */
export class Agent {
  /** Unique identifier for the agent */
  readonly id: string;

  /** Display name of the agent */
  readonly name: string;

  /** Descriptive text about the agent's purpose or capabilities */
  readonly description: string;

  /**
   * Creates a new Agent instance.
   *
   * @param id - Unique identifier for the agent
   * @param name - Display name of the agent
   * @param description - Descriptive text about the agent's purpose or capabilities
   */
  constructor(id: string, name: string, description: string) {
    this.id = id;
    this.name = name;
    this.description = description;
  }

  /**
   * Creates a new Agent instance from a JSON string which must include the
   * properties (id,name,description).
   *
   * @param json - JSON representation of the Agent.
   * @throws {Error} If the JSON is missing any required properties (id, name, description)
   * @returns A new Agent instance created from the JSON data
   */
  static fromJSON(json: string): Agent {
    // Parse JSON string into an object
    const data = JSON.parse(json);

    // Check for required properties
    if (!data.id) {
      throw new Error("Missing required property: id");
    }
    if (!data.name) {
      throw new Error("Missing required property: name");
    }
    if (!data.description) {
      throw new Error("Missing required property: description");
    }

    // Create and return a new Agent instance
    return new Agent(data.id, data.name, data.description);
  }
}
