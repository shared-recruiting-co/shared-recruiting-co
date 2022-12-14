pkgs=("./cloudfunctions/collect_examples" "./cloudfunctions/email_push_notifications" "./cloudfunctions/full_email_sync" "./cloudfunctions/watch_emails")

for pkg in "${pkgs[@]}"
do
  pushd $pkg
  echo "Updating $(basename $pkg)"
  go get github.com/shared-recruiting-co/shared-recruiting-co/libs/gmail@main
  go get github.com/shared-recruiting-co/shared-recruiting-co/libs/db@main
  go mod tidy
  popd
done
