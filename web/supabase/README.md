<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `supabase` directory is for managing the Supabase project configuration, database schema changes, and testing. It includes a `.gitignore` file, `config.toml` for project settings, an `extensions.sql` file for creating a SQL extension, a `migrations` directory for managing database schema changes, and a `tests` directory for testing the database functionality.

### Files
#### .gitignore
This file contains a list of files and directories that should be ignored by Git when tracking changes in the `web/supabase` directory. The `.gitignore` file specifically excludes the `.branches` and `.temp` directories from version control.

#### config.toml
This file (`config.toml`) contains the configuration settings for the Supabase project, including the `project_id`, `port` numbers for the API, database, and Supabase Studio, as well as settings for email testing, authentication, and more. It also includes settings for the Supabase authentication service, such as allowing/disallowing new user signups via email, requiring confirmation for email changes, and using external OAuth providers.

#### extensions.sql
This file contains a SQL script that creates an extension called "moddatetime" in the "extensions" schema, if it does not already exist. This extension likely provides functionality related to modifying date and time values.

### Directories
#### migrations
The `migrations` directory is for managing changes to the Supabase database schema. It contains SQL scripts that create, alter, or drop tables, columns, constraints, policies, triggers, and indexes. These changes include adding new columns, tables, and views, setting default values, enabling row level security, and granting permissions to different roles. The directory also contains scripts that enable Supabase Realtime and create functions for inserting or updating data. The directory includes specific files that create or modify tables, views, and policies related to candidate and recruiter information, job interests, and job postings.

#### tests
The `tests` directory is for testing the functionality and correctness of the Supabase database. It includes tests for the infrastructure, helper functions, and specific tables such as `candidate_company_inbound`. The tests ensure the correct creation, updating, and functioning of the tables and their triggers.

<!--- END SELFDOC --->