export function Fire(eventName: string, detail: Object = {}) {
  const event = new CustomEvent(eventName, {
    detail,
    bubbles: true,
    composed: true,
  });
  window.dispatchEvent(event);
}

// FirstToast dispatches a 'toast' event. 
// First parameter is element doing the dispatch and second param is the detail of the event.
export function FireToast(t: Element, detail: Object) {
  const event = new CustomEvent('toast', {
    detail,
    bubbles: true,
    composed: true,
  });
  t.dispatchEvent(event);
}

export function FireError() {

}


// removes composed error
declare global {
  interface CustomEventInit {
    readonly composed: boolean;
  }
}