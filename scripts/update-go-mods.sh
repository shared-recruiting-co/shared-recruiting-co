libs=("./libs/src")
pkgs=("./cloudfunctions/email_push_notifications" "./cloudfunctions/full_email_sync" "./cloudfunctions/watch_emails" "./cloudfunctions/adhoc" "./cloudfunctions/populate_jobs")
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
