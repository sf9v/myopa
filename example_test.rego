package example

page_managers = [
	{
		"user_id": "admin-1234",
		"page_id": "page-1234",
		"role": "admin",
	},
	{
		"user_id": "editor-1234",
		"page_id": "page-1234",
		"role": "editor",
	},
]

test_is_anon {
	is_anon with input as {"user": "anon"}
}

test_not_is_anon {
	not is_anon with input as {"user": "user-1234"}
}

test_allow_page_admin_to_read_page {
	allow with input as {"action": "read", "object": {"type": "page", "id": "page-1234"}, "user": "admin-1234"}
		 with data.page_managers as page_managers
}

test_allow_page_admin_to_update_page {
	allow with input as {"action": "update", "object": {"type": "page", "id": "page-1234"}, "user": "admin-1234"}
		 with data.page_managers as page_managers
}

test_allow_page_admin_to_delete_page {
	allow with input as {"action": "delete", "object": {"type": "page", "id": "page-1234"}, "user": "admin-1234"}
		 with data.page_managers as page_managers
}

test_allow_page_editor_to_read_page {
	allow with input as {"action": "read", "object": {"type": "page", "id": "page-1234"}, "user": "editor-1234"}
		 with data.page_managers as page_managers
}

test_allow_page_editor_to_update_page {
	allow with input as {"action": "update", "object": {"type": "page", "id": "page-1234"}, "user": "editor-1234"}
		 with data.page_managers as page_managers
}

test_not_allow_page_editor_to_delete_page {
	not allow with input as {"action": "delete", "object": {"type": "page", "id": "page-1234"}, "user": "editor-1234"}
		 with data.page_managers as page_managers
}
