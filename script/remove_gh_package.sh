#!/bin/bash

# Really crappy script that works to delete private packages stored on Github Packages
# Intended to simplify deleting packages that are counting against your limit
# By Troy Fontaine (github.com/troyfontaine)
# First displays the private packages name then the version finally the ID needed to delete it
# Then prompts you if you want to delete the packages based on the ID

GITHUB_TOKEN=$SUPERSECRETTOKEN
REPO_OWNER=$YOURGITHUBUSERNAME

graphqlJson() {
  local query="$1"; shift

  curl -s -H "Authorization: bearer $GITHUB_TOKEN" -X POST -H "Accept: application/vnd.github.v3+json" -d '{"query":"'"$query"'"}' 'https://api.github.com/graphql'
}

graphqlDelete() {
    local query="$1"; shift

    curl -s -H "Accept: application/vnd.github.package-deletes-preview+json" -H "Authorization: bearer $GITHUB_TOKEN" -X POST -d '{"query":"'"$query"'"}' 'https://api.github.com/graphql'
}

deletePackageID() {
    PACKAGE_ID="$1"
    local query="$(cat <<EOF | sed 's/"/\\"/g' | tr '\n\r' '  '
mutation {
    deletePackageVersion(
      input:{packageVersionId:"$PACKAGE_ID"}
    )
    { success }
}

EOF
)"

  RESPONSE=$(graphqlDelete "$query")
  echo "$RESPONSE"
}

listPackageIDs() {

  local query="$(cat <<EOF | sed 's/"/\\"/g' | tr '\n\r' '  '
query {
    user(login:"$REPO_OWNER") {
        registryPackagesForQuery(first: 10, query:"is:private") {
            totalCount nodes {
                nameWithOwner versions(first: 10) {
                  nodes {
                    id version
                  }
                }
                }
            }
        }
    }

EOF
)"

  PACKAGE_LIST=$(graphqlJson "$query")
  echo -e "Package Name\t\t\t\t\tVersion\t\t\t\t\tPackage ID"
  echo $PACKAGE_LIST | jq -r '.data.user.registryPackagesForQuery.nodes[] | "\(.nameWithOwner)\t\t\t\(.versions.nodes[].version)\t\t\t\(.versions.nodes[].id)"'
  ID_LIST=$(echo $PACKAGE_LIST | jq -r '.data.user.registryPackagesForQuery.nodes[].versions.nodes[].id')
}

purgePackage() {
  for ID in $ID_LIST
  do
    echo -e "Purge package with ID: '$ID'?"
    select yn in "Yes" "No"; do
        case $yn in
            Yes ) deletePackageID $ID; break;;
            No ) exit;;
        esac
    done
  done
}

listPackageIDs
purgePackage
