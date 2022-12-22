import { browser } from '$app/environment';
import { writable } from 'svelte/store';

const defaultValue: Map<string, string> = new Map();
const initialValue: Map<string, string> = browser ? new Map(Object.entries(JSON.parse(window.localStorage.getItem('player') || "{}"))) ?? defaultValue : defaultValue;

const playerStore = writable<Map<string, string>>(initialValue);

playerStore.subscribe((value) => {
    if (browser) {
        const json = JSON.stringify(Object.fromEntries(value));
        window.localStorage.setItem('player', json);
    }
});

export default playerStore;