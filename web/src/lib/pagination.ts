const DEFAULT_PAGE_SIZE = 10;

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
	let prevPageUrl = prevPageValid ? new URL(url.href) : null;
	if (prevPageValid) {
		prevPageUrl.searchParams.set('page', currentResultsPage - 1);
		prevPageUrl = prevPageUrl.href
	}

	let nextPageUrl = nextPageValid ? new URL(url.href) : null;
	if (nextPageValid) {
		nextPageUrl.searchParams.set('page', currentResultsPage + 1);
		nextPageUrl = nextPageUrl.href
	}


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
