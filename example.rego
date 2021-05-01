package example.authz

# actions: create, read, update, delete
# types: page, product
# page roles: admin, editor

default allow = false

# any user can a page
allow {
	input.action == "read"
	input.object.type == "page"
}

# only a page manager can update a page
allow {
	input.action == "update"
	input.object.type == "page"

	page := data.pages[_]
	page.id == input.object.id

	page_manager := data.page_managers[_]
	page_manager.page_id == page.id
	page_manager.user_id == input.user
}

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
