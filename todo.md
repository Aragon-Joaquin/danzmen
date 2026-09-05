#### todo: 
- [ ] On 8+ items, show remaining ones and implement `danzmen list 2` to show the second page. @priority(crucial)
- [ ] IMPROVE ARCHITECTURE: @priority(high)
    - [ ] implement flags as their own struct in a Value any/interface so then i can do type casting
    - [ ] improve the db queries holy, its a mess
    - [x] render.go
- [ ] testing @priority(high) render.go
    - [ ] db
    - [ ] toml parsing
    - [ ] flags
- [ ] separate `danzmen toggle --monthly 1` and `danzmen toggle --long 1`, `danzmen toggle 1` appends the --monthly flag automatically @priority(high)
- [ ] implement `danzmen add --monthly 1 +10`, `danzmen add --long 1 -4` @priority(high)
- [ ] when `danzmen add` or `danzmen toggle`. check if the task is in todays month, if not, dont execute the query and throw an error @priority(high)|
- [ ] cache the results to avoid making sql calls if the tasks dont change @priority(medium) 
- [ ] implement `danzmen reset` to reset the db. @priority(medium) 
- [ ] customize the float precision of the monthly_times_done @priority(medium)
- [ ] make possible to add custom priority/metatags + colors @priority(medium)
- [ ] make a notification of which tasks were pending last month @priority(medium)
- [ ] make a size between SIZE_MEDIUM and SIZE_BIG and reorder the items to occupy more space @priority(medium)
- [ ] make a streak popup. add a sql table called "monthly_progress" to store it. [x] on selectOrCreate query dont increment id if it failed silently @priority(low)
- [ ] optimze db (create indexes) @priority(low)
- [ ] make estimates on how many times needs to be done per week (+ make it an option) (└─>) @priority(low)
- [ ] make option to disable alerts about unrecognizable fields @priority(low)
- [ ] finish "check" mode with: @priority(low)
    - [ ] paginator
    - [ ] keymap help
- [ ] make a `danzmen help command` @priority(low)

## Archive

- [x] reset task counter on a new month @priority(high)
- [x] fields that are not readed (like [august.]) are logged as not valid fields @priority(high) 
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


### Test Plan
i asked chudgpt on what to add, i'll be manually adding these... sadly.

##### db (db/)
Init() hardcodes a home-dir path — tests bypass it with an in-memory helper (sqlite3 :memory:).

- TestCreateDatabase — Schema created; all 5 tables exist; re-running is idempotent.
- TestInsertOrSelectYear_MonthID — Same month/year returns same id; different returns new.
- TestCreateIfNotExistsMonthlyTasks — Inserts tasks + records, dedupes by name, returns full join with cfg.
- TestUpdateCompletedMonthlyTask — Toggles completed 0↔1; only affects matching year_month + id.
- TestUpdateCompletedLongTask — Same toggle behavior for long tasks path.
- TestInsertOrSelectLongTermTasks — Inserts/dedupes long tasks, returns scan data + cfg; errors on empty input.

##### config / TOML (config/)
ParseTOML reads a fixed user path — tests stub the file. Note getNonRepeatableMonthlyTasks / getNonRepetableLongTermTasks use log.Fatalln, so invalid-input tests need the Fatalf-intercept pattern (or refactor to return errors).

- TestParseTOML_FileMissing — Returns default empty cfg, no error.
- TestParseTOML_Valid — Parses month/longterm, fills parsed slices.
- TestParseTOML_InvalidTOML — Malformed file returns error.
- TestGetNonRepeatableMonthlyTasks — Only current month + every included; dedupes names; handles string vs table form.
- TestGetNonRepetableLongTermTasks — Dedupes, validates expiry notation + priority.

##### flags (flags/)
ParseOptions reads global os.Args and printHelp writes to stdout — tests swap os.Args / capture output.

- TestParseOptions_NoArgs — Returns PROGRAM_HELP.
- TestParseOptions_HelpListCheck — Each subcommand maps to correct PROGRAM_* type.
- TestParseOptions_Toggle — toggle {id} sets Args + Value; missing id → error; -long sets Value.
- TestParseOptions_Unknown — Unknown first arg → HELP.
- TestPrintHelp — Returns PROGRAM_HELP.
- TestFlagToggle — Needs db; toggles matching id, "Id not found" path, invalid id.

##### Pre-existing blockers to fix first
- db.UpdateCompletedTask
- flags/flags.go:45 — target_long declared but unused (build fails).
- FlagToggle uses os.Exit — refactor to return errors for testability.
