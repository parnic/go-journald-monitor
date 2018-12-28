# go-journald-monitor

Runs journalctl and prints any lines that have appeared since the last time the program was run to stdout. The log cursor is saved in a file named lastCursor (or lastCursor-unit if run for a specific unit). If the program is being run for the first time (that is, with no known last cursor position), no output is printed and the cursor position is saved.

## Arguments

-unit [unit name]
This filters output to one specific systemd service.

## Restrictions

Since the program simply calls the 'journalctl' process, 'journalctl' must be installed and usable (in the path) on the host.

## Usage

Set this in a cron job to be automatically notified via email when a service prints any log output.
