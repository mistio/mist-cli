#! /usr/bin/env bash
VERSION=$CI_COMMIT_TAG
bins=$(curl -G -d "prefix=cli/${VERSION}/bin/" -d "delimeter=/" -s  https://content-storage.googleapis.com/storage/v1/b/mist-downloads/o | jq '.items[] | {name: .name, url: ("https://dl.mist.io/" + .name)}' | jq -s)
bins_tuples=()
for bin in $(echo "${bins}" | jq -c '.[]'); do
name=$(echo $bin | jq '.name' | tr -d '\"' | awk -F/ '{print $(NF-2)" "$(NF-1)" "$NF}')
url=$(echo $bin | jq '.url' | tr -d '\"')
if ! echo $name | grep -q -E '\.sha256$'; then
    sha256=$(curl -s "$url.sha256")
    bins_tuples+=( "${name} ${url} ${sha256}" )
fi
done
oses_order=("linux" "darwin" "windows" "freebsd" "openbsd" "netbsd")
tables=""
IFS=$'\n'
for os in ${oses_order[@]}; do
IFS=$'\n'
uppercase_os="$(tr '[:lower:]' '[:upper:]' <<< ${os:0:1})${os:1}"
table="#### ${uppercase_os}\nARCH | SHA256\n------------ | -------------\n"
for bin_tuple in ${bins_tuples[@]}; do
    case $bin_tuple in $os*)
    echo $bin_tuple
    IFS=$' '
    bin_array=($bin_tuple)
    table="${table}[${bin_array[1]}](${bin_array[3]}) | ${bin_array[4]}\n"
    esac
done
tables="${tables}${table}\n"
done
echo $tables
echo -e $tables >> release.md