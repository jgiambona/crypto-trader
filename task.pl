#!/usr/bin/perl -w

use 5.018;
use strict;
use warnings;

use Getopt::Long qw(GetOptions);
use Pod::Usage qw(pod2usage);

my $verbose = 0;
my $man =0;
my $help = 0;

GetOptions (
  help => \$help,
  verbose => \$verbose,
  man => \$man) or pod2usage(2);
pod2usage(1) if $help;
pod2usage(-verbose => 2) if $man;


__END__

=encoding utf8

=head1 NAME

task.pl - Run curl task for trader.

=head1 SYPNOSIS

task.pl [options]

 Options:
   -help               brief help messages
   -verbose            verbose all output
   -man                full documentation

=head1 OPTIONS

=over 4

=item B<-help>

Print a brief help message and exits.

=item B<-verbose>

Be more verbose.

=item B<-man>

Prints the manual page and exits.

=back

=head1 DESCRIPTION

B<task.pl> will read arguments to run curl tasks.

For more info visit L<Happiness|http://vastorigins.net>.

=cut

