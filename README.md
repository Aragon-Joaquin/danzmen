> [!NOTE]  
> Pretty much on beta. not even for a v0.0.1 and reiterating over and over with api changes

#### todo: 
- [x] sqlite
- [x] make different list delegates (simple and check)
- [ ] make a streak popup. add a sql table called "date_objectives" to store it.
- [ ] add support for date specific tasks (however it isn't the purpose of this)
- [x] on selectOrCreate query dont increment id if it failed silently
- [x] make the flags on main more easily manageable
- [ ] make height be same as the list height IF list.height() < MAX_HEIGHT
- [ ] stop using bubbles/list and make my own components
- [ ] add a secondary list for long term tasks. ex: 

| Daily  | Long term |
| --------------- | --------------- |
| [ ] Go to gym | [ ] Finish X project |
| [ ] Read 10 pages | [ ] Get a j*b |

- [ ] finish loop.*d tasks
- [ ] show tasks as horizontal
- [x] parse the tasks from the toml to lowercase (make them case insensitive) (+ trimmed space)

## danzmen
more like a tui daily agenda using a .toml config file.

the idea is:
- when you open a terminal instead of opening your `riced fastfest with a pokemon`, or whatever bloat you use, it reminds you what daily tasks you need to do.
- after few days you probably wont even care or remove this, but at that point i expect you to KNOW your routine, making this program accomplish its mission.
- yes, im using this cuz pen and paper is too hard

```toml
#week starts on sunday
start = "sunday"

# each Monday
[day.monday]
tasks = [ "run 500km", "read the man pages" ]

# starting from Monday, every 2 days until Friday
[loop.2d]
tasks = ["goto gym"] # Monday, Wednesday, Friday
start = "monday"
end = "friday"

[loop.1w] # 1 week
tasks = ["shower"] # every week, starting from Sunday


#ERROR: cannot loop if the loop is greater than the difference between days/weeks
# from Monday to Tuesday there's a 1 day and its declare to repeat each 4 days
[loop.4d]
tasks = ["this is an invalid format"]
start = "monday"
end = "tuesday" 

# everyday!
[loop.1d] 
tasks = ["live"]
```

for now, the .config file is subject to changes since im horrible at naming things.

