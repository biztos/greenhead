#!/usr/bin/env perl
#
# licenses.pl -- print licenses for documentation.
#
# (quite sloppy right now!)
#
# This should be run after adding dependencies, or as part of CI.

use strict;
use warnings;

use JSON::PP qw(decode_json);

my $mod_in = `go mod edit -json`;
my $d      = decode_json($mod_in) or die $!;

my @req = grep { !$_->{Indirect} } @{ $d->{Require} };

my %have    = map { $_->{Path} => 1 } @req;
my $license = {};
my $link    = {};

my $lic_in = `go-licenses report ./...`;

for my $line ( split /\n/, $lic_in ) {
    my ( $pkg, $url, $name ) = split /,/, $line;
    next unless $have{$pkg};

    $license->{$pkg} = $name;
    $link->{$pkg}    = $url;
}

print <<'_HEAD';
---

# Included Software Packages

The following packages, with the licenses as shown, are included in this
software.

_HEAD
for my $item (@req) {
    my $pkg = $item->{Path};
    next unless $license->{$pkg};
    printf "## %s -- %s\n\n%s\n\n", $pkg, $license->{$pkg}, $link->{$pkg};
}

