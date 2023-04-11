// Label is an enum for all @SRC labels
export enum Label {
	SRC = '@SRC',
	// Jobs
	Jobs = '@SRC/Jobs',
	JobOpportunity = '@SRC/Jobs/Opportunity',
	JobInterested = '@SRC/Jobs/Interested',
	JobNotInterested = '@SRC/Jobs/Not Interested',
	JobSaved = '@SRC/Jobs/Saved'
}

// Labels maps a label to its ID
export type Labels = Record<Label, string>;
