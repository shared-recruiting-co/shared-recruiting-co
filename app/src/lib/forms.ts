export const isValidUrl = (str: string): boolean => {
	try {
		new URL(str);
		return true;
	} catch (error) {
		// If the URL is invalid, an error will be thrown
		return false;
	}
};

export const getTrimmedFormValue = (data: FormData, key: string): string => {
	const value = data.get(key);
	if (typeof value === 'string') {
		return value.trim();
	}
	return '';
};

export const getFormCheckboxValue = (data: FormData, key: string): boolean => {
	const value = data.get(key);
	return value === 'on';
};
