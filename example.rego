package example

# actions: create, read, update, delete
# types: page, product
# page manager roles: admin, editor

default allow = false

is_anon {
	input.user == "anon"
}

# non-anon users can create a page
allow {
	not is_anon

	input.action == "create"
	input.object.type == "page"
}

# anyown can read any page
allow {
	input.action == "read"
	input.object.type == "page"
}

# page managers can update the page
allow {
	not is_anon

	input.action == "update"
	input.object.type == "page"
	page_id := input.object.id
	user_id := input.user

	page_manager := data.page_managers[_]
	page_manager.page_id == page_id
	page_manager.user_id == user_id
}

# page admins can delete the page
allow {
	not is_anon

    input.action == "delete"
    input.object.type == "page"
	page_id := input.object.id
	user_id := input.user

    page_manager := data.page_managers[_]
    page_manager.page_id == page_id
    page_manager.user_id == user_id
    page_manager.role == "admin"
}
