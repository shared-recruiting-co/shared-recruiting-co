libs=("./libs/src")
pkgs=(
"./cloudfunctions/candidate_email_sync"
"./cloudfunctions/gmail_subscription"
"./cloudfunctions/adhoc"
"./cloudfunctions/populate_jobs"
"./cloudfunctions/candidate_gmail_push_notifications"
"./cloudfunctions/recruiter_gmail_push_notifications"
"./cloudfunctions/candidate_gmail_messages"
"./cloudfunctions/recruiter_gmail_messages"
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
