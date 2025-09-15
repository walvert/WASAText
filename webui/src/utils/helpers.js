export function getUserInitials(username) {
	if (!username) return '?';

	const words = username.split(' ');
	if (words.length >= 2) {
		return (words[0].charAt(0) + words[1].charAt(0)).toUpperCase();
	}
	return username.charAt(0).toUpperCase();
}

export function formatMessageTime(timestamp) {
	if (!timestamp) return '';

	const date = new Date(timestamp);
	const now = new Date();

	if (date.toDateString() === now.toDateString()) {
		return date.toLocaleTimeString([], {
			hour: '2-digit',
			minute: '2-digit',
			hour12: false
		});
	}

	const yesterday = new Date(now);
	yesterday.setDate(now.getDate() - 1);
	if (date.toDateString() === yesterday.toDateString()) {
		return 'Yesterday ' + date.toLocaleTimeString([], {
			hour: '2-digit',
			minute: '2-digit',
			hour12: false
		});
	}

	const weekAgo = new Date(now);
	weekAgo.setDate(now.getDate() - 7);
	if (date > weekAgo) {
		return date.toLocaleDateString([], {weekday: 'short'}) + ' ' +
			date.toLocaleTimeString([], {hour: '2-digit', minute: '2-digit', hour12: false});
	}

	if (date.getFullYear() === now.getFullYear()) {
		return date.toLocaleDateString([], {month: 'short', day: 'numeric'}) + ' ' +
			date.toLocaleTimeString([], {hour: '2-digit', minute: '2-digit', hour12: false});
	}

	return date.toLocaleDateString() + ' ' +
		date.toLocaleTimeString([], {hour: '2-digit', minute: '2-digit', hour12: false});
}
