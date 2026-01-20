import { writable } from 'svelte/store';

export const servers = writable([]);
export const expanded = writable(new Set());
export const selected = writable(null);

export function setServersFromArray(serverArray, expandAll) {
	if (expandAll) {
		const ids = [];
		for (const s of serverArray) {
			ids.push(s.id);
		}
		expanded.update(() => new Set(ids));
	}

	servers.set(serverArray);
}

export function toggleExpanded(id) {
	if (expanded.has(id)) {
		expanded.delete(id);
	} else {
		expanded.add(id);
	}
}
