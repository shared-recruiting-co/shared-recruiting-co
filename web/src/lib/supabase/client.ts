export enum UserEmailStats {
  EmailsProcessed = "emails_processed",
  JobsDetected = "jobs_detected",
}

export type UserEmailSettings = {
  is_active: boolean;
  auto_archive?: boolean;
  auto_contribute?: boolean;
};
