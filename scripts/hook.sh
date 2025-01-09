#!/bin/sh

# Script to be ran by pomogo on every event.
# The following environment variables are included:
# 
# POMO_EVENT: Reason why the hook was called:
#   - Error: on error.
#   - EndOfState: When a state time ends.
#   - Play: On successful start or resume event.
#   - Pause: On successful pause request.
#   - Stop: On successful stop.
#
# POMO_STATUS: When in error is the error message. The current status otherwise. It may be:
#   - Work
#   - ShortBreak
#   - LongBreak
#
# POMO_AT: Iso Date of the moment the event was triggered.

case $POMO_EVENT in
    "Pause") 
        notify-send "⏸ Pause" "Paused $POMOGO_STATUS"
        ;;
    "Stop") 
        notify-send "⏹ Pause" "Paused $POMOGO_STATUS"
        ;;
    "Error") 
        notify-send "⚠ Error" "$POMOGO_STATUS"
        ;;
    *) # Play, EndOfState
        notify-send "$POMOGO_STATUS"
        ;;
esac

