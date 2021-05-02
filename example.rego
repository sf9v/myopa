package example.authz

# actions: create, read, update, delete
# types: page, product
# page manager roles: admin, editor

default allow = false

# any non-anonymous user can create a page
allow {
	input.action == "create"
	input.object.type == "page"
	input.user != "anon"
}

# any user can a read page
allow {
	input.action == "read"
	input.object.type == "page"
}

# only a page manager can update page
allow {
	input.action == "update"
	input.object.type == "page"

	page := data.pages[_]
	page.id == input.object.id

	page_manager := data.page_managers[_]
	page_manager.page_id == page.id
	page_manager.user_id == input.user
}

# only page admins can delete a page
allow {
    input.action == "delete"
    input.object.type == "page"

    page := data.pages[_]
    page.id == input.object.id

    page_manager := data.page_managers[_]
    page_manager.page_id == page.id
    page_manager.user_id == input.user
    page_manager.role == "admin"
}
