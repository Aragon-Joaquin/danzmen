#### todo: 

- [ ] testing @priority(high)
    - [ ] db
    - [ ] toml parsing
    - [ ] flags
- [ ] separate `danzmen toggle --monthly 1` and `danzmen toggle --long 1`, `danzmen toggle 1` appends the --monthly flag automatically @priority(high)
- [ ] customize the float precision of the monthly_times_done @priority(medium)
- [ ] make possible to add custom priority/metatags + colors @priority(medium)
- [ ] make a notification of which tasks were pending last month @priority(medium)
- [ ] make a streak popup. add a sql table called "monthly_progress" to store it. [x] on selectOrCreate query dont increment id if it failed silently @priority(low)
- [ ] optimze db (create indexes) @priority(low)
- [ ] make estimates on how many times needs to be done per week (└─>) @priority(low)
- [ ] make option to disable alerts about unrecognizable fields @priority(low)
- [ ] finish "check" mode with: @priority(low)
    - [ ] paginator
    - [ ] keymap help

## Archive

- [x] create the config.toml file with a bash script and ln -s to $HOME/desktop inside the Makefile @install flag
- [?] orm??? or use a pattern to fix this mess lol @shelved(Reason: "bloat)
- [x] sqlite
- [x] make different list delegates (simple and check)
- [x] make the flags on main more easily manageable
- [x] make height be same as the list height IF list.height() < MAX_HEIGHT
- [x] stop using bubbles/list and make my own components
- [x] add a secondary list for long term tasks. ex: 
    | Daily(NOW monthly)  | Long term |
    | --------------- | --------------- |
    | □ Go to gym | □ Finish X project |
    | □ Read 10 pages | □ Get a j*b |

- [x] show tasks as horizontal
- [x] parse the tasks from the toml to lowercase (make them case insensitive) (+ trimmed space)
- [x] improve toml error to be more explicit
- [x] restart tasks state per month
- [x] show non completed tasks first
- [x] allow both ONCE or repetitive tasks. ex:
    - \["shower", "lose 10kg"\]
    - \[{name = "walk 10km", times = 10, metric = "km"},
        {name = "live 31 days", times = 31, metric = "day"}\]
    - \["bath", {name = "saludate :D", times = 10}\]
