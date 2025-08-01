resource "ctfdcm_challenge_dynamiciac" "chall" {
  name        = "Some challenge"
  category    = "cat"
  description = "..."
  value       = 500
  decay       = 20
  minimum     = 50
  state       = "visible"

  shared   = true
  scenario = "localhost:5000/some/scenario:v0.1.0"
}

resource "ctfd_user" "pandatix" {
  name     = "PandatiX"
  email    = "lucastesson@protonmail.com"
  password = "password"
}

resource "ctfd_team" "ctfer" {
  name     = "CTFer.io"
  email    = "ctfer-io@protonmail.com"
  password = "ctfer"
  members = [
    ctfd_user.pandatix.id,
  ]
  captain = ctfd_user.pandatix.id
}

resource "ctfdcm_instance" "ist" {
  challenge_id = ctfdcm_challenge_dynamiciac.chall.id
  source_id    = ctfd_team.ctfer.id
}
