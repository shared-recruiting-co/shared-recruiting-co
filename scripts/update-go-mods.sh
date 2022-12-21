libs=("./libs/gmail" "./libs/db")
pkgs=("./cloudfunctions/email_push_notifications" "./cloudfunctions/full_email_sync" "./cloudfunctions/watch_emails" "./cloudfunctions/adhoc")
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
  go get -u github.com/shared-recruiting-co/shared-recruiting-co/libs/gmail@${sha}
  go get -u github.com/shared-recruiting-co/shared-recruiting-co/libs/db@${sha}
  go mod tidy
  popd
done
