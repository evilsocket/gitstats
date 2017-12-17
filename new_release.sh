#!/bin/bash
# nothing to see here, just a utility i use to create new releases ^_^

CURRENT_VERSION=$(cat version.go | grep VERSION | cut -d '"' -f 2)

echo -n "Current version is $CURRENT_VERSION, select new version: "
read NEW_VERSION
echo "Creating version $NEW_VERSION ..."

echo "Updating version.go"
sed -i "s/$CURRENT_VERSION/$NEW_VERSION/g" version.go

git add version.go
git commit -m "Releasing v$NEW_VERSION"
git push

git tag -a v$NEW_VERSION -m "Release v$NEW_VERSION"
git push origin v$NEW_VERSION

rm -rf dist

echo "All done, just run goreleaser now ^_^"
