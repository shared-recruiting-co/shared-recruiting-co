export const debounce = (func: (...args: any[]) => void, wait: number) => {
  let debounceTimeout: NodeJS.Timeout;
  return function executedFunction(...args: any[]) {
    clearTimeout(debounceTimeout);
    debounceTimeout = setTimeout(() => func(...args), wait);
  };
};

export type FormErrors = Record<string, string>;

export const formError = (e: typeof FormErrors, field: string) => {
  return e[field] || "";
};

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
  if (typeof value === "string") {
    return value.trim();
  }
  return "";
};

export const getFormCheckboxValue = (data: FormData, key: string): boolean => {
  const value = data.get(key);
  return value === "on";
};
