<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `database` directory is for SQL tests and procedures related to the Supabase database. It includes tests for the infrastructure, helper functions, and specific tables such as `candidate_company_inbound`. The tests ensure the correct creation, updating, and functioning of the tables and their triggers.

### Files
#### 0_setup.test.sql
This file contains SQL code that creates procedures for logging in and out of a user account, as well as a procedure for logging in as an anonymous user. It also includes a test to ensure that the procedures are working properly.

#### canary.test.sql
This file contains a test called "Canary Test" which checks if the test infrastructure is working properly. It includes three tests: checking if pgtap is installed, if the Supabase database is set up correctly, and if helper functions are available.

#### candidate_company_inbound.test.sql
This file contains SQL tests for the `candidate_company_inbound` table. The tests verify the creation, updating, and correct functioning of the table's triggers. The tests include verifying the correct insertion of candidate and recruiter information, updating the `job_id` and `candidate_id` fields, and checking for correct updates when a job is unassigned from a recruiter outbound template.

#### get_user_profile_by_email.test.sql
This file tests the `get_user_profile_by_email` helper function. It inserts a user and user profile, and then checks if the user profile can be fetched correctly with and without an oauth token. It also tests if the correct user profile is fetched when there are multiple oauth tokens with different emails.

#### get_user_profile_recruiter_by_email.test.sql
This file contains SQL code to test the `get_recruiter_by_email` helper function. It inserts a user, company, and recruiter, and then checks if the recruiter can be fetched correctly using their email with and without an OAuth token.

#### list_similar_recruiter_outbound_templates.test.sql
This file contains SQL tests for the `list_similar_recruiter_outbound_templates` helper function. The tests set various parameters such as `company_id`, `user_id`, `user_email`, `template_id`, `job_id`, `first_name`, `last_name`, `template_subject`, `template_body`, and `template_normalized_content`. The tests check if the function returns the correct template ID for identical content and null for different content.

<!--- END SELFDOC --->