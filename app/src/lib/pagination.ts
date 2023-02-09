export const getPaginationPages = (current: number, total: number): string[] => {
	if (total <= 5) {
		return Array.from({ length: total }, (_, i) => (i + 1).toString());
	}
	switch (current) {
		case 1:
			return ['1', '2', '...', (total - 1).toString(), total.toString()];
		case 2:
			return ['1', '2', '3', '...', (total - 1).toString(), total.toString()];
		case 3:
			// special case for 6
			if (total === 6) {
				return ['1', '2', '3', '4', '5', '6'];
			}
			return ['1', '2', '3', '4', '...', (total - 1).toString(), total.toString()];
		case total - 2:
			// special case for 6
			if (total === 6) {
				return ['1', '2', '3', '4', '5', '6'];
			}
			return [
				'1',
				'2',
				'...',
				(total - 3).toString(),
				(total - 2).toString(),
				(total - 1).toString(),
				total.toString()
			];
		case total - 1:
			return ['1', '2', '...', (total - 2).toString(), (total - 1).toString(), total.toString()];
		case total:
			return ['1', '2', '...', (total - 1).toString(), total.toString()];
		default:
			return [
				'1',
				'2',
				...(current - 1 == 3 ? [] : ['...']),
				(current - 1).toString(),
				current.toString(),
				(current + 1).toString(),
				...(current + 1 == total - 2 ? [] : ['...']),
				(total - 1).toString(),
				total.toString()
			];
	}
};
