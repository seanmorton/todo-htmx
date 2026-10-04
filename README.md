# todo-htmx
A dead-simple TODO app designed for personal use. I used to use Asana but it got too cluttered for me.

Built with:
* Go templ + HTMX to make it nice and snappy.
* SQLite for its standalone simplicity.
* [Simple.css](https://simplecss.org/) since I'm awful at CSS and would rather avoid it altogether.

https://github.com/user-attachments/assets/b276cfa4-ee3e-4ab1-a418-3fba21490614

## A note on timezones
Given this is intended to be a self hosted app, the server's local timezone is used for all date handling.

### TODO
* Tasks
    * calendar view
    * clickable links in desc
* Tasks
    * project level default filter config
* Users/Assignees
    * CRUD
    * profile pictures for assignee bubbles on tasks
    * login?
    * default new task assignee to current user
    * private projects
* Design/Chores
    * upgrade to go 1.27
    * upgrade htmx + templ
    * isloate JSON policy marshing to db layer
    * support more hx-push-url for consistent reload experience
    * switch to https://gitlab.com/cznic/sqlite to avoid cgo
    * button loading states
    * schema migrations
