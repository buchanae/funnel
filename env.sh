#summary
fun task list | jq '[.[] | select(.name | contains("LIB17"))] | group_by(.name) | .[] | "\(.[0].name) \([.[].state] | sort | join(", "))"'

# better summary
fun task list | jq '[.[] | select(.name | contains("LIB17"))] | group_by(.name) | .[] | {"name": .[0].name, "states": [.[].state] | sort}'


# not done
fun task list | jq '[.[] | select(.name | contains("LIB17"))] | group_by(.name) | .[] | {"name": .[0].name, "states": [.[].state] | sort} | select(.states | contains(["COMPLETE"]) | not)'

alias fun='funnel -S "http://104.154.16.61/"'



curl -s http://35.193.33.100/v1/nodes | jq '[.nodes[] | select(.task_ids | length == 0) | select(.id | contains("east") | not) | .id][:50] | join(",")' -r
del=$(curl -s http://35.193.33.100/v1/nodes | jq '[.nodes[] | select(.task_ids | length == 0) | select(.id | contains("east") | not) | .id] | join(",")' -r )

fun task list | jq '[.[] | select(.name | contains("LIB17"))] | .[].state' -r | sort | uniq -c
