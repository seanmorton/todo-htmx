# todo-htmx
A dead-simple TODO app designed for personal use. I used to use Asana but it got too cluttered for me.

Built with:
* Go+HTMX to make it nice and snappy.
* SQLite for its standalone simplicity.
* [Simple.css](https://simplecss.org/) since I'm awful at CSS and would rather avoid it altogether.

### Demo

### TODO
* Tasks
    * Search
* Users/Assignees
    * CRUD
    * profile pictures for assignee bubbles on tasks
    * login?
    * default new task assignee to current user
    * private projects
* "Pinned" tasks
    * stays at top, no due date
* Design
    * support more hx-push-url for consistent reload experience
    * switch to https://gitlab.com/cznic/sqlite to avoid cgo
* Config
    * user timezone for capturing due date
