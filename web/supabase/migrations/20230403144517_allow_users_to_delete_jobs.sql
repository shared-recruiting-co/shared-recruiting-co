create policy "Enable delete for users based on user_id"
on "public"."user_email_job"
as permissive
for delete
to public
using ((auth.uid() = user_id));
