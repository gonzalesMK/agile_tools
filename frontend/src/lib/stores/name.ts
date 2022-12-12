import { browser } from '$app/environment';
import { writable } from 'svelte/store';

const defaultValue = '';
const initialValue = browser ? window.localStorage.getItem('name') ?? defaultValue : defaultValue;

const nameStore = writable<string>(initialValue);

nameStore.subscribe((value) => {
    if (browser) {
        window.localStorage.setItem('name', value);
    }
});

export default nameStore;