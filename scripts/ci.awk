BEGIN {
  fail=0;
}
{
  if ($1 != "ok" && $1 != "?") {
    fail++;
  }
}
END {
  if (fail == 0) {
    print "build passing";
  } else {
    print "build failure"
  }
}
