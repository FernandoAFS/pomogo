#!/bin/sh

# Very primitive but useful pomogo dmenu script

dmenu=dmenu

action=$(echo "status\nplay\npause\nskip\nstop" | $dmenu) 
status=$(pomogo client $action)

if [ $? -ne 0 ]; then
    notify-send "⚠ Error" "Pomogo server not available"
    exit 1
fi

getProp(){
    echo $status | jq -r ".$1"
}

case $(getProp State) in
    "Work") 
        $()
        notify-send "👷 Working" "$(getProp TimeLeft) left"
        ;;
    "ShortBreak") 
        notify-send "⏲  Short break" "Back in $(getProp TimeLeft)"
        ;;
    "LongBreak") 
        notify-send "🍵 Long break" "Enjoy for $(getProp TimeLeft)"
        ;;
    "Paused") 
        notify-send "⏸ Paused" "Play to resume session"
        ;;
    "Stopped") 
        notify-send "⏹ Stopped" "Play to start a new session"
        ;;
    *)
        >&2 echo "Unrecognized $(getProp State) status..."
        exit 1
        ;;
esac
