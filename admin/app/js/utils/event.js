export const UserUpdated = 'UserUpdated';
export function Fire(eventName, detail = {}) {
    const event = new CustomEvent(eventName, {
        detail,
        bubbles: true,
        composed: true,
    });
    window.dispatchEvent(event);
}
export function FireToast(t, detail) {
    const event = new CustomEvent('toast', {
        detail,
        bubbles: true,
        composed: true,
    });
    t.dispatchEvent(event);
}
export function FireError() {
}
