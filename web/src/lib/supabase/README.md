<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `supabase` directory is for TypeScript interfaces, enums, and composite types used in the Supabase library, including the structure of tables, views, and functions in the database, as well as the shape of arguments and return values for various functions. It also contains files that export functions related to email settings and refreshing Google access tokens.

### Files
#### client.server.ts
This file exports a function `getRefreshedGoogleAccessToken` that takes a Supabase client and an email as arguments and returns a Promise that resolves to a refreshed Google access token. The function retrieves the Google refresh token from the Supabase database, refreshes the access token using the refresh token, and returns the new access token.

#### client.ts
This file exports a `UserEmailSettings` type and a `JobInterest` enum. The `UserEmailSettings` type contains properties related to email settings for a user, while the `JobInterest` enum contains options for expressing interest in a job.

#### types.ts
This file contains TypeScript interfaces, enums, and composite types used in the Supabase library, including the structure of tables, views, and functions in the database, as well as the shape of arguments and return values for various functions.

<!--- END SELFDOC --->