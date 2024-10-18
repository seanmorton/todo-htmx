# todo-htmx
Yet another personal todo app written with go templ+htmx. I used to use Asana for personal stuff, but it's UI has gotten too cluttered and I wanted to learn more htmx.

## TODO
* Tasks
    * Search
* Users/Assignees
    * CRUD
    * profile pictures for assignee bubbles on tasks
    * login?
    * default new task assignee to current user
* Projects
    * CRUD
* Improve layout/styling
    * Themes (add https://botoxparty.github.io/XP.css/)
* "Pinned" tasks
    * stays at top, no due date
* Design
    * support more hx-push-url for consistent reload experience
    * switch to https://gitlab.com/cznic/sqlite to avoid cgo
* Config
    * user timezone for capturing due date
