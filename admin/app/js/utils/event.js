export var UserUpdated = 'UserUpdated';
export function Fire(eventName, detail) {
    if (detail === void 0) { detail = {}; }
    var event = new CustomEvent(eventName, {
        detail: detail,
        bubbles: true,
        composed: true,
    });
    window.dispatchEvent(event);
}
export function FireToast(t, detail) {
    var event = new CustomEvent('toast', {
        detail: detail,
        bubbles: true,
        composed: true,
    });
    t.dispatchEvent(event);
}
export function FireError() {
}
