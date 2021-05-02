package example

default allow = false

# objects: page, post
# page roles: anon, admin, editor
# actions: create, read, update, delete

page_role_grants = {
	"admin": {
		"page": ["read", "update", "delete"],
		"post": ["create", "read", "update", "delete"],
	},
	"editor": {
		"page": ["read", "update"],
		"post": ["create", "read", "update", "delete"],
	},
}

# is role granted to perform action on object
is_page_role_granted(role, action, object) {
	page_role_grants[role][object][_] == action
}

# is user anonymous
is_anon {
	input.user == "anon"
}

# allow users to create page
allow {
	not is_anon

	input.object.type == "page"
	input.action == "create"
}

# allow all users to read page
allow {
	input.object.type == "page"
	input.action == "read"
}

# allow users to perform the the allowed action on their role
allow {
	not is_anon

	page_manager := data.page_managers[_]
	page_manager.page_id == input.object.id
	page_manager.user_id == input.user
	is_page_role_granted(page_manager.role, input.action, input.object.type)
}
