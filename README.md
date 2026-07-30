> [!WARNING]  
> Pretty much on beta. not even for a v0.0.1 and reiterating over and over with api changes

#### todo: 
- [x] sqlite
- [x] make different list delegates (simple and check)
- [ ] make a streak popup. add a sql table called "monthly_progress" to store it.
- [x] on selectOrCreate query dont increment id if it failed silently
- [x] make the flags on main more easily manageable
- [ ] make height be same as the list height IF list.height() < MAX_HEIGHT
- [x] stop using bubbles/list and make my own components
- [x] add a secondary list for long term tasks. ex: 

| Daily(NOW monthly)  | Long term |
| --------------- | --------------- |
| □ Go to gym | □ Finish X project |
| □ Read 10 pages | □ Get a j*b |

- [x] show tasks as horizontal
- [x] parse the tasks from the toml to lowercase (make them case insensitive) (+ trimmed space)
- [ ] create the config.toml file with a bash script and ln -s to $HOME/desktop inside the Makefile @install flag
- [ ] improve toml error to be more explicit
- [x] restart tasks state per month
- [x] show non completed tasks first
- [ ] optimze db (create indexes)
- [ ] make estimates on how many times needs to be done per week (└─>)
- [x] allow both ONCE or repetitive tasks. ex:
    - \["shower", "lose 10kg"\]
    - \[{name = "walk 10km", times = 10, metric = "km"},
        {name = "live 31 days", times = 31, metric = "day"}\]
    - \["bath", {name = "saludate :D", times = 10}\]
- [ ] make possible to add custom priority/metatags + colors
- [ ] make a notification of which tasks were pending last month
- [ ] make option to disable alerts about unrecognizable fields

## danzmen
more like a tui agenda using a .toml config file.

the idea is:
- when you open a terminal instead of opening your `riced fastfest with a pokemon`, or whatever bloat you use, it reminds you what monthly/longterm tasks you need to do.
- the better solution would be using a list.todo.md with [checkmate.nvim](https://github.com/bngarren/checkmate.nvim) (not sponsored btw) but whatever fits your bloat
- yes, im using this cuz pen and paper is too hard

```toml

[longterm]
tasks = [
  {ends = "30/12/2027", name = "Get marry idk", priority = "low"} # date format is DD/MM/YYYY
  {ends = "500d", name = "Get one follower in github lol", priority = "high"} # you can also do this but it gets registered once you run the program. its also kinda buggy
]

[month.july]
tasks = [
    {name = "DONT bath"}, #one time only!
    {name = "watch 300 hours of vtubers", times = 300 } # keep track of your progress being 300 the objective!
]

# every month it gets refreshed
[month.every]
tasks = [ 
        "talk to a wom*n" # you can combine both!!! same as {name = ""}
        {name = "farm *any gacha currency*", times = 1337.55} 
        {name = "consume 120 petabytes of reels", times = 120, metric = "PB"} # and specify the metric of what we're talking about! (it just to be reminded of, provides NO functionality)
    ]

```
