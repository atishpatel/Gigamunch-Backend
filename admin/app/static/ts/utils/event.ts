export const UserUpdated = 'UserUpdated';


export function Fire(eventName: string, detail: Object = {}) {
  const event = new CustomEvent(eventName, {
    detail,
    bubbles: true,
    composed: true,
  });
  window.dispatchEvent(event);
}

export function FireTost(detail: Object) {
  const event = new CustomEvent('toast', {
    detail,
    bubbles: true,
    composed: true,
  });
  window.dispatchEvent(event);
}

export function FireError() {

}
