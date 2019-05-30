#!/bin/sh

find() {
  where=$1
  what=$2

  cd $where
  rg $what --vimgrep | awk -F ':' '{printf "git blame --porcelain -L %s,%s %s\n",$2,$2,$1}' > blames.sh
  chmod +x blames.sh
  ./blames.sh > blames
}

find2() {
  where=$1
  what=$2
  code=$3
  resultFile=$4

  mkdir -p /tmp/blames

  # first find files that should be inspected
  # git blame is faster to do on the whole file than to run on multiple lines for the same file

  cd $where
  rg $what -l > files
  rm -f $resultFile

  while IFS= read -r file
  do
    git blame -w -e -n -l -f --date=iso8601 $file | rg $what >> $resultFile
  done < files
  rm files
  wc -l $resultFile
}

main() {
  where=$1
  what=$2
  projectCode=`date +%s`
  # projectCode=$3
  # resultFile=$4
  resultFile=$3

  # time find $where $what
  find2 $where $what $projectCode $resultFile
  # blame $
}

main $@
