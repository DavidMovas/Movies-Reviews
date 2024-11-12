data "moviereviews_user" "john" {
  username = "john"
}

resource "moviereviews_user_role" "john" {
  user_id = data.moviereviews_user.john.id
  role    = "editor"
}