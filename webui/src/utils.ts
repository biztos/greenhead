/* utils.ts - I can't believe I'm writing this in 2025! 🤡

*/

type elemSel = string | HTMLElement;

/**
 * Briefly adds and removes a "flash" class to an element to create a
 * flash effect, removing it soon after.
 *
 * @param selector A CSS selector string or an HTMLElement
 */
export function flash(sel: elemSel): void {
  const element = elem(sel);
  element.classList.add("flash");
  setTimeout(() => {
    element.classList.remove("flash");
  }, 100);
}

export function deselectAll(): void {
  // https://stackoverflow.com/a/6562764
  if (window.getSelection) {
    window.getSelection()?.removeAllRanges();
  } else {
    console.log("sorry, TS won't let us fall back to document.selection");
  }
}

/**
 * Hides an HTMLElement by adding a "hidden" CSS class.
 *
 * @param sel A CSS selector string or an HTMLElement
 */
export function hide(sel: elemSel): void {
  elem(sel).classList.add("hidden");
}

/**
 * Unhides (shows) an HTMLElement by removing the "hidden" CSS class.
 *
 * @param sel A CSS selector string or an HTMLElement
 */
export function show(sel: elemSel): void {
  elem(sel).classList.remove("hidden");
}

/**
 * Returns an HTMLElement either by selecting it with a CSS selector or
 * returning the element directly.
 *
 * @param sel A CSS selector string or an HTMLElement
 * @returns The found or provided HTMLElement
 * @throws Error if the element cannot be found by the provided selector
 */
export function elem(sel: elemSel): HTMLElement {
  if (typeof sel === "string") {
    const selected = document.querySelector(sel as string);
    if (selected === null) {
      throw new Error(`Element not found by selector: ${sel}`);
    }
    return selected as HTMLElement;
  } else {
    return sel as HTMLElement;
  }
}

type HTMLAbleElement =
  | HTMLButtonElement
  | HTMLInputElement
  | HTMLSelectElement
  | HTMLTextAreaElement
  | HTMLOptionElement
  | HTMLOptGroupElement
  | HTMLFieldSetElement;

/**
 * Sets the disabled property of an element to true.
 *
 * The element must be castable to HTMLAbleElement.
 *
 * @param sel A CSS selector string or an HTMLElement
 */
export function disable(sel: elemSel): void {
  const d = elem(sel) as HTMLAbleElement;
  d.disabled = true;
}

/**
 * Sets the disabled property of an element to false.
 *
 * The element must be castable to HTMLAbleElement.
 *
 * @param sel A CSS selector string or an HTMLElement
 */
export function enable(sel: elemSel): void {
  const d = elem(sel) as HTMLAbleElement;
  d.disabled = false;
}
