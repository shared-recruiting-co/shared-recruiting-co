export type UserEmailSettings = {
	is_active: boolean;
	auto_archive?: boolean;
	auto_contribute?: boolean;
};

export enum JobInterest {
	Interested = 'interested',
	NotInterested = 'not_interested',
	Saved = 'saved'
}
