create policy "Users can delete their own jobs"
on "public"."user_email_job"
as permissive
for delete
to public
using ((auth.uid() = user_id));
