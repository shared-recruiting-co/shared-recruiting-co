pkgs=("./cloudfunctions/email_push_notifications" "./cloudfunctions/full_email_sync" "./cloudfunctions/watch_emails")
sha=$(git rev-parse origin/main)

for pkg in "${pkgs[@]}"
do
  pushd $pkg
  echo "Updating $(basename $pkg)"
  go get -u github.com/shared-recruiting-co/shared-recruiting-co/libs/gmail@${sha}
  go get -u github.com/shared-recruiting-co/shared-recruiting-co/libs/db@${sha}
  go mod tidy
  popd
done
