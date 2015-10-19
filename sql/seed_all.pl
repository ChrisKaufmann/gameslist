#!/usr/bin/perl
use warnings;
use strict;
my $db=shift || usage();

my @files;
@files=`ls -1 seed*.sql`;
chomp(@files);
print "seed_consoles\n";
system("mysql -uroot $db < seed_consoles.sql");
foreach my $file(@files)
{
	next if $file eq "seed_consoles.sql";
	print "$file\n";
	my $cmd="mysql -uroot $db < $file";
	system("$cmd")
}
print("Done.\n");

sub usage { print("usage $0 <db name>\n");exit }

