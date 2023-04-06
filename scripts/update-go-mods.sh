libs=("./libs/src")
pkgs=(
"./cloudfunctions/gmail_subscription"
"./cloudfunctions/adhoc"
"./cloudfunctions/candidate_email_sync"
"./cloudfunctions/candidate_gmail_push_notifications"
"./cloudfunctions/candidate_gmail_messages"
"./cloudfunctions/candidate_gmail_label_changes"
"./cloudfunctions/recruiter_email_sync"
"./cloudfunctions/recruiter_gmail_push_notifications"
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
  echo "Testing $(basename $pkg)"
  go test
  popd
done