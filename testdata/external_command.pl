#!/usr/bin/env perl
#
# external_command.pl -- toy program for testing external tools.
# -------------------
# Echoes args, one per line, with output modified by flags as follows:
#
# --seed=S   - seed ID with this the real number S
# --header=H - print H before echoing (can specify multiple headers)
# --indent=N - indent each line N spaces
# --prefix=P - print P after indent and before args
# --stdin    - echo standard input after args
# --stderr   - echo to standard error instead of standard output
#
# Headers always go to standard output. The first line is always the ID.
#
# Short flags can be used.
#
# The --stdin flag should be used as a configured PreArg.

use strict;
use warnings;

use feature 'say';
use Getopt::Long;
use Digest::MD5 qw(md5_hex);

my $seed   = 0;
my $indent = 0;
my $prefix = "";
my @headers;
my $stderr;
my $stdin;

MAIN: {

    GetOptions(
        "seed=f"   => \$seed,       # numeric (float)
        "indent=i" => \$indent,     # numeric (integer)
        "prefix=s" => \$prefix,     # string
        "header=s" => \@headers,    # array
        "stderr"   => \$stderr,     # flag
        "stdin"    => \$stdin,      # flag
    ) or die("Error in command line arguments\n");

    srand( int($seed) ) if $seed;
    say md5_hex( rand $seed );
    say $_ for @headers;
    say_what($_) for @ARGV;

    if ($stdin) {
        while ( my $line = <STDIN> ) {
            say_what($line);
        }
    }
}

sub say_what {
    my $what = shift @_;
    chomp $what;
    if ($stderr) {
        say STDERR " " x $indent, $prefix, $what;
    }
    else {
        say STDOUT " " x $indent, $prefix, $what;
    }

}

