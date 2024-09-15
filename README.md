# todo-htmx
Yet another personal todo app written with go templ+htmx. I used to use Asana for personal stuff, but it's UI has gotten too cluttered and I wanted to learn more htmx.

## TODO
* Search
* Users/Assignees
    * profile pictures for assignee bubbles on tasks
    * username/pass
    * default new task assignee to current user
* Project
    * CRUD
* Improve layout/styling
    * Themes (add https://botoxparty.github.io/XP.css/)
    * Show due date in red if past due
* "Pinned" tasks
    * stays at top, no due date
* Design
    * support more hx-push-url for consistent reload experience
    * switch to https://gitlab.com/cznic/sqlite to avoid cgo
* Config
    * user timezone for displaying time
