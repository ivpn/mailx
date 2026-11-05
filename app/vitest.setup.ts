// Node 25+ defines its own localStorage/sessionStorage globals, which are
// undefined unless Node runs with --localstorage-file. Vitest's jsdom
// environment does not replace globals Node already defines, so tests would
// see Node's storage instead of jsdom's. Point both at jsdom's storage.
const { window: domWindow } = (globalThis as unknown as { jsdom: { window: Window } }).jsdom

for (const key of ['localStorage', 'sessionStorage'] as const) {
    Object.defineProperty(globalThis, key, {
        value: domWindow[key],
        configurable: true,
        writable: true,
    })
}
