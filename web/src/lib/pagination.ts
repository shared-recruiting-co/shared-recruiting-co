const DEFAULT_PAGE_SIZE = 10;

/**
 * Helper function to return the previous or next page url for the below pagination function
 *
 * @param {URL} currentPageUrl - The URL of the current page
 * @param {number} currentResultsPage - The current page number
 * @param {boolean} updatedPageValid - Whether the updated page is valid or not
 * @param {boolean} updatedPageIsPrev - True if the desired url is the previous page, otherwise next
 * @returns {string|null} - The updated page URL, or null if the updated page is not valid
 */
const getUpdatedPageUrl = (currentPageUrl, currentResultsPage, updatedPageValid, updatedPageIsPrev) => {

	// set the udpate page URL to null if not valid, otherwise set it to the current url
	let updatedPageUrl = updatedPageValid ? new URL(currentPageUrl.href) : null;

	// if the page is valid update the newly created url to be the new page url
	if (updatedPageValid) {

		// if were getting the previous page we subtract 1, otherwise add 1
		const pageParamModifier = (updatedPageIsPrev) ? -1 : 1

		// get the updated page URL
		updatedPageUrl.searchParams.set('page', currentResultsPage + pageParamModifier);
		updatedPageUrl = updatedPageUrl.href
	}

	// return the updated page URL
	return updatedPageUrl
}

export type Pagination = {
	currentResultsPage: number;
	resultsPerPage: number;
	resultsToFetchStart: number;
	resultsToFetchEnd: number;
	resultShowingFirst: number;
	resultShowingLast: number;
	resultsCount: number;
	pagesDisplay: Array<string>
	pagesCount: number;
	prevPageValid: boolean;
	prevPageUrl: string | null;
	nextPageValid: boolean;
	nextPageUrl: string | null;
};

/**
 * Returns an object representing pagination information for the given inputs
 *
 * @param {URL} url - The URL of the current page
 * @param {number} resultsCount - The total number of results fetched from the DB
 * @param {number} resultsPerPage - The number of results to display per page (default above)
 * @returns {Pagination} - An object representing pagination information for the search results
 */
export const getPagePagination = (url: URL, resultsCount: number, resultsPerPage = DEFAULT_PAGE_SIZE ): Pagination => {

	// from the given URL, grab which page of results is currently selected (default 1)
	const currentResultsPage = parseInt(url.searchParams.get('page')) || 1;

	// define the start / stop range of rows we want to retreive from the DB for the page
	const resultsToFetchStart = (currentResultsPage - 1) * resultsPerPage;
	const resultsToFetchEnd = resultsToFetchStart + resultsPerPage;

	// of the results, the first / last displayed on the page
	const resultShowingFirst = (currentResultsPage - 1) * resultsPerPage + 1;
	const resultShowingLast = Math.min(
		currentResultsPage * resultsPerPage,
		resultsCount
	);

	// the total number of pages
	const pagesCount = Math.ceil(resultsCount / resultsPerPage)

	// determine, based on the current page, if the prev / next pages are valid
	const prevPageValid = (currentResultsPage > 1);
	const nextPageValid = (currentResultsPage < pagesCount);

	// if the prev / next page are valid, store what page they should correspond to
	let prevPageUrl = getUpdatedPageUrl(url, currentResultsPage, prevPageValid, true);
	let nextPageUrl = getUpdatedPageUrl(url, currentResultsPage, nextPageValid, false);

	// get the array that will be displayed as the pages are select
	const pagesDisplay = getPaginationPages(currentResultsPage, pagesCount);

	// construct the pagination object
	const pagination: Pagination = {
		currentResultsPage: currentResultsPage,
		resultsPerPage: resultsPerPage,
		resultsToFetchStart: resultsToFetchStart,
		resultsToFetchEnd: resultsToFetchEnd,
		resultShowingFirst: resultShowingFirst,
		resultShowingLast: resultShowingLast,
		resultsCount: resultsCount,
		pagesDisplay: pagesDisplay,
		pagesCount: pagesCount,
		prevPageValid: prevPageValid,
		prevPageUrl: prevPageUrl,
		nextPageValid: nextPageValid,
		nextPageUrl: nextPageUrl
	  };
	
	return pagination;
};

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
