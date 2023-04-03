libs=("./libs/src")
pkgs=(
"./cloudfunctions/gmail_subscription"
"./cloudfunctions/adhoc"
"./cloudfunctions/candidate_email_sync"
"./cloudfunctions/candidate_gmail_push_notifications"
"./cloudfunctions/candidate_gmail_messages"
"./cloudfunctions/recruiter_email_sync"
"./cloudfunctions/recruiter_gmail_push_notifications"
"./cloudfunctions/recruiter_gmail_messages"
"./cloudfunctions/scrape_job_listings"
)
sha=$(git rev-parse origin/main)

for lib in "${libs[@]}"; do
  pushd $lib
  echo "Updating $(basename $lib)"
  go mod tidy
  popd
done

for pkg in "${pkgs[@]}"
do
  pushd $pkg
  echo "Updating $(basename $pkg)"
  go get -u github.com/shared-recruiting-co/shared-recruiting-co/libs/src@${sha}
  go mod tidy
  popd
done
